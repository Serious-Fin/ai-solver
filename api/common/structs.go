package common

type TestCase struct {
	Id             int      `json:"id"`
	Inputs         []string `json:"inputs"`
	ExpectedOutput string   `json:"output"`
}
