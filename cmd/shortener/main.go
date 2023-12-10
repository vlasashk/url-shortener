package main

import "github.com/rs/zerolog"

func main() {
	logger := zerolog.Logger{}
	logger.Info().Msg("Hi")
}
