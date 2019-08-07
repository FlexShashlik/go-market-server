package main

type Login struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type SignUp struct {
	Email     string `form:"email" json:"email" binding:"required"`
	Password  string `form:"password" json:"password" binding:"required"`
	FirstName string `form:"first_name" json:"first_name" binding:"required"`
	LastName  string `form:"last_name" json:"last_name" binding:"required"`
}

type User struct {
	ID        string `form:"id" json:"id"`
	Email     string `form:"email" json:"email"`
	Hash      []byte `form:"hash" json:"hash"`
	Salt      string `form:"salt" json:"salt"`
	FirstName string `form:"first_name" json:"first_name"`
	LastName  string `form:"last_name" json:"last_name"`
	Role      string `form:"role" json:"role"`
	JTI       string `form:"jti" json:"jti"`
}

type Product struct {
	ID             string `form:"id" json:"id"`
	Name           string `form:"name" json:"name"`
	Price          string `form:"price" json:"price"`
	ImageExtension string `form:"image_extension" json:"image_extension"`
	SubCatalogID   string `form:"sub_catalog_id" json:"sub_catalog_id"`
}

type Catalog struct {
	ID   int64  `form:"id" json:"id"`
	Name string `form:"name" json:"name"`
}

type SubCatalog struct {
	ID        int64  `form:"id" json:"id"`
	Name      string `form:"name" json:"name"`
	CatalogID string `form:"catalog_id" json:"catalog_id"`
}

type Color struct {
	ID  int64 `json:"id"`
	RAL int64 `json:"ral"`
}

type Covering struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
