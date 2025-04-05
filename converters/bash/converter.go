package bash

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/monstermichl/typeshell/parser"
	"github.com/monstermichl/typeshell/transpiler"
)

type converter struct {
	interpreter       string
	code              []string
	varCounter        int
	sliceLenHelperSet bool
}

func New() *converter {
	return &converter{
		interpreter: "/bin/bash",
		code:        []string{},
	}
}

func varEvaluationString(name string) string {
	return fmt.Sprintf("${%s}", name)
}

func sliceAssignmentString(name string, index string, value string) string {
	return fmt.Sprintf("eval %s_%s=\"%s\"", name, index, value)
}

func (c *converter) BoolToString(value bool) string {
	if value {
		return "1"
	}
	return "0"
}

func (c *converter) IntToString(value int) string {
	return fmt.Sprintf("%d", value)
}

func (c *converter) StringToString(value string) string {
	return value
}

func (c *converter) Dump() (string, error) {
	return strings.Join(c.code, "\n"), nil
}

func (c *converter) ProgramStart() error {
	c.addLine(fmt.Sprintf("#!%s", c.interpreter))
	return nil
}

func (c *converter) ProgramEnd() error {
	return nil
}

func (c *converter) VarDefinition(name string, value string) error {
	c.addLine(fmt.Sprintf("%s=%s", name, value))
	return nil
}

func (c *converter) VarAssignment(name string, value string) error {
	c.addLine(fmt.Sprintf("%s=%s", name, value))
	return nil
}

func (c *converter) SliceAssignment(name string, index string, value string) error {
	c.addLine(sliceAssignmentString(varEvaluationString(name), index, value)) // TODO: Find out if using varEvaluationString here is a good idea because name might not be a variable.
	return nil
}

func (c *converter) FuncStart(name string, params []string) error {
	c.addLine(fmt.Sprintf("%s() {", name))

	for i, param := range params {
		c.addLine(fmt.Sprintf("local %s=$%d", param, i+1))
	}
	return nil
}

func (c *converter) FuncEnd(name string) error {
	c.addLine("}")
	return nil
}

func (c *converter) Return(value string, valueType parser.ValueType) error {
	if valueType.DataType() != parser.DATA_TYPE_VOID {
		c.Print([]string{value})
	}
	c.addLine("return")
	return nil
}

func (c *converter) IfStart(condition string) error {
	return c.ifStart(condition, "if")
}

func (c *converter) IfEnd() error {
	c.addLine("fi")
	return nil
}

func (c *converter) ElseIfStart(condition string) error {
	return c.ifStart(condition, "elif")
}

func (c *converter) ElseIfEnd() error {
	return nil
}

func (c *converter) ElseStart() error {
	c.addLine("else")
	return nil
}

func (c *converter) ElseEnd() error {
	return nil
}

func (c *converter) ForStart() error {
	c.addLine("while true; do")
	return nil
}

func (c *converter) ForCondition(condition string) error {
	c.addLine(fmt.Sprintf("if [ %s -ne %s ]; then break; fi", condition, c.BoolToString(true)))
	return nil
}

func (c *converter) ForEnd() error {
	c.addLine("done")
	return nil
}

func (c *converter) Break() error {
	c.addLine("break") // TODO: Break within switch. This might be a good solution -> https://stackoverflow.com/a/30874026.
	return nil
}

func (c *converter) Continue() error {
	c.addLine("continue")
	return nil
}

func (c *converter) Print(values []string) error {
	c.addLine(fmt.Sprintf("echo \"%s\"", strings.Join(values, " ")))
	return nil
}

func (c *converter) Nop() error {
	c.addLine(": # No operation")
	return nil
}

