package http

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/rudyraharjo/emurojaah/admin/termcondition"
	"github.com/rudyraharjo/emurojaah/models"
)

// HTTPHandlerTermCondition struct
type HTTPHandlerTermCondition struct {
	termconditionService termcondition.Service
}

// NewTermConditionHTTPHandler func handler route
func NewTermConditionHTTPHandler(r *gin.Engine, middleware *jwt.GinJWTMiddleware, termconditionService termcondition.Service) {

	handler := HTTPHandlerTermCondition{termconditionService}

	termConditionReq := r.Group("/api/admin/term-condition")
	termConditionReq.Use(middleware.MiddlewareFunc())
	{
		termConditionReq.GET("/list", handler.HandlerList)
		termConditionReq.POST("/update", handler.HandlerUpdate)
	}

}

// HandlerUpdate function
func (h *HTTPHandlerTermCondition) HandlerUpdate(ctx *gin.Context) {

	var params models.TermCondition
	_ = ctx.Bind(&params)

	if params.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params ID",
		})
		return
	}

	data, err := h.termconditionService.Update(params)
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
			"message": "success update Term&Condition",
			"data":    data,
		})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "failed update Term&Condition",
			"data":    params,
		})
		return
	}

}

// HandlerList Function
func (h *HTTPHandlerTermCondition) HandlerList(ctx *gin.Context) {

	data := h.termconditionService.GetListTermConditions()

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
