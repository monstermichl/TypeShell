package transpiler

import (
	"fmt"

	"github.com/monstermichl/typeshell/parser"
)

type ifBranchType int8

const (
	IF ifBranchType = iota
	ELSEIF
	ELSE
)

type expressionResult struct {
	values []string
}

func newExpressionResult(values ...string) expressionResult {
	return expressionResult{values: values}
}

func (r expressionResult) firstValue() string {
	length := len(r.values)

	if length == 0 {
		return ""
	}
	return r.values[0]
}

type transpiler struct {
	converter Converter
}

func New() transpiler {
	return transpiler{}
}

func (t *transpiler) Transpile(path string, converter Converter) (string, error) {
	p := parser.New()
	ast, err := p.Parse(path)

	if err != nil {
		return "", err
	}
	t.converter = converter
	err = t.evaluate(ast)

	if err != nil {
		return "", err
	}
	return t.converter.Dump()
}

func (t *transpiler) evaluateIndex(index parser.Expression, valueUsed bool) (expressionResult, error) {
	return t.evaluateExpression(index, true)
}

func (t *transpiler) evaluateProgram(program parser.Program) error {
	var err error
	t.converter.ProgramStart()

	for _, statementTemp := range program.Body() {
		err = t.evaluate(statementTemp)

		if err != nil {
			return err
		}
	}

	if err == nil {
		t.converter.ProgramEnd()
	}
	return err
}

func (t *transpiler) evaluateBooleanLiteral(literal parser.BooleanLiteral, valueUsed bool) (expressionResult, error) {
	return newExpressionResult(t.converter.BoolToString(literal.Value())), nil
}

func (t *transpiler) evaluateIntegerLiteral(literal parser.IntegerLiteral, valueUsed bool) (expressionResult, error) {
	return newExpressionResult(t.converter.IntToString(literal.Value())), nil
}

func (t *transpiler) evaluateStringLiteral(literal parser.StringLiteral, valueUsed bool) (expressionResult, error) {
	return newExpressionResult(t.converter.StringToString(literal.Value())), nil
}

func (t *transpiler) evaluateOperation(
	operation parser.Operation,
	callout func(left string, operator string, right string, valueType parser.ValueType, valueUsed bool) (string, error),
	valueUsed bool,
) (expressionResult, error) {
	left, err := t.evaluateExpression(operation.Left(), true)

	if err != nil {
		return expressionResult{}, err
	}
	right, err := t.evaluateExpression(operation.Right(), true)

	if err != nil {
		return expressionResult{}, err
	}
	s, err := callout(left.firstValue(), operation.Operator(), right.firstValue(), operation.Left().ValueType(), valueUsed)

	if err != nil {
		return expressionResult{}, err
	}
	return newExpressionResult(s), nil
}

func (t *transpiler) evaluateUnaryOperation(operation parser.UnaryOperation, valueUsed bool) (expressionResult, error) {
	expr := operation.Expression()
	result, err := t.evaluateExpression(expr, true)

	if err != nil {
		return expressionResult{}, err
	}
	s, err := t.converter.UnaryOperation(result.firstValue(), operation.Operator(), expr.ValueType(), valueUsed)

	if err != nil {
		return expressionResult{}, err
	}
	return newExpressionResult(s), nil
}

func (t *transpiler) evaluateBinaryOperation(operation parser.BinaryOperation, valueUsed bool) (expressionResult, error) {
	return t.evaluateOperation(operation, t.converter.BinaryOperation, valueUsed)
}

func (t *transpiler) evaluateCompareOperation(operation parser.Comparison, valueUsed bool) (expressionResult, error) {
	return t.evaluateOperation(operation, t.converter.Comparison, valueUsed)
}

func (t *transpiler) evaluateLogicalOperation(operation parser.LogicalOperation, valueUsed bool) (expressionResult, error) {
	return t.evaluateOperation(operation, t.converter.LogicalOperation, valueUsed)
}

func (t *transpiler) evaluateBreak(print parser.Break) error {
	return t.converter.Break()
}

func (t *transpiler) evaluateContinue(print parser.Continue) error {
	return t.converter.Continue()
}

func (t *transpiler) evaluatePrint(print parser.Print) error {
	values := []string{}

	for _, expr := range print.Expressions() {
		result, err := t.evaluateExpression(expr, true)

		if err != nil {
			return err
		}
		values = append(values, result.firstValue())
	}
	return t.converter.Print(values)
}

