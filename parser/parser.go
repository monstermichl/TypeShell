package parser

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"unicode"

	"github.com/monstermichl/typeshell/lexer"
)

var typeMapping = map[lexer.VarType]DataType{
	lexer.DATA_TYPE_BOOLEAN: DATA_TYPE_BOOLEAN,
	lexer.DATA_TYPE_INTEGER: DATA_TYPE_INTEGER,
	lexer.DATA_TYPE_STRING:  DATA_TYPE_STRING,
	lexer.DATA_TYPE_ERROR:   DATA_TYPE_STRING, // error is internally just a string to make heandling easier.
}

type scope string

const (
	SCOPE_PROGRAM  scope = "program"
	SCOPE_FUNCTION scope = "function"
	SCOPE_IF       scope = "if"
	SCOPE_FOR      scope = "for"
	SCOPE_SWITCH   scope = "switch"
)

func scopesToString(scopes []scope) []string {
	strings := make([]string, len(scopes))

	for i, scope := range scopes {
		strings[i] = string(scope)
	}
	return strings
}

type context struct {
	imports    map[string]string             // Maps import aliases to file hashes.
	variables  map[string]Variable           // Stores the variable name to variable relation.
	functions  map[string]FunctionDefinition // Stores the function name to function relation.
	scopeStack []scope                       // Stores the current scopes.
}

func (c context) currentScope() scope {
	return c.scopeStack[len(c.scopeStack)-1]
}

func (c context) global() bool {
	return c.currentScope() == SCOPE_PROGRAM
}

func (c context) findScope(s scope) bool {
	for i := len(c.scopeStack) - 1; i >= 0; i-- {
		if c.scopeStack[i] == s {
			return true
		}
	}
	return false
}

func (c context) buildPrefixedName(name string, prefix string, global bool, checkExistence bool) (string, error) {
	name = strings.TrimSpace(name)

	if len(name) == 0 {
		return "", errors.New("no name provided")
	}
	prefix = strings.TrimSpace(prefix)

	if len(prefix) > 0 && global {
		hash, exists := c.imports[prefix]

		if checkExistence && !exists {
			return "", fmt.Errorf(`prefix "%s" not found`, prefix)
		}
		name = buildPrefixedName(hash, name)
	}
	return name, nil
}

func (c context) addImport(alias string, hash string) error {
	c.imports[alias] = hash
	return nil
}

func (c context) addVariables(prefix string, global bool, variables ...Variable) error {
	for _, variable := range variables {
		prefixedName, err := c.buildPrefixedName(variable.Name(), prefix, global, false)

		if err != nil {
			return err
		}
		c.variables[prefixedName] = variable
	}
	return nil
}

func (c context) addFunctions(prefix string, global bool, functions ...FunctionDefinition) error {
	for _, function := range functions {
		prefixedName, err := c.buildPrefixedName(function.Name(), prefix, global, false)

		if err != nil {
			return err
		}
		c.functions[prefixedName] = function
	}
	return nil
}

func (c context) findImport(alias string) (string, bool) {
	hash, exists := c.imports[alias]
	return hash, exists
}

func (c context) findVariable(name string, prefix string, global bool) (Variable, bool) {
	prefixedName, err := c.buildPrefixedName(name, prefix, global, true)

	if err != nil {
		return Variable{}, false
	}
	variable, exists := c.variables[prefixedName]
	return variable, exists
}

func (c context) findFunction(name string, prefix string) (FunctionDefinition, bool) {
	prefixedName, err := c.buildPrefixedName(name, prefix, true, true)

	if err != nil {
		return FunctionDefinition{}, false
	}
	function, exists := c.functions[prefixedName]
	return function, exists
}

func (c context) clone() context {
	return context{
		imports:    maps.Clone(c.imports),
		variables:  maps.Clone(c.variables),
		functions:  maps.Clone(c.functions),
		scopeStack: slices.Clone(c.scopeStack),
	}
}

type evaluatedImport struct {
	alias string
	path  string
}

type evaluatedValues struct {
	values []Expression
}

func (ev evaluatedValues) isMultiReturnCall() (bool, Call) {
	var call Call
	values := ev.values
	multi := false

	if len(values) == 1 {
		callTemp, ok := values[0].(Call)

		if ok && len(callTemp.ReturnTypes()) > 1 {
			multi = true
			call = callTemp
		}
	}
	return multi, call
}

type blockCallback func(statements []Statement, last bool) error

type Parser struct {
	tokens    []lexer.Token
	index     int
	path      string
	prefix    string
	currFunc  string
	usedFuncs map[string][]string // Stores which function (key) calls which functions (values).
}

func New() Parser {
	return Parser{
		usedFuncs: map[string][]string{},
	}
}

func (p *Parser) Parse(path string) (Program, error) {
	return p.parse(path, false)
}

func (p *Parser) parse(path string, imported bool) (Program, error) {
	// If path is relative, make it absolute.
	if !filepath.IsAbs(path) {
		pathTemp, err := filepath.Abs(path)

		if err != nil {
			return Program{}, err
		}
		path = pathTemp
	}

	// Make sure path exists.
	if _, err := os.Stat(path); err != nil {
		return Program{}, err
	}
	source, err := os.ReadFile(path)

	if err != nil {
		return Program{}, err
	}
	tokens, err := lexer.Tokenize(string(source))

	if err != nil {
		return Program{}, err
	}
	p.index = 0
	p.tokens = tokens
	p.path = path
	p.prefix = ""

	// If it's an imported file, use source hash as prefix.
	if imported {
		h := sha256.New()
		h.Write(source)

		p.prefix = fmt.Sprintf("%x", h.Sum(nil))[0:7] // Only use the 7 first characters (inspired by Git).
	}
	program, err := p.evaluateProgram()

	if err != nil {
		return Program{}, err
	}

	// If this is the original program, removed unused stuff.
	if !imported {
		return p.cleanProgram(program)
	}
	return program, nil
}

func allowedBinaryOperators(t ValueType) []BinaryOperator {
	operators := []BinaryOperator{}

	if !t.IsSlice() {
		switch t.DataType() {
		case DATA_TYPE_INTEGER:
			operators = []BinaryOperator{BINARY_OPERATOR_MULTIPLICATION, BINARY_OPERATOR_DIVISION, BINARY_OPERATOR_MODULO, BINARY_OPERATOR_ADDITION, BINARY_OPERATOR_SUBTRACTION}
		case DATA_TYPE_STRING:
			operators = []BinaryOperator{BINARY_OPERATOR_ADDITION}
		default:
			// For other types no operations are permitted.
		}
	}
	return operators
}

func allowedCompareOperators(t ValueType) []CompareOperator {
	operators := []CompareOperator{}

	if !t.IsSlice() {
		switch t.DataType() {
		case DATA_TYPE_BOOLEAN:
			operators = []CompareOperator{COMPARE_OPERATOR_EQUAL, COMPARE_OPERATOR_NOT_EQUAL}
		case DATA_TYPE_INTEGER:
			operators = []CompareOperator{COMPARE_OPERATOR_EQUAL, COMPARE_OPERATOR_NOT_EQUAL, COMPARE_OPERATOR_LESS, COMPARE_OPERATOR_LESS_OR_EQUAL, COMPARE_OPERATOR_GREATER, COMPARE_OPERATOR_GREATER_OR_EQUAL}
		case DATA_TYPE_STRING:
			operators = []CompareOperator{COMPARE_OPERATOR_EQUAL, COMPARE_OPERATOR_NOT_EQUAL, COMPARE_OPERATOR_LESS, COMPARE_OPERATOR_LESS_OR_EQUAL, COMPARE_OPERATOR_GREATER, COMPARE_OPERATOR_GREATER_OR_EQUAL}
		default:
			// For other types no operations are permitted.
		}
	}
	return operators
}

func defaultVarValue(valueType ValueType) (Expression, error) {
	dataType := valueType.DataType()

	if !valueType.IsSlice() {
		switch dataType {
		case DATA_TYPE_BOOLEAN:
			return BooleanLiteral{}, nil
		case DATA_TYPE_INTEGER:
			return IntegerLiteral{}, nil
		case DATA_TYPE_STRING:
			return StringLiteral{}, nil
		}
	} else {
		return SliceInstantiation{dataType: dataType}, nil
	}
	return nil, fmt.Errorf("no default value found for type %s", valueType.String())
}

func incrementDecrementStatement(variable Variable, increment bool) Statement {
	operation := BINARY_OPERATOR_ADDITION

	if !increment {
		operation = BINARY_OPERATOR_SUBTRACTION
	}
	return VariableAssignment{
		variables: []Variable{variable},
		values: []Expression{
			BinaryOperation{
				left: VariableEvaluation{
					Variable: variable,
				},
				operator: operation,
				right:    IntegerLiteral{value: 1},
			},
		},
	}
}

func isPublic(name string) bool {
	if len(name) > 0 {
		return unicode.IsUpper([]rune(name)[0]) // https://www.reddit.com/r/golang/comments/11cig0a/comment/ja371qd/?utm_source=share&utm_medium=web3x&utm_name=web3xcss&utm_term=1&utm_content=share_button
	}
	return false
}

func buildPrefixedName(prefix string, funcName string) string {
	if len(prefix) > 0 {
		prefix = fmt.Sprintf("%s_", prefix)

		// Only prefix if it doesn't already have the prefix.
		if !strings.HasPrefix(funcName, prefix) {
			funcName = fmt.Sprintf("%s%s", prefix, funcName)
		}
	}
	return funcName
}

func (p *Parser) atError(what string, token lexer.Token) error {
	return fmt.Errorf("%s at row %d, column %d: %s", what, token.Row(), token.Column(), p.path)
}

func (p *Parser) expectedError(what string, token lexer.Token) error {
	return p.atError(fmt.Sprintf("expected %s", what), token)
}

func (p Parser) peek() lexer.Token {
	return p.peekAt(0)
}

func (p Parser) peekAt(add uint) lexer.Token {
	index := p.index + int(add)
	tokens := p.tokens
	token := lexer.Token{}

	if index < len(tokens) {
		token = tokens[index]
	}
	return token
}

func (p Parser) findAllowed(searchTokenType lexer.TokenType, allowed ...lexer.TokenType) (lexer.Token, error) {
	tokens := p.tokens

	for i := p.index; i < len(tokens); i++ {
		token := tokens[i]
		tokenType := token.Type()

		if tokenType == searchTokenType {
			return token, nil
		}

		if !slices.Contains(allowed, tokenType) {
			return lexer.Token{}, fmt.Errorf(`found illegal token "%d" before "%d"`, tokenType, searchTokenType)
		}
	}
	return lexer.Token{}, fmt.Errorf(`token type "%d" not found`, searchTokenType)
}

func (p Parser) findBefore(searchTokenType lexer.TokenType, before ...lexer.TokenType) (lexer.Token, error) {
	tokens := p.tokens

	for i := p.index; i < len(tokens); i++ {
		token := tokens[i]
		tokenType := token.Type()

		if tokenType == searchTokenType {
			return token, nil
		}

		for _, tokenTypeTemp := range before {
			if tokenTypeTemp == tokenType {
				return lexer.Token{}, fmt.Errorf(`found "%d" before "%d"`, tokenTypeTemp, tokenType)
			}
		}
	}
	return lexer.Token{}, fmt.Errorf(`token type "%d" not found`, searchTokenType)
}

