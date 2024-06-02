package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/drodrigues3/jmeter-k8s-starterkit/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routes(router *gin.Engine, db *gorm.DB, cfg *config.Config) {

	// Serve the index page
	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {

		// Get all JMXFiles
		allJMXFiles := GetAllJMXFiles(db)

		listCsvFiles, _ := ListFilesWithPath(cfg.Scenarios.Path + "/" + cfg.Scenarios.DefaultDirectories.Dataset)

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"JMXFiles":     allJMXFiles,
			"CsvFiles":     listCsvFiles,
			"success_type": c.Request.URL.Query().Get("success_type"),
			"error_type":   c.Request.URL.Query().Get("error_type"),
		})
	})

	// Handle with run operation
	router.POST("/pre-run", func(c *gin.Context) {
		PreRun(c, db)
	})
	// Handle with run operation
	router.GET("/run", func(c *gin.Context) {
		Run(c, db)
	})

	// Handle upload form
	router.GET("/upload", func(c *gin.Context) {

		c.HTML(http.StatusOK, "upload.tmpl", gin.H{
			"error_type": c.Request.URL.Query().Get("error_type"),
		})
	})

	// Handle upload request
	router.POST("/upload", func(c *gin.Context) {
		Upload(c, db, cfg)
	})

	router.GET("/events", func(c *gin.Context) {
		flusher := c.Writer.(http.Flusher)

		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("Access-Control-Allow-Origin", "*")

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				data := fmt.Sprintf("data: %s\n\n", time.Now().Format("2006-01-02 15:04:05"))
				fmt.Fprintf(c.Writer, data)
				flusher.Flush()
			case <-c.Request.Context().Done():
				fmt.Println("Connection closed")
				return
			}
		}
	})
}
