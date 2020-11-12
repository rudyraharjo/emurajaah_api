package http

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"

	"github.com/rudyraharjo/emurojaah/user"

	"github.com/rudyraharjo/emurojaah/content"
	"github.com/rudyraharjo/emurojaah/models"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
)

type HttpHandlerContent struct {
	ContentService content.Service
	UserService    user.Service
}

func NewContentHttpHandler(r *gin.Engine, middleware *jwt.GinJWTMiddleware, contentService content.Service, userService user.Service) {
	handler := HttpHandlerContent{contentService, userService}

	termReq := r.Group("/api/content/term")
	termReq.GET("/", handler.HandlerTermCondition)

	splashScreenReq := r.Group("/api/content/splash_screen")
	splashScreenReq.GET("/list", handler.HandlerListSplashScreen)

	contentKotaReq := r.Group("/api/content/ibukota")
	contentKotaReq.GET("/list", handler.HandlerListIbuKota)

	contentRegionRoute := r.Group("/api/content/region")
	contentRegionRoute.Use()
	{
		contentRegionRoute.GET("/list-prov", handler.HandlerListProvinces)
		contentRegionRoute.POST("/list-city", handler.HandlerListCities)
		contentRegionRoute.GET("/total-user-by-province", handler.HanlerTotalUserByProvince)
	}

	contentBoardingRoute := r.Group("/api/content/boarding-page")
	contentBoardingRoute.GET("/list", handler.HandlerListBoardingPage)
	contentBoardingRoute.Use(middleware.MiddlewareFunc())
	{
		contentBoardingRoute.GET("/list-all", handler.HandlerListAllBoardingPage)
		contentBoardingRoute.POST("/add", handler.HandlerAddBoardingPage)
		contentBoardingRoute.POST("/update", handler.HandlerUpdateBoardingPage)
		contentBoardingRoute.POST("/deactivate", handler.HandlerDeactivate)
		contentBoardingRoute.POST("/update-status", handler.HandlerUpdateStatus)
	}

	quranReq := r.Group("/api/content/quran")
	quranReq.Use(middleware.MiddlewareFunc())
	{
		quranReq.POST("/all", handler.HandlerQuran)
		quranReq.POST("/add", handler.HandlerAddSurahQuran)
		quranReq.POST("/refresh", handler.HandlerRefreshQuran)
		quranReq.POST("/surah", handler.HandleListSurah)
		quranReq.POST("/save", handler.HandleSaveReadingQuran)
		quranReq.POST("/paging", handler.HandleListPagingQuran)
		quranReq.POST("/detail-per-surah", handler.HandleDetailPerSurahQuran)
	}

	homePageReq := r.Group("/api/content/home")
	homePageReq.Use(middleware.MiddlewareFunc())
	{
		homePageReq.POST("/all", handler.HandlerHomePageContent)
		homePageReq.POST("/banners", handler.HandlerHomePageContentBanner)
		homePageReq.POST("/groups", handler.HandlerHomePageContentGroups)
		homePageReq.POST("/quotes", handler.HandlerHomePageContentQuotes)
		homePageReq.POST("/global-group-status", handler.HandlerHomePageContentGlobalGroupStatus)
	}

	contentadminRoute := r.Group("/api/content/admin")
	contentadminRoute.Use(cors.Default())
	contentadminRoute.Use(middleware.MiddlewareFunc())
	{
		contentadminRoute.POST("/global-total-khatam", handler.GetGlobalTotalKhatam)
		//contentadminRoute.GET("/")
	}

}

func (h *HttpHandlerContent) GetGlobalTotalKhatam(ctx *gin.Context) {

	TotalKhatam, err := h.ContentService.GetTotalKhatamAllGroup()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   TotalKhatam,
	})
}

func (h *HttpHandlerContent) HandlerAddSurahQuran(ctx *gin.Context) {
	var req models.AddQuranRequest
	_ = ctx.Bind(&req)

	err := h.ContentService.AddQuranSurah(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success add surah",
	})
}

func (h *HttpHandlerContent) HandlerRefreshQuran(ctx *gin.Context) {
	err := h.ContentService.GetQuranFromAPI()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success refresh surah",
	})
}

