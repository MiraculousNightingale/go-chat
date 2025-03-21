package main

import (
	"fmt"
	"http_server/internal/http"
	"log/slog"
)

func main() {
	fmt.Println("Hello! This will be a websocket server one day!")

	slog.SetLogLoggerLevel(slog.LevelDebug)

	hs := http.New(
		http.WithHost("localhost"),
		http.WithPort("8080"),
	)

	hs.DebugPrint()

	hs.Start()
}
