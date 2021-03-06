package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/dghubble/gologin/v2/facebook"
	"github.com/dghubble/gologin/v2/google"
	"github.com/dghubble/sessions"
	"github.com/gorilla/mux"
	"runhill.cz/db"
	"runhill.cz/utils"
)

const (
	sessionSecret    = "example cookie signing secret"
	sessionUserKey   = "sessionID"
	sessionUsername  = "sessionName"
	sessionFirstName = "sessionFirstName"
	sessionLastName  = "sessionLastName"
	sessionEmail     = "sessionEmail"
	sessionVerify    = "sessionVerify"
	sessionIdo       = "sessionIdo"
	sessionOauth     = "sessionOauth"
)

//var sessionUrl *sessions.Session

func loginHandler(res http.ResponseWriter, req *http.Request) {
	/*
		oauth2ConfigFb := &oauth2.Config{
			ClientID:     "872931336611060",
			ClientSecret: "df28dc60eaa5e26d625f19806207322d",
			RedirectURL:  "http://localhost:8080/google/callback",
			Endpoint:     facebookOAuth2.Endpoint,
			Scopes:       []string{"profile", "email"},
		}*/

	//facebook.StateHandler(stateConfig, facebook.LoginHandler(oauth2ConfigFb, nil))
	sessionUrl := utils.SessionStore.New("url")
	sessionUrl.Values["url"] = fmt.Sprint(req.URL)
	sessionUrl.Values["refererBeforeLast"] = req.Referer()
	sessionUrl.Save(res)

	session, _ := utils.SessionStore.Get(req, utils.SessionName)
	fmt.Println(session)
	if session != nil {
		http.Redirect(res, req, "/", http.StatusFound)
		//utils.SessionStore.Destroy(res, utils.SessionName)
	}

	utils.ExecuteTemplate(res, "login.html", struct {
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
		Years []string
	}{
		Title: "Login",
		Login: utils.SessionExists(utils.SessionName, req),
		Years: utils.YearsArr(16),
	})
}

type Person struct {
	Ido              string
	Firstname        string
	Surname          string
	Gender           string
	Email            string
	Birdthyear       string
	Password         *string
	PasswordFromForm string
	Auth             byte
}

type PersonFromForm struct {
	Firstname  string
	Lastname   string
	Gender     string
	Email      string
	Birdthyear string
	Password   string
}

type passwordChange struct {
	PasswordOld        string
	PasswordNew        string
	PasswordNewConfirm string
	PasswordFromDb     string
}

