package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testComplexProgramSuccess(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		// Returns int doubled and a fake error if input is negative
		func doubleIfPositive(n int) (int, error) {
			if n < 0 {
				return 0, "negative number"
			}
			return n * 2, ""
		}

		// Takes bool and string, returns a new string and a flag
		func tagStatus(status bool, label string) (string, bool) {
			if status {
				return "OK: " + label, true
			}
			return "FAIL: " + label, false
		}

		// Returns true if all integers are non-zero
		func allNonZero(nums []int) bool {
			for i := 0; i < len(nums); i++ {
				if nums[i] == 0 {
					return false
				}
			}
			return true
		}

		// Takes a slice of bools, returns count of "true"
		func countTrue(bools []bool) int {
			count := 0
			for i := 0; i < len(bools); i++ {
				if bools[i] {
					count++
				}
			}
			return count
		}

		// Returns a constant slice of ints
		func getConstants() []int {
			return []int{1, 2, 3}
		}

		// Test doubleIfPositive
		a, err := doubleIfPositive(10)
		print("Double 10:", a, "Error:", err)

		b, err := doubleIfPositive(-5)
		print("Double -5:", b, "Error:", err)

		// Test tagStatus
		s, ok := tagStatus(true, "Ready")
		print("Status:", s, "OK:", ok)

		s, ok = tagStatus(false, "Error")
		print("Status:", s, "OK:", ok)

		// Test allNonZero
		vals := []int{1, 2, 3}
		print("All non-zero:", allNonZero(vals))

		vals = []int{1, 0, 2}
		print("All non-zero:", allNonZero(vals))

		// Test countTrue
		flags := []bool{true, false, true, true}
		print("Count true:", countTrue(flags))

		// Test getConstants
		constants := getConstants()
		print("Constants:", constants[0], constants[1], constants[2])
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "Double 10: 20 Error: \nDouble -5: 0 Error: negative number\nStatus: OK: Ready OK: 1\nStatus: FAIL: Error OK: 0\nAll non-zero: 1\nAll non-zero: 0\nCount true: 3\nConstants: 1 2 3", output)
	})
}
