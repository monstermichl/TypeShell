package batch

import (
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
	sliceLenSetHelper     helperName = "_sls"  // Slice length set
	sliceLenGetHelper     helperName = "_slg"  // Slice length get
	sliceAssignmentHelper helperName = "_sah"  // Slice assignment
	sliceCopyHelper       helperName = "_sch"  // Slice copy
	stringSubscriptHelper helperName = "_stsh" // String subscript
	stringLengthHelper    helperName = "_stlh" // String length
	stringEscapeHelper    helperName = "_seh"  // String escape
	echoHelper            helperName = "_ech"  // Echo
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
	startCode                     []string
	helperCode                    []string
	globalCode                    []string
	previousFunctionName          string
	functionsCode                 [][]string
	endCode                       []string
	varCounter                    int
	ifCounter                     int
	forCounter                    int
	endLabels                     []string
	funcs                         []funcInfo
	funcCounter                   int
	fors                          []forInfo
	ifs                           []ifInfo
	lfSet                         bool
	appCallHelperRequired         bool
	readHelperRequired            bool
	sliceAssignmentHelperRequired bool
	sliceCopyHelperRequired       bool
	sliceLenSetHelperRequired     bool
	sliceLenGetHelperRequired     bool
	stringSubscriptHelperRequired bool
	stringLenHelperRequired       bool
	fileWriteHelperRequired       bool
	echoHelperRequired            bool
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

func forLabel(count int) string {
	return fmt.Sprintf(":_f%d", count)
}

func (c *converter) StringToString(value string) string {
	c.addLf()

	// Replace "\n" with Batch newline value.
	return strings.ReplaceAll(value, "\n", "!LF!")
}

func (c *converter) Dump() (string, error) {
	functionsCode := []string{}

	// Collect function code.
	for i := len(c.functionsCode) - 1; i >= 0; i-- {
		functionsCode = append(functionsCode, c.functionsCode[i]...)
	}
	allCode := append([]string{}, c.startCode...)
	allCode = append(allCode, c.helperCode...)
	allCode = append(allCode, functionsCode...)
	allCode = append(allCode, c.globalCode...)
	allCode = append(allCode, c.endCode...)
	allCode = append(allCode, "") // Add a terminating newline.
	indent := 0

	for i, line := range allCode {
		if strings.HasPrefix(line, ")") {
			indent--
		}
		if strings.HasSuffix(line, "(") {
			indent++
		}
		if indent < 0 {
			indent = 0
		}
		allCode[i] = fmt.Sprintf("%s%s", strings.Repeat(" ", indent), line)
	}
	return strings.Join(allCode, "\r\n"), nil
}

func (c *converter) ProgramStart() error {
	c.addStartLine("@echo off")
	c.addStartLine("setlocal EnableDelayedExpansion")
	c.addStartLine("setlocal")
	c.addStartLine(`set "_e=0"`)

	return nil
}

func (c *converter) ProgramEnd() error {
	if c.fileWriteHelperRequired {
		c.addHelper("file write", fileWriteHelper,
			`set "_a=>"`,
			fmt.Sprintf(`if "%%2" equ "%s" set "_a=>>"`, transpiler.BoolToString(true)),
			fmt.Sprintf(`for /f "delims=" %%%%i in ("!%s!") do (`, funcArgVar(0)),
			`for /f "delims=" %%j in ('echo %%i!_a! %~1') do (`,
			"rem", // NOP (https://stackoverflow.com/a/7191952). The for-loop is just required to expand !_a! delayed.
			`)`,
			`set "_a=>>"`,
			")",
		)
	}

	if c.appCallHelperRequired {
		c.addHelper("app call", appCallHelper,
			`set "_h="`,
			`set "_te="`, // Temporary error variable.
			fmt.Sprintf(`for /f "delims=" %%%%i in ('cmd /V:ON /C "!%s! & echo ^!errorlevel^!"') do (`, funcArgVar(0)), // Use the carets in ^!errorlevel^! to make sure errorlevel is expanded within cmd.
			`if defined _h set "_h=!_h!!LF!"`,
			`set "_h=!_h!!_te!"`,
			`set _te=%%i`,
			")",
		)
	}

	if c.readHelperRequired {
		c.addHelper("read", fileReadHelper,
			`set "_h="`,
			`for /f "delims=" %%i in (%~1) do (`,
			`if defined _h set "_h=!_h!!LF!"`,
			`set "_h=!_h!%%i"`,
			")",
		)
	}

	if c.sliceCopyHelperRequired {
		c.sliceLenGetHelperRequired = true
		c.sliceLenSetHelperRequired = true
		c.sliceAssignmentHelperRequired = true

		// %1: Destination slice
		// %2: Source slice
		c.addHelper("slice copy", sliceCopyHelper,
			`set "_i=0"`,
			c.callFuncString(sliceLenGetHelper, []string{}, "%2"),
			":_sch_loop",
			`if "!_i!" lss "!_len!" (`,
			`for /f "delims=" %%i in ("%2_!_i!") do set "_v=!%%i!"`,
			c.sliceAssignmentString("!%1!", "!_i!", "!_v!", false),
			`set /A "_i=!_i!+1"`,
			"goto :_sch_loop",
			")",
			c.callFuncString(sliceLenSetHelper, []string{}, "!%1!", "!_i!"),
		)
	}

	if c.sliceAssignmentHelperRequired {
		c.sliceLenGetHelperRequired = true
		c.sliceLenSetHelperRequired = true
		c.sliceAssignmentHelperRequired = true

		// %1: Slice name
		// %2: Assigned index
		// %3: Default value
		// arg0: Assigned value
		c.addHelper("slice assignment", sliceAssignmentHelper,
			c.callFuncString(sliceLenGetHelper, []string{}, "!%1!"), // Get current slice length.
			`set "_i=!_len!"`,
			":_sah_loop",
			`if "!_i!" lss "%2" (`,
			c.sliceAssignmentString("!%1!", "!_i!", "%3", false),
			`set /A "_i=!_i!+1"`,
			"goto :_sah_loop",
			") else (",
			`set /A "_len=%2+1"`,
			c.callFuncString(sliceLenSetHelper, []string{}, "!%1!", "!_len!"),
			")",
			c.sliceAssignmentString("!%1!", "%2", fmt.Sprintf("!%s!", funcArgVar(0)), false),
		)
	}

	if c.sliceLenSetHelperRequired {
		c.addHelper("slice length set", sliceLenSetHelper,
			`set "%1_len=%2"`,
		)
	}

	if c.sliceLenGetHelperRequired {
		c.addHelper("slice length get", sliceLenGetHelper,
			`set "_len=!%1_len!"`,
		)
	}

	if c.stringSubscriptHelperRequired {
		c.addHelper("string subscript", stringSubscriptHelper,
			`set /A "_sh=(%2-%1)+1"`,
			fmt.Sprintf(`set "_sub=!%s:~%%1,%%_sh%%!"`, funcArgVar(0)), // https://stackoverflow.com/a/636391
		)
	}

	if c.stringLenHelperRequired {
		c.addHelper("string length", stringLengthHelper,
			"set _l=0",
			":_stlhl",
			fmt.Sprintf(`if "!%s:~%%_l%%!" equ "" goto :_stlhle`, funcArgVar(0)), // https://www.geeksforgeeks.org/batch-script-string-length/
			`set /A "_l=%_l%+1"`,
			"goto :_stlhl",
			":_stlhle",
		)
	}

	if c.echoHelperRequired {
		v := c.varEvaluationString(funcArgVar(0), true)
		c.addHelper("echo", echoHelper,
			fmt.Sprintf(`if "%s" neq "" (echo %s) else echo.`, v, v), // echo. could be problematic (see discussion: https://stackoverflow.com/a/20691061).
		)
	}
	c.addEndLine(":end")
	c.addEndLine("endlocal & exit /B %_e%")
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

func (c *converter) SliceAssignment(name string, index string, value string, defaultValue string, global bool) error {
	c.sliceAssignmentHelperRequired = true
	c.callFunc(sliceAssignmentHelper, []string{value}, c.varName(name, global), index, defaultValue)
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
	c.addLine("exit /B")
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
	label := c.nextForLabel()
	c.nextEndLabel()
	c.fors = append(c.fors, forInfo{
		label: label,
	})
	c.addLine(fmt.Sprintf(`set "%s="`, c.mustCurrentForVar()))
	c.addLine(label)
	return nil
}

func (c *converter) ForIncrementStart() error {
	c.addLine(fmt.Sprintf(`if defined %s (`, c.mustCurrentForVar()))
	return nil
}

func (c *converter) ForIncrementEnd() error {
	c.addLine(")")
	c.addLine(fmt.Sprintf(`set "%s=1"`, c.mustCurrentForVar()))
	return nil
}

func (c *converter) ForCondition(condition string) error {
	c.addLine(fmt.Sprintf(`if "%s" equ "%s" (`, condition, transpiler.BoolToString(true)))
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
	c.addLine(fmt.Sprintf("goto %s", c.mustCurrentForLabel()))
	return nil
}

func (c *converter) Print(values []string) error {
	c.callEchoFunc(values...)
	return nil
}

func (c *converter) Panic(value string) error {
	c.callEchoFunc(value)
	c.addLine(`set "_e=1"`)
	c.addLine("goto :end")
	return nil
}

func (c *converter) WriteFile(path string, content string, append string) error {
	c.fileWriteHelperRequired = true

	// Use global variable to pass content to write file helper because Batch doesn't
	// support newline passing because it splits arguments at newlines.
	c.callFunc(fileWriteHelper, []string{content}, path, append)

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
			fmt.Sprintf("if %s equ %s (%s) else %s",
				expr,
				transpiler.BoolToString(true),
				c.varAssignmentString(helper, transpiler.BoolToString(false), false),
				c.varAssignmentString(helper, transpiler.BoolToString(true), false),
			),
		)
	default:
		return "", fmt.Errorf(`unknown unary operator "%s"`, operator)
	}
	return c.VarEvaluation(helper, valueUsed, false)
}

