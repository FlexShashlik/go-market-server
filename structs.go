package main

// Login struct
type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// User demo struct
type User struct {
	UserName  string
	FirstName string
	LastName  string
	JTI       string
}

// Product demo struct
type Product struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Price          string `json:"price"`
	ImageExtension string `json:"image_extension"`
}
