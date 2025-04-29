package batch

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/monstermichl/typeshell/parser"
	"github.com/monstermichl/typeshell/transpiler"
)

type helperName = string

const (
	appCallHelper         helperName = "_ach"  // App call
	fileReadHelper        helperName = "_frh"  // File write
	fileWriteHelper       helperName = "_fwh"  // File read
	sliceCopyHelper       helperName = "_sch"  // Slice copy
	sliceLengthHelper     helperName = "_sllh" // Slice length
	stringSubscriptHelper helperName = "_stsh" // String subscript
	stringLengthHelper    helperName = "_stlh" // String length
)

type funcInfo struct {
	name string
}

type forInfo struct {
	label string
}

type ifInfo struct {
	label string
}

type converter struct {
	globalCode                    []string
	previousFunctionName          string
	functionsCode                 [][]string
	endCode                       []string
	varCounter                    int
	ifCounter                     int
	whileCounter                  int
	endLabels                     []string
	funcs                         []funcInfo
	funcCounter                   int
	fors                          []forInfo
	ifs                           []ifInfo
	lfSet                         bool
	appCallHelperRequired         bool
	readHelperRequired            bool
	sliceLenHelperRequired        bool
	sliceCopyHelperRequired       bool
	stringSubscriptHelperRequired bool
	stringLenHelperRequired       bool
	fileWriteHelperRequired       bool
}

func New() *converter {
	return &converter{}
}

func returnValVar(subscript int) string {
	return fmt.Sprintf("_rv%d", subscript)
}

func funcArgVar(subscript int) string {
	return fmt.Sprintf("_fa%d", subscript)
}

func escapeString(s string) string {
	escapes := []string{"%", "^", "&", "<", ">", "|", "'", "`", ",", ";", "=", "(", ")"} // https://www.robvanderwoude.com/escapechars.php
	updated := ""

	for i := range s {
		char := string(s[i])

		if slices.Contains(escapes, char) {
			char = fmt.Sprintf("^%s", char)
		}
		updated = fmt.Sprintf("%s%s", updated, char)
	}
	return updated
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
	c.addLf()

	// Replace "\n" with Batch newline value.
	return strings.ReplaceAll(value, "\\n", "!LF!")
}

func (c *converter) Dump() (string, error) {
	functionsCode := []string{}

	// Append function code in the reverse order because Batch only
	// knows "function"-labels that are further down the code.
	for i := len(c.functionsCode) - 1; i >= 0; i-- {
		functionsCode = append(functionsCode, c.functionsCode[i]...)
	}
	return strings.Join([]string{
		strings.Join(c.globalCode, "\n"),
		strings.Join(functionsCode, "\n"),
		strings.Join(c.endCode, "\n"),
	}, "\n"), nil
}

func (c *converter) ProgramStart() error {
	c.addLine("@echo off")
	c.addLine("setlocal EnableDelayedExpansion")
	c.addLine("setlocal")
	return nil
}