func (c *converter) BinaryOperation(left string, operator parser.BinaryOperator, right string, valueType parser.ValueType, valueUsed bool) (string, error) {
	helper := c.nextHelperVar()
	notAllowedError := func() (string, error) {
		return "", fmt.Errorf("binary operation %s is not allowed on type %s", operator, valueType.String())
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
		c.addLine(fmt.Sprintf(`set /A "%s=%s%s%s"`, c.varName(helper, false), left, operator, right))
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
	quote := ""

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
			quote = `"` // Strings shall be quoted.
		}
	}

	if len(operatorString) == 0 {
		return "", fmt.Errorf("comparison %s is not allowed on type %s", operator, valueType.String())
	}
	helper := c.nextHelperVar()
	c.addLine(
		fmt.Sprintf(`if %s%s%s %s %s%s%s (%s) else %s`,
			quote,
			left,
			quote,
			operatorString,
			quote,
			right,
			quote,
			c.varAssignmentString(helper, transpiler.BoolToString(true), false),
			c.varAssignmentString(helper, transpiler.BoolToString(false), false),
		),
	)
	return c.VarEvaluation(helper, valueUsed, false)
}

func (c *converter) LogicalOperation(left string, operator parser.LogicalOperator, right string, valueType parser.ValueType, valueUsed bool) (string, error) {
	var line string
	trueString := transpiler.BoolToString(true)
	falseString := transpiler.BoolToString(false)
	helper := c.nextHelperVar()
	trueAssignment := c.varAssignmentString(helper, trueString, false)
	falseAssignment := c.varAssignmentString(helper, falseString, false)

	switch operator {
	case parser.LOGICAL_OPERATOR_AND:
		line = fmt.Sprintf(`if %s equ %s (if %s equ %s (%s) else %s) else %s`,
			left,
			trueString,
			right,
			trueString,
			trueAssignment,
			falseAssignment,
			falseAssignment,
		)
	case parser.LOGICAL_OPERATOR_OR:
		line = fmt.Sprintf(`if %s equ %s (%s) else if %s equ %s (%s) else %s`,
			left,
			trueString,
			trueAssignment,
			right,
			trueString,
			trueAssignment,
			falseAssignment,
		)
	default:
		return "", fmt.Errorf(`unknown logical operator "%s"`, operator)
	}

	c.addLine(line)
	return c.VarEvaluation(helper, valueUsed, false)
}