func loginEmailHandler() http.Handler {
	fn := func(res http.ResponseWriter, req *http.Request) {

		var person Person

		decoder := json.NewDecoder(req.Body)
		//	var person PersonFromForm
		err := decoder.Decode(&person)

		if err != nil {
			panic(err)
		}

		res.Header().Set("Content-Type", "application/json")

		if len(*&person.PasswordFromForm) > 0 && len(person.Email) > 0 { //pokud je zadan?? email i heslo, hl??d?? to javascript
			var rowExists bool
			sql1 := "SELECT EXISTS(SELECT * FROM osoby WHERE email LIKE ?)"
			err1 := db.Mdb.QueryRow(sql1, person.Email).Scan(&rowExists)

			if err1 != nil {
				fmt.Println("err")
				return // je??t?? nev??m, jestli to ??e??it takhle, chce to hlub???? v??zkum
			}

			if rowExists == false { //pokud v db takov?? email/user nen??
				jsonResponse, err2 := json.Marshal(utils.HttpResponse{Status: "error", Code: 12})
				if err2 != nil {
					//nejaka chyba json.Marshal vy tu asi mela prijit
				}
				res.Write(jsonResponse)
				//utils.Message(res, req, "alert-danger", "Registrace na z??klad?? uveden??ho emailu neexistuje, zkuste to znovu")
			} else { // pokud v db takov?? email/user je
				sql2 := "SELECT ido,password,jmeno,authorization FROM osoby WHERE email LIKE ?"
				err2 := db.Mdb.QueryRow(sql2, person.Email).Scan(&person.Ido, &person.Password, &person.Firstname, &person.Auth)
				if err2 != nil {
					fmt.Println(err2)
					return
				}

				if person.Auth == 0 {
					//http.Redirect(res, req, "/message?from=loginnoauthorise&alert=danger", http.StatusFound)
					jsonResponse, err2 := json.Marshal(utils.HttpResponse{Status: "error", Code: 13})
					if err2 != nil {
						//nejaka chyba json.Marshal vy tu asi mela prijit
					}
					res.Write(jsonResponse) //????et je??t?? nebyl funk??n??, nejprve je t??eba dokon??it autorizaci, kter?? v??m byla zasl??na emailem
				} else {
					if person.Password == nil {
						jsonResponse, err3 := json.Marshal(utils.HttpResponse{Status: "error", Code: 15})
						if err3 != nil {
							//nejaka chyba json.Marshal vy tu asi mela prijit
						}
						res.Write(jsonResponse) // je to ok
					} else {

						if utils.ComparePasswords(person.Password, person.PasswordFromForm) == true { //pokud je heslo ok
							session := utils.SessionStore.New(utils.SessionName)
							session.Values[sessionFirstName] = person.Firstname
							session.Values[sessionIdo] = person.Ido
							session.Values[sessionEmail] = person.Email
							session.Values[sessionVerify] = true
							session.Values[sessionOauth] = "false" //jako, ??e to nen?? od Google, apod.. je to debiln??, v??m, ale kdy?? jsem tady dal bolean, tak to nefungovalo
							session.Save(res)

							sessionUrl, _ := utils.SessionStore.Get(req, "url")

							jsonResponse, err2 := json.Marshal(utils.HttpResponse{Status: "ok", Code: 11, Firstname: person.Firstname, Refererer: fmt.Sprint(sessionUrl.Values["refererBeforeLast"])})
							if err2 != nil {
								//nejaka chyba json.Marshal vy tu asi mela prijit
							}
							utils.SessionStore.Destroy(res, "url")
							res.Write(jsonResponse) // je to ok

						} else { //pokud heslo nesouhlas??
							jsonResponse, err2 := json.Marshal(utils.HttpResponse{Status: "error", Code: 14})
							if err2 != nil {
								//nejaka chyba json.Marshal vy tu asi mela prijit
							}
							res.Write(jsonResponse) //Uveden?? heslo k emailu  nen?? spr??vn??, zkuste to znovu.

						}
					}
				}
			}
		}
	}
	return http.HandlerFunc(fn)
}

/**
 *	funkce, kter?? je vol??na po p??ihl????en?? p??es Google
 */
func loginOauthGoogleHandler() http.Handler {
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
			sql2 := "SELECT ido,jmeno,password,authorization FROM osoby WHERE email LIKE ?"
			err2 := db.Mdb.QueryRow(sql2, googleUser.Email).Scan(&person.Ido, &person.Firstname, &person.Password, &person.Auth)
			if err2 != nil {
				fmt.Println(err2)
				return
			}

			if person.Password != nil { //pokud tam je heslo, tzn. u?? tam je registrace, ale ne p??es OAuth
				res.Write([]byte("Registrace pod emailem " + googleUser.Email + " u?? existuje a je pou??it zp??sob p??ihla??ov??n?? pomoc?? hesla, tud???? nen?? mo??no pou????t p??ihl????en?? pomoc?? t??et?? strany"))

			} else { //pokud tam heslo nen??, tak se m????eme p??ihl??sit
				if person.Auth == 0 {
					res.Write([]byte("Zat??m nebyla provedena verifikace ????tu, kter?? byla zasl??na na email " + googleUser.Email))
				} else {
					session.Values[sessionVerify] = true
					session.Values[sessionFirstName] = person.Firstname
					session.Values[sessionIdo] = person.Ido
					session.Values[sessionOauth] = "true"
					session.Values[sessionEmail] = googleUser.Email
					session.Save(res)

					sessionUrl, _ := utils.SessionStore.Get(req, "url")
					http.Redirect(res, req, fmt.Sprint(sessionUrl.Values["refererBeforeLast"]), http.StatusFound)
					utils.SessionStore.Destroy(res, "url")
				}
			}

		} else {
			session.Values[sessionVerify] = false
			session.Values[sessionFirstName] = googleUser.GivenName
			session.Values[sessionLastName] = googleUser.FamilyName
			session.Values[sessionEmail] = googleUser.Email
			session.Save(res)
			http.Redirect(res, req, "/registration-ouath", http.StatusFound)
		}

	}
	return http.HandlerFunc(fn)
}