func (c *converter) ProgramEnd() error {
	if c.fileWriteHelperRequired {
		c.addHelper("file write", fileWriteHelper,
			"set _a=",
			fmt.Sprintf("for /f \"delims=\" %%%%i in (\"!%s!\") do (", funcArgVar(0)),
			"if defined _a (echo %%i >> %~1) else (",
			"echo %%i > %~1",
			"set _a=1",
			")",
			")",
		)
	}

	if c.appCallHelperRequired {
		c.addHelper("app call", appCallHelper,
			"set _h=",
			fmt.Sprintf("for /f \"delims=\" %%%%i in ('call !%s!') do (", funcArgVar(0)),
			"if defined _h set \"_h=!_h!!LF!\"",
			"set \"_h=!_h!%%i\"",
			")",
		)
	}

	if c.readHelperRequired {
		c.addHelper("read", fileReadHelper,
			"set _h=",
			"for /f \"delims=\" %%i in (%~1) do (",
			"if defined _h set \"_h=!_h!!LF!\"",
			"set \"_h=!_h!%%i\"",
			")",
		)
	}

	if c.sliceCopyHelperRequired {
		c.addHelper("slice copy", sliceCopyHelper,
			"set _i=0",
			"call :_sllh %2", // Call slice length helper.
			":_schl",
			"if !_i! lss !_l! (",
			"for /f \"delims=\" %%i in (\"%2[!_i!]\") do set \"_v=!%%i!\"",
			c.sliceAssignmentString("!%1!", "!_i!", "!_v!", false),
			"set /A \"_i=!_i!+1\"",
			"goto :_schl",
			")",
		)
		c.sliceLenHelperRequired = true
	}

	if c.sliceLenHelperRequired {
		c.addHelper("slice length", sliceLengthHelper,
			"set _l=0",
			":_sllhl",
			"if not defined %1[%_l%] goto :_sllhle",
			"set /A \"_l=%_l%+1\"",
			"goto :_sllhl",
			":_sllhle",
		)
	}

	if c.stringSubscriptHelperRequired {
		c.addHelper("string subscript", stringSubscriptHelper,
			fmt.Sprintf("set \"_sub=!%s:~%%1%%,1!\"", funcArgVar(0)), // https://stackoverflow.com/a/636391
		)
	}

	if c.stringLenHelperRequired {
		c.addHelper("string length", stringLengthHelper,
			"set _l=0",
			":_stlhl",
			fmt.Sprintf("if \"!%s:~%%_l%%!\" equ \"\" goto :_stlhle", funcArgVar(0)), // https://www.geeksforgeeks.org/batch-script-string-length/
			"set /A \"_l=%_l%+1\"",
			"goto :_stlhl",
			":_stlhle",
		)
	}

	c.addEndLine("endlocal")
	return nil
}

func (c *converter) VarDefinition(name string, value string, global bool) error {
	c.addLine(c.varAssignmentString(name, value, global))
	return nil
}

func (c *converter) VarAssignment(name string, value string, global bool) error {
	c.addLine(c.varAssignmentString(name, value, global))
	return nil
}

func (c *converter) SliceAssignment(name string, index string, value string, global bool) error {
	// TODO: Find out if global is used correctly here.
	c.addLine(c.sliceAssignmentString(c.varEvaluationString(name, global), index, value, global)) // TODO: Find out if using varEvaluationString here is a good idea because name might not be a variable.
	return nil
}

func (c *converter) FuncStart(name string, params []string, returnTypes []parser.ValueType) error {
	c.funcCounter++
	c.funcs = append(c.funcs, funcInfo{
		name: name,
	})
	c.addLine(fmt.Sprintf(":: %s function begin", name))
	c.addLine(fmt.Sprintf("goto :_eo_%s", name))
	c.addLine(fmt.Sprintf(":%s", name))

	for i, param := range params {
		c.addLine(c.varAssignmentString(param, c.varEvaluationString(funcArgVar(i), true), false))
	}
	return nil
}

func (c *converter) FuncEnd() error {
	name := c.mustCurrentFuncInfo().name

	c.addLine(fmt.Sprintf(":_ret_%s", name))
	c.addLine("exit /B 0")
	c.addLine(fmt.Sprintf(":_eo_%s", name))
	c.addLine(fmt.Sprintf(":: %s function end", name))

	lastIndex := len(c.funcs) - 1
	c.funcs = slices.Delete(c.funcs, lastIndex, lastIndex+1)
	return nil
}

func (c *converter) Return(values []transpiler.ReturnValue) error {
	currFunc := c.mustCurrentFuncInfo()

	for i, value := range values {
		c.VarDefinition(returnValVar(i), value.Value(), true)
	}
	c.addLine(fmt.Sprintf("goto :_ret_%s", currFunc.name))
	return nil
}

func (c *converter) IfStart(condition string) error {
	c.ifs = append(c.ifs, ifInfo{
		label: c.nextIfLabel(),
	})
	return c.ifStart(condition, "")
}

func (c *converter) IfEnd() error {
	label := c.mustCurrentIfInfo().label

	c.addLine(fmt.Sprintf("goto %s", label))
	c.addLine(")")
	c.addLine(label)

	lastIndex := len(c.ifs) - 1
	c.ifs = slices.Delete(c.ifs, lastIndex, lastIndex+1)

	return nil
}

