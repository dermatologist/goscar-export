# goscar-export : OSCAR EMR EForm Export (csv) to FHIR

[![forthebadge](https://forthebadge.com/images/badges/built-with-love.svg) **for OSCAR EMR** ](https://oscar-emr.com/)

## About

This is a simple application to convert a CSV file to a FHIR bundle and post it to a FHIR server in Golang. The OSCAR EMR has an EForm export tool that exports EForms to a CSV file that can be downloaded. This tool can load that CSV file to a FHIR server for consolidated analysis.
**This tool can be used with any CSV, if columns specified below (CSV format section) are present.** 

## Use Cases

This is useful for family practice groups with multiple OSCAR EMR instances. Analysts at each site can use this to send data to a central FHIR server for centralized data analysis and reporting. Public health agencies using OSCAR or similar health information systems can use this to consolidate data collection.

## How to build

First *go get* all dependencies
This package includes three tools (*Go build* them separately from the cmd folder):

* Fhirpost: The application for posting the csv fie to the FHIR server
* Serverfhir: A simple FHIR server for testing (requires mongodb). We recommend using [PHIS-DW](https://github.com/E-Health/fhir-server-phis-dw) for production.
* Report: A simple application for descriptive statistics on the csv file

## Format of the CSV file

 Using vocabulary such as SNOMED for field names in the E-Form is very useful for consolidated analysis.

Each record should have: 

* *demographicNo* → The patient ID
* *dateCreated* 
* *efmfid* → The ID of the eform
* *fdid* → The ID of the each form field.

 **(The Eform export csv of OSCAR typically has all these fields and requires no further processing)**

## Mapping
* Bundle with unique patients. All columns mapped to observations.
* Submitter mapped to Practitioner.
* Document type *bundle* with *composition* as the first entry
* Unique fullUrls are generated.
* PatientID is location + demographicNo
* Budle of 1 composition, 1 practitioner, 1 or more patients, and many observations
* Validates with R4 schema

## How to use:

* Change the settings in .env 
* You can compile this for Windows, Mac or Linux. Check the fhirmap.go file and make any desired changes. You should be able to figure out the mapping rules from this file. 
* It reads data.csv file from the same folder by default. *(can be specified by the -file commandline argument: fhirpost -file=data.csv)*
* Start mongodb and run server and fhirpost in separate windows for testing.
* On windows, you can just double-click executables to run. (Closes automatically after run)

## [Import](https://e-health.github.io/goscar-export/pkg/github.com/E-Health/goscar-export/oscutil/index.html)

## Privacy and security:

This application does not encrypt the data. Use it only in a secure network. 

## Disclaimer:

This is an experimental application. Use it at your own risk.
Pull requests welcome. Refer to CONTRIBUTING.md

## Contributors
* [Bell Eapen](http://nuchange.ca) | [![Twitter Follow](https://img.shields.io/twitter/follow/beapen?style=social)](https://twitter.com/beapen)

