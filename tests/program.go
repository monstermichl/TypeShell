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

func testComplexProgram3Success(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		// === Manual string slice growth ===
		func growStr(slice []string, val string) []string {
			var out []string
			for i, unused := range slice {
				out[i] = slice[i]
			}
			out[len(out)] = val
			return out
		}

		// === Manual append for string slices ===
		func appendStr(slice []string, val string) []string {
			newSlice := []string{}
			for i, unused := range slice {
				newSlice = growStr(newSlice, slice[i])
			}
			newSlice = growStr(newSlice, val)
			return newSlice
		}

		// === Manual int slice growth ===
		func growInt(slice []int, val int) []int {
			var out []int
			for i, unused := range slice {
				out[i] = slice[i]
			}
			out[len(out)] = val
			return out
		}

		// === Manual append for int slices ===
		func appendInt(slice []int, val int) []int {
			newSlice := []int{}
			for i, unused := range slice {
				newSlice = growInt(newSlice, slice[i])
			}
			newSlice = growInt(newSlice, val)
			return newSlice
		}

		// === Simple int to string without rune ops ===
		func intToString(n int) string {
			if n == 0 {
				return "0"
			}
			digits := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
			powers := []int{1000000000,100000000,10000000,1000000,100000,10000,1000,100,10,1}
			s := ""
			started := false
			for i := 0; i < len(powers); i++ {
				d := n / powers[i]
				n = n % powers[i]
				if d != 0 || started || powers[i] == 1 {
					started = true
					s = s + digits[d]
				}
			}
			return s
		}

		// === Label construction ===
		func buildMessage(label string, total int, ok bool) string {
			prefix := "FAIL:"
			if ok {
				prefix = "OK:"
			}
			return prefix + " " + label + " = " + intToString(total)
		}

		// === Conditional summation ===
		func sumIfAllPositive(nums []int) (int, bool) {
			sum := 0
			for i := 0; i < len(nums); i++ {
				if nums[i] <= 0 {
					return 0, false
				}
				sum = sum + nums[i]
			}
			return sum, true
		}

		// === Combined summary function ===
		func summarize(name string, values []int) (string, bool) {
			total, ok := sumIfAllPositive(values)
			msg := buildMessage(name, total, ok)
			return msg, ok
		}

		// === Label validation ===
		func validateLabels(labels []string) (bool, string) {
			for i := 0; i < len(labels); i++ {
				if len([]string{labels[i]}) == 0 {
					return false, "empty label"
				}
			}
			return true, "all labels valid"
		}

		// === Label generation for positives ===
		func filterAndLabel(nums []int, tag string) ([]string, error) {
			var result []string
			for i := 0; i < len(nums); i++ {
				if nums[i] > 0 {
					label := tag + "_" + intToString(nums[i])
					result = appendStr(result, label)
				}
			}
			if len(result) == 0 {
				return nil, "no positive numbers"
			}
			return result, ""
		}

		// === High-level wrapper ===
		func reportCard(title string, scores []int, active bool) (string, error) {
			if !active {
				return "", "inactive user"
			}
			msg, ok := summarize(title, scores)
			if !ok {
				return msg, "some scores were non-positive"
			}
			return msg, ""
		}

		msg, ok := summarize("Math", []int{5, 10, 15})
		print("Summary 1:", msg, "Success:", ok)

		msg, ok = summarize("Science", []int{10, -1, 4})
		print("Summary 2:", msg, "Success:", ok)

		valid, valErr := validateLabels([]string{"OK", "Go"})
		print("Labels valid:", valid, "Message:", valErr)

		labels, err := filterAndLabel([]int{1, 2, -3}, "Score")
		if err != "" {
			print("Error:", err)
		} else {
			for i := 0; i < len(labels); i++ {
				print("Label", i, "=", labels[i])
			}
		}

		rep, err := reportCard("English", []int{3, 5, 7}, true)
		print("ReportCard 1:", rep, "Error:", err)

		rep, err = reportCard("History", []int{0, 2}, true)
		print("ReportCard 2:", rep, "Error:", err)

		rep, err = reportCard("Art", []int{9, 10}, false)
		print("ReportCard 3:", rep, "Error:", err)
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "Summary 1: OK: Math = 30 Success: 1\nSummary 2: FAIL: Science = 0 Success: 0\nLabels valid: 1 Message: all labels valid\nLabel 0 = Score_1\nLabel 1 = Score_2\nReportCard 1: OK: English = 15 Error: \nReportCard 2: FAIL: History = 0 Error: some scores were non-positive\nReportCard 3:  Error: inactive user", output)
	})
}

func testComplexProgram4Success(t *testing.T, transpilerFunc transpilerFunc) {
	transpilerFunc(t, `
		// Helper: Get element index for 2D logical position (row-major)
		func idx(row int, col int, cols int) int {
			return row*cols + col
		}

		// Step 1: Generate scores as 1D slice representing 2D matrix
		func generateScores(students int, subjects int) []int {
			total := students * subjects
			scores := []int{}
			for i := 0; i < total; i++ {
				scores[i] = 0 // fill with zeros initially
			}

			for r := 0; r < students; r++ {
				for c := 0; c < subjects; c++ {
					// Example scoring logic
					score := (r + 1) * (c + 2)
					scores[idx(r, c, subjects)] = score
				}
			}
			return scores
		}

		// Step 2: Label a single student’s scores from 1D slice + indexing
		func labelStudentScores(scores []int, student int, subjects int) []string {
			labels := []string{}
			for i := 0; i < subjects; i++ {
				labels[len(labels)] = "Student" + itoa(student+1) + "_Sub" + itoa(i+1) + "=" + itoa(scores[idx(student, i, subjects)])
			}
			return labels
		}

		// Step 3: Summarize student's labeled scores
		func summarizeStudent(labels []string) string {
			count := 0
			last := ""
			for i, l := range labels {
				count = i + 1
				last = l
			}
			return "Count=" + itoa(count) + ", Last=" + last
		}

		// Step 4: Evaluate classroom by iterating students
		func evaluateClassroom(scores []int, students int, subjects int) []string {
			summaries := []string{}
			for s := 0; s < students; s++ {
				labels := labelStudentScores(scores, s, subjects)
				summary := summarizeStudent(labels)
				summaries[len(summaries)] = summary
			}
			return summaries
		}

		students := 3
		subjects := 4
		scores := generateScores(students, subjects)
		summaries := evaluateClassroom(scores, students, subjects)

		for i := 0; i < len(summaries); i++ {
			print("Student", itoa(i+1)+":", summaries[i])
		}
	`, func(output string, err error) {
		require.Nil(t, err)
		require.Equal(t, "Student 1: Count=4, Last=Student1_Sub4=5\nStudent 2: Count=4, Last=Student2_Sub4=10\nStudent 3: Count=4, Last=Student3_Sub4=15", output)
	})
}
