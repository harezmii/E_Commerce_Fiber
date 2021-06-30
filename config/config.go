package config

import (
	"github.com/joho/godotenv"
	"os"
)

func InitEnv() bool {
	errEnv := godotenv.Load(".env")
	if errEnv != nil {
		return false
	}
	return true
}
func GetEnv(envName string) string {
	if InitEnv() {
		return	os.Getenv(envName)
	}
	return ""
}