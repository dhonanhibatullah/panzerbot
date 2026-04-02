package main

import (
	"os"

	domainappstubs "github.com/dhonanhibatullah/panzerbot/backend/internal/domain/app/stubs"
)

func main() {
	os.Exit(domainappstubs.Run())
}
