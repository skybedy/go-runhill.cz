package routes

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/dghubble/gologin/v2/google"
	"github.com/gorilla/mux"
	"runhill.cz/db"
	"runhill.cz/utils"
)

const (
	sessionSecret    = "example cookie signing secret"
	sessionUserKey   = "googleID"
	sessionUsername  = "googleName"
	sessionFirstName = "sessionFirstName"
	sessionLastName  = "googleLastName"
	sessionEmail     = "googleEmail"
	sessionVerify    = "sessionVerify"
	sessionIdo       = "sessionIdo"
)

func loginOptionsHandler(res http.ResponseWriter, req *http.Request) {
	/*
		oauth2ConfigFb := &oauth2.Config{
			ClientID:     "872931336611060",
			ClientSecret: "df28dc60eaa5e26d625f19806207322d",
			RedirectURL:  "http://localhost:8080/google/callback",
			Endpoint:     facebookOAuth2.Endpoint,
			Scopes:       []string{"profile", "email"},
		}*/

	//facebook.StateHandler(stateConfig, facebook.LoginHandler(oauth2ConfigFb, nil))

	utils.ExecuteTemplate(res, "login-options.html", struct {
		Title string
		Login interface{}
	}{
		Title: "Login",
		Login: utils.SessionExists(utils.SessionName, req),
	})
}

func signinFormHandler(res http.ResponseWriter, req *http.Request) {

	utils.ExecuteTemplate(res, "signin-form.html", struct {
		Title string
		Login interface{}
	}{
		Title: "Login",
		Login: utils.SessionExists(utils.SessionName, req),
	})
}

type Person struct {
	Ido        string
	Firstname  string
	Surname    string
	Gender     string
	Email      string
	Birdthyear string
	Password   string
	Oauth      *string
	Auth       byte
}

func LoginEmailHandler() http.Handler {
	fn := func(res http.ResponseWriter, req *http.Request) {
		var person Person
		if len(req.FormValue("password")) > 0 && len(req.FormValue("email")) > 0 {
			var rowExists bool
			sql1 := "SELECT EXISTS(SELECT * FROM osoby WHERE email LIKE ?)"
			err1 := db.Mdb.QueryRow(sql1, req.FormValue("email")).Scan(&rowExists)
			if err1 != nil {
				fmt.Println("err")
				return // ještě nevím, jestli to řešit takhle, chce to hlubší výzkum
			}
			if rowExists == false {
				utils.Message(res, req, "alert-danger", "Registrace na základě uvedeného emailu neexistuje, zkuste to znovu")
			} else {
				sql2 := "SELECT ido,password,jmeno,authorization FROM osoby WHERE email LIKE ?"
				err2 := db.Mdb.QueryRow(sql2, req.FormValue("email")).Scan(&person.Ido, &person.Password, &person.Firstname, &person.Auth)
				if err2 != nil {
					fmt.Println(err2)
					return
				}

				if person.Auth == 0 {
					http.Redirect(res, req, "/message?from=loginnoauthorise&alert=danger", http.StatusFound)
				}

				if utils.ComparePasswords(person.Password, req.FormValue("password")) == true {
					session := utils.SessionStore.New(utils.SessionName)
					session.Values[sessionFirstName] = person.Firstname
					session.Values[sessionIdo] = person.Ido
					session.Values[sessionEmail] = req.FormValue("email")
					session.Values[sessionVerify] = true
					session.Save(res)
					refererRoute := strings.Split(req.Referer(), "/")
					if refererRoute[3] == "login-form" || refererRoute[3] == "login-form#" || refererRoute[3] == "login" || refererRoute[3] == "login-options" || refererRoute[3] == "account" {
						http.Redirect(res, req, "/", http.StatusFound)
					} else {
						http.Redirect(res, req, req.Referer(), http.StatusFound)
					}
				} else {
					utils.Message(res, req, "alert-danger", "Uvedené heslo k emailu "+req.FormValue("email")+" není správné, zkuste to znovu.")
				}
			}
		}
	}
	return http.HandlerFunc(fn)
}

