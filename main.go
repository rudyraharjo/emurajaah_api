package main

import (
	"fmt"
	"log"
	netHttp "net/http"
	"os"

	cors "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/robfig/cron"

	// Admin

	// Point
	_PointDelivery "github.com/rudyraharjo/emurojaah/admin/point/delivery/http"
	_PointRepo "github.com/rudyraharjo/emurojaah/admin/point/repository"
	_PointService "github.com/rudyraharjo/emurojaah/admin/point/service"

	// Banner
	_BannerDelivery "github.com/rudyraharjo/emurojaah/admin/banner/delivery/http"
	_BannerRepo "github.com/rudyraharjo/emurojaah/admin/banner/repository"
	_BannerService "github.com/rudyraharjo/emurojaah/admin/banner/service"

	//SplashScreen
	_SplashScreenDelivery "github.com/rudyraharjo/emurojaah/admin/splashscreen/delivery/http"
	_SplashScreenRepo "github.com/rudyraharjo/emurojaah/admin/splashscreen/repository"
	_SplashScreenService "github.com/rudyraharjo/emurojaah/admin/splashscreen/service"

	//Quote
	_QuoteDelivery "github.com/rudyraharjo/emurojaah/admin/quote/delivery/http"
	_QuoteRepo "github.com/rudyraharjo/emurojaah/admin/quote/repository"
	_QuoteService "github.com/rudyraharjo/emurojaah/admin/quote/service"

	//Country
	_CountryDelivery "github.com/rudyraharjo/emurojaah/admin/country/delivery/http"
	_CountryRepo "github.com/rudyraharjo/emurojaah/admin/country/repository"
	_CountryService "github.com/rudyraharjo/emurojaah/admin/country/service"

	// Province
	_ProvinceDelivery "github.com/rudyraharjo/emurojaah/admin/province/delivery/http"
	_ProvinceRepo "github.com/rudyraharjo/emurojaah/admin/province/repository"
	_ProvinceService "github.com/rudyraharjo/emurojaah/admin/province/service"

	// Cities
	_CityDelivery "github.com/rudyraharjo/emurojaah/admin/city/delivery/http"
	_CityRepo "github.com/rudyraharjo/emurojaah/admin/city/repository"
	_CityService "github.com/rudyraharjo/emurojaah/admin/city/service"

	//TermCondition
	_TermConditionDelivery "github.com/rudyraharjo/emurojaah/admin/termcondition/delivery/http"
	_TermConditionRepo "github.com/rudyraharjo/emurojaah/admin/termcondition/repository"
	_TermConditionService "github.com/rudyraharjo/emurojaah/admin/termcondition/service"

	// admGroup
	_admGroupDelivery "github.com/rudyraharjo/emurojaah/admin/group/delivery/http"
	_admGroupRepo "github.com/rudyraharjo/emurojaah/admin/group/repository"
	_admGroupService "github.com/rudyraharjo/emurojaah/admin/group/service"

	// end Admin

	_contentDelivery "github.com/rudyraharjo/emurojaah/content/delivery/http"
	_contentRepo "github.com/rudyraharjo/emurojaah/content/repository"
	_contentService "github.com/rudyraharjo/emurojaah/content/service"

	_credentialDelivery "github.com/rudyraharjo/emurojaah/credential/delivery/http"
	_credentialRepo "github.com/rudyraharjo/emurojaah/credential/repository"
	_credentialService "github.com/rudyraharjo/emurojaah/credential/service"

	_groupDelivery "github.com/rudyraharjo/emurojaah/group/delivery/http"
	_groupRepo "github.com/rudyraharjo/emurojaah/group/repository"
	_groupService "github.com/rudyraharjo/emurojaah/group/service"

	authMiddleware "github.com/rudyraharjo/emurojaah/middleware"

	_miscDelivery "github.com/rudyraharjo/emurojaah/misc/delivery/http"
	_miscRepo "github.com/rudyraharjo/emurojaah/misc/repository"
	_miscService "github.com/rudyraharjo/emurojaah/misc/service"

	_userDelivery "github.com/rudyraharjo/emurojaah/user/delivery/http"
	_userRepo "github.com/rudyraharjo/emurojaah/user/repository"
	_userService "github.com/rudyraharjo/emurojaah/user/service"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(`app-config.json`)
	err := viper.ReadInConfig()

	if err != nil {
		panic(err)
	}

	if viper.GetBool(`debug`) {
		fmt.Println("Service RUN on DEBUG mode")
	}

}