func (c *converter) VarEvaluation(name string, valueUsed bool, global bool) (string, error) {
	return c.varEvaluationString(name, global), nil
}

func (c *converter) SliceInstantiation(values []string, valueUsed bool) (string, error) {
	c.addLine(`set /A "_dvc=!_dvc!+1"`) // Dynamic variable counter.
	helper := c.nextHelperVar()
	c.VarAssignment(helper, "_dv!_dvc!", false)

	c.sliceAssignmentHelperRequired = true
	c.callFunc(sliceLenSetHelper, []string{}, c.varEvaluationString(helper, false), strconv.Itoa(len(values)))

	// Init slice values.
	for i, value := range values {
		c.addLine(c.sliceAssignmentString(c.varEvaluationString(helper, false), strconv.Itoa(i), value, false))
	}
	return c.varEvaluationString(helper, false), nil
}

func (c *converter) SliceEvaluation(name string, index string, valueUsed bool) (string, error) {
	helper := c.nextHelperVar()

	// A for-loop is required because the evaluation wouldn't work with the following code as expected.
	// It always put out "_h0_0" instead of "4".
	//
	// set a1=_h0
	// set _h0_0=4
	// set x=!a1!_0
	// echo !x!
	//
	// ChatGpt (yes, I'm a bit ashamed about it but I used it) told me the following:
	// In your Batch script, the issue arises because set x=!a1!_0 does not expand !a1! before
	// accessing [0]. Instead, it treats !a1!_0 as a literal string, so x is assigned the value
	// _h0_0, not 4. Batch scripts do not support indirect variable expansion in a straightforward
	// way. However, you can work around this by using for /f to evaluate the variable dynamically
	c.addLine(
		fmt.Sprintf(`for /f "delims=" %%%%i in ("%s_%s") do set "%s=!%%%%i!"`,
			name, // TODO: Is global flag even required here? Because value is already passed to function.
			index,
			c.varName(helper, false),
		),
	)
	return c.VarEvaluation(helper, valueUsed, false)
}

