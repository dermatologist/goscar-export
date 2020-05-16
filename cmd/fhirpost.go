package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/E-Health/goscar"
	oscutil2 "github.com/E-Health/goscar-export/src/oscutil"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"text/template"
)

func main() {
	// Commandline flags
	filePtr := flag.String("file", "data.csv", "The csv file to process")
	flag.Parse()

	usage := `

Usage:

fhirpost -file=output.csv

	`

	settings := oscutil2.DefaultSettings()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if *filePtr != "" {
		r, err := os.Open(*filePtr)
		if err != nil {
			fmt.Print(usage)
			os.Exit(1)
		}
		csvMap := goscar.CSVToMap(r)
		b, err := json.Marshal(oscutil2.MapToFHIR(csvMap, settings))
		if err != nil {
			fmt.Println(err)
			return
		}
		// fmt.Println(string(b))
		postResource(string(b), settings)
	} else {
		fmt.Print(usage)
		os.Exit(1)
	}

}

// postResource : Posts FHIR resource to the API
func postResource(jsonStr string, settings oscutil2.Settings) {
	url := settings.FHIR_SERVER
	//fmt.Println("URL:>", settings.FHIR_SERVER)
	//fmt.Print("FHIR:> ", jsonStr)
	err := ioutil.WriteFile("data.json", []byte(jsonStr), 0644)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonStr)))
	if err != nil {
		panic(err)
	} else {
		req.Header.Set("X-Custom-Header", "myvalue")
		req.Header.Set("Content-Type", "application/json")
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Files are provided as a slice of strings.
	paths := []string{
		"./template/post.tmpl",
	}
	name := path.Base(paths[0])
	t := template.Must(template.New(name).ParseFiles(paths...))
	err = t.Execute(os.Stdout, settings)
	if err != nil {
		panic(err)
	}

	fmt.Println("response Status:", resp.Status)

	//fmt.Println("response Headers:", resp.Header)
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(body))
}
