package main

type Login struct {
	Email     string `form:"email" json:"email" binding:"required"`
	Password  string `form:"password" json:"password" binding:"required"`
}

// User demo struct
type User struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
	JTI       string `json:"jti"`
}

// Product demo struct
type Product struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Price          string `json:"price"`
	ImageExtension string `json:"image_extension"`
}
