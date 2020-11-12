package http

import (
	"net/http"
	"strings"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/rudyraharjo/emurojaah/models"
	"github.com/rudyraharjo/emurojaah/user"
)

type HttpUserHandler struct {
	UserService user.Service
	MiddleWare  *jwt.GinJWTMiddleware
}

func NewUserHttpHandler(r *gin.Engine, middleware *jwt.GinJWTMiddleware, userService user.Service) {

	handler := &HttpUserHandler{userService, middleware}

	r.POST("/api/login", middleware.LoginHandler)
	r.POST("/api/register", handler.RegisterUserHandler)
	r.POST("/api/login/google", handler.HandlerLoginGoogle)

	userRoute := r.Group("/api/user")
	userRoute.Use(middleware.MiddlewareFunc())
	{
		userRoute.POST("/readsurahisdone", handler.GetReadIsDone)
		userRoute.POST("/survey/juzz", handler.HandlerSurveyJuzz)
		userRoute.POST("/profile", handler.GetProfile)
		userRoute.POST("/reward", handler.HandlerGetUserReward)
		userRoute.POST("/profile/overview", handler.HandlerGetUserOverview)
		userRoute.POST("/profile/update", handler.HandlerEditProfile)
		userRoute.POST("/destroy-token-firebase", handler.HandlerDestroyTokenFirebase)
	}

	useradminRoute := r.Group("/api/user/admin")
	useradminRoute.Use(cors.Default())
	useradminRoute.Use(middleware.MiddlewareFunc())
	{
		useradminRoute.POST("/list-member", handler.GetListMember)
		useradminRoute.POST("/list", handler.GetListAdmin)
		useradminRoute.POST("/profile", handler.GetProfileAdm)
		useradminRoute.GET("/total-member", handler.GetTotalMember)
		// useradminRoute.POST("/profile/update", handler.HandlerEditAdminProfile)
		useradminRoute.GET("/survey/juzz", handler.HandlerDashboardSurveyJuzz)
	}

	rewardReq := r.Group("/api/reward")
	rewardReq.Use(middleware.MiddlewareFunc())
	{
		rewardReq.POST("/all", handler.HandlerGetReward)
	}
}

// GetReadIsDone set surah read is done
func (h *HttpUserHandler) GetReadIsDone(ctx *gin.Context) {

	var params models.ReadIsDone
	_ = ctx.Bind(&params)

	isDone, err := h.UserService.CheckSurahReadingIsDone(params)

	var IsReading bool
	switch isDone {
	case 1:
		IsReading = true
	default:
		IsReading = false
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error binding",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"is_done": IsReading,
	})
	return
}

func (h *HttpUserHandler) RegisterUserHandler(ctx *gin.Context) {

	var params models.RequestRegister
	_ = ctx.Bind(&params)

	if params.FullName == "" || params.Email == "" || params.Password == "" || params.IDCity == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params ",
		})
		return
	}

	if params.TokenFirebase == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "Error connection please try again ",
		})
		return
	}

	if exCode := h.UserService.CheckIsAliasExist(strings.ToLower(params.Email)); exCode == 1 {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusUnauthorized,
			"message": "email already used",
		})
		return
	}

	userId, errCreateUser := h.UserService.RegisterUser(params)
	if errCreateUser != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error insert",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success register",
		"user_id": userId,
	})
}

func (h *HttpUserHandler) GetProfile(ctx *gin.Context) {
	var params models.RequestGetProfile
	_ = ctx.Bind(&params)

	if params.Alias == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params",
		})
		return
	}

	data, err := h.UserService.GetUserBasicInfoByAlias(strings.ToLower(params.Alias))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "ok",
		"data":    data,
	})
}

func (h *HttpUserHandler) HandlerSurveyJuzz(ctx *gin.Context) {
	var params models.RequestJuzzSurvey
	_ = ctx.Bind(&params)

	if params.UserId == 0 || params.Answer == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params",
		})
		return
	}

	err := h.UserService.AddJuzzSurvey(params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error insert",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success add survey",
	})
}

func (h *HttpUserHandler) HandlerGetUserReward(ctx *gin.Context) {
	var params models.RequestByUserId
	//var PointReward models.UserReward

	_ = ctx.Bind(&params)

	if params.UserId == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params",
		})
		return
	}

	totalPoint, points, err := h.UserService.GetUserPointReward(params.UserId)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":      http.StatusOK,
		"message":     "ok",
		"total_point": totalPoint,
		"data":        points,
	})

}

func (h *HttpUserHandler) HandlerGetUserOverview(ctx *gin.Context) {
	var params models.RequestByUserId
	_ = ctx.Bind(&params)

	if params.UserId == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params",
		})
		return
	}

	profile, _ := h.UserService.GetUserBasicInfoById(params.UserId)
	personalStat := h.UserService.GetUserReadStatistic(params.UserId)
	publicStat := h.UserService.GetPublicGroupStatistic(params.UserId)

	data := models.ResponseGetProfileOverview{
		Profile:                *profile,
		PersonalReadStatus:     personalStat,
		GlobalGroupTotalKhatam: publicStat,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "ok",
		"data":    data,
	})
}