func (t *transpiler) evaluateWrite(write parser.Write) error {
	path := write.Path()
	valueType := path.ValueType()

	if !valueType.IsString() {
		return fmt.Errorf("expected string but got %s as read path", valueType.ToString())
	}
	result, err := t.evaluateExpression(path, true)

	if err != nil {
		return err
	}
	pathString := result.firstValue()
	data := write.Data()
	valueType = data.ValueType()

	if !valueType.IsString() {
		return fmt.Errorf("expected string but got %s as data", valueType.ToString())
	}
	result, err = t.evaluateExpression(data, true)

	if err != nil {
		return err
	}
	dataString := result.firstValue()
	appendExpr := write.Append()
	appendString := t.converter.BoolToString(false)

	if appendExpr != nil {
		append := write.Append()
		valueType = append.ValueType()

		if !valueType.IsBool() {
			return fmt.Errorf("expected bool but got %s as append flag", valueType.ToString())
		}
		result, err = t.evaluateExpression(append, true)

		if err != nil {
			return err
		}
		appendString = result.firstValue()
	}
	return t.converter.WriteFile(pathString, dataString, appendString)
}

func (t *transpiler) evaluateIf(ifStatement parser.If) error {
	conv := t.converter
	ifBranch := ifStatement.IfBranch()
	result, err := t.evaluateExpression(ifBranch.Condition(), true)

	if err != nil {
		return err
	}
	elifConditions := []string{}

	// Evaluate all conditions first to make sure they are added in front of the first if.
	for _, branch := range ifStatement.ElseIfBranches() {
		result, err := t.evaluateExpression(branch.Condition(), true)

		if err != nil {
			return err
		}
		elifConditions = append(elifConditions, result.firstValue())
	}
	err = conv.IfStart(result.firstValue())

	if err != nil {
		return err
	}
	err = t.evaluateBlock(ifBranch)

	if err != nil {
		return err
	}

	for i, branch := range ifStatement.ElseIfBranches() {
		err = conv.ElseIfStart(elifConditions[i])

		if err != nil {
			return err
		}
		err = t.evaluateBlock(branch)

		if err != nil {
			return err
		}
		err = conv.ElseIfEnd()

		if err != nil {
			return err
		}
	}

	if ifStatement.HasElse() {
		err = conv.ElseStart()

		if err != nil {
			return err
		}
		err = t.evaluateBlock(ifStatement.Else())

		if err != nil {
			return err
		}
		err = conv.ElseEnd()

		if err != nil {
			return err
		}
	}
	return conv.IfEnd()
}

func (t *transpiler) evaluateFor(forStatement parser.For) error {
	conv := t.converter
	init := forStatement.Init()

	if init != nil {
		err := t.evaluate(init)

		if err != nil {
			return err
		}
	}
	err := conv.ForStart()

	if err != nil {
		return err
	}
	result, err := t.evaluateExpression(forStatement.Condition(), true)

	if err != nil {
		return err
	}
	err = conv.ForCondition(result.firstValue())

	if err != nil {
		return err
	}
	err = t.evaluateBlock(forStatement)

	if err != nil {
		return err
	}
	increment := forStatement.Increment()

	if increment != nil {
		err = t.evaluate(increment)

		if err != nil {
			return err
		}
	}
	return conv.ForEnd()
}

func (t *transpiler) evaluateVarDefinition(definition parser.VariableDefinition) error {
	for i, variable := range definition.Variables() {
		result, err := t.evaluateExpression(definition.Values()[i], true)

		if err != nil {
			return err
		}
		err = t.converter.VarDefinition(variable.Name(), result.firstValue(), variable.Global())

		if err != nil {
			return err
		}
	}
	return nil
}

func (t *transpiler) evaluateVarDefinitionCallAssignment(definition parser.VariableDefinitionCallAssignment) error {
	result, err := t.evaluateExpression(definition.Call(), true)

	if err != nil {
		return err
	}
	variables := definition.Variables()
	variablesLen := len(variables)
	values := result.values
	valuesLen := len(values)

	if valuesLen != variablesLen {
		return fmt.Errorf("require %d values but got %d", variablesLen, valuesLen)
	}

	for i, variable := range variables {
		err = t.converter.VarDefinition(variable.Name(), values[i], variable.Global())

		if err != nil {
			return err
		}
	}
	return nil
}

func (t *transpiler) evaluateVarAssignment(assignment parser.VariableAssignment) error {
	for i, variable := range assignment.Variables() {
		result, err := t.evaluateExpression(assignment.Values()[i], true)

		if err != nil {
			return err
		}
		err = t.converter.VarDefinition(variable.Name(), result.firstValue(), variable.Global())

		if err != nil {
			return err
		}
	}
	return nil
}

