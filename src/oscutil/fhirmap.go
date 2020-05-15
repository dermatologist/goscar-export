package oscutil

import (
	"encoding/json"
	"github.com/E-Health/goscar"
	"github.com/E-Health/goscar-export/internal/oscutil"
	"github.com/google/uuid"
	"github.com/samply/golang-fhir-models/fhir-models/fhir"
	"os"
	"regexp"
	"strconv"
	"time"
)

// FhirObservation : Extended Observation as the existing one does not have valueString and valueInteger
type FhirObservation struct {
	Id           string                 `json:"id,omitempty"`
	Identifier   []fhir.Identifier      `json:"identifier,omitempty"`
	Subject      fhir.Reference         `json:"subject,omitempty"`
	ValueString  string                 `json:"valueString,omitempty"`
	ValueInteger int                    `json:"valueInteger,omitempty"`
	ResourceType string                 `json:"resourceType"`
	Status       fhir.ObservationStatus `json:"status"`
	Code         fhir.CodeableConcept   `json:"code"`
}

// MapToFHIR : maps the csvMap to a FHIR bundle.
func MapToFHIR(_csvMapValid []map[string]string) fhir.Bundle {
	var composition = fhir.Composition{}
	var patient = fhir.Patient{}
	var observation = FhirObservation{} // Extended observation: See above for definition
	var i1 = fhir.Identifier{}
	var bundle = fhir.Bundle{}
	var reference = fhir.Reference{}
	var bundleEntry = fhir.BundleEntry{}
	var codableConcept = fhir.CodeableConcept{}
	var bundleType = fhir.BundleType(fhir.BundleTypeDocument) // The bundle is a document. The first resource is a Composition.
	var practitioner = fhir.Practitioner{}
	id := uuid.New()
	dt := time.Now()
	patients := []string{}

	location := os.Getenv("GOSCAR_LOCATION")
	username := os.Getenv("USER_NAME")
	mySeparator := os.Getenv("GOSCAR_ID_SEPARATOR")
	myUrn := os.Getenv("GOSCAR_URN")
	myForm := os.Getenv("GOSCAR_FORM_NAME")
	myVocabulary := os.Getenv("GOSCAR_SYSTEM_VOCABULARY")

	mySystem := os.Getenv("GOSCAR_SYSTEM")
	mySystemEntry := os.Getenv("GOSCAR_SYSTEM_ENTRY")
	mySystemTimestamp := os.Getenv("GOSCAR_SYSTEM_TIMESTAMP")
	mySystemClinic := os.Getenv("GOSCAR_SYSTEM_CLINIC")

	toIgnore := []string{"id", "fdid", "dateCreated", "eform_link", "StaffSig", "SubmitButton", "efmfid"}

	// Single bundle
	myTitle := location + mySeparator + username
	i1.System = &mySystem
	i1.Value = &myTitle
	bundle.Identifier = &i1
	bundle.Type = bundleType // The bundle is a document. The first resource is a Composition.
	//bundleEntry.Id = &mySystem
	bundleTimestamp := dt.UTC().Format("2006-01-02T15:04:05Z")
	bundle.Timestamp = &bundleTimestamp

	// Create a practitioner who is the author and the subject of composition
	practitionerId := location + mySeparator + username
	practitionerRefId := "Practitioner/" + practitionerId
	practitioner.Id = &practitionerId

	// Single composition
	composition.Identifier = &i1
	composition.Status = fhir.CompositionStatus(fhir.CompositionStatusFinal) // Required

	// Random UUID for composition
	compositionId := id.String()
	composition.Id = &compositionId
	composition.Title = myTitle + mySeparator + myForm
	composition.Date = dt.Format("2006-01-02")

	// Set author as author and subject
	reference.Reference = &practitionerRefId
	composition.Author = append(composition.Author, reference)
	composition.Subject = &reference
	codableText := myUrn + "E-Form" + mySeparator + myForm
	codableConcept.Text = &codableText
	composition.Type = codableConcept

	// Add composition
	bundleEntry.Resource, _ = composition.MarshalJSON()
	myCompositionEntry := myUrn + compositionId
	bundleEntry.FullUrl = &myCompositionEntry
	bundle.Entry = append(bundle.Entry, bundleEntry)

	// Add practitioner
	bundleEntry.Resource, _ = practitioner.MarshalJSON()
	myPractitionerEntry := myUrn + practitionerId
	bundleEntry.FullUrl = &myPractitionerEntry
	bundle.Entry = append(bundle.Entry, bundleEntry)

	// Get headers from the first row
	headers := _csvMapValid[0]
	for _, record := range _csvMapValid {

		// Each record has a patient (ID is unique for location)
		patientId := location + mySeparator + record["demographicNo"]
		refPatientId := "Patient/" + patientId
		identifier := fhir.Identifier{}
		identifier.System = &mySystem
		identifier.Value = &patientId
		_identifier := []fhir.Identifier{}
		_identifier = append(_identifier, identifier)
		patient.Identifier = _identifier
		patient.Id = &patientId // Needed for Reference

		// Each value is an observation
		for header, myval := range headers {
			// Function call to get the type of header -> number or string
			headerStat := goscar.GetStats(header, oscutil.RecordCount, _csvMapValid)

			// Add the form field identifier
			identifier := fhir.Identifier{}
			toAdd1 := header
			identifier.System = &mySystemEntry
			identifier.Value = &toAdd1
			_identifier = append(_identifier, identifier)

			// Add the timestamp identifier
			identifier = fhir.Identifier{}
			toAdd2 := record["dateCreated"]
			identifier.System = &mySystemTimestamp
			identifier.Value = &toAdd2
			_identifier = append(_identifier, identifier)

			// Add the form field identifier
			identifier = fhir.Identifier{}
			toAdd3 := location
			identifier.System = &mySystemClinic
			identifier.Value = &toAdd3
			_identifier = append(_identifier, identifier)

			observation.Identifier = _identifier

			// Create a unique ID for observation to be added to key to generate the final ID
			observationId := location + mySeparator +
				record["efmfid"] + mySeparator +
				record["fdid"] + mySeparator +
				record["dateCreated"] + mySeparator + header
			reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
			observationId = reg.ReplaceAllString(observationId, "")
			observation.Id = observationId
			if !goscar.IsMember(header, toIgnore) {
				if headerStat["num"] > 0 {
					vI, _ := strconv.Atoi(myval)
					observation.ValueInteger = vI
					// observation.ValueString = ""
					// Else treat it like a string
				} else {
					observation.ValueString = myval
					// observation.ValueInteger = 0
				}
			}
			reference.Reference = &refPatientId // Observation refers to the /Patient/id
			observation.Subject = reference
			observation.ResourceType = "Observation"
			observation.Status = fhir.ObservationStatus(fhir.ObservationStatusRegistered) // Required
			codableConcept.Id = &myVocabulary
			codableConcept.Text = &header
			observation.Code = codableConcept
			// @TODO To switch after debug
			// Unique ID
			id := uuid.New()
			myUuid := myUrn + id.String()
			myPatientEntry := myUuid + "_patient"
			myObservationEntry := myUuid + "_observation"
			// Add patient if not added already
			if !goscar.IsMember(patientId, patients) {
				bundleEntry.Resource, _ = patient.MarshalJSON()
				bundleEntry.FullUrl = &myPatientEntry
				bundle.Entry = append(bundle.Entry, bundleEntry)
				patients = append(patients, patientId)
			}
			// Observation
			// bundleEntry.Id = &mySystem
			bundleEntry.Resource, _ = json.Marshal(observation)
			bundleEntry.FullUrl = &myObservationEntry
			bundle.Entry = append(bundle.Entry, bundleEntry)
			observation = FhirObservation{} // Clear values

		}

	}

	return bundle
}
