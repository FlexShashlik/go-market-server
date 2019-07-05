package main

import (
	"database/sql"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/logger"
)

var identityKey = "id"

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
		IdentityKey:      identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				jti, _ := GenerateRandomString(10)

				_, err = db.Exec("update users set jti = ? where id = ?", jti, 1)

				if err != nil {
					logger.Errorf("[DB Query : SetJTI] %v; jti = %v", err, jti)
				}

				return jwt.MapClaims{
					identityKey: v.UserName,
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
				UserName: claims["id"].(string),
				JTI: claims["jti"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals Login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password

			if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
				return &User{
					UserName:  userID,
					LastName:  "Flexyan",
					FirstName: "Sevada",
				}, nil
			}

			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			var currentJTI string

			row := db.QueryRow("select jti from users where id = ?", 1)
			err := row.Scan(&currentJTI)

			if err != nil {
				logger.Errorf("[DB Query : FetchJTI : row.Scan] %v; jti = %v", err, currentJTI)
				c.JSON(
					http.StatusNotImplemented,
					gin.H{
						"status":  http.StatusNotImplemented,
						"message": err.Error()})
			} else {
				if v, ok := data.(*User); ok && v.UserName == "admin" && v.JTI == currentJTI {
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

	router.POST("/login", authMiddleware.LoginHandler)

	router.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	auth := router.Group("/auth")
	// Refresh time can be longer than token timeout
	auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.GET("/hello", HelloHandler)
	}

	v1 := router.Group("/api/v1/")
	{
		v1.POST("products/", CreateProduct)
		v1.GET("products/", FetchAllProducts)
		v1.GET("products/:id", FetchSingleProduct)
		v1.PUT("products/:id", UpdateProduct)
		v1.DELETE("products/:id", DeleteProduct)
	}
	router.Run()
}
