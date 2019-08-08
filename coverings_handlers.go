package main

import (
	"net/http"
	"github.com/google/logger"
	"github.com/gin-gonic/gin"
)

func FetchColors(c *gin.Context) {
	var colors []Color

	rows, err := db.Query("select id, ral from colors")

	if err != nil {
		logger.Errorf("[DB Query : FetchColors] %v", err)
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": err.Error(),
			})
	} else {
		for rows.Next() {
			c := Color{}

			err := rows.Scan(&c.ID, &c.RAL)

			if err != nil {
				logger.Errorf("[DB Query : FetchColors : rows.Scan] %v", err)
				continue
			}

			colors = append(colors, c)
		}

		logger.Infof("Colors fetched")
		c.JSON(http.StatusOK, colors)
	}
}

func FetchCoverings(c *gin.Context) {
	var coverings []Covering

	rows, err := db.Query("select id, name from coverings")

	if err != nil {
		logger.Errorf("[DB Query : FetchCoverings] %v", err)
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": err.Error(),
			})
	} else {
		for rows.Next() {
			c := Covering{}

			err := rows.Scan(&c.ID, &c.Name)

			if err != nil {
				logger.Errorf("[DB Query : FetchCoverings : rows.Scan] %v", err)
				continue
			}

			coverings = append(coverings, c)
		}

		logger.Infof("Coverings fetched")
		c.JSON(http.StatusOK, coverings)
	}
}

func FetchColorsByCovering(c *gin.Context) {
	var colors []Color

	rows, err := db.Query(
		"select colors.id, colors.ral from combinations inner join colors on combinations.color_id = colors.id where covering_id = ?", 
		c.Param("id"),
	)

	if err != nil {
		logger.Errorf("[DB Query : FetchColorsByCovering] %v", err)
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": err.Error(),
			})
	} else {
		for rows.Next() {
			c := Color{}

			err := rows.Scan(&c.ID, &c.RAL)

			if err != nil {
				logger.Errorf("[DB Query : FetchColorsByCovering : rows.Scan] %v", err)
				continue
			}

			colors = append(colors, c)
		}

		logger.Infof("Colors fetched")
		c.JSON(http.StatusOK, colors)
	}
}

func FetchCoveringsByColor(c *gin.Context) {
	var coverings []Covering

	rows, err := db.Query(
		"select coverings.id, coverings.name from combinations inner join coverings on combinations.covering_id = coverings.id where color_id = ?", 
		c.Param("id"),
	)

	if err != nil {
		logger.Errorf("[DB Query : FetchCoveringsByColor] %v", err)
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": err.Error(),
			})
	} else {
		for rows.Next() {
			c := Covering{}

			err := rows.Scan(&c.ID, &c.Name)

			if err != nil {
				logger.Errorf("[DB Query : FetchCoveringsByColor : rows.Scan] %v", err)
				continue
			}

			coverings = append(coverings, c)
		}

		logger.Infof("Coverings fetched")
		c.JSON(http.StatusOK, coverings)
	}
}