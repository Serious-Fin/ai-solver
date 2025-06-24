package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type TestCase struct {
	Inputs         []string
	ExpectedOutput string
}

func main() {
	testTemplate := `func TestTwoSum{{ID}}(t *testing.T) {
	got := twoSum({{INPUT0}}, {{INPUT1}})
	want := {{OUTPUT}}
	if !areEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}`

	userCode := `func twoSum(nums []int, target int) []int {
	numMap := make(map[int]int)

	for i, num := range nums {
		complement := target - num

		if prevIndex, ok := numMap[complement]; ok {
			return []int{prevIndex, i}
		}

		numMap[num] = i
	}

	return nil
}`

	testCases := []TestCase{
		{
			Inputs: []string{
				"[]int{2, 7, 11, 15}",
				"9",
			},
			ExpectedOutput: "[]int{0, 1}",
		},
		{
			Inputs: []string{
				"[]int{1, 2, 3, 4, 4}",
				"8",
			},
			ExpectedOutput: "[]int{3, 4}",
		},
	}

	helperFuncs := `func areEqual(got []int, want []int) bool {
	if len(got) != 2 {
		return false
	}

	if got[0] == want[0] && got[1] == want[1] {
		return true
	}
	if got[0] == want[1] && got[1] == want[0] {
		return true
	}

	return false
}
`

	CreateTestFile("foo_test.go", userCode, testTemplate, testCases, helperFuncs)
}

var fileStartTemplate = `package main
import "testing"
`

func CreateTestFile(filename string, userCode string, testTemplate string, testCases []TestCase, helperFuncs string) {
	file, err := os.Create(filename)
	check(err)
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s\n", fileStartTemplate))
	check(err)
	_, err = file.WriteString(fmt.Sprintf("%s\n", userCode))
	check(err)

	for testCaseIndex, testCase := range testCases {
		newTestCode := testTemplate
		newTestCode = strings.Replace(newTestCode, "{{ID}}", strconv.Itoa(testCaseIndex), 1)
		newTestCode = strings.Replace(newTestCode, "{{OUTPUT}}", testCase.ExpectedOutput, 1)
		for inputIndex, input := range testCase.Inputs {
			newTestCode = strings.Replace(newTestCode, fmt.Sprintf("{{INPUT%d}}", inputIndex), input, 1)
		}
		_, err = file.WriteString(fmt.Sprintf("%s\n", newTestCode))
		check(err)
	}
	_, err = file.WriteString(helperFuncs)
	check(err)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