func (t *transpiler) evaluateVarAssignmentCallAssignment(assignment parser.VariableAssignmentCallAssignment) error {
	result, err := t.evaluateExpression(assignment.Call(), true)

	if err != nil {
		return err
	}
	variables := assignment.Variables()
	variablesLen := len(variables)
	values := result.values
	valuesLen := len(values)

	if valuesLen != variablesLen {
		return fmt.Errorf("require %d values but got %d", variablesLen, valuesLen)
	}

	for i, variable := range variables {
		err = t.converter.VarDefinition(variable.Name(), values[i], variable.Global())

		if err != nil {
			return err
		}
	}
	return nil
}

func (t *transpiler) evaluateSliceAssignment(assignment parser.SliceAssignment) error {
	indexResult, err := t.evaluateIndex(assignment.Index(), true)

	if err != nil {
		return err
	}
	valueResult, err := t.evaluateExpression(assignment.Value(), true)

	if err != nil {
		return err
	}
	return t.converter.SliceAssignment(assignment.Name(), indexResult.firstValue(), valueResult.firstValue(), assignment.Global())
}

func (t *transpiler) evaluateVarEvaluation(evaluation parser.VariableEvaluation, valueUsed bool) (expressionResult, error) {
	s, err := t.converter.VarEvaluation(evaluation.Name(), valueUsed, evaluation.Global())

	if err != nil {
		return expressionResult{}, err
	}
	return newExpressionResult(s), nil
}

func (t *transpiler) evaluateSliceEvaluation(evaluation parser.SliceEvaluation, valueUsed bool) (expressionResult, error) {
	value := evaluation.Value()
	result, err := t.evaluateExpression(value, true)

	if err != nil {
		return expressionResult{}, err
	}
	valueString := result.firstValue()
	result, err = t.evaluateIndex(evaluation.Index(), true)

	if err != nil {
		return expressionResult{}, err
	}
	s, err := t.converter.SliceEvaluation(valueString, result.firstValue(), valueUsed)

	if err != nil {
		return expressionResult{}, err
	}
	return newExpressionResult(s), err
}

func (t *transpiler) evaluateStringSubscript(subscript parser.StringSubscript, valueUsed bool) (expressionResult, error) {
	indexResult, err := t.evaluateIndex(subscript.Index(), true)

	if err != nil {
		return expressionResult{}, err
	}
	value := subscript.Value()
	str, err := t.evaluateExpression(value, true)

	if err != nil {
		return expressionResult{}, err
	}
	s, err := t.converter.StringSubscript(str.firstValue(), indexResult.firstValue(), valueUsed)

	if err != nil {
		return expressionResult{}, err
	}
	return newExpressionResult(s), err
}

func (t *transpiler) evaluateGroup(group parser.Group, valueUsed bool) (expressionResult, error) {
	return t.evaluateExpression(group.Child(), valueUsed)
}

func (t *transpiler) evaluateBlock(block parser.Block) error {
	body := block.Body()
	length := len(body)

	if length == 0 {
		return t.converter.Nop()
	}
	for index, statement := range body {
		err := t.evaluate(statement)

		if err != nil {
			return err
		}
		if index == length-1 {
			return nil
		}
	}
	return nil
}

func (t *transpiler) evaluateReturn(returnStatement parser.Return) error {
	returnValues := []ReturnValue{}

	for _, expr := range returnStatement.Values() {
		result, err := t.evaluateExpression(expr, true)

		if err != nil {
			return err
		}
		returnValues = append(returnValues, ReturnValue{
			value:     result.firstValue(),
			valueType: expr.ValueType(),
		})
	}
	return t.converter.Return(returnValues)
}

func (t *transpiler) evaluateFunctionDefinition(functionDefinition parser.FunctionDefinition) error {
	name := functionDefinition.Name()
	params := []string{}

	for _, param := range functionDefinition.Params() {
		params = append(params, param.Name())
	}
	conv := t.converter
	err := conv.FuncStart(name, params, functionDefinition.ReturnTypes())

	if err != nil {
		return err
	}
	err = t.evaluateBlock(functionDefinition)

	if err != nil {
		return err
	}
	return conv.FuncEnd()
}

func (t *transpiler) evaluateFunctionCall(functionCall parser.FunctionCall, valueUsed bool) (expressionResult, error) {
	name := functionCall.Name()
	args := []string{}

	for _, arg := range functionCall.Args() {
		result, err := t.evaluateExpression(arg, true)

		if err != nil {
			return expressionResult{}, err
		}
		args = append(args, result.firstValue())
	}
	returnTypes := functionCall.ReturnTypes()
	values, err := t.converter.FuncCall(name, args, returnTypes, valueUsed)

	if err != nil {
		return expressionResult{}, err
	}
	returnTypesLen := len(returnTypes)
	valuesLen := len(values)

	if valueUsed && valuesLen != returnTypesLen {
		return expressionResult{}, fmt.Errorf("function \"%s\" must return %d values but returned %d", name, returnTypesLen, valuesLen)
	}
	return newExpressionResult(values...), nil
}