/**
 *	funkce, kter?? je vol??na po p??ihl????en?? p??es facebbok
 */
func loginOauthFacebookHandler() http.Handler {
	fn := func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		facebookUser, err := facebook.UserFromContext(ctx)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		// 2. Implement a success handler to issue some form of session
		session := utils.SessionStore.New(utils.SessionName)

		if checkUserExists(facebookUser.Email) == true {
			var person Person
			sql2 := "SELECT ido,jmeno,password,authorization FROM osoby WHERE email LIKE ?"
			err2 := db.Mdb.QueryRow(sql2, facebookUser.Email).Scan(&person.Ido, &person.Firstname, &person.Password, &person.Auth)
			if err2 != nil {
				fmt.Println(err2)
				return
			}

			if person.Password != nil { //pokud tam je heslo, tzn. u?? tam je registrace, ale ne p??es OAuth
				res.Write([]byte("Registrace pod emailem " + facebookUser.Email + " u?? existuje a je pou??it zp??sob p??ihla??ov??n?? pomoc?? hesla, tud???? nen?? mo??no pou????t p??ihl????en?? pomoc?? t??et?? strany"))

			} else { //pokud tam heslo nen??, tak se m????eme p??ihl??sit
				if person.Auth == 0 {
					res.Write([]byte("Zat??m nebyla provedena verifikace ????tu, kter?? byla zasl??na na email " + facebookUser.Email))
				} else {
					session.Values[sessionVerify] = true
					session.Values[sessionFirstName] = person.Firstname
					session.Values[sessionIdo] = person.Ido
					session.Values[sessionOauth] = "true"
					session.Values[sessionEmail] = facebookUser.Email
					session.Save(res)

					sessionUrl, _ := utils.SessionStore.Get(req, "url")
					http.Redirect(res, req, fmt.Sprint(sessionUrl.Values["refererBeforeLast"]), http.StatusFound)
					utils.SessionStore.Destroy(res, "url")
				}
			}

		} else {
			splitName := strings.Split(facebookUser.Name, " ")
			session.Values[sessionVerify] = false
			session.Values[sessionFirstName] = splitName[0]
			session.Values[sessionLastName] = splitName[1]
			session.Values[sessionEmail] = facebookUser.Email
			fmt.Println(session)

			session.Save(res)
			http.Redirect(res, req, "/registration-ouath", http.StatusFound)
		}

	}
	return http.HandlerFunc(fn)
}

func checkUserExistsHandler(res http.ResponseWriter, req *http.Request) {
	var exists bool
	sql1 := "SELECT EXISTS(SELECT * FROM osoby WHERE email LIKE ?)"
	row := db.Mdb.QueryRow(sql1, req.URL.Query().Get("email"))
	err := row.Scan(&exists)
	if err != nil {
		panic(err)
	}
	json.NewEncoder(res).Encode(exists)
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
	session, _ := utils.SessionStore.Get(req, utils.SessionName)
	if session != nil {
		if req.Method == "GET" {
			utils.SessionStore.Destroy(res, utils.SessionName)
		}

		refererRoute := strings.Split(req.Referer(), "/")
		if refererRoute[3] == "login-form" || refererRoute[3] == "login-form#" || refererRoute[3] == "login" || refererRoute[3] == "account-edit" || refererRoute[3] == "password-change" {
			http.Redirect(res, req, "/", http.StatusFound)
		} else {
			http.Redirect(res, req, req.Referer(), http.StatusFound)
		}
	} else {
		http.Redirect(res, req, "/login", http.StatusFound)
	}
}

