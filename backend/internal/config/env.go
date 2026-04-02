package config

import "github.com/dhonanhibatullah/panzerbot/backend/internal/utils"

var (
	AppPort            int
	CorsAllowedOrigins []string
)

func LoadEnv() (err error) {
	err = utils.EnvLoad(".env")
	if err != nil {
		return err
	}

	defAppPort := 6767
	AppPort, err = utils.EnvGetInt("APP_PORT", &defAppPort)
	if err != nil {
		return err
	}

	CorsAllowedOrigins, err = utils.EnvGetStrings("FRONTEND_BASE_URL", nil)
	if err != nil {
		return err
	}

	return nil
}
