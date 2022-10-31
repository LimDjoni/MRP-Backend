package helper

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func GetEnvWithKey(key string) string {

	err := godotenv.Load(".env")

	if err != nil {
		dir, errLoad := filepath.Abs(filepath.Dir(os.Args[0]))
		if errLoad != nil {
			fmt.Println("Error loading .env filepath directory")
		}
		environmentPath := filepath.Join(dir, ".env")
		errLoad = godotenv.Load(environmentPath)
	}

	return os.Getenv(key)
}