func (p *Parser) eat() lexer.Token {
	token := p.peek()
	p.index++

	return token
}

func (p *Parser) isShortVarInit() bool {
	_, err := p.findAllowed(lexer.SHORT_INIT_OPERATOR, lexer.IDENTIFIER, lexer.COMMA)

	// Short initialization is an arbitrary number of identifiers and commas plus the short init operator (e.g. x, y := ...).
	return err == nil
}

func (p *Parser) checkNewVariableNameToken(token lexer.Token, ctx context) error {
	name := token.Value()
	_, exists := ctx.findVariable(name, p.prefix, ctx.global())

	if exists {
		return p.atError(fmt.Sprintf("variable %s has already been defined", name), token)
	}
	return nil
}

func (p *Parser) getUsedFuncs(startFunc string) []string {
	usedFuncs := []string{}
	startFunc = strings.TrimSpace(startFunc)

	if usedFuncsTemp, exists := p.usedFuncs[startFunc]; exists {
		if len(startFunc) > 0 && !slices.Contains(usedFuncs, startFunc) {
			usedFuncs = append(usedFuncs, startFunc)
		}

		for _, usedFuncTemp := range usedFuncsTemp {
			if !slices.Contains(usedFuncs, usedFuncTemp) {
				usedFuncs = append(usedFuncs, usedFuncTemp)
			}
			usedSubFuncs := p.getUsedFuncs(usedFuncTemp)

			for _, usedSubFunc := range usedSubFuncs {
				if !slices.Contains(usedFuncs, usedSubFunc) {
					usedFuncs = append(usedFuncs, usedSubFunc)
				}
			}
		}
	}
	return usedFuncs
}

func (p *Parser) cleanProgram(program Program) (Program, error) {
	statements := program.Body()
	usedFuncs := p.getUsedFuncs("")

	// Remove all functions that are not being used.
	statements = slices.DeleteFunc(statements, func(stmt Statement) bool {
		switch stmt.StatementType() {
		case STATEMENT_TYPE_FUNCTION_DEFINITION:
			function := stmt.(FunctionDefinition)
			return !slices.Contains(usedFuncs, function.Name())
		}
		return false
	})
	return Program{
		body: statements,
	}, nil
}

func (p *Parser) evaluateVarNames() ([]lexer.Token, error) {
	nameTokens := []lexer.Token{}

	for {
		nextToken := p.eat()

		if nextToken.Type() != lexer.IDENTIFIER {
			return nil, p.expectedError("variable name", nextToken)
		}
		nameTokens = append(nameTokens, nextToken)
		nextToken = p.peek()

		if nextToken.Type() != lexer.COMMA {
			break
		}
		p.eat() // Eat comma token.
	}
	return nameTokens, nil
}

func (p *Parser) evaluateValues(ctx context) (evaluatedValues, error) {
	expressions := []Expression{}

	for {
		exprToken := p.peek()
		expr, err := p.evaluateExpression(ctx)

		if err != nil {
			return evaluatedValues{}, err
		}
		expressions = append(expressions, expr)
		nextToken := p.peek()
		returnValuesLength := -1
		funcName := ""

		// If expression is a function, check if it returns a value.
		if expr.StatementType() == STATEMENT_TYPE_FUNCTION_CALL {
			call := expr.(FunctionCall)
			returnValuesLength = len(call.ReturnTypes())
			funcName = call.Name()

			if returnValuesLength == 0 {
				return evaluatedValues{}, p.expectedError(fmt.Sprintf(`return value from function "%s"`, funcName), exprToken)
			}
		}
		// Check if other values follow.
		if nextToken.Type() != lexer.COMMA {
			break
		}
		p.eat() // Eat comma token.

		// If other values follow, function must only return one value.
		if returnValuesLength > 1 {
			return evaluatedValues{}, p.expectedError(fmt.Sprintf(`only one return value from function "%s"`, funcName), exprToken)
		}
	}
	return evaluatedValues{
		values: expressions,
	}, nil
}

func (p *Parser) evaluateBuiltInFunction(tokenType lexer.TokenType, keyword string, minArgs int, maxArg int, ctx context, stmtCallout func(keywordToken lexer.Token, expressions []Expression) (Statement, error)) (Statement, error) {
	keywordToken := p.eat()

	if keywordToken.Type() != tokenType {
		return nil, p.expectedError(fmt.Sprintf("%s-keyword", keyword), keywordToken)
	}
	nextToken := p.eat()

	// Make sure after the print call comes a  opening round bracket.
	if nextToken.Type() != lexer.OPENING_ROUND_BRACKET {
		return nil, p.expectedError(`"("`, nextToken)
	}
	expressions := []Expression{}
	nextToken = p.peek()

	// Evaluate arguments if it's a print call with arguments.
	if nextToken.Type() != lexer.CLOSING_ROUND_BRACKET {
		for {
			expr, err := p.evaluateExpression(ctx)

			if err != nil {
				return nil, err
			}
			expressions = append(expressions, expr)
			nextToken = p.peek()
			nextTokenType := nextToken.Type()

			if nextTokenType == lexer.COMMA {
				p.eat()
			} else if nextTokenType == lexer.CLOSING_ROUND_BRACKET {
				break
			} else {
				return nil, p.expectedError(`"," or ")"`, nextToken)
			}
		}
	}
	expressionsLength := len(expressions)

	if minArgs < 0 {
		minArgs = 0
	}
	if expressionsLength < minArgs {
		return nil, p.expectedError(fmt.Sprintf("at least %d arguments for %s", minArgs, keyword), keywordToken)
	}
	if maxArg >= 0 && expressionsLength > maxArg {
		return nil, p.expectedError(fmt.Sprintf("a maximum of %d arguments for %s", minArgs, keyword), keywordToken)
	}
	nextToken = p.eat()

	// Make sure print call is terminated with a closing round bracket.
	if nextToken.Type() != lexer.CLOSING_ROUND_BRACKET {
		return nil, p.expectedError(`")"`, nextToken)
	}
	return stmtCallout(keywordToken, expressions)
}

func (p *Parser) evaluateProgram() (Program, error) {
	ctx := context{
		imports:   map[string]string{},
		variables: map[string]Variable{},
		functions: map[string]FunctionDefinition{},
	}
	statements, err := p.evaluateImports(ctx)

	if err != nil {
		return Program{}, err
	}

	// Add own hash to imports for easier mapping handling.
	prefix := p.prefix
	err = ctx.addImport(prefix, prefix)

	if err != nil {
		return Program{}, err
	}
	statementsTemp, err := p.evaluateBlockContent([]lexer.TokenType{lexer.EOF}, nil, ctx, SCOPE_PROGRAM)

	if err != nil {
		return Program{}, err
	}
	statements = append(statements, statementsTemp...)

	return Program{
		body: statements,
	}, nil
}

func (p *Parser) evaluateImports(ctx context) ([]Statement, error) {
	var nextToken lexer.Token
	statementsTemp := []Statement{}

	// Skip empty characters.
	for {
		nextToken = p.peek()
		if !slices.Contains([]lexer.TokenType{lexer.NEWLINE}, nextToken.Type()) {
			break
		}
		p.eat()
	}

	if nextToken.Type() == lexer.IMPORT {
		p.eat()
		nextToken := p.peek()
		multiple := nextToken.Type() == lexer.OPENING_ROUND_BRACKET

		if multiple {
			p.eat()
			nextToken = p.eat()

			if nextToken.Type() != lexer.NEWLINE {
				return nil, p.expectedError("newline", nextToken)
			}
		}

		for {
			imp, err := p.evaluateImport()

			if err != nil {
				return nil, err
			}
			path := imp.path
			alias := imp.alias
			absPath := path

			// If path is relative, create an absolute path by combining the loaded path with the import path.
			if !filepath.IsAbs(absPath) {
				absPath = filepath.Join(filepath.Dir(p.path), absPath)
			}
			aliasLen := len(alias)

			// If path doesn't exist, try to find it in the standard library.
			if _, err := os.Stat(absPath); err != nil {
				ex, err := os.Executable()

				if err != nil {
					return nil, err
				}
				pathWithoutExt := strings.TrimSuffix(path, filepath.Ext(path))
				absPathTemp := filepath.Join(filepath.Dir(ex), "std", fmt.Sprintf("%s.tsh", pathWithoutExt)) // Standart library is at <executable-path>/std.

				// If path exists, use it.
				if _, err := os.Stat(absPath); err != nil {
					absPath = absPathTemp

					if aliasLen == 0 {
						alias = filepath.Base(pathWithoutExt)
					}
				}
			} else if aliasLen == 0 {
				// If it's not a standard library path, an alias must be provided.
				return nil, fmt.Errorf(`an alias must be provided for the local import "%s" in "%s"`, path, p.path)
			}
			importParser := New()
			importedProg, err := importParser.parse(absPath, true)

			if err != nil {
				return nil, err
			}

			if _, exists := ctx.findImport(alias); exists {
				return nil, fmt.Errorf(`import alias "%s" already exists`, alias)
			}
			err = ctx.addImport(alias, importParser.prefix)

			if err != nil {
				return nil, err
			}
			statementsTemp = append(statementsTemp, importedProg.Body()...)

			// Import-parser funcs with current parser funcs.
			for funcName, usedFuncs := range importParser.usedFuncs {
				if foundUsedFuncs, exists := p.usedFuncs[funcName]; !exists {
					p.usedFuncs[funcName] = usedFuncs
				} else {
					for _, usedFunc := range foundUsedFuncs {
						if !slices.Contains(foundUsedFuncs, usedFunc) {
							p.usedFuncs[funcName] = append(p.usedFuncs[funcName], usedFunc)
						}
					}
				}
			}

			nextToken = p.peek()
			nextTokenType := nextToken.Type()

			if !multiple {
				break
			} else if nextTokenType == lexer.CLOSING_ROUND_BRACKET {
				p.eat()
				break
			} else if slices.Contains([]lexer.TokenType{lexer.IDENTIFIER, lexer.STRING_LITERAL}, nextTokenType) {
				// Nothing to do, parse next import in the next cycle.
			} else {
				return nil, p.expectedError(`")"`, nextToken)
			}
		}
	}
	statements := []Statement{}

	// Add functions add variables.
	for _, statement := range statementsTemp {
		exists := false

		switch statement.StatementType() {
		case STATEMENT_TYPE_VAR_DEFINITION:
			definedVariable := statement.(VariableDefinition)

			for _, variable := range definedVariable.Variables() {
				name := variable.Name()

				if _, exists = ctx.variables[name]; !exists && variable.Public() {
					ctx.variables[name] = variable
				}
			}
		case STATEMENT_TYPE_FUNCTION_DEFINITION:
			definedFunction := statement.(FunctionDefinition)
			name := definedFunction.Name()

			if _, exists = ctx.functions[name]; !exists && definedFunction.Public() {
				ctx.functions[name] = definedFunction
			}
		}

		// Prevent code duplication.
		if !exists {
			statements = append(statements, statement)
		}
	}
	return statements, nil
}

