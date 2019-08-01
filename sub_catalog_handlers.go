package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/logger"
)

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

func CreateSubCatalog(c *gin.Context) {
	var subCatalog SubCatalog

	if err := c.ShouldBind(&subCatalog); err != nil {
		logger.Errorf("[CreateSubCatalog] %v", err)
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": err.Error(),
			})
	}

	result, err := db.Exec("insert into sub_catalog (name, catalog_id) values (?, ?)",
		subCatalog.Name,
		subCatalog.CatalogID)

	if err != nil {
		logger.Errorf("[DB Query : CreateSubCatalog] %v", err)
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": err.Error(),
			})
	} else {
		subCatalog.ID, err = result.LastInsertId()

		if err != nil {
			logger.Errorf("[DB Query : CreateSubCatalog : LastInsertID] %v; ", err)
			c.JSON(
				http.StatusNotImplemented,
				gin.H{
					"status":  http.StatusNotImplemented,
					"message": err.Error(),
				})
		} else {
			logger.Infof("SubCatalog [%v] created", subCatalog)
			c.JSON(http.StatusCreated, gin.H{"ID": subCatalog.ID})
		}
	}
}

func UpdateSubCatalog(c *gin.Context) {
	var subCatalog SubCatalog

	if err := c.ShouldBind(&subCatalog); err != nil {
		logger.Errorf("[UpdateSubCatalog] %v", err)
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": err.Error(),
			})
	}

	subCatalog.ID, _ = strconv.ParseInt(c.Param("id"), 10, 64)

	_, err := db.Exec(
		"update sub_catalog SET name = ?, catalog_id = ? where id = ?",
		subCatalog.Name,
		subCatalog.CatalogID,
		subCatalog.ID,
	)

	if err != nil {
		logger.Errorf("[DB Query : UpdateSubCatalog] %v; SubCatalog = %v", err, subCatalog)
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": err.Error(),
			})
	} else {
		logger.Infof("SubCatalog updated to %v", subCatalog)
		c.JSON(
			http.StatusOK,
			gin.H{
				"status":  http.StatusOK,
				"message": "Updated successfully!",
			})
	}
}

func DeleteSubCatalog(c *gin.Context) {
	_, err := db.Exec("delete from sub_catalog where id = ?", c.Param("id"))

	if err != nil {
		logger.Errorf("[DB Query : DeleteSubCatalog] %v; ID = %v", err, c.Param("id"))
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": err.Error(),
			})
	} else {
		logger.Infof("SubCatalog %v deleted", c.Param("id"))
		c.JSON(
			http.StatusOK,
			gin.H{
				"status":  http.StatusOK,
				"message": "Deleted successfully!",
			})
	}
}
