package oscutil

import (
	"encoding/json"
	"fmt"
	"github.com/E-Health/goscar"
	"github.com/samply/golang-fhir-models/fhir-models/fhir"
	"os"
	"strconv"
	"github.com/google/uuid"
)

// FhirObservation : Extended Observation as the existing one does not have valueString and valueInteger
type FhirObservation struct {
	Id           string            `json:"id,omitempty"`
	Identifier   []fhir.Identifier `json:"identifier,omitempty"`
	Subject      fhir.Reference    `json:"subject,omitempty"`
	ValueString  string            `json:"valueString,omitempty"`
	ValueInteger int               `json:"valueInteger,omitempty"`
	ResourceType string            `json:"resourceType"`
	Status       fhir.ObservationStatus `json:"status"`
	Code         fhir.CodeableConcept             `json:"code"`
}

// MapToFHIR : maps the csvMap to a FHIR bundle.
func MapToFHIR(_csvMapValid []map[string]string) fhir.Bundle {
	var composition = fhir.Composition{}
	var patient = fhir.Patient{}
	var observation = FhirObservation{} // Extended observation: See above for definition
	var identifier = fhir.Identifier{}
	var bundle = fhir.Bundle{}
	var reference = fhir.Reference{}
	var bundleEntry = fhir.BundleEntry{}
	var codableConcept = fhir.CodeableConcept{}
	var bundleType = fhir.BundleType(fhir.BundleTypeDocument) // The bundle is a document. The first resource is a Composition.
	id := uuid.New()

	location := os.Getenv("GOSCAR_LOCATION")
	username := os.Getenv("USER_NAME")
	mySystem := os.Getenv("GOSCAR_SYSTEM")
	mySeparator := os.Getenv("GOSCAR_ID_SEPARATOR")
	toIgnore := []string{"id", "fdid", "dateCreated", "eform_link", "StaffSig", "SubmitButton", "efmfid"}

	// Single composition
	myValue := location + mySeparator + username
	identifier.System = &mySystem
	identifier.Value = &myValue
	composition.Identifier = &identifier
	composition.Status = fhir.CompositionStatus(fhir.CompositionStatusFinal) // Required
	compositionId := id.String()
	composition.Id = &compositionId
	reference.Reference = &username
	composition.Author = append(composition.Author, reference)
	reference.Reference = &username // Observation refers to the patient
	composition.Subject = &reference

	// Single bundle
	identifier.System = &mySystem
	identifier.Value = &location
	bundle.Identifier = &identifier
	bundle.Type = bundleType // The bundle is a document. The first resource is a Composition.
	//bundleEntry.Id = &mySystem
	bundleEntry.Resource, _ = composition.MarshalJSON()
	bundle.Entry = append(bundle.Entry, bundleEntry)

	// Get headers from the first row
	headers := _csvMapValid[0]
	for _, record := range _csvMapValid {

		// Each record has a patient (ID is unique for location)
		patientId := location + mySeparator + record["demographicNo"]
		identifier.System = &mySystem
		identifier.Value = &myValue
		_identifier := []fhir.Identifier{}
		_identifier = append(_identifier, identifier)
		patient.Identifier = _identifier
		patient.Id = &patientId // Needed for Reference

		// Each value is an observation
		for header, myval := range headers {
			// Function call to get the type of header -> number or string
			headerStat := goscar.GetStats(header, RecordCount, _csvMapValid)
			identifier.System = &mySystem
			identifier.Value = &header
			_identifier := []fhir.Identifier{}
			_identifier = append(_identifier, identifier)
			observation.Identifier = _identifier
			// Create a unique ID for observation to be added to key to generate the final ID
			observationId := location + mySeparator +
				record["efmfid"] + mySeparator +
				record["fdid"] + mySeparator +
				record["dateCreated"] + mySeparator + header
			observation.Id = observationId
			if !goscar.IsMember(header, toIgnore) {
				if headerStat["num"] > 0 {
					vI, _ := strconv.Atoi(myval)
					observation.ValueInteger = vI
					observation.ValueString = ""
					// Else treat it like a string
				} else {
					observation.ValueString = myval
					observation.ValueInteger = 0
				}
			}
			reference.Reference = &patientId // Observation refers to the patient
			observation.Subject = reference
			observation.ResourceType = "Observation"
			observation.Status = fhir.ObservationStatus(fhir.ObservationStatusRegistered) // Required
			observation.Code = codableConcept
			// @TODO To switch after debug

			// // Patient
			// // bundleEntry.Id = &mySystem
			// bundleEntry.Resource, _ = patient.MarshalJSON()
			// bundle.Entry = append(bundle.Entry, bundleEntry)
			// // Observation
			// // bundleEntry.Id = &mySystem
			// // fmt.Println(*observation.ValueString)
			// bundleEntry.Resource, _ = json.Marshal(observation)
			// fmt.Print(string(bundleEntry.Resource))
			// bundle.Entry = append(bundle.Entry, bundleEntry)
		}
			// Patient
			// bundleEntry.Id = &mySystem
			bundleEntry.Resource, _ = patient.MarshalJSON()
			bundle.Entry = append(bundle.Entry, bundleEntry)
			// Observation
			// bundleEntry.Id = &mySystem
			// fmt.Println(*observation.ValueString)
			bundleEntry.Resource, _ = json.Marshal(observation)
			fmt.Print(string(bundleEntry.Resource))
			bundle.Entry = append(bundle.Entry, bundleEntry)
	}

	return bundle
}
