package syslog

import (
	"log"

	"gorm.io/gorm"
)

type SystemLog struct {
	ID        uint   `gorm:"primarykey"`
	Level     string `gorm:"size:20;not null"`
	Message   string `gorm:"type:text;not null"`
	Source    string `gorm:"size:100"`
	Stack     string `gorm:"type:text"`
	CreatedAt interface{} `gorm:"autoCreateTime"`
}

func (SystemLog) TableName() string {
	return "system_logs"
}

var globalDB *gorm.DB

func Init(db *gorm.DB) {
	globalDB = db
	db.AutoMigrate(&SystemLog{})
}

func Info(message, source string) {
	writeLog("info", message, source)
}

func Warn(message, source string) {
	writeLog("warn", message, source)
}

func Error(message, source string) {
	writeLog("error", message, source)
}

func writeLog(level, message, source string) {
	if globalDB == nil {
		return
	}
	go func() {
		if err := globalDB.Create(&SystemLog{Level: level, Message: message, Source: source}).Error; err != nil {
			log.Printf("Failed to write system log: %v", err)
		}
	}()
}
