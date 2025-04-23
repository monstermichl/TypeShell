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
	imports    map[string]string
	variables  map[string]Variable
	functions  map[string]FunctionDefinition
	scopeStack []scope
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

func (c context) buildPrefixedName(name string, prefix string) (string, error) {
	name = strings.TrimSpace(name)
	prefix = strings.TrimSpace(prefix)

	if len(prefix) > 0 {
		hash, exists := c.imports[prefix]

		if !exists {
			return "", fmt.Errorf("prefix \"%s\" not found", prefix)
		}
		name = buildPrefixedName(hash, name)
	}
	return name, nil
}

func (c context) findImport(alias string) (string, bool) {
	hash, exists := c.imports[alias]
	return hash, exists
}

func (c context) findVariable(name string, prefix string) (Variable, bool) {
	prefixedName, err := c.buildPrefixedName(name, prefix)

	if err != nil {
		return Variable{}, false
	}
	variable, exists := c.variables[prefixedName]
	return variable, exists
}

func (c context) findFunction(name string, prefix string) (FunctionDefinition, bool) {
	prefixedName, err := c.buildPrefixedName(name, prefix)

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

func (ev evaluatedValues) isMultiReturnFuncCall() (bool, FunctionCall) {
	var call FunctionCall
	values := ev.values
	multi := false

	if len(values) == 1 && values[0].StatementType() == STATEMENT_TYPE_FUNCTION_CALL {
		callTemp := values[0].(FunctionCall)

		if len(callTemp.ReturnTypes()) > 1 {
			multi = true
			call = callTemp
		}
	}
	return multi, call
}

type blockCallback func(statements []Statement, last bool) error

type Parser struct {
	tokens []lexer.Token
	index  int
	path   string
	prefix string
}

func New() Parser {
	return Parser{}
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

	// Create prefix.
	prefix := ""

	if imported {
		h := sha256.New()
		h.Write(source)

		prefix = fmt.Sprintf("%x", h.Sum(nil))[0:7] // Only use the 7 first characters (inspired by Git).
	}
	p.prefix = strings.TrimSpace(prefix)

	return p.evaluateProgram()
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
	return nil, fmt.Errorf("no default value found for type %s", valueType.ToString())
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
				operator:  operation,
				right:     IntegerLiteral{value: 1},
				valueType: NewValueType(DATA_TYPE_INTEGER, false),
			},
		},
	}
}

func buildPrefixedName(alias string, funcName string) string {
	if len(alias) > 0 {
		prefix := fmt.Sprintf("%s_", alias)

		// Only prefix if it doesn't already have the prefix.
		if !strings.HasPrefix(funcName, prefix) {
			funcName = fmt.Sprintf("%s%s", prefix, funcName)
		}
	}
	return funcName
}

func (p *Parser) buildPrefixedName(funcName string) string {
	return buildPrefixedName(p.prefix, funcName)
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
			return lexer.Token{}, fmt.Errorf("found illegal token \"%d\" before \"%d\"", tokenType, searchTokenType)
		}
	}
	return lexer.Token{}, fmt.Errorf("token type \"%d\" not found", searchTokenType)
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
				return lexer.Token{}, fmt.Errorf("found \"%d\" before \"%d\"", tokenTypeTemp, tokenType)
			}
		}
	}
	return lexer.Token{}, fmt.Errorf("token type \"%d\" not found", searchTokenType)
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
	_, exists := ctx.findVariable(name, "")

	if exists {
		return p.atError(fmt.Sprintf("variable %s has already been defined", name), token)
	}
	return nil
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
				return evaluatedValues{}, p.expectedError(fmt.Sprintf("return value from function \"%s\"", funcName), exprToken)
			}
		}
		// Check if other values follow.
		if nextToken.Type() != lexer.COMMA {
			break
		}
		p.eat() // Eat comma token.

		// If other values follow, function must only return one value.
		if returnValuesLength > 1 {
			return evaluatedValues{}, p.expectedError(fmt.Sprintf("only one return value from function \"%s\"", funcName), exprToken)
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
		return nil, p.expectedError("\"(\"", nextToken)
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
				return nil, p.expectedError("\",\" or \")\"", nextToken)
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
		return nil, p.expectedError("\")\"", nextToken)
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
	statementsTemp, err := p.evaluateBlockContent(lexer.EOF, nil, ctx, SCOPE_PROGRAM)

	if err != nil {
		return Program{}, err
	}
	statements = append(statements, statementsTemp...)

	return Program{
		body: statements,
	}, nil
}

