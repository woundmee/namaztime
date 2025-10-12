package main

import (
	"log"
	"namaztimeApi/bots/telegram"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Не удалось загрузить переменные окружения!")
	}

	telegram.Bot()
}
