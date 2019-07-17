package main

// Login demo struct
type Login struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// SignUp demo struct
type SignUp struct {
	Email     string `form:"email" json:"email" binding:"required"`
	Password  string `form:"password" json:"password" binding:"required"`
	FirstName string `form:"first_name" json:"first_name" binding:"required"`
	LastName  string `form:"last_name" json:"last_name" binding:"required"`
}

// User demo struct
type User struct {
	ID        int64  `json:"id"`
	Email     string `json:"email"`
	Hash      []byte `json:"hash"`
	Salt      string `json:"salt"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Role      string `json:"role"`
	JTI       string `json:"jti"`
}

// Product demo struct
type Product struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Price          int64  `json:"price"`
	ImageExtension string `json:"image_extension"`
}
