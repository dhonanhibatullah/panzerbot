package withstubs

import (
	"os"

	domainappwithstubs "github.com/dhonanhibatullah/panzerbot/backend/internal/domain/app/withstubs"
)

func main() {
	os.Exit(domainappwithstubs.Run())
}
