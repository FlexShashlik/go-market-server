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