# TypeShell
TypeShell is a Go-like programming language that transpiles down to Batch or Bash.

## Basics
### Variables
Supported variable types are *bool*, *int* and *string*.

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
if a == 5 {
    // Do something.
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

### Operators
```golang
// Arithmetical operators.
2 + 5
2 - 5
2 * 5
2 / 5
2 % 5
```

```golang
// Compare operators.
2 == 5
2 != 5
2 > 5
2 >= 5
2 < 5
2 <= 5
```

```golang
// Logical operators.
!true
true && true
true || true
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
// To capture the output, just assign the call chain to a variable.
x := @ls("-a") | @sort()
```

### Builtin
```golang
// Returns the length of a slice.
len(slice)
```

```golang
// Prints the passed arguments to stdout.
print()
```

```golang
// Asks for user input.
input()
input("name: ")
```

```golang
// Copies values from srcSlice to dstSlice. Returns the copied length.
copy(dstSlice, srcSlice)
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
