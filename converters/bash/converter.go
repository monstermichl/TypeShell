package bash

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/monstermichl/typeshell/parser"
	"github.com/monstermichl/typeshell/transpiler"
)

type funcInfo struct {
	name string
}

type converter struct {
	interpreter                   string
	startCode                     []string
	code                          []string
	varCounter                    int
	funcs                         []funcInfo
	funcCounter                   int
	sliceLenHelperRequired        bool
	sliceCopyHelperRequired       bool
	stringSubscriptHelperRequired bool
}

func New() *converter {
	return &converter{
		interpreter: "/bin/bash",
		code:        []string{},
	}
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
	return strings.Join([]string{
		strings.Join(c.startCode, "\n"),
		strings.Join(c.code, "\n"),
	}, "\n"), nil
}

func (c *converter) ProgramStart() error {
	c.addStartLine(fmt.Sprintf("#!%s", c.interpreter))
	return nil
}

func (c *converter) ProgramEnd() error {
	if c.sliceCopyHelperRequired {
		c.sliceLenHelperRequired = true

		c.addHelper("slice copy", "_sch",
			"local _i=0",
			"local _l=$(_slh ${2})",
			"local _n=$(eval \"echo \\${${1}}\")",
			"while [ ${_i} -lt ${_l} ]; do",
			fmt.Sprintf("local _v=%s", c.sliceEvaluationString("2", "${_i}", false)),
			c.sliceAssignmentString("${_n}", "${_i}", "${_v}", false),
			"_i=$(expr ${_i} + 1)",
			"done",
		)
	}

	if c.sliceLenHelperRequired {
		c.addHelper("slice length", "_slh",
			"local _l=0",
			"while true; do",
			"eval \"local _t=\\${$1_${_l}}\"",
			"if [ -z \"${_t}\" ]; then break; fi", // https://stackoverflow.com/a/13864829 (didn't work with +x (probably due to the underscore of the variable)).
			"_l=$(expr ${_l} + 1)",
			"done",
			"echo ${_l}",
		)
	}

	if c.stringSubscriptHelperRequired {
		c.addHelper("substring", "_ssh",
			"echo \"${1}\" | cut -c $(expr ${2} \\+ 1)", // Cut index starts at 1, therefore 1 must be added to 0-based subscript.
		)
	}
	return nil
}

func (c *converter) VarDefinition(name string, value string, global bool) error {
	return c.VarAssignment(name, value, global)
}

func (c *converter) VarAssignment(name string, value string, global bool) error {
	length := len(value)

	if length > 0 {
		if string(value[length-1]) != "\"" {
			value = fmt.Sprintf("%s\"", value)
		}
		if string(value[0]) != "\"" {
			value = fmt.Sprintf("\"%s", value)
		}
	}
	c.addLine(fmt.Sprintf("%s=%s", c.varName(name, global), value))
	return nil
}

func (c *converter) SliceAssignment(name string, index string, value string, global bool) error {
	// TODO: Find out if global is correctly used here.
	c.addLine(c.sliceAssignmentString(c.varEvaluationString(name, global), index, value, global)) // TODO: Find out if using varEvaluationString here is a good idea because name might not be a variable.
	return nil
}

func (c *converter) FuncStart(name string, params []string, returnTypes []parser.ValueType) error {
	c.funcs = append(c.funcs, funcInfo{
		name: name,
	})
	c.funcCounter++
	c.addLine(fmt.Sprintf("%s() {", name))

	for i, param := range params {
		c.VarAssignment(param, fmt.Sprintf("$%d", i+1), false)
	}
	return nil
}

func (c *converter) FuncEnd() error {
	c.addLine("}")

	lastIndex := len(c.funcs) - 1
	c.funcs = slices.Delete(c.funcs, lastIndex, lastIndex+1)
	return nil
}

