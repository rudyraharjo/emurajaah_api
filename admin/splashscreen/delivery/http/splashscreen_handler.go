package http

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/rudyraharjo/emurojaah/admin/splashscreen"
	"github.com/rudyraharjo/emurojaah/models"
)

// HTTPHandlerSplashScreen struct
type HTTPHandlerSplashScreen struct {
	splashscreenService splashscreen.Service
}

// NewSplashScreenHTTPHandler func handler route
func NewSplashScreenHTTPHandler(r *gin.Engine, middleware *jwt.GinJWTMiddleware, splashscreenService splashscreen.Service) {

	handler := HTTPHandlerSplashScreen{splashscreenService}

	splashScreenReq := r.Group("/api/admin/splashscreen")
	splashScreenReq.Use(middleware.MiddlewareFunc())
	{
		splashScreenReq.GET("/list", handler.HandlerListSplashScreen)
		splashScreenReq.GET("/list-all", handler.HandlerListAllSplashScreen)
		splashScreenReq.POST("/add", handler.HandlerAddSplashScreen)
		splashScreenReq.POST("/update", handler.HandlerUpdateSplashScreen)
		splashScreenReq.POST("/delete", handler.HandlerDeleteSplashScreen)
		splashScreenReq.POST("/inactive", handler.HandlerInActiveSplashScreen)
		splashScreenReq.POST("/isactive", handler.HandlerActiveSplashScreen)
	}

}

// HandlerInActiveSplashScreen func inactive point
func (h *HTTPHandlerSplashScreen) HandlerInActiveSplashScreen(ctx *gin.Context) {
	var params models.RequestIDPoint
	_ = ctx.Bind(&params)

	err := h.splashscreenService.UpdateToinactive(params.ID)

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

// HTTPHandlerSplashScreen func active point
func (h *HTTPHandlerSplashScreen) HandlerActiveSplashScreen(ctx *gin.Context) {
	var params models.RequestIDPoint
	_ = ctx.Bind(&params)

	err := h.splashscreenService.UpdateToActive(params.ID)

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

// HandlerListSplashScreen Function
func (h *HTTPHandlerSplashScreen) HandlerListSplashScreen(ctx *gin.Context) {

	data := h.splashscreenService.GetSplashScreenList()

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

func (h *HTTPHandlerSplashScreen) HandlerListAllSplashScreen(ctx *gin.Context) {

	data := h.splashscreenService.GetSplashScreenListAll()

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

// HandlerAddSplashScreen Function
func (h *HTTPHandlerSplashScreen) HandlerAddSplashScreen(ctx *gin.Context) {

	var params models.SplashScreen
	_ = ctx.Bind(&params)

	InsertID, err := h.splashscreenService.AddSplashScreen(params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success add Splash Screen",
		"ID":      InsertID,
	})

}

// HandlerUpdateSplashScreen Function
func (h *HTTPHandlerSplashScreen) HandlerUpdateSplashScreen(ctx *gin.Context) {

	var params models.SplashScreen
	_ = ctx.Bind(&params)

	if params.ID == 0 || params.ImageURL == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params Splash Screen",
		})
		return
	}

	data, err := h.splashscreenService.UpdateSplashScreen(params)
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
			"message": "success update Splash Screen",
			"data":    data,
		})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "failed update Splash Screen",
			"data":    params,
		})
		return
	}

}

// HandlerDeleteSplashScreen Function
func (h *HTTPHandlerSplashScreen) HandlerDeleteSplashScreen(ctx *gin.Context) {

	var params models.RequestIDSplashScreen
	_ = ctx.Bind(&params)

	if params.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params Splash Screen",
		})
		return
	}

	err := h.splashscreenService.DeleteSplashScreen(params.ID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success delete Splash Screen",
	})

}
