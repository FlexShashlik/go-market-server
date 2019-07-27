package main

import (
	"strconv"
	"crypto/sha1"
	"net/http"

	"golang.org/x/crypto/pbkdf2"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	"github.com/google/logger"
)

// HelloHandler demo
func HelloHandler(c *gin.Context) {
	claims := jwt.ExtractClaims(c)
	c.JSON(200, gin.H{
		"id":   claims["id"],
		"jti":  claims["jti"],
		"text": "Hello World.",
	})
}

// FetchCatalog fetches catalog
func FetchCatalog(c *gin.Context) {
	var catalog []Catalog

	rows, err := db.Query("select id, name from catalog")
	if err != nil {
		logger.Errorf("[DB Query : FetchCatalog] %v", err)
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": err.Error(),
			})
	} else {
		for rows.Next() {
			c := Catalog{}

			err := rows.Scan(&c.ID, &c.Name)

			if err != nil {
				logger.Errorf("[DB Query : FetchCatalog : rows.Scan] %v", err)
				continue
			}
			catalog = append(catalog, c)
		}
		logger.Infof("Catalog fetched")
		c.JSON(http.StatusOK, catalog)
	}
}

// FetchSubCatalog fetches sub_catalog
func FetchSubCatalog(c *gin.Context) {
	var subCatalog []SubCatalog

	rows, err := db.Query("select id, name, catalog_id from sub_catalog")
	if err != nil {
		logger.Errorf("[DB Query : FetchSubCatalog] %v", err)
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": err.Error(),
			})
	} else {
		for rows.Next() {
			sc := SubCatalog{}

			err := rows.Scan(&sc.ID, &sc.Name, &sc.CatalogID)

			if err != nil {
				logger.Errorf("[DB Query : FetchSubCatalog : rows.Scan] %v", err)
				continue
			}
			subCatalog = append(subCatalog, sc)
		}
		logger.Infof("SubCatalog fetched")
		c.JSON(http.StatusOK, subCatalog)
	}
}

// CreateUser registers new user
func CreateUser(c *gin.Context) {
	logger.Info("[CreateUser] attempt to create new user")

	var signUp SignUp
	var user User
	if err := c.ShouldBind(&signUp); err != nil {
		logger.Errorf("[CreateUser] %v", err)
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": err.Error(),
			})
	}

	if IsSignUpDataValid(signUp) {
		user.Email = signUp.Email

		// Check if user with this email already exists
		_, err := FetchUserByEmail(user.Email)
		if err == nil {
			// User already exists
			logger.Error("[DB Query : CreateUser] User alredy exists!")
			c.JSON(
				http.StatusNotImplemented,
				gin.H{
					"status":  http.StatusNotImplemented,
					"message": "User with this email already exists!",
				})
		} else {
			user.FirstName = signUp.FirstName
			user.LastName = signUp.LastName

			user.Salt, _ = GenerateRandomString(20)
			user.Hash = pbkdf2.Key([]byte(signUp.Password), []byte(user.Salt), 4096, 256, sha1.New)

			result, err := db.Exec(
				"insert into users (email, hash, salt, first_name, last_name) values (?, ?, ?, ?, ?)",
				user.Email,
				user.Hash,
				user.Salt,
				user.FirstName,
				user.LastName,
			)

			if err != nil {
				logger.Errorf("[DB Query : CreateUser] %v", err)
				c.JSON(
					http.StatusNotImplemented,
					gin.H{
						"status":  http.StatusNotImplemented,
						"message": err.Error(),
					})
			} else {
				userID, err := result.LastInsertId()
				if err != nil {
					logger.Errorf("[DB Query : CreateUser : LastInsertID] %v; ", err)
					c.JSON(
						http.StatusNotImplemented,
						gin.H{
							"status":  http.StatusNotImplemented,
							"message": err.Error(),
						})
				} else {
					logger.Infof("User [%v] created", user)
					c.JSON(http.StatusCreated, gin.H{"ID": userID})
				}
			}
		}
	} else {
		logger.Error("[CreateUser] sign up data is invalid!")
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": "Sign up data is invalid!",
			})
	}
}

// FetchAllProducts fetches all products
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

// FetchProductsBySubCatalog fetches products by subcatalog
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

// CreateProduct creates new product
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

// UpdateProduct updates product
func UpdateProduct(c *gin.Context) {
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

	product.ID = c.Param("id")

	_, err := db.Exec(
		"update products SET name = ?, price = ?, image_extension = ? where id = ?",
		product.Name,
		product.Price,
		product.ImageExtension,
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

		logger.Infof("Product %v updated to name = %v, price = %v, image_extension = %v",
			product.ID,
			product.Name,
			product.Price,
			product.ImageExtension,
		)
		c.JSON(
			http.StatusOK,
			gin.H{
				"status":  http.StatusOK,
				"message": "Product updated successfully!",
			})
	}
}

// DeleteProduct deletes product
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

// CreateCatalog stupid linter
func CreateCatalog(c *gin.Context) {
	var catalog Catalog

	if err := c.ShouldBind(&catalog); err != nil {
		logger.Errorf("[CreateCatalog] %v", err)
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": err.Error(),
			})
	}

	result, err := db.Exec("insert into catalog (name) values (?)",
							catalog.Name)

	if err != nil {
		logger.Errorf("[DB Query : CreateCatalog] %v", err)
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": err.Error(),
			})
	} else {
		catalog.ID, err = result.LastInsertId()
		if err != nil {
			logger.Errorf("[DB Query : CreateCatalog : LastInsertID] %v; ", err)
			c.JSON(
				http.StatusNotImplemented,
				gin.H{
					"status":  http.StatusNotImplemented,
					"message": err.Error(),
				})
		} else {
			logger.Infof("Catalog [%v] created", catalog)
			c.JSON(http.StatusCreated, gin.H{"ID": catalog.ID})
		}
	}
}
