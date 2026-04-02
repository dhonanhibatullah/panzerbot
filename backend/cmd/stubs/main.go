package main

import (
	"os"

	_ "github.com/dhonanhibatullah/panzerbot/backend/docs/swagger"
	domainappstubs "github.com/dhonanhibatullah/panzerbot/backend/internal/domain/app/stubs"
)

// @title Panzerbot API
// @version 1.0
// @description The API documentation for Panzerbot API endpoints.
func main() {
	os.Exit(domainappstubs.Run())
}
