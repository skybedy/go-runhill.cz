package routes

import (
	"net/http"

	"github.com/dghubble/gologin/v2"
	"github.com/dghubble/gologin/v2/google"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	googleOAuth2 "golang.org/x/oauth2/google"
	"runhill.cz/utils"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter()

	oauth2Config := &oauth2.Config{
		ClientID:     utils.GoogleClientId,
		ClientSecret: utils.GoogleClientSecret,
		RedirectURL:  utils.GoogleRedirectUrl,
		Endpoint:     googleOAuth2.Endpoint,
		Scopes:       []string{"profile", "email"},
	}
	stateConfig := gologin.DebugOnlyCookieConfig

	router.HandleFunc("/", indexHandler).Methods("GET")
	router.HandleFunc("/login-options", loginOptionsHandler).Methods("GET")
	router.HandleFunc("/signin-form", signinFormHandler).Methods("GET")
	router.HandleFunc("/account-summary", accountSummaryHandler).Methods("GET")
	router.Handle("/account-delete/{action}", accountDeleteHandler()).Methods("GET")
	router.HandleFunc("/logout", logoutHandler).Methods("GET")
	router.HandleFunc("/login-form", loginFormHandler).Methods("GET")
	router.HandleFunc("/verify/{verifystr}", verifyHandler).Methods("GET")
	router.HandleFunc("/registration", registrationHandler).Methods("POST")
	router.HandleFunc("/edit-person", editPersonHandler).Methods("POST")
	router.HandleFunc("/filetesty", Filetesty).Methods("POST")
	router.HandleFunc("/message", messageHandler).Methods("GET")

	router.Handle("/login", LoginEmailHandler()).Methods("POST")

	//router.HandleFunc("/login", LoginEmailHandler).Methods("POST")
	router.Handle("/login", LoginEmailHandler()).Methods("POST")

	router.Handle("/login-google", google.StateHandler(stateConfig, google.LoginHandler(oauth2Config, nil))).Methods("GET")
	//router.Handle("/login-facebook", google.StateHandler(stateConfig, google.LoginHandler(oauth2Config, nil))).Methods("GET")
	router.Handle("/afterlogin", google.StateHandler(stateConfig, google.CallbackHandler(oauth2Config, loginHandler(), nil))).Methods("GET")
	//staticFileDirectory := http.Dir("/var/www/timechip.cz/go-www.timechip.cz/static")
	staticFileDirectory := http.Dir("./static")
	staticFileHandler := http.StripPrefix("/static/", http.FileServer(staticFileDirectory))
	router.PathPrefix("/static/").Handler(staticFileHandler).Methods("GET")
	return router
}
