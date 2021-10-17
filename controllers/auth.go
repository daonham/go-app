package controllers

import (
	"net/http"

	"github.com/daonham/go-app/forms"
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

func RefreshToken(c *gin.Context) {
	var tokenForm forms.Token

	if err := c.ShouldBindJSON(&tokenForm); err != nil {
		c.JSON(http.StatusBadRequest, helper.ResponseError(http.StatusBadRequest, "Refresh token is required", err.Error(), helper.EmptyObj{}))
		return
	}

	tokens, message, err := models.RefreshToken(tokenForm)

	if err != nil {
		c.JSON(http.StatusUnauthorized, helper.ResponseError(http.StatusBadRequest, message, err.Error(), helper.EmptyObj{}))
		return
	}

	c.JSON(http.StatusOK, tokens)
}
