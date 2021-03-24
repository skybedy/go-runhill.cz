package routes

import (
	"net/http"

	"github.com/dghubble/gologin/v2"
	"github.com/dghubble/gologin/v2/facebook"
	"github.com/dghubble/gologin/v2/google"
	"golang.org/x/oauth2"
	facebookOAuth2 "golang.org/x/oauth2/facebook"
	googleOAuth2 "golang.org/x/oauth2/google"
	"runhill.cz/utils"
)

func Login(w http.ResponseWriter, r *http.Request) {

	oauth2Config := &oauth2.Config{
		ClientID:     "990650209650-tidphg8cnge229cd3888st5jhkfk1g73.apps.googleusercontent.com",
		ClientSecret: "-sosFxmYeGyU-Pe6QutUvGlo",
		RedirectURL:  "http://localhost:8080/google/callback",
		Endpoint:     googleOAuth2.Endpoint,
		Scopes:       []string{"profile", "email"},
	}

	oauth2ConfigFb := &oauth2.Config{
		ClientID:     "872931336611060",
		ClientSecret: "df28dc60eaa5e26d625f19806207322d",
		RedirectURL:  "http://localhost:8080/google/callback",
		Endpoint:     facebookOAuth2.Endpoint,
		Scopes:       []string{"profile", "email"},
	}

	stateConfig := gologin.DebugOnlyCookieConfig
	google.StateHandler(stateConfig, google.LoginHandler(oauth2Config, nil))

	facebook.StateHandler(stateConfig, facebook.LoginHandler(oauth2ConfigFb, nil))

	utils.ExecuteTemplate(w, "login.html", struct {
		Title string
	}{
		Title: "Login",
	})
}
