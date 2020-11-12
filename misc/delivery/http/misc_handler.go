package http

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/rudyraharjo/emurojaah/misc"
)

type HttpHandlerMisc struct {
	MiscService misc.Service
}

func NewMiscHttpHandler(r *gin.Engine, middleware *jwt.GinJWTMiddleware, miscService misc.Service) {
	handler := HttpHandlerMisc{miscService}

	groupMisc := r.Group("/api/misc")
	groupMisc.Use(middleware.MiddlewareFunc())
	{
		groupMisc.POST("/address/province", handler.HandlerProvinceList)
	}
}

func (h *HttpHandlerMisc) HandlerProvinceList(ctx *gin.Context) {
	province := h.MiscService.GetProvinceList()

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "ok",
		"data":    province,
	})
}
