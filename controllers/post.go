package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/daonham/go-app/database"
	"github.com/daonham/go-app/helper"
	"github.com/gin-gonic/gin"
)

type Post struct {
	Id        int       `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Published bool      `json:"published"`
	AuthorId  int       `json:"authorId"`
}

func GetPosts(c *gin.Context) {
	db := database.ConnectDataBase()

	rows, err := db.Query("SELECT id, createdAt, updatedAt, title, content, published, authorId FROM post")

	defer db.Close()

	if err != nil {
		resError := helper.ResponseError(http.StatusBadRequest, "Error when connect DB", err.Error(), helper.EmptyObj{})
		c.IndentedJSON(http.StatusBadRequest, resError)
		return
	}

	posts := []Post{}

	for rows.Next() {
		var id, authorId int
		var title, content string
		var createdAt, updatedAt time.Time
		var published bool

		err = rows.Scan(&id, &createdAt, &updatedAt, &title, &content, &published, &authorId)
		if err != nil {
			resError := helper.ResponseError(http.StatusBadRequest, "Error when connect DB", err.Error(), helper.EmptyObj{})
			c.IndentedJSON(http.StatusBadRequest, resError)
		}

		post := Post{}

		post.Id = id
		post.CreatedAt = createdAt
		post.UpdatedAt = updatedAt
		post.Title = title
		post.Content = content
		post.Published = published
		post.AuthorId = authorId

		posts = append(posts, post)
	}

	c.IndentedJSON(http.StatusOK, posts)
}

func GetPost(c *gin.Context) {
	cid := c.Param("id")
	pid, err := strconv.Atoi(cid)
	if err != nil {
		resError := helper.ResponseError(http.StatusBadRequest, "Please check Id param", err.Error(), helper.EmptyObj{})
		c.IndentedJSON(http.StatusBadRequest, resError)
		return
	}

	db := database.ConnectDataBase()

	rows, err := db.Query("SELECT id, createdAt, updatedAt, title, content, published, authorId FROM post WHERE id=?", pid)

	defer db.Close()

	if err != nil {
		resError := helper.ResponseError(http.StatusBadRequest, "Error when connect DB", err.Error(), helper.EmptyObj{})
		c.IndentedJSON(http.StatusBadRequest, resError)
		return
	}

	post := Post{}

	for rows.Next() {
		var id, authorId int
		var title, content string
		var createdAt, updatedAt time.Time
		var published bool

		err = rows.Scan(&id, &createdAt, &updatedAt, &title, &content, &published, &authorId)
		if err != nil {
			resError := helper.ResponseError(http.StatusBadRequest, "Error when connect DB", err.Error(), helper.EmptyObj{})
			c.IndentedJSON(http.StatusBadRequest, resError)
		}

		post.Id = id
		post.CreatedAt = createdAt
		post.UpdatedAt = updatedAt
		post.Title = title
		post.Content = content
		post.Published = published
		post.AuthorId = authorId
	}

	c.IndentedJSON(http.StatusOK, post)
}

func CreatePost(c *gin.Context) {

	type CreatePost struct {
		Title     string `form:"title" json:"title" binding:"required"`
		Content   string `form:"content" json:"content"`
		Published bool   `form:"published" json:"published"`
		AuthorId  int    `form:"authorId" json:"authorId" binding:"required"`
	}

	var create CreatePost

	if err := c.ShouldBindJSON(&create); err == nil {
		db := database.ConnectDataBase()

		insert, err := db.Prepare("INSERT INTO post(updatedAt, title, content, published, authorId) VALUES(?,?,?,?,?)")

		defer db.Close()

		if err != nil {
			resError := helper.ResponseError(http.StatusBadRequest, "Error when insert to database", err.Error(), helper.EmptyObj{})
			c.IndentedJSON(http.StatusBadRequest, resError)
			return
		}

		timeStamp := time.Now().UTC().Format("2006-01-02 15:04:05") // Format timestamp in mysql

		res, err := insert.Exec(timeStamp, create.Title, create.Content, create.Published, create.AuthorId)
		if err != nil {
			resError := helper.ResponseError(http.StatusBadRequest, "Error when insert to database", err.Error(), helper.EmptyObj{})
			c.IndentedJSON(http.StatusBadRequest, resError)
			return
		}

		lastId, err := res.LastInsertId()
		if err != nil {
			resError := helper.ResponseError(http.StatusBadRequest, "Cannot get LastInsertId", err.Error(), helper.EmptyObj{})
			c.IndentedJSON(http.StatusBadRequest, resError)
			return
		}

		data := gin.H{
			"id": lastId,
		}

		resSuccess := helper.ResponseSuccess(http.StatusCreated, "Post insert successfully", data)
		c.IndentedJSON(http.StatusCreated, resSuccess)
	} else {
		res := helper.ResponseError(http.StatusBadRequest, "Error: Should Bind Json", err.Error(), helper.EmptyObj{})
		c.IndentedJSON(http.StatusBadRequest, res)
	}
}

func UpdatePost(c *gin.Context) {
	cid := c.Param("id")
	id, err := strconv.Atoi(cid)
	if err != nil {
		resError := helper.ResponseError(http.StatusBadRequest, "Check Id param", err.Error(), helper.EmptyObj{})
		c.IndentedJSON(http.StatusBadRequest, resError)
		return
	}

	type UpdatePost struct {
		Title     string `form:"title" json:"title" binding:"required"`
		Content   string `form:"content" json:"content"`
		Published bool   `form:"published" json:"published"`
		AuthorId  int    `form:"authorId" json:"authorId" binding:"required"`
	}

	var update UpdatePost

	if err := c.ShouldBindJSON(&update); err == nil {
		db := database.ConnectDataBase()

		edit, err := db.Prepare("UPDATE post SET updatedAt=?, title=?, content=?, published=?, authorId=? WHERE id=?")

		defer db.Close()

		if err != nil {
			resError := helper.ResponseError(http.StatusBadRequest, "Update data failed", err.Error(), helper.EmptyObj{})
			c.IndentedJSON(http.StatusBadRequest, resError)
			return
		}

		timeStamp := time.Now().UTC().Format("2006-01-02 15:04:05") // Format timestamp in mysql

		res, err := edit.Exec(timeStamp, update.Title, update.Content, update.Published, update.AuthorId, id)
		if err != nil {
			resError := helper.ResponseError(http.StatusBadRequest, "Update data failed", err.Error(), helper.EmptyObj{})
			c.IndentedJSON(http.StatusBadRequest, resError)
			return
		}

		row, err := res.RowsAffected()
		if err != nil {
			resError := helper.ResponseError(http.StatusBadRequest, "Get rows affected error", err.Error(), helper.EmptyObj{})
			c.IndentedJSON(http.StatusBadRequest, resError)
			return
		}

		rows := int(row)

		data := gin.H{
			"rows": rows,
		}

		resSuccess := helper.ResponseSuccess(http.StatusOK, "Update post successfully", data)

		c.IndentedJSON(http.StatusOK, resSuccess)
	} else {
		resError := helper.ResponseError(http.StatusBadRequest, "Error", err.Error(), helper.EmptyObj{})
		c.IndentedJSON(http.StatusBadRequest, resError)
		return
	}
}

func DeletePost(c *gin.Context) {
	cid := c.Param("id")

	id, err := strconv.Atoi(cid)
	if err != nil {
		resError := helper.ResponseError(http.StatusBadRequest, "Check Id param", err.Error(), helper.EmptyObj{})
		c.IndentedJSON(http.StatusBadRequest, resError)
		return
	}

	db := database.ConnectDataBase()

	delete, err := db.Prepare("DELETE FROM post WHERE id=?")

	defer db.Close()

	if err != nil {
		resError := helper.ResponseError(http.StatusBadRequest, "Delete in database error", err.Error(), helper.EmptyObj{})
		c.IndentedJSON(http.StatusBadRequest, resError)
		return
	}

	res, err := delete.Exec(id)

	if err != nil {
		resError := helper.ResponseError(http.StatusBadRequest, "Delete in database error", err.Error(), helper.EmptyObj{})
		c.IndentedJSON(http.StatusBadRequest, resError)
		return
	}

	row, err := res.RowsAffected() // Get number rows deleted
	if err != nil {
		resError := helper.ResponseError(http.StatusBadRequest, "Error get row affected", err.Error(), helper.EmptyObj{})
		c.IndentedJSON(http.StatusBadRequest, resError)
		return
	}

	rows := int(row)

	data := gin.H{
		"rows": rows,
	}

	resSuccess := helper.ResponseSuccess(http.StatusOK, "Delete post successfully", data)

	c.IndentedJSON(http.StatusOK, resSuccess)
}
