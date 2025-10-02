# TypeShell
TypeShell is a Go-like programming language that transpiles down to Batch or Bash.

```cmd
rem Transpile helloworld.tsh to Batch and Bash and write the scripts to the current directory.
tsh.exe -i helloworld.tsh -t batch -t bash -o .
```

## Example
```golang
// helloworld.tsh

func hello() string {
    return "hello"
}

func buildGreeting(p string) string {
    return hello() + " " + p
}

greeting := buildGreeting("world")

print(greeting) // Prints "Hello World" to the console.
```

## Basics
### Variables
Supported variable types are *bool*, *int*, *string* and *error*.

```golang
// Variable definition with default value.
var a int
var b, c int
```

```golang
// Variable definition with assigned value.
var a int = 5
var b, c int = divisionWithRemainder(5, 2)
```

```golang
// Variable definition via grouping.
var (
    a = 5
    b, c int = divisionWithRemainder(5, 2)
)
```

```golang
// Variable definition short form.
a := 5
b, c := 5, 6
d, e := divisionWithRemainder(5, 2)
```

### Constants
```golang
// Constant definition.
const a = 0
const b, c = 1, 2
```

```golang
// Constant definition via grouping.
const (
    a = -1
    b = iota
    c
)
```

### Control flow
```golang
// If-statement.
if a < 5 {
    // Do something.
} else if > 5 {
    // Do something.
} else {
    // Do something else.
}
```

```golang
// Switch-statement.
switch a {
case 5:
    // Do something.
case 6:
    // Do something.
default:
    // Do something else.
}

switch true {
case a == 5:
    // Do something.
default:
    // Do something else.
}

switch {
case false:
    // Do something.
default:
    // Do something else.
}
```

```golang
// For-loop.
for {
    // Do something.
}

for a == 5 {
    // Do something.
}

for i := 0; i < 5; i++ {
    // Do something.
}

// For-range-loop (supported for slices, strings and integers).
for i, v := range s {
    // Do something.
}

for i := range 5 {
    // Do something.
}
```

### Functions
```golang
// Function definition.
func division(a int, b int) int {
    return a / b
}

func divisionWithRemainder(a int, b int) (int, int) {
    return division(a, b), a % b
}

// Function call.
division(5, 2)
```

```golang
// Define function with optional parameters.
func sum(vals ...int) int {
    sum := 0

    for _, val := range vals {
        sum += val
    }
    return sum
}

// Call function with an arbitrary amount of arguments.
sum(1, 2, 3, 4)
```

### Slices
```golang
// Slice creation.
s := []int{}
```

```golang
// Slice creation with values.
s := []int{1, 2, 3}
```

```golang
// Slice assignment.
s[0] = 10
```

```golang
// Slice evaluation.
v := s[0]
```

```golang
// Slice length.
l := len(s)
```

```golang
// Slice iteration.
for i := 0; i < len(s); i++ {
    v := s[i]
}
```

### Types
TypeShell supports the definition of types.

```golang
// Define a type.
type myType int

var a myType
a = myType(24)
```

```golang
// Define an alias.
type myType = int

var a myType
a = 24
```

```golang
// Define a struct.
type myStruct struct {
    a    string
    b    string
    c, d int
}
```

### Structs
```golang
// Define a simple struct.
type myStruct struct {
    greeting string
    planet   string
}

// Define a new struct variable.
var s myStruct = myStruct{greeting: "Hello", planet: "World"}

// Assign a value to a struct field.
s.planet = "Mars"

// Retrieve values from the struct.
print(s.greeting, s.planet)
```

#### Nested structs
```golang
// Define structs.
type myNestedStruct struct {
    planet string
}

type myStruct struct {
    greeting string
    info     myNestedStruct
}

// Define a new struct variable.
var s myStruct = myStruct{
    greeting: "Hello",
    info: myNestedStruct{
        planet: "World",
    },
}

// Assign a value to a nested struct field.
s.info.planet = "Mars"

// Retrieve values from the structs.
print(s.greeting, s.info.planet)
```

#### Function receivers
Functions can have a struct as receiver. This way the function can be called directly on the struct and the corresponding struct is passed to the function automatically.

```golang
// Define a struct.
type myStruct struct {
    greeting string
}

// Define a function with the struct as receiver.
func (s myStruct) greet() {
	print(s.greeting)
}

s := myStruct{greeting: "Hello"}

// Call struct function.
s.greet()
```

#### Struct pointers
Structs can also be passed to functions as pointers. This holds true for parameters and receivers.

```golang
// Define a structs.
type myStruct struct {
    greeting string
}

type mySecondStruct struct {
	planet string
}

// Define a function with a pointer receiver and a pointer parameter.
func (g *myStruct) greet(p *mySecondStruct) {
	print(g.greeting, p.planet)
	
	// Update receiver and parameter value.
	g.greeting = "Bye"
	p.planet = "Mars"
}

s1 := myStruct{greeting: "Hello"}
s2 := mySecondStruct{planet: "World"}

// Call struct function.
s1.greet(s2) // Prints "Hello World".
s1.greet(s2) // Prints "Bye Mars".
```

