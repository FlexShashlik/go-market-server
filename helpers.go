package main

import (
	"net/http"
	"strconv"

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

// FetchUser fetches user info
func FetchUser(email string) (*User, error) {
	var user User

	row := db.QueryRow("select id, email, hash, salt, first_name, last_name, role from users where email = ?", email)
	err := row.Scan(&user.ID, &user.Email, &user.Hash, &user.Salt, &user.FirstName, &user.LastName, &user.Role)

	if err != nil {
		logger.Errorf("[DB Query : FetchUser : row.Scan] %v; email = %v", err, email)
		return nil, err
	}

	logger.Infof("User %v fetched", email)
	return &user, nil
}
