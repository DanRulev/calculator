package handler

import "github.com/gin-gonic/gin"

type Handler struct {
	calculator calculatorServiceI
}

func New(calculator calculatorServiceI) *Handler {
	return &Handler{
		calculator: calculator,
	}
}

func (h *Handler) InitRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	calc := router.Group("/")
	{
		calc.GET("/", func(ctx *gin.Context) {
			router.LoadHTMLGlob("templates/*")
			h.GetOperator(ctx)
		})
		calc.POST("/", h.Operator)
	}

	return router
}
