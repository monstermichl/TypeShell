package parser

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/monstermichl/typeshell/lexer"
)

type scope string

const (
	SCOPE_PROGRAM  scope = "program"
	SCOPE_FUNCTION scope = "function"
	SCOPE_IF       scope = "if"
	SCOPE_FOR      scope = "for"
	SCOPE_SWITCH   scope = "switch"
	SCOPE_CONST    scope = "const"
)

func scopesToString(scopes []scope) []string {
	strings := make([]string, len(scopes))

	for i, scope := range scopes {
		strings[i] = string(scope)
	}
	return strings
}

type initValue struct {
	nameToken  lexer.Token
	name       string
	valueToken lexer.Token
	value      Expression
}

type builtinArg struct {
	token lexer.Token
	expr  Expression
}

type context struct {
	imports     map[string]string               // Maps import aliases to file hashes.
	types       Importables[Type]               // Stores the declared types.
	namedValues Importables[NamedValue]         // Stores the variable/constant name to variable/constant relation.
	functions   Importables[FunctionDefinition] // Stores the function name to function relation.
	scopeStack  []scope                         // Stores the current scopes.
	layer       int
	iotaCounter int
}

func newContext(prefix string) context {
	c := context{
		imports:     map[string]string{},
		layer:       -1, // Init at -1 since program increases it right away.
		iotaCounter: 0,
	}

	// Add elementary types.
	c.addType(NewValueType(NewTypeBool(), false))
	c.addType(NewValueType(NewTypeInt(), false))
	c.addType(NewValueType(NewTypeString(), false))
	c.addType(NewValueType(NewTypeError(), false))

	return c
}

func (c *context) pushScope(scope scope) error {
	c.scopeStack = append(c.scopeStack, scope)

	// If new const scope is pushed, reset iota counter.
	if scope == SCOPE_CONST {
		c.iotaCounter = 0
	}
	return nil
}

