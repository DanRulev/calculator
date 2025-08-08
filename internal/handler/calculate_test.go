package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"calculator-go/internal/domain/model"
	"calculator-go/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_Operator_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := service.NewCalcService()
	h := New(svc)
	router := h.InitRouter()

	tests := []struct {
		name       string
		input      model.CalcInput
		wantResult string
	}{
		{
			name: "addition",
			input: model.CalcInput{
				Num1:     "10",
				Num2:     "20",
				Operator: "+",
			},
			wantResult: `{"result":"30"}`,
		},
		{
			name: "subtraction",
			input: model.CalcInput{
				Num1:     "50",
				Num2:     "15",
				Operator: "-",
			},
			wantResult: `{"result":"35"}`,
		},
		{
			name: "multiplication",
			input: model.CalcInput{
				Num1:     "6",
				Num2:     "7",
				Operator: "*",
			},
			wantResult: `{"result":"42"}`,
		},
		{
			name: "division",
			input: model.CalcInput{
				Num1:     "100",
				Num2:     "4",
				Operator: "/",
			},
			wantResult: `{"result":"25"}`,
		},
		{
			name: "division fractional",
			input: model.CalcInput{
				Num1:     "1",
				Num2:     "3",
				Operator: "/",
			},
			wantResult: `{"result":"0.333333"}`,
		},
		{
			name: "power 2^10",
			input: model.CalcInput{
				Num1:     "2",
				Num2:     "10",
				Operator: "^",
			},
			wantResult: `{"result":"1024"}`,
		},
		{
			name: "factorial 0!",
			input: model.CalcInput{
				Num1:     "0",
				Operator: "!",
			},
			wantResult: `{"result":"1"}`,
		},
		{
			name: "factorial 1!",
			input: model.CalcInput{
				Num1:     "1",
				Operator: "!",
			},
			wantResult: `{"result":"1"}`,
		},
		{
			name: "factorial 5!",
			input: model.CalcInput{
				Num1:     "5",
				Operator: "!",
			},
			wantResult: `{"result":"120"}`,
		},
		{
			name: "factorial 10!",
			input: model.CalcInput{
				Num1:     "10",
				Operator: "!",
			},
			wantResult: `{"result":"3628800"}`,
		},
		{
			name: "sqrt",
			input: model.CalcInput{
				Num1:     "25",
				Operator: "sqrt",
			},
			wantResult: `{"result":"5"}`,
		},
		{
			name: "sqrt 2 (irrational)",
			input: model.CalcInput{
				Num1:     "2",
				Operator: "sqrt",
			},
			wantResult: `{"result":"1.414214"}`,
		},
		{
			name: "negative addition",
			input: model.CalcInput{
				Num1:     "-10",
				Num2:     "25",
				Operator: "+",
			},
			wantResult: `{"result":"15"}`,
		},
		{
			name: "negative multiplication",
			input: model.CalcInput{
				Num1:     "-5",
				Num2:     "6",
				Operator: "*",
			},
			wantResult: `{"result":"-30"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := toJSON(t, tt.input)

			req, _ := http.NewRequest("POST", "/", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			assert.Equal(t, tt.wantResult, strings.TrimSpace(w.Body.String()))
		})
	}
}

func TestHandler_Operator_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)

	svc := service.NewCalcService()
	h := New(svc)
	router := h.InitRouter()

	tests := []struct {
		name       string
		input      model.CalcInput
		wantStatus int
		wantError  string
	}{
		{
			name: "divide by zero",
			input: model.CalcInput{
				Num1:     "10",
				Num2:     "0",
				Operator: "/",
			},
			wantStatus: http.StatusBadRequest,
			wantError:  "division by zero",
		},
		{
			name: "invalid number",
			input: model.CalcInput{
				Num1:     "abc",
				Operator: "+",
			},
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid number",
		},
		{
			name: "invalid json",
			input: model.CalcInput{
				Num1:     "",
				Operator: "",
			},
			wantStatus: http.StatusBadRequest,
			wantError:  "Invalid request",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := toJSON(t, tt.input)

			req, _ := http.NewRequest("POST", "/", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)

			require.NotEmpty(t, w.Body)
		})
	}
}

func toJSON(t *testing.T, v interface{}) []byte {
	data, err := json.Marshal(v)
	require.NoError(t, err)
	return data
}
