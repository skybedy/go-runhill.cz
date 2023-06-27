package routes

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"runhill.cz/db"
	"runhill.cz/utils"
)

type Gpx struct {
	XMLName xml.Name `xml:"gpx"`
	Gpx     []Trk    `xml:"trk"`
}

type Trk struct {
	XMLName xml.Name `xml:"trk"`
	Trk     []Trkseg `xml:"trkseg"`
}

type Trkseg struct {
	XMLName xml.Name `xml:"trkseg"`
	Trkseg  []Trkpt  `xml:"trkpt"`
}

type Trkpt struct {
	XMLName xml.Name `xml:"trkpt"`
	Lat     string   `xml:"lat,attr"`
	Lon     string   `xml:"lon,attr"`
	Ele     string   `xml:"ele"`
	Time    string   `xml:"time"`
}

type Tcd struct {
	XMLName xml.Name     `xml:"TrainingCenterDatabase"`
	Tcd     []Activities `xml:"Activities"`
}

type Activities struct {
	XMLName    xml.Name   `xml:"Activities"`
	Activities []Activity `xml:"Activity"`
}

type Activity struct {
	XMLName  xml.Name `xml:"Activity"`
	Activity []Lap    `xml:"Lap"`
}

type Lap struct {
	XMLName xml.Name `xml:"Lap"`
	Lap     []Track  `xml:"Track"`
}

type Track struct {
	XMLName xml.Name     `xml:"Track"`
	Track   []Trackpoint `xml:"Trackpoint"`
}

type Trackpoint struct {
	XMLName        xml.Name `xml:"Trackpoint"`
	Time           string   `xml:"Time"`
	AltitudeMeters string   `xml:"AltitudeMeters"`
	Position       Position `xml:"Position"`
}

type Position struct {
	XMLName          xml.Name `xml:"Position"`
	LatitudeDegrees  string   `xml:"LatitudeDegrees"`
	LongitudeDegrees string   `xml:"LongitudeDegrees"`
}

func parseXML(xmlFile string) {
	/*
		xmlFile, err := os.Open("static/files/test.gpx")
		if err != nil {
			fmt.Println(err)
		}
		defer xmlFile.Close()*/

	var gpx Gpx
	var startTime time.Time
	var finishTime time.Time
	xml.Unmarshal([]byte(xmlFile), &gpx)
	if len(gpx.Gpx) > 0 {
		x := gpx.Gpx[0].Trk[0].Trkseg
		for i := 0; i < len(x); i++ {
			sql1 := "INSERT INTO tracks (timedate,latitude,langtitude,altitude) VALUES('" + convertUtcToDb(x[i].Time) + "','" + x[i].Lat + "','" + x[i].Lon + "','" + x[i].Ele + "')"
			_, err := db.Mdb.Exec(sql1)
			if err != nil {
				panic(err.Error())
			}
			if i == 0 {
				startTime, _ = time.Parse("2006-01-02T15:04:05.000Z", x[i].Time)
			}
			if i == (len(x) - 1) {
				finishTime, _ = time.Parse("2006-01-02T15:04:05.000Z", x[i].Time)
			}
		}
		s := finishTime.Unix() - startTime.Unix()
		fmt.Println(utils.SecToTime(int(s)))

		fmt.Println("GPX vloženo")
	} else {
		var tcd Tcd
		xml.Unmarshal([]byte(xmlFile), &tcd)
		x := tcd.Tcd[0].Activities[0].Activity[0].Lap[0].Track
		for i := 0; i < len(x); i++ {
			sql1 := "INSERT INTO tracks (timedate,latitude,langtitude,altitude) VALUES('" + convertUtcToDb(x[i].Time) + "','" + x[i].Position.LatitudeDegrees + "','" + x[i].Position.LongitudeDegrees + "','" + x[i].AltitudeMeters + "')"
			_, err := db.Mdb.Exec(sql1)
			if err != nil {
				panic(err.Error())
			}
			if i == 0 {
				startTime, _ = time.Parse("2006-01-02T15:04:05.000Z", x[i].Time)
			}
			if i == (len(x) - 1) {
				finishTime, _ = time.Parse("2006-01-02T15:04:05.000Z", x[i].Time)
			}

		}

		y := finishTime.Unix() - startTime.Unix()
		fmt.Println(utils.SecToTime(int(y)))

	}

}

func fileFromForm(w http.ResponseWriter, r *http.Request) string {
	r.ParseMultipartForm(32 << 20) // limit your max input length!
	file, header, err := r.FormFile("gpxFile")
	if err != nil {
		panic(err)
	}

	/*
	* jen test pro mé účely jak by to šlo bez bufferu, vynechak jsem tu error
	* bytestring, _ := ioutil.ReadAll(file)
	* fmt.Println(string(bytestring))
	*
	**/

	name := strings.Split(header.Filename, ".") // jméno souboru z Header
	fmt.Println(name)                           //
	var buf bytes.Buffer                        // nadefinujeme buffer pro pole bytů (alespoň tak to chápu já)
	io.Copy(&buf, file)                         // zkopírujeme obsah souboru do bufferu
	fileContent := buf.String()                 // převedeme to na string
	buf.Reset()                                 // vyprázníme buffer
	return fileContent
}

func Filetesty(w http.ResponseWriter, r *http.Request) {
	//neco := fileFromForm(w, r)
	parseXML(fileFromForm(w, r))
}

func convertUtcToDb(Utc string) string {
	t, err := time.Parse("2006-01-02T15:04:05.000Z", Utc)
	if err != nil {
		panic(err)
	}
	return t.Format("2006-01-02 15:04:05")
}
