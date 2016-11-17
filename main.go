package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/redijedi/dropbeat/beater"
)

func main() {
	err := beat.Run("dropbeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
