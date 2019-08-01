package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/logger"
)

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

func UpdateCatalog(c *gin.Context) {
	var catalog Catalog

	if err := c.ShouldBind(&catalog); err != nil {
		logger.Errorf("[UpdateCatalog] %v", err)
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": err.Error(),
			})
	}

	catalog.ID, _ = strconv.ParseInt(c.Param("id"), 10, 64)

	_, err := db.Exec(
		"update catalog SET name = ? where id = ?",
		catalog.Name,
		catalog.ID,
	)

	if err != nil {
		logger.Errorf("[DB Query : UpdateCatalog] %v; catalog = %v", err, catalog)
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": err.Error(),
			})
	} else {
		logger.Infof("Catalog updated to %v", catalog)
		c.JSON(
			http.StatusOK,
			gin.H{
				"status":  http.StatusOK,
				"message": "Updated successfully!",
			})
	}
}

func DeleteCatalog(c *gin.Context) {
	_, err := db.Exec("delete from catalog where id = ?", c.Param("id"))

	if err != nil {
		logger.Errorf("[DB Query : DeleteCatalog] %v; ID = %v", err, c.Param("id"))
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": err.Error(),
			})
	} else {
		logger.Infof("Catalog %v deleted", c.Param("id"))
		c.JSON(
			http.StatusOK,
			gin.H{
				"status":  http.StatusOK,
				"message": "Deleted successfully!",
			})
	}
}
