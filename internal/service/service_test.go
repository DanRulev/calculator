package service

import (
	"testing"

	"calculator-go/internal/domain/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalcService_Calculator_Addition(t *testing.T) {
	s := NewCalcService()

	input := model.CalcInput{
		Num1:     "10",
		Num2:     "20",
		Operator: "+",
	}

	result, err := s.Calculator(input)
	require.NoError(t, err)

	assert.Equal(t, "30", result.Text('f', -1))
}

func TestCalcService_Calculator_Subtraction(t *testing.T) {
	s := NewCalcService()

	input := model.CalcInput{
		Num1:     "50",
		Num2:     "15",
		Operator: "-",
	}

	result, err := s.Calculator(input)
	require.NoError(t, err)

	assert.Equal(t, "35", result.Text('f', -1))
}

func TestCalcService_Calculator_Multiplication(t *testing.T) {
	s := NewCalcService()

	input := model.CalcInput{
		Num1:     "12",
		Num2:     "15",
		Operator: "*",
	}

	result, err := s.Calculator(input)
	require.NoError(t, err)

	assert.Equal(t, "180", result.Text('f', -1))
}

func TestCalcService_Calculator_Division(t *testing.T) {
	s := NewCalcService()

	input := model.CalcInput{
		Num1:     "100",
		Num2:     "3",
		Operator: "/",
	}

	result, err := s.Calculator(input)
	require.NoError(t, err)

	// 100 / 3 ≈ 33.333333
	assert.Equal(t, "33.333333", result.Text('f', 6))
}

func TestCalcService_Calculator_DivideByZero(t *testing.T) {
	s := NewCalcService()

	input := model.CalcInput{
		Num1:     "10",
		Num2:     "0",
		Operator: "/",
	}

	result, err := s.Calculator(input)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "division by zero")
}

func TestCalcService_Calculator_Power(t *testing.T) {
	s := NewCalcService()

	input := model.CalcInput{
		Num1:     "2",
		Num2:     "10",
		Operator: "^",
	}

	result, err := s.Calculator(input)
	require.NoError(t, err)

	assert.Equal(t, "1024", result.Text('f', -1))
}

func TestCalcService_Calculator_Factorial(t *testing.T) {
	s := NewCalcService()

	input := model.CalcInput{
		Num1:     "5",
		Operator: "!",
	}

	result, err := s.Calculator(input)
	require.NoError(t, err)

	assert.Equal(t, "120", result.Text('f', -1))
}

func TestCalcService_Calculator_Factorial_Negative(t *testing.T) {
	s := NewCalcService()

	input := model.CalcInput{
		Num1:     "-5",
		Operator: "!",
	}

	result, err := s.Calculator(input)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "negative numbers")
}

func TestCalcService_Calculator_Factorial_NonInteger(t *testing.T) {
	s := NewCalcService()

	input := model.CalcInput{
		Num1:     "5.5",
		Operator: "!",
	}

	result, err := s.Calculator(input)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "integers")
}

func TestCalcService_Calculator_Sqrt(t *testing.T) {
	s := NewCalcService()

	input := model.CalcInput{
		Num1:     "16",
		Operator: "sqrt",
	}

	result, err := s.Calculator(input)
	require.NoError(t, err)

	assert.Equal(t, "4", result.Text('f', -1))
}

func TestCalcService_Calculator_Sqrt_Negative(t *testing.T) {
	s := NewCalcService()

	input := model.CalcInput{
		Num1:     "-4",
		Operator: "sqrt",
	}

	result, err := s.Calculator(input)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "negative number")
}

func TestCalcService_Calculator_InvalidOperator(t *testing.T) {
	s := NewCalcService()

	input := model.CalcInput{
		Num1:     "10",
		Num2:     "5",
		Operator: "%",
	}

	result, err := s.Calculator(input)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unsupported operator")
}

func TestCalcService_Calculator_InvalidNumber(t *testing.T) {
	s := NewCalcService()

	input := model.CalcInput{
		Num1:     "not-a-number",
		Num2:     "5",
		Operator: "+",
	}

	result, err := s.Calculator(input)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid number")
}
