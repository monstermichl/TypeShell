# TypeShell changelog

## v0.2.0 - xxxx-xx-xx
### New
- Add type definition support (https://github.com/monstermichl/TypeShell/issues/45).
- Add constants support (https://github.com/monstermichl/TypeShell/issues/41).
- Add struct support (https://github.com/monstermichl/TypeShell/issues/51).
- Allow passing of structs to functions as pointer (https://github.com/monstermichl/TypeShell/issues/63).
- Add function receiver support (https://github.com/monstermichl/TypeShell/issues/64).
- Add support for remote imports (https://github.com/monstermichl/TypeShell/issues/46).
- Include std in release (https://github.com/monstermichl/TypeShell/issues/43).
- Add for-range support for numbers (https://github.com/monstermichl/TypeShell/issues/51).
- Make subscript available on string literals (https://github.com/monstermichl/TypeShell/issues/15).
- Add support for one type specification for consecutive parameters with same type (https://github.com/monstermichl/TypeShell/issues/38).
- Add os std library.

### Improvements
- Performance improvements of std strings.

### Fixes
- Multiple imports lead to error (https://github.com/monstermichl/TypeShell/issues/40).
- If a variable is re-defined with other variables, it's not considered a new variable (https://github.com/monstermichl/TypeShell/issues/42).
- Print with spaces only shows echo off message (Batch) (https://github.com/monstermichl/TypeShell/issues/53).
- Subtraction is not considered subtraction if number comes directly after minus (https://github.com/monstermichl/TypeShell/issues/62).
- Using imported variables/constants doesn't work (https://github.com/monstermichl/TypeShell/issues/49).
  
## v0.1.1 - 2025-07-05
### New
- Allow execution of non-local/non-PATH executables/scripts (see https://github.com/monstermichl/TypeShell?tab=readme-ov-file#programsscripts) (https://github.com/monstermichl/TypeShell/issues/35).

### Fixes
- Handle exclamation marks correctly in string literals (file read/writes don't work correctly yet) (Batch) (https://github.com/monstermichl/TypeShell/issues/33).

## v0.1.1 - 2025-06-21
### New
- Add raw string support (https://github.com/monstermichl/TypeShell/issues/31).
- Add support for single variable in for-range loops (https://github.com/monstermichl/TypeShell/issues/29).
- Add Split function to std/strings.tsh.
- Add find example for Windows and Linux.

### Fixes
- Len on empty string (Batch) (https://github.com/monstermichl/TypeShell/issues/32).
- Substring newline handling (Bash).
- Negation handling (https://github.com/monstermichl/TypeShell/issues/30).

## v0.1.0 - 2025-06-19
### New
- First release
