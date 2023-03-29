package controller

import (
	"context"
	"github.com/gin-contrib/pprof"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"linkshortener/i18n"
	"linkshortener/lib/tool"
	"linkshortener/log"
	"linkshortener/model"
	"linkshortener/setting"
	"linkshortener/statikFS"
	"net/http"
	"sync"
	"time"
)

var router *gin.Engine

func ReqLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		endTime := time.Now()

		latencyTime := endTime.Sub(startTime)

		reqMethod := c.Request.Method

		reqUri := c.Request.RequestURI

		statusCode := c.Writer.Status()

		clientIP := c.ClientIP()

		log.InfoPrint("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}

func InitRouter() {
	router.GET("ping", func(c *gin.Context) {
		model.SuccessResponse(c, map[string]interface{}{
			"msg": "pong",
		})
	}) //Service Test Interface

	router.GET("/s/:hash", Redirect) //Short link redirection

	router.GET("/api/captcha", Captcha)             //Generate captcha code
	router.POST("/api/generate_link", GenerateLink) //Create link
	router.POST("/api/stats_link", StatsLink)       //Link statistics
	router.POST("/api/delete_link", DeleteLink)     //Delete link

	if setting.Cfg.HTTP.FilesDirEmbed { //Static files
		router.NoRoute(gin.WrapH(http.FileServer(statikFS.StatikFS))) //Use of embedded resources
	} else {
		router.NoRoute(gin.WrapH(http.FileServer(http.Dir(setting.Cfg.HTTP.FilesDirURI)))) //Use of external resources
	}
}

func NewLimiter(reqRate rate.Limit, reqBurst int, reqTimeout time.Duration) gin.HandlerFunc {
	limiters := &sync.Map{}

	return func(c *gin.Context) {
		if c.FullPath() != "" {
			key := c.ClientIP()
			limit, _ := limiters.LoadOrStore(key, rate.NewLimiter(reqRate, reqBurst))

			ctx, cancel := context.WithTimeout(c, reqTimeout)
			defer cancel()

			if err := limit.(*rate.Limiter).Wait(ctx); err != nil {
				localizer := i18n.GetLocalizer(c)
				model.FailureResponse(c, http.StatusTooManyRequests, http.StatusTooManyRequests, localizer.GetMessage("tooManyRequests", nil), "")
			}
		}
		c.Next()
	}
}

func InitController() {

	if setting.Cfg.RunMode == "dev" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router = gin.New()

	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.DebugPrint("%v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}
	router.Use(ReqLogger())

	if setting.Cfg.HTTPLimiter.EnableLimiter {
		router.Use(NewLimiter(rate.Limit(setting.Cfg.HTTPLimiter.LimitRate), setting.Cfg.HTTPLimiter.LimitBurst, time.Duration(setting.Cfg.HTTPLimiter.Timeout)*time.Millisecond))
	}

	SessionSecret := tool.GetToken(16)
	if !setting.Cfg.HTTP.RandomSessionSecret {
		SessionSecret = setting.Cfg.HTTP.SessionSecret
	}
	store := memstore.NewStore([]byte(SessionSecret))
	router.Use(sessions.Sessions("session", store))

	if setting.Cfg.RunMode == "dev" {
		pprof.Register(router) //debug
	}
}

func RunServer() {
	log.InfoPrint("Listening and serving HTTP on %s", setting.Cfg.HTTP.Listen)
	err := router.Run(setting.Cfg.HTTP.Listen)
	if err != nil {
		log.PanicPrint("Start Web Server Fail: %s", err)
	}
}
