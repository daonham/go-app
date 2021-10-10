package models

import (
	"errors"
	"time"

	"github.com/daonham/go-app/database"
	"github.com/daonham/go-app/forms"
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

func GetPosts() (posts []Post, message string, err error) {
	db := database.ConnectDB()

	rows, err := db.Query("SELECT id, createdAt, updatedAt, title, content, published, authorId FROM post")

	defer db.Close()

	if err != nil {
		return posts, "Error when connect DB", err
	}

	for rows.Next() {
		var id, authorId int
		var title, content string
		var createdAt, updatedAt time.Time
		var published bool

		err = rows.Scan(&id, &createdAt, &updatedAt, &title, &content, &published, &authorId)
		if err != nil {
			return posts, "Error when connect DB", err
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

	return posts, "Get posts successfully", nil
}

func GetPost(pid int) (post Post, message string, err error) {
	db := database.ConnectDB()

	row := db.QueryRow("SELECT id, createdAt, updatedAt, title, content, published, authorId FROM post WHERE id=?", pid)

	defer db.Close()

	var id, authorId int
	var title, content string
	var createdAt, updatedAt time.Time
	var published bool

	err = row.Scan(&id, &createdAt, &updatedAt, &title, &content, &published, &authorId)
	if err != nil {
		return post, "No rows in result", err
	}

	post.Id = id
	post.CreatedAt = createdAt
	post.UpdatedAt = updatedAt
	post.Title = title
	post.Content = content
	post.Published = published
	post.AuthorId = authorId

	return post, "Get post successfully", nil
}

func CreatePost(create forms.PostForm) (id int64, message string, err error) {
	db := database.ConnectDB()

	insert, err := db.Prepare("INSERT INTO post(updatedAt, title, content, published, authorId) VALUES(?,?,?,?,?)")

	defer db.Close()

	if err != nil {
		return 0, "Error when connect DB", err
	}

	timeStamp := time.Now().UTC().Format("2006-01-02 15:04:05") // Format timestamp in MySQL

	res, err := insert.Exec(timeStamp, create.Title, create.Content, create.Published, create.AuthorId)
	if err != nil {
		return id, "Create post successfully", err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return 0, "Error when connect DB", err
	}

	return lastId, "Create post successfully", nil
}

func UpdatePost(id int, update forms.PostForm) (rows int64, message string, err error) {
	db := database.ConnectDB()

	edit, err := db.Prepare("UPDATE post SET updatedAt=?, title=?, content=?, published=?, authorId=? WHERE id=?")

	defer db.Close()

	if err != nil {
		return 0, "Error when connect DB", err
	}

	timeStamp := time.Now().UTC().Format("2006-01-02 15:04:05") // Format timestamp in MySQL

	res, err := edit.Exec(timeStamp, update.Title, update.Content, update.Published, update.AuthorId, id)

	if err != nil {
		return 0, "Error when connect DB", err
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

func DeletePost(id int) (rows int64, message string, err error) {
	db := database.ConnectDB()

	delete, err := db.Prepare("DELETE FROM post WHERE id=?")

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
		return 0, "Error: Update 0 records", errors.New("updated 0 records")
	}

	return row, "Delete post successfully", nil
}