func registrationOauthHandler(res http.ResponseWriter, req *http.Request) {
	var inputPassword bool
	//pokud by tady byly session, tak vymazat
	session, err := utils.SessionStore.Get(req, utils.SessionName)
	if err != nil {
		fmt.Println(err)
	}

	if session != nil {
		http.Redirect(res, req, "/", http.StatusFound)
	}

	if req.Referer() == "" {
		inputPassword = false
	}

	//utils.SessionExists(utils.SessionName, req) - m??l jsem to tu, nev??m pro??

	utils.ExecuteTemplate(res, "registration-oauth.html", struct {
		Title     string
		FirstName string
		LastName  string
		Email     string
		Login     interface{}

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

	//session := utils.SessionStore.New("registration")

	//sessionUrl := utils.SessionStore.New("url")
	//sessionUrl.Values["url"] = fmt.Sprint(req.URL)
	//sessionUrl.Save(res)
	session,_ := utils.SessionStore.Get(req,utils.SessionName)
	if session != nil {
			http.Redirect(res, req, "/", http.StatusFound)
	}




	if req.Method == "POST" {
		res.Header().Set("Content-Type", "application/json")
		decoder := json.NewDecoder(req.Body)
		var person PersonFromForm
		err := decoder.Decode(&person)

		if err != nil {
			panic(err)
		}

		var rowExists bool
		sql3 := "SELECT EXISTS(SELECT * FROM osoby WHERE email LIKE ?)"
		err3 := db.Mdb.QueryRow(sql3, person.Email).Scan(&rowExists)
		if err3 != nil {
			fmt.Println("err")
		}
		if rowExists == true { //pokud u?? n??hodou takov?? user existuje
			//session.Values["email"] = person.Email
			//session.Values["hlaska"] = "email_exists_already"
			//session.Save(res)
			//http.Redirect(res, req, "/registration", http.StatusFound)
			//res.Write([]byte("Registrace pod emailem " + person.Email + " u?? existuje, <a href=\"/registration\">zaregistrujte se znovu</a> a pou??ijte jin?? email"))
			jsonResponse, err4 := json.Marshal(utils.HttpResponse{Status: "error", Code: 11})
			if err4 != nil {

			}
			res.Write(jsonResponse)

		} else { //pokud ne
			var sql1 string
			password := utils.PasswordGenerator(person.Password)
			if len(person.Password) > 0 {
				sql1 = "INSERT INTO osoby (jmeno,prijmeni,pohlavi,rocnik,email,jmeno_bd,prijmeni_bd,password) VALUES('" + person.Firstname + "','" + person.Lastname + "'" +
					",'" + person.Gender + "'," + person.Birdthyear + ",'" + person.Email + "',toSlug('" + person.Firstname + "'),toSlug('" + person.Lastname + "'),'" + password + "')"
			} else {
				sql1 = "INSERT INTO osoby (jmeno,prijmeni,pohlavi,rocnik,email,jmeno_bd,prijmeni_bd) VALUES('" + person.Firstname + "','" + person.Lastname + "'" +
					",'" + person.Gender + "'," + person.Birdthyear + ",'" + person.Email + "',toSlug('" + person.Firstname + "'),toSlug('" + person.Lastname + "'))"
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

			utils.SendingEmail(person.Email, "Verifikace", "Autorizaci dokon????te kliknut??m na tento odktaz "+authorizationUrl)

			//session.Values["email"] = person.Email
			//session.Save(res)
			//http.Redirect(res, req, "/registration", http.StatusFound)
			//response := []byte("kokokokok")
			//res.Write(response)

			jsonResponse, err2 := json.Marshal(utils.HttpResponse{Status: "ok", Code: 1})
			if err2 != nil {

			}
			res.Write(jsonResponse)

		}
	} else {
		utils.ExecuteTemplate(res, "registration.html", struct {
			Title string
			Login interface{}
			Years []string
		}{
			Title: "Login",
			Login: utils.SessionExists(utils.SessionName, req),
			Years: utils.YearsArr(16),
		})

	}

}

func notFoundHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusNotFound)
	utils.ExecuteTemplate(res, "404.html", struct {
		Title string
		Login interface{}
	}{
		Title: "404, str??nka nenalezena",
		Login: utils.SessionExists(utils.SessionName, req),
	})
}

func verifyHandler(res http.ResponseWriter, req *http.Request) {
	var message string
	var title string
	var alertType string
	vars := mux.Vars(req)
	var ido int64
	var authorization bool
	sql1 := "SELECT osoby.ido,osoby.authorization FROM osoby,verify_registration WHERE osoby.ido = verify_registration.ido AND verify_str LIKE ?"
	row := db.Mdb.QueryRow(sql1, vars["verifystr"])
	switch err := row.Scan(&ido, &authorization); err {

	case sql.ErrNoRows:
		res.WriteHeader(http.StatusNotFound)
		message = "Z nezn??m??ch d??vod?? tato verifikace na serveru nen??, domn??v??te-li se, ??e jde o chybu, kontaktujte n??s pros??m emailem"
		title = "Chyba"
		alertType = "danger"
	case nil:
		if authorization == false {
			sql2 := "UPDATE osoby SET authorization = 1 WHERE ido = " + strconv.FormatInt(ido, 10)
			_, err2 := db.Mdb.Exec(sql2)
			if err2 != nil {
				panic(err.Error())
			}
			message = "Super, klaplo to a te?? se u?? m????ete regul??rn?? p??ihl??sit"
			title = "??sp????n?? verifikace"
			alertType = "primary"

		} else {
			res.WriteHeader(http.StatusNotFound)
			message = "Tato verifikace ji?? byla provedena a nen?? mo??n?? ji prov??st znovu"
			title = "Chyba"
			alertType = "danger"
		}
	default:
		panic(err)
	}

	utils.ExecuteTemplate(res, "verify.html", struct {
		Title     string
		Login     interface{}
		Message   string
		AlertType string
	}{
		Title:     title,
		Login:     utils.SessionExists(utils.SessionName, req),
		Message:   message,
		AlertType: alertType,
	})

}

func accountDeleteHandler() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		session, _ := utils.SessionStore.Get(req, utils.SessionName)
		if session != nil {
			if req.Method == "GET" {
				utils.ExecuteTemplate(res, "account-delete.html", struct {
					Title string
					Login interface{}
				}{
					Title: "",
					Login: utils.SessionExists(utils.SessionName, req),
				})
			} else if req.Method == "DELETE" {
				var jsonResponse []byte

				sql1 := "DELETE FROM osoby WHERE ido = " + fmt.Sprint(session.Values[sessionIdo])
				_, err1 := db.Mdb.Exec(sql1)
				if err1 != nil {
					jsonResponse, _ = json.Marshal(map[string]string{"status": "error"})
				} else {
					jsonResponse, _ = json.Marshal(map[string]string{"status": "ok"})
				}
				utils.SessionStore.Destroy(res, utils.SessionName)
				res.Write(jsonResponse)

			}
		} else {
			http.Redirect(res, req, "/login", http.StatusFound)
		}
	})
}

