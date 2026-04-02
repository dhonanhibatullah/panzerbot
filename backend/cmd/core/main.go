package main

import (
	"os"

	_ "github.com/dhonanhibatullah/panzerbot/backend/docs/swagger"
	domainappcore "github.com/dhonanhibatullah/panzerbot/backend/internal/domain/app/core"
)

// @title Panzerbot API
// @version 1.0
// @description The API documentation for Panzerbot API endpoints.
func main() {
	os.Exit(domainappcore.Run())
}
