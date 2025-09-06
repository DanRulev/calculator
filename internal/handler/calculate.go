package handler

import (
	"math/big"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type calculatorServiceI interface {
	Evaluate(expr string) (*big.Float, error)
}

func (h *Handler) GetOperator(c *gin.Context) {
	c.HTML(http.StatusOK, "calculator.html", nil)
}

func (h *Handler) Operator(c *gin.Context) {
	var expr string

	err := c.BindJSON(&expr)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	result, err := h.calculator.Evaluate(expr)
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
	s := f.Text('f', 6)

	if strings.Contains(s, ".") {
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
	}

	return s
}