func (t *transpiler) evaluateAppCall(call parser.AppCall, valueUsed bool) (expressionResult, error) {
	convertedCalls := []AppCall{}
	nextCall := &call

	for nextCall != nil {
		name := nextCall.Name()
		args := []string{}

		for _, arg := range nextCall.Args() {
			result, err := t.evaluateExpression(arg, true)

			if err != nil {
				return expressionResult{}, err
			}
			args = append(args, result.firstValue())
		}
		convertedCalls = append(convertedCalls, AppCall{
			name: name,
			args: args,
		})
		nextCall = nextCall.Next()
	}
	s, err := t.converter.AppCall(convertedCalls, valueUsed)

	if err != nil {
		return expressionResult{}, err
	}
	return newExpressionResult(s), nil
}

func (t *transpiler) evaluateSliceInstantiation(instantiation parser.SliceInstantiation, valueUsed bool) (expressionResult, error) {
	values := []string{}

	for _, expr := range instantiation.Values() {
		result, err := t.evaluateExpression(expr, true)

		if err != nil {
			return expressionResult{}, err
		}
		values = append(values, result.firstValue())
	}
	s, err := t.converter.SliceInstantiation(values, valueUsed)

	if err != nil {
		return expressionResult{}, err
	}
	return newExpressionResult(s), nil
}

func (t *transpiler) evaluateInput(input parser.Input, valueUsed bool) (expressionResult, error) {
	promptString := ""
	prompt := input.Prompt()

	if prompt != nil {
		result, err := t.evaluateExpression(prompt, valueUsed)

		if err != nil {
			return expressionResult{}, err
		}
		promptString = result.firstValue()
	}
	s, err := t.converter.Input(promptString, valueUsed)

	if err != nil {
		return expressionResult{}, err
	}
	return newExpressionResult(s), nil
}

func (t *transpiler) evaluateCopy(copy parser.Copy, valueUsed bool) (expressionResult, error) {
	expr, err := t.evaluateExpression(copy.Source(), true)

	if err != nil {
		return expressionResult{}, err
	}
	destination := copy.Destination()
	amount, err := t.converter.Copy(destination.Name(), expr.firstValue(), valueUsed, destination.Global())

	if err != nil {
		return expressionResult{}, err
	}
	return expressionResult{
		values: []string{amount},
	}, nil
}

func (t *transpiler) evaluateLen(len parser.Len, valueUsed bool) (expressionResult, error) {
	expr := len.Expression()
	valueType := expr.ValueType()
	result, err := t.evaluateExpression(expr, true)

	if err != nil {
		return expressionResult{}, err
	}
	s := ""

	if valueType.IsString() {
		s, err = t.converter.StringLen(result.firstValue(), valueUsed)
	} else {
		s, err = t.converter.SliceLen(result.firstValue(), valueUsed)
	}

	if err != nil {
		return expressionResult{}, err
	}
	return newExpressionResult(s), nil
}

func (t *transpiler) evaluateRead(read parser.Read, valueUsed bool) (expressionResult, error) {
	path := read.Path()
	valueType := path.ValueType()

	if !valueType.IsString() {
		return expressionResult{}, fmt.Errorf("expected string but got %s as read path", valueType.ToString())
	}
	result, err := t.evaluateExpression(path, true)

	if err != nil {
		return expressionResult{}, err
	}
	s, err := t.converter.ReadFile(result.firstValue(), valueUsed)

	if err != nil {
		return expressionResult{}, err
	}
	return newExpressionResult(s), nil
}