func (h *HttpUserHandler) HandlerEditProfile(ctx *gin.Context) {
	var params models.RequestEditProfile
	_ = ctx.Bind(&params)

	if params.UserId == 0 || params.Email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params",
		})
		return
	}

	//reqLogin := models.RequestLogin{
	//	Alias:    params.Email,
	//	Password: params.Password,
	//	Type:     "EMAIL",
	//}

	//isPasswordCorrect, _ := h.UserService.LoginUser(reqLogin)
	//
	//if !isPasswordCorrect {
	//	ctx.JSON(http.StatusBadRequest, gin.H{
	//		"status":  http.StatusNotAcceptable,
	//		"message": "password incorrect",
	//	})
	//	return
	//}

	errUpdate := h.UserService.UpdateUserProfile(params)
	if errUpdate != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "data saved!",
	})

}

func (h *HttpUserHandler) HandlerDestroyTokenFirebase(ctx *gin.Context) {

	var params models.RequestByUserId
	_ = ctx.Bind(&params)

	if params.UserId == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusForbidden,
			"message": "missing params",
		})
		return
	}

	//DeleteTokenFirebase
	errDeleteToken := h.UserService.DeleteTokenFirebase(params.UserId)

	//fmt.Print("errDeleteToken => ", errDeleteToken)

	if errDeleteToken != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error",
		})
		return

	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Success Destroy Token Firebase",
	})
	return
}

func (h *HttpUserHandler) HandlerLoginGoogle(ctx *gin.Context) {

	var params models.RequestLoginWithGoogle
	_ = ctx.Bind(&params)

	if params.Alias == "" || params.Password != "k3yF0rL0g1n3mur0j44h" || params.TokenFirebase == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "missing params",
		})
		return
	}

	isAliasExist := h.UserService.CheckIsAliasExist(strings.ToLower(params.Alias))

	if isAliasExist == 2 {
		// register new user

		//hashedPass, _ := h.UserService.HashPassword(params.Password)
		//
		//reqRegister := models.RequestRegister{
		//	Email:     params.Alias,
		//	Password:  hashedPass,
		//	FullName:  params.FullName,
		//	Address:   "",
		//	BirthDate: "1990-01-01",
		//	Gender:    "M",
		//}
		//
		//_, _ = h.UserService.RegisterUser(reqRegister)

		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusForbidden,
			"message": "user not registered!",
		})
		return

	}

	reqLogin := models.RequestLogin{
		Alias:         strings.ToLower(params.Alias),
		Password:      params.Password,
		Type:          "GOOGLE",
		TokenFirebase: params.TokenFirebase,
	}

	code, success, fail := h.UserService.LoginWithRequest(reqLogin)

	if code == 200 {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "Ok",
			"token":   success.Token,
		})
		return
	} else if code == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "System error",
		})
		return
	} else {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusUnauthorized,
			"message": fail.Message,
		})
		return
	}
}

// GetProfileAdm GetProfileAdm
func (h *HttpUserHandler) GetProfileAdm(ctx *gin.Context) {

	var params models.RequestGetProfile
	_ = ctx.Bind(&params)

	if params.Alias == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params",
		})
		return
	}

	UserAdmin, UserDetAdmin, err := h.UserService.GetUserAdminInfoByAlias(strings.ToLower(params.Alias))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Data Not Found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":          http.StatusOK,
		"message":         "ok",
		"dataAdmin":       UserAdmin,
		"dataAdminDetail": UserDetAdmin,
	})
}

// GetListAdmin Function
func (h *HttpUserHandler) GetListAdmin(ctx *gin.Context) {

	ListAdmin := h.UserService.GetUserAdminList()
	data := models.ResponseListAdmin{
		UserListAdmin: ListAdmin,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   data,
	})
}

func (h *HttpUserHandler) GetListMember(ctx *gin.Context) {
	data := h.UserService.GetUserMemberList()

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   data,
	})

}

// GetTotalMember Function
func (h *HttpUserHandler) GetTotalMember(ctx *gin.Context) {

	data, err := h.UserService.GetTotalMember()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Data Not Found, Please Check Connection",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":       http.StatusOK,
		"total_member": data,
	})
}

// HandlerGetReward Function
func (h *HttpUserHandler) HandlerGetReward(ctx *gin.Context) {

	totalPoint, points, err := h.UserService.GetPointReward()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":      http.StatusOK,
		"message":     "ok",
		"total_point": totalPoint,
		"data":        points,
	})
}

// HandlerDashboardSurveyJuzz func
func (h *HttpUserHandler) HandlerDashboardSurveyJuzz(ctx *gin.Context) {

	data, _ := h.UserService.SurveyJuzzGrouping()
	if len(data) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "Data Not Found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success get list survey",
		"data":    data,
	})

}