func (p *Parser) evaluateImport() (evaluatedImport, error) {
	nextToken := p.eat()
	var alias string

	if nextToken.Type() == lexer.IDENTIFIER {
		alias = nextToken.Value()
		nextToken = p.eat()
	}

	if nextToken.Type() != lexer.STRING_LITERAL {
		return evaluatedImport{}, p.expectedError("import path", nextToken)
	}
	path := nextToken.Value()
	nextToken = p.eat()

	if !slices.Contains([]lexer.TokenType{lexer.NEWLINE, lexer.EOF}, nextToken.Type()) {
		return evaluatedImport{}, p.expectedError("newline or end-of-file", nextToken)
	}
	return evaluatedImport{
		alias,
		path,
	}, nil
}

func (p *Parser) evaluateBlockBegin() error {
	beginToken := p.eat()

	if beginToken.Type() != lexer.OPENING_CURLY_BRACKET {
		return p.expectedError("block begin", beginToken)
	}
	newlineToken := p.eat()

	if newlineToken.Type() != lexer.NEWLINE {
		return p.expectedError("newline", newlineToken)
	}
	return nil
}

func (p *Parser) evaluateBlockContent(terminationTokenTypes []lexer.TokenType, callback blockCallback, ctx context, scope scope) ([]Statement, error) {
	var err error

	statements := []Statement{}
	loop := true
	callCallback := func() error {
		errTemp := err

		if errTemp == nil && callback != nil {
			errTemp = callback(statements, !loop)
		}
		return errTemp
	}

	// Clone context to avoid modification of the original.
	ctx = ctx.clone()

	// Add scope to context.
	ctx.scopeStack = append(ctx.scopeStack, scope)

	for loop {
		token := p.peek()
		tokenType := token.Type()
		var stmt Statement

		if slices.Contains(terminationTokenTypes, tokenType) {
			loop = false // Just break on termination token.
		} else {
			switch tokenType {
			case lexer.NEWLINE:
				// Ignore termination tokens as they are handled after the switch.
			default:
				stmt, err = p.evaluateStatement(ctx)
				prefix := p.prefix

				if err != nil {
					break
				}
				global := ctx.global()

				switch stmt.StatementType() {
				case STATEMENT_TYPE_VAR_DEFINITION:
					// Store new variable.
					err = ctx.addVariables(prefix, global, stmt.(VariableDefinition).Variables()...)

					if err != nil {
						return nil, err
					}
				case STATEMENT_TYPE_VAR_DEFINITION_CALL_ASSIGNMENT:
					// Store new variable.
					err = ctx.addVariables(prefix, global, stmt.(VariableDefinitionCallAssignment).Variables()...)

					if err != nil {
						return nil, err
					}
				case STATEMENT_TYPE_FUNCTION_DEFINITION:
					// Store new function.
					err = ctx.addFunctions(prefix, global, stmt.(FunctionDefinition))

					if err != nil {
						return nil, err
					}
				}
			}
		}

		if err != nil {
			break
		}

		if !loop {
			err = callCallback()
			break
		}

		if stmt != nil {
			statements = append(statements, stmt)
			err = callCallback()

			if err != nil {
				break
			}
		}
		terminationToken := p.peek()

		// Expect newline or termination token.
		if terminationToken.Type() == lexer.NEWLINE {
			p.eat()
		} else if !slices.Contains(terminationTokenTypes, terminationToken.Type()) {
			err = p.expectedError("termination token", terminationToken)
			break
		}
	}
	return statements, err
}

func (p *Parser) evaluateBlockEnd() error {
	endToken := p.eat()

	if endToken.Type() != lexer.CLOSING_CURLY_BRACKET {
		return p.expectedError("block end", endToken)
	}
	return nil
}

func (p *Parser) evaluateBlock(callback blockCallback, ctx context, scope scope) ([]Statement, error) {
	err := p.evaluateBlockBegin()

	if err != nil {
		return nil, err
	}
	statements, err := p.evaluateBlockContent([]lexer.TokenType{lexer.CLOSING_CURLY_BRACKET}, callback, ctx, scope)

	if err != nil {
		return nil, err
	}
	err = p.evaluateBlockEnd()

	if err != nil {
		return nil, err
	}
	return statements, nil
}

func (p *Parser) evaluateValueType() (ValueType, error) {
	nextToken := p.peek()
	evaluatedType := NewValueType(DATA_TYPE_UNKNOWN, false)

	// Evaluate if value type is a slice type.
	if nextToken.Type() == lexer.OPENING_SQUARE_BRACKET {
		p.eat()             // Eat opening square bracket.
		nextToken = p.eat() // Eat closing square bracket.

		if nextToken.Type() != lexer.CLOSING_SQUARE_BRACKET {
			return evaluatedType, p.expectedError(`"]"`, nextToken)
		}
		nextToken = p.peek()
		evaluatedType.isSlice = true
	}

	// Evaluate data type.
	if nextToken.Type() != lexer.DATA_TYPE {
		return evaluatedType, p.expectedError("data type", nextToken)
	}
	p.eat() // Eat data type token.
	dataType, exists := typeMapping[nextToken.Value()]

	if !exists {
		return evaluatedType, p.expectedError("valid data type", nextToken)
	}
	evaluatedType.dataType = dataType
	return evaluatedType, nil
}

func (p *Parser) evaluateVarDefinition(ctx context) (Statement, error) {
	// Possible variable declarations/definitions:
	// var v int
	// var v int = 1
	// var v = 1
	// v := 1
	isShortVarInit := p.isShortVarInit()

	// Eat "var" token only, if the variable is not defined using the short init operator (:=).
	if !isShortVarInit {
		varToken := p.eat()

		if varToken.Type() != lexer.VAR_DEFINITION {
			return nil, p.expectedError("variable definition", varToken)
		}
	}
	nameTokens, err := p.evaluateVarNames()

	if err != nil {
		return nil, err
	}
	nameTokensLength := len(nameTokens)
	firstNameToken := nameTokens[0]

	// Check if all variables are already defined.
	if nameTokensLength > 1 {
		alreadyDefined := 0

		for _, nameToken := range nameTokens {
			err := p.checkNewVariableNameToken(nameToken, ctx)

			if err != nil {
				// Only allow "re-definition" of variable via the short init operator.
				if !isShortVarInit {
					return nil, err
				}
				alreadyDefined++
			}
		}

		if alreadyDefined == nameTokensLength {
			return nil, p.atError("no new variables", firstNameToken)
		}
	} else {
		err := p.checkNewVariableNameToken(firstNameToken, ctx)

		if err != nil {
			return nil, err
		}
	}
	specifiedType := NewValueType(DATA_TYPE_UNKNOWN, false)

	if isShortVarInit {
		nextToken := p.eat() // Eat short init operator.

		if nextToken.Type() != lexer.SHORT_INIT_OPERATOR {
			return nil, p.expectedError("short initialization operator", nextToken)
		}
	} else {
		nextToken := p.peek()

		// If next token starts a type definition, evaluate value type.
		if slices.Contains([]lexer.TokenType{lexer.DATA_TYPE, lexer.OPENING_SQUARE_BRACKET}, nextToken.Type()) {
			specifiedTypeTemp, err := p.evaluateValueType()

			if err != nil {
				return nil, err
			}
			specifiedType = specifiedTypeTemp
			nextToken = p.peek()
		}
		nextTokenType := nextToken.Type()
		dataType := specifiedType.DataType()

		// If no data type has been specified and no value is being assigned, return an error.
		if dataType == DATA_TYPE_UNKNOWN && nextTokenType != lexer.ASSIGN_OPERATOR {
			return nil, p.expectedError("data type or value assignment", nextToken)
		} else if nextTokenType == lexer.ASSIGN_OPERATOR {
			p.eat()
		}
	}
	nextToken := p.peek()
	nextTokenType := nextToken.Type()
	variables := []Variable{}

	// Fill variables slice (might not contain the final type after this step).
	for _, nameToken := range nameTokens {
		prefix := p.prefix
		global := ctx.global()
		name := nameToken.Value()
		variable, exists := ctx.findVariable(name, prefix, global)
		variableValueType := variable.ValueType()

		// If the variable already exists, make sure it has the same type as the specified type.
		if exists && specifiedType.DataType() != DATA_TYPE_UNKNOWN && !specifiedType.Equals(variableValueType) {
			return nil, p.atError(fmt.Sprintf(`variable "%s" already exists but has type %s`, name, variableValueType.String()), nextToken)
		}
		storedName := name

		if global {
			storedName = buildPrefixedName(prefix, name)
		}
		variables = append(variables, NewVariable(storedName, specifiedType, global, isPublic(name)))
	}
	values := []Expression{}

	// TODO: Improve check (avoid NEWLINE and EOF check).
	if nextTokenType != lexer.NEWLINE && nextTokenType != lexer.EOF {
		evaluatedVals, err := p.evaluateValues(ctx)

		if err != nil {
			return nil, err
		}
		values = evaluatedVals.values
		valuesTypes := []ValueType{}
		isMultiReturnFuncCall, call := evaluatedVals.isMultiReturnCall()

		// If multi-return function, get function return types, else get value types.
		if isMultiReturnFuncCall {
			valuesTypes = call.ReturnTypes()
		} else {
			for _, valueTemp := range values {
				valuesTypes = append(valuesTypes, valueTemp.ValueType())
			}
		}
		valuesTypesLen := len(valuesTypes)
		variablesLen := len(variables)

		// Check if the amount of values is equal to the amount of variable names.
		if valuesTypesLen != variablesLen {
			pluralInit := ""
			pluralValues := ""

			if valuesTypesLen != 1 {
				pluralInit = "s"
			}
			if variablesLen != 1 {
				pluralValues = "s"
			}
			return nil, p.atError(fmt.Sprintf("got %d initialisation value%s but %d variable%s", valuesTypesLen, pluralInit, variablesLen, pluralValues), nextToken)
		}

		// If a type has been specified, make sure the returned types fit this type.
		if specifiedType.DataType() != DATA_TYPE_UNKNOWN {
			for _, valueType := range valuesTypes {
				if !valueType.Equals(specifiedType) {
					return nil, p.expectedError(fmt.Sprintf("%s but got %s", specifiedType.String(), valueType.String()), nextToken)
				}
			}
		}

		// Check if variables exist and if, check if the types match.
		for i, variable := range variables {
			valueValueType := valuesTypes[i]
			variableValueType := variable.ValueType()

			if variableValueType.DataType() == DATA_TYPE_UNKNOWN {
				variables[i].valueType = valueValueType // Use index here to make sure the original variable is modified, not the copy.
			} else if !variableValueType.Equals(valueValueType) {
				return nil, p.expectedError(fmt.Sprintf("%s but got %s for variable %s", variableValueType.String(), valueValueType.String(), variable.Name()), nextToken)
			}
		}

		// If it's a function call multi assignment, build return value here.
		if isMultiReturnFuncCall {
			call := VariableDefinitionCallAssignment{
				variables,
				call,
			}
			return call, nil
		}
	}

	// If no value has been specified, define default value.
	if len(values) == 0 {
		for _, variable := range variables {
			value, err := defaultVarValue(variable.ValueType())

			if err != nil {
				return nil, err
			}
			values = append(values, value)
		}
	}
	variable := VariableDefinition{
		variables,
		values,
	}
	return variable, nil
}

