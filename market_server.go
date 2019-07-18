package main

import (
	"bytes"
	"crypto/sha1"
	"database/sql"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/crypto/pbkdf2"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/logger"
)

var db *sql.DB

func main() {
	var verbose = flag.Bool("verbose", false, "print info level logs to stdout")
	flag.Parse()

	lf, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		logger.Fatalf("Failed to open log file: %v", err)
	}
	defer lf.Close()

	defer logger.Init("Logger", *verbose, true, lf).Close()

	logger.SetFlags(log.LstdFlags)
	logger.Info("I'm starting!")

	// Connect to the DB:
connect:

	db, err = sql.Open("mysql", "admin:veryStrongnt@tcp(127.0.0.1:3306)/market")
	if err != nil {
		logger.Errorf("Connect to the DB failed: %v", err)
		logger.Info("Tryna do it one more time...")
		goto connect
	} else {
		logger.Info("Connection has been established!")
	}
	defer db.Close()

	// Validate DSN data:
	err = db.Ping()
	if err != nil {
		logger.Fatalf("DSN is incorrect: %v", err)
	} else {
		logger.Info("Validation has been successfully passed!")
	}

	gin.ForceConsoleColor()

	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	router := gin.Default()

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:            "Flex Market",
		PrivKeyFile:      "../jwt.key",
		PubKeyFile:       "../jwt_pub.key",
		SigningAlgorithm: "RS512",
		Timeout:          time.Hour,
		MaxRefresh:       time.Hour,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			jti, _ := GenerateRandomString(10)

			logger.Infof("User = %v", data.(*User))

			_, err = db.Exec("update users set jti = ? where id = ?", jti, data.(*User).ID)

			if err != nil {
				logger.Errorf("[DB Query : SetJTI] %v; jti = %v", err, jti)
			}

			return jwt.MapClaims{
				"id": data.(*User).ID,
				// The "jti" (JWT ID) claim provides a unique identifier for the JWT.
				// https://tools.ietf.org/html/rfc7519#section-4.1.7
				"jti": jti,
			}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				ID:  int64(claims["id"].(float64)),
				JTI: claims["jti"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var login Login

			if err := c.ShouldBind(&login); err != nil {
				logger.Errorf("[Login] %v", err)
				return "", jwt.ErrMissingLoginValues
			}

			logger.Infof("[Login attempt] = %v", login.Email)
			user, err := FetchUserByEmail(login.Email)

			if err == nil {
				hash := pbkdf2.Key([]byte(login.Password), []byte(user.Salt), 4096, 256, sha1.New)

				if bytes.Equal(hash, user.Hash) {
					return user, nil
				}
			}

			logger.Infof("[Auth : Wrong Credentials]")
			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			user, err := FetchUserByID(data.(*User).ID)

			if err != nil {
				logger.Errorf("[DB Query : FetchUserByID : row.Scan] %v; id = %v", err, data.(*User).ID)
				c.JSON(
					http.StatusNotImplemented,
					gin.H{
						"status":  http.StatusNotImplemented,
						"message": err.Error(),
					})
			} else {
				if user.Role == "admin" && user.JTI == data.(*User).JTI {
					logger.Infof("Access granted for user [%v]", user.ID)
					return true
				}
			}
			logger.Infof("Access denied for user [%v]", user.ID)
			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",
		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	anon := router.Group("/api/v1/")
	{
		anon.GET("catalog/", FetchCatalog)
		anon.GET("sub_catalog/", FetchSubCatalog)
		anon.GET("products/", FetchAllProducts)
		anon.GET("products/:id", FetchSingleProduct)
		anon.POST("sign_up/", CreateUser)
	}

	auth := router.Group("/api/v1/auth/")
	{
		auth.POST("login/", authMiddleware.LoginHandler)
		auth.GET("refresh_token/", authMiddleware.RefreshHandler)
	}

	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("hello/", HelloHandler)
	}

	admin := router.Group("/api/v1/admin/").Use(authMiddleware.MiddlewareFunc())
	{
		admin.POST("products/", CreateProduct)
		admin.PUT("products/:id", UpdateProduct)
		admin.DELETE("products/:id", DeleteProduct)
	}

	router.Run()
}