func accountEditHandler(res http.ResponseWriter, req *http.Request) {
	session, _ := utils.SessionStore.Get(req, utils.SessionName)

	if req.Method == "GET" {

		var person Person
		if session == nil {
			http.Redirect(res, req, "/login", http.StatusFound)
		} else {
			sql1 := "SELECT ido,jmeno,prijmeni,rocnik, pohlavi, email FROM osoby WHERE ido = ?"
			err1 := db.Mdb.QueryRow(sql1, fmt.Sprint(session.Values[sessionIdo])).Scan(&person.Ido, &person.Firstname, &person.Surname, &person.Birdthyear, &person.Gender, &person.Email)
			if err1 != nil {
				//fmt.Println(err)
				log.Fatal(err1)
			}

			utils.ExecuteTemplate(res, "account-edit.html", struct {
				Title   string
				Login   interface{}
				Person  Person
				Years   []string
				Session *sessions.Session
			}{
				Title:   "",
				Login:   utils.SessionExists(utils.SessionName, req),
				Person:  person,
				Years:   utils.YearsArr(16),
				Session: session,
			})
		}

	} else if req.Method == "POST" {
		res.Header().Set("Content-Type", "application/json")
		decoder := json.NewDecoder(req.Body)
		var person Person
		err := decoder.Decode(&person)

		if err != nil {
			panic(err)
		}

		sql1 := "UPDATE osoby SET jmeno='" + person.Firstname + "',prijmeni='" + person.Surname +
			"',rocnik=" + person.Birdthyear + ",pohlavi='" + person.Gender + "',jmeno_bd=toSlug('" + person.Firstname + "'),prijmeni_bd=toSlug('" + person.Surname + "')" +
			" WHERE ido = " + person.Ido
		_, err1 := db.Mdb.Exec(sql1)

		var jsonResponse []byte
		if err1 != nil {
			//panic(err.Error())
			jsonResponse, _ = json.Marshal(map[string]string{"status": "error"})
		} else {
			session.Values[sessionFirstName] = person.Firstname
			session.Save(res)
			jsonResponse, _ = json.Marshal(map[string]string{"status": "ok"})
		}
		res.Write(jsonResponse)

	}

}

