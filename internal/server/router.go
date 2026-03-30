// Package server provides HTTP server setup and routing
package server

import (
	"github.com/gin-gonic/gin"
	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/config"
	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/pkg/logger"
	"github.com/wan-h/awesome-digital-human-live2d/go-backend/internal/server/handlers"
)

const (
	GlobalPrefix = "/adh"
)

// SetupRouter creates and configures the Gin router
func SetupRouter(cfg *config.Config) *gin.Engine {
	if cfg.Common.LogLevel == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(RequestLogger())
	router.Use(CORSMiddleware())

	setupRoutes(router)

	return router
}

// setupRoutes configures all API routes
func setupRoutes(router *gin.Engine) {
	agentHandler := handlers.NewAgentHandler()
	ttsHandler := handlers.NewTTSHandler()

	v0 := router.Group(GlobalPrefix)
	{
		agent := v0.Group("/agent/v0")
		{
			agent.GET("/engine", agentHandler.GetEngineList)
			agent.GET("/engine/default", agentHandler.GetDefaultEngine)
			agent.GET("/engine/:engine", agentHandler.GetEngineParams)
			agent.POST("/engine/:engine", agentHandler.CreateConversation)
			agent.POST("/engine", agentHandler.StreamInfer)
		}

		tts := v0.Group("/tts/v0")
		{
			tts.GET("/engine", ttsHandler.GetEngineList)
			tts.GET("/engine/default", ttsHandler.GetDefaultEngine)
			tts.GET("/engine/:engine", ttsHandler.GetEngineParams)
			tts.GET("/engine/:engine/voice", ttsHandler.GetVoiceList)
			tts.POST("/engine", ttsHandler.Infer)
		}
	}
}

// RequestLogger logs HTTP requests
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Infof("[HTTP] %s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
	}
}

// CORSMiddleware handles CORS headers
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Request-Id, User-Id")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