func (h *HttpHandlerContent) HandleListSurah(ctx *gin.Context) {
	var params models.RequestQuran
	_ = ctx.Bind(&params)

	tokenString := ctx.Request.Header.Get("Authorization")

	// token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
	// 	if jwt.GetSigningMethod("HS256") != token.Method {
	// 		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	// 	}

	// 	return []byte("secret"), nil
	// })

	if params.Type == "" || params.Index == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params",
		})
		return
	}

	data := h.ContentService.GetSurahByCategory(params)

	ctx.JSON(http.StatusOK, gin.H{
		"status":      http.StatusOK,
		"message":     "ok",
		"data":        data,
		"tokenString": tokenString,
	})

}

func (h *HttpHandlerContent) HandleListPagingQuran(ctx *gin.Context) {
	var params models.RequestQuranPaging
	_ = ctx.Bind(&params)

	if params.Type == "" || params.Index == 0 || params.Limit == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params",
		})
		return
	}

	data := h.ContentService.GetQuranByCategoryWithPaging(params)

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "ok",
		"data":    data,
	})
}

func (h *HttpHandlerContent) HandleSaveReadingQuran(ctx *gin.Context) {
	var params models.RequestSaveReadingQuran
	_ = ctx.Bind(&params)

	if params.ID == 0 || params.UserID == 0 || params.ContentIndex == 0 || params.GroupID == 0 || params.GroupType == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params",
		})
		return
	}

	err := h.ContentService.SaveUserReadingActivity(params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "System Error Please Try Again",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "successfully save reading history",
	})

}

// HandlerHomePageContent func home
func (h *HttpHandlerContent) HandlerHomePageContent(ctx *gin.Context) {
	var params models.RequestByUserId
	_ = ctx.Bind(&params)

	if params.UserId == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params",
		})
		return
	}

	if params.TokenFirebase != "" {

		checkTokenExits := h.UserService.CheckTokenFirebase(params.UserId)

		//fmt.Printf("Token %s ", checkTokenExits)

		if checkTokenExits == 0 {
			//Add Token
			_, err := h.UserService.AddTokenUser(params.UserId, params.TokenFirebase)
			if err != nil {
				fmt.Print(err)
			}
		} else {

			deleteToken := h.UserService.DeleteTokenFirebase(params.UserId)
			if deleteToken != nil {
				fmt.Print(deleteToken)
			}

			_, AddTokenUser := h.UserService.AddTokenUser(params.UserId, params.TokenFirebase)
			if AddTokenUser != nil {
				fmt.Print(AddTokenUser)
			}
		}
	}

	resp := h.ContentService.HandlerHomePageContent(params.UserId)

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "ok",
		"data":    resp,
	})
}

func (h *HttpHandlerContent) HandlerHomePageContentGlobalGroupStatus(ctx *gin.Context) {
	var params models.RequestByUserId
	_ = ctx.Bind(&params)

	if params.UserId == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params",
		})
		return
	}

	if params.TokenFirebase != "" {

		checkTokenExits := h.UserService.CheckTokenFirebase(params.UserId)

		//fmt.Printf("Token %s ", checkTokenExits)

		if checkTokenExits == 0 {
			//Add Token
			_, err := h.UserService.AddTokenUser(params.UserId, params.TokenFirebase)
			if err != nil {
				fmt.Print(err)
			}
		} else {

			deleteToken := h.UserService.DeleteTokenFirebase(params.UserId)
			if deleteToken != nil {
				fmt.Print(deleteToken)
			}

			_, AddTokenUser := h.UserService.AddTokenUser(params.UserId, params.TokenFirebase)
			if AddTokenUser != nil {
				fmt.Print(AddTokenUser)
			}
		}
	}

	resp := h.ContentService.HandlerHomePageContentGlobalGroupStatus(params.UserId)

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "ok",
		"data":    resp,
	})
}

func (h *HttpHandlerContent) HandlerHomePageContentQuotes(ctx *gin.Context) {
	var params models.RequestByUserId
	_ = ctx.Bind(&params)

	if params.UserId == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params",
		})
		return
	}

	resp := h.ContentService.HandlerHomePageContentQuotes(params.UserId)

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "ok",
		"data":    resp,
	})
}

func (h *HttpHandlerContent) HandlerHomePageContentGroups(ctx *gin.Context) {

	var params models.RequestByUserId
	_ = ctx.Bind(&params)

	if params.UserId == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params",
		})
		return
	}

	resp := h.ContentService.HandlerHomePageContentGroups(params.UserId)

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "ok",
		"data":    resp,
	})

}

func (h *HttpHandlerContent) HandlerHomePageContentBanner(ctx *gin.Context) {
	var params models.RequestByUserId
	_ = ctx.Bind(&params)

	if params.UserId == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params",
		})
		return
	}

	resp := h.ContentService.HandlerHomePageContentBanner(params.UserId)

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "ok",
		"data":    resp,
	})

}

