package http

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/rudyraharjo/emurojaah/admin/country"
	"github.com/rudyraharjo/emurojaah/models"
)

// HTTPHandlerCountry struct
type HTTPHandlerCountry struct {
	countryService country.Service
}

// NewQuoteHTTPHandler func handler route
func NewCountryHTTPHandler(r *gin.Engine, middleware *jwt.GinJWTMiddleware, countryService country.Service) {

	handler := HTTPHandlerCountry{countryService}

	CountryReq := r.Group("/api/admin/country")
	CountryReq.Use(middleware.MiddlewareFunc())
	{
		CountryReq.GET("/list", handler.HandlerListCountry)
		CountryReq.POST("/add", handler.HandlerAddCountry)
		CountryReq.POST("/delete", handler.HandlerDeleteCountry)
		CountryReq.POST("/update", handler.HandlerUpdateCountry)
	}

}

// HTTPHandlerCountry Function
func (h *HTTPHandlerCountry) HandlerListCountry(ctx *gin.Context) {

	data := h.countryService.GetCountryList()

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
func (h *HTTPHandlerCountry) HandlerAddCountry(ctx *gin.Context) {
	var params models.Country
	_ = ctx.Bind(&params)

	if params.CountryName == "" || params.CountryCode == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params Country",
		})
		return
	}

	InsertID, err := h.countryService.AddCountry(params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success add Country",
		"ID":      InsertID,
	})

}

// HandlerUpdateCountry function
func (h *HTTPHandlerCountry) HandlerUpdateCountry(ctx *gin.Context) {
	var params models.Country
	_ = ctx.Bind(&params)

	if params.ID == 0 || params.CountryCode == "" || params.CountryName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params Country",
		})
		return
	}

	data, err := h.countryService.UpdateCountry(params)
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
			"message": "success update Country",
			"data":    data,
		})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "failed update Country",
			"data":    params,
		})
		return
	}
}

// HandlerDeleteCountry function
func (h *HTTPHandlerCountry) HandlerDeleteCountry(ctx *gin.Context) {
	var params models.RequestIDCountry
	_ = ctx.Bind(&params)

	data, errCheck := h.countryService.GetCountryByID(params.ID)

	if errCheck != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error, please try again",
		})
		return
	}

	if len(data) > 0 {

		err := h.countryService.DeletedCountry(params.ID)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "system error",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Success Deleted Country",
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