func (p *Parser) evaluateCompoundAssignment(ctx context) (Statement, error) {
	nameTokens, err := p.evaluateVarNames()

	if err != nil {
		return nil, err
	}
	nameToken := nameTokens[0]
	namesLen := len(nameTokens)

	if namesLen > 1 {
		return nil, p.expectedError("a single variable on the left side", nameToken)
	}
	assignToken := p.eat()

	// Check assign token.
	if assignToken.Type() != lexer.COMPOUND_ASSIGN_OPERATOR {
		return nil, p.expectedError(`"+=", "-=", "*=", "/=" or "%="`, assignToken)
	}
	valuesToken := p.peek()
	evaluatedVals, err := p.evaluateValues(ctx)

	if err != nil {
		return nil, err
	}
	isMultiReturnFuncCall, call := evaluatedVals.isMultiReturnCall()
	values := evaluatedVals.values
	valuesTypes := []ValueType{}

	// If it's a multi return function call evaluate how many values are returned by the function.
	if isMultiReturnFuncCall {
		valuesTypes = call.ReturnTypes()
	} else {
		for _, value := range values {
			valuesTypes = append(valuesTypes, value.ValueType())
		}
	}
	valuesTypesLen := len(valuesTypes)

	if valuesTypesLen > 1 {
		return nil, p.expectedError("a single value on the right side", valuesToken)
	}
	name := nameToken.Value()

	// Make sure variable has been defined.
	definedVariable, exists := ctx.findVariable(name, p.prefix, ctx.global())

	if !exists {
		return nil, p.atError(fmt.Sprintf("variable %s has not been defined", name), nameToken)
	}
	valueType := valuesTypes[0]
	expectedValueType := definedVariable.ValueType()

	if valueType != expectedValueType {
		return nil, p.expectedError(fmt.Sprintf("%s but got %s", expectedValueType.String(), valueType.String()), valuesToken)
	}
	assignOperator := assignToken.Value()
	binaryOperator := string(assignOperator[0])

	if !slices.Contains(allowedBinaryOperators(valueType), binaryOperator) {
		return nil, p.expectedError(fmt.Sprintf(`valid %s compound assign operator but got "%s"`, valueType.String(), assignOperator), assignToken)
	}
	return VariableAssignment{
		variables: []Variable{definedVariable},
		values: []Expression{
			BinaryOperation{
				left:     VariableEvaluation{definedVariable},
				operator: binaryOperator,
				right:    values[0],
			},
		},
	}, nil
}

func (p *Parser) evaluateVarAssignment(ctx context) (Statement, error) {
	nameTokens, err := p.evaluateVarNames()

	if err != nil {
		return nil, err
	}
	assignToken := p.eat()

	// Check assign token.
	if assignToken.Type() != lexer.ASSIGN_OPERATOR {
		return nil, p.expectedError(`"="`, assignToken)
	}
	valuesToken := p.peek()
	evaluatedVals, err := p.evaluateValues(ctx)

	if err != nil {
		return nil, err
	}
	isMultiReturnFuncCall, call := evaluatedVals.isMultiReturnCall()
	valuesTypes := []ValueType{}

	// If it's a multi return function call evaluate how many values are returned by the function.
	if isMultiReturnFuncCall {
		valuesTypes = call.ReturnTypes()
	} else {
		for _, value := range evaluatedVals.values {
			valuesTypes = append(valuesTypes, value.ValueType())
		}
	}
	namesLen := len(nameTokens)
	valuesTypesLen := len(valuesTypes)

	// Make sure variables and values match in length.
	if namesLen != valuesTypesLen {
		return nil, p.atError(fmt.Sprintf("got %d values but %d variables", valuesTypesLen, namesLen), valuesToken)
	}
	variables := []Variable{}

	for i, nameToken := range nameTokens {
		name := nameToken.Value()

		// Make sure variable has been defined.
		definedVariable, exists := ctx.findVariable(name, p.prefix, ctx.global())

		if !exists {
			return nil, p.atError(fmt.Sprintf("variable %s has not been defined", name), nameToken)
		}
		valueType := valuesTypes[i]
		expectedValueType := definedVariable.ValueType()

		if valueType != expectedValueType {
			return nil, p.expectedError(fmt.Sprintf("%s but got %s", expectedValueType.String(), valueType.String()), valuesToken)
		}
		variables = append(variables, NewVariable(name, valueType, ctx.global(), isPublic(name)))
	}

	if isMultiReturnFuncCall {
		return VariableAssignmentCallAssignment{
			variables,
			call,
		}, nil
	}
	return VariableAssignment{
		variables: variables,
		values:    evaluatedVals.values,
	}, nil
}

func (p *Parser) evaluateParams(ctx context) ([]Variable, error) {
	params := []Variable{}

	for {
		nameToken := p.peek()
		nameTokenType := nameToken.Type()

		// If closing bracket has been discovered, all parameters have been parsed.
		if nameTokenType == lexer.CLOSING_ROUND_BRACKET {
			break
		}
		if nameTokenType != lexer.IDENTIFIER {
			return params, p.expectedError("parameter name", nameToken)
		}
		p.eat()

		name := nameToken.Value()
		_, exists := ctx.findVariable(name, p.prefix, false)

		if exists {
			return params, fmt.Errorf("scope already contains a variable with the name %s", name)
		}
		valueType, err := p.evaluateValueType()

		if err != nil {
			return nil, err
		}
		nextToken := p.peek()
		nextTokenType := nextToken.Type()

		if nextTokenType != lexer.COMMA && nextTokenType != lexer.CLOSING_ROUND_BRACKET {
			return params, p.expectedError(`"," or ")"`, nextToken)
		} else if nextTokenType == lexer.COMMA {
			p.eat()
		}
		params = append(params, NewVariable(name, valueType, false, false))
	}
	return params, nil
}

func (p *Parser) evaluateFunctionDefinition(ctx context) (Statement, error) {
	functionToken := p.eat()

	if !ctx.global() {
		return nil, p.expectedError("function definition at top level", functionToken)
	}
	if functionToken.Type() != lexer.FUNCTION_DEFINITION {
		return nil, p.expectedError("function definition", functionToken)
	}
	nameToken := p.eat()

	if nameToken.Type() != lexer.IDENTIFIER {
		return nil, p.expectedError("function name", nameToken)
	}
	name := nameToken.Value()

	// Make sure no function exists with the same name.
	_, exists := ctx.findFunction(name, p.prefix)

	if exists {
		return nil, p.expectedError("unique function name", nameToken)
	}
	openingBrace := p.peek()
	params := []Variable{}

	// Clone context to avoid modification of the original.
	ctx = ctx.clone()

	// Remove all variables which are not global.
	maps.DeleteFunc(ctx.variables, func(_ string, v Variable) bool {
		return !v.Global()
	})

	// If no parameters are given, the brackets are optional.
	if openingBrace.Type() == lexer.OPENING_ROUND_BRACKET {
		var err error

		p.eat()
		params, err = p.evaluateParams(ctx)

		if err != nil {
			return nil, err
		}
		closingBrace := p.eat()

		if closingBrace.Type() != lexer.CLOSING_ROUND_BRACKET {
			return nil, p.expectedError(`")"`, closingBrace)
		}
	}
	returnTypeToken := p.peek()
	multiple := false
	returnTypes := []ValueType{}

	if returnTypeToken.Type() == lexer.OPENING_ROUND_BRACKET {
		p.eat()
		returnTypeToken = p.peek()
		multiple = true
	}

	for {
		// Check if a return type has been specified.
		if slices.Contains([]lexer.TokenType{lexer.DATA_TYPE, lexer.OPENING_SQUARE_BRACKET}, returnTypeToken.Type()) {
			returnTypeTemp, err := p.evaluateValueType()

			if err != nil {
				return nil, err
			}
			returnTypes = append(returnTypes, returnTypeTemp)
		}

		if !multiple {
			break
		}
		nextToken := p.eat()
		nextTokenType := nextToken.Type()

		if nextTokenType == lexer.CLOSING_ROUND_BRACKET {
			break
		} else if nextTokenType != lexer.COMMA {
			return nil, p.expectedError(`"," or ")"`, nextToken)
		}
		returnTypeToken = p.peek()
	}

	// Add parameters to variables.
	for _, param := range params {
		err := ctx.addVariables(p.prefix, false, param)

		if err != nil {
			return nil, err
		}
	}
	prefixedName := buildPrefixedName(p.prefix, name)

	// Make sure sub-statements know in which function they are currently in.
	p.currFunc = prefixedName

	statements, err := p.evaluateBlock(func(statements []Statement, last bool) error {
		var errTemp error
		var lastStatement Statement
		length := len(statements)

		if length > 0 {
			lastStatement = statements[length-1]
		}

		if len(returnTypes) > 0 {
			// If a return value is required, the last statement must be a return statement.
			if last {
				// TODO: Add token position to errors to raise clearer error messages.
				if lastStatement == nil || lastStatement.StatementType() != STATEMENT_TYPE_RETURN {
					errTemp = fmt.Errorf(`function "%s" requires a return statement at the end of the block`, name)
				} else if returnStatement := lastStatement.(Return); len(returnStatement.Values()) != len(returnTypes) {
					errTemp = fmt.Errorf(`function "%s" requires %d return values but returns %d`, name, len(returnTypes), len(returnStatement.Values()))
				} else {
					for i, returnValue := range returnStatement.Values() {
						returnType := returnTypes[i]
						returnValueType := returnValue.ValueType()

						if !returnValueType.Equals(returnType) {
							errTemp = fmt.Errorf(`function "%s" returns %s but expects %s`, name, returnValueType.String(), returnType.String())
							break
						}
					}
				}
			}
		} else if lastStatement != nil && lastStatement.StatementType() == STATEMENT_TYPE_RETURN {
			errTemp = fmt.Errorf("function %s must not have a return statement", name)
		}
		return errTemp
	}, ctx, SCOPE_FUNCTION)

	if err != nil {
		return nil, err
	}
	p.currFunc = ""

	return FunctionDefinition{
		name:        prefixedName,
		returnTypes: returnTypes,
		params:      params,
		body:        statements,
		public:      isPublic(name),
	}, nil
}

func (p *Parser) evaluateReturn(ctx context) (Statement, error) {
	returnToken := p.eat()

	if !ctx.findScope(SCOPE_FUNCTION) {
		return nil, p.expectedError(fmt.Sprintf("return within %s-scope", SCOPE_FUNCTION), returnToken)
	}
	if returnToken.Type() != lexer.RETURN {
		return nil, p.expectedError("return-keyword", returnToken)
	}
	evaluatedVals, err := p.evaluateValues(ctx)

	if err != nil {
		return nil, err
	}
	return Return{
		values: evaluatedVals.values,
	}, nil
}

