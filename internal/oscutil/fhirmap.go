package oscutil

import (
	"github.com/E-Health/goscar"
	"github.com/samply/golang-fhir-models/fhir-models/fhir"
	"strconv"
)

// Extended as the existing one does not have valueString and valueInteger
type FhirObservation struct {
	fhir.Observation
	valueString  *string `json:"value-string,omitempty"`
	valueInteger *int    `json:"value-integer,omitempty"`
}

func maptofhir(csvMapValid []map[string]string) fhir.Bundle {
	var composition = fhir.Composition{}
	var patient = fhir.Patient{}
	var observation = FhirObservation{} // Extended observation: See above for definition
	var identifier = fhir.Identifier{}
	var bundle = fhir.Bundle{}
	var reference = fhir.Reference{}
	var bundleEntry = fhir.BundleEntry{}
	var bundleType = fhir.BundleType(fhir.BundleTypeDocument) // The bundle is a document. The first resource is a Composition.

	mySystem := "My System" // Get this from .env

	toIgnore := []string{"id", "fdid", "dateCreated", "eform_link", "StaffSig", "SubmitButton", "efmfid"}

	// Single composition
	myValue := "Some Value"
	identifier.System = &mySystem
	identifier.Value = &myValue
	composition.Identifier = &identifier

	// Single bundle
	myValue = "Some Value"
	identifier.System = &mySystem
	identifier.Value = &myValue
	bundle.Identifier = &identifier
	bundle.Type = bundleType // The bundle is a document. The first resource is a Composition.
	bundleEntry.Id = &mySystem
	bundleEntry.Resource, _ = composition.MarshalJSON()
	bundle.Entry = append(bundle.Entry, bundleEntry)

	// Get headers from the first row
	headers := CsvMap[0]
	for _, record := range csvMapValid {
		myValue := record["demographicNo"]
		identifier.System = &mySystem
		identifier.Value = &myValue
		_identifier := []fhir.Identifier{}
		_identifier = append(_identifier, identifier)
		patient.Identifier = _identifier
		for header, value := range headers {
			// Function call to get the type of header -> number or string
			headerStat := goscar.GetStats(header, RecordCount, CsvMapValid)
			identifier.System = &mySystem
			identifier.Value = &header
			_identifier := []fhir.Identifier{}
			_identifier = append(_identifier, identifier)
			observation.Identifier = _identifier
			if !goscar.IsMember(header, toIgnore) {
				// If the value is a number
				if headerStat["num"] > 0 {
					vI, _ := strconv.Atoi(value)
					observation.valueInteger = &vI
					// Else treat it like a string
				} else {
					observation.valueString = &value
				}
			}
			reference.Reference = &myValue
			observation.Subject = &reference
		}
		// Patient
		bundleEntry.Id = &mySystem
		bundleEntry.Resource, _ = patient.MarshalJSON()
		bundle.Entry = append(bundle.Entry, bundleEntry)
		// Observation
		bundleEntry.Id = &mySystem
		bundleEntry.Resource, _ = observation.MarshalJSON()
		bundle.Entry = append(bundle.Entry, bundleEntry)
	}

	return bundle
}
