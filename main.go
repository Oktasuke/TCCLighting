package main

import (
	"bytes"
	"flag"
	"net/http"
	"time"

	"github.com/Oktasuke/TCCLighting/controllers"
	"github.com/Oktasuke/TCCLighting/models"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
)
var config = models.NewConfig()

func main() {
	/**
	 * config loading
	 */
	var (
		cf string
	)
	flag.StringVar(&cf, "conf", "resources/config.toml", "TOML形式のconfigファイルのPATH")
	flag.Parse()
	_, err := toml.DecodeFile(cf, &config)
	if err != nil || config.ServerInfo.ListenPort == "" {
		glog.Fatal("Listen_port must be set")
	}
	li := controllers.NewLightSwitcher(config.ShopInfo, config.WeMoInfo)
	/**
	 * router setting
	 */
	gin.SetMode(config.ServerInfo.GinMode)

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")

	/**
	 * healthCheck reciver
	 */
	router.GET("/healthcheck", func(c *gin.Context) {
		c.HTML(http.StatusOK, "healthcheck.tmpl.html", gin.H{
			"status":   http.StatusOK,
			"date":     time.Now(),
			"response": "healthCheck is OK",
		})
	})

	/**
	 * facebook Webhooke receiver
	 */
	router.POST("/facebook", func(c *gin.Context) {
		glog.Info("facebook POST call back was received.")
		go c.String(http.StatusOK, "")

		bufbody := new(bytes.Buffer)
		bufbody.ReadFrom(c.Request.Body)
		li.IlluminateCtrl(bufbody.Bytes(), controllers.FACE_BOOK)
	})

	/**
	 * facebook Webhook setting up call this url
	 * refs : https://developers.facebook.com/docs/graph-api/webhooks#setupget
	 */
	router.GET("/facebook", func(c *gin.Context) {
		if c.Query("hub.mode") == "subscribe" && c.Query("hub.verify_token") == config.FacebookInfo.VerifyToken {
			glog.Info("facebook Webhook subscribe.")
			c.String(http.StatusOK, c.Query("hub.challenge"))
		} else {
			c.String(http.StatusBadRequest, "UnKnown Rquest")
			glog.Warning("UnKnown Request was received")
		}
	})

	router.Run(":" + config.ServerInfo.ListenPort)
}