func (c *converter) Return(values []transpiler.ReturnValue) error {
	for i, value := range values {
		c.VarDefinition(fmt.Sprintf("_rv%d", i), value.Value(), true)
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

func (c *converter) UnaryOperation(expr string, operator parser.UnaryOperator, valueType parser.ValueType, valueUsed bool) (string, error) {
	helper := c.nextHelperVar()

	switch operator {
	case parser.UNARY_OPERATOR_NEGATE:
		c.VarAssignment(helper,
			fmt.Sprintf("$(if [ \"%s\" -eq \"%s\" ]; then echo %s; else echo %s; fi)",
				expr,
				c.BoolToString(true),
				c.BoolToString(false),
				c.BoolToString(true),
			),
			false,
		)
	default:
		return "", fmt.Errorf("unknown unary operator \"%s\"", operator)
	}
	return c.VarEvaluation(helper, valueUsed, false)
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
		c.VarAssignment(helper, fmt.Sprintf("$(expr %s \\%s %s)", left, operator, right), false) // Backslash is required for * operator to prevent pattern expansion (https://www.shell-tips.com/bash/math-arithmetic-calculation/#using-the-expr-command-line).
	case parser.DATA_TYPE_STRING:
		switch operator {
		case parser.BINARY_OPERATOR_ADDITION:
			c.VarAssignment(helper, fmt.Sprintf("\"%s%s\"", left, right), false)
		default:
			return notAllowedError()
		}
	default:
		return notAllowedError()
	}
	return c.VarEvaluation(helper, valueUsed, false)
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

	c.VarAssignment(
		helper,
		fmt.Sprintf("$(if [ \"%s\" %s \"%s\" ]; then echo %s; else echo %s; fi)",
			left,
			operatorString,
			right,
			c.BoolToString(true),
			c.BoolToString(false),
		),
		false,
	)
	return c.VarEvaluation(helper, valueUsed, false)
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

	c.VarAssignment(
		helper,
		fmt.Sprintf("$(if [ \"%s\" -eq \"%s\" ] %s [ \"%s\" -eq \"%s\" ]; then echo %s; else echo %s; fi)",
			left,
			trueString,
			operatorString,
			right,
			trueString,
			trueString,
			c.BoolToString(false),
		),
		false,
	)
	return c.VarEvaluation(helper, valueUsed, false)
}

func (c *converter) VarEvaluation(name string, valueUsed bool, global bool) (string, error) {
	return c.varEvaluationString(name, global), nil
}

func (c *converter) SliceInstantiation(values []string, valueUsed bool) (string, error) {
	helper := c.nextHelperVar()

	for i, value := range values {
		c.addLine(c.sliceAssignmentString(helper, strconv.Itoa(i), value, false))
	}
	return helper, nil
}

func (c *converter) SliceEvaluation(name string, index string, valueUsed bool, global bool) (string, error) {
	helper := c.nextHelperVar()
	c.VarAssignment(
		helper,
		c.sliceEvaluationString(name, index, global),
		false,
	) // TODO: Find out if using varEvaluationString here is a good idea because name might not be a variable.

	return c.VarEvaluation(helper, valueUsed, false)
}

func (c *converter) SliceLen(name string, valueUsed bool, global bool) (string, error) {
	helper := c.nextHelperVar()
	c.sliceLenHelperRequired = true
	// TODO: Handle global flag.
	c.VarAssignment(helper, c.sliceLenString(name), false) // TODO: Find out if using varEvaluationString here is a good idea because name might not be a variable.

	return c.VarEvaluation(helper, valueUsed, false)
}

func (c *converter) StringSubscript(name string, index string, valueUsed bool, global bool) (string, error) {
	helper := c.nextHelperVar()

	c.VarAssignment(helper, fmt.Sprintf("$(_ssh %s %s)", c.varEvaluationString(name, global), index), false) // https://www.baeldung.com/linux/bash-substring#1-using-thecut-command
	c.stringSubscriptHelperRequired = true

	return c.varEvaluationString(helper, false), nil
}

func (c *converter) StringLen(value string, valueUsed bool, global bool) (string, error) {
	helper := c.nextHelperVar()

	c.VarAssignment(helper, value, false)
	c.VarAssignment(helper, fmt.Sprintf("${#%s}", c.varName(helper, false)), false)

	return c.VarEvaluation(helper, valueUsed, false)
}

func (c *converter) Group(value string, valueUsed bool) (string, error) {
	return fmt.Sprintf("(%s)", value), nil
}

func (c *converter) FuncCall(name string, args []string, returnTypes []parser.ValueType, valueUsed bool) ([]string, error) {
	argsCopy := args

	for i, arg := range argsCopy {
		argsCopy[i] = fmt.Sprintf("\"%s\"", arg)
	}
	returnValues := []string{}
	c.addLine(fmt.Sprintf("%s %s", name, strings.Join(argsCopy, " "))) // TODO: Remove general parameter quoting.

	if valueUsed {
		for i, _ := range returnTypes {
			helper := c.nextHelperVar()

			c.VarDefinition(helper, c.varEvaluationString(fmt.Sprintf("_rv%d", i), true), false)
			eval, _ := c.VarEvaluation(helper, valueUsed, false)
			returnValues = append(returnValues, eval)
		}
	}
	return returnValues, nil
}

func (c *converter) AppCall(calls []transpiler.AppCall, valueUsed bool) (string, error) {
	callsCopy := calls
	callStrings := []string{}

	for _, call := range callsCopy {
		argsCopy := call.Args()

		for j, arg := range argsCopy {
			// argument is a variable or contains whitespaces, quote it.
			if strings.HasPrefix(arg, "$") || len(strings.Split(arg, " ")) > 1 {
				arg = fmt.Sprintf("\"%s\"", arg)
			}
			argsCopy[j] = arg
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

		c.VarDefinition(helper, callString, false)
		return c.VarEvaluation(helper, valueUsed, false)
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
	return c.VarEvaluation(helper, valueUsed, false)
}

func (c *converter) Copy(destination string, source string, valueUsed bool, global bool) (string, error) {
	destination = c.varName(destination, global)

	c.addLine(fmt.Sprintf("_sch %s %s", destination, source))
	c.sliceCopyHelperRequired = true

	helper := c.nextHelperVar()
	c.VarAssignment(helper, c.sliceLenString(c.varEvaluationString(destination, true)), false)

	return c.varEvaluationString(helper, false), nil
}

func (c *converter) varName(name string, global bool) string {
	if c.inFunction() && !global {
		name = fmt.Sprintf("f%d_%s", c.funcCounter, name)
	}
	return name
}

func (c *converter) varEvaluationString(name string, global bool) string {
	return fmt.Sprintf("${%s}", c.varName(name, global))
}

func (c *converter) sliceAssignmentString(name string, index string, value string, global bool) string {
	// TODO: Handle global flag.
	return fmt.Sprintf("eval %s_%s=\"%s\"", name, index, value)
}

func (c *converter) sliceEvaluationString(name string, index string, global bool) string {
	return fmt.Sprintf("$(eval \"echo \\${%s_%s}\")",
		c.varEvaluationString(name, global),
		index,
	)
}

func (c *converter) sliceLenString(name string) string {
	return fmt.Sprintf("$(_slh %s)", name)
}

func (c *converter) ifStart(condition string, startWord string) error {
	c.addLine(fmt.Sprintf("%s [ %s -eq %s ]; then", startWord, condition, c.BoolToString(true)))
	return nil
}

func (c *converter) inFunction() bool {
	return len(c.funcs) > 0
}

func (c *converter) addHelper(helperType string, functionName string, code ...string) {
	c.addStartLine(fmt.Sprintf("# global %s helper", helperType))
	c.addStartLine(fmt.Sprintf("%s() {", functionName))

	for _, line := range code {
		c.addStartLine(line)
	}
	c.addStartLine("}")
}

func (c *converter) addStartLine(line string) {
	c.startCode = append(c.startCode, line)
}

func (c *converter) addLine(line string) {
	c.code = append(c.code, line)
}

func (c *converter) nextHelperVar() string {
	helperVar := fmt.Sprintf("h%d", c.varCounter)
	c.varCounter++

	return helperVar
}
