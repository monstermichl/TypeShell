package transpiler

import (
	"fmt"
	"strconv"

	"github.com/monstermichl/typeshell/parser"
)

type ifBranchType int8

const (
	IF ifBranchType = iota
	ELSEIF
	ELSE
)

type transpiler struct {
	ast       parser.Statement
	converter Converter
}

func New(ast parser.Statement) transpiler {
	return transpiler{
		ast: ast,
	}
}

func (t *transpiler) Transpile(converter Converter) (string, error) {
	t.converter = converter
	err := t.evaluate(t.ast)

	if err != nil {
		return "", err
	}
	return t.converter.Dump()
}

func (t *transpiler) evaluateIndex(index parser.Expression, valueUsed bool) (int, error) {
	i := -1
	indexString, err := t.evaluateExpression(index, true)

	if err != nil {
		return i, err
	}
	i, err = strconv.Atoi(indexString)

	if err != nil {
		return i, err
	}
	return i, nil
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

func (t *transpiler) evaluateBooleanLiteral(literal parser.BooleanLiteral, valueUsed bool) (string, error) {
	return t.converter.BoolToString(literal.Value()), nil
}

func (t *transpiler) evaluateIntegerLiteral(literal parser.IntegerLiteral, valueUsed bool) (string, error) {
	return t.converter.IntToString(literal.Value()), nil
}

func (t *transpiler) evaluateStringLiteral(literal parser.StringLiteral, valueUsed bool) (string, error) {
	return t.converter.StringToString(literal.Value()), nil
}

func (t *transpiler) evaluateOperation(
	operation parser.Operation,
	callout func(left string, operator string, right string, valueType parser.ValueType, valueUsed bool) (string, error),
	valueUsed bool,
) (string, error) {
	left, err := t.evaluateExpression(operation.Left(), true)

	if err != nil {
		return "", err
	}
	right, err := t.evaluateExpression(operation.Right(), true)

	if err != nil {
		return "", err
	}
	return callout(left, operation.Operator(), right, operation.Left().ValueType(), valueUsed)
}

func (t *transpiler) evaluateBinaryOperation(operation parser.BinaryOperation, valueUsed bool) (string, error) {
	return t.evaluateOperation(operation, t.converter.BinaryOperation, valueUsed)
}

func (t *transpiler) evaluateCompareOperation(operation parser.Comparison, valueUsed bool) (string, error) {
	return t.evaluateOperation(operation, t.converter.Comparison, valueUsed)
}

func (t *transpiler) evaluateLogicalOperation(operation parser.LogicalOperation, valueUsed bool) (string, error) {
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
		value, err := t.evaluateExpression(expr, true)

		if err != nil {
			return err
		}
		values = append(values, value)
	}
	return t.converter.Print(values)
}