func passwordChangeHandler(res http.ResponseWriter, req *http.Request) {
	session, _ := utils.SessionStore.Get(req, utils.SessionName)
	if req.Method == "POST" {
		var jsonResponse []byte
		res.Header().Set("Content-Type", "application/json")

		decoder := json.NewDecoder(req.Body)
		var pch passwordChange
		err := decoder.Decode(&pch)
		if err != nil {
			panic(err)
		}

		sql1 := "SELECT password FROM osoby WHERE ido LIKE ?"
		err1 := db.Mdb.QueryRow(sql1, session.Values[sessionIdo]).Scan(&pch.PasswordFromDb)
		if err1 != nil {
			fmt.Println(err1)
			return
		}

		if utils.ComparePasswords(&pch.PasswordFromDb, pch.PasswordOld) != true {
			jsonResponse, _ = json.Marshal(map[string]interface{}{"status": "error", "code": 21})
		} else {
			if pch.PasswordNew != pch.PasswordNewConfirm {
				jsonResponse, _ = json.Marshal(map[string]interface{}{"status": "error", "code": 22})
			} else {
				password := utils.PasswordGenerator(pch.PasswordNew)
				sql2 := "UPDATE osoby SET password = '" + password + "' WHERE ido = " + fmt.Sprint(session.Values[sessionIdo])
				_, err2 := db.Mdb.Exec(sql2)

				if err2 != nil {
					panic(err.Error())
				}
				jsonResponse, _ = json.Marshal(map[string]interface{}{"status": "ok"})
			}
		}
		res.Write(jsonResponse)
	} else {
		if session != nil {
			utils.ExecuteTemplate(res, "password-change.html", struct {
				Title string
				Login interface{}
			}{
				Title: "Zm??na hesla",
				Login: utils.SessionExists(utils.SessionName, req),
			})
		} else {
			http.Redirect(res, req, "/login", http.StatusFound)
		}

	}

}
