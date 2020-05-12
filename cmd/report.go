package main

import (
	"flag"
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/E-Health/goscar"
	"log"
	"os"
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
		r, err := mysqlConnect()
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

	g.SetManagerFunc(layout)

	goscar.FindDuplicates(csvMap)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}