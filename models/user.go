package models

import (
	"errors"
	"time"

	"github.com/daonham/go-app/database"
	"github.com/daonham/go-app/forms"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        int       `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	Role      string    `json:"role"`
}

func GetUsers() (users []User, message string, err error) {
	db := database.ConnectDB()

	rows, err := db.Query("SELECT id, email, name, createdAt, role FROM user")

	defer db.Close()

	if err != nil {
		return users, "Something went wrong, please try again later", err
	}

	for rows.Next() {
		var id int
		var name, email, role string
		var createdAt time.Time

		err = rows.Scan(&id, &email, &name, &createdAt, &role)
		if err != nil {
			return users, "Something went wrong, please try again later", err
		}

		user := User{}

		user.Id = id
		user.Email = email
		user.Name = name
		user.CreatedAt = createdAt
		user.Role = role

		users = append(users, user)
	}

	return users, "Select list users successfully", nil
}

func GetUser(uid int) (user User, message string, err error) {
	db := database.ConnectDB()

	row := db.QueryRow("SELECT id, email, name, createdAt, role FROM user WHERE id=?", uid)

	defer db.Close()

	var id int
	var email, name, role string
	var createdAt time.Time

	err = row.Scan(&id, &email, &name, &createdAt, &role)
	if err != nil {
		return user, "No rows in result", err
	}

	user.Id = id
	user.Email = email
	user.Name = name
	user.CreatedAt = createdAt
	user.Role = role

	return user, "Get user successfully", nil
}

func CreateUser(create forms.UserForm) (id int64, message string, err error) {

	db := database.ConnectDB()

	check := db.QueryRow("SELECT id FROM user WHERE email=?", create.Email)

	err = check.Scan(&id)

	if err == nil {
		return 0, "Email already exists", errors.New("email already exists")
	}

	bytePassword := []byte(create.Pass)
	hashedPassword, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)

	insert, err := db.Prepare("INSERT INTO user(email, name, pass, role) VALUES(?,?,?,?)")

	defer db.Close()

	if err != nil {
		return 0, "Something went wrong, please try again later", err
	}

	res, err := insert.Exec(create.Email, create.Name, string(hashedPassword), create.Role)
	if err != nil {
		return id, "Something went wrong, please try again later", err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, "Error get LastInsertId", err
	}

	return lastId, "Create user successfully", nil
}

func UpdateUser(uid int, update forms.UserForm) (rows int64, message string, err error) {

	db := database.ConnectDB()

	edit, err := db.Prepare("UPDATE user SET email=?, name=?, pass=?, role=? WHERE id=?")

	defer db.Close()

	if err != nil {
		return 0, "Something went wrong, please try again later", err
	}

	bytePassword := []byte(update.Pass)
	hashedPassword, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)

	res, err := edit.Exec(update.Email, update.Name, string(hashedPassword), update.Role, uid)

	if err != nil {
		return 0, "Please check your email address", err
	}

	row, err := res.RowsAffected() // Get records rows deleted

	if err != nil {
		return 0, "Error: Get rows", err
	}

	if row == 0 {
		return 0, "Error: Update 0 records", errors.New("updated 0 records")
	}

	return row, "Update post successfully", nil
}

func DeleteUser(id int) (rows int64, message string, err error) {
	db := database.ConnectDB()

	delete, err := db.Prepare("DELETE FROM user WHERE id=?")

	defer db.Close()

	if err != nil {
		return 0, "Error when connect DB", err
	}

	res, err := delete.Exec(id)

	if err != nil {
		return 0, "Error when connect DB", err
	}

	row, err := res.RowsAffected() // Get records rows deleted

	if err != nil {
		return 0, "Error get row affected", err
	}

	if row == 0 {
		return 0, "Error: Cannot delete this post", errors.New("deleted 0 records")
	}

	return row, "Delete post successfully", nil
}

func Login(form forms.LoginForm) (user *User, token Token, message string, err error) {
	db := database.ConnectDB()

	query := db.QueryRow("SELECT id, email, name, pass, createdAt, role FROM user WHERE email=?", form.Email)

	defer db.Close()

	var id int
	var name, email, pass, role string
	var createdAt time.Time

	err = query.Scan(&id, &email, &name, &pass, &createdAt, &role)

	if err != nil {
		return nil, token, "Something went wrong, please try again later", err
	}

	user = &User{
		Id:        id,
		Email:     email,
		Name:      name,
		CreatedAt: createdAt,
		Role:      role,
	}

	//Compare the password form and database if match
	bytePassword := []byte(form.Pass)
	byteHashedPassword := []byte(pass)

	err = bcrypt.CompareHashAndPassword(byteHashedPassword, bytePassword)

	if err != nil {
		return nil, token, "Please check your password", err
	}

	//Generate the JWT auth token
	tokenDetails, err := CreateToken(user.Id)

	if err != nil {
		return nil, token, "Cannot create token", err
	}

	token.AccessToken = tokenDetails.AccessToken
	token.RefreshToken = tokenDetails.RefreshToken

	return user, token, "Login Done", nil
}

func Register(form forms.RegisterForm) (user *User, message string, err error) {
	db := database.ConnectDB()

	check := db.QueryRow("SELECT id FROM user WHERE email=?", form.Email)

	var id int

	err = check.Scan(&id)

	if err == nil {
		return nil, "Email already exists", errors.New("email already exists")
	}

	bytePassword := []byte(form.Pass)
	hashedPassword, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)

	if err != nil {
		return nil, "Something went wrong, please try again later", err
	}

	insert, err := db.Prepare("INSERT INTO user(email, name, pass) VALUES(?,?,?)")

	defer db.Close()

	if err != nil {
		return nil, "Something went wrong, please try again later", err
	}

	res, err := insert.Exec(form.Email, form.Name, string(hashedPassword))
	if err != nil {
		return nil, "Something went wrong, please try again later", err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return nil, "Error get LastInsertId", err
	}

	user = &User{
		Id:    int(lastId),
		Email: form.Email,
		Name:  form.Name,
		Role:  "USER",
	}

	return user, "Register Done", err
}