func (c *converter) ElseIfStart(condition string) error {
	c.addLine(fmt.Sprintf("goto %s", c.mustCurrentIfInfo().label))
	return c.ifStart(condition, ") else ")
}

func (c *converter) ElseIfEnd() error {
	return nil
}

func (c *converter) ElseStart() error {
	c.addLine(fmt.Sprintf("goto %s", c.mustCurrentIfInfo().label))
	c.addLine(") else (")
	return nil
}

func (c *converter) ElseEnd() error {
	return nil
}

func (c *converter) ForStart() error {
	label := c.nextWhileLabel()
	c.nextEndLabel()
	c.fors = append(c.fors, forInfo{
		label: label,
	})
	c.addLine(label)
	return nil
}

func (c *converter) ForCondition(condition string) error {
	c.addLine(fmt.Sprintf("if \"%s\" equ \"%s\" (", condition, c.BoolToString(true)))
	return nil
}

func (c *converter) ForEnd() error {
	w := c.fors[len(c.fors)-1]

	c.addLine(fmt.Sprintf("goto %s", w.label))
	c.addLine(")")

	endLabel := c.popEndLabel()
	lastIndex := len(c.fors) - 1
	c.fors = slices.Delete(c.fors, lastIndex, lastIndex+1)

	c.addLine(endLabel)

	return nil
}

func (c *converter) Break() error {
	c.addLine(fmt.Sprintf("goto %s", c.mustCurrentEndLabel()))
	return nil
}

func (c *converter) Continue() error {
	return errors.New("continue has not been implemented yet")
}

func (c *converter) Print(values []string) error {
	c.addLine(fmt.Sprintf("echo %s", escapeString(strings.Join(values, " "))))
	return nil
}

func (c *converter) WriteFile(path string, content string, append string) error {
	c.fileWriteHelperRequired = true
	// TODO: Consider append.
	// Use global variable to pass content to write file helper because Batch doesn't
	// support newline passing because it splits arguments at newlines.
	c.callFunc(fileWriteHelper, []string{content}, path)

	return nil
}

func (c *converter) Nop() error {
	c.addLine("rem No operation")
	return nil
}

func (c *converter) UnaryOperation(expr string, operator parser.UnaryOperator, valueType parser.ValueType, valueUsed bool) (string, error) {
	helper := c.nextHelperVar()

	switch operator {
	case parser.UNARY_OPERATOR_NEGATE:
		c.addLine(
			fmt.Sprintf("if \"%s\" equ \"%s\" (%s) else %s",
				expr,
				c.BoolToString(true),
				c.varAssignmentString(helper, c.BoolToString(false), false),
				c.varAssignmentString(helper, c.BoolToString(true), false),
			),
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
			parser.BINARY_OPERATOR_ADDITION,
			parser.BINARY_OPERATOR_SUBTRACTION:
			// These operations are fine.
		case parser.BINARY_OPERATOR_MODULO:
			operator = fmt.Sprintf("%%%s", operator) // Modulo needs to be escaped as "%" is used to dereference variables in Batch.
		default:
			return notAllowedError()
		}
		c.addLine(fmt.Sprintf("set /A \"%s=%s%s%s\"", c.varName(helper, false), left, operator, right))
	case parser.DATA_TYPE_STRING:
		switch operator {
		case parser.BINARY_OPERATOR_ADDITION:
			c.VarAssignment(helper, fmt.Sprintf("%s%s", left, right), false)
		default:
			return notAllowedError()
		}
	default:
		return notAllowedError()
	}
	return c.VarEvaluation(helper, valueUsed, false)
}

