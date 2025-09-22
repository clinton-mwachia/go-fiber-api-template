package utils

import (
	"log"
	"os"
)

func EnsureUploadsFolder() {
	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		if err := os.Mkdir("uploads", os.ModePerm); err != nil {
			log.Fatal("Failed to create uploads folder:", err)
		}
		log.Println("✅ Created uploads folder")
	}
}