// HandleDetailPerSurahQuran Function
func (h *HttpHandlerContent) HandleDetailPerSurahQuran(ctx *gin.Context) {
	var req models.RequestQuranBySurahId
	_ = ctx.Bind(&req)

	data := h.ContentService.GetSurahByIDSurah(req.SurahID)

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "ok",
		"data":    data,
	})
}

// HandlerQuran function
func (h *HttpHandlerContent) HandlerQuran(ctx *gin.Context) {

	data := h.ContentService.GetAllQuran()

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "ok",
		"data":    data,
	})

}

func (h *HttpHandlerContent) HandlerTermCondition(ctx *gin.Context) {
	data := h.ContentService.GetTermCondition()
	if len(data) > 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "ok",
			"data":    data,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "no data",
	})

}

// HandlerListSplashScreen Function
func (h *HttpHandlerContent) HandlerListSplashScreen(ctx *gin.Context) {

	data := h.ContentService.GetSplashScreenList()

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
			"data":    data,
		})
		return
	}

}

// HandlerListAllBoardingPage function
func (h *HttpHandlerContent) HandlerListAllBoardingPage(ctx *gin.Context) {

	data := h.ContentService.GetListAllBoarding()

	ctx.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   data,
	})

}

// HandlerListBoardingPage func
func (h *HttpHandlerContent) HandlerListBoardingPage(ctx *gin.Context) {

	data := h.ContentService.GetListBoardingPageIsActive()

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "ok",
		"data":    data,
	})
}

// HandlerAddBoardingPage func
func (h *HttpHandlerContent) HandlerAddBoardingPage(ctx *gin.Context) {

	var params models.BoardingPage
	_ = ctx.Bind(&params)

	InsertID, err := h.ContentService.AddBoardingPage(params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success add Boarding Page",
		"ID":      InsertID,
	})

}

// HandlerUpdateBoardingPage function
func (h *HttpHandlerContent) HandlerUpdateBoardingPage(ctx *gin.Context) {

	var params models.BoardingPage
	_ = ctx.Bind(&params)

	if params.ID == 0 || params.ImageURL == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params Boarding page",
		})
		return
	}

	data, err := h.ContentService.UpdateBoardingPage(params)
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
			"message": "success update Boarding page",
			"data":    data,
		})
		return
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "failed update Boarding page",
			"data":    params,
		})
		return
	}

}

// HandlerUpdateStatus func
func (h *HttpHandlerContent) HandlerUpdateStatus(ctx *gin.Context) {
	var params models.RequestIDBoardingPage
	_ = ctx.Bind(&params)

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success update Boarding page",
	})
	return

}

// HandlerDeactivate function
func (h *HttpHandlerContent) HandlerDeactivate(ctx *gin.Context) {
	var params models.RequestIDBoardingPage
	_ = ctx.Bind(&params)

	if params.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusForbidden,
			"message": "missing params ID",
		})
		return
	}

	err := h.ContentService.DeactivateBoardingPage(params.ID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "system error",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "success delete BoardingPage",
	})
}

// HandlerListProvinces func list Provinces
func (h *HttpHandlerContent) HandlerListProvinces(ctx *gin.Context) {

	data := h.ContentService.GetListProvinces()

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
			"data":    data,
		})
		return
	}
}

// HandlerListCities func list Cities
func (h *HttpHandlerContent) HandlerListCities(ctx *gin.Context) {

	provinceID := 0
	var params models.RequestIDCityByProvinceID
	_ = ctx.Bind(&params)

	if params.ID != 0 {
		provinceID = params.ID
	}

	data := h.ContentService.GetListCities(provinceID)

	if len(data) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "no data",
			"data":    data,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "ok",
		"data":    data,
	})
}

// HandlerListIbuKota func get list ibukota
func (h *HttpHandlerContent) HandlerListIbuKota(ctx *gin.Context) {

	data := h.ContentService.GetListIbukota()

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
			"data":    data,
		})
		return
	}

}

// HanlerTotalUserByProvince function for map
func (h *HttpHandlerContent) HanlerTotalUserByProvince(ctx *gin.Context) {
	data, IsSuccess := h.ContentService.GetTotalUserByProvince()

	if !IsSuccess {
		ctx.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "no data",
			"data":    data,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "ok",
		"data":    data,
	})
}
