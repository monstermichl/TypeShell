package parser

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/monstermichl/typeshell/lexer"
)

var typeMapping = map[lexer.VarType]DataType{
	lexer.DATA_TYPE_BOOLEAN: DATA_TYPE_BOOLEAN,
	lexer.DATA_TYPE_INTEGER: DATA_TYPE_INTEGER,
	lexer.DATA_TYPE_STRING:  DATA_TYPE_STRING,
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
	variables  map[string]Variable
	functions  map[string]FunctionDefinition
	scopeStack []scope
}

func (c context) currentScope() scope {
	return c.scopeStack[len(c.scopeStack)-1]
}

func (c context) findScope(s scope) bool {
	for i := len(c.scopeStack) - 1; i >= 0; i-- {
		if c.scopeStack[i] == s {
			return true
		}
	}
	return false
}

type blockCallback func(statements []Statement, last bool) error

type Parser struct {
	tokens []lexer.Token
	index  int
}

func New(tokens []lexer.Token) Parser {
	return Parser{
		tokens: tokens,
	}
}

func (p *Parser) Parse() (Program, error) {
	return p.evaluateProgram()
}

func expectedError(what string, token lexer.Token) error {
	return fmt.Errorf("expected %s at row %d, column %d", what, token.Row(), token.Column())
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

func (p *Parser) eat() lexer.Token {
	token := p.peek()
	p.index++

	return token
}

func (p *Parser) isShortVarInit() bool {
	// Short initialization is an identifier plus the short init operator (e.g. x := 123).
	return p.peek().Type() == lexer.IDENTIFIER && p.peekAt(1).Type() == lexer.SHORT_INIT_OPERATOR
}

func (p *Parser) checkNewVariableNameToken(token lexer.Token, ctx context) error {
	name := token.Value()
	_, exists := ctx.variables[name]

	if exists {
		return fmt.Errorf("variable %s has already been defined at row %d, column %d", name, token.Row(), token.Column())
	}
	return nil
}

func (p *Parser) evaluateBuiltInFunction(tokenType lexer.TokenType, keyword string, minArgs int, maxArg int, ctx context, stmtCallout func(keywordToken lexer.Token, expressions []Expression) (Statement, error)) (Statement, error) {
	keywordToken := p.eat()

	if keywordToken.Type() != tokenType {
		return nil, expectedError(fmt.Sprintf("%s-keyword", keyword), keywordToken)
	}
	nextToken := p.eat()

	// Make sure after the print call comes a  opening round bracket.
	if nextToken.Type() != lexer.OPENING_ROUND_BRACKET {
		return nil, expectedError("\"(\"", nextToken)
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
				return nil, expectedError("\",\" or \")\"", nextToken)
			}
		}
	}
	expressionsLength := len(expressions)

	if minArgs < 0 {
		minArgs = 0
	}
	if expressionsLength < minArgs {
		return nil, expectedError(fmt.Sprintf("at least %d arguments for %s", minArgs, keyword), keywordToken)
	}
	if maxArg >= 0 && expressionsLength > maxArg {
		return nil, expectedError(fmt.Sprintf("a maximum of %d arguments for %s", minArgs, keyword), keywordToken)
	}
	nextToken = p.eat()

	// Make sure print call is terminated with a closing round bracket.
	if nextToken.Type() != lexer.CLOSING_ROUND_BRACKET {
		return nil, expectedError("\")\"", nextToken)
	}
	return stmtCallout(keywordToken, expressions)
}

func (p *Parser) evaluateProgram() (Program, error) {
	ctx := context{
		variables: map[string]Variable{},
		functions: map[string]FunctionDefinition{},
	}
	statements, err := p.evaluateBlockContent(lexer.EOF, nil, ctx, SCOPE_PROGRAM)

	if err != nil {
		return Program{}, err
	}
	return Program{
		body: statements,
	}, nil
}

