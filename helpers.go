package main

import (
	"net/http"
	"regexp"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/google/logger"
)

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

	filename := "/var/www/html/images/" + product.ID + "." + product.ImageExtension
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

func FetchUserByEmail(email string) (*User, error) {
	var user User

	row := db.QueryRow("select id, email, hash, salt, first_name, last_name, role from users where email = ?", email)
	err := row.Scan(&user.ID, &user.Email, &user.Hash, &user.Salt, &user.FirstName, &user.LastName, &user.Role)

	if err != nil {
		return nil, err
	}

	logger.Infof("User [%v] fetched", email)
	return &user, nil
}

func FetchUserByID(id string) (*User, error) {
	var user User

	row := db.QueryRow("select id, email, hash, salt, first_name, last_name, role, jti from users where id = ?", id)
	err := row.Scan(&user.ID, &user.Email, &user.Hash, &user.Salt, &user.FirstName, &user.LastName, &user.Role, &user.JTI)

	if err != nil {
		return nil, err
	}

	logger.Infof("User [%v] fetched", id)
	return &user, nil
}

func IsSignUpDataValid(data SignUp) bool {
	var firstName, lastName, email, password bool

	isEmailValid := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	var number, upper, lower bool
	letters := 0
	for _, c := range data.Password {
		switch {
		case unicode.IsNumber(c):
			number = true
			letters++
		case unicode.IsUpper(c):
			upper = true
			letters++
		case unicode.IsLower(c):
			lower = true
			letters++
		}
	}

	firstName = len(data.FirstName) > 0
	lastName = len(data.LastName) > 0
	email = isEmailValid.MatchString(data.Email)
	password = number && upper && lower && letters >= 6

	return firstName && lastName && email && password
}
