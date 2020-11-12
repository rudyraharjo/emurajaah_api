package http

import (
	"fmt"
	"net/http"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/rudyraharjo/emurojaah/group"
	"github.com/rudyraharjo/emurojaah/models"
)

type HttpHandlerGroup struct {
	GroupService group.Service
}

func NewGroupHttphandler(r *gin.Engine, middleware *jwt.GinJWTMiddleware, groupService group.Service) {
	handler := HttpHandlerGroup{groupService}

	groupReq := r.Group("/api/group")
	groupReq.Use(middleware.MiddlewareFunc())
	{
		groupReq.GET("/all-list", handler.HandlerGroupAllList)
		groupReq.POST("/list-type", handler.HandlerGroupListType)
		groupReq.POST("/update-member-reading", handler.HandlerUpdateMemberReadingIndex)
		groupReq.POST("/join", handler.HandlerJoinGroup)
		groupReq.POST("/join-new", handler.HandlerJoinGroupNew)
		groupReq.POST("/joingroupbyemail", handler.HandlerJoinGroupByEmail)
		groupReq.POST("/joinbulkgroup", handler.HandlerJoinBulkGroup)
		groupReq.POST("/list", handler.HandlerGroupList)
		groupReq.POST("/members", handler.HandlerGroupMemberList)
		groupReq.POST("/members/paging", handler.HandlerGroupMemberListWithOffsetAndLimit)
		groupReq.POST("/members-by-status", handler.HandlerGroupMemberListByStatus)
		groupReq.POST("/leave", handler.HandleExitGroup)
		groupReq.POST("/leave-reading-group", handler.HandleLeaveReadingGroup)
	}
}

func (h *HttpHandlerGroup) HandlerUpdateMemberReadingIndex(ctx *gin.Context) {
	h.GroupService.HandleUpdateMemberReadingIndex()

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success Update Member Reading + 1",
	})
}

func (h *HttpHandlerGroup) HandlerJoinGroupByEmail(ctx *gin.Context) {

	var params models.RequestJoinGroupByEmail
	if errBind := ctx.Bind(&params); errBind != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error binding",
		})
		return
	}

	if params.GroupID == 0 || params.Email == "" || params.GroupType == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params",
		})
		return
	}

	newMember, Group, successJoin, code := h.GroupService.JoinGroupByEmail(params.GroupID, params.Email, params.GroupType)

	if !successJoin && code == 1 {

		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Failed Join Please Check Connection & Try Again",
		})
		return

	} else if !successJoin && code == 2 {

		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Failed Join Email not found",
		})
		return

	} else if !successJoin && code == 3 {

		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Failed Join Group ID not found",
		})
		return

	} else {

		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "success join group",
			"data": gin.H{
				"id_group_member": newMember.ID,
				"group_id":        newMember.GroupID,
				"current_index":   newMember.CurrentIndex,
				"no_urut":         Group[0].NoGroupIndex,
				"max_member":      Group[0].MaxMember,
				"current_member":  Group[0].CurrentMember,
			},
		})
		return

	}

}

func (h *HttpHandlerGroup) HandlerJoinGroup(ctx *gin.Context) {
	var params models.RequestJoinGroup
	if errBind := ctx.Bind(&params); errBind != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error binding",
		})
		return
	}

	if params.UserID == 0 || params.GroupType == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params",
		})
		return
	}

	newMember, Group, successJoin := h.GroupService.JoinGroup(params.UserID, params.GroupType)

	if !successJoin {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Failed Join Please Check Connection & Try Again",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success join group",
		"data": gin.H{
			"id_group_member": newMember.ID,
			"group_id":        newMember.GroupID,
			"current_index":   newMember.CurrentIndex,
			"no_urut":         Group[0].NoGroupIndex,
			"max_member":      Group[0].MaxMember,
			"current_member":  Group[0].CurrentMember,
		},
	})
}

// HandlerJoinBulkGroup functionHandler
func (h *HttpHandlerGroup) HandlerJoinBulkGroup(ctx *gin.Context) {
	var params models.RequestJoinBulkGroup

	if errBind := ctx.Bind(&params); errBind != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error binding",
		})
		return
	}

	if len(params.DataBulk) > 0 && params.GroupID != 0 {

		Group, successJoin, code := h.GroupService.JoinBulkGroup(params)

		if !successJoin && code == 1 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Failed Join , Group ID Not Found or Full Member",
			})
			return
		} else if !successJoin && code == 2 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Failed Join params email not registered",
			})
			return
		} else {

			ctx.JSON(http.StatusOK, gin.H{
				"status":  http.StatusOK,
				"message": fmt.Sprintf("%d Success Join in Group ID %d", len(params.DataBulk), Group[0].NoGroupIndex),
			})
			return

		}

		// if code == 1 {
		// 	ctx.JSON(http.StatusOK, gin.H{
		// 		"status":  http.StatusOK,
		// 		"message": "Group ID Not found",
		// 		"ID":      params.GroupID,
		// 	})
		// 	return
		// }
		// if code == 2 {
		// 	ctx.JSON(http.StatusOK, gin.H{
		// 		"status":  http.StatusOK,
		// 		"message": "Failed Join Please check params",
		// 	})
		// 	return
		// }
		// ctx.JSON(http.StatusOK, gin.H{
		// 	"status":   http.StatusOK,
		// 	"code":     code,
		// 	"dataTemp": params,
		// })

	}

	// ctx.JSON(http.StatusBadRequest, gin.H{
	// 	"status":  http.StatusBadRequest,
	// 	"message": "missing params",
	// })
	// return

}

