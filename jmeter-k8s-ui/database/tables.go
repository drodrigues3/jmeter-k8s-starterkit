package database

import (
	"gorm.io/gorm"
)

type JmeterDb struct {
	gorm.Model
	JmxFile        string
	Namespace      string
	InjectorNumber int
	CsvSplit       int
	EnableReport   bool
}

type JMXFilesListDb struct {
	gorm.Model
	NameFile     string
	UniqNameFile string
}
