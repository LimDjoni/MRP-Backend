package helper

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
)

func GetEnvWithKey(key string) string {

	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Error loading .env file")
		dir, errLoad := filepath.Abs(filepath.Dir(os.Args[0]))
		if errLoad != nil {
			fmt.Println("Error loading .env filepath directory")
		}
		environmentPath := filepath.Join(dir, ".env")
		errLoad = godotenv.Load(environmentPath)
	}

	return os.Getenv(key)
}
