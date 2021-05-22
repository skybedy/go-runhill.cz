package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/dghubble/gologin/v2"
	"github.com/dghubble/gologin/v2/facebook"
	"github.com/dghubble/gologin/v2/google"
	"golang.org/x/oauth2"
	facebookOAuth2 "golang.org/x/oauth2/facebook"
	googleOAuth2 "golang.org/x/oauth2/google"
	"runhill.cz/utils"
)

var message string
var alert string

func NewRouter() *mux.Router {
	router := mux.NewRouter()

	googleOauth2Config := &oauth2.Config{
		ClientID:     utils.GoogleClientId,
		ClientSecret: utils.GoogleClientSecret,
		RedirectURL:  utils.GoogleRedirectUrl,
		Endpoint:     googleOAuth2.Endpoint,
		Scopes:       []string{"profile", "email"},
	}

	facebookOauth2Config := &oauth2.Config{
		ClientID:     utils.FacebookClientId,
		ClientSecret: utils.FacebookClientSecret,
		RedirectURL:  utils.FacebookRedirectUrl,
		Endpoint:     facebookOAuth2.Endpoint,
		Scopes:       []string{"email"},
	}

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
	stateConfig := gologin.DebugOnlyCookieConfig

	router.Handle("/login-google", google.StateHandler(stateConfig, google.LoginHandler(googleOauth2Config, nil))).Methods("GET")
	router.Handle("/login-google-callback", google.StateHandler(stateConfig, google.CallbackHandler(googleOauth2Config, loginOauthGoogleHandler(), nil))).Methods("GET")
	router.Handle("/login-facebook", facebook.StateHandler(stateConfig, facebook.LoginHandler(facebookOauth2Config, nil))).Methods("GET")
	router.Handle("/login-facebook-callback", facebook.StateHandler(stateConfig, facebook.CallbackHandler(facebookOauth2Config, loginOauthFacebookHandler(), nil))).Methods("GET")

	//staticFileDirectory := http.Dir("/var/www/runhill.cz/static")
	staticFileDirectory := http.Dir(utils.StaticPath)
	staticFileHandler := http.StripPrefix("/static/", http.FileServer(staticFileDirectory))
	router.PathPrefix("/static/").Handler(staticFileHandler).Methods("GET")
	return router
}