func (h *HttpHandlerGroup) HandlerJoinGroupNew(ctx *gin.Context) {

	var params models.RequestJoinGroup
	if errBind := ctx.Bind(&params); errBind != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error binding",
		})
		return
	}

	if params.UserID == 0 || params.GroupType == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params",
		})
		return
	}

	newMember, Group, successJoin, code := h.GroupService.JoinGroupNew(params.UserID, params.GroupType)

	if !successJoin && code == 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Failed Join Please Check Connection & Try Again",
		})
		return
	} else if !successJoin && code == 2 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Failed Join User Not Found",
		})
		return
	} else {

		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "success join group",
			"data": gin.H{
				"id_group_member": newMember.ID,
				"group_id":        newMember.GroupID,
				"current_index":   newMember.CurrentIndex,
				"no_urut":         Group[0].NoGroupIndex,
				"max_member":      Group[0].MaxMember,
				"current_member":  Group[0].CurrentMember,
			},
		})

	}

}

func (h *HttpHandlerGroup) HandlerGroupAllList(ctx *gin.Context) {

	data, _ := h.GroupService.GetListAllGroup()

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "ok",
		"data":    data,
	})

}

func (h *HttpHandlerGroup) HandlerGroupListType(ctx *gin.Context) {
	var req models.RequestListTypeGroupPaging
	_ = ctx.Bind(&req)

	if req.UserID == 0 || req.GroupType == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params",
		})
		return
	}

	data, _ := h.GroupService.GetListGroupType(req.UserID, req.GroupType, req.Offset, req.Limit)

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "ok",
		"data":    data,
	})

}

func (h *HttpHandlerGroup) HandlerGroupList(ctx *gin.Context) {
	var req models.RequestListGroup
	_ = ctx.Bind(&req)

	if req.UserID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params",
		})
		return
	}

	data, _ := h.GroupService.GetListUserGroup(req.UserID)

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "ok",
		"data":    data,
	})
}

func (h *HttpHandlerGroup) HandlerGroupMemberList(ctx *gin.Context) {
	var req models.RequestListMemberGroup
	_ = ctx.Bind(&req)

	if req.UserId == 0 || req.GroupId == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params",
		})
		return
	}

	data := h.GroupService.GetListGroupMemberWithName(req.GroupId)

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "ok",
		"data":    data,
	})
}

func (h *HttpHandlerGroup) HandlerGroupMemberListWithOffsetAndLimit(ctx *gin.Context) {
	var params models.RequestListMemberGroupWithOffsetLimit
	_ = ctx.Bind(&params)

	if params.UserId == 0 || params.GroupId == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params",
		})
		return
	}

	members := h.GroupService.GetListGroupMemberWithNameAndOffsetLimit(params.GroupId, params.Offset, params.Limit)
	totalMembers := h.GroupService.TotalGroupMember(params.GroupId)

	ctx.JSON(http.StatusOK, gin.H{
		"status":        http.StatusOK,
		"message":       "ok",
		"data":          members,
		"total_members": totalMembers,
	})
}

func (h *HttpHandlerGroup) HandlerGroupMemberListByStatus(ctx *gin.Context) {
	var params models.RequestListMemberGroupByStatus
	_ = ctx.Bind(&params)

	if params.UserId == 0 || params.GroupId == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params",
		})
		return
	}

	members := h.GroupService.GetListMemberGroupByStatusAndPaging(params.GroupId, params.Offset, params.Limit, params.Status)
	totalMembers := h.GroupService.TotalGroupMemberByStatus(params.GroupId, params.Status)

	ctx.JSON(http.StatusOK, gin.H{
		"status":        http.StatusOK,
		"message":       "ok",
		"data":          members,
		"total_members": totalMembers,
	})
}

// HandleLeaveReadingGroup func
func (h *HttpHandlerGroup) HandleLeaveReadingGroup(ctx *gin.Context) {
	var req models.RequestLeaveReading
	_ = ctx.Bind(&req)

	if req.ID == 0 || req.UserID == 0 || req.GroupID == 0 || req.GroupType == "" || req.ContentIndex == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params",
		})
		return
	}

	Count := h.GroupService.GetListGroupMembersByUserIDAndGroupID(req.UserID, req.GroupID)

	if Count == 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Anda tidak dapat meninggalkan group saat ini",
		})
		return
	}

	err := h.GroupService.LeaveReadingGroup(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "system error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success exit group",
	})

}
func (h *HttpHandlerGroup) HandleExitGroup(ctx *gin.Context) {
	var req models.RequestExitGroup
	_ = ctx.Bind(&req)

	if req.UserId == 0 || req.GroupId == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params",
		})
		return
	}

	/*
		data, _ := h.GroupService.GetListUserGroup(req.UserId)

		if len(data) == 1 {

			ctx.JSON(http.StatusBadRequest, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Anda tidak dapat meninggalkan group saat ini",
			})
			return

		}
	*/

	err := h.GroupService.ExitGroup(req.GroupId, req.UserId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "system error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success exit group",
	})
}
