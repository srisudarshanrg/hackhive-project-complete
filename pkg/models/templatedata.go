package models

type TemplateData struct {
	CustomErrors map[string]string
	CustomInfo   map[string]string
	Data         map[string]interface{}
	CSRFToken    string
}