func (p *Parser) evaluateBlockBegin() error {
	beginToken := p.eat()

	if beginToken.Type() != lexer.OPENING_CURLY_BRACKET {
		return expectedError("block begin", beginToken)
	}
	newlineToken := p.eat()

	if newlineToken.Type() != lexer.NEWLINE {
		return expectedError("newline", newlineToken)
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
	handleVarDefinition := func() (Statement, error) {
		stmtTemp, errTemp := p.evaluateVarDefinition(ctx)

		// Add new variable to scope variables.
		if errTemp == nil {
			variable := stmtTemp.(VariableDefinition).variable
			ctx.variables[variable.Name()] = variable
		}
		return stmtTemp, errTemp
	}

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
		case lexer.VAR_DEFINITION:
			stmt, err = handleVarDefinition()
		case lexer.FUNCTION_DEFINITION:
			stmt, err = p.evaluateFunctionDefinition(ctx)

			// Add new function to scope functions.
			if err == nil {
				function := stmt.(FunctionDefinition)
				ctx.functions[function.Name()] = function
			}
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
				stmt, err = handleVarDefinition()
			} else {
				// If token is identifier it could be a slice assignment.
				if token.Type() == lexer.IDENTIFIER {
					variable, exists := ctx.variables[token.Value()]

					// If variable has been defined and is a slice, handles slice assignment.
					if exists && variable.ValueType().IsSlice() {
						stmt, err = p.evaluateSliceAssignment(ctx)
					}
				}

				if err == nil && stmt == nil {
					stmt, err = p.evaluateExpression(ctx)
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
		} else if terminationToken.Type() != terminationTokenType {
			err = expectedError("termination token", terminationToken)
			break
		}
	}
	return statements, err
}

func (p *Parser) evaluateBlockEnd() error {
	endToken := p.eat()

	if endToken.Type() != lexer.CLOSING_CURLY_BRACKET {
		return expectedError("block end", endToken)
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

func (p *Parser) evaluateValueType(ctx context) (ValueType, error) {
	nextToken := p.peek()
	evaluatedType := NewValueType(DATA_TYPE_UNKNOWN, false)

	// Evaluate if value type is a slice type.
	if nextToken.Type() == lexer.OPENING_SQUARE_BRACKET {
		p.eat()             // Eat opening square bracket.
		nextToken = p.eat() // Eat closing square bracket.

		if nextToken.Type() != lexer.CLOSING_SQUARE_BRACKET {
			return evaluatedType, expectedError("\"]\"", nextToken)
		}
		nextToken = p.peek()
		evaluatedType.isSlice = true
	}

	// Evaluate data type.
	if nextToken.Type() != lexer.DATA_TYPE {
		return evaluatedType, expectedError("data type", nextToken)
	}
	p.eat() // Eat data type token.
	dataType, exists := typeMapping[nextToken.Value()]

	if !exists {
		return evaluatedType, expectedError("valid data type", nextToken)
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
			return nil, expectedError("variable definition", varToken)
		}
	}
	nameToken := p.eat()

	if nameToken.Type() != lexer.IDENTIFIER {
		return nil, expectedError("variable name", nameToken)
	}
	err := p.checkNewVariableNameToken(nameToken, ctx)

	if err != nil {
		return nil, err
	}
	name := nameToken.Value()
	specifiedType := NewValueType(DATA_TYPE_UNKNOWN, false)

	if isShortVarInit {
		nextToken := p.eat() // Eat short init operator.

		if nextToken.Type() != lexer.SHORT_INIT_OPERATOR {
			return nil, expectedError("short initialization operator", nextToken)
		}
	} else {
		nextToken := p.peek()

		// If next token starts a type definition, evaluate value type.
		if slices.Contains([]lexer.TokenType{lexer.DATA_TYPE, lexer.OPENING_SQUARE_BRACKET}, nextToken.Type()) {
			specifiedTypeTemp, err := p.evaluateValueType(ctx)

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
			return nil, expectedError("data type or value assignment", nextToken)
		} else if nextTokenType == lexer.ASSIGN_OPERATOR {
			p.eat()
		}
	}
	nextToken := p.peek()
	nextTokenType := nextToken.Type()
	var value Expression

	// TODO: Improve check.
	if nextTokenType != lexer.NEWLINE && nextTokenType != lexer.EOF {
		value, err = p.evaluateExpression(ctx)

		if err != nil {
			return nil, err
		}
		valueType := value.ValueType()

		if specifiedType.DataType() != DATA_TYPE_UNKNOWN {
			if !valueType.Equals(specifiedType) {
				return nil, expectedError(fmt.Sprintf("%s but got %s", specifiedType.ToString(), valueType.ToString()), nextToken)
			}
		} else {
			specifiedType = valueType
		}
	}

	// If no value has been specified, define default value.
	if value == nil {
		value, err = defaultVarValue(specifiedType)

		if err != nil {
			return nil, err
		}
	}
	variable := VariableDefinition{
		variable: NewVariable(name, specifiedType),
		value:    value,
	}
	return variable, nil
}

func (p *Parser) evaluateVarAssignment(ctx context) (Expression, error) {
	nameToken := p.eat()

	if nameToken.Type() != lexer.IDENTIFIER {
		return nil, expectedError("variable name", nameToken)
	}
	name := nameToken.Value()

	// Make sure variable has been defined.
	definedVariable, exists := ctx.variables[name]

	if !exists {
		return nil, fmt.Errorf("variable %s has not been defined at row %d, column %d", name, nameToken.Row(), nameToken.Column())
	}

	// Check assign token.
	if p.eat().Type() != lexer.ASSIGN_OPERATOR {
		return nil, expectedError("\"=\"", nameToken)
	}
	valueToken := p.peek()
	value, err := p.evaluateExpression(ctx)

	if err != nil {
		return nil, err
	}
	assignedValueType := value.ValueType()
	expectedValueType := definedVariable.ValueType()

	if assignedValueType != expectedValueType {
		return nil, expectedError(fmt.Sprintf("%s but got %s", expectedValueType.ToString(), assignedValueType.ToString()), valueToken)
	}
	return VariableAssignment{
		variable: definedVariable,
		value:    value,
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
			return params, expectedError("parameter name", nameToken)
		}
		p.eat()

		name := nameToken.Value()
		_, exists := ctx.variables[name]

		if exists {
			return params, fmt.Errorf("scope already contains a variable with the name %s", name)
		}
		valueType, err := p.evaluateValueType(ctx)

		if err != nil {
			return nil, err
		}
		nextToken := p.peek()
		nextTokenType := nextToken.Type()

		if nextTokenType != lexer.COMMA && nextTokenType != lexer.CLOSING_ROUND_BRACKET {
			return params, expectedError("\",\" or \")\"", nextToken)
		} else if nextTokenType == lexer.COMMA {
			p.eat()
		}
		params = append(params, NewVariable(name, valueType))
	}
	return params, nil
}