func (c *converter) SliceLen(name string, valueUsed bool) (string, error) {
	helper := c.nextHelperVar()

	c.sliceLenGetHelperRequired = true
	c.callFunc(sliceLenGetHelper, []string{}, name)
	c.VarAssignment(helper, c.varEvaluationString("_len", true), false)

	return c.VarEvaluation(helper, valueUsed, false)
}

func (c *converter) StringSubscript(value string, startIndex string, endIndex string, valueUsed bool) (string, error) {
	helper := c.nextHelperVar()
	c.stringSubscriptHelperRequired = true

	c.callFunc(stringSubscriptHelper, []string{value}, startIndex, endIndex)
	c.VarAssignment(helper, c.varEvaluationString("_sub", true), false) // TODO: Is global flag even required here? Because value is already passed to function.

	return c.varEvaluationString(helper, false), nil
}

func (c *converter) StringLen(value string, valueUsed bool) (string, error) {
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

	// Make sure return values contain as many values as expected.
	for len(returnValues) < len(returnTypes) {
		returnValues = append(returnValues, "")
	}
	return returnValues, nil
}

func (c *converter) AppCall(calls []transpiler.AppCall, valueUsed bool) ([]string, error) {
	callsCopy := calls
	callStrings := []string{}

	for _, call := range callsCopy {
		argsCopy := call.Args()

		for j, arg := range argsCopy {
			// If argument is a variable or contains whitespaces, quote it.
			if strings.HasPrefix(arg, "%") || len(strings.Split(arg, " ")) > 1 {
				arg = fmt.Sprintf(`"%s"`, arg)
			}
			argsCopy[j] = arg
		}
		space := ""

		if len(argsCopy) > 0 {
			space = " "
		}
		callStrings = append(callStrings, fmt.Sprintf("%s%s%s", call.Name(), space, strings.Join(argsCopy, " ")))
	}

	if valueUsed {
		helper1 := c.nextHelperVar()
		helper2 := c.nextHelperVar()
		c.appCallHelperRequired = true

		c.addLf()
		c.callFunc(appCallHelper, []string{strings.Join(callStrings, " | ")})
		c.VarAssignment(helper1, c.varEvaluationString("_h", true), false)

		eval, err := c.VarEvaluation(helper1, valueUsed, false)

		if err != nil {
			return nil, err
		}
		c.VarAssignment(helper2, c.varEvaluationString("_te", true), false)
		return []string{eval, "", c.varEvaluationString(helper2, false)}, nil // TODO: Return stderr (https://github.com/monstermichl/TypeShell/issues/28).
	}
	c.addLine(fmt.Sprintf("call %s", strings.Join(callStrings, " | ")))
	return []string{"", "", "0"}, nil
}