func main() {

	listenPort := viper.GetString(`server.port`)
	keyMiddleware := viper.GetString(`middleware.key`)

	// -- connect gcp
	hostName := viper.GetString(`database.host`)
	port := viper.GetString(`database.port`)
	userName := viper.GetString(`database.user`)
	dbName := viper.GetString(`database.db`)
	pass := viper.GetString(`database.pass`)
	connString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", hostName, port, userName, dbName, pass, "disable")

	// -- connect local
	// hostName := "localhost"
	// userName := "postgres"
	// dbName := "db_emurojaah"
	// connString := fmt.Sprintf("postgres://%s:@%s/%s?sslmode=disable", userName, hostName, dbName)

	dbConn, err := gorm.Open("postgres", connString)

	if err != nil {
		log.Fatal(err)
	}

	defer dbConn.Close()

	dbConn.DB().SetMaxIdleConns(10)
	dbConn.DB().SetMaxOpenConns(10)
	
	gin.SetMode(gin.DebugMode)

	r := gin.New()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	miscRepo := _miscRepo.NewMiscRepository(dbConn)
	userRepo := _userRepo.NewUserRepository(dbConn)
	groupRepo := _groupRepo.NewGroupRepository(dbConn)
	contentRepo := _contentRepo.NewContentRepository(dbConn)
	credentialRepo := _credentialRepo.NewCredentialRepository(dbConn)

	// admin

	// Point Repo
	pointRepo := _PointRepo.NewPointRepository(dbConn)
	pointService := _PointService.NewPointService(pointRepo)
	// Banner Repo
	bannerRepo := _BannerRepo.NewBannerRepository(dbConn)
	bannerService := _BannerService.NewBannerService(bannerRepo)
	//SplashScreen Repo
	splashscreenRepo := _SplashScreenRepo.NewSplashScreenRepository(dbConn)
	splashscreenService := _SplashScreenService.NewSplashScreenService(splashscreenRepo)
	//Quote Repo
	quoteRepo := _QuoteRepo.NewQuoteRepository(dbConn)
	quoteService := _QuoteService.NewQuoteService(quoteRepo)
	//Country Repo
	countryRepo := _CountryRepo.NewCountryRepository(dbConn)
	countryService := _CountryService.NewCountryService(countryRepo)
	// Province Repo
	provinceRepo := _ProvinceRepo.NewProvinceRepository(dbConn)
	provinceService := _ProvinceService.NewProvinceService(provinceRepo)
	// City Repo
	cityRepo := _CityRepo.NewCityRepository(dbConn)
	cityService := _CityService.NewCityService(cityRepo)

	// TermCondition
	termconditionRepo := _TermConditionRepo.NewTermConditionRepository(dbConn)
	termconditionService := _TermConditionService.NewTermConditionService(termconditionRepo)

	// admGroup
	admGroupRepo := _admGroupRepo.NewAdmGroupRepository(dbConn)
	admGroupService := _admGroupService.NewAdmGroupService(admGroupRepo)

	//end admin

	miscService := _miscService.NewMiscService(miscRepo)
	userService := _userService.NewUSerService(userRepo)
	groupService := _groupService.NewGroupService(groupRepo)
	contentService := _contentService.NewContentService(contentRepo, userRepo, groupRepo, bannerRepo, splashscreenRepo, termconditionRepo)
	credentialService := _credentialService.NewCredentailService(credentialRepo)

	middleware, errMid := authMiddleware.InitMiddleware(keyMiddleware, userService)

	if errMid != nil {
		log.Fatal(errMid)
	}

	// admin
	_PointDelivery.NewPointHTTPHandler(r, middleware, pointService)
	_BannerDelivery.NewBannerHTTPHandler(r, middleware, bannerService)
	_SplashScreenDelivery.NewSplashScreenHTTPHandler(r, middleware, splashscreenService)
	_QuoteDelivery.NewQuoteHTTPHandler(r, middleware, quoteService)
	_CountryDelivery.NewCountryHTTPHandler(r, middleware, countryService)
	_ProvinceDelivery.NewProvinceHTTPHandler(r, middleware, provinceService)
	_CityDelivery.NewCityHTTPHandler(r, middleware, cityService)
	_TermConditionDelivery.NewTermConditionHTTPHandler(r, middleware, termconditionService)
	_admGroupDelivery.NewAdmGroupHTTPHandler(r, middleware, admGroupService)

	// end admin

	_miscDelivery.NewMiscHttpHandler(r, middleware, miscService)
	_userDelivery.NewUserHttpHandler(r, middleware, userService)
	_groupDelivery.NewGroupHttphandler(r, middleware, groupService)
	_contentDelivery.NewContentHttpHandler(r, middleware, contentService, userService)
	_credentialDelivery.NewCredentialHttpHandler(r, middleware, credentialService)

	// Dont Use
	//contentService.HandlerAutoUpdateQuranPage()
	//contentService.HandlerUpdateIDCityAndIDProv()
	// END Dont Use

	// --------Start Corn-------- //

	//groupService.HandleNotifUserReadIsNotDone()

	corn := cron.New()

	// // // 59 * * * * * Every59Second Test

	corn.AddFunc("0 58 23 * * *", func() {
		fmt.Println("23.58")
		// Default @midnight
		// fmt.Print("HandleUpdateMemberReadingIndex")
		groupService.HandleUpdateMemberReadingIndex()
	})

	// // ----- HandleNotifUserReadIsNotDone ---- //

	corn.AddFunc("0 0 4 * * *", func() {
		fmt.Println("4.00")
		groupService.HandleNotifUserReadIsNotDone()
	})

	corn.AddFunc("0 30 10 * * *", func() {
		fmt.Println("10.30")
		groupService.HandleNotifUserReadIsNotDone()
	})

	corn.AddFunc("0 20 16 * * *", func() {
		fmt.Println("16.20")
		groupService.HandleNotifUserReadIsNotDone()
	})

	corn.AddFunc("0 00 21 * * *", func() {
		fmt.Println("21.00")
		groupService.HandleNotifUserReadIsNotDone()
	})

	// // // ----- End HandleNotifUserReadIsNotDone ---- //

	corn.Start()

	// --------End Cron-------- //

	if errHTTP := netHttp.ListenAndServe(listenPort, r); errHTTP != nil {
		log.Println(errHTTP)
		os.Exit(1)
	}

}
