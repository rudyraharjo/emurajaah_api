package http

import (
	"net/http"

	"github.com/rudyraharjo/emurojaah/admin/banner"
	"github.com/rudyraharjo/emurojaah/models"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

// HTTPHandlerBanner struct
type HTTPHandlerBanner struct {
	bannerService banner.Service
}

// NewPointHTTPHandler func
func NewBannerHTTPHandler(r *gin.Engine, middleware *jwt.GinJWTMiddleware, bannerService banner.Service) {

	handler := HTTPHandlerBanner{bannerService}

	bannerReq := r.Group("/api/admin/banner")
	bannerReq.Use(middleware.MiddlewareFunc())
	{
		bannerReq.GET("/list", handler.HandlerListBanner)
		bannerReq.POST("/add", handler.HandlerAddBanner)
		bannerReq.POST("/update", handler.HandlerUpdateBanner)
		bannerReq.POST("/delete", handler.HandlerDeleteBanner)
		bannerReq.POST("/inactive", handler.HandlerInActiveBanner)
		bannerReq.POST("/isactive", handler.HandlerActiveBanner)
	}

}

// HandlerListBanner Func
func (h *HTTPHandlerBanner) HandlerListBanner(ctx *gin.Context) {

	data := h.bannerService.GetBannerList()

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "ok",
		"data":    data,
	})

}

// HandlerAddBanner func Add banner
func (h *HTTPHandlerBanner) HandlerAddBanner(ctx *gin.Context) {
	var params models.Banner
	_ = ctx.Bind(&params)

	data, err := h.bannerService.AddBanner(params)

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

// HandlerUpdateBanner func update banner
func (h *HTTPHandlerBanner) HandlerUpdateBanner(ctx *gin.Context) {
	var params models.Banner
	_ = ctx.Bind(&params)

	if params.Id == 0 || params.ImageUrl == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params banner",
		})
		return
	}

	data, err := h.bannerService.UpdateBanner(params)

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
			"message": "success update banner",
			"data":    data,
		})
		return

	} else {

		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "missing params ID banner",
			"data":    data,
		})
		return
	}
}

// HandlerInActiveBanner func inactive banner
func (h *HTTPHandlerBanner) HandlerInActiveBanner(ctx *gin.Context) {
	var params models.RequestIDBanner
	_ = ctx.Bind(&params)

	if params.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params banner",
		})
		return
	}

	err := h.bannerService.UpdateToinactive(params.ID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success inactive banner",
		"ID":      params.ID,
	})

}

// HandlerActiveBanner func active banner
func (h *HTTPHandlerBanner) HandlerActiveBanner(ctx *gin.Context) {
	var params models.RequestIDBanner
	_ = ctx.Bind(&params)

	if params.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params banner",
		})
		return
	}

	err := h.bannerService.UpdateToActive(params.ID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success active banner",
		"ID":      params.ID,
	})
}

// HandlerDeleteBanner function
func (h *HTTPHandlerBanner) HandlerDeleteBanner(ctx *gin.Context) {
	var params models.RequestIDBanner
	_ = ctx.Bind(&params)

	if params.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params banner",
		})
		return
	}

	err := h.bannerService.DeleteBanner(params.ID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success delete banner",
		"ID":      params.ID,
	})
}
