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

var message string
var alert string

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

	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/login", loginHandler).Methods("GET")
	router.HandleFunc("/signin-form", signinFormHandler).Methods("GET")
	router.HandleFunc("/password-change", passwordChangeHandler) //POST i GET
	router.Handle("/account-delete", accountDeleteHandler())
	router.HandleFunc("/logout", logoutHandler).Methods("GET")
	router.HandleFunc("/registration-ouath", registrationOauthHandler).Methods("GET")
	router.HandleFunc("/verify/{verifystr}", verifyHandler).Methods("GET")
	router.HandleFunc("/registration", registrationHandler)
	router.HandleFunc("/account-edit", accountEditHandler)
	router.HandleFunc("/filetesty", Filetesty).Methods("POST")
	router.HandleFunc("/etapy", etapyHandler).Methods("GET")
	router.HandleFunc("/checkuserexists", checkUserExistsHandler).Methods("GET")

	router.Handle("/login", loginEmailHandler()).Methods("POST")
	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)
	router.HandleFunc("/404", notFoundHandler)

	//router.HandleFunc("/login", LoginEmailHandler).Methods("POST")
	//router.Handle("/login", loginEmailHandler()).Methods("POST")

	router.Handle("/login-google", google.StateHandler(stateConfig, google.LoginHandler(oauth2Config, nil))).Methods("GET")
	//router.Handle("/login-facebook", google.StateHandler(stateConfig, google.LoginHandler(oauth2Config, nil))).Methods("GET")
	router.Handle("/loginoauth", google.StateHandler(stateConfig, google.CallbackHandler(oauth2Config, loginOauthHandler(), nil))).Methods("GET")

	//staticFileDirectory := http.Dir("/var/www/runhill.cz/static")
	staticFileDirectory := http.Dir(utils.StaticPath)
	staticFileHandler := http.StripPrefix("/static/", http.FileServer(staticFileDirectory))
	router.PathPrefix("/static/").Handler(staticFileHandler).Methods("GET")
	return router
}
