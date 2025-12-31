package parser

import (
	"crypto/sha256"
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
	imports     map[string]string // Maps import aliases to file hashes.
	scopeStack  []scope           // Stores the current scopes.
	layer       int
	iotaCounter int
}

func newContext(prefix string) context {
	c := context{
		imports:     map[string]string{},
		layer:       -1, // Init at -1 since program increases it right away.
		iotaCounter: 0,
	}
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

func (c context) findImport(alias string) (string, bool) {
	hash, exists := c.imports[alias]
	return hash, exists
}

func (c context) clone() context {
	return context{
		imports:     maps.Clone(c.imports),
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

type blockCheckCallout func(stmt Statement) bool

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
	err := fmt.Errorf("%s at row %d, column %d: %s", what, token.Row, token.Column, p.path)
	p.errors = append(p.errors, parserError{err, token})

	return err
}

func (p *Parser) expectedError(what string, token lexer.Token) error {
	return p.atError(fmt.Sprintf("expected %s", what), token)
}

func (p *Parser) expectedKeywordError(keyword string, token lexer.Token) error {
	return p.expectedError(fmt.Sprintf("%s-keyword", keyword), token)
}

func (p *Parser) expectedAssignOperatorError(token lexer.Token) error {
	return p.expectedError(fmt.Sprintf(`"%s" or "%s"`, lexer.OperatorAssign, lexer.OperatorShortAssign), token)
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

	for !slices.Contains(t, p.peek().Type) {
		p.eat()
	}
}

func (p *Parser) skipUntilNewline() {
	p.skipUntil(lexer.NEWLINE)
}

func (p *Parser) skipWhile(t ...lexer.TokenType) {
	for slices.Contains(t, p.peek().Type) {
		p.eat()
	}
}

func (p *Parser) skipNewlines() {
	p.skipWhile(lexer.NEWLINE)
}

func (p *Parser) defaultBlockCheckCallout(stmt Statement) bool {
	// if stmt.StatementType == STATEMENT_TYPE_CASE_CLAUSE {
	// 	p.atError("case clause is not permitted in this context", stmt.Token())
	// 	return false
	// }
	return true
}

func (p Parser) findAllowed(searchTokenType lexer.TokenType, allowed ...lexer.TokenType) (lexer.Token, error) {
	tokens := p.tokens

	for i := p.index; i < len(tokens); i++ {
		token := tokens[i]
		tokenType := token.Type

		if tokenType == searchTokenType {
			return token, nil
		}

		if !slices.Contains(allowed, tokenType) {
			return lexer.Token{}, fmt.Errorf(`found illegal token "%d" before "%d"`, tokenType, searchTokenType)
		}
	}
	return lexer.Token{}, fmt.Errorf(`token type "%d" not found`, searchTokenType)
}

func (p Parser) findBefore(searchTokenType lexer.TokenType, before ...lexer.TokenType) (lexer.Token, bool) {
	tokens := p.tokens
	before = append(before, lexer.EOF)

	for i := p.index; i < len(tokens); i++ {
		token := tokens[i]
		tokenType := token.Type

		if tokenType == searchTokenType {
			return token, true
		} else if slices.Contains(before, tokenType) {
			return lexer.Token{}, false
		}
	}
	return lexer.Token{}, false
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
	expressions := []Expression{}

	for {
		expr, ok := p.evaluateExpression(ctx)

		if !ok {
			return expressions, false
		}
		expressions = append(expressions, expr)
		nextToken := p.peek()

		if nextToken.Type == lexer.COMMA {
			p.eat()
		} else {
			break
		}
	}
	return expressions, true
}

// func (p *Parser) evaluateValues(ctx context) (evaluatedValues, error) {
// }

func (p *Parser) evaluateProgram() (Program, bool) {
	ctx := newContext(p.prefix)
	statements, ok := p.evaluateBlockContent(ctx, SCOPE_PROGRAM, nil)

	return Program{statements}, ok
}

// func (p *Parser) evaluateImports(ctx *context) ([]Statement, error) {

// }

// func (p *Parser) evaluateImport() (evaluatedImport, error) {

// }

func (p *Parser) evaluateBlockContent(ctx context, scope scope, checkCallout blockCheckCallout, stopKeywords ...lexer.Keyword) ([]Statement, bool) {
	statements := []Statement{}
	ok := true

	p.skipNewlines()

	if checkCallout == nil {
		checkCallout = p.defaultBlockCheckCallout
	}

	for {
		nextToken := p.peek()
		nextTokenType := nextToken.Type

		if slices.Contains([]lexer.TokenType{lexer.EOF, lexer.CLOSING_CURLY_BRACKET}, nextTokenType) || (nextTokenType == lexer.KEYWORD && slices.Contains(stopKeywords, nextToken.Value)) {
			break
		}
		stmt, okTemp := p.evaluateStatement(ctx)
		ok = ok && okTemp

		if !okTemp {
			p.skipUntilNewline() // RECOVER: Move to end of line.
		} else {
			ok = ok && checkCallout(stmt)
		}

		if stmt != nil {
			statements = append(statements, stmt)
		}
		p.skipNewlines()
	}
	return statements, ok
}

func (p *Parser) evaluateBlock(ctx context, scope scope, checkCallout blockCheckCallout, stopKeywords ...lexer.Keyword) (Block, bool) {
	openingBracketToken := p.peek()
	block := Block{token: openingBracketToken}

	defer (func() {
		// RECOVER: Skip to closing curly bracket or next section keyword.
		p.skipUntil(lexer.CLOSING_CURLY_BRACKET, lexer.SECTION_KEYWORD)

		if p.peek().Type == lexer.CLOSING_CURLY_BRACKET {
			p.eat()
		}
	})()

	if openingBracketToken.Type != lexer.OPENING_CURLY_BRACKET {
		p.expectedError(`"{"`, openingBracketToken)
	} else {
		p.eat()
		block.OpeningBracket = &openingBracketToken
	}
	statements, ok := p.evaluateBlockContent(ctx, scope, checkCallout, stopKeywords...)
	block.Statements = statements
	closingBracketToken := p.peek()

	if closingBracketToken.Type != lexer.CLOSING_CURLY_BRACKET {
		p.expectedError(`"}"`, closingBracketToken)
		ok = false
	} else {
		// Don't eat token as it's eaten by the defered function.
		block.ClosingBracket = &closingBracketToken
	}
	return block, ok
}

func (p *Parser) evaluateType(ctx context) (Type, bool) {
	ok := true
	slice := false
	nextToken := p.peek()

	// Evaluate if value type is a slice type.
	if nextToken.Type == lexer.OPENING_SQUARE_BRACKET {
		p.eat() // Eat opening square bracket.
		nextToken = p.peek()

		if nextToken.Type != lexer.CLOSING_SQUARE_BRACKET {
			p.expectedError(`"]"`, nextToken)
			ok = false
		} else {
			p.eat()
		}
		nextToken = p.peek()
		slice = true
	}
	var t Type

	// Evaluate data type.
	switch nextToken.Type {
	case lexer.IDENTIFIER:
		t, _ = p.evaluateIdentifier(ctx)
	case lexer.KEYWORD:
		if !nextToken.IsKeyword(lexer.KeywordStruct) {
			p.expectedKeywordError(lexer.KeywordStruct, nextToken)
			ok = false
		} else {
			t, ok = p.evaluateStructDeclaration(ctx)
		}
	default:
		// If still ok, now is the time to invalidate the ok-state.
		if ok {
			p.expectedIdentifierError(nextToken)
			ok = false
		}
	}

	if slice {
		t = SliceType{Type: t}
	}
	return t, ok
}

func (p *Parser) evaluateStructDeclaration(ctx context) (StructDeclaration, bool) {
	declaration := StructDeclaration{}
	keywordToken, ok := p.evaluateKeyword(lexer.KeywordStruct)

	if !ok {
		p.expectedKeywordError(lexer.KeywordStruct, keywordToken)
	} else {
		declaration.keywordToken = keywordToken
	}
	hasOpeningBracket := false

	if ok {
		openingBracketToken := p.peek()
		hasOpeningBracket = openingBracketToken.Type == lexer.OPENING_CURLY_BRACKET

		if !hasOpeningBracket {
			p.expectedError(`"{`, openingBracketToken)
			ok = false
		} else {
			p.eat()
		}
	}

	if !ok {
		p.skipUntil(lexer.OPENING_CURLY_BRACKET, lexer.NEWLINE)

		if p.peek().Type == lexer.OPENING_CURLY_BRACKET {
			hasOpeningBracket = true
		}
	}

	if hasOpeningBracket {
		for {
			p.skipNewlines()

			names, okTemp := p.evaluateExpressions(ctx)
			field := StructField{}

			if okTemp {
				field.Names = names
				field.Type, okTemp = p.evaluateType(ctx)
			}
			ok = ok && okTemp

			if !okTemp {
				p.skipUntil(lexer.CLOSING_CURLY_BRACKET, lexer.NEWLINE)
			} else {
				declaration.Fields = append(declaration.Fields, field)
			}
			p.skipNewlines()

			if slices.Contains([]lexer.TokenType{lexer.CLOSING_CURLY_BRACKET, lexer.SECTION_KEYWORD}, p.peek().Type) {
				break
			}
		}
	}

	if !ok {
		// If opening curly bracket exists, look for closing bracket or section keyword. Otherwise, just look for newline.
		if hasOpeningBracket {
			p.skipUntil(lexer.CLOSING_ROUND_BRACKET, lexer.SECTION_KEYWORD)
		} else {
			p.skipUntil(lexer.NEWLINE)
		}
	}
	closingBracketToken := p.peek()

	if closingBracketToken.Type == lexer.CLOSING_CURLY_BRACKET {
		p.eat()
	} else if ok {
		p.expectedError(`"}"`, closingBracketToken)
	}
	return declaration, ok
}

func (p *Parser) evaluateTypeDeclaration(ctx context) (Statement, bool) {
	keywordToken, ok := p.evaluateKeyword(lexer.KeywordType)
	declaration := TypeDeclaration{}

	if !ok {
		p.expectedKeywordError(lexer.KeywordType, keywordToken)
	} else {
		declaration.token = keywordToken
	}
	nameToken := p.peek()

	if ok {
		// Make sure, name token is an identifier.
		if nameToken.Type != lexer.IDENTIFIER {
			p.expectedIdentifierError(nameToken)
			ok = false
		} else {
			nameExpr, okTemp := p.evaluateIdentifier(ctx)
			ok = ok && okTemp

			if okTemp {
				declaration.Name = nameExpr
			}
		}
	}

	if ok {
		assignOpToken := p.peek()

		if assignOpToken.Type == lexer.ASSIGN_OPERATOR && assignOpToken.Value == lexer.OperatorAssign {
			p.eat()
			declaration.AssignOpToken = &assignOpToken
		}
	}

	if ok {
		declaration.Type, ok = p.evaluateType(ctx)
	}

	if !ok {
		p.skipUntil(lexer.NEWLINE)
	}
	return declaration, ok
}

func (p *Parser) evaluateVarDefinition(ctx context) (Statement, bool) {
	keywordToken, ok := p.evaluateKeyword(lexer.KeywordVar, lexer.KeywordConst)
	varStatement := Declaration{Keyword: keywordToken}
	grouped := false
	nextToken := p.peek()
	grouped = nextToken.Type == lexer.OPENING_ROUND_BRACKET

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

			if slices.Contains([]lexer.TokenType{lexer.OPENING_SQUARE_BRACKET, lexer.IDENTIFIER}, nextToken.Type) {
				spec.Type, okTemp = p.evaluateType(ctx)
			}
			nextToken = p.peek()

			if !okTemp || nextToken.Type != lexer.ASSIGN_OPERATOR {
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

		if nextToken.Type == lexer.NEWLINE {
			p.skipNewlines()
			noError = true
		}
		escape := false
		nextToken = p.peek()

		if nextToken.Type == lexer.CLOSING_ROUND_BRACKET {
			escape = true
		} else if !noError {
			p.expectedError(fmt.Sprintf(`newline or ")" but got "%s"`, nextToken.Value), nextToken)
		}

		if escape {
			break
		}
	}

	if grouped {
		p.skipUntil(lexer.CLOSING_ROUND_BRACKET, lexer.SECTION_KEYWORD)
		nextToken = p.peek()

		if nextToken.Type != lexer.CLOSING_ROUND_BRACKET {
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

func (p *Parser) evaluateVarAssignment(ctx context, params ...string) (Statement, bool) {
	names, ok := p.evaluateExpressions(ctx)
	okTemp := ok
	assignment := Assignment{Left: names}

	if !okTemp {
		p.skipUntil(lexer.ASSIGN_OPERATOR, lexer.NEWLINE)

		// If an assign operator was found, reset okTemp.
		if p.peek().Type == lexer.ASSIGN_OPERATOR {
			okTemp = true
		}
	}

	if okTemp {
		operatorToken := p.peek()
		assignment.OperatorToken = operatorToken

		if operatorToken.Type != lexer.ASSIGN_OPERATOR {
			p.expectedAssignOperatorError(operatorToken)
			okTemp = false
		} else {
			p.eat()
			assignment.Right, okTemp = p.evaluateExpressions(ctx)
		}
	}
	ok = ok && okTemp

	if !ok {
		p.skipUntilNewline()
	}
	return assignment, ok
}

func (p *Parser) evaluateFieldList(bracketsRequired bool, ctx context) (FieldList, bool) {
	openingBracketToken := p.peek()
	params := FieldList{token: openingBracketToken}
	ok := true
	oneType := false

	if openingBracketToken.Type != lexer.OPENING_ROUND_BRACKET {
		if bracketsRequired {
			p.expectedError(`"("`, openingBracketToken)
			ok = false
		}
		oneType = true // If no brackets exist, it's only allowed to specify one type.
	} else {
		p.eat()
		params.OpeningBracket = &openingBracketToken
	}
	onlyTypes := oneType
	var nextToken lexer.Token

	for {
		nextToken := p.peek()

		if nextToken.Type != lexer.IDENTIFIER {
			break
		} else {
			p.eat()
		}
		var name Expression
		var paramType Expression

		okTemp := true
		expr, _ := p.evaluateExpression(ctx)

		if onlyTypes {
			paramType = expr
		} else {
			nextToken := p.peek()

			if nextToken.Type != lexer.IDENTIFIER {
				paramType = expr
				onlyTypes = true
			} else {
				name = expr
				paramType, okTemp = p.evaluateExpression(ctx)
			}
		}
		ok = ok && okTemp
		param := Field{token: nextToken, Name: name, Type: paramType}
		params.Fields = append(params.Fields, param)
		nextToken = p.peek()

		if !okTemp || nextToken.Type != lexer.COMMA {
			break
		}
	}

	if nextToken.Type != lexer.CLOSING_ROUND_BRACKET {
		_, closingBracketExists := p.findBefore(lexer.CLOSING_ROUND_BRACKET, lexer.IDENTIFIER, lexer.OPENING_ROUND_BRACKET, lexer.OPENING_CURLY_BRACKET, lexer.NEWLINE)

		if closingBracketExists {
			p.skipUntil(lexer.CLOSING_ROUND_BRACKET)
		}
		nextToken = p.peek()
	}

	if nextToken.Type == lexer.CLOSING_ROUND_BRACKET {
		p.eat()
		params.ClosingBracket = &nextToken
	}
	return params, ok
}

func (p *Parser) evaluateFunctionDefinition(ctx context) (Statement, bool) {
	keywordToken, ok := p.evaluateKeyword(lexer.KeywordFunc)
	funcDef := FunctionDefinition{token: keywordToken}
	nextToken := p.peek()

	// Check for receiver.
	if nextToken.Type == lexer.OPENING_ROUND_BRACKET {
		var receiverParams FieldList

		receiverParams, ok = p.evaluateFieldList(true, ctx)
		params := receiverParams.Fields
		paramsLength := len(params)

		if paramsLength > 0 {
			funcDef.Receiver = &params[0]
		}

		if paramsLength > 0 {
			p.atError(fmt.Sprintf("only one receiver is permitted but got %d", paramsLength), params[paramsLength-1].Token())
			ok = false
		}
	}

	if ok {
		nameToken := p.peek()

		if nameToken.Type != lexer.IDENTIFIER {
			p.expectedIdentifierError(nameToken)
			ok = false
		} else {
			funcDef.Name, _ = p.evaluateExpression(ctx)
		}
	}

	if ok {
		funcDef.Params, ok = p.evaluateFieldList(true, ctx)
	}

	if ok {
		nextToken := p.peek()

		if slices.Contains([]lexer.TokenType{lexer.IDENTIFIER, lexer.OPENING_ROUND_BRACKET}, nextToken.Type) {

		}
	}

	if !ok {
		p.skipUntil(lexer.OPENING_CURLY_BRACKET, lexer.NEWLINE)
	}
	return funcDef, ok
}

// func (p *Parser) evaluateReturn(ctx context) (Statement, bool) {

// }

// func (p *Parser) evaluateBreak(ctx context) (Statement, bool) {

// }

// func (p *Parser) evaluateIota(ctx context) (Expression, error) {

// }

// func (p *Parser) evaluateContinue(ctx context) (Statement, bool) {

// }

func (p *Parser) evaluateIf(ctx context) (Statement, bool) {
	keywordToken, ok := p.evaluateKeyword(lexer.KeywordIf)
	ifStatement := If{token: keywordToken}
	nextToken := p.peek()

	if nextToken.Type == lexer.OPENING_CURLY_BRACKET {
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
	block, okTemp := p.evaluateBlock(ctx, SCOPE_IF, nil)
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
			stmt, okTemp = p.evaluateBlock(ctx, SCOPE_IF, nil)
		}
		ok = ok && okTemp
		ifStatement.Else = stmt
	}
	return ifStatement, ok
}

func (p *Parser) evaluateSwitch(ctx context) (Statement, bool) {
	keywordToken, ok := p.evaluateKeyword(lexer.KeywordSwitch)
	switchStatement := Switch{token: keywordToken}
	okTemp := true
	nextToken := p.peek()

	var tag Expression

	if nextToken.Type == lexer.OPENING_CURLY_BRACKET {
		tag = NewBooleanLiteral(true, nextToken)
	} else {
		tag, okTemp = p.evaluateExpression(ctx)

		if !okTemp {
			p.skipUntil(lexer.OPENING_CURLY_BRACKET, lexer.NEWLINE)
		}
	}
	ok = ok && okTemp

	switchStatement.Tag = tag
	body, okTemp := p.evaluateBlock(ctx, SCOPE_SWITCH, func(stmt Statement) bool {
		if _, isCaseClause := stmt.(CaseClause); isCaseClause {
			p.expectedError("case clause", stmt.Token())
			return false
		}
		return true
	})
	switchStatement.Body = body
	defaultCases := []CaseClause{}

	for _, stmt := range body.Statements {
		if clause, ok := stmt.(CaseClause); ok {
			defaultCases = append(defaultCases, clause)
		}
	}
	defaultCasesLength := len(defaultCases)

	if defaultCasesLength > 1 {
		for i := 1; i < defaultCasesLength; i++ {
			p.atError("only one default case is permitted", defaultCases[i].Token())
		}
	}
	return switchStatement, ok && okTemp
}

func (p *Parser) evaluateCaseClause(ctx context) (Statement, bool) {
	keywordToken, ok := p.evaluateKeyword(lexer.KeywordCase, lexer.KeywordDefault)
	defaultClause := keywordToken.Type == lexer.KEYWORD && keywordToken.Value == lexer.KeywordDefault
	clause := CaseClause{token: keywordToken}
	list := []Expression{}
	okTemp := true

	if !defaultClause {
		list, okTemp = p.evaluateExpressions(ctx)

		if !okTemp {
			p.skipUntil(lexer.COLON, lexer.NEWLINE)
		}
	}
	ok = ok && okTemp
	clause.List = list
	nextToken := p.peek()

	if nextToken.Type != lexer.COLON {
		p.expectedError(`":"`, nextToken)
		ok = false
	} else {
		p.eat()
	}
	nextToken = p.peek()
	statements, okTemp := p.evaluateBlockContent(ctx, SCOPE_SWITCH, nil, lexer.KeywordCase, lexer.KeywordDefault)

	clause.Body = Block{
		token:      nextToken,
		Statements: statements,
	}
	return clause, ok && okTemp
}

func (p *Parser) evaluateFor(ctx context) (Statement, bool) {
	var stmt Statement

	keywordToken, ok := p.evaluateKeyword(lexer.KeywordFor)
	isRange := false

	// Try to find range-keyword.
	for i := 0; ; i++ {
		nextToken := p.peekAt(uint(i))
		leave := false

		switch nextToken.Type {
		case lexer.KEYWORD:
			if nextToken.Value == lexer.KeywordRange {
				isRange = true
				leave = true
			}
		case lexer.NEWLINE, lexer.EOF:
			leave = true
		}

		if leave {
			break
		}
	}
	skipUntilBlock := func() {
		p.skipUntil(lexer.OPENING_CURLY_BRACKET, lexer.NEWLINE)
	}

	if isRange {
		identifiers, okTemp := p.evaluateExpressions(ctx)
		rangeStatement := ForRange{token: keywordToken}

		if !okTemp {
			p.skipUntil(lexer.KEYWORD)
		} else {
			lenNames := len(identifiers)
			rangeStatement.Key = identifiers[0]

			if lenNames > 1 {
				rangeStatement.Value = identifiers[1]
			}

			if lenNames > 2 {
				p.expectedError("at most 2 identifiers", identifiers[2].Token())
				okTemp = false
			}
		}
		ok = ok && okTemp
		assignOperatorToken := p.peek()

		if assignOperatorToken.Type != lexer.ASSIGN_OPERATOR {
			p.expectedAssignOperatorError(assignOperatorToken)
			okTemp = false
		} else {
			p.eat() // Eat assign-operator token.
			rangeStatement.AssignOpToken = assignOperatorToken
		}

		if okTemp {
			rangeToken := p.peek()

			if rangeToken.Value != lexer.KeywordRange {
				p.expectedKeywordError(lexer.KeywordRange, rangeToken)
				okTemp = false
			} else {
				p.eat() // Eat range token.
				rangeStatement.X, okTemp = p.evaluateExpression(ctx)
			}
		}

		if !okTemp {
			skipUntilBlock()
		}
		rangeStatement.Body, okTemp = p.evaluateBlock(ctx, SCOPE_FOR, nil)
		ok = ok && okTemp
		stmt = rangeStatement
	} else {
		var condition Expression
		forStatement := For{token: keywordToken}
		_, foundSemicolon := p.findBefore(lexer.SEMICOLON, lexer.OPENING_CURLY_BRACKET, lexer.NEWLINE)
		okTemp := true

		if foundSemicolon {
			nextToken := p.peek()

			if nextToken.Type != lexer.SEMICOLON {
				forStatement.Init, okTemp = p.evaluateStatement(ctx)

				if !okTemp {
					p.skipUntil(lexer.SEMICOLON) // Skip until semicolon. This works for sure because otherwise foundSemicolon would not be true.
				}
			}
			ok = ok && okTemp
			okTemp = true // Reset.
			p.eat()       // Eat semicolon.
			nextToken = p.peek()

			if nextToken.Type == lexer.SEMICOLON {
				condition = NewBooleanLiteral(true, nextToken)
			} else {
				condition, okTemp = p.evaluateExpression(ctx)
			}

			if okTemp {
				nextToken = p.peek()
				okTemp = nextToken.Type == lexer.SEMICOLON

				if !okTemp {
					p.expectedError(`";"`, nextToken)
				} else {
					p.eat() // Eat semicolon.
					nextToken = p.peek()

					if !slices.Contains([]lexer.TokenType{lexer.OPENING_CURLY_BRACKET, lexer.NEWLINE}, nextToken.Type) {
						forStatement.Post, okTemp = p.evaluateExpression(ctx)
					}
				}
			}
			ok = ok && okTemp
		} else {
			nextToken := p.peek()

			if nextToken.Type == lexer.OPENING_CURLY_BRACKET {
				condition = NewBooleanLiteral(true, nextToken)
			} else {
				condition, okTemp = p.evaluateExpression(ctx)
			}
		}
		ok = ok && okTemp
		forStatement.Condition = condition

		if !okTemp {
			skipUntilBlock()
		}
		forStatement.Body, okTemp = p.evaluateBlock(ctx, SCOPE_FOR, nil)
		ok = ok && okTemp
		stmt = forStatement
	}
	return stmt, ok
}

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

func (p *Parser) evaluateIdentifier(ctx context) (Identifier, bool) {
	identifierToken := p.peek()
	identifier := Identifier{}
	ok := false

	if identifierToken.Type == lexer.IDENTIFIER {
		p.eat()

		identifier.token = identifierToken
		identifier.Name = identifierToken.Value

		ok = true
	}
	return identifier, ok
}

func (p *Parser) evaluateSingleExpression(ctx context) (Expression, bool) {
	var expr Expression = nil

	nextToken := p.peek()
	nextTokenType := nextToken.Type
	value := nextToken.Value

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
		expr, _ = p.evaluateIdentifier(ctx)
	// Group.
	case lexer.OPENING_ROUND_BRACKET:
		var child Expression

		p.eat()
		openingBracket := nextToken
		child, okTemp := p.evaluateExpression(ctx)

		if okTemp {
			var closingBracket *lexer.Token
			nextToken := p.peek()

			if nextToken.Type != lexer.CLOSING_ROUND_BRACKET {
				p.expectedError(`")"`, nextToken)
			} else {
				p.eat()
				closingBracket = &nextToken
			}
			expr = NewGroup(child, openingBracket, closingBracket)
		}
	default:
		p.atError(fmt.Sprintf("unknown expression token type %d (%s)", nextTokenType, value), nextToken)
	}
	ok := expr != nil

	// If the expression has been evaluated correctly, see if it's followed by something usable.
	if ok {
		nextToken = p.peek()

		switch nextToken.Type {
		case lexer.DOT:
			p.eat()
			selectorToken := p.peek()

			if selectorToken.Type != lexer.IDENTIFIER {
				p.expectedIdentifierError(selectorToken)
				ok = false
			} else {
				p.eat()
				expr = NewSelector(expr, selectorToken.Value)
			}
		case lexer.OPENING_SQUARE_BRACKET:
			leftBracket := p.eat()
			indexExpr, okTemp := p.evaluateExpression(ctx)
			nextToken = p.peek()
			ok = ok && okTemp

			if nextToken.Type != lexer.CLOSING_SQUARE_BRACKET {
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

	if unaryOperatorToken.Type == lexer.UNARY_OPERATOR {
		p.eat()
		value := unaryOperatorToken.Value

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
	value := token.Value
	tokenType := token.Type

	switch tokenType {
	case lexer.KEYWORD, lexer.SECTION_KEYWORD:
		switch value {
		case lexer.KeywordType:
			stmt, ok = p.evaluateTypeDeclaration(ctx)
		case lexer.KeywordVar, lexer.KeywordConst:
			stmt, ok = p.evaluateVarDefinition(ctx)
		case lexer.KeywordFunc:
			stmt, ok = p.evaluateFunctionDefinition(ctx)
		// case lexer.KeywordReturn:
		// 	stmt, ok = p.evaluateReturn(ctx)
		case lexer.KeywordIf:
			stmt, ok = p.evaluateIf(ctx)
		case lexer.KeywordSwitch:
			stmt, ok = p.evaluateSwitch(ctx)
		case lexer.KeywordCase, lexer.KeywordDefault:
			stmt, ok = p.evaluateCaseClause(ctx)
		case lexer.KeywordFor:
			stmt, ok = p.evaluateFor(ctx)
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
		default:
			p.atError(fmt.Sprintf("unknown statement token type %d (%s)", tokenType, value), token)
		}
	default:
		_, isAssignemnt := p.findBefore(lexer.ASSIGN_OPERATOR, lexer.NEWLINE)

		if isAssignemnt {
			stmt, ok = p.evaluateVarAssignment(ctx)
		} else {
			// Assume it's an expression.
			stmt, ok = p.evaluateExpression(ctx)
		}
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

		if operatorToken.Type == lexer.NUMBER_LITERAL {
			rightExpression, ok = p.evaluateSingleExpression(ctx)

			if ok {
				numberLiteral := rightExpression.(IntegerLiteral)

				if numberLiteral.Value < 0 {
					operator = "+"
				}
			}
		} else {
			operator = operatorToken.Value
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
