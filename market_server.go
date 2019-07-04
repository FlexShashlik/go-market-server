package main

import (
	"time"
	"database/sql"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/logger"
)

// Login struct
type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

var identityKey = "id"

// User demo
type User struct {
	UserName  string
	FirstName string
	LastName  string
}

// Product struct
type Product struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Price          string `json:"price"`
	ImageExtension string `json:"image_extension"`
}

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
		Realm:       "Flex Market",
		PrivKeyFile: "jwtRS256.key",
		PubKeyFile: "jwtRS256pub.key",
		SigningAlgorithm: "RS512",
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				UserName: claims["id"].(string),
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
			if v, ok := data.(*User); ok && v.UserName == "admin" {
				return true
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

// UploadImage loads image from post request to memory
func UploadImage(product *Product, c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		logger.Errorf("[DB Query : CreateProduct : FormFile()] %v", err)
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": err.Error()})
		return
	}

	filename := "images/" + strconv.FormatInt(product.ID, 10) + "." + product.ImageExtension
	if err := c.SaveUploadedFile(file, filename); err != nil {
		logger.Errorf("[DB Query : CreateProduct : SaveUploadedFile()] %v", err)
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": err.Error()})
		return
	}
}

/* Handlers */

func HelloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	user, _ := c.Get(identityKey)
	c.JSON(200, gin.H{
		"userID":   claims["id"],
		"userName": user.(*User).UserName,
		"text":     "Hello World.",
	})
}

// CreateProduct creates new product
func CreateProduct(c *gin.Context) {
	product := Product{
		Name:           c.PostForm("name"),
		Price:          c.PostForm("price"),
		ImageExtension: c.PostForm("image_extension"),
	}

	result, err := db.Exec(
		"insert into products (name, price, image_extension) values (?, ?, ?)",
		product.Name,
		product.Price,
		product.ImageExtension)

	if err != nil {
		logger.Errorf("[DB Query : CreateProduct] %v", err)
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": err.Error()})
	} else {
		product.ID, err = result.LastInsertId()
		if err != nil {
			logger.Errorf("[DB Query : CreateProduct : LastInsertID] %v; ", err)
			c.JSON(
				http.StatusNotImplemented,
				gin.H{
					"status":  http.StatusNotImplemented,
					"message": err.Error()})
		} else {
			UploadImage(&product, c)

			logger.Infof("Product [%v] created", product)
			c.JSON(http.StatusCreated, gin.H{"productID": product.ID})
		}
	}
}

// FetchAllProducts fetches all products
func FetchAllProducts(c *gin.Context) {
	var products []Product

	rows, err := db.Query("select id, name, price, image_extension from products")
	if err != nil {
		logger.Errorf("[DB Query : FetchAllProducts] %v", err)
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": err.Error()})
	} else {
		for rows.Next() {
			p := Product{}

			err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.ImageExtension)

			if err != nil {
				logger.Errorf("[DB Query : FetchAllProducts : rows.Scan] %v", err)
				continue
			}
			products = append(products, p)
		}
		logger.Infof("Products fetched")
		c.JSON(http.StatusOK, products)
	}
}

// FetchSingleProduct fetches single product
func FetchSingleProduct(c *gin.Context) {
	var product Product
	productID := c.Param("id")

	row := db.QueryRow("select id, name, price, image_extension from products where id = ?", productID)
	err := row.Scan(&product.ID, &product.Name, &product.Price, &product.ImageExtension)

	if err != nil {
		logger.Errorf("[DB Query : FetchSingleProduct : row.Scan] %v; productID = %v", err, productID)
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": err.Error()})
	} else {
		logger.Infof("Product %v fetched", productID)
		c.JSON(http.StatusOK, product)
	}
}

// UpdateProduct updates product
func UpdateProduct(c *gin.Context) {
	product := Product{}

	product.ID, _ = strconv.ParseInt(c.Param("id"), 10, 64)
	product.Name = c.PostForm("name")
	product.Price = c.PostForm("price")
	product.ImageExtension = c.PostForm("image_extension")

	_, err := db.Exec(
		"update products SET name = ?, price = ?, image_extension = ? where id = ?",
		product.Name,
		product.Price,
		product.ImageExtension,
		product.ID)

	if err != nil {
		logger.Errorf("[DB Query : UpdateProduct] %v; product = %v", err, product)
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": err.Error()})
	} else {
		UploadImage(&product, c)

		logger.Infof("Product %v updated to name = %v, price = %v, image_extension = %v", product.ID, product.Name, product.Price, product.ImageExtension)
		c.JSON(
			http.StatusOK,
			gin.H{
				"status":  http.StatusOK,
				"message": "Product updated successfully!"})
	}
}

// DeleteProduct deletes product
func DeleteProduct(c *gin.Context) {
	productID, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	_, err := db.Exec("delete from products where id = ?", productID)

	if err != nil {
		logger.Errorf("[DB Query : DeleteProduct] %v; productID = %v", err, productID)
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": err.Error()})
	} else {
		logger.Infof("Product %v deleted", productID)
		c.JSON(
			http.StatusOK,
			gin.H{
				"status":  http.StatusOK,
				"message": "Product deleted successfully!"})
	}
}
