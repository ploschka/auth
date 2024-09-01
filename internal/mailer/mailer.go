package mailer

import (
	"fmt"

	log "github.com/ploschka/auth/internal/logger"
)

func SendIpWarning(email string, newIp string) error {
	log.Info(fmt.Sprintf("Warning sended to %s", email))
	return nil
}