### Programs/Scripts
```golang
// Programs/Scripts are called by stating the name preceded by an @.
@dir("/b")
```

```golang
// Similar to Bash/Batch, the output can be piped into another program/script.
@dir("/b") | @sort("/r")
```

```golang
// To capture the output, just assign the call chain to variables.
stdout, stderr, code := @dir("/b") | @sort("/r")
```

```golang
// To specify the path to a program/script, a string literal is used.
@`helper\dir.bat`("/b") // Equivalent to @"helper\\dir.bat"("/b")
```

### Imports
TypeShell does not support import of packages like Go does, but it supports single file imports. If an imported script is not a [standard "library" script](https://github.com/monstermichl/TypeShell/tree/main/std), an alias needs to be defined.

```golang
// Relative file import.
import hp "helper.tsh"

hp.HelperFunc()
```

```golang
// Standard "library" import.
import (
    "strings"
)

print(strings.Contains("Hello World", "World")) // Prints 1.
```

It's also possible to include files from a remote source. If the loaded code is a Go-file, the package statement is automatically removed to make it usable as a TypeShell file.
```golang
// Remote import.
import (
    ext "https://exampleserver.com/somegofile.tsh"
)

ext.SomeFunction()
```

### Builtin
```golang
// Returns the length of a slice or a string.
len(slice)
len(str)
```

```golang
// Prints the passed arguments to stdout.
print(arg0, arg1, ...)
```

```golang
// Asks for user input.
input()
input(promptString)
```

```golang
// Copies values from srcSlice to dstSlice. Returns the copied length.
copy(dstSlice, srcSlice)
```

```golang
// Reads file content.
read(path)
```

```golang
// Writes file content.
write(path, contentString)
write(path, contentString, appendBool)
```

```golang
// Checks if a path exists.
exists(path)
```

```golang
// Converts an integer to a string.
itoa(str)
```

```golang
// Kills the program with an error.
panic(err)
```

### ⚠️ Unsafe
It's possible to add native code directly to the output by using the *unsafe*-builtin. However, this should be avoided if possible as it can introduce unwanted side-effects. The builtin's first argument must be a string literal which contains the executable code. All other arguments can be of type *bool*, *int* or *string*. The transpiler parses the string literal and replaces all placeholders (e.g. {0}) with the corresponding positional argument at transpile time.

To pass data into the native code, placeholders with either only the positional number (e.g. "{0}") or "i:" followed by the positional number (e.g. "{i:0}") should be used.

To get data out of the native code, placeholders with "o:" followed by the positional number (e.g. "{o:0}") must be used. **IMPORTANT**: The passed expression must be a variable.

In the following example, the variable *file* is passed as input to the native Batch code, while the variable *date* is passed as output argument. After the execution, *date* holds the value evaluated by the native code.

```golang
file := "test.tsh"
var date string

unsafe(`for %%U in ({i:0}) do (set "{o:1}=%%~tU")`, file, date)
```

Results in

```batch
set "file_0=test.tsh"
set "date_0="
for %%U in (!file_0!) do (set "date_0=%%~tU")
```

## Caveats
### Condition evaluation
In contrast to many other programming languages, TypeShell evaluates all conditions before the actual statement. This is done to handle the limitations of Batch/Bash. HINT: This is also true for switch-evaluations since switchs are internally converted to ifs.

```golang
if a == 1 && b == 1 {
    // Do something.
}
```

```golang
// How it's handled internally.
h1 := a == 1
h2 := b == 2
h3 := a && b

if h3 {
    // Do something.
}
```

### Error and nil
In TypeShell error is just a string type and nil is an empty string. However, they are still supported to provide developers with the possibility to use the typical Go workflow of error checking.

```golang
err := func()

if err != nil {
    // Do something.
}
```

### Functions
- Functions must be defined before being used.
- Recursions are not supported yet.

### Slices
If a slice index does not exist on assignment, it and its intermediate indices are created.

```golang
s := []string{"Hello"}

s[2] = "World"

print(s[0]) // Prints "Hello".
print(s[1]) // Prints "".
print(s[2]) // Prints "World".
```

### Types
Type definitions which result in slices are not supported yet.

```golang
type myType []string // Not supported.
```

### Pointers
Pointers are only supported as function parameters and only for structs.

### Performance
Try to avoid function calls like *len* in for-conditions if possible since they are evaluated for each iteration.

```golang
// Don't do this.
for i := 0; i < len(s); i++ {
    ...
}

// Do this instead.
lenS := len(s)

for i := 0; i < lenS; i++ {
    ...
}
```

## Visual Studio Code
There is no extension for VSCode yet. However, since the code is very Go-like, adding the ".tsh" extension to the settings should serve as a first workaround.
- Open VSCode.
- Go to File -> Preferences -> Settings.
- Seach for "file associations".
- Add "*.tsh" to the list and associate it with Go.