func (c *converter) BinaryOperation(left string, operator parser.BinaryOperator, right string, valueType parser.ValueType, valueUsed bool) (string, error) {
	helper := c.nextHelperVar()
	notAllowedError := func() (string, error) {
		return "", fmt.Errorf("binary operation %s is not allowed on type %s", operator, valueType.ToString())
	}

	if valueType.IsSlice() {
		return notAllowedError()
	}

	switch valueType.DataType() {
	case parser.DATA_TYPE_INTEGER:
		switch operator {
		case parser.BINARY_OPERATOR_MULTIPLICATION,
			parser.BINARY_OPERATOR_DIVISION,
			parser.BINARY_OPERATOR_MODULO,
			parser.BINARY_OPERATOR_ADDITION,
			parser.BINARY_OPERATOR_SUBTRACTION:
			// These operations are fine.
		default:
			return notAllowedError()
		}
		c.VarAssignment(helper, fmt.Sprintf("$(expr %s \\%s %s)", left, operator, right)) // Backslash is required for * operator to prevent pattern expansion (https://www.shell-tips.com/bash/math-arithmetic-calculation/#using-the-expr-command-line).
	case parser.DATA_TYPE_STRING:
		switch operator {
		case parser.BINARY_OPERATOR_ADDITION:
			c.VarAssignment(helper, fmt.Sprintf("\"%s%s\"", left, right))
		default:
			return notAllowedError()
		}
	default:
		return notAllowedError()
	}
	return c.VarEvaluation(helper, valueUsed)
}

func (c *converter) Comparison(left string, operator parser.CompareOperator, right string, valueType parser.ValueType, valueUsed bool) (string, error) {
	var operatorString string

	if !valueType.IsSlice() {
		switch valueType.DataType() {
		case parser.DATA_TYPE_BOOLEAN:
			switch operator {
			case parser.COMPARE_OPERATOR_EQUAL:
				operatorString = "-eq"
			case parser.COMPARE_OPERATOR_NOT_EQUAL:
				operatorString = "-ne"
			}
		case parser.DATA_TYPE_INTEGER:
			switch operator {
			case parser.COMPARE_OPERATOR_EQUAL:
				operatorString = "-eq"
			case parser.COMPARE_OPERATOR_NOT_EQUAL:
				operatorString = "-ne"
			case parser.COMPARE_OPERATOR_GREATER:
				operatorString = "-gt"
			case parser.COMPARE_OPERATOR_GREATER_OR_EQUAL:
				operatorString = "-ge"
			case parser.COMPARE_OPERATOR_LESS:
				operatorString = "-lt"
			case parser.COMPARE_OPERATOR_LESS_OR_EQUAL:
				operatorString = "-le"
			}
		case parser.DATA_TYPE_STRING:
			switch operator {
			case parser.COMPARE_OPERATOR_EQUAL:
				operatorString = "=="
			case parser.COMPARE_OPERATOR_NOT_EQUAL:
				operatorString = "!="
			}
		}
	}

	if len(operatorString) == 0 {
		return "", fmt.Errorf("comparison %s is not allowed on type %s", operator, valueType.ToString())
	}
	helper := c.nextHelperVar()

	c.VarAssignment(helper, fmt.Sprintf("$(if [ \"%s\" %s \"%s\" ]; then echo %s; else echo %s; fi)", left, operatorString, right, c.BoolToString(true), c.BoolToString(false)))
	return c.VarEvaluation(helper, valueUsed)
}

func (c *converter) LogicalOperation(left string, operator parser.LogicalOperator, right string, valueType parser.ValueType, valueUsed bool) (string, error) {
	var operatorString string

	switch operator {
	case parser.LOGICAL_OPERATOR_AND:
		operatorString = "&&"
	case parser.LOGICAL_OPERATOR_OR:
		operatorString = "||"
	default:
		return "", fmt.Errorf("unknown logical operator \"%s\"", operator)
	}
	trueString := c.BoolToString(true)
	helper := c.nextHelperVar()

	c.VarAssignment(helper, fmt.Sprintf("$(if [ \"%s\" -eq \"%s\" ] %s [ \"%s\" -eq \"%s\" ]; then echo %s; else echo %s; fi)", left, trueString, operatorString, right, trueString, trueString, c.BoolToString(false)))
	return c.VarEvaluation(helper, valueUsed)
}

func (c *converter) VarEvaluation(name string, valueUsed bool) (string, error) {
	return varEvaluationString(name), nil
}

