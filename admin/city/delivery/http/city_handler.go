package http

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/rudyraharjo/emurojaah/admin/city"
	"github.com/rudyraharjo/emurojaah/models"
)

// HTTPHandlerCity struct
type HTTPHandlerCity struct {
	cityService city.Service
}

// NewCityHTTPHandler func handler route
func NewCityHTTPHandler(r *gin.Engine, middleware *jwt.GinJWTMiddleware, cityService city.Service) {

	handler := HTTPHandlerCity{cityService}

	CityReq := r.Group("/api/admin/city")
	CityReq.Use(middleware.MiddlewareFunc())
	{
		CityReq.GET("/list", handler.HandlerListCity)
		CityReq.POST("/add", handler.HandlerAddCity)
		CityReq.POST("/delete", handler.HandlerDeleteCity)
		CityReq.POST("/update", handler.HandlerUpdateCity)
	}

}

// HandlerListCity Function
func (h *HTTPHandlerCity) HandlerListCity(ctx *gin.Context) {

	data := h.cityService.GetCityList()

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

// HandlerAddCity function
func (h *HTTPHandlerCity) HandlerAddCity(ctx *gin.Context) {
	var params models.City
	_ = ctx.Bind(&params)

	if params.ProvinceID == 0 || params.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params Cities",
		})
		return
	}

	InsertID, err := h.cityService.AddCity(params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success add City",
		"ID":      InsertID,
	})

}

// HandlerUpdateCity function
func (h *HTTPHandlerCity) HandlerUpdateCity(ctx *gin.Context) {
	var params models.City
	_ = ctx.Bind(&params)

	if params.ID == 0 || params.ProvinceID == 0 || params.Name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params Cities",
		})
		return
	}

	data, err := h.cityService.UpdateCity(params)
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
			"message": "success update Cities",
			"data":    data,
		})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "failed update Cities",
			"data":    params,
		})
		return
	}
}

// HandlerDeleteCity function
func (h *HTTPHandlerCity) HandlerDeleteCity(ctx *gin.Context) {
	var params models.RequestIDProvince
	_ = ctx.Bind(&params)

	data, errCheck := h.cityService.GetCityByID(params.ID)

	if errCheck != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error, please try again",
		})
		return
	}

	if len(data) > 0 {

		err := h.cityService.DeletedCity(params.ID)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "system error",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Success Deleted City",
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