func (p *Parser) evaluateBreak(ctx context) (Statement, error) {
	breakToken := p.eat()
	breakScopes := []scope{SCOPE_FOR, SCOPE_SWITCH}
	scopeOk := false

	for _, breakScope := range breakScopes {
		if ctx.findScope(breakScope) {
			scopeOk = true
			break
		}
	}

	if !scopeOk {
		return nil, p.expectedError(fmt.Sprintf("break statement within %s-scope", strings.Join(scopesToString(breakScopes), "- or ")), breakToken)
	}
	return Break{}, nil
}

func (p *Parser) evaluateContinue(ctx context) (Statement, error) {
	continueToken := p.eat()
	breakScopes := []scope{SCOPE_FOR}
	scopeOk := false

	for _, breakScope := range breakScopes {
		if ctx.findScope(breakScope) {
			scopeOk = true
			break
		}
	}

	if !scopeOk {
		return nil, p.expectedError(fmt.Sprintf("continue statement within %s-scope", strings.Join(scopesToString(breakScopes), "- or ")), continueToken)
	}
	return Continue{}, nil
}

func (p *Parser) evaluateIf(ctx context) (Statement, error) {
	var ifStatement If

	for i := 0; true; i++ {
		ifRequired := i == 0
		nextToken := p.peek()
		nextTokenType := nextToken.Type()
		evaluateCondition := true
		var condition Expression

		// "if" needs to start with if-token.
		if ifRequired {
			if nextTokenType != lexer.IF {
				return nil, p.expectedError("if-keyword", nextToken)
			}
			p.eat()
		} else {
			if nextTokenType != lexer.ELSE {
				break
			}
			p.eat()

			if p.peek().Type() != lexer.IF {
				evaluateCondition = false
			} else {
				p.eat()
			}
		}

		if evaluateCondition {
			conditionToken := p.peek()
			expr, err := p.evaluateExpression(ctx)

			if err != nil {
				return nil, err
			}
			if !expr.ValueType().IsBool() {
				return nil, p.expectedError("boolean expression", conditionToken)
			}
			condition = expr
		}
		statements, err := p.evaluateBlock(nil, ctx, SCOPE_IF)

		if err != nil {
			return nil, err
		}

		// During the first iteration, the initial if statement has to be created.
		if ifRequired {
			ifStatement = If{
				ifBranch: IfBranch{
					condition: condition,
					body:      statements,
				},
			}
		} else {
			// If condition has not been evaluated, it is the else-branch, otherwise it's an else-if-branch.
			if !evaluateCondition {
				ifStatement.elseBranch = Else{
					body: statements,
				}
			} else {
				ifStatement.elifBranches = append(ifStatement.elifBranches, IfBranch{
					condition: condition,
					body:      statements,
				})
			}
		}
	}
	return ifStatement, nil
}

func (p *Parser) evaluateSwitch(ctx context) (Statement, error) {
	switchToken := p.eat()

	if switchToken.Type() != lexer.SWITCH {
		return nil, p.expectedError("switch-keyword", switchToken)
	}
	var switchExpr Expression
	var err error

	exprToken := p.peek()

	if exprToken.Type() == lexer.OPENING_CURLY_BRACKET {
		switchExpr = BooleanLiteral{true}
	} else {
		switchExpr, err = p.evaluateExpression(ctx)

		if err != nil {
			return nil, err
		}
	}
	switchExprValueType := switchExpr.ValueType()

	if switchExprValueType.IsSlice() {
		return nil, p.atError("slices are not allowed in switch statements", exprToken)
	}
	beginToken := p.eat()

	if beginToken.Type() != lexer.OPENING_CURLY_BRACKET {
		return nil, p.expectedError(`{`, beginToken)
	}
	nextToken := p.eat()

	if nextToken.Type() != lexer.NEWLINE {
		return nil, p.expectedError("newline", nextToken)
	}
	fakeIf := If{
		ifBranch: IfBranch{
			condition: BooleanLiteral{false}, // Use a fake if-branch that isn't entered if only a default branch has been set in switch.
			body:      []Statement{},
		},
	}
	useMock := true
	nextToken = p.peek()
	defaultSet := false

	// While switch has not been terminated, evaluate cases.
	for nextToken.Type() != lexer.CLOSING_CURLY_BRACKET {
		var compareExpr Expression
		var compareExprToken lexer.Token

		switch nextToken.Type() {
		case lexer.CASE:
			p.eat() // Eat case-token.
			compareExprToken = p.peek()
			exprTemp, err := p.evaluateExpression(ctx)

			if err != nil {
				return nil, err
			}
			compareExpr = exprTemp
		case lexer.DEFAULT:
			p.eat() // Eat default-token.
		default:
			return nil, p.expectedError(`"case", "default" or "}"`, nextToken)
		}
		colonToken := p.eat()

		if colonToken.Type() != lexer.COLON {
			return nil, p.expectedError(`":"`, colonToken)
		}
		statements, err := p.evaluateBlockContent([]lexer.TokenType{lexer.CASE, lexer.DEFAULT, lexer.CLOSING_CURLY_BRACKET}, nil, ctx, SCOPE_SWITCH)

		if err != nil {
			return nil, err
		}

		// Check if non-default case.
		if compareExpr != nil {
			compareExprValueType := compareExpr.ValueType()

			if !switchExprValueType.Equals(compareExprValueType) {
				return nil, p.atError(fmt.Sprintf("%s value cannot be compared with switch's %s value", compareExprValueType.String(), switchExprValueType.String()), compareExprToken)
			}
			ifBranch := IfBranch{
				condition: NewComparison(switchExpr, COMPARE_OPERATOR_EQUAL, compareExpr),
				body:      statements,
			}

			// If fake-if has not been overwritten, overwrite it now.
			if useMock {
				fakeIf.ifBranch = ifBranch
				useMock = false
			} else {
				fakeIf.elifBranches = append(fakeIf.elifBranches, ifBranch)
			}
		} else if !defaultSet {
			fakeIf.elseBranch = Else{
				body: statements,
			}
			defaultSet = true
		} else {
			return nil, p.atError("multiple default cases are not allowed", nextToken)
		}
		nextToken = p.peek()
	}
	p.eat() // Eat "}" token.

	if nextToken.Type() != lexer.CLOSING_CURLY_BRACKET {
		return nil, p.expectedError(`"}"`, nextToken)
	}
	return fakeIf, nil
}

func (p *Parser) evaluateFor(ctx context) (Statement, error) {
	forToken := p.eat()

	if forToken.Type() != lexer.FOR {
		return nil, p.expectedError("for-keyword", forToken)
	}
	var stmt Statement
	nextToken := p.peek()
	nextTokenType := nextToken.Type()
	nextAfterNextTokenType := p.peekAt(1).Type()

	// Clone context to avoid modification of the original.
	ctx = ctx.clone()

	// If next token is an identifier and the one after it a comma or a short-init operator and range keyword, parse a for-range statement.
	if nextTokenType == lexer.IDENTIFIER && (nextAfterNextTokenType == lexer.COMMA || (nextAfterNextTokenType == lexer.SHORT_INIT_OPERATOR && p.peekAt(2).Type() == lexer.RANGE)) {
		p.eat()
		err := p.checkNewVariableNameToken(nextToken, ctx)

		if err != nil {
			return nil, err
		}
		indexVarName := nextToken.Value()
		nextToken = p.peek()
		valueVarName := ""

		if nextToken.Type() == lexer.COMMA {
			p.eat()
			nextToken = p.eat()

			if nextToken.Type() != lexer.IDENTIFIER {
				return nil, p.expectedError("identifier", nextToken)
			}
			err = p.checkNewVariableNameToken(nextToken, ctx)

			if err != nil {
				return nil, err
			}
			valueVarName = nextToken.Value()
		}
		nextToken = p.eat()
		hasNamedVar := len(valueVarName) > 0

		if nextToken.Type() != lexer.SHORT_INIT_OPERATOR {
			return nil, p.expectedError(`":=" or ","`, nextToken)
		}
		nextToken = p.eat()

		if nextToken.Type() != lexer.RANGE {
			return nil, p.expectedError("range-keyword", nextToken)
		}
		nextToken := p.peek()
		iterableExpression, err := p.evaluateExpression(ctx)

		if err != nil {
			return nil, err
		}
		iterableValueType := iterableExpression.ValueType()
		indexVar := NewVariable(indexVarName, NewValueType(DATA_TYPE_INTEGER, false), false, false)
		var iterableEvaluation Expression

		if iterableValueType.IsSlice() {
			iterableEvaluation = SliceEvaluation{
				value:    iterableExpression,
				index:    VariableEvaluation{indexVar},
				dataType: iterableValueType.DataType(),
			}
		} else if iterableValueType.IsString() {
			iterableEvaluation = StringSubscript{
				value:      iterableExpression,
				startIndex: VariableEvaluation{indexVar},
			}
		} else {
			return nil, p.expectedError("slice or string", nextToken)
		}
		iterableValueType.isSlice = false // Make sure the value var is not a slice.
		forRangeStatements := []Statement{}

		// Add count variable.
		ctx.addVariables(p.prefix, false, indexVar)

		// If no value variable has been provided, there's no need to add it.
		if hasNamedVar {
			valueVar := NewVariable(valueVarName, iterableValueType, false, false)

			// Add value variable.
			ctx.addVariables(p.prefix, false, valueVar)

			forRangeStatements = []Statement{
				VariableAssignment{
					variables: []Variable{valueVar},
					values:    []Expression{iterableEvaluation},
				},
			}
		}

		init := VariableAssignment{
			variables: []Variable{indexVar},
			values:    []Expression{IntegerLiteral{0}},
		}
		condition := Comparison{
			left:     VariableEvaluation{indexVar},
			operator: COMPARE_OPERATOR_LESS,
			right:    Len{iterableExpression},
		}
		increment := incrementDecrementStatement(indexVar, true)
		statements, err := p.evaluateBlock(nil, ctx, SCOPE_FOR)

		if err != nil {
			return nil, err
		}

		stmt = For{
			init:      init,
			condition: condition,
			increment: increment,
			body:      append(forRangeStatements, statements...),
		}
	} else {
		var init Statement
		var condition Expression
		var increment Statement

		conditionToken := nextToken
		trueCondition := BooleanLiteral{value: true}

		// If next token is already a curly brackets, it's an endless loop without a condition.
		// Therefore create a fake condition.
		if nextTokenType == lexer.OPENING_CURLY_BRACKET {
			condition = trueCondition
		} else if _, err := p.findBefore(lexer.SEMICOLON, lexer.OPENING_CURLY_BRACKET); err == nil {
			// If a semicolon was found before the curly bracket, consider for as a three-part for-loop.
			nextToken := p.peek()

			// If the next token is not a semicolon, consider it a statement.
			if nextToken.Type() != lexer.SEMICOLON {
				init, err = p.evaluateStatement(ctx)

				if err != nil {
					return nil, err
				}
				switch init.StatementType() {
				case STATEMENT_TYPE_VAR_DEFINITION:
					// Store new variable.
					err = ctx.addVariables(p.prefix, false, init.(VariableDefinition).Variables()...)

					if err != nil {
						return nil, err
					}
				case STATEMENT_TYPE_VAR_DEFINITION_CALL_ASSIGNMENT:
					// Store new variable.
					err = ctx.addVariables(p.prefix, false, init.(VariableDefinitionCallAssignment).Variables()...)

					if err != nil {
						return nil, err
					}
				case STATEMENT_TYPE_VAR_ASSIGNMENT:
				default:
					return nil, p.expectedError("variable assignment or variable definition", nextToken)
				}
			}
			nextToken = p.eat()

			// Next token must be a semicolon.
			if nextToken.Type() != lexer.SEMICOLON {
				return nil, p.expectedError(`";"`, nextToken)
			}
			nextToken = p.peek()
			conditionToken = nextToken

			// If the next token is not a semicolon, consider it a condition.
			if nextToken.Type() != lexer.SEMICOLON {
				condition, err = p.evaluateExpression(ctx)

				if err != nil {
					return nil, err
				}
			} else {
				condition = trueCondition
			}
			nextToken = p.eat()

			// Next token must be a semicolon.
			if nextToken.Type() != lexer.SEMICOLON {
				return nil, p.expectedError(`";"`, nextToken)
			}
			nextToken = p.peek()

			if nextToken.Type() != lexer.OPENING_CURLY_BRACKET {
				increment, err = p.evaluateStatement(ctx)

				if err != nil {
					return nil, err
				}
				switch increment.StatementType() {
				case STATEMENT_TYPE_VAR_ASSIGNMENT:
				default:
					return nil, p.expectedError("variable assignment", nextToken)
				}
			}
		} else {
			exprTemp, err := p.evaluateExpression(ctx)

			if err != nil {
				return nil, err
			}
			condition = exprTemp
		}

		if !condition.ValueType().IsBool() {
			return nil, p.expectedError("boolean expression", conditionToken)
		}
		statements, err := p.evaluateBlock(nil, ctx, SCOPE_FOR)

		if err != nil {
			return nil, err
		}
		stmt = For{
			init:      init,
			condition: condition,
			increment: increment,
			body:      statements,
		}
	}
	return stmt, nil
}