func (p *Parser) evaluateFunctionDefinition(ctx context) (Statement, error) {
	functionToken := p.eat()

	if ctx.currentScope() != SCOPE_PROGRAM {
		return nil, expectedError("function definition at top level", functionToken)
	}
	if functionToken.Type() != lexer.FUNCTION_DEFINITION {
		return nil, expectedError("function definition", functionToken)
	}
	nameToken := p.eat()
	name := nameToken.Value()

	if nameToken.Type() != lexer.IDENTIFIER {
		return nil, expectedError("function name", nameToken)
	}

	// Make sure no function exists with the same name.
	_, exists := ctx.functions[name]

	if exists {
		return nil, expectedError("unique function name", nameToken)
	}
	openingBrace := p.peek()
	params := []Variable{}

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
			return nil, expectedError("closing bracket", closingBrace)
		}
	}
	returnTypeToken := p.peek()
	returnType := NewValueType(DATA_TYPE_VOID, false)

	// Check if a return type has been specified.
	if returnTypeToken.Type() == lexer.DATA_TYPE {
		returnTypeTemp, err := p.evaluateValueType(ctx)

		if err != nil {
			return nil, err
		}
		returnType = returnTypeTemp
	}

	// Create ctx copy and add params as variables.
	ctxCopy := ctx
	for _, param := range params {
		ctxCopy.variables[param.Name()] = param
	}

	statements, err := p.evaluateBlock(func(statements []Statement, last bool) error {
		var errTemp error
		var lastStatement Statement
		length := len(statements)

		if length > 0 {
			lastStatement = statements[length-1]
		}

		if returnType.DataType() != DATA_TYPE_VOID {
			// If a return value is required, the last statement must be a return statement.
			if last {
				if lastStatement == nil || lastStatement.StatementType() != STATEMENT_TYPE_RETURN {
					errTemp = fmt.Errorf("function %s requires a return statement at the end of the block", name)
				} else if returnedType := lastStatement.(Return).value.ValueType(); !returnedType.Equals(returnType) {
					errTemp = fmt.Errorf("function %s returns %s but expects %s", name, returnedType.ToString(), returnType.ToString())
				}
			}
		} else if lastStatement != nil && lastStatement.StatementType() == STATEMENT_TYPE_RETURN {
			errTemp = fmt.Errorf("function %s must not have a return statement", name)
		}
		return errTemp
	}, ctxCopy, SCOPE_FUNCTION)

	if err != nil {
		return nil, err
	}
	return FunctionDefinition{
		name:      name,
		valueType: returnType,
		params:    params,
		body:      statements,
	}, nil
}