func loginHandler() http.Handler { //nemelo by to být něco jako login soc síte login?
	fn := func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		googleUser, err := google.UserFromContext(ctx)

		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		// 2. Implement a success handler to issue some form of session
		session := utils.SessionStore.New(utils.SessionName)

		if checkUserExists(googleUser.Email) == true {
			var person Person
			//var oauth
			sql2 := "SELECT ido,jmeno,oauth,authorization FROM osoby WHERE email LIKE ?"
			err2 := db.Mdb.QueryRow(sql2, googleUser.Email).Scan(&person.Ido, &person.Firstname, &person.Oauth, &person.Auth)
			if err2 != nil {
				fmt.Println(err2)
				return
			}

			if person.Oauth == nil {
				http.Redirect(res, req, "/message?from=login&alert=danger", http.StatusFound)
			}

			if person.Auth == 0 {
				http.Redirect(res, req, "/message?from=loginnoauthorise&alert=danger", http.StatusFound)
			}

			session.Values[sessionVerify] = true
			session.Values[sessionFirstName] = person.Firstname
			session.Values[sessionIdo] = person.Ido
			session.Values[sessionEmail] = googleUser.Email
			session.Save(res)
			http.Redirect(res, req, req.Referer(), http.StatusFound)
		} else {
			session.Values[sessionVerify] = false
			session.Values[sessionFirstName] = googleUser.GivenName
			session.Values[sessionLastName] = googleUser.FamilyName
			session.Values[sessionEmail] = googleUser.Email
			session.Values["oauth"] = "G"
			session.Save(res)
			http.Redirect(res, req, "/login-form", http.StatusFound)
		}

	}
	return http.HandlerFunc(fn)
}

func checkUserExists(email string) bool {
	var exists bool
	sql1 := "SELECT EXISTS(SELECT * FROM osoby WHERE email LIKE ?)"
	row := db.Mdb.QueryRow(sql1, email)
	err := row.Scan(&exists)
	if err != nil {
		panic(err)
	}
	return exists
}

// logoutHandler destroys the session on POSTs and redirects to home.
func logoutHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		utils.SessionStore.Destroy(res, utils.SessionName)
	}

	refererRoute := strings.Split(req.Referer(), "/")
	if refererRoute[3] == "login-form" || refererRoute[3] == "login-form#" || refererRoute[3] == "login" || refererRoute[3] == "account" {
		http.Redirect(res, req, "/", http.StatusFound)
	} else {
		http.Redirect(res, req, req.Referer(), http.StatusFound)
	}
}

func loginFormHandler(res http.ResponseWriter, req *http.Request) {
	var inputPassword bool
	session, err := utils.SessionStore.Get(req, utils.SessionName)
	if err != nil {
		fmt.Println(err)
	}

	if req.Referer() == "" {
		inputPassword = false
	}

	utils.SessionExists(utils.SessionName, req)

	utils.ExecuteTemplate(res, "login-form.html", struct {
		Title         string
		FirstName     string
		LastName      string
		Email         string
		Login         interface{}
		InputPassword bool
		Oauth         string
		Years         []string
	}{
		Title:         "Login",
		FirstName:     fmt.Sprint(session.Values[sessionFirstName]),
		LastName:      fmt.Sprint(session.Values[sessionLastName]),
		Email:         fmt.Sprint(session.Values[sessionEmail]),
		Login:         utils.SessionExists(utils.SessionName, req),
		InputPassword: inputPassword,
		Oauth:         fmt.Sprint(session.Values["oauth"]),
		Years:         utils.YearsArr(16),
	})
}

func registrationHandler(res http.ResponseWriter, req *http.Request) {
	var sql1 string
	password := utils.PasswordGenerator(req.FormValue("password"))
	if len(req.FormValue("password")) > 0 {
		sql1 = "INSERT INTO osoby (jmeno,prijmeni,pohlavi,rocnik,email,jmeno_bd,prijmeni_bd,password) VALUES('" + req.FormValue("firstname") + "','" + req.FormValue("lastname") + "'" +
			",'" + req.FormValue("gender") + "'," + req.FormValue("birdthyear") + ",'" + req.FormValue("email") + "',toSlug('" + req.FormValue("firstname") + "'),toSlug('" + req.FormValue("lastname") + "'),'" + password + "')"
	} else {
		sql1 = "INSERT INTO osoby (jmeno,prijmeni,pohlavi,rocnik,email,jmeno_bd,prijmeni_bd,oauth) VALUES('" + req.FormValue("firstname") + "','" + req.FormValue("lastname") + "'" +
			",'" + req.FormValue("gender") + "'," + req.FormValue("birdthyear") + ",'" + req.FormValue("email") + "',toSlug('" + req.FormValue("firstname") + "'),toSlug('" + req.FormValue("lastname") + "'),'" + req.FormValue("oauth") + "')"
	}
	dbres, err := db.Mdb.Exec(sql1)

	if err != nil {
		panic(err.Error())
	}
	lastID, err := dbres.LastInsertId()
	authorizationStr := utils.RandStr(124, 0, 60)
	authorizationUrl := utils.ServerWebname + "/verify/" + authorizationStr
	sql2 := "INSERT INTO verify_registration (ido,verify_str) VALUES (" + strconv.FormatInt(lastID, 10) + ",'" + authorizationStr + "')"
	_, err1 := db.Mdb.Exec(sql2)
	if err1 != nil {
		panic(err.Error())
	}

	utils.SendingEmail("visalajka@seznam.cz", "Verifikace", "Autorizaci dokončíte kliknutím na tento odktaz "+authorizationUrl)

	utils.ExecuteTemplate(res, "after-registration.html", struct {
		Title string
		Login interface{}
		Email string
	}{
		Title: "After...",
		Login: utils.SessionExists(utils.SessionName, req),
		Email: req.FormValue("email"),
	})

}

func verifyHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Println("tu")
	vars := mux.Vars(req)
	var ido int64
	var authorization bool
	sql1 := "SELECT osoby.ido,osoby.authorization FROM osoby,verify_registration WHERE osoby.ido = verify_registration.ido AND verify_str LIKE ?"
	row := db.Mdb.QueryRow(sql1, vars["verifystr"])
	switch err := row.Scan(&ido, &authorization); err {

	case sql.ErrNoRows:
		fmt.Println("Neplatná autorizace")
	case nil:
		if authorization == false {
			sql2 := "UPDATE osoby SET authorization = 1 WHERE ido = " + strconv.FormatInt(ido, 10)
			_, err2 := db.Mdb.Exec(sql2)
			if err2 != nil {
				panic(err.Error())
			}
			http.Redirect(res, req, "/message?from=authorizationtrue&alert=primary", http.StatusFound)

		} else {
			http.Redirect(res, req, "/message?from=authorizationmisunderstanding&alert=primary", http.StatusFound)
		}
	default:
		panic(err)
	}

}

func accountSummaryHandler(res http.ResponseWriter, req *http.Request) {
	var person Person
	session, err := utils.SessionStore.Get(req, utils.SessionName)
	if err != nil {
		fmt.Println(err)
		http.Redirect(res, req, "/", http.StatusFound)
	} else {
		sql1 := "SELECT jmeno,prijmeni,rocnik, pohlavi, email FROM osoby WHERE ido = ?"
		err1 := db.Mdb.QueryRow(sql1, fmt.Sprint(session.Values[sessionIdo])).Scan(&person.Firstname, &person.Surname, &person.Birdthyear, &person.Gender, &person.Email)
		if err1 != nil {
			log.Fatal(err1)
		}

		utils.ExecuteTemplate(res, "account-summary.html", struct {
			Title  string
			Login  interface{}
			Person Person
		}{
			Title:  "",
			Login:  utils.SessionExists(utils.SessionName, req),
			Person: person,
		})
	}
}

func accountDeleteHandler() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		session, err := utils.SessionStore.Get(req, utils.SessionName)
		if err != nil {
			fmt.Println(err)
			http.Redirect(res, req, "/", http.StatusFound)
		} else {
			vars := mux.Vars(req)
			action := vars["action"]

			if action == "delete" {
				sql1 := "DELETE FROM osoby WHERE ido = " + fmt.Sprint(session.Values[sessionIdo])
				_, err1 := db.Mdb.Exec(sql1)

				if err1 != nil {
					fmt.Println(err)
					//panic(err.Error())
				}

				utils.SessionStore.Destroy(res, utils.SessionName)
				http.Redirect(res, req, "/message?from=accountdelete&alert=primary", http.StatusFound)

			} else {

				utils.ExecuteTemplate(res, "account-delete.html", struct {
					Title  string
					Login  interface{}
					Action string
				}{
					Title:  "",
					Login:  utils.SessionExists(utils.SessionName, req),
					Action: action,
				})
			}

		}

	})

}

func editPersonHandler(res http.ResponseWriter, req *http.Request) {
	session, err := utils.SessionStore.Get(req, utils.SessionName)
	if err != nil {
		fmt.Println(err)
		http.Redirect(res, req, "/", http.StatusFound)
	} else {

		sql1 := "UPDATE osoby SET jmeno='" + req.FormValue("firstname") + "',prijmeni='" + req.FormValue("surname") +
			"',rocnik=" + req.FormValue("birdthyear") + ",pohlavi='" + req.FormValue("gender") + "',jmeno_bd=toSlug('" + req.FormValue("firstname") + "'),prijmeni_bd=toSlug('" + req.FormValue("surname") + "')" +
			" WHERE ido = " + fmt.Sprint(session.Values[sessionIdo])
		_, err1 := db.Mdb.Exec(sql1)

		if err1 != nil {
			panic(err.Error())
		}
		utils.ExecuteTemplate(res, "message.html", struct {
			Title   string
			Login   interface{}
			Message string
		}{
			Title:   "After...",
			Login:   utils.SessionExists(utils.SessionName, req),
			Message: "Údaje byly změněny a uloženy",
		})
	}
}