func (c *context) popScope() error {
	i := len(c.scopeStack) - 1
	c.scopeStack = slices.Delete(c.scopeStack, i, i+1)
	return nil
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

func (c *context) incrementIota() {
	c.iotaCounter++
}

func (c context) addImport(alias string, hash string) error {
	c.imports[alias] = hash
	return nil
}

func (c *context) addType(valueType ValueType) error {
	t := valueType.Type()
	name := t.Name()
	_, exists := c.types.find(name, t.Prefix())

	if exists {
		return fmt.Errorf("type %s has already been defined", name) // TODO: Find out if this should be connected to the layer.
	}
	if valueType.IsSlice() {
		// TODO: Add support.
		return errors.New("slices are not allowed yet in type declarations")
	}
	c.types.add(t)
	return nil
}

func (c *context) addNamedValues(namedValues ...NamedValue) error {
	c.namedValues.add(namedValues...)
	return nil
}

func (c *context) addFunctions(functions ...FunctionDefinition) {
	c.functions.add(functions...)
}

func (c context) findImport(alias string) (string, bool) {
	hash, exists := c.imports[alias]
	return hash, exists
}

func (c context) findType(typeName string, prefix string) (Type, bool) {
	return c.types.find(typeName, prefix)
}

func (c context) findNamedValue(name string, prefix string) (NamedValue, bool) {
	return c.namedValues.find(name, prefix)
}

func (c context) findFunction(name string, prefix string) (FunctionDefinition, bool) {
	return c.functions.find(name, prefix)
}

func (c context) clone() context {
	return context{
		imports:     maps.Clone(c.imports),
		types:       c.types.clone(),
		namedValues: c.namedValues.clone(), // TODO: Make sure this is appropriate cloning because each entry contains a slice.
		functions:   c.functions.clone(),
		scopeStack:  slices.Clone(c.scopeStack),
		layer:       c.layer,
		iotaCounter: c.iotaCounter,
	}
}

type evaluatedImport struct {
	alias string
	path  string
}

type evaluatedValues struct {
	values []Expression
	tokens []lexer.Token
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

type parserError struct {
	Err   error
	Token lexer.Token
}

type Parser struct {
	tokens    []lexer.Token
	errors    []parserError
	index     int
	path      string
	prefix    string
	currFunc  *FunctionDefinition
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

func (p Parser) Errors() []parserError {
	return p.errors
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

		prefix := fmt.Sprintf("%x", h.Sum(nil))[0:7] // Only use the 7 first characters (inspired by Git).

		// If prefix starts with a number, prepend an "x" to make sure it starts with a letter.
		if regexp.MustCompile(`^\d`).MatchString(prefix) {
			prefix = fmt.Sprintf("x%s", prefix)
		}
		p.prefix = prefix
	}
	program, _ := p.evaluateProgram()

	if len(p.errors) > 0 {
		err = p.errors[0].Err
	}
	return program, err
}

// func incrementDecrementStatement(variable Variable, increment bool) Statement {
// 	operation := BINARY_OPERATOR_ADDITION

// 	if !increment {
// 		operation = BINARY_OPERATOR_SUBTRACTION
// 	}
// 	return VariableAssignmentValueAssignment{
// 		variables: []Variable{variable},
// 		values: []Expression{
// 			BinaryOperation{
// 				left: VariableEvaluation{
// 					Variable: variable,
// 				},
// 				operator: operation,
// 				right:    IntegerLiteral{value: 1},
// 			},
// 		},
// 	}
// }

func (p *Parser) atError(what string, token lexer.Token) error {
	err := fmt.Errorf("%s at row %d, column %d: %s", what, token.Row(), token.Column(), p.path)
	p.errors = append(p.errors, parserError{err, token})

	return err
}

func (p *Parser) expectedError(what string, token lexer.Token) error {
	return p.atError(fmt.Sprintf("expected %s", what), token)
}

func (p *Parser) expectedKeywordError(keyword string, token lexer.Token) error {
	return p.expectedError(fmt.Sprintf("%s-keyword", keyword), token)
}

func (p *Parser) expectedIdentifierError(token lexer.Token) error {
	return p.expectedError("identifier", token)
}

func (p *Parser) expectedNewlineError(token lexer.Token) error {
	return p.expectedError("newline", token)
}

func (p *Parser) constantError(constant string, token lexer.Token) error {
	return p.atError(fmt.Sprintf("cannot assign a value to constant %s", constant), token)
}

func (p *Parser) notDefinedError(what string, name string, token lexer.Token) error {
	return p.atError(fmt.Sprintf("%s %s has not been defined", what, name), token)
}

func (p *Parser) variableNotDefinedError(variable string, token lexer.Token) error {
	return p.notDefinedError("variable", variable, token)
}

func (p *Parser) typeNotDefinedError(t string, token lexer.Token) error {
	return p.notDefinedError("type", t, token)
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

func (p *Parser) skipUntil(t ...lexer.TokenType) {
	t = append(t, lexer.EOF) // Make sure skipping stops eventually.

	for !slices.Contains(t, p.peek().Type()) {
		p.eat()
	}
}

func (p *Parser) skipUntilNewline() {
	p.skipUntil(lexer.NEWLINE)
}

func (p *Parser) skipWhile(t ...lexer.TokenType) {
	for slices.Contains(t, p.peek().Type()) {
		p.eat()
	}
}

func (p *Parser) skipNewlines() {
	p.skipWhile(lexer.NEWLINE)
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

func (p *Parser) evaluateKeyword(keywords ...lexer.Keyword) (lexer.Token, bool) {
	found := false
	keywordToken := p.peek()

	for _, keyword := range keywords {
		if keywordToken.IsKeyword(keyword) {
			found = true
			break
		}
	}

	if found {
		p.eat()
	} else {
		for i, keyword := range keywords {
			keywords[i] = fmt.Sprintf(`"%s"`, keyword)
		}
		p.expectedKeywordError(strings.Join(keywords, " or "), keywordToken)
	}
	return keywordToken, found
}

func (p *Parser) evaluateExpressions(ctx context) ([]Expression, bool) {
	names := []Expression{}

	for {
		expr, ok := p.evaluateExpression(ctx)

		if !ok {
			return names, false
		}
		names = append(names, expr)
		nextToken := p.peek()

		if nextToken.Type() == lexer.COMMA {
			p.eat()
		} else {
			break
		}
	}
	return names, true
}

// func (p *Parser) evaluateValues(ctx context) (evaluatedValues, error) {
// }

func (p *Parser) evaluateProgram() (Program, bool) {
	ctx := newContext(p.prefix)
	statements, ok := p.evaluateBlockContent(ctx, SCOPE_PROGRAM)

	return Program{statements}, ok
}

// func (p *Parser) evaluateImports(ctx *context) ([]Statement, error) {

// }

// func (p *Parser) evaluateImport() (evaluatedImport, error) {

// }

func (p *Parser) evaluateBlockContent(ctx context, scope scope) ([]Statement, bool) {
	statements := []Statement{}
	ok := true

	p.skipNewlines()

	for !slices.Contains([]lexer.TokenType{lexer.EOF, lexer.CLOSING_CURLY_BRACKET}, p.peek().Type()) {
		stmt, okTemp := p.evaluateStatement(ctx)
		ok = ok && okTemp

		if !okTemp {
			p.skipUntilNewline() // RECOVER: Move to end of line.
		}
		if stmt != nil {
			statements = append(statements, stmt)
		}
		p.skipNewlines()
	}
	return statements, ok
}

func (p *Parser) evaluateBlock(ctx context, scope scope) (Block, bool) {
	openingBracketToken := p.peek()
	block := Block{}

	defer (func() {
		// RECOVER: Skip to closing curly bracket or next section keyword.
		p.skipUntil(lexer.CLOSING_CURLY_BRACKET, lexer.SECTION_KEYWORD)

		if p.peek().Type() == lexer.CLOSING_CURLY_BRACKET {
			p.eat()
		}
	})()

	if openingBracketToken.Type() != lexer.OPENING_CURLY_BRACKET {
		p.expectedError(`"{"`, openingBracketToken)
	} else {
		p.eat()
		block.OpeningBracket = &openingBracketToken
	}
	statements, ok := p.evaluateBlockContent(ctx, scope)
	block.Statements = statements
	closingBracketToken := p.peek()

	if closingBracketToken.Type() != lexer.CLOSING_CURLY_BRACKET {
		p.expectedError(`"}"`, closingBracketToken)
		ok = false
	} else {
		// Don't eat token as it's eaten by the defered function.
		block.ClosingBracket = &closingBracketToken
	}
	return block, ok
}

// func (p *Parser) evaluateValueType(ctx context, importAliases ...string) (ValueType, error) {

// }

// func (p *Parser) evaluateStructDefinition(name string, ctx context) (StructDefinition, error) {

// }

// func (p *Parser) evaluateTypeDeclaration(ctx context) (Statement, bool) {

// }

func (p *Parser) evaluateVarDefinition(ctx context) (Statement, bool) {
	keywordToken, ok := p.evaluateKeyword(lexer.KeywordVar, lexer.KeywordConst)
	varStatement := VariableDeclaration{Keyword: keywordToken}
	grouped := false
	nextToken := p.peek()
	grouped = nextToken.Type() == lexer.OPENING_ROUND_BRACKET

	if grouped {
		varStatement.OpeningBracket = &nextToken
		p.eat()
		p.skipNewlines()
	}
	skipUntilClosingOrNewline := func() {
		until := []lexer.TokenType{lexer.NEWLINE}

		if grouped {
			until = append(until, lexer.CLOSING_ROUND_BRACKET)
		}
		p.skipUntil(until...)
	}

	for {
		names, okTemp := p.evaluateExpressions(ctx)
		spec := ValueSpec{Names: names}

		if okTemp {
			nextToken := p.peek()

			if nextToken.Type() == lexer.IDENTIFIER {
				t, _ := p.evaluateExpression(ctx)
				spec.Type = t
			}
			nextToken = p.peek()

			if nextToken.Type() != lexer.ASSIGN_OPERATOR {
				p.expectedError(lexer.OperatorAssign, nextToken)
				skipUntilClosingOrNewline()
				okTemp = false
			} else {
				p.eat()
			}
		} else {
			skipUntilClosingOrNewline()
		}

		if okTemp {
			var values []Expression
			values, okTemp = p.evaluateExpressions(ctx)
			spec.Values = values
		} else {
			skipUntilClosingOrNewline()
		}
		ok = ok && okTemp
		varStatement.Specs = append(varStatement.Specs, spec)

		if !grouped {
			break
		}
		nextToken := p.peek()
		noError := false

		if nextToken.Type() == lexer.NEWLINE {
			p.skipNewlines()
			noError = true
		}
		escape := false
		nextToken = p.peek()

		if nextToken.Type() == lexer.CLOSING_ROUND_BRACKET {
			escape = true
		} else if !noError {
			p.expectedError(fmt.Sprintf(`newline or ")" but got "%s"`, nextToken.Value()), nextToken)
		}

		if escape {
			break
		}
	}

	if grouped {
		p.skipUntil(lexer.CLOSING_ROUND_BRACKET, lexer.SECTION_KEYWORD)
		nextToken = p.peek()

		if nextToken.Type() != lexer.CLOSING_ROUND_BRACKET {
			p.expectedError(`")"`, nextToken)
		} else {
			varStatement.ClosingBracket = &nextToken
			p.eat()
		}
	}
	return varStatement, ok
}

// func (p *Parser) evaluateCompoundAssignment(importAlias string, ctx context) (Statement, error) {

// }

// func (p *Parser) evaluateVarAssignment(importAlias string, ctx context) (Statement, error) {

// }

// func (p *Parser) evaluateParams(ctx context) ([]Param, error) {

// }

// func (p *Parser) evaluateFunctionDefinition(ctx context) (Statement, bool) {

// }

// func (p *Parser) evaluateReturn(ctx context) (Statement, bool) {

// }

// func (p *Parser) evaluateBreak(ctx context) (Statement, bool) {

// }

// func (p *Parser) evaluateIota(ctx context) (Expression, error) {

// }

// func (p *Parser) evaluateContinue(ctx context) (Statement, bool) {

// }

func (p *Parser) evaluateIf(ctx context) (Statement, bool) {
	ifStatement := If{}
	_, ok := p.evaluateKeyword(lexer.KeywordIf)
	nextToken := p.peek()

	if nextToken.Type() == lexer.OPENING_CURLY_BRACKET {
		p.expectedError("expression", nextToken)
		ok = false
	} else {
		expr, okTemp := p.evaluateExpression(ctx)
		ok = ok && okTemp
		ifStatement.Condition = expr

		if !ok {
			p.skipUntil(lexer.OPENING_CURLY_BRACKET, lexer.NEWLINE)
		}
	}
	block, okTemp := p.evaluateBlock(ctx, SCOPE_IF)
	ok = ok && okTemp
	ifStatement.Body = block
	nextToken = p.peek()

	if nextToken.IsKeyword(lexer.KeywordElse) {
		p.eat()
		nextToken := p.peek()
		var stmt Statement

		// If another if-keyword follows, evaluate if, otherwise evaluate else.
		if nextToken.IsKeyword(lexer.KeywordIf) {
			stmt, okTemp = p.evaluateIf(ctx)
		} else {
			stmt, okTemp = p.evaluateBlock(ctx, SCOPE_IF)
		}
		ok = ok && okTemp
		ifStatement.Else = stmt
	}
	return ifStatement, ok
}

// func (p *Parser) evaluateSwitch(ctx context) (Statement, bool) {

// }

// func (p *Parser) evaluateFor(ctx context) (Statement, bool) {

// }

// func (p *Parser) evaluateTypeDefinition(importAlias string, ctx context) (Expression, error) {

// }

// func (p *Parser) evaluateStructFields(importAlias string, stopOnLastStruct bool, ctx context) (Expression, StructField, error) {

// }

// func (p *Parser) evaluateStructFieldsFromExpression(importAlias string, structExpression Expression, structExpressionToken lexer.Token, stopOnLastStruct bool, ctx context) (Expression, StructField, error) {

// }

// func (p *Parser) evaluateStructEvaluation(importAlias string, ctx context) (Expression, error) {

// }

// func (p *Parser) evaluateVarEvaluation(importAlias string, ctx context) (Expression, error) {

// }

// func (p *Parser) evaluateImportAlias(ctx context) (string, lexer.Token, error) {

// }

func (p *Parser) evaluateSingleExpression(ctx context) (Expression, bool) {
	var expr Expression

	nextToken := p.peek()
	nextTokenType := nextToken.Type()
	value := nextToken.Value()

	switch nextTokenType {
	// String literal.
	case lexer.BOOL_LITERAL:
		p.eat()
		b, err := strconv.ParseBool(value)

		if err != nil {
			p.atError(err.Error(), nextToken)
		} else {
			expr = NewBooleanLiteral(b, nextToken)
		}
	// Number literal.
	case lexer.NUMBER_LITERAL:
		p.eat()
		i, err := strconv.Atoi(value)

		if err != nil {
			p.atError(err.Error(), nextToken)
		} else {
			expr = NewIntegerLiteral(i, nextToken)
		}
	// String literal.
	case lexer.STRING_LITERAL:
		p.eat()
		expr = NewStringLiteral(value, nextToken)
	// Identifier.
	case lexer.IDENTIFIER:
		p.eat()
		expr = NewIdentifier(value, nextToken)
	// Group.
	case lexer.OPENING_ROUND_BRACKET:
		var child Expression

		p.eat()
		openingBracket := nextToken
		child, okTemp := p.evaluateExpression(ctx)

		if okTemp {
			var closingBracket *lexer.Token
			nextToken := p.peek()

			if nextToken.Type() != lexer.CLOSING_ROUND_BRACKET {
				p.expectedError(`")"`, nextToken)
			} else {
				p.eat()
				closingBracket = &nextToken
			}
			expr = NewGroup(child, &openingBracket, closingBracket)
		}
	default:
		p.atError(fmt.Sprintf("unknown token type %d (%s)", nextTokenType, value), nextToken)
	}
	ok := expr != nil

	// If the expression has been evaluated correctly, see if it's followed by something usable.
	if ok {
		nextToken = p.peek()

		switch nextToken.Type() {
		case lexer.DOT:
			p.eat()
			selectorToken := p.peek()

			if selectorToken.Type() != lexer.IDENTIFIER {
				p.expectedIdentifierError(selectorToken)
				ok = false
			} else {
				p.eat()
				expr = NewSelector(expr, selectorToken.Value())
			}
		case lexer.OPENING_SQUARE_BRACKET:
			leftBracket := p.eat()
			indexExpr, okTemp := p.evaluateExpression(ctx)
			nextToken = p.peek()
			ok = ok && okTemp

			if nextToken.Type() != lexer.CLOSING_SQUARE_BRACKET {
				p.expectedError(`"]"`, nextToken)
				ok = false
			} else {
				p.eat()
				expr = NewIndex(expr, indexExpr, &leftBracket, &nextToken)
			}
		}
	}
	return expr, ok
}

// -----------------------------------------------------------------------------------------------
// This section defines the operator precedence. Call operators with higher precende first as
// in a function because higher precedence means it must be processed further down the chain.
// Learnt a lot about priority handling from this video https://www.youtube.com/watch?v=aAvL2BTHf60.
// Precedence is the same as in Go (https://go.dev/ref/spec#Operator_precedence).
func (p *Parser) evaluateUnaryOperation(ctx context) (Expression, bool) {
	unaryOperatorToken := p.peek()
	negate := false
	ok := true

	if unaryOperatorToken.Type() == lexer.UNARY_OPERATOR {
		p.eat()
		value := unaryOperatorToken.Value()

		switch value {
		case lexer.UnaryOperatorNegate:
			negate = true
		default:
			p.atError(fmt.Sprintf("unknown unary operator %s", value), unaryOperatorToken)
			ok = false
		}
	}
	expr, okTemp := p.evaluateSingleExpression(ctx)
	ok = ok && okTemp

	if negate {
		expr = NewUnaryOperation(expr, unaryOperatorToken)
	}
	return expr, ok
}

func (p *Parser) evaluateMultiplication(ctx context) (Expression, bool) {
	return p.evaluateBinaryOperation(ctx, []BinaryOperator{lexer.OperatorMultiplication}, p.evaluateUnaryOperation)
}

func (p *Parser) evaluateAddition(ctx context) (Expression, bool) {
	return p.evaluateBinaryOperation(ctx, []BinaryOperator{lexer.OperatorAddition}, p.evaluateMultiplication)
}

func (p *Parser) evaluateComparison(ctx context) (Expression, bool) {
	return p.evaluateBinaryOperation(ctx, []BinaryOperator{lexer.ComparisonLess, lexer.ComparisonLessOrEqual, lexer.ComparisonGreater, lexer.ComparisonGreaterOrEqual}, p.evaluateAddition)
}

func (p *Parser) evaluateEqualityComparison(ctx context) (Expression, bool) {
	return p.evaluateBinaryOperation(ctx, []BinaryOperator{lexer.ComparisonEqual, lexer.ComparisonNotEqual}, p.evaluateComparison)
}

func (p *Parser) evaluateLogicalAnd(ctx context) (Expression, bool) {
	return p.evaluateBinaryOperation(ctx, []BinaryOperator{lexer.OperatorAnd}, p.evaluateEqualityComparison)
}

func (p *Parser) evaluateLogicalOr(ctx context) (Expression, bool) {
	return p.evaluateBinaryOperation(ctx, []BinaryOperator{lexer.OperatorOr}, p.evaluateLogicalAnd)
}

func (p *Parser) evaluateExpression(ctx context) (Expression, bool) {
	return p.evaluateLogicalOr(ctx)
}

func (p *Parser) evaluateStatement(ctx context) (Statement, bool) {
	var stmt Statement
	var ok bool

	token := p.peek()

	switch token.Type() {
	case lexer.KEYWORD, lexer.SECTION_KEYWORD:
		switch token.Value() {
		// case lexer.KeywordType:
		// 	stmt, ok = p.evaluateTypeDeclaration(ctx)
		case lexer.KeywordVar, lexer.KeywordConst:
			stmt, ok = p.evaluateVarDefinition(ctx)
			fmt.Println("----xx", ok)
		// case lexer.KeywordFunc:
		// 	stmt, ok = p.evaluateFunctionDefinition(ctx)
		// case lexer.KeywordReturn:
		// 	stmt, ok = p.evaluateReturn(ctx)
		case lexer.KeywordIf:
			stmt, ok = p.evaluateIf(ctx)
			// case lexer.KeywordSwitch:
			// 	stmt, ok = p.evaluateSwitch(ctx)
			// case lexer.KeywordFor:
			// 	stmt, ok = p.evaluateFor(ctx)
			// case lexer.KeywordBreak:
			// 	stmt, ok = p.evaluateBreak(ctx)
			// case lexer.KeywordContinue:
			// 	stmt, ok = p.evaluateContinue(ctx)
			// TODO: Handle in type checker or somewhere else.
			// case lexer.KeywordPrint:
			// 	stmt, err = p.evaluatePrint(ctx)
			// case lexer.KeywordWrite:
			// 	stmt, err = p.evaluateWrite(ctx)
			// case lexer.KeywordPanic:
			// 	stmt, err = p.evaluatePanic(ctx)
			// case lexer.KeywordUnsafe:
			// 	stmt, err = p.evaluateUnsafe(ctx)
		}
	default:
		// Assume it's an expression.
		stmt, ok = p.evaluateExpression(ctx)
	}
	return stmt, ok
}

// -----------------------------------------------------------------------------------------------

func (p *Parser) evaluateBinaryOperation(ctx context, allowedOperators []BinaryOperator, higherPrioOperation func(ctx context) (Expression, bool)) (Expression, bool) {
	leftExpression, ok := higherPrioOperation(ctx)

	if !ok {
		return leftExpression, false
	}

	for {
		var rightExpression Expression

		operatorToken := p.peek()
		operator := ""

		if operatorToken.Type() == lexer.NUMBER_LITERAL {
			rightExpression, ok = p.evaluateSingleExpression(ctx)

			if ok {
				numberLiteral := rightExpression.(IntegerLiteral)

				if numberLiteral.Value < 0 {
					operator = "+"
				}
			}
		} else {
			operator = operatorToken.Value()
		}

		// If operator found, process it.
		if !slices.Contains(allowedOperators, operator) {
			break
		}
		p.eat() // Eat operator token.

		rightExpression, ok = higherPrioOperation(ctx)
		leftExpression = NewBinaryOperation(leftExpression, operatorToken, rightExpression)
	}
	return leftExpression, ok
}

// func (p *Parser) evaluateArguments(typeName string, name string, params []Param, receiver Expression, ctx context) ([]Expression, []Expression, error) {

// }

// func (p *Parser) evaluateFunctionCall(importAlias string, receiver Expression, ctx context) (Call, error) {

// }

// func (p *Parser) evaluateAppCall(ctx context) (Call, error) {

// }

// func (p *Parser) evaluateInitializationValues(checkCallout func(initValue initValue) error, ctx context) ([]initValue, error) {
// }

// func (p *Parser) evaluateSliceInitialization(ctx context) (Expression, error) {
// }

// func (p *Parser) evaluateStructInitialization(importAlias string, ctx context) (Expression, error) {

// }

// func (p *Parser) evaluateChaining(expr Expression, ctx context) (Expression, bool, error) {
// }

// func (p *Parser) evaluateSubscript(importAlias string, ctx context) (Expression, error) {
// }

// func (p *Parser) evaluateSubscriptFromExpression(value Expression, valueToken lexer.Token, ctx context) (Expression, error) {
// }

// func (p *Parser) evaluateSliceAssignment(importAlias string, ctx context) (Statement, error) {
// }

// func (p *Parser) evaluateStructAssignment(importAlias string, ctx context) (Statement, error) {
// }

// func (p *Parser) evaluateIncrementDecrement(importAlias string, ctx context) (Statement, error) {
// }
