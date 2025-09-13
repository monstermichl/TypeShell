package transpiler

import (
	"errors"
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

func BoolToString(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

func IntToString(i int) string {
	return strconv.Itoa(i)
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

func (t *transpiler) evaluateValueTypeDefaultValue(valueType parser.ValueType) (string, error) {
	var defaultValue string
	conv := t.converter

	switch valueType.Type().Kind() {
	case parser.TypeKindBool:
		defaultValue = BoolToString(false)
	case parser.TypeKindInt:
		defaultValue = IntToString(0)
	case parser.TypeKindString:
		defaultValue = conv.StringToString("")
	case parser.TypeKindStruct:
		structDefinition, valid := valueType.Type().(parser.StructDefinition)

		if !valid {
			return "", errors.New("struct declaration could not be evaluated")
		}
		values := []parser.StructValue{}

		for _, field := range structDefinition.Fields() {
			fieldValueType := field.ValueType()
			defaultValueTemp, err := t.evaluateValueTypeDefaultValue(fieldValueType)

			if err != nil {
				return "", err
			}
			values = append(values, parser.NewStructValue(field.Name(), fieldValueType, parser.NewStringLiteral(defaultValueTemp)))
		}

		// Create helper struct definition.
		structInitialization := parser.NewStructInitialization(structDefinition, values...)
		result, err := t.evaluateExpression(structInitialization, true)

		if err != nil {
			return "", err
		}
		defaultValue = result.firstValue()
	default:
		return "", fmt.Errorf(`no default value defined for %s`, valueType.String())
	}
	return defaultValue, nil
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

func (t *transpiler) evaluateTypeDefinition(instantiation parser.TypeDefinition, valueUsed bool) (expressionResult, error) {
	return t.evaluateExpression(instantiation.Value(), valueUsed)
}

func (t *transpiler) evaluateBooleanLiteral(literal parser.BooleanLiteral, valueUsed bool) (expressionResult, error) {
	return newExpressionResult(BoolToString(literal.Value())), nil
}

func (t *transpiler) evaluateIntegerLiteral(literal parser.IntegerLiteral, valueUsed bool) (expressionResult, error) {
	return newExpressionResult(IntToString(literal.Value())), nil
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

func (t *transpiler) evaluateBreak() error {
	return t.converter.Break()
}

func (t *transpiler) evaluateContinue() error {
	return t.converter.Continue()
}

func (t *transpiler) evaluatePrint(print parser.Print) error {
	values := []string{}

	for _, expr := range print.Expressions() {
		result, err := t.evaluateExpression(expr, true)

		if err != nil {
			return err
		}
		values = append(values, result.values...)
	}
	return t.converter.Print(values)
}

func (t *transpiler) evaluatePanic(panic parser.Panic) error {
	result, err := t.evaluateExpression(panic.Expression(), true)

	if err != nil {
		return err
	}
	return t.converter.Panic(fmt.Sprintf("panic: %s", result.firstValue()))
}

func (t *transpiler) evaluateWrite(write parser.Write) error {
	path := write.Path()
	valueType := path.ValueType()

	if !valueType.IsString() {
		return fmt.Errorf("expected string but got %s as read path", valueType.String())
	}
	result, err := t.evaluateExpression(path, true)

	if err != nil {
		return err
	}
	pathString := result.firstValue()
	data := write.Data()
	valueType = data.ValueType()

	if !valueType.IsString() {
		return fmt.Errorf("expected string but got %s as data", valueType.String())
	}
	result, err = t.evaluateExpression(data, true)

	if err != nil {
		return err
	}
	dataString := result.firstValue()
	appendExpr := write.Append()
	appendString := BoolToString(false)

	if appendExpr != nil {
		append := write.Append()
		valueType = append.ValueType()

		if !valueType.IsBool() {
			return fmt.Errorf("expected bool but got %s as append flag", valueType.String())
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
	increment := forStatement.Increment()

	if increment != nil {
		err = conv.ForIncrementStart()

		if err != nil {
			return err
		}
		err = t.evaluate(increment)

		if err != nil {
			return err
		}
		err = conv.ForIncrementEnd()

		if err != nil {
			return err
		}
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
	return conv.ForEnd()
}

func (t *transpiler) evaluateExpressionAssignment(assignedExpression parser.Expression) (expressionResult, error) {
	result, err := t.evaluateExpression(assignedExpression, true)
	value := result.firstValue()

	if err != nil {
		return expressionResult{}, err
	}

	switch evaluationType := assignedExpression.ValueType().Type().(type) {
	case parser.StructDefinition:
		newStruct, err := t.converter.StructInitialization([]StructValue{}, true)

		if err != nil {
			return expressionResult{}, err
		}

		// If expression is a struct, the values need to be copied to avoid manipulation of the original.
		for _, field := range evaluationType.Fields() {
			fieldName := field.Name()
			fieldValue, err := t.converter.StructEvaluation(value, fieldName, true)

			if err != nil {
				return expressionResult{}, nil
			}
			err = t.converter.StructAssignment(newStruct, fieldName, fieldValue, false)

			if err != nil {
				return expressionResult{}, err
			}
		}
		evaluatedValue, err := t.converter.VarEvaluation(newStruct, true, false)

		if err != nil {
			return expressionResult{}, err
		}
		value = evaluatedValue
	}
	return newExpressionResult(value), nil
}

func (t *transpiler) evaluateNamedValuesDefinition(definition parser.NamedValuesDefinition) error {
	for _, assignment := range definition.Assignments() {
		err := t.evaluate(assignment)

		if err != nil {
			return err
		}
	}
	return nil
}

func (t *transpiler) evaluateConstDefinition(definition parser.ConstDefinition) error {
	variables := []parser.Variable{}

	// Map const definition to var definition since constant check has already been performed by parser.
	for _, constant := range definition.Constants() {
		variables = append(variables, parser.NewVariable(constant.LayerName(), constant.ValueType(), constant.Layer(), constant.Public()))
	}
	return t.evaluateVarDefinition(parser.NewVariableDefinition(variables, definition.Values()))
}

func (t *transpiler) evaluateVarDefinition(definition parser.VariableDefinitionValueAssignment) error {
	for i, variable := range definition.Variables() {
		result, err := t.evaluateExpressionAssignment(definition.Values()[i])

		if err != nil {
			return err
		}
		err = t.converter.VarDefinition(variable.LayerName(), result.firstValue(), variable.Global())

		if err != nil {
			return err
		}
	}
	return nil
}

func (t *transpiler) evaluateVarDefinitionCallAssignment(definition parser.VariableDefinitionCallAssignment) error {
	result, err := t.evaluateExpression(definition.Call(), true) // TODO: Handle issue https://github.com/monstermichl/TypeShell/issues/54 here as well?

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
		err = t.converter.VarDefinition(variable.LayerName(), values[i], variable.Global())

		if err != nil {
			return err
		}
	}
	return nil
}

func (t *transpiler) evaluateVarAssignment(assignment parser.VariableAssignmentValueAssignment) error {
	for i, variable := range assignment.Variables() {
		result, err := t.evaluateExpressionAssignment(assignment.Values()[i])

		if err != nil {
			return err
		}
		err = t.converter.VarDefinition(variable.LayerName(), result.firstValue(), variable.Global())

		if err != nil {
			return err
		}
	}
	return nil
}

func (t *transpiler) evaluateVarAssignmentCallAssignment(assignment parser.VariableAssignmentCallAssignment) error {
	result, err := t.evaluateExpression(assignment.Call(), true) // TODO: Handle issue https://github.com/monstermichl/TypeShell/issues/54 here as well?

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
		err = t.converter.VarDefinition(variable.LayerName(), values[i], variable.Global())

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
	value := assignment.Value()
	result, err := t.evaluateExpressionAssignment(value)

	if err != nil {
		return err
	}
	valueType := value.ValueType()
	defaultValue, err := t.evaluateValueTypeDefaultValue(valueType)

	if err != nil {
		return err
	}
	return t.converter.SliceAssignment(assignment.LayerName(), indexResult.firstValue(), result.firstValue(), defaultValue, assignment.Global())
}

func (t *transpiler) evaluateStructAssignment(assignment parser.StructAssignment) error {
	value := assignment.Value()
	valueResult, err := t.evaluateExpressionAssignment(value.Value())

	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	return t.converter.StructAssignment(assignment.LayerName(), value.Name(), valueResult.firstValue(), assignment.Global())
}

func (t *transpiler) evaluateConstEvaluation(evaluation parser.ConstEvaluation, valueUsed bool) (expressionResult, error) {
	// Map const evaluation to var evaluation since constant evaluation works the same.
	varEvaluation := parser.NewVariableEvaluation(evaluation.LayerName(), evaluation.ValueType(), evaluation.Layer(), evaluation.Public())

	return t.evaluateVarEvaluation(varEvaluation, valueUsed)
}

func (t *transpiler) evaluateVarEvaluation(evaluation parser.VariableEvaluation, valueUsed bool) (expressionResult, error) {
	s, err := t.converter.VarEvaluation(evaluation.LayerName(), valueUsed, evaluation.Global())

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

func (t *transpiler) evaluateStructEvaluation(evaluation parser.StructEvaluation, valueUsed bool) (expressionResult, error) {
	s, err := t.converter.StructEvaluation(evaluation.LayerName(), evaluation.Field().Name(), valueUsed)

	if err != nil {
		return expressionResult{}, err
	}
	return newExpressionResult(s), err
}

func (t *transpiler) evaluateStringSubscript(subscript parser.StringSubscript, valueUsed bool) (expressionResult, error) {
	startIndexResult, err := t.evaluateIndex(subscript.StartIndex(), true)

	if err != nil {
		return expressionResult{}, err
	}
	endIndexResult, err := t.evaluateIndex(subscript.EndIndex(), true)

	if err != nil {
		return expressionResult{}, err
	}
	value := subscript.Value()
	str, err := t.evaluateExpression(value, true)

	if err != nil {
		return expressionResult{}, err
	}
	s, err := t.converter.StringSubscript(str.firstValue(), startIndexResult.firstValue(), endIndexResult.firstValue(), valueUsed)

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
		params = append(params, param.LayerName())
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
		result, err := t.evaluateExpressionAssignment(arg)

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
	values, err := t.converter.AppCall(convertedCalls, valueUsed)

	if err != nil {
		return expressionResult{}, err
	}
	return newExpressionResult(values...), nil
}

func (t *transpiler) evaluateSliceInstantiation(instantiation parser.SliceInstantiation, valueUsed bool) (expressionResult, error) {
	values := []string{}

	for _, expr := range instantiation.Values() {
		result, err := t.evaluateExpressionAssignment(expr)

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

func (t *transpiler) evaluateStructInitialization(definition parser.StructInitialization, valueUsed bool) (expressionResult, error) {
	values := []StructValue{}

	for _, value := range definition.Values() {
		result, err := t.evaluateExpressionAssignment(value.Value())

		if err != nil {
			return expressionResult{}, err
		}
		values = append(values, StructValue{
			name:  value.Name(),
			value: result.firstValue(),
		})
	}
	s, err := t.converter.StructInitialization(values, valueUsed)

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
	amount, err := t.converter.Copy(destination.LayerName(), expr.firstValue(), valueUsed, destination.Global())

	if err != nil {
		return expressionResult{}, err
	}
	return expressionResult{
		values: []string{amount},
	}, nil
}

func (t *transpiler) evaluateItoa(itoa parser.Itoa, valueUsed bool) (expressionResult, error) {
	result, err := t.evaluateExpression(itoa.Value(), true)

	if err != nil {
		return expressionResult{}, err
	}
	return expressionResult{
		values: []string{result.firstValue()},
	}, nil
}

func (t *transpiler) evaluateExists(exists parser.Exists, valueUsed bool) (expressionResult, error) {
	expr, err := t.evaluateExpression(exists.Path(), true)

	if err != nil {
		return expressionResult{}, err
	}
	s, err := t.converter.Exists(expr.firstValue(), valueUsed)

	if err != nil {
		return expressionResult{}, err
	}
	return expressionResult{
		values: []string{s},
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
		return expressionResult{}, fmt.Errorf("expected string but got %s as read path", valueType.String())
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
	case parser.STATEMENT_TYPE_TYPE_DECLARATION:
		return nil // Nothing to handle here, types are just relevant for the parser.
	case parser.STATEMENT_TYPE_NAMED_VALUES_DEFINITION:
		return t.evaluateNamedValuesDefinition(statement.(parser.NamedValuesDefinition))
	case parser.STATEMENT_TYPE_CONST_DEFINITION:
		return t.evaluateConstDefinition(statement.(parser.ConstDefinition))
	case parser.STATEMENT_TYPE_VAR_DEFINITION_VALUE_ASSIGNMENT:
		return t.evaluateVarDefinition(statement.(parser.VariableDefinitionValueAssignment))
	case parser.STATEMENT_TYPE_VAR_DEFINITION_CALL_ASSIGNMENT:
		return t.evaluateVarDefinitionCallAssignment(statement.(parser.VariableDefinitionCallAssignment))
	case parser.STATEMENT_TYPE_VAR_ASSIGNMENT_VALUE_ASSIGNMENT:
		return t.evaluateVarAssignment(statement.(parser.VariableAssignmentValueAssignment))
	case parser.STATEMENT_TYPE_VAR_ASSIGNMENT_CALL_ASSIGNMENT:
		return t.evaluateVarAssignmentCallAssignment(statement.(parser.VariableAssignmentCallAssignment))
	case parser.STATEMENT_TYPE_SLICE_ASSIGNMENT:
		return t.evaluateSliceAssignment(statement.(parser.SliceAssignment))
	case parser.STATEMENT_TYPE_STRUCT_ASSIGNMENT:
		return t.evaluateStructAssignment(statement.(parser.StructAssignment))
	case parser.STATEMENT_TYPE_FUNCTION_DEFINITION:
		return t.evaluateFunctionDefinition(statement.(parser.FunctionDefinition))
	case parser.STATEMENT_TYPE_RETURN:
		return t.evaluateReturn(statement.(parser.Return))
	case parser.STATEMENT_TYPE_IF:
		return t.evaluateIf(statement.(parser.If))
	case parser.STATEMENT_TYPE_FOR:
		return t.evaluateFor(statement.(parser.For))
	case parser.STATEMENT_TYPE_BREAK:
		return t.evaluateBreak()
	case parser.STATEMENT_TYPE_CONTINUE:
		return t.evaluateContinue()
	case parser.STATEMENT_TYPE_PRINT:
		return t.evaluatePrint(statement.(parser.Print))
	case parser.STATEMENT_TYPE_PANIC:
		return t.evaluatePanic(statement.(parser.Panic))
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
	case parser.STATEMENT_TYPE_TYPE_DEFINITION:
		return t.evaluateTypeDefinition(expression.(parser.TypeDefinition), valueUsed)
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
	case parser.STATEMENT_TYPE_CONST_EVALUATION:
		return t.evaluateConstEvaluation(expression.(parser.ConstEvaluation), valueUsed)
	case parser.STATEMENT_TYPE_VAR_EVALUATION:
		return t.evaluateVarEvaluation(expression.(parser.VariableEvaluation), valueUsed)
	case parser.STATEMENT_TYPE_SLICE_EVALUATION:
		return t.evaluateSliceEvaluation(expression.(parser.SliceEvaluation), valueUsed)
	case parser.STATEMENT_TYPE_STRUCT_EVALUATION:
		return t.evaluateStructEvaluation(expression.(parser.StructEvaluation), valueUsed)
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
	case parser.STATEMENT_TYPE_STRUCT_DEFINITION:
		return t.evaluateStructInitialization(expression.(parser.StructInitialization), valueUsed)
	case parser.STATEMENT_TYPE_INPUT:
		return t.evaluateInput(expression.(parser.Input), valueUsed)
	case parser.STATEMENT_TYPE_COPY:
		return t.evaluateCopy(expression.(parser.Copy), valueUsed)
	case parser.STATEMENT_TYPE_EXISTS:
		return t.evaluateExists(expression.(parser.Exists), valueUsed)
	case parser.STATEMENT_TYPE_ITOA:
		return t.evaluateItoa(expression.(parser.Itoa), valueUsed)
	case parser.STATEMENT_TYPE_LEN:
		return t.evaluateLen(expression.(parser.Len), valueUsed)
	case parser.STATEMENT_TYPE_READ:
		return t.evaluateRead(expression.(parser.Read), valueUsed)
	}
	return expressionResult{}, fmt.Errorf("unknown expression type %s", expressionType)
}
