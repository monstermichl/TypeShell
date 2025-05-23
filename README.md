# TypeShell
TypeShell is a Go-like programming language that transpiles down to Batch or Bash.

## Basics
### Variables
Supported variable types are *bool*, *int*, *string* and *error*.

```golang
// Variable definition with default value.
var a int
var a, b int
```

```golang
// Variable definition with assigned value.
var a int = 5
var a, b int = divisionWithRemainder(5, 2)
```

```golang
// Variable definition short form.
a := 5
a, b := 5, 6
a, b := divisionWithRemainder(5, 2)
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

### Programs
```golang
// Programs are called by stating the program name preceded by an @.
@ls("-a")
```

```golang
// Similar to Bash/Batch, the output can be piped into another program.
@ls("-a") | @sort()
```

```golang
// To capture the output, just assign the call chain to variables.
stdout, stderr, code := @ls("-a") | @sort()
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
In contrast to many other programming languages, TypeShell evaluates all conditions before the actual statement. This is done to handle the limitations of Batch/Bash.

```golang
if a == 1 && b == 1 {
    // Do something.
}
```

```golang
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
