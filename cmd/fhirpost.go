package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"text/template"

	"github.com/E-Health/goscar"
	"github.com/E-Health/goscar-export/oscutil"
	"github.com/joho/godotenv"
)

type Config struct {
	mySettings oscutil.Settings
}

// New returns a new Config struct
func New() *Config {
	return &Config{
		mySettings: oscutil.Settings{
			GOSCAR_LOCATION:          getEnv("GOSCAR_LOCATION", "MyClinic"),
			USER_NAME:                getEnv("USER_NAME", "MyName"),
			GOSCAR_INPUT_FILE:        getEnv("GOSCAR_INPUT_FILE", "data.csv"),
			GOSCAR_OUTPUT_FILE:       getEnv("GOSCAR_OUTPUT_FILE", "data.json"),
			GOSCAR_SYSTEM:            getEnv("GOSCAR_SYSTEM", "http://canehealth.com/goscar"),
			FHIR_SERVER:              getEnv("FHIR_SERVER", "http://localhost:3001"),
			GOSCAR_ID_SEPARATOR:      getEnv("GOSCAR_ID_SEPARATOR", "-"),
			GOSCAR_URN:               getEnv("GOSCAR_URN", "urn:uuid:"),
			GOSCAR_FORM_NAME:         getEnv("GOSCAR_FORM_NAME", "MyForm"),
			GOSCAR_SYSTEM_ENTRY:      getEnv("GOSCAR_SYSTEM_ENTRY", "http://canehealth.com/goscar/entry"),
			GOSCAR_SYSTEM_TIMESTAMP:  getEnv("GOSCAR_SYSTEM_TIMESTAMP", "http://canehealth.com/goscar/timestamp"),
			GOSCAR_SYSTEM_CLINIC:     getEnv("GOSCAR_SYSTEM_CLINIC", "http://canehealth.com/goscar/clinic"),
			GOSCAR_SYSTEM_VOCABULARY: getEnv("GOSCAR_SYSTEM_VOCABULARY", "SNOMED-CT"),
		},
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func main() {
	settings := oscutil.DefaultSettings()
	err := godotenv.Load()
	conf := New() // config
	settings = conf.mySettings

	// Commandline flags
	filePtr := flag.String("file", "data.csv", "The csv file to process")
	flag.Parse()

	usage := `

		Usage:

			fhirpost -file=data.csv

	`

	if err != nil {
		log.Fatal("Error loading .env file")
	}
	if *filePtr != "" {
		r, err := os.Open(*filePtr)
		if err != nil {
			fmt.Print(usage)
			os.Exit(1)
		}
		settings.GOSCAR_INPUT_FILE = *filePtr
		csvMap := goscar.CSVToMap(r)
		b, err := json.Marshal(oscutil.MapToFHIR(csvMap, settings))
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
func postResource(jsonStr string, settings oscutil.Settings) {
	url := settings.FHIR_SERVER
	if settings.GOSCAR_OUTPUT_FILE != "" {
		_ = ioutil.WriteFile("data.json", []byte(jsonStr), 0644)
	}
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
