package main

import (
	"database/sql"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/logger"
	"github.com/appleboy/gin-jwt"
)

// Product struct
type Product struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Price          string `json:"price"`
	ImageExtension string `json:"image_extension"`
}

const logPath = "server.log"

var verbose = flag.Bool("verbose", false, "print info level logs to stdout")
var db *sql.DB

func main() {
	flag.Parse()

	lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
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
