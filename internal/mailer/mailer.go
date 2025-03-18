package mailer

import "embed"

const (
	FromName = "Sodia"
	maxRetries = 3
	UserWelcomeTemplate = "user_invitation.go.tmpl"
)

//go:embed "template"
var FS embed.FS

type Client interface {
	Send(templateFile, username, email string, data any, isSandbox bool) error
}