func (p *Parser) evaluateReturn(ctx context) (Statement, error) {
	returnToken := p.eat()

	if !ctx.findScope(SCOPE_FUNCTION) {
		return nil, expectedError(fmt.Sprintf("return within %s-scope", SCOPE_FUNCTION), returnToken)
	}
	if returnToken.Type() != lexer.RETURN {
		return nil, expectedError("return-keyword", returnToken)
	}
	value, err := p.evaluateExpression(ctx)

	if err != nil {
		return nil, err
	}
	return Return{
		value: value,
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
		return nil, expectedError(fmt.Sprintf("break statement within %s-scope", strings.Join(scopesToString(breakScopes), "- or ")), breakToken)
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
		return nil, expectedError(fmt.Sprintf("continue statement within %s-scope", strings.Join(scopesToString(breakScopes), "- or ")), continueToken)
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
				return nil, expectedError("if-keyword", nextToken)
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
				return nil, expectedError("boolean expression", conditionToken)
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
		return nil, expectedError("for-keyword", forToken)
	}
	nextToken := p.peek()

	// If next token is an identifier, parse a for-range statement.
	if nextToken.Type() == lexer.IDENTIFIER {
		p.eat()
		err := p.checkNewVariableNameToken(nextToken, ctx)

		if err != nil {
			return nil, err
		}
		indexVarName := nextToken.Value()
		nextToken = p.eat()

		if nextToken.Type() != lexer.COMMA {
			return nil, expectedError("\",\"", nextToken)
		}
		nextToken = p.eat()
		err = p.checkNewVariableNameToken(nextToken, ctx)

		if err != nil {
			return nil, err
		}
		valueVarName := nextToken.Value()
		nextToken = p.eat()

		if nextToken.Type() != lexer.SHORT_INIT_OPERATOR {
			return nil, expectedError("\":=\"", nextToken)
		}
		nextToken = p.eat()

		if nextToken.Type() != lexer.RANGE {
			return nil, expectedError("range-keyword", nextToken)
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
			return nil, expectedError("slice", nextToken)
		}
		sliceValueType := sliceExpression.ValueType()
		indexVar := NewVariable(indexVarName, NewValueType(DATA_TYPE_INTEGER, false))
		valueVar := NewVariable(valueVarName, sliceValueType)

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

		// Convert for-range-loop to for-loop.
		indexEvaluation := VariableEvaluation{indexVar.Name(), indexVar.ValueType()}
		condition := NewComparison(indexEvaluation, COMPARE_OPERATOR_LESS, Len{expression: sliceExpression})
		assignment := VariableAssignment{
			variable: valueVar,
			value: SliceEvaluation{
				name:     "test",
				index:    indexEvaluation,
				dataType: sliceValueType.DataType(),
			},
		}
		updatedStatements := []Statement{assignment}
		updatedStatements = append(updatedStatements, statements...)
		updatedStatements = append(updatedStatements, VariableAssignment{
			variable: indexVar,
			value: BinaryOperation{
				left:      indexEvaluation,
				operator:  BINARY_OPERATOR_ADDITION,
				right:     IntegerLiteral{value: 1},
				valueType: indexEvaluation.ValueType(),
			},
		})

		x := For{
			condition: condition,
			body:      updatedStatements,
		}
		fmt.Println(x)
		return x, nil
	} else {
		expr, err := p.evaluateExpression(ctx)

		if err != nil {
			return nil, err
		}
		if !expr.ValueType().IsBool() {
			return nil, expectedError("boolean expression", nextToken)
		}
		statements, err := p.evaluateBlock(nil, ctx, SCOPE_FOR)

		if err != nil {
			return nil, err
		}
		return For{
			condition: expr,
			body:      statements,
		}, nil
	}
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
			return nil, expectedError("closing bracket", closingToken)
		}

	// Handle slice instantiation.
	case lexer.OPENING_SQUARE_BRACKET:
		expr, err = p.evaluateSliceInstantiation(ctx)

	// Handle input.
	case lexer.INPUT:
		expr, err = p.evaluateInput(ctx)

	// Handle len.
	case lexer.LEN:
		expr, err = p.evaluateLen(ctx)

	// Handle identifiers.
	case lexer.IDENTIFIER:
		nextToken := p.peekAt(1)

		// If the current token is an identifier and the next is an opening
		// brace, it's a function call, if the next is an assignment operator,
		// it's an assignment, if the next is a smaller sign, it's an app call,
		// otherwise it's a variable evaluation.
		switch nextToken.Type() {
		case lexer.OPENING_ROUND_BRACKET:
			expr, err = p.evaluateFunctionCall(ctx)
		case lexer.ASSIGN_OPERATOR:
			expr, err = p.evaluateVarAssignment(ctx)
		// case lexer.OPENING_CURLY_BRACKET:
		// 	expr, err = p.evaluateAppCall(ctx)
		case lexer.OPENING_SQUARE_BRACKET:
			expr, err = p.evaluateSliceEvaluation(ctx)
		default:
			p.eat() // Eat identifier token.
			name := token.Value()
			variable, exists := ctx.variables[name]

			if !exists {
				err = fmt.Errorf("variable %s has not been defined at row %d, column %d", name, nextToken.Row(), nextToken.Column())
			} else {
				expr = VariableEvaluation{
					name:      name,
					valueType: variable.ValueType(),
				}
			}
		}

	default:
		return nil, fmt.Errorf("unknown expression type %d \"%s\" at row %d, column %d", tokenType, value, token.Row(), token.Column())
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
func (p *Parser) evaluateMultiplication(ctx context) (Expression, error) {
	return p.evaluateBinaryOperation(ctx, []BinaryOperator{BINARY_OPERATOR_MULTIPLICATION, BINARY_OPERATOR_DIVISION, BINARY_OPERATOR_MODULO}, p.evaluateSingleExpression)
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
			return nil, expectedError(fmt.Sprintf("same binary operation types but got %s and %s", leftType.ToString(), rightType.ToString()), operatorToken)
		}
		allowedOperators = allowedBinaryOperators(leftType)

		if !slices.Contains(allowedOperators, operator) {
			return nil, expectedError(fmt.Sprintf("valid %s operator but got \"%s\"", leftType.ToString(), operator), operatorToken)
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
			return nil, expectedError(fmt.Sprintf("same comparison types but got %s and %s", leftType.DataType(), rightType.ToString()), operatorToken)
		}
		allowedOperators := allowedCompareOperators(leftType)

		if !slices.Contains(allowedOperators, operator) {
			return nil, expectedError(fmt.Sprintf("valid %s operator but got \"%s\"", leftType.ToString(), operator), operatorToken)
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
			return nil, expectedError("boolean value", conditionToken)
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
		return nil, expectedError("opening bracket", openingBraceToken)
	}
	nextToken := p.peek()
	args := []Expression{}
	paramsLength := len(params)
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
			return nil, expectedError(fmt.Sprintf("parameter %s (%s) but got %s", lastParamType.ToString(), param.Name(), lastArgType.ToString()), argToken)
		}
		nextToken = p.peek()
		tokenType := nextToken.Type()

		if !slices.Contains([]lexer.TokenType{lexer.COMMA, lexer.CLOSING_ROUND_BRACKET}, tokenType) {
			err = expectedError("comma or closing bracket", nextToken)
			break
		} else if tokenType == lexer.COMMA {
			p.eat()
		}
	}

	// Check for the appropriate arguments amount.
	argsLength := len(args)

	if len(args) != paramsLength {
		return nil, argsLengthError(argsLength)
	}

	if err != nil {
		return nil, err
	}
	closingBraceToken := p.eat()

	if closingBraceToken.Type() != lexer.CLOSING_ROUND_BRACKET {
		return nil, expectedError("closing bracket", closingBraceToken)
	}
	return args, nil
}

