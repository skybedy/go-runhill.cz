package utils

import "github.com/dghubble/sessions"

var GoogleClientId string
var GoogleClientSecret string
var GoogleRedirectUrl string
var ServerWebname string
var SessionName string
var SessionStore = sessions.NewCookieStore([]byte("example cookie signing secret"), nil)
var SmtpPort int
var SmtpServer string
var EmailCharset string
var EmailFrom string
var EmailFromName string

type HttpResponse struct {
	Status    string `json:"status"`
	Code      int    `json:"code"`
	Firstname string `json:"firstname"`
	Refererer string `json:"referer"`
}
