package middleware

import (
	"fmt"
	"strconv"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/rudyraharjo/emurojaah/models"
	"github.com/rudyraharjo/emurojaah/user"
)

func InitMiddleware(key string, userService user.Service) (*jwt.GinJWTMiddleware, error) {

	var identityKey = "id"

	authMiddleWare, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "authentication zone",
		Key:         []byte(key),
		Timeout:     time.Hour * 720,
		MaxRefresh:  time.Hour * 720,
		IdentityKey: identityKey,

		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.RequestLogin); ok {
				return jwt.MapClaims{
					identityKey: v.Alias,
				}
			}
			return jwt.MapClaims{}
		},

		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &models.RequestLogin{
				Alias: claims["id"].(string),
			}
		},

		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals models.RequestLogin

			if errAuthenticator := c.ShouldBind(&loginVals); errAuthenticator != nil {
				return "", jwt.ErrMissingLoginValues
			}

			if loginVals.Type == "EMAIL" {

				if loginVals.TokenFirebase == "" {
					return "", jwt.ErrMissingLoginValues
				}

				isMatch, userId := userService.LoginUser(loginVals)

				if isMatch && userId != 0 {

					// CheckTokenisExits
					tokenExits := userService.CheckTokenFirebase(userId)

					//fmt.Print("tokenExits Type Email => ", tokenExits)

					if tokenExits == 0 {
						//Add Token
						_, AddTokenUser := userService.AddTokenUser(userId, loginVals.TokenFirebase)
						if AddTokenUser != nil {
							fmt.Print(AddTokenUser)
						}
						//fmt.Println(UserID)
					} else {

						deleteToken := userService.DeleteTokenFirebase(userId)
						if deleteToken != nil {
							fmt.Print(deleteToken)
						}

						_, AddTokenUser := userService.AddTokenUser(userId, loginVals.TokenFirebase)
						if AddTokenUser != nil {
							fmt.Print(AddTokenUser)
						}

					}

					return &models.RequestLogin{
						Alias: strconv.Itoa(userId),
					}, nil
				}

			} else if loginVals.Type == "User_Backend" {

				isMatch, userID := userService.LoginUserAdmin(loginVals)

				if isMatch && userID != 0 {

					return &models.RequestLogin{
						Alias: strconv.Itoa(userID),
					}, nil

				}

			} else {

				if loginVals.TokenFirebase == "" {
					return "", jwt.ErrMissingLoginValues
				}

				isMatch, userId := userService.LoginUserWithGoogle(loginVals)

				if isMatch && userId != 0 {

					// CheckTokenisExits
					tokenExits := userService.CheckTokenFirebase(userId)

					//fmt.Print("tokenExits No Type=> ", tokenExits)

					if tokenExits == 0 {
						//Add Token
						_, AddTokenUser := userService.AddTokenUser(userId, loginVals.TokenFirebase)
						if AddTokenUser != nil {
							fmt.Print(AddTokenUser)
						}
						//fmt.Println(UserID)
					} else {

						deleteToken := userService.DeleteTokenFirebase(userId)
						if deleteToken != nil {
							fmt.Print(deleteToken)
						}

						_, AddTokenUser := userService.AddTokenUser(userId, loginVals.TokenFirebase)
						if AddTokenUser != nil {
							fmt.Print(AddTokenUser)
						}
					}

					return &models.RequestLogin{
						Alias: strconv.Itoa(userId),
					}, nil
				}
			}

			return nil, jwt.ErrFailedAuthentication
		},

		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*models.RequestLogin); ok {

				convertedID, errCnf := strconv.Atoi(v.Alias)

				if errCnf != nil {
					fmt.Println(errCnf)
					return false
				}

				res := userService.CheckIsUserExist(convertedID)

				if res == 1 {
					return true
				}

			}
			return false
		},

		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": "message",
			})
		},

		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	return authMiddleWare, err

}
