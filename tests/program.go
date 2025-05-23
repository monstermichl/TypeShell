package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func testComplexProgram1Success(t *testing.T, transpilerFunc transpilerFunc) {
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

func testComplexProgram2Success(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		// Formats a string message from two ints
		func numberMessage(original int, doubled int) string {
			return "Input: " + itoa(original) + ", Doubled: " + itoa(doubled)
		}

		// Doubles a number if it's non-negative
		func doubleIfPositive(n int) (int, error) {
			if n < 0 {
				return 0, "negative number"
			}
			return n * 2, ""
		}

		// Uses doubleIfPositive and formats a message
		func processNumber(n int) (string, error) {
			doubled, err := doubleIfPositive(n)
			if err != "" {
				return "", err
			}
			return numberMessage(n, doubled), ""
		}

		// Helper used in countPositive
		func isPositive(n int) bool {
			return n > 0
		}

		// Counts how many are positive using helper
		func countPositive(nums []int) int {
			count := 0
			for i := 0; i < len(nums); i++ {
				if isPositive(nums[i]) {
					count++
				}
			}
			return count
		}

		// Aggregates a report using multiple helper functions
		func statusReport(values []int) string {
			count := countPositive(values)

			if count == len(values) {
				return "All values are positive"
			}
			return "Some values are not positive"
		}

		// Test processNumber
		msg, err := processNumber(5)
		print("Process 5:", msg, "Error:", err)

		msg, err = processNumber(-1)
		print("Process -1:", msg, "Error:", err)

		// Test countPositive
		vals := []int{3, -2, 7, 0}
		print("Positive count:", countPositive(vals))

		// Test statusReport
		print("Report 1:", statusReport([]int{1, 2, 3}))
		print("Report 2:", statusReport([]int{1, -1, 3}))
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "Process 5: Input: 5, Doubled: 10 Error: \nProcess -1:  Error: negative number\nPositive count: 2\nReport 1: All values are positive\nReport 2: Some values are not positive", output)
	})
}
