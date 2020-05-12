package main

import (
	"github.com/E-Health/goscar-export/internal/oscutil"
	"flag"
	"fmt"
	"github.com/E-Health/goscar"
	"github.com/jroimartin/gocui"
	"log"
	"os"
	"time"
)

//This is how you declare a global variable
var csvMap, csvMapValid []map[string]string
var recordCount int
var sshHost, sshUser, sshPass, dbUser, dbPass, dbHost, dbName, dateFrom, dateTo, filePtr *string
var sshPort, fid *int
var includeAll *bool

func main() {
	// Commandline flags
	sshHost = flag.String("sshhost", "", "The SSH host")
	sshPort = flag.Int("sshport", 22, "The port number")
	sshUser = flag.String("sshuser", "ssh-user", "ssh user")
	sshPass = flag.String("sshpass", "ssh-pass", "SSH Password")
	dbUser = flag.String("dbuser", "dbuser", "The db user")
	dbPass = flag.String("dbpass", "dbpass", "The db password")
	dbHost = flag.String("dbhost", "localhost:3306", "The db host")
	dbName = flag.String("dbname", "oscar", "The database name")
	dateFrom = flag.String("datefrom", "oscar", "The start date")
	dateTo = flag.String("dateto", "oscar", "The end date")
	fid = flag.Int("fid", 1, "The eform ID")
	filePtr = flag.String("file", "", "The csv file to process")
	includeAll = flag.Bool("include", false, "Include all records")
	flag.Parse()

	usage := `

Usage:

oscar_helper -file=output.csv

oscar_helper -sshhost=xxx -sshport=22 -sshuser=xxx -sshpass=xxx -dbuser=xxx -dbpass=xxx -dbname=xxx -dbhost=localhost -datefrom=YYYY-MM-DD -dateto=YYYY-MM-DD -fid=1 -include

	`
	if *filePtr != "" {
		r, err := os.Open(*filePtr)
		if err != nil {
			log.Panicln(err)
		}
		csvMap = goscar.CSVToMap(r)

	} else if *sshHost != "" {
		r, err := oscutil.MysqlConnect()
		if err != nil {
			log.Panicln(err)
		}

		csvMap = goscar.MysqlToMap(r)
	} else {
		fmt.Print(usage)
		os.Exit(1)
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true

	g.SetManagerFunc(oscutil.Layout)

	findDuplicates(csvMap)

	if err := oscutil.Keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}

// TODO change to the one from goscar
func findDuplicates(csvMap []map[string]string) {
	var latest bool
	var included bool
	var demographicNo []string
	for _, v := range csvMap {
		latest = false
		included = true
		for k2, v2 := range v {
			if k2 == "eft_latest" && v2 == "1" {
				latest = true
			}
			if k2 == "dateCreated" {
				dateCreated, _ := time.Parse("2006-01-02", v2)
				_dateFrom, _ := time.Parse("2006-01-02", *dateFrom)
				_dateTo, _ := time.Parse("2006-01-02", *dateTo)
				if len(*dateFrom) > 0 && len(*dateTo) > 0 && !goscar.InTimeSpan(_dateFrom, _dateTo, dateCreated) {
					included = false
				}
			}
			if k2 == "demographic_no" {
				if !goscar.IsMember(v2, demographicNo){
					demographicNo = append(demographicNo, v2)
					latest = true
				}
			}
			if *includeAll {
				latest = true
			}
		}
		if latest && !included {
			csvMapValid = append(csvMapValid, v)
			recordCount++
		}
	}
}