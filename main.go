package main

import (
	"bytes"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/Oktasuke/TCCLighting/controllers"
	"github.com/Oktasuke/TCCLighting/models"

	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
)

func main() {
	/**
	 * load config
	 */
	config := models.NewConfig()
	var (
		cf string
	)
	flag.StringVar(&cf, "conf", "resources/config.toml", "TOML形式のconfigファイルのPATH")
	flag.Parse()
	_, err := toml.DecodeFile(cf, &config)
	if err != nil || config.Server.Listen_port == "" {
		log.Fatal("Listen_port must be set")
	}

	/**
	 * router setting
	 */
	gin.SetMode(config.Server.Gin_mode)

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
			"response": "heal check is OK",
		})
	})

	/**
	 * facebook Webhooke reciver
	 */
	router.POST("/facebook", func(c *gin.Context) {
		log.Printf("[INFO] facebook POST call back received.")
		go c.String(http.StatusOK, "")

		bufbody := new(bytes.Buffer)
		bufbody.ReadFrom(c.Request.Body)
		nls := controllers.NewLightSwitcher()
		if nls.IsIlluminate(bufbody.Bytes()) {
			log.Print("[INFO] ON")
		} else {
			log.Print("[INFO] OFF")
		}
	})

	/**
	 * facebook Webhook setting up call this url
	 * refs : https://developers.facebook.com/docs/graph-api/webhooks#setupget
	 */
	router.GET("/facebook", func(c *gin.Context) {
		if c.Query("hub.mode") == "subscribe" && c.Query("hub.verify_token") == config.Facebook.Verify_token {
			log.Printf("[INFO] called facebook WebHook subscribe.")
			c.String(http.StatusOK, c.Query("hub.challenge"))
		} else {
			c.String(http.StatusBadRequest, "UnKnown Rquest")
			log.Printf("[WARN] reciveing UnKnown Rquest")
		}
	})

	router.Run(":" + config.Server.Listen_port)
}
