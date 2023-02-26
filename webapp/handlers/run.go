package handlers

import (
	"net/http"

	"github.com/drodrigues3/jmeter-k8s-starterkit/database"
	"github.com/drodrigues3/jmeter-k8s-starterkit/log"
	"github.com/gin-gonic/gin"
)

type LoginForm struct {
	JmxFile        string `form:"jmx-file" binding:"required"`
	Namespace      string `form:"namespace" binding:"required"`
	InjectorNumber int    `form:"injector-number" binding:"required"`
	CsvSplit       int    `form:"csv-split" binding:"required"`
	EnableReport   string `form:"enable-report" `
}

func PreRun(c *gin.Context) {

	var args LoginForm
	var EnableReport bool

	err := c.ShouldBind(&args)

	if err != nil {
		log.Error().Err(err).Interface("dic", args).Msg("Was not possible to bind form values with Types defined in the code")
		c.Redirect(http.StatusSeeOther, "/?error_type=Bind")
		return
	}

	// Check if EnableReport is enabled and convert it to boolean
	if args.EnableReport == "on" {
		EnableReport = true
	}

	jmeterCfg := database.JmeterDb{
		JmxFile:        args.JmxFile,
		Namespace:      args.Namespace,
		InjectorNumber: args.InjectorNumber,
		CsvSplit:       args.CsvSplit,
		EnableReport:   EnableReport,
	}

	log.Debug().Interface("dict", jmeterCfg)

	database.Set(jmeterCfg)

	log.Info().Msg("Form data successfully saved")

	c.Redirect(http.StatusSeeOther, "/run")
}

func Run(c *gin.Context) {

	args := database.Get()
	c.HTML(http.StatusOK, "run.tmpl", gin.H{
		"JmxFile":  args.EnableReport,
		"username": "myname",
		"Args":     ""})
}
