package main

import (
	"fmt"
	"http_server/internal/flags"
	"http_server/internal/http"
	"log/slog"
)

func main() {
	fmt.Println("Hello! This will be a websocket server one day!")

	f := flags.ParseHTTP()

	slog.SetLogLoggerLevel(slog.LevelDebug)

	hs := http.New(
		http.WithHost(f.Address),
		http.WithPort(f.Port),
	)

	hs.DebugPrint()

	hs.Start()
}
