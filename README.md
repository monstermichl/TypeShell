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
```

```golang
// Variable definition short form.
a := 5
a, b := 5, 6
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

### Input/Output
```golang
// Read user input.
x := input("number: ")
```

```golang
// Output.
print(x)
```
