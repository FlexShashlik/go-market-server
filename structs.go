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
	ID        string `json:"id"`
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
	ID             string `form:"id" json:"id"`
	Name           string `form:"name" json:"name"`
	Price          string `form:"price" json:"price"`
	ImageExtension string `form:"image_extension" json:"image_extension"`
	SubCatalogID   string `form:"sub_catalog_id" json:"sub_catalog_id"`
}

// Catalog demo struct
type Catalog struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// SubCatalog demo struct
type SubCatalog struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CatalogID string `json:"catalog_id"`
}
