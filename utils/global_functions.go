package utils

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
	"runhill.cz/db"
)

func SecToTime(secinput int) string {
	var hh string
	h := secinput / 3600

	if h > 0 {
		hh = strconv.Itoa(h) + ":"
	}

	hz := secinput % 3600
	m := strconv.Itoa(hz / 60)
	s := strconv.Itoa(hz % 60)
	if len(m) == 1 {
		m = "0" + m
	}
	if len(s) == 1 {
		s = "0" + s
	}

	return hh + m + ":" + s
}

/*
func SessionExists(SessionName string, r *http.Request) map[string]interface{} {
	arr := map[string]interface{}{"verify": false, "firstname": "", "oauth": ""} //když jsem tady pro oauth dal bolean, tak to nefungovalo
	val, _ := SessionStore.Get(r, SessionName)
	if val != nil {
		if val.Values["sessionVerify"] == true {
			arr["verify"] = true
			arr["firstname"] = fmt.Sprintf("%v", val.Values["sessionFirstName"])
			arr["oauth"] = fmt.Sprintf("%v", val.Values["sessionOauth"])
		}
	}
	return arr
}*/

func SessionExists(SessionName string, r *http.Request) map[string]interface{} {
	arr1 := map[string]interface{}{"verify": false, "firstname": "", "oauth": ""} //když jsem tady pro oauth dal bolean, tak to nefungovalo
	val, _ := SessionStore.Get(r, SessionName)
	if val != nil {
		if val.Values["sessionVerify"] == true {
			arr1["verify"] = true
			arr1["firstname"] = fmt.Sprintf("%v", val.Values["sessionFirstName"])
			arr1["oauth"] = fmt.Sprintf("%v", val.Values["sessionOauth"])
		}
	}

	var etapa Etapy
	sql1 := "SELECT id,start,cil,toSlug(CONCAT_WS(' ',start,cil)) AS seoUrl FROM etapy WHERE id = 1"
	err := db.Mdb.QueryRow(sql1).Scan(&etapa.Id, &etapa.Start, &etapa.Cil, &etapa.SeoUrl)
	if err != nil {
		panic(err.Error())
	}

	var etapyx []Etapy
	sql2 := "SELECT id,start,cil,toSlug(CONCAT_WS(' ',start,cil)) AS seoUrl FROM etapy ORDER BY id ASC"
	results, err := db.Mdb.Query(sql2)
	if err != nil {
		panic(err.Error())
	}
	for results.Next() {
		var etapax Etapy
		err = results.Scan(&etapax.Id, &etapax.Start, &etapax.Cil, &etapax.SeoUrl)
		etapyx = append(etapyx, etapax)
	}

	//arr3 := map[string]interface{}{"etapa": etapa, "neco": 1}
	//	fmt.Println(etapyx)

	//	arr2 := map[string]map[string]interface{}{"bla": etapa}
	arr3 := map[string]interface{}{"login": arr1, "etapy": EtapyList()}
	return arr3
}

func RandStr(count int, min int, max int) string {
	var str string
	rand.Seed(time.Now().UnixNano())
	Arr := [60]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "R", "Q", "S", "T", "V", "W", "X", "Y", "Z", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	for i := 1; i <= count; i++ {
		str += Arr[rand.Intn(max-min)+min]
	}
	return str
}

func PasswordGenerator(plainPassword string) string {
	bytePassword := []byte(plainPassword)
	hash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.MinCost)
	if err != nil {
		fmt.Println(err)
	}
	return string(hash[:])
}

func ComparePasswords(hashedPassword *string, plainPassword string) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice

	byteHash := []byte(*hashedPassword)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(plainPassword))
	if err != nil {
		return false
	}
	return true
}

// alertType - https://getbootstrap.com/docs/4.0/components/alerts/
func Message(res http.ResponseWriter, req *http.Request, alertType string, message string) {
	ExecuteTemplate(res, "message.html", struct {
		Title     string
		Login     interface{}
		Message   string
		AlertType string
	}{
		Title:     "",
		Login:     SessionExists(SessionName, req),
		Message:   message,
		AlertType: alertType,
	})
}

func SendingEmail(recipient string, subject string, body string) {
	ch := gomail.MessageSetting(gomail.SetCharset("UTF-8"))
	m := gomail.NewMessage(ch)
	m.SetHeader("From", EmailFromName+" <"+EmailFrom+">")
	m.SetHeader("To", recipient)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)
	d := gomail.Dialer{Host: SmtpServer, Port: SmtpPort}
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("mail poslan")
	return
}

func YearsArr(minAge int) []string {
	var str string
	currentYear := time.Now().Year()
	firstYear := currentYear - 99
	lastYear := currentYear - minAge
	for i := firstYear; i <= lastYear; i++ {
		str += strconv.Itoa(i) + ","
	}
	return strings.Split(strings.TrimRight(str, ","), ",")
}

type Etapy struct {
	Id     int
	Start  string
	Cil    string
	SeoUrl string
}

func EtapyListzal() Etapy {
	var etapa Etapy
	sql1 := "SELECT id,start,cil,toSlug(CONCAT_WS(' ',start,cil)) AS seoUrl FROM etapy WHERE id = 1"
	err := db.Mdb.QueryRow(sql1).Scan(&etapa.Id, &etapa.Start, &etapa.Cil, &etapa.SeoUrl)

	if err != nil {
		panic(err.Error())
	}
	return etapa
}

func EtapyList() []Etapy {
	var etapy []Etapy
	sql1 := "SELECT id,start,cil,toSlug(CONCAT_WS(' ',start,cil)) AS seoUrl FROM etapy ORDER BY id ASC"
	results, err := db.Mdb.Query(sql1)
	if err != nil {
		panic(err.Error())
	}
	for results.Next() {
		var etapa Etapy
		err = results.Scan(&etapa.Id, &etapa.Start, &etapa.Cil, &etapa.SeoUrl)
		etapy = append(etapy, etapa)
	}
	return etapy
}
