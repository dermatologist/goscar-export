package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/samply/golang-fhir-models/fhir-models/fhir"
	"io/ioutil"
	"net/http"
)

func main() {
	patientId := "12345"
	given := []string{"Mickey"}
	family := "Mouse"

	humanName := fhir.HumanName{
		Given:  given,
		Family: &family,
	}

	hn := []fhir.HumanName{}
	hn = append(hn, humanName)
	patient := fhir.Patient{
		Id:   &patientId,
		Name: hn,
	}
	// fmt.Print(patient)
	b, err := json.Marshal(patient)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
	postResource(string(b))
}

// postResource : Posts FHIR resource to the API
func postResource(jsonStr string) {
	url := "http://restapi3.apiary.io/notes"
	fmt.Println("URL:>", url)

	// var jsonStr = []byte(`{"title":"Buy cheese and bread for breakfast."}`)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
