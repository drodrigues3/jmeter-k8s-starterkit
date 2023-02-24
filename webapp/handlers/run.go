package handlers

import (
	"log"
	"net/http"

	"github.com/drodrigues3/jmeter-k8s-starterkit/database"
	"github.com/gin-gonic/gin"
)

type LoginForm struct {
	JmxFile        string `form:"jmx-file" binding:"required"`
	Namespace      string `form:"namespace" binding:"required"`
	InjectorNumber int    `form:"injector-number" binding:"required"`
	CsvSplit       int    `form:"csv-split" binding:"required"`
	EnableReport   bool   `form:"enable-report" `
}

func PreRun(c *gin.Context) {

	var args LoginForm
	err := c.ShouldBind(&args)

	if err != nil {
		log.Fatal(err)
	}

	cfg := database.JmeterDb{
		JmxFile:        args.JmxFile,
		Namespace:      args.Namespace,
		InjectorNumber: args.InjectorNumber,
		CsvSplit:       args.CsvSplit,
		EnableReport:   args.EnableReport,
	}
	log.Println(cfg)
	database.Set(cfg)

	c.Redirect(http.StatusSeeOther, "/run")
}

func Run(c *gin.Context) {

	args := database.Get()
	c.HTML(http.StatusOK, "run.tmpl", gin.H{
		"JmxFile":  args.JmxFile,
		"username": "myname",
		"Args":     ""})
}
