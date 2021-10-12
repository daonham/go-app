package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/daonham/go-app/forms"
	"github.com/daonham/go-app/helper"
	"github.com/daonham/go-app/models"
	"github.com/gin-gonic/gin"
)

type User struct {
	Id        int       `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	Role      string    `json:"role"`
}

func GetUsers(c *gin.Context) {
	users, message, err := models.GetUsers()

	if err != nil {
		c.IndentedJSON(http.StatusNotAcceptable, helper.ResponseError(http.StatusNotAcceptable, message, err.Error(), helper.EmptyObj{}))
		return
	}

	c.IndentedJSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	cid := c.Param("id")
	uid, err := strconv.Atoi(cid)

	if err != nil {
		c.IndentedJSON(http.StatusNotAcceptable, helper.ResponseError(http.StatusNotAcceptable, "Id param need type int", err.Error(), helper.EmptyObj{}))
		return
	}

	user, message, err := models.GetUser(uid)

	if err != nil {
		c.IndentedJSON(http.StatusNotAcceptable, helper.ResponseError(http.StatusNotAcceptable, message, err.Error(), helper.EmptyObj{}))
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context) {

	var create forms.UserForm

	if err := c.ShouldBindJSON(&create); err == nil {
		id, message, err := models.CreateUser(create)

		if err != nil {
			c.IndentedJSON(http.StatusNotAcceptable, helper.ResponseError(http.StatusNotAcceptable, message, err.Error(), helper.EmptyObj{}))
			return
		}

		data := gin.H{
			"id": id,
		}

		c.IndentedJSON(http.StatusCreated, helper.ResponseSuccess(http.StatusCreated, message, data))
	} else {
		c.IndentedJSON(http.StatusNotAcceptable, helper.ResponseError(http.StatusNotAcceptable, "Error: Should Bind Json", err.Error(), helper.EmptyObj{}))
	}
}

func UpdateUser(c *gin.Context) {
	cid := c.Param("id")
	id, err := strconv.Atoi(cid)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, helper.ResponseError(http.StatusBadRequest, "Check Id param", err.Error(), helper.EmptyObj{}))
		return
	}

	var update forms.UserForm

	if err := c.ShouldBindJSON(&update); err == nil {
		rows, message, err := models.UpdateUser(id, update)

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, helper.ResponseError(http.StatusBadRequest, message, err.Error(), helper.EmptyObj{}))
			return
		}

		data := gin.H{
			"rows": rows,
		}

		c.IndentedJSON(http.StatusOK, helper.ResponseSuccess(http.StatusOK, message, data))
	} else {
		c.IndentedJSON(http.StatusBadRequest, helper.ResponseError(http.StatusBadRequest, "Error: Should Bind Json", err.Error(), helper.EmptyObj{}))
		return
	}
}

func DeleteUser(c *gin.Context) {
	cid := c.Param("id")

	id, err := strconv.Atoi(cid)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, helper.ResponseError(http.StatusBadRequest, "Check Id param", err.Error(), helper.EmptyObj{}))
		return
	}

	rows, message, err := models.DeleteUser(id)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, helper.ResponseError(http.StatusBadRequest, message, err.Error(), helper.EmptyObj{}))
		return
	}

	data := gin.H{
		"rows": rows,
	}

	c.IndentedJSON(http.StatusOK, helper.ResponseSuccess(http.StatusOK, message, data))
}

func Login(c *gin.Context) {
	var login forms.LoginForm

	if err := c.ShouldBindJSON(&login); err == nil {
		user, token, message, err := models.Login(login)

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, helper.ResponseError(http.StatusBadRequest, message, err.Error(), helper.EmptyObj{}))
			return
		}

		data := gin.H{
			"user":  user,
			"token": token,
		}

		c.IndentedJSON(http.StatusOK, helper.ResponseSuccess(http.StatusOK, message, data))
	} else {
		c.IndentedJSON(http.StatusBadRequest, helper.ResponseError(http.StatusBadRequest, "Error: Should Bind Json", err.Error(), helper.EmptyObj{}))
		return
	}
}

func Register(c *gin.Context) {
	var register forms.RegisterForm

	if err := c.ShouldBindJSON(&register); err == nil {
		user, message, err := models.Register(register)

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, helper.ResponseError(http.StatusBadRequest, message, err.Error(), helper.EmptyObj{}))
			return
		}

		data := gin.H{
			"user": user,
		}

		c.IndentedJSON(http.StatusOK, helper.ResponseSuccess(http.StatusOK, message, data))
	} else {
		c.IndentedJSON(http.StatusBadRequest, helper.ResponseError(http.StatusBadRequest, "Error: Should Bind Json", err.Error(), helper.EmptyObj{}))
		return
	}
}
