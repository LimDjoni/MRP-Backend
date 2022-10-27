package helper

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func GetEnvWithKey(key string) string {

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	environmentPath := filepath.Join(dir, ".env")
	err = godotenv.Load(environmentPath)

	return os.Getenv(key)
}
