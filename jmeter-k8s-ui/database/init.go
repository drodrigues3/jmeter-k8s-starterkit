package database

import "gorm.io/gorm"

func AutoMigrate(db *gorm.DB) {

	db.AutoMigrate(
		&JmeterDb{},
		&JMXFilesListDb{})
}
