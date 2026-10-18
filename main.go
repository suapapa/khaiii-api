// Copyright 2025 Homin Lee <ff4500@gmail.com>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package main provides the REST API for Khaiii Korean morphological analyzer.
package main

import (
	"cmp"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"

	_ "github.com/suapapa/khaiii-api/docs"
	"github.com/suapapa/khaiii-api/pkg/khaiii"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var k *khaiii.Khaiii

// @title			Khaiii API
// @version		1.0
// @description	This is a REST API for Khaiii Korean morphological analyzer.
// @BasePath		/v1
// @schemes		http
func main() {
	secret := strings.TrimSpace(os.Getenv("KHAIII_APIKEY"))
	if secret != "" {
		log.Println("using KHAIII_APIKEY from environment")
	}

	var err error
	k, err = khaiii.New(nil)
	if err != nil {
		log.Fatalf("failed to open khaiii: %v", err)
	}
	defer k.Close()

	root := cmp.Or(os.Getenv("ROOT_PATH"), "/")
	if strings.HasSuffix(root, "/") {
		root = root[:len(root)-1]
	}

	router := gin.Default()

	v1 := router.Group(root + "/v1")
	// add middleware to check secret
	v1.Use(func(c *gin.Context) {
		if c.Request.Method == "GET" || secret == "" {
			c.Next()
			return
		}

		if c.GetHeader("Authorization") != "Bearer "+secret {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		log.Println("authorized")
		c.Next()
	})

	v1.POST("/analyze", analyzeHandler)
	v1.GET("/ping", func(c *gin.Context) {
		c.JSON(200, PingResponse{
			Message: "pong",
		})
	})
	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run() // listen and serve on 0.0.0.0:8080
}

// @Summary		Analyze text
// @Description	Perform morphological analysis on Korean text
// @Tags			Analysis
// @Accept			json
// @Produce		json
// @Param			request	body		reqInput	true	"Text to analyze"
// @Success		200		{object}	AnalyzeResponse
// @Failure		400		{object}	ErrorResponse
// @Router			/analyze [post]
func analyzeHandler(c *gin.Context) {
	var request reqInput
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, ErrorResponse{Error: "Invalid request"})
		return
	}

	AnalyzeResult := k.Analyze(request.Text, "")
	c.JSON(200, AnalyzeResponse{Data: AnalyzeResult})
}

// PingResponse represents the response for the ping endpoint
//
//	@Description	Response for the ping endpoint
type PingResponse struct {
	Message string `json:"message" example:"pong"`
}

// AnalyzeResponse represents the response for the analyze endpoint
//
//	@Description	Response for the analyze endpoint
type AnalyzeResponse struct {
	Data interface{} `json:"data"`
}

// ErrorResponse represents an error response
//
//	@Description	Error response structure
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request"`
}

// reqInput represents the input for analysis
//
//	@Description	Input structure for text analysis
type reqInput struct {
	Text string `json:"text" binding:"required" example:"안녕하세요"`
}
