package main

import "testing"

func TestSignUp1(t *testing.T) {
	result := IsSignUpDataValid(
		SignUp{
			FirstName: "",
			LastName:  "",
			Email:     "",
			Password:  "",
		})
	if result {
		t.Error("Expected false, got ", result)
	}
}

func TestSignUp2(t *testing.T) {
	result := IsSignUpDataValid(
		SignUp{
			FirstName: "Test",
			LastName:  "Testyan",
			Email:     "tests@test.com",
			Password:  "lol11",
		})
	if result {
		t.Error("Expected false, got ", result)
	}
}

func TestSignUp3(t *testing.T) {
	result := IsSignUpDataValid(
		SignUp{
			FirstName: "User",
			LastName:  "1422",
			Email:     "test@lol.com",
			Password:  "laal13",
		})
	if result {
		t.Error("Expected false, got ", result)
	}
}

func TestSignUp4(t *testing.T) {
	result := IsSignUpDataValid(
		SignUp{
			FirstName: "Test21",
			LastName:  "Useryan",
			Email:     "adm@in@lulz.com",
			Password:  "194DudJe(_",
		})
	if result {
		t.Error("Expected false, got ", result)
	}
}

func TestSignUp5(t *testing.T) {
	result := IsSignUpDataValid(
		SignUp{
			FirstName: "User",
			LastName:  "Useryan",
			Email:     "admin@admin.com",
			Password:  "notthathard",
		})
	if result {
		t.Error("Expected false, got ", result)
	}
}

func TestSignUp6(t *testing.T) {
	result := IsSignUpDataValid(
		SignUp{
			FirstName: "User",
			LastName:  "Useryan",
			Email:     "admin@admin.com",
			Password:  "HardPassword_1",
		})
	if !result {
		t.Error("Expected true, got ", result)
	}
}
