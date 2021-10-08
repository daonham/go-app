package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/daonham/go-app/database"
	"github.com/gin-gonic/gin"
)

type Post struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func GetPosts(c *gin.Context) {
	db := database.ConnectDataBase()

	rows, err := db.Query("SELECT id, title, content FROM posts")

	defer db.Close()

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"messages": "Post not found",
		})
	}

	posts := []Post{}

	for rows.Next() {
		var id int
		var title, content string

		err = rows.Scan(&id, &title, &content)
		if err != nil {
			panic(err.Error())
		}

		post := Post{}

		post.Id = id
		post.Title = title
		post.Content = content

		posts = append(posts, post)
	}

	c.IndentedJSON(http.StatusOK, posts)
}

func GetPost(c *gin.Context) {
	cid := c.Param("id")
	pid, err := strconv.Atoi(cid)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"messages": err,
		})
	}

	db := database.ConnectDataBase()

	rows, err := db.Query("SELECT id, title, content FROM posts WHERE id=?", pid)

	defer db.Close()

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{
			"messages": "Post not found",
		})
	}

	post := Post{}

	for rows.Next() {
		var id int
		var title, content string

		err = rows.Scan(&id, &title, &content)
		if err != nil {
			panic(err.Error())
		}

		post.Id = id
		post.Title = title
		post.Content = content
	}

	c.IndentedJSON(http.StatusOK, post)
}

func CreatePost(c *gin.Context) {

	type CreatePost struct {
		Title   string `form:"title" json:"title" binding:"required"`
		Content string `form:"content" json:"content" binding:"required"`
	}

	var create CreatePost

	if err := c.ShouldBindJSON(&create); err == nil {
		db := database.ConnectDataBase()

		insert, err := db.Prepare("INSERT INTO posts(title, content) VALUES(?,?)")

		defer db.Close()

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"messages": err,
			})
		}

		res, err := insert.Exec(create.Title, create.Content)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"messages": err,
			})
		}

		lastId, err := res.LastInsertId()
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"messages": err,
			})
		}

		id := int(lastId)

		c.IndentedJSON(http.StatusCreated, gin.H{
			"messages": "Post inserted",
			"id":       id,
		})
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
}

func UpdatePost(c *gin.Context) {
	cid := c.Param("id")
	id, err := strconv.Atoi(cid)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"messages": err,
		})
	}

	type UpdatePost struct {
		Title   string `form:"title" json:"title" binding:"required"`
		Content string `form:"content" json:"content" binding:"required"`
	}

	var update UpdatePost

	if err := c.ShouldBindJSON(&update); err == nil {
		db := database.ConnectDataBase()

		edit, err := db.Prepare("UPDATE posts SET title=?, content=? WHERE id=?")

		defer db.Close()

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"messages": err,
			})
		}

		res, err := edit.Exec(update.Title, update.Content, id)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"messages": err,
			})
		}

		row, err := res.RowsAffected()
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"messages": err,
			})
		}

		rows := int(row)

		c.IndentedJSON(http.StatusOK, gin.H{
			"messages": fmt.Sprintf("Update post %d successful %d", id, rows),
			"rows":     rows,
		})
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
}

func DeletePost(c *gin.Context) {
	cid := c.Param("id")

	id, err := strconv.Atoi(cid)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"messages": err,
		})
	}

	db := database.ConnectDataBase()

	delete, err := db.Prepare("DELETE FROM posts WHERE id=?")

	defer db.Close()

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"messages": err,
		})
	}

	res, err := delete.Exec(id)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"messages": err,
		})
	}

	row, err := res.RowsAffected() // Get number rows deleted
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"messages": err,
		})
	}

	rows := int(row)

	c.IndentedJSON(http.StatusOK, gin.H{
		"messages": fmt.Sprintf("Delete post %d successful %d", id, rows),
		"rows":     rows,
	})
}
