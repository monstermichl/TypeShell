# TypeShell
TypeShell is a Go-like programming language that transpiles down to Batch or Bash.

## Basics
### Variables
Supported variable types are *bool*, *int* and *string*.

```golang
// Variable definition with default value.
var x int
```

```golang
// Variable definition with assigned value.
var x int = 5
```

```golang
// Variable definition short form.
x := 5
```

### Control flow
```golang
// If-statement.
if x == 5 {
    // Do something.
}
```

```golang
// For-loop.
for x == 5 {
    // Do something.
}
```

### Functions
```golang
// Function definition.
func sum(a int, b int) int {
    return a + b
} 
```

```golang
// Function call.
sum(2, 5)
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
true || true
true && true
```

### Programs
```golang
// Programs get called by stating the program name followed by
// curly brackets. Arguments to the program are passed within
// the curly brackets.
ls{"-a"}
```

```golang
// Similar to Bash/Batch, the output can be piped into another
// program.
ls{"-a"} | sort{}
```

```golang
// To capture the output, just the call chain to a variable.
x := ls{"-a"} | sort{}
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
