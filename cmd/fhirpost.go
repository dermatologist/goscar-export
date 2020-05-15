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
)

func main() {
	// Commandline flags
	filePtr := flag.String("file", "data.csv", "The csv file to process")
	flag.Parse()

	usage := `

Usage:

fhirpost -file=output.csv

	`

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
		b, err := json.Marshal(oscutil2.MapToFHIR(csvMap))
		if err != nil {
			fmt.Println(err)
			return
		}
		// fmt.Println(string(b))
		postResource(string(b))
	} else {
		fmt.Print(usage)
		os.Exit(1)
	}

}

// postResource : Posts FHIR resource to the API
func postResource(jsonStr string) {
	url := os.Getenv("FHIR_SERVER")

	fmt.Println("URL:>", url)
	fmt.Print("FHIR:> ", jsonStr)
	err := ioutil.WriteFile("data.json", []byte(jsonStr), 0644)
	// var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
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

	fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(body))
}
