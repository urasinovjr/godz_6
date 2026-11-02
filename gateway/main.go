package main

import (
	"fmt"
	"gateway/internal/api"
	"net/http"
)

func main() {
	fmt.Println("Gateway запущен на порту 8080")
	err := http.ListenAndServe(":8080", api.Logging(api.NewRouter()))
	if err != nil {
		fmt.Println("Ошибка! Не удалось запустить сервер:", err)
	}
}
