package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func testSingleImportSuccess(t *testing.T, transpilerFunc transpilerCalloutFunc) {
	transpilerFunc(t, func(dir string) (string, error) {
		return `
			import "strings"
			print(strings.Contains("Hello World", "Wor"))
		`, nil
	}, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1", output)
	})
}

func testSeveralSingleImportsSuccess(t *testing.T, transpilerFunc transpilerCalloutFunc) {
	transpilerFunc(t, func(dir string) (string, error) {
		return `
			import strings1 "strings"
			import strings2 "strings"

			print(strings1.Contains("Hello World", "Wor"))
		`, nil
	}, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1", output)
	})
}

func testMultiImportSuccess(t *testing.T, transpilerFunc transpilerCalloutFunc) {
	transpilerFunc(t, func(dir string) (string, error) {
		return `import (
			strings1 "strings"
			strings2 "strings"
			)
			print(strings1.Contains("Hello World", "Wor"))
			print(strings2.HasPrefix("Hello World", "Hel"))
		`, nil
	}, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1\n1", output)
	})
}

func testWildlyMixedImportsSuccess(t *testing.T, transpilerFunc transpilerCalloutFunc) {
	transpilerFunc(t, func(dir string) (string, error) {
		return `import (
			strings1 "strings"
			)
			import strings2 "strings"
			import strings3 "strings"
			import (
				strings4 "strings"
				"strings"
			)
			print(strings1.Contains("Hello World", "Wor"))
			print(strings2.HasPrefix("Hello World", "Hel"))
			print(strings.Contains("Hello World", "Hel"))
			print(strings3.Contains("Hello World", "Bel"))
			print(strings4.Contains("Hello World", "Bel"))
		`, nil
	}, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1\n1\n1\n0\n0", output)
	})
}

func testImportsFromExternalSourceSuccess(t *testing.T, transpilerFunc transpilerCalloutFunc) {
	transpilerFunc(t, func(dir string) (string, error) {
		return `
			import strings "https://raw.githubusercontent.com/monstermichl/TypeShell/refs/heads/main/std/strings.tsh"

			print(strings.Contains("Hello World", "Hel"))
		`, nil
	}, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1", output)
	})
}

func testImportVariableSuccess(t *testing.T, transpilerFunc transpilerCalloutFunc) {
	value := "Hello World"

	transpilerFunc(t, func(dir string) (string, error) {
		importFile := "import.tsh"
		variable := "PublicVariable"
		err := os.WriteFile(filepath.Join(dir, importFile), []byte(variable+` := "`+value+`"`), 0700)

		if err != nil {
			return "", err
		}
		return `
			import imp "` + importFile + `"
			print(imp.` + variable + `)
		`, nil
	}, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, value, output)
	})
}

func testImportedSliceAssignmentSuccess(t *testing.T, transpilerFunc transpilerCalloutFunc) {
	values := []string{"Hello", "World"}

	transpilerFunc(t, func(dir string) (string, error) {
		importFile := "import.tsh"
		variable := "PublicVariable"
		err := os.WriteFile(filepath.Join(dir, importFile), []byte(variable+` := []string{"`+values[0]+`", "`+values[1]+`"}`), 0700)

		if err != nil {
			return "", err
		}
		return `
			import imp "` + importFile + `"
			print(imp.` + variable + `[0])
			imp.` + variable + `[0] = "Bye"
			print(imp.` + variable + `[0])
		`, nil
	}, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "Hello\nBye", output)
	})
}

func testImportSimpleTypeSuccess(t *testing.T, transpilerFunc transpilerCalloutFunc) {
	transpilerFunc(t, func(dir string) (string, error) {
		importFile := "import.tsh"
		t := "MyType"
		err := os.WriteFile(filepath.Join(dir, importFile), []byte(`
			type `+t+` string`,
		), 0700)

		if err != nil {
			return "", err
		}
		return `
			import imp "` + importFile + `"
			s := imp.` + t + `("Hello")
			print(s)
		`, nil
	}, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "Hello", output)
	})
}

func testImportStructTypeSuccess(t *testing.T, transpilerFunc transpilerCalloutFunc) {
	transpilerFunc(t, func(dir string) (string, error) {
		importFile := "import.tsh"
		t := "MyType"
		err := os.WriteFile(filepath.Join(dir, importFile), []byte(`
			type `+t+` struct {
				field string
			}`,
		), 0700)

		if err != nil {
			return "", err
		}
		return `
			import imp "` + importFile + `"
			s := imp.` + t + `{field: "Hello"}
			print(s.field)
		`, nil
	}, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "Hello", output)
	})
}

func testImportConstAssignmentFail(t *testing.T, transpilerFunc transpilerCalloutFunc) {
	value := "Hello World"
	constant := "PublicConst"

	transpilerFunc(t, func(dir string) (string, error) {
		importFile := "import.tsh"
		err := os.WriteFile(filepath.Join(dir, importFile), []byte(`const `+constant+` = "`+value+`"`), 0700)

		if err != nil {
			return "", err
		}
		return `
			import imp "` + importFile + `"
			imp.` + constant + ` = "Something else"
		`, nil
	}, func(output string, err error) {
		require.EqualError(t, shortenError(err), "cannot assign a value to constant imp."+constant)
	})
}

func testImportPrivateVariableFail(t *testing.T, transpilerFunc transpilerCalloutFunc) {
	variable := "privateVariable"
	value := "Hello World"

	transpilerFunc(t, func(dir string) (string, error) {
		importFile := "import.tsh"
		err := os.WriteFile(filepath.Join(dir, importFile), []byte(variable+` := "`+value+`"`), 0700)

		if err != nil {
			return "", err
		}
		return `
			import imp "` + importFile + `"
			print(imp.` + variable + `)
		`, nil
	}, func(output string, err error) {
		require.EqualError(t, shortenError(err), "variable imp."+variable+" has not been defined")
	})
}
