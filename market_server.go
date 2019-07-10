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

	db, err = sql.Open("mysql", "Sevada:LAliMDVIG24\\#@tcp(127.0.0.1:3306)/market")
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
		PrivKeyFile:      "jwtRS256.key",
		PubKeyFile:       "jwtRS256pub.key",
		SigningAlgorithm: "RS256",
		Timeout:          time.Hour,
		MaxRefresh:       time.Hour,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				jti, _ := GenerateRandomString(10)

				logger.Infof("User = %v", v)

				_, err = db.Exec("update users set jti = ? where id = ?", jti, v.ID)

				if err != nil {
					logger.Errorf("[DB Query : SetJTI] %v; jti = %v", err, jti)
				}

				return jwt.MapClaims{
					"id": v.ID,
					// The "jti" (JWT ID) claim provides a unique identifier for the JWT.
					// https://tools.ietf.org/html/rfc7519#section-4.1.7
					"jti": jti,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				ID:  claims["id"].(int64),
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

			user, err := FetchUser(login.Email)

			hash := pbkdf2.Key([]byte(login.Password), []byte(user.Salt), 4096, 256, sha1.New)

			if bytes.Equal(hash, user.Hash) && err == nil {
				return user, nil
			}

			logger.Infof("[Auth : Wrong Credentials]")

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			var currentJTI string

			row := db.QueryRow("select jti from users where id = ?", data.(*User).ID)
			err := row.Scan(&currentJTI)

			if err != nil {
				logger.Errorf("[DB Query : FetchJTI : row.Scan] %v; jti = %v", err, currentJTI)
				c.JSON(
					http.StatusNotImplemented,
					gin.H{
						"status":  http.StatusNotImplemented,
						"message": err.Error(),
					})
			} else {
				if v, ok := data.(*User); ok && v.Email == "test@test.com" && v.JTI == currentJTI {
					return true
				}
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

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
		anon.GET("products/", FetchAllProducts)
		anon.GET("products/:id", FetchSingleProduct)
		anon.POST("sign_up/", CreateUser)
	}

	auth := router.Group("/api/v1/auth/")
	{
		auth.POST("login/", authMiddleware.LoginHandler)
		auth.GET("refresh_token/", authMiddleware.RefreshHandler)
		auth.GET("hello/", HelloHandler)
	}

	auth.Use(authMiddleware.MiddlewareFunc())

	admin := router.Group("/api/v1/admin/")
	{
		admin.POST("products/", CreateProduct)
		admin.PUT("products/:id", UpdateProduct)
		admin.DELETE("products/:id", DeleteProduct)
	}

	admin.Use(authMiddleware.MiddlewareFunc())

	router.Run()
}