func (c *converter) Comparison(left string, operator parser.CompareOperator, right string, valueType parser.ValueType, valueUsed bool) (string, error) {
	EQUAL_OPERATOR := "equ"
	NOT_EQUAL_OPERATOR := "neq"

	var operatorString string

	if !valueType.IsSlice() {
		switch valueType.DataType() {
		case parser.DATA_TYPE_BOOLEAN:
			switch operator {
			case parser.COMPARE_OPERATOR_EQUAL:
				operatorString = EQUAL_OPERATOR
			case parser.COMPARE_OPERATOR_NOT_EQUAL:
				operatorString = NOT_EQUAL_OPERATOR
			}
		case parser.DATA_TYPE_INTEGER:
			switch operator {
			case parser.COMPARE_OPERATOR_EQUAL:
				operatorString = EQUAL_OPERATOR
			case parser.COMPARE_OPERATOR_NOT_EQUAL:
				operatorString = NOT_EQUAL_OPERATOR
			case parser.COMPARE_OPERATOR_GREATER:
				operatorString = "gtr"
			case parser.COMPARE_OPERATOR_GREATER_OR_EQUAL:
				operatorString = "geq"
			case parser.COMPARE_OPERATOR_LESS:
				operatorString = "lss"
			case parser.COMPARE_OPERATOR_LESS_OR_EQUAL:
				operatorString = "leq"
			}
		case parser.DATA_TYPE_STRING:
			switch operator {
			case parser.COMPARE_OPERATOR_EQUAL:
				operatorString = EQUAL_OPERATOR
			case parser.COMPARE_OPERATOR_NOT_EQUAL:
				operatorString = NOT_EQUAL_OPERATOR
			}
		}
	}

	if len(operatorString) == 0 {
		return "", fmt.Errorf("comparison %s is not allowed on type %s", operator, valueType.ToString())
	}
	helper := c.nextHelperVar()
	c.addLine(
		fmt.Sprintf("if \"%s\" %s \"%s\" (%s) else %s",
			left,
			operatorString,
			right,
			c.varAssignmentString(helper, c.BoolToString(true), false),
			c.varAssignmentString(helper, c.BoolToString(false), false),
		),
	)
	return c.VarEvaluation(helper, valueUsed, false)
}

func (c *converter) LogicalOperation(left string, operator parser.LogicalOperator, right string, valueType parser.ValueType, valueUsed bool) (string, error) {
	var line string
	trueString := c.BoolToString(true)
	falseString := c.BoolToString(false)
	helper := c.nextHelperVar()
	trueAssignment := c.varAssignmentString(helper, trueString, false)
	falseAssignment := c.varAssignmentString(helper, falseString, false)

	switch operator {
	case parser.LOGICAL_OPERATOR_AND:
		line = fmt.Sprintf("if \"%s\" equ \"%s\" if \"%s\" equ \"%s\" (%s) else %s",
			left,
			trueString,
			right,
			trueString,
			trueAssignment,
			falseAssignment,
		)
	case parser.LOGICAL_OPERATOR_OR:
		line = fmt.Sprintf("if \"%s\" equ \"%s\" (%s) else if \"%s\" equ \"%s\" (%s) else %s",
			left,
			trueString,
			trueAssignment,
			right,
			trueString,
			trueAssignment,
			falseAssignment,
		)
	default:
		return "", fmt.Errorf("unknown logical operator \"%s\"", operator)
	}

	c.addLine(line)
	return c.VarEvaluation(helper, valueUsed, false)
}

func (c *converter) VarEvaluation(name string, valueUsed bool, global bool) (string, error) {
	return c.varEvaluationString(name, global), nil
}

func (c *converter) SliceInstantiation(values []string, valueUsed bool) (string, error) {
	helper := c.nextHelperVar()

	// Init slice values.
	for i, value := range values {
		c.addLine(c.sliceAssignmentString(helper, strconv.Itoa(i), value, false))
	}
	return helper, nil
}

