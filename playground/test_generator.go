package main

import (
	"fmt"
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

	userCode := ``

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

	CreateTestFile("foo_test.go", userCode, testTemplate, testCases)
}

func CreateTestFile(filename string, userCode string, testTemplate string, testCases []TestCase) {
	fmt.Printf("Creating and opening file %s", filename)
	fmt.Println("Adding \"package main\" and \"import \"testing\"\" to the file")
	fmt.Printf("Adding user code:\n%s\n", userCode)

	for testCaseIndex, testCase := range testCases {
		newTestCode := testTemplate
		newTestCode = strings.Replace(newTestCode, "{{ID}}", string(testCaseIndex), 1)
		newTestCode = strings.Replace(newTestCode, "{{OUTPUT}}", testCase.ExpectedOutput, 1)
		for inputIndex, input := range testCase.Inputs {
			newTestCode = strings.Replace(newTestCode, fmt.Sprintf("{{INPUT%d}}", inputIndex), input, 1)
		}
		fmt.Printf("Adding test:\n%s\n", newTestCode)
	}
}