func (p *Parser) evaluateFunctionCall(ctx context) (Call, error) {
	nameToken := p.eat()
	name := nameToken.Value()

	if nameToken.Type() != lexer.IDENTIFIER {
		return nil, expectedError("function identifier", nameToken)
	}

	// Make sure function has been defined.
	definedFunction, exists := ctx.functions[name]

	if !exists {
		return nil, fmt.Errorf("function %s has not been defined at row %d, column %d", name, nameToken.Row(), nameToken.Column())
	}
	args, err := p.evaluateArguments("function", name, definedFunction.params, ctx)

	if err != nil {
		return nil, err
	}
	return FunctionCall{
		name:      name,
		arguments: args,
		valueType: definedFunction.ValueType(),
	}, nil
}

func (p *Parser) evaluateAppCall(ctx context) (Call, error) {
	nextToken := p.eat()
	name := nextToken.Value()

	if nextToken.Type() != lexer.IDENTIFIER {
		return nil, expectedError("program identifier", nextToken)
	}
	nextToken = p.eat()

	// Make sure after the program call comes an opening curly bracket.
	if nextToken.Type() != lexer.OPENING_CURLY_BRACKET {
		return nil, expectedError("\"{\"", nextToken)
	}
	args := []Expression{}
	nextToken = p.peek()

	// Evaluate arguments if it's a program call with arguments.
	if nextToken.Type() != lexer.CLOSING_CURLY_BRACKET {
		for {
			arg, err := p.evaluateExpression(ctx)

			if err != nil {
				return nil, err
			}
			args = append(args, arg)
			nextToken = p.peek()
			nextTokenType := nextToken.Type()

			if nextTokenType == lexer.COMMA {
				p.eat()
			} else if nextTokenType == lexer.CLOSING_CURLY_BRACKET {
				break
			} else {
				return nil, expectedError("\",\" or \"}\"", nextToken)
			}
		}
	}
	nextToken = p.eat()

	// Make sure program call is terminated with a closing curly bracket.
	if nextToken.Type() != lexer.CLOSING_CURLY_BRACKET {
		return nil, expectedError("\"}\"", nextToken)
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
	sliceValueType, err := p.evaluateValueType(ctx)

	if err != nil {
		return nil, err
	}
	if !sliceValueType.IsSlice() {
		return nil, expectedError(fmt.Sprintf("slice type but got %s", sliceValueType.ToString()), nextToken)
	}
	nextToken = p.eat()

	if nextToken.Type() != lexer.OPENING_CURLY_BRACKET {
		return nil, expectedError("\"{\"", nextToken)
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

			if !valueDataType.Equals(sliceValueType) {
				return nil, fmt.Errorf("%s cannot not be added to %s at row %d, column %d", valueDataType.ToString(), sliceValueType.ToString(), valueToken.Row(), valueToken.Column())
			}
			values = append(values, expr)
			nextToken = p.peek()
			nextTokenType := nextToken.Type()

			if nextTokenType == lexer.COMMA {
				p.eat()
			} else if nextTokenType == lexer.CLOSING_CURLY_BRACKET {
				break
			} else {
				fmt.Println(nextToken)
				return nil, expectedError("\",\" or \"}\"", nextToken)
			}
		}
	}
	nextToken = p.eat()

	if nextToken.Type() != lexer.CLOSING_CURLY_BRACKET {
		return nil, expectedError("\"}\"", nextToken)
	}
	return SliceInstantiation{
		dataType: sliceValueType.DataType(),
		values:   values,
	}, nil
}

