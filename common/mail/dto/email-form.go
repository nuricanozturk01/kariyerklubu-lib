package dto

const KariyerKlubuEmail = "kariyerklubu.com@gmail.com"
const KariyerKlubuName = "kariyerklubu"

type EmailTemplateForm struct {
	From       string
	To         string
	Title      string
	Body       string
	Name       string
	Variables  map[string]any
	TemplateID int
}

type EmailForm struct {
	From     string
	To       string
	Name     string
	Subject  string
	Text     string
	HTMLPart string
}
