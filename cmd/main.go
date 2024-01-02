package main

import (
	"os"

	"github.com/taylormonacelli/navyloss"
)

func main() {
	code := navyloss.Execute()
	os.Exit(code)
}