func (p *Parser) evaluateVarEvaluation(ctx context) (Expression, error) {
	identifierToken := p.eat() // Eat identifier token.

	if identifierToken.Type() != lexer.IDENTIFIER {
		return nil, p.expectedError("identifier", identifierToken)
	}
	name := identifierToken.Value()
	variable, exists := ctx.findVariable(name, p.prefix, ctx.global())

	if !exists {
		return nil, p.atError(fmt.Sprintf("variable %s has not been defined", name), identifierToken)
	}
	return VariableEvaluation{
		Variable: variable,
	}, nil
}

func (p *Parser) evaluateSingleExpression(ctx context) (Expression, error) {
	var err error
	var expr Expression

	token := p.peek()
	tokenType := token.Type()
	value := token.Value()

	switch tokenType {

	// Handle literals.
	case lexer.BOOL_LITERAL:
		p.eat() // Eat bool token.
		b, err := strconv.ParseBool(value)

		if err != nil {
			return nil, err
		}
		expr = BooleanLiteral{
			value: b,
		}
	case lexer.NUMBER_LITERAL:
		p.eat() // Eat number token.
		// TODO: Implement float handling.
		integer, err := strconv.Atoi(value)

		if err != nil {
			return nil, err
		}
		expr = IntegerLiteral{
			value: integer,
		}
	case lexer.NIL_LITERAL:
		p.eat()                // Eat string token.
		expr = StringLiteral{} // nil is an empty string literal.
	case lexer.STRING_LITERAL:
		p.eat() // Eat string token.
		expr = StringLiteral{
			value: value,
		}

	// Handle groups.
	case lexer.OPENING_ROUND_BRACKET:
		p.eat() // Eat opening bracket.
		child, err := p.evaluateExpression(ctx)

		if err != nil {
			return nil, err
		}
		expr = Group{
			child: child,
		}
		closingToken := p.eat()

		if closingToken.Type() != lexer.CLOSING_ROUND_BRACKET {
			return nil, p.expectedError(`")"`, closingToken)
		}

	// Handle slice instantiation.
	case lexer.OPENING_SQUARE_BRACKET:
		expr, err = p.evaluateSliceInstantiation(ctx)

	// Handle input.
	case lexer.INPUT:
		expr, err = p.evaluateInput(ctx)

	// Handle read.
	case lexer.READ:
		expr, err = p.evaluateRead(ctx)

	// Handle copy.
	case lexer.COPY:
		expr, err = p.evaluateCopy(ctx)

	// Handle itoa.
	case lexer.ITOA:
		expr, err = p.evaluateItoa(ctx)

	// Handle exists.
	case lexer.EXISTS:
		expr, err = p.evaluateExists(ctx)

	// Handle len.
	case lexer.LEN:
		expr, err = p.evaluateLen(ctx)

	// Handle app call.
	case lexer.AT:
		expr, err = p.evaluateAppCall(ctx)

	// Handle identifiers.
	case lexer.IDENTIFIER:
		nextToken := p.peekAt(1)

		// If the current token is an identifier and the next is an opening
		// round bracket or a dot, it's a function call if the next is an
		// opening square bracket, it's a slice evaluation, otherwise it's
		// a variable evaluation.
		switch nextToken.Type() {
		case lexer.OPENING_ROUND_BRACKET, lexer.DOT:
			expr, err = p.evaluateFunctionCall(ctx)
		case lexer.OPENING_SQUARE_BRACKET:
			expr, err = p.evaluateSubscript(ctx)
		default:
			expr, err = p.evaluateVarEvaluation(ctx)
		}

	default:
		return nil, p.atError(fmt.Sprintf(`unknown expression type %d "%s"`, tokenType, value), token)
	}

	if err != nil {
		return nil, err
	}
	return expr, nil
}

// -----------------------------------------------------------------------------------------------
// This section defines the operator precedence. Call operators with higher precende first as
// in a function because higher precedence means it must be processed further down the chain.
// Learnt a lot about priority handling from this video https://www.youtube.com/watch?v=aAvL2BTHf60.
// Precedence is the same as in Go (https://go.dev/ref/spec#Operator_precedence).
func (p *Parser) evaluateUnaryOperation(ctx context) (Expression, error) {
	nextToken := p.peek()
	negate := false

	if nextToken.Type() == lexer.UNARY_OPERATOR {
		// Use nested if for possible future unary operators.
		if nextToken.Value() == UNARY_OPERATOR_NEGATE {
			negate = true
			p.eat()
		}
	}
	valueToken := p.peek()
	expr, err := p.evaluateSingleExpression(ctx)

	if err != nil {
		return nil, err
	}

	if negate {
		valueType := expr.ValueType()

		if !valueType.IsBool() {
			return nil, p.expectedError("boolean value", valueToken)
		}
		return UnaryOperation{
			expr:      expr,
			operator:  UNARY_OPERATOR_NEGATE,
			valueType: expr.ValueType(),
		}, nil
	}
	return expr, nil
}

func (p *Parser) evaluateMultiplication(ctx context) (Expression, error) {
	return p.evaluateBinaryOperation(ctx, []BinaryOperator{BINARY_OPERATOR_MULTIPLICATION, BINARY_OPERATOR_DIVISION, BINARY_OPERATOR_MODULO}, p.evaluateUnaryOperation)
}

func (p *Parser) evaluateAddition(ctx context) (Expression, error) {
	return p.evaluateBinaryOperation(ctx, []BinaryOperator{BINARY_OPERATOR_ADDITION, BINARY_OPERATOR_SUBTRACTION}, p.evaluateMultiplication)
}

func (p *Parser) evaluateLogicalAnd(ctx context) (Expression, error) {
	return p.evaluateLogicalOperation(ctx, LOGICAL_OPERATOR_AND, p.evaluateComparison)
}

func (p *Parser) evaluateLogicalOr(ctx context) (Expression, error) {
	return p.evaluateLogicalOperation(ctx, LOGICAL_OPERATOR_OR, p.evaluateLogicalAnd)
}

func (p *Parser) evaluateExpression(ctx context) (Expression, error) {
	return p.evaluateLogicalOr(ctx)
}

func (p *Parser) evaluateStatement(ctx context) (Statement, error) {
	var stmt Statement
	var err error

	token := p.peek()
	tokenType := token.Type()

	switch tokenType {
	case lexer.VAR_DEFINITION:
		stmt, err = p.evaluateVarDefinition(ctx)
	case lexer.FUNCTION_DEFINITION:
		stmt, err = p.evaluateFunctionDefinition(ctx)
	case lexer.RETURN:
		stmt, err = p.evaluateReturn(ctx)
	case lexer.IF:
		stmt, err = p.evaluateIf(ctx)
	case lexer.SWITCH:
		stmt, err = p.evaluateSwitch(ctx)
	case lexer.FOR:
		stmt, err = p.evaluateFor(ctx)
	case lexer.BREAK:
		stmt, err = p.evaluateBreak(ctx)
	case lexer.CONTINUE:
		stmt, err = p.evaluateContinue(ctx)
	case lexer.PRINT:
		stmt, err = p.evaluatePrint(ctx)
	case lexer.WRITE:
		stmt, err = p.evaluateWrite(ctx)
	case lexer.PANIC:
		stmt, err = p.evaluatePanic(ctx)
	default:
		// Variable initialization also starts with identifier but is a statement (e.g. x := 1234).
		if p.isShortVarInit() {
			stmt, err = p.evaluateVarDefinition(ctx)
		} else {
			// If token is identifier it could be a slice assignment, an increment or a decrement.
			if token.Type() == lexer.IDENTIFIER {
				switch p.peekAt(1).Type() {
				case lexer.INCREMENT_OPERATOR, lexer.DECREMENT_OPERATOR:
					stmt, err = p.evaluateIncrementDecrement(ctx)
				case lexer.COMPOUND_ASSIGN_OPERATOR:
					stmt, err = p.evaluateCompoundAssignment(ctx)
				case lexer.ASSIGN_OPERATOR, lexer.COMMA:
					stmt, err = p.evaluateVarAssignment(ctx)
				default:
					// Handle slice assignment.
					variable, exists := ctx.findVariable(token.Value(), p.prefix, ctx.global())

					// If variable has been defined and is a slice, handles slice assignment.
					if exists && variable.ValueType().IsSlice() {
						stmt, err = p.evaluateSliceAssignment(ctx)
					}
				}
			}

			if err == nil && stmt == nil {
				stmt, err = p.evaluateExpression(ctx)
			}
		}
	}
	return stmt, err
}

// -----------------------------------------------------------------------------------------------

