package http

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/rudyraharjo/emurojaah/admin/province"
	"github.com/rudyraharjo/emurojaah/models"
)

// HTTPHandlerProvince struct
type HTTPHandlerProvince struct {
	provinceService province.Service
}

// NewQuoteHTTPHandler func handler route
func NewProvinceHTTPHandler(r *gin.Engine, middleware *jwt.GinJWTMiddleware, provinceService province.Service) {

	handler := HTTPHandlerProvince{provinceService}

	ProvinceReq := r.Group("/api/admin/province")
	ProvinceReq.Use(middleware.MiddlewareFunc())
	{
		ProvinceReq.GET("/list", handler.HandlerListProvince)
		ProvinceReq.POST("/add", handler.HandlerAddProvince)
		ProvinceReq.POST("/delete", handler.HandlerDeleteProvince)
		ProvinceReq.POST("/update", handler.HandlerUpdateProvince)
	}

}

// HandlerListProvince Function
func (h *HTTPHandlerProvince) HandlerListProvince(ctx *gin.Context) {

	data := h.provinceService.GetProvinceList()

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

// HandlerAddProvince function
func (h *HTTPHandlerProvince) HandlerAddProvince(ctx *gin.Context) {
	var params models.Province
	_ = ctx.Bind(&params)

	if params.Name == "" || params.InternationID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params Province",
		})
		return
	}

	InsertID, err := h.provinceService.AddProvince(params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success add Province",
		"ID":      InsertID,
	})

}

// HandlerUpdateProvince function
func (h *HTTPHandlerProvince) HandlerUpdateProvince(ctx *gin.Context) {
	var params models.Province
	_ = ctx.Bind(&params)

	if params.ID == 0 || params.Name == "" || params.InternationID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params Province",
		})
		return
	}

	data, err := h.provinceService.UpdateProvince(params)
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
			"message": "success update Province",
			"data":    data,
		})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "failed update Province",
			"data":    params,
		})
		return
	}
}

// HandlerDeleteProvince function
func (h *HTTPHandlerProvince) HandlerDeleteProvince(ctx *gin.Context) {
	var params models.RequestIDProvince
	_ = ctx.Bind(&params)

	data, errCheck := h.provinceService.GetProvinceByID(params.ID)

	if errCheck != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error, please try again",
		})
		return
	}

	if len(data) > 0 {

		err := h.provinceService.DeletedProvince(params.ID)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusInternalServerError,
				"message": "system error",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Success Deleted Province",
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
