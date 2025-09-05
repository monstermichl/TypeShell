package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testTypeDeclarationAndDefinitionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		type testType int

		var a testType
		a = testType(1)

		print(a)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1", output)
	})
}

func testTypeDeclarationAndAssignmentFail(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		type testType int

		var a testType
		a = 1

		print(a)
	`, func(output string, err error) {
		require.EqualError(t, shortenError(err), "expected testType but got int")
	})
}

func testTypeAliasAndAssignmentSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		type testType = int

		var a testType
		a = 1

		print(a)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1", output)
	})
}

func testTypeDeclarationAndDefinitionInFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func main() {
			type testType int

			var a testType
			a = testType(1)

			print(a)
		}
		main()
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1", output)
	})
}

func testTypeDeclarationAndAssignmentInFunctionFail(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func main() {
			type testType int

			var a testType
			a = 1

			print(a)
		}
		main()
	`, func(output string, err error) {
		require.EqualError(t, shortenError(err), "expected testType but got int")
	})
}

func testTypeAliasAndAssignmentInFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func main() {
			type testType = int

			var a testType
			a = 1

			print(a)
		}
		main()
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "1", output)
	})
}

func testTypeDeclaredInFunctionUsedOutsideFail(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		func main() {
			type testType int
		}
		var a testType
	`, func(output string, err error) {
		require.EqualError(t, shortenError(err), "expected valid data type")
	})
}

func testTypeDeclaredInIfUsedOutsideFail(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		if true {
			type testType int
		}
		var a testType
	`, func(output string, err error) {
		require.EqualError(t, shortenError(err), "expected valid data type")
	})
}

func testDeclareTypeTwiceFail(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		type myType int
		type myType string
	`, func(output string, err error) {
		require.EqualError(t, shortenError(err), "myType has already been defined")
	})
}

func testPassDeclaredTypeToFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		type myType string
		
		func test(param myType) {
			print(param)
		}
		test(myType("test"))
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "test", output)
	})
}

func testPassBaseTypeToFunctionFail(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		type myType string
		
		func test(param myType) {
			print(param)
		}
		test("test")
	`, func(output string, err error) {
		require.EqualError(t, shortenError(err), "expected parameter of type myType (param) but got string")
	})
}

func testPassValueWithSameBaseTypeToFunctionSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		type myType1 = string
		type myType2 = string
		type myType3 = myType2
		type myType4 = myType3
		
		func test(param myType2) {
			print(param)
		}
		test(myType4("test"))
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "test", output)
	})
}

func testAssignDifferentDefinedTypeFail(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		type myType1 string
		type myType2 string
		
		var a myType1
		a = myType2("test")
	`, func(output string, err error) {
		require.EqualError(t, shortenError(err), "expected myType1 but got myType2")
	})
}
