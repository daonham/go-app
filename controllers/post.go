package controllers

import (
	"net/http"
	"strconv"

	"github.com/daonham/go-app/forms"
	"github.com/daonham/go-app/helper"
	"github.com/daonham/go-app/models"
	"github.com/gin-gonic/gin"
)

func GetPosts(c *gin.Context) {
	posts, message, err := models.GetPosts()

	if err != nil {
		c.IndentedJSON(http.StatusNotAcceptable, helper.ResponseError(http.StatusNotAcceptable, message, err.Error(), helper.EmptyObj{}))
		return
	}

	c.IndentedJSON(http.StatusOK, posts)
}

func GetPost(c *gin.Context) {
	cid := c.Param("id")
	pid, err := strconv.Atoi(cid)

	if err != nil {
		c.IndentedJSON(http.StatusNotAcceptable, helper.ResponseError(http.StatusNotAcceptable, "Id param need type int", err.Error(), helper.EmptyObj{}))
		return
	}

	post, message, err := models.GetPost(pid)

	if err != nil {
		c.IndentedJSON(http.StatusNotAcceptable, helper.ResponseError(http.StatusNotAcceptable, message, err.Error(), helper.EmptyObj{}))
		return
	}

	c.IndentedJSON(http.StatusOK, post)
}

func CreatePost(c *gin.Context) {
	var create forms.PostForm

	if err := c.ShouldBindJSON(&create); err == nil {
		id, message, err := models.CreatePost(create)

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

func UpdatePost(c *gin.Context) {
	cid := c.Param("id")
	id, err := strconv.Atoi(cid)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, helper.ResponseError(http.StatusBadRequest, "Check Id param", err.Error(), helper.EmptyObj{}))
		return
	}

	var update forms.PostForm

	if err := c.ShouldBindJSON(&update); err == nil {
		rows, message, err := models.UpdatePost(id, update)

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

func DeletePost(c *gin.Context) {
	cid := c.Param("id")

	id, err := strconv.Atoi(cid)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, helper.ResponseError(http.StatusBadRequest, "Check Id param", err.Error(), helper.EmptyObj{}))
		return
	}

	rows, message, err := models.DeletePost(id)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, helper.ResponseError(http.StatusBadRequest, message, err.Error(), helper.EmptyObj{}))
		return
	}

	data := gin.H{
		"rows": rows,
	}

	c.IndentedJSON(http.StatusOK, helper.ResponseSuccess(http.StatusOK, message, data))
}
