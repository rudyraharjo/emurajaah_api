package http

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/rudyraharjo/emurojaah/admin/point"
	"github.com/rudyraharjo/emurojaah/models"
)

// HTTPHandlerPoint struct
type HTTPHandlerPoint struct {
	PointService point.Service
}

// NewPointHTTPHandler func
func NewPointHTTPHandler(r *gin.Engine, middleware *jwt.GinJWTMiddleware, pointService point.Service) {

	handler := HTTPHandlerPoint{pointService}

	pointReq := r.Group("/api/admin/point")
	pointReq.Use(middleware.MiddlewareFunc())
	{
		pointReq.GET("/list", handler.HandlerListPoint)
		pointReq.POST("/add", handler.HandlerAddPoint)
		pointReq.POST("/update", handler.HandlerUpdatePoint)
		pointReq.POST("/inactive", handler.HandlerInActivePoint)
		pointReq.POST("/isactive", handler.HandlerActivePoint)
	}

}

// HandlerListPoint Func
func (h *HTTPHandlerPoint) HandlerListPoint(ctx *gin.Context) {

	data := h.PointService.GetPointList()

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "ok",
		"data":    data,
	})

}

// HandlerAddPoint func Add Point
func (h *HTTPHandlerPoint) HandlerAddPoint(ctx *gin.Context) {
	var params models.RequestAddPoint
	_ = ctx.Bind(&params)

	data, err := h.PointService.AddPoint(params)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error",
			"data":    data,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success add banner",
		"data":    data,
	})
}

// HandlerUpdatePoint func update point
func (h *HTTPHandlerPoint) HandlerUpdatePoint(ctx *gin.Context) {
	var params models.ResponsePoint
	_ = ctx.Bind(&params)

	if params.ID == 0 || params.Type == "" || params.Point == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params point",
		})
		return
	}

	data, err := h.PointService.UpdatePoint(params)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error",
			"data":    data,
		})
		return

	} else if len(data) > 0 {

		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "success update point",
			"data":    data,
		})
		return

	} else {

		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "missing params ID point",
			"data":    data,
		})
		return
	}
}

// HandlerInActivePoint func inactive point
func (h *HTTPHandlerPoint) HandlerInActivePoint(ctx *gin.Context) {
	var params models.RequestIDPoint
	_ = ctx.Bind(&params)

	err := h.PointService.UpdateToinactive(params.ID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success inactive point",
		"ID":      params.ID,
	})

}

// HandlerActivePoint func active point
func (h *HTTPHandlerPoint) HandlerActivePoint(ctx *gin.Context) {
	var params models.RequestIDPoint
	_ = ctx.Bind(&params)

	err := h.PointService.UpdateToActive(params.ID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success active point",
		"ID":      params.ID,
	})
}
