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
// constant definition via grouping.
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

// For-range-loop (supported for slices and strings).
for i, v := range s {
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
```

```golang
// Function call.
sum(2, 5)
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
import (
    hp "helper.tsh"
)

hp.HelperFunc()
```

```golang
// Standard "library" import.
import (
    "strings"
)

print(strings.Contains("Hello World", "World")) // Prints 1.
```

### Type declarations
TypeShell supports the declaration of types. However, types which result in slices are not supported yet.

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

## Visual Studio Code
There is no extension for VSCode yet. However, since the code is very Go-like, adding the ".tsh" extension to the settings should serve as a first workaround.
- Open VSCode.
- Go to File -> Preferences -> Settings.
- Seach for "file associations".
- Add "*.tsh" to the list and associate it with Go.
