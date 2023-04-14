package config

import (
	"github.com/alexedwards/scs/v2"
	"html/template"
	"log"
)

// Create an AppConfig type to store configurations

type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	Inproduction  bool
	Session       *scs.SessionManager
}
