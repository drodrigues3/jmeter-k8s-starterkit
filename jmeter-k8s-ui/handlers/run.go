package handlers

import (
	"net/http"

	"github.com/drodrigues3/jmeter-k8s-starterkit/database"
	"github.com/drodrigues3/jmeter-k8s-starterkit/log"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SetupJmeter struct {
	JmxFile        string `form:"jmx-file" binding:"required"`
	Namespace      string `form:"namespace" binding:"required"`
	InjectorNumber int    `form:"injector-number" binding:"required"`
	CsvSplit       int    `form:"csv-split" binding:"required"`
	EnableReport   string `form:"enable-report" `
}

func PreRun(c *gin.Context, db *gorm.DB) {

	var args SetupJmeter
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

	err = db.Create(&jmeterCfg).Error

	if err != nil {
		log.Error().Err(err).Interface("dic", args).Msg("Was not possible to save data in database")
		c.Redirect(http.StatusInternalServerError, "/?error_type=DatabaseAdd")
		return
	}

	log.Info().Msg("Form data successfully saved")

	c.Redirect(http.StatusSeeOther, "/run")
}

func Run(c *gin.Context, db *gorm.DB) {
	var JmeterDb *database.JmeterDb
	err := db.Find(&JmeterDb).Last(&database.JmeterDb{}).Error

	if err != nil {
		log.Error().Err(err).Msg("Was not possible retrive data from database")
		c.Redirect(http.StatusInternalServerError, "/?error_type=DatabaseGetAll")
		return
	}

	c.HTML(http.StatusOK, "run.tmpl", gin.H{
		"JmxFile":  JmeterDb.EnableReport,
		"username": "myname",
		"Args":     ""})
}