func (t *transpiler) evaluate(statement parser.Statement) error {
	statementType := statement.StatementType()

	switch statementType {
	case parser.STATEMENT_TYPE_PROGRAM:
		return t.evaluateProgram(statement.(parser.Program))
	case parser.STATEMENT_TYPE_VAR_DEFINITION:
		return t.evaluateVarDefinition(statement.(parser.VariableDefinition))
	case parser.STATEMENT_TYPE_VAR_DEFINITION_CALL_ASSIGNMENT:
		return t.evaluateVarDefinitionCallAssignment(statement.(parser.VariableDefinitionCallAssignment))
	case parser.STATEMENT_TYPE_VAR_ASSIGNMENT:
		return t.evaluateVarAssignment(statement.(parser.VariableAssignment))
	case parser.STATEMENT_TYPE_VAR_ASSIGNMENT_CALL_ASSIGNMENT:
		return t.evaluateVarAssignmentCallAssignment(statement.(parser.VariableAssignmentCallAssignment))
	case parser.STATEMENT_TYPE_SLICE_ASSIGNMENT:
		return t.evaluateSliceAssignment(statement.(parser.SliceAssignment))
	case parser.STATEMENT_TYPE_FUNCTION_DEFINITION:
		return t.evaluateFunctionDefinition(statement.(parser.FunctionDefinition))
	case parser.STATEMENT_TYPE_RETURN:
		return t.evaluateReturn(statement.(parser.Return))
	case parser.STATEMENT_TYPE_IF:
		return t.evaluateIf(statement.(parser.If))
	case parser.STATEMENT_TYPE_FOR:
		return t.evaluateFor(statement.(parser.For))
	case parser.STATEMENT_TYPE_BREAK:
		return t.evaluateBreak(statement.(parser.Break))
	case parser.STATEMENT_TYPE_CONTINUE:
		return t.evaluateContinue(statement.(parser.Continue))
	case parser.STATEMENT_TYPE_PRINT:
		return t.evaluatePrint(statement.(parser.Print))
	case parser.STATEMENT_TYPE_WRITE:
		return t.evaluateWrite(statement.(parser.Write))
	default:
		expression, ok := statement.(parser.Expression)

		if !ok {
			return fmt.Errorf("statement is not an expression (%v)", statement)
		}
		_, err := t.evaluateExpression(expression, false)
		return err
	}
}

func (t *transpiler) evaluateExpression(expression parser.Expression, valueUsed bool) (expressionResult, error) {
	expressionType := expression.StatementType()

	switch expressionType {
	case parser.STATEMENT_TYPE_BOOL_LITERAL:
		return t.evaluateBooleanLiteral(expression.(parser.BooleanLiteral), valueUsed)
	case parser.STATEMENT_TYPE_INT_LITERAL:
		return t.evaluateIntegerLiteral(expression.(parser.IntegerLiteral), valueUsed)
	case parser.STATEMENT_TYPE_STRING_LITERAL:
		return t.evaluateStringLiteral(expression.(parser.StringLiteral), valueUsed)
	case parser.STATEMENT_TYPE_UNARY_OPERATION:
		return t.evaluateUnaryOperation(expression.(parser.UnaryOperation), valueUsed)
	case parser.STATEMENT_TYPE_BINARY_OPERATION:
		return t.evaluateBinaryOperation(expression.(parser.BinaryOperation), valueUsed)
	case parser.STATEMENT_TYPE_COMPARISON:
		return t.evaluateCompareOperation(expression.(parser.Comparison), valueUsed)
	case parser.STATEMENT_TYPE_LOGICAL_OPERATION:
		return t.evaluateLogicalOperation(expression.(parser.LogicalOperation), valueUsed)
	case parser.STATEMENT_TYPE_VAR_EVALUATION:
		return t.evaluateVarEvaluation(expression.(parser.VariableEvaluation), valueUsed)
	case parser.STATEMENT_TYPE_SLICE_EVALUATION:
		return t.evaluateSliceEvaluation(expression.(parser.SliceEvaluation), valueUsed)
	case parser.STATEMENT_TYPE_STRING_SUBSCRIPT:
		return t.evaluateStringSubscript(expression.(parser.StringSubscript), valueUsed)
	case parser.STATEMENT_TYPE_GROUP:
		return t.evaluateGroup(expression.(parser.Group), valueUsed)
	case parser.STATEMENT_TYPE_FUNCTION_CALL:
		return t.evaluateFunctionCall(expression.(parser.FunctionCall), valueUsed)
	case parser.STATEMENT_TYPE_APP_CALL:
		return t.evaluateAppCall(expression.(parser.AppCall), valueUsed)
	case parser.STATEMENT_TYPE_SLICE_INSTANTIATION:
		return t.evaluateSliceInstantiation(expression.(parser.SliceInstantiation), valueUsed)
	case parser.STATEMENT_TYPE_INPUT:
		return t.evaluateInput(expression.(parser.Input), valueUsed)
	case parser.STATEMENT_TYPE_COPY:
		return t.evaluateCopy(expression.(parser.Copy), valueUsed)
	case parser.STATEMENT_TYPE_LEN:
		return t.evaluateLen(expression.(parser.Len), valueUsed)
	case parser.STATEMENT_TYPE_READ:
		return t.evaluateRead(expression.(parser.Read), valueUsed)
	}
	return expressionResult{}, fmt.Errorf("unknown expression type %s", expressionType)
}
