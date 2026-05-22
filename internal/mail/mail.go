package mail

import (
	"context"
	"fmt"
	"log/slog"
)

type Sender interface {
	SendMagicLink(ctx context.Context, to, link string) error
}

// New returns a Sender based on driver. Only "log" is wired in v0.1; SMTP /
// Postmark drivers land in v0.2 when we go to real email.
func New(driver, from string, logger *slog.Logger) Sender {
	switch driver {
	case "log", "":
		return &logSender{from: from, log: logger}
	default:
		// Unknown driver — fail closed by logging loudly. We'd rather print the
		// link in dev than silently drop it.
		logger.Warn("unknown mail driver, falling back to log", "driver", driver)
		return &logSender{from: from, log: logger}
	}
}

type logSender struct {
	from string
	log  *slog.Logger
}

func (l *logSender) SendMagicLink(_ context.Context, to, link string) error {
	l.log.Info("magic-link email (log driver)",
		"from", l.from,
		"to", to,
		"link", link)
	fmt.Println("---- MAGIC LINK ----")
	fmt.Println("to:  ", to)
	fmt.Println("open:", link)
	fmt.Println("--------------------")
	return nil
}
