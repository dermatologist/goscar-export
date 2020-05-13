package main

import (
	"github.com/E-Health/goscar-export/internal/oscutil"
	"flag"
	"fmt"
	"github.com/E-Health/goscar"
	"github.com/jroimartin/gocui"
	"log"
	"os"
)

//This is how you declare a global variable
var csvMap, csvMapValid []map[string]string
var recordCount int
var dateFrom, dateTo, filePtr *string
var fid *int
var includeAll *bool

func main() {
	// Commandline flags
	dateFrom = flag.String("datefrom", "2016-01-02", "The start date")
	dateTo = flag.String("dateto", "2020-01-02", "The end date")
	fid = flag.Int("fid", 1, "The eform ID")
	filePtr = flag.String("file", "data.csv", "The csv file to process")
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
			fmt.Print(usage)
			os.Exit(1)
		}
		csvMap = goscar.CSVToMap(r)

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

	oscutil.CsvMap = csvMap
	oscutil.CsvMapValid, oscutil.RecordCount = goscar.FindDuplicates(csvMap, *dateFrom, *dateTo, *includeAll)

	g.SetManagerFunc(oscutil.Layout)


	if err := oscutil.Keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}

