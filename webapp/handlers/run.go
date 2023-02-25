package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/drodrigues3/jmeter-k8s-starterkit/database"
	"github.com/gin-gonic/gin"
)

type LoginForm struct {
	JmxFile        string `form:"jmx-file" binding:"required"`
	Namespace      string `form:"namespace" binding:"required"`
	InjectorNumber string `form:"injector-number" binding:"required"`
	CsvSplit       string `form:"csv-split" binding:"required"`
	EnableReport   string `form:"enable-report" `
}

func PreRun(c *gin.Context) {

	var args LoginForm
	var EnableReport bool

	err := c.ShouldBind(&args)

	if err != nil {
		log.Fatal(err)
	}

	// Validate InjectorNumber and CsvSplit
	injectorNumber, err := strconv.Atoi(args.InjectorNumber)
	if err != nil {
		// Redirect to an error page if InjectorNumber is not an integer
		c.Redirect(http.StatusSeeOther, "/?error_type=injector")
		return
	}

	CsvSplit, err := strconv.Atoi(args.CsvSplit)
	if err != nil {
		// Redirect to an error page if CsvSplit is not an integer
		c.Redirect(http.StatusSeeOther, "/?error_type=CsvSplit")
		return
	}

	// Check if EnableReport is enabled and convert it to boolean
	if args.EnableReport == "on" {
		EnableReport = true
	}

	cfg := database.JmeterDb{
		JmxFile:        args.JmxFile,
		Namespace:      args.Namespace,
		InjectorNumber: injectorNumber,
		CsvSplit:       CsvSplit,
		EnableReport:   EnableReport,
	}
	log.Println(cfg)
	database.Set(cfg)

	c.Redirect(http.StatusSeeOther, "/run")
}

func Run(c *gin.Context) {

	args := database.Get()
	c.HTML(http.StatusOK, "run.tmpl", gin.H{
		"JmxFile":  args.EnableReport,
		"username": "myname",
		"Args":     ""})
}