func (p *Parser) evaluateImports(ctx context) ([]Statement, error) {
	nextToken := p.peek()
	statementsTemp := []Statement{}

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

			// If path is relative, create an absolute path by combining the loaded path with the import path.
			if !filepath.IsAbs(path) {
				path = filepath.Join(filepath.Dir(p.path), path)
			}
			importParser := New()
			importedProg, err := importParser.parse(path, true)

			if err != nil {
				return nil, err
			}
			alias := imp.alias

			if _, exists := ctx.findImport(alias); exists {
				return nil, fmt.Errorf("import alias \"%s\" already exists", alias)
			}
			ctx.imports[alias] = importParser.prefix
			statementsTemp = append(statementsTemp, importedProg.Body()...)

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
				return nil, p.expectedError("\")\"", nextToken)
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

				if _, exists = ctx.variables[name]; !exists {
					ctx.variables[name] = variable
				}
			}
		case STATEMENT_TYPE_FUNCTION_DEFINITION:
			definedFunction := statement.(FunctionDefinition)
			name := definedFunction.Name()

			if _, exists = ctx.functions[name]; !exists {
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

func (p *Parser) evaluateBlockContent(terminationTokenType lexer.TokenType, callback blockCallback, ctx context, scope scope) ([]Statement, error) {
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

		switch tokenType {
		case terminationTokenType:
			// Just break on termination token.
			loop = false
		case lexer.NEWLINE:
			// Ignore termination tokens as they are handled after the switch.
		default:
			stmt, err = p.evaluateStatement(ctx)

			if err != nil {
				break
			}
			switch stmt.StatementType() {
			case STATEMENT_TYPE_VAR_DEFINITION:
				// Store new variable.
				for _, variable := range stmt.(VariableDefinition).Variables() {
					ctx.variables[p.buildPrefixedName(variable.Name())] = variable
				}
			case STATEMENT_TYPE_VAR_DEFINITION_CALL_ASSIGNMENT:
				// Store new variable.
				for _, variable := range stmt.(VariableDefinitionCallAssignment).Variables() {
					ctx.variables[p.buildPrefixedName(variable.Name())] = variable
				}
			case STATEMENT_TYPE_FUNCTION_DEFINITION:
				// Store new function.
				function := stmt.(FunctionDefinition)
				ctx.functions[function.Name()] = function
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
		} else if terminationToken.Type() != terminationTokenType {
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
	statements, err := p.evaluateBlockContent(lexer.CLOSING_CURLY_BRACKET, callback, ctx, scope)

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
			return evaluatedType, p.expectedError("\"]\"", nextToken)
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
		name := nameToken.Value()
		prefixedName := p.buildPrefixedName(name)
		variable, exists := ctx.findVariable(name, p.prefix)
		variableValueType := variable.ValueType()

		// If the variable already exists, make sure it has the same type as the specified type.
		if exists && specifiedType.DataType() != DATA_TYPE_UNKNOWN && !specifiedType.Equals(variableValueType) {
			return nil, p.atError(fmt.Sprintf("variable \"%s\" already exists but has type %s", prefixedName, variableValueType.ToString()), nextToken)
		}
		storedName := name
		global := ctx.global()

		if global {
			storedName = prefixedName
		}
		variables = append(variables, NewVariable(storedName, specifiedType, global))
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
		isMultiReturnFuncCall, call := evaluatedVals.isMultiReturnFuncCall()

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
			return nil, p.atError(fmt.Sprintf("got %d initialisation values but %d variables", valuesTypesLen, variablesLen), nextToken)
		}

		// If a type has been specified, make sure the returned types fit this type.
		if specifiedType.DataType() != DATA_TYPE_UNKNOWN {
			for _, valueType := range valuesTypes {
				if !valueType.Equals(specifiedType) {
					return nil, p.expectedError(fmt.Sprintf("%s but got %s", specifiedType.ToString(), valueType.ToString()), nextToken)
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
				return nil, p.expectedError(fmt.Sprintf("%s but got %s for variable %s", variableValueType.ToString(), valueValueType.ToString(), variable.Name()), nextToken)
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

func (p *Parser) evaluateVarAssignment(ctx context) (Statement, error) {
	nameTokens, err := p.evaluateVarNames()

	if err != nil {
		return nil, err
	}
	assignToken := p.eat()

	// Check assign token.
	if assignToken.Type() != lexer.ASSIGN_OPERATOR {
		return nil, p.expectedError("\"=\"", assignToken)
	}
	valuesToken := p.peek()
	evaluatedVals, err := p.evaluateValues(ctx)

	if err != nil {
		return nil, err
	}
	isMultiReturnFuncCall, call := evaluatedVals.isMultiReturnFuncCall()
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
		definedVariable, exists := ctx.variables[name]

		if !exists {
			return nil, p.atError(fmt.Sprintf("variable %s has not been defined", name), nameToken)
		}
		valueType := valuesTypes[i]
		expectedValueType := definedVariable.ValueType()

		if valueType != expectedValueType {
			return nil, p.expectedError(fmt.Sprintf("%s but got %s", expectedValueType.ToString(), valueType.ToString()), valuesToken)
		}
		variables = append(variables, NewVariable(name, valueType, ctx.global()))
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
		_, exists := ctx.variables[name]

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
			return params, p.expectedError("\",\" or \")\"", nextToken)
		} else if nextTokenType == lexer.COMMA {
			p.eat()
		}
		params = append(params, NewVariable(name, valueType, false))
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
			return nil, p.expectedError("closing bracket", closingBrace)
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
			return nil, p.expectedError("\",\" or \")\"", nextToken)
		}
		returnTypeToken = p.peek()
	}

	// Add parameters to variables.
	for _, param := range params {
		ctx.variables[param.Name()] = param
	}

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
					errTemp = fmt.Errorf("function \"%s\" requires a return statement at the end of the block", name)
				} else if returnStatement := lastStatement.(Return); len(returnStatement.Values()) != len(returnTypes) {
					errTemp = fmt.Errorf("function \"%s\" requires %d return values but returns %d", name, len(returnTypes), len(returnStatement.Values()))
				} else {
					for i, returnValue := range returnStatement.Values() {
						returnType := returnTypes[i]
						returnValueType := returnValue.ValueType()

						if !returnValueType.Equals(returnType) {
							errTemp = fmt.Errorf("function \"%s\" returns %s but expects %s", name, returnValueType.ToString(), returnType.ToString())
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

	return FunctionDefinition{
		name:        p.buildPrefixedName(name),
		returnTypes: returnTypes,
		params:      params,
		body:        statements,
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

func (p *Parser) evaluateFor(ctx context) (Statement, error) {
	forToken := p.eat()

	if forToken.Type() != lexer.FOR {
		return nil, p.expectedError("for-keyword", forToken)
	}
	var stmt Statement
	nextToken := p.peek()
	nextTokenType := nextToken.Type()

	// Clone context to avoid modification of the original.
	ctx = ctx.clone()

	// If next token is an identifier and the one after it a comma, parse a for-range statement.
	if nextTokenType == lexer.IDENTIFIER && p.peekAt(1).Type() == lexer.COMMA {
		p.eat()
		err := p.checkNewVariableNameToken(nextToken, ctx)

		if err != nil {
			return nil, err
		}
		indexVarName := nextToken.Value()
		nextToken = p.eat()

		if nextToken.Type() != lexer.COMMA {
			return nil, p.expectedError("\",\"", nextToken)
		}
		nextToken = p.eat()
		err = p.checkNewVariableNameToken(nextToken, ctx)

		if err != nil {
			return nil, err
		}
		valueVarName := nextToken.Value()
		nextToken = p.eat()

		if nextToken.Type() != lexer.SHORT_INIT_OPERATOR {
			return nil, p.expectedError("\":=\"", nextToken)
		}
		nextToken = p.eat()

		if nextToken.Type() != lexer.RANGE {
			return nil, p.expectedError("range-keyword", nextToken)
		}

		// To make transpilation easier, only allow a variable-identifier here instead of an expression.
		// This is necessary to have an identifier for the slice for converting the for-range-loop into
		// a for-loop.
		// sliceIdentifierToken := p.eat()
		// err = p.checkNewVariableNameToken(sliceIdentifierToken, ctx)
		nextToken := p.peek()
		sliceExpression, err := p.evaluateExpression(ctx)

		if err != nil {
			return nil, err
		}
		// sliceName := sliceIdentifierToken.Value()
		// sliceVariable := ctx.variables[sliceName]

		if !sliceExpression.ValueType().isSlice {
			return nil, p.expectedError("slice", nextToken)
		}
		sliceValueType := sliceExpression.ValueType()
		indexVar := NewVariable(indexVarName, NewValueType(DATA_TYPE_INTEGER, false), false)
		valueVar := NewVariable(valueVarName, sliceValueType, false)

		// Add block variables.
		ctx.variables[indexVarName] = indexVar
		ctx.variables[valueVarName] = valueVar

		statements, err := p.evaluateBlock(nil, ctx, SCOPE_FOR)

		// Remove block variables.
		delete(ctx.variables, indexVarName)
		delete(ctx.variables, valueVarName)

		if err != nil {
			return nil, err
		}
		stmt = ForRange{
			indexVar: indexVar,
			valueVar: valueVar,
			slice:    sliceExpression,
			body:     statements,
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
					for _, variable := range init.(VariableDefinition).Variables() {
						ctx.variables[variable.Name()] = variable
					}
				case STATEMENT_TYPE_VAR_DEFINITION_CALL_ASSIGNMENT:
					// Store new variable.
					for _, variable := range init.(VariableDefinitionCallAssignment).Variables() {
						ctx.variables[variable.Name()] = variable
					}
				case STATEMENT_TYPE_VAR_ASSIGNMENT:
				default:
					return nil, p.expectedError("variable assignment or variable definition", nextToken)
				}
			}
			nextToken = p.eat()

			// Next token must be a semicolon.
			if nextToken.Type() != lexer.SEMICOLON {
				return nil, p.expectedError("\";\"", nextToken)
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
				return nil, p.expectedError("\";\"", nextToken)
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
			return nil, p.expectedError("closing bracket", closingToken)
		}

	// Handle slice instantiation.
	case lexer.OPENING_SQUARE_BRACKET:
		expr, err = p.evaluateSliceInstantiation(ctx)

	// Handle input.
	case lexer.INPUT:
		expr, err = p.evaluateInput(ctx)

	// Handle copy.
	case lexer.COPY:
		expr, err = p.evaluateCopy(ctx)

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
			p.eat() // Eat identifier token.
			name := token.Value()
			variable, exists := ctx.findVariable(name, p.prefix)

			if !exists {
				err = p.atError(fmt.Sprintf("variable %s has not been defined", name), nextToken)
			} else {
				expr = VariableEvaluation{
					Variable: variable,
				}
			}
		}

	default:
		return nil, p.atError(fmt.Sprintf("unknown expression type %d \"%s\"", tokenType, value), token)
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

	if nextToken.Value() == UNARY_OPERATOR_NEGATE {
		negate = true
		p.eat()
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
	case lexer.FOR:
		stmt, err = p.evaluateFor(ctx)
	case lexer.BREAK:
		stmt, err = p.evaluateBreak(ctx)
	case lexer.CONTINUE:
		stmt, err = p.evaluateContinue(ctx)
	case lexer.PRINT:
		stmt, err = p.evaluatePrint(ctx)
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
				case lexer.ASSIGN_OPERATOR, lexer.COMMA:
					stmt, err = p.evaluateVarAssignment(ctx)
				default:
					// Handle slice assignment.
					variable, exists := ctx.findVariable(token.Value(), p.prefix)

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
	operatorToken := p.peek()
	operator := operatorToken.Value()

	if operatorToken.Type() == lexer.BINARY_OPERATOR && slices.Contains(allowedOperators, operator) {
		p.eat() // Eat operator token.
		rightExpression, err := p.evaluateBinaryOperation(ctx, allowedOperators, higherPrioOperation)

		if err != nil {
			return nil, err
		}
		leftType := leftExpression.ValueType()
		rightType := rightExpression.ValueType()

		if !leftType.Equals(rightType) {
			return nil, p.expectedError(fmt.Sprintf("same binary operation types but got %s and %s", leftType.ToString(), rightType.ToString()), operatorToken)
		}
		allowedOperators = allowedBinaryOperators(leftType)

		if !slices.Contains(allowedOperators, operator) {
			return nil, p.expectedError(fmt.Sprintf("valid %s operator but got \"%s\"", leftType.ToString(), operator), operatorToken)
		}
		return BinaryOperation{
			left:      leftExpression,
			operator:  operator,
			right:     rightExpression,
			valueType: leftExpression.ValueType(),
		}, nil
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
			return nil, p.expectedError(fmt.Sprintf("same comparison types but got %s and %s", leftType.DataType(), rightType.ToString()), operatorToken)
		}
		allowedOperators := allowedCompareOperators(leftType)

		if !slices.Contains(allowedOperators, operator) {
			return nil, p.expectedError(fmt.Sprintf("valid %s operator but got \"%s\"", leftType.ToString(), operator), operatorToken)
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
	operatorToken := p.peek()

	if operatorToken.Type() == lexer.LOGICAL_OPERATOR && operatorToken.Value() == operator {
		if !leftExpression.ValueType().IsBool() {
			return nil, p.expectedError("boolean value", conditionToken)
		}
		p.eat() // Eat operator token.
		operatorValue := operatorToken.Value()
		rightExpression, errTemp := p.evaluateLogicalOperation(ctx, operator, higherPrioOperation)

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
		return nil, p.expectedError("opening bracket", openingBraceToken)
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
				return nil, p.expectedError(fmt.Sprintf("parameter %s (%s) but got %s", lastParamType.ToString(), param.Name(), lastArgType.ToString()), argToken)
			}
		}
		nextToken = p.peek()
		tokenType := nextToken.Type()

		if !slices.Contains([]lexer.TokenType{lexer.COMMA, lexer.CLOSING_ROUND_BRACKET}, tokenType) {
			err = p.expectedError("comma or closing bracket", nextToken)
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
		return nil, p.expectedError("closing bracket", closingBraceToken)
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
	return FunctionCall{
		name:        definedFunction.Name(),
		arguments:   args,
		returnTypes: definedFunction.ReturnTypes(),
	}, nil
}

func (p *Parser) evaluateAppCall(ctx context) (Call, error) {
	nextToken := p.eat()

	if nextToken.Type() != lexer.AT {
		return nil, p.expectedError("\"@\"", nextToken)
	}
	nextToken = p.eat()
	name := nextToken.Value()

	if nextToken.Type() != lexer.IDENTIFIER {
		return nil, p.expectedError("program identifier", nextToken)
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
		return nil, p.expectedError(fmt.Sprintf("slice type but got %s", sliceValueType.ToString()), nextToken)
	}
	nextToken = p.eat()

	if nextToken.Type() != lexer.OPENING_CURLY_BRACKET {
		return nil, p.expectedError("\"{\"", nextToken)
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
				return nil, p.atError(fmt.Sprintf("%s cannot not be added to %s", valueDataType.ToString(), sliceElementValueType.ToString()), valueToken)
			}
			values = append(values, expr)
			nextToken = p.peek()
			nextTokenType := nextToken.Type()

			if nextTokenType == lexer.COMMA {
				p.eat()
			} else if nextTokenType == lexer.CLOSING_CURLY_BRACKET {
				break
			} else {
				return nil, p.expectedError("\",\" or \"}\"", nextToken)
			}
		}
	}
	nextToken = p.eat()

	if nextToken.Type() != lexer.CLOSING_CURLY_BRACKET {
		return nil, p.expectedError("\"}\"", nextToken)
	}
	return SliceInstantiation{
		dataType: sliceValueType.DataType(),
		values:   values,
	}, nil
}

func (p *Parser) evaluateSubscript(ctx context) (Expression, error) {
	nameToken := p.eat()

	if nameToken.Type() != lexer.IDENTIFIER {
		return nil, p.expectedError("identifier", nameToken)
	}
	name := nameToken.Value()
	variable, exists := ctx.findVariable(name, p.prefix)

	if !exists {
		return nil, p.atError(fmt.Sprintf("variable %s has not been defined", name), nameToken)
	}
	valueType := variable.ValueType()
	isSlice := valueType.IsSlice()

	if !isSlice && valueType.DataType() != DATA_TYPE_STRING {
		return nil, p.expectedError(fmt.Sprintf("slice or string but variable is of type %s", valueType.ToString()), nameToken)
	}
	nextToken := p.eat()

	if nextToken.Type() != lexer.OPENING_SQUARE_BRACKET {
		return nil, p.expectedError("\"[\"", nextToken)
	}
	nextToken = p.peek()
	index, err := p.evaluateExpression(ctx)

	if err != nil {
		return nil, err
	}
	indexValueType := index.ValueType()

	if !indexValueType.IsInt() {
		return nil, p.expectedError(fmt.Sprintf("%s as index but got %s", DATA_TYPE_INTEGER, indexValueType.ToString()), nextToken)
	}
	nextToken = p.eat()

	if nextToken.Type() != lexer.CLOSING_SQUARE_BRACKET {
		return nil, p.expectedError("\"]\"", nextToken)
	}

	if !isSlice {
		return StringSubscript{
			Variable: variable,
			index:    index,
		}, nil
	}
	return SliceEvaluation{
		Variable: variable,
		index:    index,
		dataType: variable.ValueType().DataType(),
	}, nil
}

func (p *Parser) evaluateSliceAssignment(ctx context) (Statement, error) {
	nameToken := p.eat()

	if nameToken.Type() != lexer.IDENTIFIER {
		return nil, p.expectedError("slice variable", nameToken)
	}
	name := nameToken.Value()
	variable, exists := ctx.findVariable(name, p.prefix)

	if !exists {
		return nil, p.atError(fmt.Sprintf("variable %s has not been defined", name), nameToken)
	}
	variableValueType := variable.ValueType()

	if !variableValueType.IsSlice() {
		return nil, p.expectedError(fmt.Sprintf("slice but variable is of type %s", variableValueType.ToString()), nameToken)
	}
	nextToken := p.eat()

	if nextToken.Type() != lexer.OPENING_SQUARE_BRACKET {
		return nil, p.expectedError("\"[\"", nextToken)
	}
	nextToken = p.peek()
	index, err := p.evaluateExpression(ctx)

	if err != nil {
		return nil, err
	}
	indexValueType := index.ValueType()

	if !indexValueType.IsInt() {
		return nil, p.expectedError(fmt.Sprintf("%s as index but got %s", DATA_TYPE_INTEGER, indexValueType.ToString()), nextToken)
	}
	nextToken = p.eat()

	if nextToken.Type() != lexer.CLOSING_SQUARE_BRACKET {
		return nil, p.expectedError("\"]\"", nextToken)
	}
	nextToken = p.eat()

	if nextToken.Type() != lexer.ASSIGN_OPERATOR {
		return nil, p.expectedError("\"=\"", nameToken)
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
	definedVariable, exists := ctx.findVariable(name, p.prefix)

	if !exists {
		return nil, p.atError(fmt.Sprintf("variable %s has not been defined", name), identifierToken)
	}
	valueType := definedVariable.ValueType()

	if !valueType.IsInt() {
		return nil, p.expectedError(fmt.Sprintf("%s but got %s", NewValueType(DATA_TYPE_INTEGER, false).ToString(), valueType.ToString()), identifierToken)
	}
	operationToken := p.eat()
	increment := true

	switch operationToken.Type() {
	case lexer.INCREMENT_OPERATOR:
		// Nothing to do.
	case lexer.DECREMENT_OPERATOR:
		increment = false
	default:
		return nil, p.expectedError("\"++\" or \"--\"", operationToken)
	}
	return incrementDecrementStatement(definedVariable, increment), nil
}

func (p *Parser) evaluateLen(ctx context) (Expression, error) {
	expr, err := p.evaluateBuiltInFunction(lexer.LEN, "len", 1, 1, ctx, func(keywordToken lexer.Token, expressions []Expression) (Statement, error) {
		expr := expressions[0]

		if !expr.ValueType().isSlice {
			return nil, p.expectedError("slice", keywordToken)
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
			return nil, p.atError(fmt.Sprintf("got %s as destination but %s as source", dstType.ToString(), srcType.ToString()), keywordToken)
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