func (c *converter) SliceEvaluation(name string, index string, valueUsed bool, global bool) (string, error) {
	helper := c.nextHelperVar()

	// A for-loop is required because the evaluation wouldn't work with the following code as expected.
	// It always put out "_h0[0]" instead of "4".
	//
	// set a1=_h0
	// set _h0[0]=4
	// set x=!a1![0]
	// echo !x!
	//
	// ChatGpt (yes, I'm a bit ashamed about it but I used it) told me the following:
	// In your Batch script, the issue arises because set x=!a1![0] does not expand !a1! before
	// accessing [0]. Instead, it treats !a1![0] as a literal string, so x is assigned the value
	// _h0[0], not 4. Batch scripts do not support indirect variable expansion in a straightforward
	// way. However, you can work around this by using for /f to evaluate the variable dynamically
	c.addLine(
		fmt.Sprintf("for /f \"delims=\" %%%%i in (\"%s[%s]\") do set \"%s=!%%%%i!\"",
			name, // TODO: Is global flag even required here? Because value is already passed to function.
			index,
			c.varName(helper, false),
		),
	)
	return c.VarEvaluation(helper, valueUsed, false)
}

func (c *converter) SliceLen(name string, valueUsed bool, global bool) (string, error) {
	helper := c.nextHelperVar()
	c.sliceLenHelperRequired = true

	c.callFunc(sliceLengthHelper, []string{}, name)
	c.VarAssignment(helper, c.varEvaluationString("_l", true), false)

	return c.VarEvaluation(helper, valueUsed, false)
}

func (c *converter) StringSubscript(value string, index string, valueUsed bool, global bool) (string, error) {
	helper := c.nextHelperVar()
	c.stringSubscriptHelperRequired = true

	c.callFunc(stringSubscriptHelper, []string{value}, index)
	c.VarAssignment(helper, c.varEvaluationString("_sub", true), false) // TODO: Is global flag even required here? Because value is already passed to function.

	return c.varEvaluationString(helper, false), nil
}

func (c *converter) StringLen(value string, valueUsed bool, global bool) (string, error) {
	helper := c.nextHelperVar()
	c.stringLenHelperRequired = true

	c.callFunc(stringLengthHelper, []string{value})
	c.VarAssignment(helper, c.varEvaluationString("_l", true), false)

	return c.VarEvaluation(helper, valueUsed, false)
}

func (c *converter) Group(value string, valueUsed bool) (string, error) {
	return fmt.Sprintf("(%s)", value), nil
}

func (c *converter) FuncCall(name string, args []string, returnTypes []parser.ValueType, valueUsed bool) ([]string, error) {
	returnValues := []string{}
	c.callFunc(name, args)

	if valueUsed {
		for i := range returnTypes {
			helper := c.nextHelperVar()
			c.VarDefinition(helper, c.varEvaluationString(returnValVar(i), true), false)
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
			argsCopy[j] = fmt.Sprintf("\"%s\"", arg)
		}
		space := ""

		if len(argsCopy) > 0 {
			space = " "
		}
		callStrings = append(callStrings, fmt.Sprintf("%s%s%s", call.Name(), space, strings.Join(argsCopy, " ")))
	}

	if valueUsed {
		helper := c.nextHelperVar()
		c.appCallHelperRequired = true

		c.addLf()
		c.callFunc(appCallHelper, []string{strings.Join(callStrings, " ^| ")})
		c.VarAssignment(helper, c.varEvaluationString("_h", true), false)

		return c.VarEvaluation(helper, valueUsed, false)
	}
	c.addLine(fmt.Sprintf("call %s", strings.Join(callStrings, " | ")))
	return "", nil
}

func (c *converter) Input(prompt string, valueUsed bool) (string, error) {
	helper := c.nextHelperVar()
	c.addLine(fmt.Sprintf("set /p \"%s=%s\"", helper, prompt))
	return c.VarEvaluation(helper, valueUsed, false)
}

func (c *converter) Copy(destination string, source string, valueUsed bool, global bool) (string, error) {
	c.sliceCopyHelperRequired = true
	c.callFunc(sliceCopyHelper, []string{}, c.varName(destination, global), source)

	return c.varEvaluationString("_l", true), nil
}

func (c *converter) ReadFile(path string, valueUsed bool) (string, error) {
	helper := c.nextHelperVar()
	c.readHelperRequired = true

	c.addLf()
	c.callFunc(fileReadHelper, []string{}, path)
	c.VarAssignment(helper, c.varEvaluationString("_h", true), false)

	return c.VarEvaluation(helper, valueUsed, false)
}