func (p *Parser) evaluateBinaryOperation(ctx context, allowedOperators []BinaryOperator, higherPrioOperation func(ctx context) (Expression, error)) (Expression, error) {
	if higherPrioOperation == nil {
		return nil, errors.New("missing higher precedence callout")
	}
	// Call higherPrioOperation first as it has higher precedence and higher precedence means it must
	// be processed further down the chain. Learnt a lot about priority handling from this video
	// https://www.youtube.com/watch?v=aAvL2BTHf60.
	leftExpression, err := higherPrioOperation(ctx)

	if err != nil {
		return nil, err
	}

	// To implement associativity, use for-loop and keep appending same prio-expressions.
	for {
		operatorToken := p.peek()
		operator := operatorToken.Value()

		if operatorToken.Type() != lexer.BINARY_OPERATOR || !slices.Contains(allowedOperators, operator) {
			break
		}
		p.eat() // Eat operator token.
		rightExpression, err := higherPrioOperation(ctx)

		if err != nil {
			return nil, err
		}
		leftType := leftExpression.ValueType()
		rightType := rightExpression.ValueType()

		if !leftType.Equals(rightType) {
			return nil, p.expectedError(fmt.Sprintf("same binary operation types but got %s and %s", leftType.String(), rightType.String()), operatorToken)
		}
		allowedTypeOperators := allowedBinaryOperators(leftType)

		if !slices.Contains(allowedTypeOperators, operator) {
			return nil, p.expectedError(fmt.Sprintf(`valid %s operator but got "%s"`, leftType.String(), operator), operatorToken)
		}
		leftExpression = BinaryOperation{
			left:     leftExpression,
			operator: operator,
			right:    rightExpression,
		}
	}
	return leftExpression, nil
}

func (p *Parser) evaluateComparison(ctx context) (Expression, error) {
	// Call evaluateAddition first as it has higher precedence and higher precedence means it must
	// be processed further down the chain. Learnt a lot about priority handling from this video
	// https://www.youtube.com/watch?v=aAvL2BTHf60.
	leftExpression, err := p.evaluateAddition(ctx)

	if err != nil {
		return nil, err
	}
	operatorToken := p.peek()
	operator := operatorToken.Value()

	if operatorToken.Type() == lexer.COMPARE_OPERATOR {
		p.eat() // Eat operator token.
		rightExpression, err := p.evaluateComparison(ctx)

		if err != nil {
			return nil, err
		}
		leftType := leftExpression.ValueType()
		rightType := rightExpression.ValueType()

		if !leftType.Equals(rightType) {
			return nil, p.expectedError(fmt.Sprintf("same comparison types but got %s and %s", leftType.DataType(), rightType.String()), operatorToken)
		}
		allowedOperators := allowedCompareOperators(leftType)

		if !slices.Contains(allowedOperators, operator) {
			return nil, p.expectedError(fmt.Sprintf(`valid %s operator but got "%s"`, leftType.String(), operator), operatorToken)
		}
		return NewComparison(leftExpression, operator, rightExpression), nil
	}
	return leftExpression, nil
}

func (p *Parser) evaluateLogicalOperation(ctx context, operator LogicalOperator, higherPrioOperation func(ctx context) (Expression, error)) (Expression, error) {
	conditionToken := p.peek()
	leftExpression, err := higherPrioOperation(ctx)

	if err != nil {
		return nil, err
	}

	for {
		operatorToken := p.peek()

		if operatorToken.Type() != lexer.LOGICAL_OPERATOR || operatorToken.Value() != operator {
			break
		}

		if !leftExpression.ValueType().IsBool() {
			return nil, p.expectedError("boolean value", conditionToken)
		}
		p.eat() // Eat operator token.
		operatorValue := operatorToken.Value()
		rightExpression, errTemp := higherPrioOperation(ctx)

		if errTemp != nil {
			return nil, errTemp
		}
		leftExpression = LogicalOperation{
			left:     leftExpression,
			operator: operatorValue,
			right:    rightExpression,
		}
	}
	return leftExpression, nil
}

func (p *Parser) evaluateArguments(typeName string, name string, params []Variable, ctx context) ([]Expression, error) {
	var err error
	openingBraceToken := p.eat()

	if openingBraceToken.Type() != lexer.OPENING_ROUND_BRACKET {
		return nil, p.expectedError(`"("`, openingBraceToken)
	}
	nextToken := p.peek()
	args := []Expression{}
	ignoreParams := params == nil // If params is nil, arguments will not be checked for length or type.
	paramsLength := 0

	if !ignoreParams {
		paramsLength = len(params)
	}
	argsLengthError := func(amount int) error {
		return fmt.Errorf("%s %s expects %d parameters but got %d", typeName, name, paramsLength, amount)
	}

	// While next-token is not closing brace, evaluate arguments.
	for nextToken.Type() != lexer.CLOSING_ROUND_BRACKET {
		var expr Expression
		argToken := nextToken
		expr, err = p.evaluateExpression(ctx)

		if err != nil {
			return nil, err
		}
		args = append(args, expr)

		if !ignoreParams {
			argsLength := len(args)

			// Make sure arguments have not been exceeded.
			if argsLength > paramsLength {
				return nil, argsLengthError(argsLength)
			}

			// Make sure argument type fits parameter type.
			lastArgsIndex := argsLength - 1
			param := params[lastArgsIndex]
			lastParamType := param.ValueType()
			lastArgType := expr.ValueType()

			if !lastParamType.Equals(lastArgType) {
				return nil, p.expectedError(fmt.Sprintf("parameter %s (%s) but got %s", lastParamType.String(), param.Name(), lastArgType.String()), argToken)
			}
		}
		nextToken = p.peek()
		tokenType := nextToken.Type()

		if !slices.Contains([]lexer.TokenType{lexer.COMMA, lexer.CLOSING_ROUND_BRACKET}, tokenType) {
			err = p.expectedError(`"," or ")"`, nextToken)
			break
		} else if tokenType == lexer.COMMA {
			p.eat()
		}
	}

	// Check for the appropriate arguments amount.
	if !ignoreParams {
		argsLength := len(args)

		if len(args) != paramsLength {
			return nil, argsLengthError(argsLength)
		}
	}

	if err != nil {
		return nil, err
	}

	closingBraceToken := p.eat()

	if closingBraceToken.Type() != lexer.CLOSING_ROUND_BRACKET {
		return nil, p.expectedError(`")"`, closingBraceToken)
	}
	return args, nil
}

func (p *Parser) evaluateFunctionCall(ctx context) (Call, error) {
	nextToken := p.eat()
	dotToken := p.peek()
	alias := ""

	// If next token is a dot, it's an include-function call.
	if dotToken.Type() == lexer.DOT {
		p.eat()
		alias = nextToken.Value()
		nextToken = p.eat()
	}

	if nextToken.Type() != lexer.IDENTIFIER {
		return nil, p.expectedError("function identifier", nextToken)
	}
	name := nextToken.Value()
	prefix := p.prefix
	dotedName := name

	// If it's an include-function call, use provided alias.
	if len(alias) > 0 {
		prefix = alias
		dotedName = fmt.Sprintf("%s.%s", alias, name)
	}

	// Make sure function has been defined.
	definedFunction, exists := ctx.findFunction(name, prefix)

	if !exists {
		return nil, p.atError(fmt.Sprintf("function %s has not been defined", dotedName), nextToken)
	}
	args, err := p.evaluateArguments("function", dotedName, definedFunction.params, ctx)

	if err != nil {
		return nil, err
	}
	name = definedFunction.Name()
	currFunc := p.currFunc

	// Keep track of used functions.
	if _, exists := p.usedFuncs[currFunc]; !exists {
		p.usedFuncs[currFunc] = []string{}
	}
	if !slices.Contains(p.usedFuncs[currFunc], name) {
		p.usedFuncs[currFunc] = append(p.usedFuncs[currFunc], name)
	}

	return FunctionCall{
		name:        name,
		arguments:   args,
		returnTypes: definedFunction.ReturnTypes(),
	}, nil
}

func (p *Parser) evaluateAppCall(ctx context) (Call, error) {
	nextToken := p.eat()

	if nextToken.Type() != lexer.AT {
		return nil, p.expectedError(`"@"`, nextToken)
	}
	nextToken = p.eat()
	name := nextToken.Value()

	switch nextToken.Type() {
	case lexer.IDENTIFIER, lexer.STRING_LITERAL:
		// Nothing to do, those cases are valid.
	default:
		return nil, p.expectedError("program identifier or string literal", nextToken)
	}
	args, err := p.evaluateArguments("program", name, nil, ctx)

	if err != nil {
		return nil, err
	}
	call := AppCall{
		name: name,
		args: args,
	}

	if p.peek().Type() == lexer.PIPE {
		p.eat() // Eat pipe token.
		nextCall, err := p.evaluateAppCall(ctx)

		if err != nil {
			return nil, err
		}
		nextAppCall := nextCall.(AppCall)
		call.next = &nextAppCall
	}
	return call, nil
}

func (p *Parser) evaluateSliceInstantiation(ctx context) (Expression, error) {
	nextToken := p.peek()
	sliceValueType, err := p.evaluateValueType()

	if err != nil {
		return nil, err
	}
	if !sliceValueType.IsSlice() {
		return nil, p.expectedError(fmt.Sprintf("slice type but got %s", sliceValueType.String()), nextToken)
	}
	nextToken = p.eat()

	if nextToken.Type() != lexer.OPENING_CURLY_BRACKET {
		return nil, p.expectedError(`"{"`, nextToken)
	}
	nextToken = p.peek()
	values := []Expression{}

	// Evaluate slice initialization values.
	if nextToken.Type() != lexer.CLOSING_CURLY_BRACKET {
		for {
			valueToken := p.peek()
			expr, err := p.evaluateExpression(ctx)

			if err != nil {
				return nil, err
			}
			valueDataType := expr.ValueType()
			sliceElementValueType := sliceValueType
			sliceElementValueType.isSlice = false

			if !valueDataType.Equals(sliceElementValueType) {
				return nil, p.atError(fmt.Sprintf("%s cannot not be added to %s", valueDataType.String(), sliceElementValueType.String()), valueToken)
			}
			values = append(values, expr)
			nextToken = p.peek()
			nextTokenType := nextToken.Type()

			if nextTokenType == lexer.COMMA {
				p.eat()
			} else if nextTokenType == lexer.CLOSING_CURLY_BRACKET {
				break
			} else {
				return nil, p.expectedError(`"," or "}"`, nextToken)
			}
		}
	}
	nextToken = p.eat()

	if nextToken.Type() != lexer.CLOSING_CURLY_BRACKET {
		return nil, p.expectedError(`"}"`, nextToken)
	}
	return SliceInstantiation{
		dataType: sliceValueType.DataType(),
		values:   values,
	}, nil
}

