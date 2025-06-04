package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/powiedl/myGoWebApplication/internal/models"
)

// AppConfig is a struct holding the application's configuration
type AppConfig struct {
	TemplateCache map[string]*template.Template
	UseCache      bool
	InProduction  bool
	Session       *scs.SessionManager
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	Basedir       string
	MailChan      chan models.MailData
}