func (c *converter) callFunc(name string, globalArgs []string, args ...string) {
	for i, arg := range globalArgs {
		c.VarAssignment(funcArgVar(i), arg, true)
	}
	c.addLine(fmt.Sprintf("call :%s %s", strings.TrimLeft(name, ":"), strings.Join(args, " ")))
}

func (c *converter) varName(name string, global bool) string {
	if c.inFunction() && !global {
		name = fmt.Sprintf("f%d_%s", c.funcCounter, name)
	}
	return name
}

func (c *converter) varAssignmentString(name string, value string, global bool) string {
	return fmt.Sprintf("set \"%s=%s\"", c.varName(name, global), value)
}

func (c *converter) varEvaluationString(name string, global bool) string {
	return fmt.Sprintf("!%s!", c.varName(name, global))
}

func (c *converter) ifStart(condition string, startAddition string) error {
	c.addLine(fmt.Sprintf("%sif \"%s\" equ \"%s\" (", startAddition, condition, c.BoolToString(true)))
	return nil
}

func (c *converter) addLine(line string) {
	if c.inFunction() {
		currFunc := c.mustCurrentFuncInfo()
		currFuncName := currFunc.name
		prevFuncName := c.previousFunctionName

		if currFuncName != prevFuncName {
			c.previousFunctionName = currFuncName
			c.functionsCode = append(c.functionsCode, []string{})
		}
		index := len(c.functionsCode) - 1
		c.functionsCode[index] = append(c.functionsCode[index], line)
	} else {
		c.globalCode = append(c.globalCode, line)
	}
}

func (c *converter) addEndLine(line string) {
	c.endCode = append(c.endCode, line)
}

func (c *converter) nextWhileLabel() string {
	label := fmt.Sprintf(":_w%d", c.whileCounter)
	c.whileCounter++

	return label
}

func (c *converter) nextIfLabel() string {
	label := fmt.Sprintf(":_i%d", c.ifCounter)
	c.ifCounter++

	return label
}

func (c *converter) addHelper(helperType string, label string, code ...string) {
	label = strings.TrimLeft(label, ":")
	endLabel := fmt.Sprintf(":_eo_%s", label)

	c.addEndLine(fmt.Sprintf(":: global %s helper begin", helperType))
	c.addEndLine(fmt.Sprintf("goto %s", endLabel))
	c.addEndLine(fmt.Sprintf(":%s", label))

	for _, line := range code {
		c.addEndLine(line)
	}
	c.addEndLine(endLabel)
	c.addEndLine("exit /B 0")
	c.addEndLine(fmt.Sprintf(":: global %s helper end", helperType))
}

func (c *converter) inFunction() bool {
	return len(c.funcs) > 0
}

func (c *converter) mustCurrentFuncInfo() funcInfo {
	return c.funcs[len(c.funcs)-1]
}

func (c *converter) mustCurrentIfInfo() ifInfo {
	return c.ifs[len(c.ifs)-1]
}

func (c *converter) mustCurrentEndLabel() string {
	return c.endLabels[len(c.endLabels)-1]
}

func (c *converter) popEndLabel() string {
	lastIndex := len(c.endLabels) - 1
	label := c.mustCurrentEndLabel()
	c.endLabels = slices.Delete(c.endLabels, lastIndex, lastIndex+1)

	return label
}

func (c *converter) nextEndLabel() string {
	c.endLabels = append(c.endLabels, fmt.Sprintf(":_e%d", len(c.endLabels)))
	return c.mustCurrentEndLabel()
}

func (c *converter) nextHelperVar() string {
	helperVar := fmt.Sprintf("_h%d", c.varCounter)
	c.varCounter++

	return helperVar
}

func (c *converter) sliceAssignmentString(name string, index string, value string, global bool) string {
	// TODO: Is global flag even required here? Because value is already passed to function.
	return fmt.Sprintf("set \"%s[%s]=%s\"", name, index, value)
}

func (c *converter) addLf() {
	if !c.lfSet {
		c.addLine("(set LF=^") // https://stackoverflow.com/a/60389149
		c.addLine("")
		c.addLine(")")

		c.lfSet = true
	}
}
