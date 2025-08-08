package handler

import (
	"calculator-go/internal/domain/model"
	"math/big"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type calculatorServiceI interface {
	Calculator(calc model.CalcInput) (*big.Float, error)
}

func (h *Handler) GetOperator(c *gin.Context) {
	c.HTML(http.StatusOK, "calculator.html", nil)
}

func (h *Handler) Operator(c *gin.Context) {
	var input model.CalcInput

	err := c.BindJSON(&input)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	result, err := h.calculator.Calculator(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	outputStr := formatBigFloat(result)

	c.JSON(http.StatusOK, gin.H{
		"result": outputStr,
	})
}

func formatBigFloat(f *big.Float) string {
	// Форматируем с 6 знаками после запятой
	s := f.Text('f', 6)

	// Если есть точка
	if strings.Contains(s, ".") {
		// Убираем хвостовые нули: 12.340000 → 12.34
		s = strings.TrimRight(s, "0")
		// Убираем точку, если стала последней: 12. → 12
		s = strings.TrimRight(s, ".")
	}

	return s
}