func (p *Parser) evaluateSliceEvaluation(ctx context) (Expression, error) {
	nameToken := p.eat()

	if nameToken.Type() != lexer.IDENTIFIER {
		return nil, expectedError("slice variable", nameToken)
	}
	name := nameToken.Value()
	variable, exists := ctx.variables[name]

	if !exists {
		return nil, fmt.Errorf("variable %s has not been defined at row %d, column %d", name, nameToken.Row(), nameToken.Column())
	}
	if !variable.ValueType().IsSlice() {
		return nil, expectedError(fmt.Sprintf("slice but variable is of type %s", variable.ValueType().ToString()), nameToken)
	}
	nextToken := p.eat()

	if nextToken.Type() != lexer.OPENING_SQUARE_BRACKET {
		return nil, expectedError("\"[\"", nextToken)
	}
	nextToken = p.peek()
	index, err := p.evaluateExpression(ctx)

	if err != nil {
		return nil, err
	}
	indexValueType := index.ValueType()

	if !indexValueType.IsInt() {
		return nil, expectedError(fmt.Sprintf("%s as index but got %s", DATA_TYPE_INTEGER, indexValueType.ToString()), nextToken)
	}
	nextToken = p.eat()

	if nextToken.Type() != lexer.CLOSING_SQUARE_BRACKET {
		return nil, expectedError("\"]\"", nextToken)
	}
	return SliceEvaluation{
		name:     name,
		index:    index,
		dataType: variable.ValueType().DataType(),
	}, nil
}

