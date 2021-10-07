package controllers

import (
	"github.com/daonham/go-app/database"
	"github.com/gin-gonic/gin"
)

type Post struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func ReadPost(c *gin.Context) {
	db := database.ConnectDataBase()

	rows, err := db.Query("SELECT id, title, content FROM posts WHERE id = " + c.Param("id"))

	if err != nil {
		c.JSON(500, gin.H{
			"messages": "Story not found",
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

	c.JSON(200, post)
}
