package main

import (
	"net/http"
	"strconv"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/google/logger"
)

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

// HelloHandler demo
func HelloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	c.JSON(200, gin.H{
		"id": claims["id"],
		"text":   "Hello World.",
	})
}

// CreateUser registers new user
func CreateUser(c *gin.Context) {
	// Implement this shit
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
