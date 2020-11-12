package http

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/rudyraharjo/emurojaah/admin/group"
	"github.com/rudyraharjo/emurojaah/models"
)

// HTTPHandlerGroup
type HTTPHandlerAdmGroup struct {
	admGroupService group.Service
}

// NewAdmGroupHTTPHandler func
func NewAdmGroupHTTPHandler(r *gin.Engine, middleware *jwt.GinJWTMiddleware, admGroupService group.Service) {

	handler := HTTPHandlerAdmGroup{admGroupService}

	groupReq := r.Group("/api/admin/group")
	groupReq.Use(middleware.MiddlewareFunc())
	{
		groupReq.GET("/list", handler.HandlerListGroup)
		groupReq.POST("/detail-member", handler.HandlerListGroupMember)
		groupReq.POST("/add", handler.HandlerAddGroup)
		groupReq.POST("/send-notif-bygroup", handler.HandlerSendNotif)
		groupReq.GET("/refresh-duplicate-group", handler.HandlerRefreshGroup) // delete duplicate group_member
	}
}

// HandlerAddGroup func
func (h *HTTPHandlerAdmGroup) HandlerAddGroup(ctx *gin.Context) {

	var params models.RequestTypeGroup
	_ = ctx.Bind(&params)

	if params.GroupType == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params group type",
		})
		return
	}

	data := h.admGroupService.GenerateAddGroup(params.GroupType)

	if len(data) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params group type",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "ok",
		"data":    data,
	})

}

// HandlerRefreshGroup function
func (h *HTTPHandlerAdmGroup) HandlerRefreshGroup(ctx *gin.Context) {

	c, err := h.admGroupService.DeleteDuplicateGroupMembers()

	if c == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "failed refresh group duplicate",
		})
		return
	}

	if c == 2 {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Duplicate Data Group Members Not Found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "OK",
	})

}

// HandlerListGroup Func
func (h *HTTPHandlerAdmGroup) HandlerListGroup(ctx *gin.Context) {

	data := h.admGroupService.GetListGroups()

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "ok",
		"data":    data,
	})

}

func (h *HTTPHandlerAdmGroup) HandlerListGroupMember(ctx *gin.Context) {

	var params models.RequestIDGroup
	_ = ctx.Bind(&params)

	if params.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params banner",
		})
		return
	}

	data := h.admGroupService.GetListGroupMember(params.ID)

	if len(data) > 0 {

		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "ok",
			"data":    data,
		})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "no data",
		})
		return
	}

}

// HandlerSendNotif function
func (h *HTTPHandlerAdmGroup) HandlerSendNotif(ctx *gin.Context) {

	var params models.RequestIDGroup
	_ = ctx.Bind(&params)

	if params.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params ID",
		})
		return
	}

	code, err := h.admGroupService.SendNotifBelomBaca(params.ID)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "ada sesuatu dengan dirinya..",
		})
		return
	}

	if code != 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "ada sesuatu dengan dirinya..",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "ok",
	})
}