func (c *converter) Input(prompt string, valueUsed bool) (string, error) {
	helper := c.nextHelperVar()
	c.addLine(fmt.Sprintf(`set /p "%s=%s"`, helper, prompt))
	return c.VarEvaluation(helper, valueUsed, false)
}

func (c *converter) Copy(destination string, source string, valueUsed bool, global bool) (string, error) {
	c.sliceCopyHelperRequired = true
	c.callFunc(sliceCopyHelper, []string{}, c.varName(destination, global), source)

	c.callFunc(sliceLenGetHelper, []string{}, c.varEvaluationString(destination, global))
	return c.varEvaluationString("_len", true), nil
}

func (c *converter) Exists(path string, valueUsed bool) (string, error) {
	helper := c.nextHelperVar()

	c.addLine(fmt.Sprintf(`if exist "%s" (%s) else %s`,
		path,
		c.varAssignmentString(helper, transpiler.BoolToString(true), false),
		c.varAssignmentString(helper, transpiler.BoolToString(false), false),
	))
	return c.VarEvaluation(helper, valueUsed, false)
}

func (c *converter) ReadFile(path string, valueUsed bool) (string, error) {
	helper := c.nextHelperVar()
	c.readHelperRequired = true

	c.addLf()
	c.callFunc(fileReadHelper, []string{}, path)
	c.VarAssignment(helper, c.varEvaluationString("_h", true), false)

	return c.VarEvaluation(helper, valueUsed, false)
}

func (c *converter) callFuncString(name string, globalArgs []string, args ...string) string {
	for i, arg := range globalArgs {
		c.VarAssignment(funcArgVar(i), arg, true)
	}
	return fmt.Sprintf("call :%s %s", strings.TrimLeft(name, ":"), strings.Join(args, " "))
}

func (c *converter) callFunc(name string, globalArgs []string, args ...string) {
	c.addLine(c.callFuncString(name, globalArgs, args...))
}

func (c *converter) callEchoFunc(values ...string) {
	c.echoHelperRequired = true
	c.addLine(c.callFuncString(echoHelper, []string{strings.Join(values, " ")}))
}

func (c *converter) varName(name string, global bool) string {
	if c.inFunction() && !global {
		name = fmt.Sprintf("f%d_%s", c.funcCounter, name)
	}
	return name
}

func (c *converter) varAssignmentString(name string, value string, global bool) string {
	return fmt.Sprintf(`set "%s=%s"`, c.varName(name, global), value)
}

func (c *converter) varEvaluationString(name string, global bool) string {
	return fmt.Sprintf("!%s!", c.varName(name, global))
}

func (c *converter) ifStart(condition string, startAddition string) error {
	c.addLine(fmt.Sprintf(`%sif "%s" equ "%s" (`, startAddition, condition, transpiler.BoolToString(true)))
	return nil
}

func (c *converter) addStartLine(line string) {
	c.startCode = append(c.startCode, line)
}

func (c *converter) addHelperLine(line string) {
	c.helperCode = append(c.helperCode, line)
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

func (c *converter) mustCurrentForLabel() string {
	return forLabel(c.forCounter - 1)
}

func (c *converter) mustCurrentForVar() string {
	return fmt.Sprintf("_fv%d", c.forCounter-1)
}

func (c *converter) nextForLabel() string {
	label := forLabel(c.forCounter)
	c.forCounter++

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

	c.addHelperLine(fmt.Sprintf(":: global %s helper begin", helperType))
	c.addHelperLine(fmt.Sprintf("goto %s", endLabel))
	c.addHelperLine(fmt.Sprintf(":%s", label))

	for _, line := range code {
		c.addHelperLine(line)
	}
	c.addHelperLine("exit /B")
	c.addHelperLine(endLabel)
	c.addHelperLine(fmt.Sprintf(":: global %s helper end", helperType))
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
	return fmt.Sprintf(`set "%s_%s=%s"`, name, index, value)
}

func (c *converter) addLf() {
	if !c.lfSet {
		c.addStartLine("(set LF=^") // https://stackoverflow.com/a/60389149
		c.addStartLine("")
		c.addStartLine(")")

		c.lfSet = true
	}
}
