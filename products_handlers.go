package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/logger"
)

func FetchAllProducts(c *gin.Context) {
	var products []Product

	rows, err := db.Query("select id, name, price, image_extension, sub_catalog_id from products")

	if err != nil {
		logger.Errorf("[DB Query : FetchAllProducts] %v", err)
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": err.Error(),
			})
	} else {
		for rows.Next() {
			p := Product{}

			err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.ImageExtension, &p.SubCatalogID)

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

func FetchProductsBySubCatalog(c *gin.Context) {
	var products []Product

	subCatalogID := c.Param("sub_catalog_id")

	rows, err := db.Query("select id, name, price, image_extension, sub_catalog_id from products where sub_catalog_id = ?", subCatalogID)

	if err != nil {
		logger.Errorf("[DB Query : FetchProductsBySubCatalog] %v", err)
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": err.Error(),
			})
	} else {
		for rows.Next() {
			p := Product{}

			err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.ImageExtension, &p.SubCatalogID)

			if err != nil {
				logger.Errorf("[DB Query : FetchProductsBySubCatalog : rows.Scan] %v", err)
				continue
			}

			products = append(products, p)
		}

		logger.Infof("Products by subcatalog [%v] fetched", subCatalogID)
		c.JSON(http.StatusOK, products)
	}
}

func CreateProduct(c *gin.Context) {
	var product Product

	if err := c.ShouldBind(&product); err != nil {
		logger.Errorf("[CreateProduct] %v", err)
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": err.Error(),
			})
	}

	result, err := db.Exec(
		"insert into products (name, price, image_extension, sub_catalog_id) values (?, ?, ?, ?)",
		product.Name,
		product.Price,
		product.ImageExtension,
		product.SubCatalogID)

	if err != nil {
		logger.Errorf("[DB Query : CreateProduct] %v", err)
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": err.Error(),
			})
	} else {
		productID, err := result.LastInsertId()

		if err != nil {
			logger.Errorf("[DB Query : CreateProduct : LastInsertID] %v; ", err)
			c.JSON(
				http.StatusNotImplemented,
				gin.H{
					"status":  http.StatusNotImplemented,
					"message": err.Error(),
				})
		} else {
			product.ID = strconv.FormatInt(productID, 10)

			UploadImage(&product, c)

			logger.Infof("Product [%v] created", product)
			c.JSON(http.StatusCreated, gin.H{"ID": product.ID})
		}
	}
}

func UpdateProduct(c *gin.Context) {
	var product Product

	if err := c.ShouldBind(&product); err != nil {
		logger.Errorf("[UpdateProduct] %v", err)
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": err.Error(),
			})
	}

	product.ID = c.Param("id")

	_, err := db.Exec(
		"update products SET name = ?, price = ?, image_extension = ?, sub_catalog_id = ? where id = ?",
		product.Name,
		product.Price,
		product.ImageExtension,
		product.SubCatalogID,
		product.ID,
	)

	if err != nil {
		logger.Errorf("[DB Query : UpdateProduct] %v; product = %v", err, product)
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": err.Error(),
			})
	} else {
		UploadImage(&product, c)

		logger.Infof("Product updated to %v", product)
		c.JSON(
			http.StatusOK,
			gin.H{
				"status":  http.StatusOK,
				"message": "Product updated successfully!",
			})
	}
}

func DeleteProduct(c *gin.Context) {
	_, err := db.Exec("delete from products where id = ?", c.Param("id"))

	if err != nil {
		logger.Errorf("[DB Query : DeleteProduct] %v; productID = %v", err, c.Param("id"))
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": err.Error(),
			})
	} else {
		logger.Infof("Product %v deleted", c.Param("id"))
		c.JSON(
			http.StatusOK,
			gin.H{
				"status":  http.StatusOK,
				"message": "Product deleted successfully!",
			})
	}
}
