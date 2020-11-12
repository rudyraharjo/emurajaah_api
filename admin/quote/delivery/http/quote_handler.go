package http

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	quote "github.com/rudyraharjo/emurojaah/admin/quote"
	"github.com/rudyraharjo/emurojaah/models"
)

// HTTPHandlerQuote struct
type HTTPHandlerQuote struct {
	quoteService quote.Service
}

// NewQuoteHTTPHandler func handler route
func NewQuoteHTTPHandler(r *gin.Engine, middleware *jwt.GinJWTMiddleware, quoteService quote.Service) {

	handler := HTTPHandlerQuote{quoteService}

	QuoteReq := r.Group("/api/admin/quote")
	QuoteReq.Use(middleware.MiddlewareFunc())
	{
		QuoteReq.GET("/list", handler.HandlerListQuote)
		QuoteReq.POST("/add", handler.HandlerAddQuote)
		QuoteReq.POST("/delete", handler.HandlerDeleteQuote)
		QuoteReq.POST("/update", handler.HandlerUpdateQuote)
		QuoteReq.POST("/update-status", handler.HandlerUpdateStatusQuote)
	}

}

// HandlerListQuote Function
func (h *HTTPHandlerQuote) HandlerListQuote(ctx *gin.Context) {

	data := h.quoteService.GetQuotesList()

	if len(data) > 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "ok",
			"data":    data,
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "no data",
			"data":    data,
		})
	}

}

// HandlerAddQuote function
func (h *HTTPHandlerQuote) HandlerAddQuote(ctx *gin.Context) {
	var params models.Quote
	_ = ctx.Bind(&params)

	if params.Message == "" || params.Author == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params Quote",
		})
		return
	}

	InsertID, err := h.quoteService.AddQuote(params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success add Quote",
		"ID":      InsertID,
	})

}

// HandlerUpdateQuote function
func (h *HTTPHandlerQuote) HandlerUpdateQuote(ctx *gin.Context) {
	var params models.Quote
	_ = ctx.Bind(&params)

	if params.Id == 0 || params.Message == "" || params.Author == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params Quote",
		})
		return
	}

	data, err := h.quoteService.UpdateQuote(params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error",
			"data":    data,
		})
		return
	}

	if len(data) > 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "success update Quote",
			"data":    data,
		})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "failed update Quote",
			"data":    params,
		})
		return
	}
}

// HandlerDeleteQuote function
func (h *HTTPHandlerQuote) HandlerDeleteQuote(ctx *gin.Context) {
	var params models.RequestIDQuote
	_ = ctx.Bind(&params)

	data, errCheck := h.quoteService.GetQuoteByID(params.ID)

	if errCheck != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error, please try again",
		})
		return
	}

	if len(data) > 0 {

		err := h.quoteService.DeletedQuote(params.ID)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "system error",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Success Deleted Quote",
			"ID":      params.ID,
		})
		return

	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Missing parrams",
		"ID":      params.ID,
	})

}

// HandlerUpdateStatusQuote func
func (h *HTTPHandlerQuote) HandlerUpdateStatusQuote(ctx *gin.Context) {
	var params models.RequestIDQuote
	_ = ctx.Bind(&params)
	msg := ""

	data, errCheck := h.quoteService.GetQuoteByID(params.ID)

	if errCheck != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error, please try again",
		})
		return
	}

	if len(data) > 0 {

		err := h.quoteService.UpdateStatus(params.ID)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "system error",
			})
			return
		}

		if data[0].IsActive == 0 {
			msg = "success active Quote"
		} else {
			msg = "success inactive Quote"
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": msg,
			"ID":      params.ID,
		})
		return

	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Missing parrams",
		"ID":      params.ID,
	})

}
