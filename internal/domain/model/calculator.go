package model

type CalcInput struct {
	Num1     string `json:"num1"`
	Num2     string `json:"num2,omitempty"`
	Operator string `json:"operator"`
}