func (p *Parser) evaluateSliceAssignment(ctx context) (Statement, error) {
	nameToken := p.eat()

	if nameToken.Type() != lexer.IDENTIFIER {
		return nil, expectedError("slice variable", nameToken)
	}
	name := nameToken.Value()
	variable, exists := ctx.variables[name]

	if !exists {
		return nil, fmt.Errorf("variable %s has not been defined at row %d, column %d", name, nameToken.Row(), nameToken.Column())
	}
	variableValueType := variable.ValueType()

	if !variableValueType.IsSlice() {
		return nil, expectedError(fmt.Sprintf("slice but variable is of type %s", variableValueType.ToString()), nameToken)
	}
	nextToken := p.eat()

	if nextToken.Type() != lexer.OPENING_SQUARE_BRACKET {
		return nil, expectedError("\"[\"", nextToken)
	}
	nextToken = p.peek()
	index, err := p.evaluateExpression(ctx)

	if err != nil {
		return nil, err
	}
	indexValueType := index.ValueType()

	if !indexValueType.IsInt() {
		return nil, expectedError(fmt.Sprintf("%s as index but got %s", DATA_TYPE_INTEGER, indexValueType.ToString()), nextToken)
	}
	nextToken = p.eat()

	if nextToken.Type() != lexer.CLOSING_SQUARE_BRACKET {
		return nil, expectedError("\"]\"", nextToken)
	}
	nextToken = p.eat()

	if nextToken.Type() != lexer.ASSIGN_OPERATOR {
		return nil, expectedError("\"=\"", nameToken)
	}
	valueToken := p.peek()
	value, err := p.evaluateExpression(ctx)

	if err != nil {
		return nil, err
	}
	variableDataType := variableValueType.DataType()
	assignedDataType := value.ValueType().DataType()

	if variableDataType != assignedDataType {
		return nil, expectedError(fmt.Sprintf("%s value but got %s", variableDataType, assignedDataType), valueToken)
	}
	return SliceAssignment{
		name:  name,
		index: index,
		value: value,
	}, nil
}

func (p *Parser) evaluateLen(ctx context) (Expression, error) {
	expr, err := p.evaluateBuiltInFunction(lexer.LEN, "len", 1, 1, ctx, func(keywordToken lexer.Token, expressions []Expression) (Statement, error) {
		expr := expressions[0]

		if !expr.ValueType().isSlice {
			return nil, expectedError("slice", keywordToken)
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
