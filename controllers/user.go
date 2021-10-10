package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/daonham/go-app/database"
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
	db := database.ConnectDB()

	rows, err := db.Query("SELECT id, email, name, createdAt, role FROM user")

	defer db.Close()

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status": gin.H{
				"code":    http.StatusBadRequest,
				"message": "Post not found",
			},
		})
	}

	users := []User{}

	for rows.Next() {
		var id int
		var name, email, role string
		var createdAt time.Time

		err = rows.Scan(&id, &email, &name, &createdAt, &role)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"status": gin.H{
					"code":    http.StatusBadRequest,
					"message": err,
				},
			})
		}

		user := User{}

		user.Id = id
		user.Email = email
		user.Name = name
		user.CreatedAt = createdAt
		user.Role = role

		users = append(users, user)
	}

	c.IndentedJSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	cid := c.Param("id")
	uid, err := strconv.Atoi(cid)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status": gin.H{
				"code":    http.StatusBadRequest,
				"message": err,
			},
		})
	}

	db := database.ConnectDB()

	rows, err := db.Query("SELECT id, email, name, createdAt, role FROM user WHERE id=?", uid)

	defer db.Close()

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status": gin.H{
				"code":    http.StatusBadRequest,
				"message": "Post not found",
			},
		})
	}

	user := User{}

	for rows.Next() {
		var id int
		var name, email, role string
		var createdAt time.Time

		err = rows.Scan(&id, &email, &name, &createdAt, &role)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"status": gin.H{
					"code":    http.StatusBadRequest,
					"message": err,
				},
			})
		}

		user.Id = id
		user.Email = email
		user.Name = name
		user.CreatedAt = createdAt
		user.Role = role
	}

	c.IndentedJSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context) {

	type CreateUser struct {
		Email string `form:"email" json:"email" binding:"required"`
		Name  string `form:"name" json:"name" binding:"required"`
		Pass  string `form:"pass" json:"pass" binding:"required"`
		Role  string `form:"role" json:"role"`
	}

	var create CreateUser

	if err := c.ShouldBindJSON(&create); err == nil {
		db := database.ConnectDB()

		insert, err := db.Prepare("INSERT INTO user(email, name, pass, createdAt, role) VALUES(?,?,?,?,?)")

		defer db.Close()

		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"status": gin.H{
					"code":    http.StatusBadRequest,
					"message": err,
				},
			})
		}

		timeStamp := time.Now().UTC().Format("2006-01-02 15:04:05") // Format timestamp in mysql

		res, err := insert.Exec(create.Email, create.Name, create.Pass, timeStamp, create.Role)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"status": gin.H{
					"code":    http.StatusBadRequest,
					"message": err,
				},
			})
		}

		lastId, err := res.LastInsertId()
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{
				"status": gin.H{
					"code":    http.StatusBadRequest,
					"message": err,
				},
			})
		}

		id := int(lastId)

		c.IndentedJSON(http.StatusCreated, gin.H{
			"status": gin.H{
				"code":    http.StatusCreated,
				"message": "Post inserted",
			},
			"id": id,
		})
	} else {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"status": gin.H{
				"code":    http.StatusBadRequest,
				"message": err.Error(),
			},
		})
	}
}
