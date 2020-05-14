package oscutil

import (
	"github.com/E-Health/goscar"
	"github.com/samply/golang-fhir-models/fhir-models/fhir"
	"os"
	"strconv"
)

// FhirObservation : Extended Observation as the existing one does not have valueString and valueInteger
type FhirObservation struct {
	fhir.Observation
	ValueString  *string `json:"value-string,omitempty"`
	ValueInteger *int    `json:"value-integer,omitempty"`
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
	var bundleType = fhir.BundleType(fhir.BundleTypeDocument) // The bundle is a document. The first resource is a Composition.

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
		for header, value := range headers {
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
			observation.Id = &observationId
			if !goscar.IsMember(header, toIgnore) {
				// If the value is a number
				if headerStat["num"] > 0 {
					vI, _ := strconv.Atoi(value)
					observation.ValueInteger = &vI
					// Else treat it like a string
				} else {
					observation.ValueString = &value
				}
			}
			reference.Reference = &patientId // Observation refers to the patient
			observation.Subject = &reference
			// Patient
			// bundleEntry.Id = &mySystem
			bundleEntry.Resource, _ = patient.MarshalJSON()
			bundle.Entry = append(bundle.Entry, bundleEntry)
			// Observation
			// bundleEntry.Id = &mySystem
			bundleEntry.Resource, _ = observation.MarshalJSON()
			bundle.Entry = append(bundle.Entry, bundleEntry)
		}

	}

	return bundle
}
