package service

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalcService_Evaluate(t *testing.T) {
	t.Parallel()

	type args struct {
		expr string
	}
	tests := []struct {
		name       string
		args       args
		want       *big.Float
		wantErr    bool
		wantErrMsg string
	}{
		{
			name:    "simple addition",
			args:    args{expr: "2 + 3"},
			want:    big.NewFloat(5),
			wantErr: false,
		},
		{
			name:    "simple subtraction",
			args:    args{expr: "10 - 4"},
			want:    big.NewFloat(6),
			wantErr: false,
		},
		{
			name:    "multiplication and division",
			args:    args{expr: "6 * 7 / 3"},
			want:    big.NewFloat(14),
			wantErr: false,
		},
		{
			name:    "power operation",
			args:    args{expr: "2 ^ 3"},
			want:    big.NewFloat(8),
			wantErr: false,
		},
		{
			name:    "factorial",
			args:    args{expr: "5!"},
			want:    big.NewFloat(120),
			wantErr: false,
		},
		{
			name:    "parentheses",
			args:    args{expr: "(2 + 3) * 4"},
			want:    big.NewFloat(20),
			wantErr: false,
		},
		{
			name:    "negative number",
			args:    args{expr: "-5"},
			want:    big.NewFloat(-5),
			wantErr: false,
		},
		{
			name:    "floating point",
			args:    args{expr: "1.5 + 2.5"},
			want:    big.NewFloat(4.0),
			wantErr: false,
		},
		{
			name:    "complex expression",
			args:    args{expr: "((2 + 3)! - 1) / 2 ^ 2"},
			want:    big.NewFloat(29.75),
			wantErr: false,
		},

		// --- Ошибочные случаи ---
		{
			name:       "division by zero",
			args:       args{expr: "10 / 0"},
			wantErr:    true,
			wantErrMsg: "integer division by zero",
		},
		{
			name:       "invalid factorial non-integer",
			args:       args{expr: "3.5!"},
			wantErr:    true,
			wantErrMsg: "factorial of non-integer number is not defined",
		},
		{
			name:       "factorial negative",
			args:       args{expr: "(-3)!"},
			wantErr:    true,
			wantErrMsg: "factorial of negative number is not defined",
		},
		{
			name:       "unmatched left parenthesis",
			args:       args{expr: "(2 + 3"},
			wantErr:    true,
			wantErrMsg: "expecting closing parenthesis",
		},
		{
			name:       "unmatched right parenthesis",
			args:       args{expr: "2 + 3)"},
			wantErr:    true,
			wantErrMsg: "unexpected token",
		},
		{
			name:       "empty expression",
			args:       args{expr: ""},
			wantErr:    true,
			wantErrMsg: "unexpected token",
		},
		{
			name:       "only spaces",
			args:       args{expr: "   "},
			wantErr:    true,
			wantErrMsg: "unexpected token",
		},
		{
			name:       "invalid character",
			args:       args{expr: "2 + @"},
			wantErr:    true,
			wantErrMsg: "invalid character",
		},
		{
			name:       "trailing tokens",
			args:       args{expr: "2 + 3 4"},
			wantErr:    true,
			wantErrMsg: "expecting end of expression",
		},
		{
			name:       "non integer exponent",
			args:       args{expr: "2 ^ 1.5"},
			wantErr:    true,
			wantErrMsg: "unexpected expression",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service := NewCalcService()
			result, err := service.Evaluate(tt.args.expr)

			if tt.wantErr {
				require.Error(t, err)
				if tt.wantErrMsg != "" {
					require.Contains(t, err.Error(), tt.wantErrMsg)
				}
				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)

			// Сравниваем с точностью
			cmp := result.Cmp(tt.want)
			assert.Zero(t, cmp, "expected %s, got %s", tt.want.Text('g', 10), result.Text('g', 10))
		})
	}
}
