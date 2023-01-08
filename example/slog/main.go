package main

import (
	"golang.org/x/exp/slog"
	"os"
)

func main() {
	//textHandler := slog.NewTextHandler(os.Stdout)
	//logger := slog.New(textHandler)

	jsonHandler := slog.NewJSONHandler(os.Stdout) // 👈
	logger := slog.New(jsonHandler)

	logger.Info("Go is the best language!")
}