func (p *Parser) evaluateSubscript(ctx context) (Expression, error) {
	var value Expression
	var err error

	valueToken := p.peek()

	switch valueToken.Type() {
	case lexer.IDENTIFIER:
		value, err = p.evaluateVarEvaluation(ctx)
	case lexer.STRING_LITERAL:
		value, err = p.evaluateExpression(ctx)
	default:
		return nil, p.expectedError("string or variable", valueToken)
	}

	if err != nil {
		return nil, err
	}
	valueType := value.ValueType()
	isSlice := valueType.IsSlice()

	if !isSlice && valueType.DataType() != DATA_TYPE_STRING {
		return nil, p.expectedError("slice or string", valueToken)
	}
	nextToken := p.eat()

	if nextToken.Type() != lexer.OPENING_SQUARE_BRACKET {
		return nil, p.expectedError(`"["`, nextToken)
	}
	nextToken = p.peek()
	startToken := nextToken
	gotRange := nextToken.Type() == lexer.COLON
	var startIndex Expression

	if gotRange {
		p.eat()
		startIndex = IntegerLiteral{0}
		startToken = p.peek()
	} else {
		startIndex, err = p.evaluateExpression(ctx)
	}

	if err != nil {
		return nil, err
	}
	startIndexValueType := startIndex.ValueType()

	if !startIndexValueType.IsInt() {
		return nil, p.expectedError(fmt.Sprintf("%s as start-index but got %s", DATA_TYPE_INTEGER, startIndexValueType.String()), startToken)
	}
	nextToken = p.peek()

	if nextToken.Type() == lexer.COLON {
		if gotRange {
			return nil, p.expectedError("only one colon", nextToken)
		}
		p.eat()
		gotRange = true
	}

	if gotRange && isSlice {
		return nil, p.atError("subscript range is not supported for slices", valueToken)
	}
	nextToken = p.peek()
	endToken := nextToken
	endIndex := startIndex

	if nextToken.Type() == lexer.CLOSING_SQUARE_BRACKET {
		// If range but no end-index is provided, create one by using Len-expression.
		if gotRange {
			endIndex = BinaryOperation{
				left:     Len{value},
				operator: BINARY_OPERATOR_SUBTRACTION,
				right:    IntegerLiteral{1},
			}
		}
		p.eat() // Eat square bracket.
	} else {
		endIndex, err = p.evaluateExpression(ctx)

		if err != nil {
			return nil, err
		}
		nextToken = p.eat()

		if nextToken.Type() != lexer.CLOSING_SQUARE_BRACKET {
			return nil, p.expectedError(`"]"`, nextToken)
		}

		// End-index is not included.
		endIndex = BinaryOperation{
			left:     endIndex,
			operator: BINARY_OPERATOR_SUBTRACTION,
			right:    IntegerLiteral{1},
		}
	}
	endIndexValueType := endIndex.ValueType()

	if !endIndexValueType.IsInt() {
		return nil, p.expectedError(fmt.Sprintf("%s as stop-index but got %s", DATA_TYPE_INTEGER, endIndexValueType.String()), endToken)
	}

	if !isSlice {
		return StringSubscript{
			value:      value,
			startIndex: startIndex,
			endIndex:   endIndex,
		}, nil
	}
	return SliceEvaluation{
		value:    value,
		index:    startIndex,
		dataType: valueType.DataType(),
	}, nil
}

func (p *Parser) evaluateSliceAssignment(ctx context) (Statement, error) {
	nameToken := p.eat()

	if nameToken.Type() != lexer.IDENTIFIER {
		return nil, p.expectedError("slice variable", nameToken)
	}
	name := nameToken.Value()
	variable, exists := ctx.findVariable(name, p.prefix, ctx.global())

	if !exists {
		return nil, p.atError(fmt.Sprintf("variable %s has not been defined", name), nameToken)
	}
	variableValueType := variable.ValueType()

	if !variableValueType.IsSlice() {
		return nil, p.expectedError(fmt.Sprintf("slice but variable is of type %s", variableValueType.String()), nameToken)
	}
	nextToken := p.eat()

	if nextToken.Type() != lexer.OPENING_SQUARE_BRACKET {
		return nil, p.expectedError(`"["`, nextToken)
	}
	nextToken = p.peek()
	index, err := p.evaluateExpression(ctx)

	if err != nil {
		return nil, err
	}
	indexValueType := index.ValueType()

	if !indexValueType.IsInt() {
		return nil, p.expectedError(fmt.Sprintf("%s as index but got %s", DATA_TYPE_INTEGER, indexValueType.String()), nextToken)
	}
	nextToken = p.eat()

	if nextToken.Type() != lexer.CLOSING_SQUARE_BRACKET {
		return nil, p.expectedError(`"]"`, nextToken)
	}
	nextToken = p.eat()

	if nextToken.Type() != lexer.ASSIGN_OPERATOR {
		return nil, p.expectedError(`"="`, nameToken)
	}
	valueToken := p.peek()
	value, err := p.evaluateExpression(ctx)

	if err != nil {
		return nil, err
	}
	variableDataType := variableValueType.DataType()
	assignedDataType := value.ValueType().DataType()

	if variableDataType != assignedDataType {
		return nil, p.expectedError(fmt.Sprintf("%s value but got %s", variableDataType, assignedDataType), valueToken)
	}
	return SliceAssignment{
		Variable: variable,
		index:    index,
		value:    value,
	}, nil
}

func (p *Parser) evaluateIncrementDecrement(ctx context) (Statement, error) {
	identifierToken := p.eat()

	if identifierToken.Type() != lexer.IDENTIFIER {
		return nil, p.expectedError("identifier", identifierToken)
	}
	name := identifierToken.Value()
	definedVariable, exists := ctx.findVariable(name, p.prefix, ctx.global())

	if !exists {
		return nil, p.atError(fmt.Sprintf("variable %s has not been defined", name), identifierToken)
	}
	valueType := definedVariable.ValueType()

	if !valueType.IsInt() {
		return nil, p.expectedError(fmt.Sprintf("%s but got %s", NewValueType(DATA_TYPE_INTEGER, false).String(), valueType.String()), identifierToken)
	}
	operationToken := p.eat()
	increment := true

	switch operationToken.Type() {
	case lexer.INCREMENT_OPERATOR:
		// Nothing to do.
	case lexer.DECREMENT_OPERATOR:
		increment = false
	default:
		return nil, p.expectedError(`"++" or "--"`, operationToken)
	}
	return incrementDecrementStatement(definedVariable, increment), nil
}

func (p *Parser) evaluateLen(ctx context) (Expression, error) {
	expr, err := p.evaluateBuiltInFunction(lexer.LEN, "len", 1, 1, ctx, func(keywordToken lexer.Token, expressions []Expression) (Statement, error) {
		expr := expressions[0]
		valueType := expr.ValueType()

		if !valueType.IsSlice() && !valueType.IsString() {
			return nil, p.expectedError("slice or string", keywordToken)
		}
		return Len{
			expression: expr,
		}, nil
	})

	if err != nil {
		return nil, err
	}
	return expr.(Len), nil
}

func (p *Parser) evaluateInput(ctx context) (Expression, error) {
	expr, err := p.evaluateBuiltInFunction(lexer.INPUT, "input", 0, 1, ctx, func(keywordToken lexer.Token, expressions []Expression) (Statement, error) {
		var expr Expression

		if len(expressions) > 0 {
			expr = expressions[0]
		}
		return Input{
			prompt: expr,
		}, nil
	})

	if err != nil {
		return nil, err
	}
	return expr.(Input), nil
}

func (p *Parser) evaluatePrint(ctx context) (Statement, error) {
	return p.evaluateBuiltInFunction(lexer.PRINT, "print", 0, -1, ctx, func(keywordToken lexer.Token, expressions []Expression) (Statement, error) {
		return Print{
			expressions: expressions,
		}, nil
	})
}

func (p *Parser) evaluatePanic(ctx context) (Statement, error) {
	return p.evaluateBuiltInFunction(lexer.PANIC, "panic", 1, 1, ctx, func(keywordToken lexer.Token, expressions []Expression) (Statement, error) {
		return Panic{
			expression: expressions[0],
		}, nil
	})
}

func (p *Parser) evaluateCopy(ctx context) (Expression, error) {
	expr, err := p.evaluateBuiltInFunction(lexer.COPY, "copy", 2, 2, ctx, func(keywordToken lexer.Token, expressions []Expression) (Statement, error) {
		expressionsLen := len(expressions)

		if expressionsLen < 1 {
			return nil, p.expectedError("destination slice as first argument", keywordToken)
		} else if expressionsLen < 2 {
			return nil, p.expectedError("source slice as second argument", keywordToken)
		}
		dst := expressions[0]
		src := expressions[1]
		dstType := dst.ValueType()
		srcType := src.ValueType()

		if !dstType.IsSlice() || dst.StatementType() != STATEMENT_TYPE_VAR_EVALUATION {
			return nil, p.expectedError("slice variable as first argument", keywordToken)
		} else if !srcType.IsSlice() {
			return nil, p.expectedError("slice as second argument", keywordToken)
		} else if !dstType.Equals(srcType) {
			return nil, p.atError(fmt.Sprintf("got %s as destination but %s as source", dstType.String(), srcType.String()), keywordToken)
		}
		dstSlice := dst.(VariableEvaluation)

		// To copy a slice, just create a for-loop.
		return Copy{
			destination: dstSlice.Variable,
			source:      src,
		}, nil
	})

	if err != nil {
		return nil, err
	}
	return expr.(Copy), nil
}

func (p *Parser) evaluateItoa(ctx context) (Expression, error) {
	expr, err := p.evaluateBuiltInFunction(lexer.ITOA, "itoa", 1, 1, ctx, func(keywordToken lexer.Token, expressions []Expression) (Statement, error) {
		value := expressions[0]

		if !value.ValueType().IsInt() {
			return nil, p.expectedError("integer", keywordToken)
		}
		return Itoa{
			value: value,
		}, nil
	})

	if err != nil {
		return nil, err
	}
	return expr.(Itoa), nil
}

func (p *Parser) evaluateExists(ctx context) (Expression, error) {
	expr, err := p.evaluateBuiltInFunction(lexer.EXISTS, "exists", 1, 1, ctx, func(keywordToken lexer.Token, expressions []Expression) (Statement, error) {
		path := expressions[0]

		if !path.ValueType().IsString() {
			return nil, p.expectedError("path string", keywordToken)
		}
		return Exists{
			path: path,
		}, nil
	})

	if err != nil {
		return nil, err
	}
	return expr.(Exists), nil
}

func (p *Parser) evaluateRead(ctx context) (Expression, error) {
	expr, err := p.evaluateBuiltInFunction(lexer.READ, "read", 1, 1, ctx, func(keywordToken lexer.Token, expressions []Expression) (Statement, error) {
		path := expressions[0]

		if !path.ValueType().IsString() {
			return nil, p.expectedError("file path string as first parameter", keywordToken)
		}
		return Read{
			path: path,
		}, nil
	})

	if err != nil {
		return nil, err
	}
	return expr.(Read), nil
}

func (p *Parser) evaluateWrite(ctx context) (Expression, error) {
	expr, err := p.evaluateBuiltInFunction(lexer.WRITE, "write", 2, 3, ctx, func(keywordToken lexer.Token, expressions []Expression) (Statement, error) {
		path := expressions[0]

		if !path.ValueType().IsString() {
			return nil, p.expectedError("file path string as first parameter", keywordToken)
		}
		data := expressions[1]

		if !path.ValueType().IsString() {
			return nil, p.expectedError("data string as second parameter", keywordToken)
		}
		var append Expression = BooleanLiteral{false}

		if len(expressions) > 2 {
			append = expressions[2]

			if !append.ValueType().IsBool() {
				return nil, p.expectedError("append boolean as third parameter", keywordToken)
			}
		}
		return Write{
			path:   path,
			data:   data,
			append: append,
		}, nil
	})

	if err != nil {
		return nil, err
	}
	return expr.(Write), nil
}
