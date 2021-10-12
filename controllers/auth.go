package controllers

import (
	"net/http"

	"github.com/daonham/go-app/helper"
	"github.com/daonham/go-app/models"
	"github.com/gin-gonic/gin"
)

func TokenValid(c *gin.Context) {
	tokenAuth, err := models.ExtractTokenMetadata(c.Request)

	//Token either expired or not valid
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, helper.ResponseError(http.StatusUnauthorized, "Please login to continue", err.Error(), helper.EmptyObj{}))
		return
	}

	userID := tokenAuth.UserID

	c.Set("userID", userID)
}
