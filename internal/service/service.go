package service

import (
	"calculator-go/internal/domain/model"
	"fmt"
	"math"
	"math/big"
)

type CalcService struct{}

func NewCalcService() *CalcService {
	return &CalcService{}
}

// Внутренние типы
var (
	zero = big.NewFloat(0)
	one  = big.NewFloat(1)
	two  = big.NewFloat(2)
	ten  = big.NewFloat(10)
	prec = uint(512) // точность для big.Float (в битах)
)

func (c *CalcService) Calculator(calc model.CalcInput) (*big.Float, error) {
	// Парсим числа
	num1, ok := big.NewFloat(0).SetString(calc.Num1)
	if !ok {
		return nil, fmt.Errorf("invalid number: %s", calc.Num1)
	}
	num1.SetPrec(prec)

	var num2 *big.Float
	if calc.Num2 != "" {
		n, ok := big.NewFloat(0).SetString(calc.Num2)
		if !ok {
			return nil, fmt.Errorf("invalid number: %s", calc.Num2)
		}
		n.SetPrec(prec)
		num2 = n
	}

	// Установим режим округления
	num1.SetMode(big.ToZero)
	if num2 != nil {
		num2.SetMode(big.ToZero)
	}

	var result *big.Float

	switch calc.Operator {
	case "+":
		result = new(big.Float).Add(num1, num2)
	case "-":
		result = new(big.Float).Sub(num1, num2)
	case "*":
		result = new(big.Float).Mul(num1, num2)
	case "/":
		if num2.Sign() == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		result = new(big.Float).Quo(num1, num2)
	case "^":
		res, err := powerFloat(num1, num2)
		if err != nil {
			return nil, err
		}
		result = res
	case "!":
		res, err := factorialBig(num1)
		if err != nil {
			return nil, err
		}
		result = res
	case "sqrt":
		if num1.Sign() < 0 {
			return nil, fmt.Errorf("cannot compute sqrt of negative number")
		}
		result = new(big.Float).Sqrt(num1)
	default:
		return nil, fmt.Errorf("unsupported operator: %s", calc.Operator)
	}

	// Ограничиваем точность результата
	result.SetPrec(prec)
	result.SetMode(big.ToNearestEven)

	return result, nil
}

func powerFloat(base, exp *big.Float) (*big.Float, error) {
	// Если степень целая — можно использовать умножение
	// Иначе — используем exp(ln(base) * exp), но с big.Float это сложно
	// Упрощение: если степень — целое число, используем умножение

	// Попробуем преобразовать exp в int64
	f64, _ := exp.Float64()
	if math.IsInf(f64, 0) {
		return nil, fmt.Errorf("exponent too large")
	}

	if exp.IsInt() {
		i, _ := exp.Int64()
		if i >= 0 && i <= 10000 { // ограничение для безопасности
			result := big.NewFloat(1).SetPrec(prec).SetMode(big.ToZero)
			baseCopy := new(big.Float).Set(base).SetPrec(prec)
			for j := int64(0); j < i; j++ {
				result.Mul(result, baseCopy)
			}
			return result, nil
		} else if i < 0 {
			pos := -i
			result := big.NewFloat(1).SetPrec(prec)
			baseCopy := new(big.Float).Set(base).SetPrec(prec)
			for j := int64(0); j < pos; j++ {
				result.Mul(result, baseCopy)
			}
			return new(big.Float).Quo(one, result), nil
		}
	}

	// Для нецелых степеней — можно использовать log/exp, но это сложно и медленно
	// Пока вернём ошибку
	return nil, fmt.Errorf("non-integer exponent not supported yet: %s", exp.String())
}

func factorialBig(n *big.Float) (*big.Float, error) {
	// Проверяем, что число неотрицательное
	if n.Sign() < 0 {
		return nil, fmt.Errorf("factorial is not defined for negative numbers")
	}

	// Проверяем, что это целое число
	if !n.IsInt() {
		return nil, fmt.Errorf("factorial is only defined for integers")
	}

	// Преобразуем в int64
	i, acc := n.Int64()
	if acc != big.Exact {
		return nil, fmt.Errorf("value too large or not integer")
	}

	if i > 500 { // или другое разумное число
		return nil, fmt.Errorf("factorial input too large: %d", i)
	}

	// Вычисляем факториал
	result := big.NewFloat(1).SetPrec(prec)
	temp := new(big.Float).SetPrec(prec)

	for j := int64(1); j <= i; j++ {
		temp.SetInt64(j)
		result.Mul(result, temp)
	}

	return result, nil
}
