package database

import (
	"sync"
)

var (
	instance *JmeterDb
	once     sync.Once
)

type JmeterDb struct {
	Id             string
	JmxFile        string
	Namespace      string
	InjectorNumber int
	CsvSplit       int
	EnableReport   bool
}

func Connect() *JmeterDb {

	once.Do(func() {

		instance = &JmeterDb{Id: "Running1"}
	})

	return instance
}

func Get() *JmeterDb {

	instance = Connect()

	return instance
}

func Set(cfg JmeterDb) bool {

	instance = Connect()

	instance.JmxFile = cfg.JmxFile
	instance.Namespace = cfg.Namespace
	instance.InjectorNumber = cfg.InjectorNumber
	instance.CsvSplit = cfg.CsvSplit
	instance.EnableReport = cfg.EnableReport

	return true

}