func (c *converter) SliceInstantiation(values []string, valueUsed bool) (string, error) {
	helper := c.nextHelperVar()

	for i, value := range values {
		c.addLine(sliceAssignmentString(helper, strconv.Itoa(i), value))
	}
	return helper, nil
}

func (c *converter) SliceEvaluation(name string, index string, valueUsed bool) (string, error) {
	helper := c.nextHelperVar()
	c.VarAssignment(helper, fmt.Sprintf("$(eval \"echo \\${%s_%s}\")", varEvaluationString(name), index)) // TODO: Find out if using varEvaluationString here is a good idea because name might not be a variable.

	return c.VarEvaluation(helper, valueUsed)
}

func (c *converter) SliceLen(name string, valueUsed bool) (string, error) {
	if !c.sliceLenHelperSet {
		c.addLine("_sl() {")
		c.addLine("local _l=0")
		c.addLine("while true; do")
		c.addLine("eval \"local _t=\\${$1_${_l}}\"")
		c.addLine("if [ -z \"${_t}\" ]; then break; fi") // https://stackoverflow.com/a/13864829 (didn't work with +x (probably due to the underscore of the variable)).
		c.addLine("_l=$(expr ${_l} + 1)")
		c.addLine("done")
		c.addLine("echo ${_l}")
		c.addLine("}")

		c.sliceLenHelperSet = true
	}
	helper := c.nextHelperVar()
	c.VarAssignment(helper, fmt.Sprintf("$(_sl %s)", name)) // TODO: Find out if using varEvaluationString here is a good idea because name might not be a variable.

	return c.VarEvaluation(helper, valueUsed)
}

func (c *converter) Group(value string, valueUsed bool) (string, error) {
	return fmt.Sprintf("(%s)", value), nil
}

func (c *converter) FuncCall(name string, args []string, valueType parser.ValueType, valueUsed bool) (string, error) {
	returnsValue := valueType.DataType() != parser.DATA_TYPE_VOID
	argsCopy := args

	for i, arg := range argsCopy {
		argsCopy[i] = fmt.Sprintf("\"%s\"", arg)
	}
	call := fmt.Sprintf("%s %s", name, strings.Join(argsCopy, " "))

	if returnsValue && valueUsed {
		helper := c.nextHelperVar()
		call = fmt.Sprintf("$(%s)", call)

		c.VarDefinition(helper, call)
		return c.VarEvaluation(helper, valueUsed)
	}
	c.addLine(call)
	return "", nil
}

func (c *converter) AppCall(calls []transpiler.AppCall, valueUsed bool) (string, error) {
	callsCopy := calls
	callStrings := []string{}

	for _, call := range callsCopy {
		argsCopy := call.Args()

		for j, arg := range argsCopy {
			argsCopy[j] = fmt.Sprintf("\"%s\"", arg)
		}
		space := ""

		if len(argsCopy) > 0 {
			space = " "
		}
		callStrings = append(callStrings, fmt.Sprintf("%s%s%s", call.Name(), space, strings.Join(argsCopy, " ")))
	}
	callString := strings.Join(callStrings, " | ")

	if valueUsed {
		callString = fmt.Sprintf("$(%s)", callString)
		helper := c.nextHelperVar()

		c.VarDefinition(helper, callString)
		return c.VarEvaluation(helper, valueUsed)
	}
	c.addLine(callString)
	return "", nil
}

func (c *converter) Input(prompt string, valueUsed bool) (string, error) {
	helper := c.nextHelperVar()

	if len(prompt) > 0 {
		prompt = fmt.Sprintf(" -p \"%s\"", prompt)
	}
	c.addLine(fmt.Sprintf("read%s %s", prompt, helper))
	return c.VarEvaluation(helper, valueUsed)
}

func (c *converter) ifStart(condition string, startWord string) error {
	c.addLine(fmt.Sprintf("%s [ %s -eq %s ]; then", startWord, condition, c.BoolToString(true)))
	return nil
}

func (c *converter) addLine(line string) {
	c.code = append(c.code, line)
}

func (c *converter) nextHelperVar() string {
	helperVar := fmt.Sprintf("h%d", c.varCounter)
	c.varCounter++

	return helperVar
}
