package main

import (
	"crypto/sha1"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/logger"
	"golang.org/x/crypto/pbkdf2"
)

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

func FetchAllUsers(c *gin.Context) {
	var users []User

	rows, err := db.Query("select id, email, first_name, last_name, role from users")

	if err != nil {
		logger.Errorf("[DB Query : FetchAllUsers] %v", err)
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": err.Error(),
			})
	} else {
		for rows.Next() {
			u := User{}

			err := rows.Scan(&u.ID, &u.Email, &u.FirstName, &u.LastName, &u.Role)

			if err != nil {
				logger.Errorf("[DB Query : FetchAllUsers : rows.Scan] %v", err)
				continue
			}

			users = append(users, u)
		}

		logger.Infof("Users fetched")
		c.JSON(http.StatusOK, users)
	}
}

func UpdateUser(c *gin.Context) {
	var user User

	if err := c.ShouldBind(&user); err != nil {
		logger.Errorf("[UpdateUser] %v", err)
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": err.Error(),
			})
	}

	user.ID = c.Param("id")

	_, err := db.Exec(
		"update users SET email = ?, first_name = ?, last_name = ?, role = ? where id = ?",
		user.Email,
		user.FirstName,
		user.LastName,
		user.Role,
		user.ID,
	)

	if err != nil {
		logger.Errorf("[DB Query : UpdateUser] %v; User = %v", err, user)
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": err.Error(),
			})
	} else {
		logger.Infof("User updated to %v", user)
		c.JSON(
			http.StatusOK,
			gin.H{
				"status":  http.StatusOK,
				"message": "Updated successfully!",
			})
	}
}

func DeleteUser(c *gin.Context) {
	_, err := db.Exec("delete from users where id = ?", c.Param("id"))

	if err != nil {
		logger.Errorf("[DB Query : DeleteUser] %v; ID = %v", err, c.Param("id"))
		c.JSON(
			http.StatusNotImplemented,
			gin.H{
				"status":  http.StatusNotImplemented,
				"message": err.Error(),
			})
	} else {
		logger.Infof("User %v deleted", c.Param("id"))
		c.JSON(
			http.StatusOK,
			gin.H{
				"status":  http.StatusOK,
				"message": "Deleted successfully!",
			})
	}
}
