package main

import (
	"net/http"
	"regexp"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/google/logger"
)

var pricePerFold int64
var brigadierCost int64
var standardCoveringPrice int64

func FetchConsts() {
	row := db.QueryRow(
		"select value from constants where name = ?",
		"price_per_fold",
	)

	err := row.Scan(&pricePerFold)

	if err != nil {
		logger.Errorf("[DB Query : FetchConsts : PricePerFold] %v", err)
	}

	logger.Infof("PricePerFold [%v] fetched", pricePerFold)

	row = db.QueryRow("select value from constants where name = ?", "brigadier_cost")
	err = row.Scan(&brigadierCost)

	if err != nil {
		logger.Errorf("[DB Query : FetchConsts : BrigadierCost] %v", err)
	}

	logger.Infof("BrigadierCost [%v] fetched", brigadierCost)
}

func FetchFirstCombo() {
	row := db.QueryRow("select price from combinations")
	err := row.Scan(&standardCoveringPrice)

	if err != nil {
		logger.Errorf("[DB Query : FetchFirstCombo] %v", err)
	}

	logger.Infof("StandardCoveringPrice [%v] fetched", standardCoveringPrice)
}

func FetchCustomCombo(colorID int64, coveringID int64) int64 {
	var customCoveringPrice int64

	row := db.QueryRow(
		"select price from combinations where color_id = ? and covering_id = ?",
		colorID,
		coveringID,
	)

	err := row.Scan(&customCoveringPrice)

	if err != nil {
		logger.Errorf("[DB Query : FetchCustomCombo] %v", err)
	}

	logger.Infof("CustomCoveringPrice [%v] fetched", customCoveringPrice)

	return customCoveringPrice
}

func EvaluateProduct(product *Product, colorID int64, coveringID int64) {
	if (pricePerFold | brigadierCost) == 0 {
		FetchConsts()
	}

	if (colorID | coveringID) == 0 {
		if standardCoveringPrice == 0 {
			FetchFirstCombo()
		}

		product.Price = (standardCoveringPrice / product.ProductPerList) + (product.Folds * pricePerFold) + brigadierCost
	} else {
		coveringPrice := FetchCustomCombo(colorID, coveringID)

		product.Price = (coveringPrice / product.ProductPerList) + (product.Folds * pricePerFold) + brigadierCost
	}
}

func UploadImage(product *Product, c *gin.Context) {
	file, err := c.FormFile("image")

	if err != nil {
		logger.Errorf("[DB Query : UploadImage : FormFile()] %v", err)
		return
	}

	filename := "C:/xampp/htdocs/images/" + product.ID + "." + product.ImageExtension
	if err := c.SaveUploadedFile(file, filename); err != nil {
		logger.Errorf("[DB Query : UploadImage : SaveUploadedFile()] %v", err)
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