func (t *transpiler) evaluateIf(ifStatement parser.If) error {
	conv := t.converter
	ifBranch := ifStatement.IfBranch()
	condition, err := t.evaluateExpression(ifBranch.Condition(), true)

	if err != nil {
		return err
	}
	elifConditions := []string{}

	// Evaluate all conditions first to make sure they are added in front of the first if.
	for _, branch := range ifStatement.ElseIfBranches() {
		condition, err := t.evaluateExpression(branch.Condition(), true)

		if err != nil {
			return err
		}
		elifConditions = append(elifConditions, condition)
	}
	err = conv.IfStart(condition)

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

func (t *transpiler) evaluateWhile(whileStatement parser.For) error {
	conv := t.converter
	err := conv.ForStart()

	if err != nil {
		return err
	}
	condition, err := t.evaluateExpression(whileStatement.Condition(), true)

	if err != nil {
		return err
	}
	err = conv.ForCondition(condition)

	if err != nil {
		return err
	}
	err = t.evaluateBlock(whileStatement)

	if err != nil {
		return err
	}
	return conv.ForEnd()
}

func (t *transpiler) evaluateVarDefinition(definition parser.VariableDefinition) error {
	value, err := t.evaluateExpression(definition.Value(), true)

	if err != nil {
		return err
	}
	return t.converter.VarDefinition(definition.Variable().Name(), value)
}

func (t *transpiler) evaluateVarAssignment(assignment parser.VariableAssignment) error {
	value, err := t.evaluateExpression(assignment.Value(), true)

	if err != nil {
		return err
	}
	return t.converter.VarAssignment(assignment.Variable().Name(), value)
}

func (t *transpiler) evaluateSliceAssignment(assignment parser.SliceAssignment) error {
	index, err := t.evaluateIndex(assignment.Index(), true)

	if err != nil {
		return err
	}
	value, err := t.evaluateExpression(assignment.Value(), true)

	if err != nil {
		return err
	}
	return t.converter.SliceAssignment(assignment.Name(), index, value)
}

func (t *transpiler) evaluateVarEvaluation(evaluation parser.VariableEvaluation, valueUsed bool) (string, error) {
	return t.converter.VarEvaluation(evaluation.Name(), valueUsed)
}

func (t *transpiler) evaluateSliceEvaluation(evaluation parser.SliceEvaluation, valueUsed bool) (string, error) {
	index, err := t.evaluateIndex(evaluation.Index(), true)

	if err != nil {
		return "", err
	}
	return t.converter.SliceEvaluation(evaluation.Name(), index, valueUsed)
}

func (t *transpiler) evaluateGroup(group parser.Group, valueUsed bool) (string, error) {
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
	value, err := t.evaluateExpression(returnStatement.Value(), true)

	if err != nil {
		return err
	}
	return t.converter.Return(value, returnStatement.Value().ValueType())
}

func (t *transpiler) evaluateFunctionDefinition(functionDefinition parser.FunctionDefinition) error {
	name := functionDefinition.Name()
	params := []string{}

	for _, param := range functionDefinition.Params() {
		params = append(params, param.Name())
	}
	conv := t.converter
	err := conv.FuncStart(name, params)

	if err != nil {
		return err
	}
	err = t.evaluateBlock(functionDefinition)

	if err != nil {
		return err
	}
	return conv.FuncEnd(name)
}

func (t *transpiler) evaluateFunctionCall(functionCall parser.FunctionCall, valueUsed bool) (string, error) {
	name := functionCall.Name()
	args := []string{}

	for _, arg := range functionCall.Args() {
		value, err := t.evaluateExpression(arg, true)

		if err != nil {
			return "", err
		}
		args = append(args, value)
	}
	return t.converter.FuncCall(name, args, functionCall.ValueType(), valueUsed)
}

func (t *transpiler) evaluateAppCall(call parser.AppCall, valueUsed bool) (string, error) {
	convertedCalls := []AppCall{}
	nextCall := &call

	for nextCall != nil {
		name := nextCall.Name()
		args := []string{}

		for _, arg := range nextCall.Args() {
			value, err := t.evaluateExpression(arg, true)

			if err != nil {
				return "", err
			}
			args = append(args, value)
		}
		convertedCalls = append(convertedCalls, AppCall{
			name: name,
			args: args,
		})
		nextCall = nextCall.Next()
	}
	return t.converter.AppCall(convertedCalls, valueUsed)
}

func (t *transpiler) evaluateInput(input parser.Input, valueUsed bool) (string, error) {
	promptString := ""
	prompt := input.Prompt()

	if prompt != nil {
		var err error
		promptString, err = t.evaluateExpression(prompt, valueUsed)

		if err != nil {
			return "", err
		}
	}
	return t.converter.Input(promptString, valueUsed)
}

func (t *transpiler) evaluateSliceInstantiation(instantiation parser.SliceInstantiation, valueUsed bool) (string, error) {
	values := []string{}

	for _, expr := range instantiation.Values() {
		value, err := t.evaluateExpression(expr, true)

		if err != nil {
			return "", err
		}
		values = append(values, value)
	}
	return t.converter.SliceInstantiation(values, valueUsed)
}

func (t *transpiler) evaluate(statement parser.Statement) error {
	statementType := statement.StatementType()

	switch statementType {
	case parser.STATEMENT_TYPE_PROGRAM:
		return t.evaluateProgram(statement.(parser.Program))
	case parser.STATEMENT_TYPE_VAR_DEFINITION:
		return t.evaluateVarDefinition(statement.(parser.VariableDefinition))
	case parser.STATEMENT_TYPE_VAR_ASSIGNMENT:
		return t.evaluateVarAssignment(statement.(parser.VariableAssignment))
	case parser.STATEMENT_TYPE_SLICE_ASSIGNMENT:
		return t.evaluateSliceAssignment(statement.(parser.SliceAssignment))
	case parser.STATEMENT_TYPE_FUNCTION_DEFINITION:
		return t.evaluateFunctionDefinition(statement.(parser.FunctionDefinition))
	case parser.STATEMENT_TYPE_RETURN:
		return t.evaluateReturn(statement.(parser.Return))
	case parser.STATEMENT_TYPE_IF:
		return t.evaluateIf(statement.(parser.If))
	case parser.STATEMENT_TYPE_FOR:
		return t.evaluateWhile(statement.(parser.For))
	case parser.STATEMENT_TYPE_BREAK:
		return t.evaluateBreak(statement.(parser.Break))
	case parser.STATEMENT_TYPE_CONTINUE:
		return t.evaluateContinue(statement.(parser.Continue))
	case parser.STATEMENT_TYPE_PRINT:
		return t.evaluatePrint(statement.(parser.Print))
	default:
		expression, ok := statement.(parser.Expression)

		if !ok {
			return fmt.Errorf("statement is not an expression (%v)", statement)
		}
		_, err := t.evaluateExpression(expression, false)
		return err
	}
}

func (t *transpiler) evaluateExpression(expression parser.Expression, valueUsed bool) (string, error) {
	expressionType := expression.StatementType()

	switch expressionType {
	case parser.STATEMENT_TYPE_BOOL_LITERAL:
		return t.evaluateBooleanLiteral(expression.(parser.BooleanLiteral), valueUsed)
	case parser.STATEMENT_TYPE_INT_LITERAL:
		return t.evaluateIntegerLiteral(expression.(parser.IntegerLiteral), valueUsed)
	case parser.STATEMENT_TYPE_STRING_LITERAL:
		return t.evaluateStringLiteral(expression.(parser.StringLiteral), valueUsed)
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
	case parser.STATEMENT_TYPE_GROUP:
		return t.evaluateGroup(expression.(parser.Group), valueUsed)
	case parser.STATEMENT_TYPE_FUNCTION_CALL:
		return t.evaluateFunctionCall(expression.(parser.FunctionCall), valueUsed)
	case parser.STATEMENT_TYPE_APP_CALL:
		return t.evaluateAppCall(expression.(parser.AppCall), valueUsed)
	case parser.STATEMENT_TYPE_INPUT:
		return t.evaluateInput(expression.(parser.Input), valueUsed)
	case parser.STATEMENT_TYPE_SLICE_INSTANTIATION:
		return t.evaluateSliceInstantiation(expression.(parser.SliceInstantiation), valueUsed)
	}
	return "", fmt.Errorf("unknown expression type %s", expressionType)
}
