# goscar-export : OSCAR EForm Export (csv) to FHIR


## About

This is a simple application to convert a CSV file to a FHIR bundle and post it to a FHIR server in Golang. The OSCAR EMR has an EForm export tool that exports EForms to a CSV file that can be downloaded. This tool can load that CSV file to a FHIR server for consolidated analysis. 

## Use Cases

This is useful for family practice groups with multiple OSCAR EMR instances. Analysts at each site can use this to send data to a central FHIR server for centralized data analysis and reporting. Public health agencies using OSCAR or similar health information systems can use this to consolidate data collection.

## How to build

First *go get* all dependencies
This package includes three tools (*Go build* them separately from the cmd folder):

* Fhirpost: The application for posting the csv fie to the FHIR server
* Serverfhir: A simple FHIR server for testing. We recommend using PHS DW for production.
* Report: A simple application for descriptive statistics on the csv file

## Format of the CSV file

 Using vocabulary such as SNOMED for field names in the E-Form is very useful for consolidated analysis.

Each record should have: 

* *demographicNo* → The patient ID
* *dateCreated* 
* *efmfid* → The ID of the eform
* *fdid* → The ID of the each form field.

 **(The Eform export csv of OSCAR typically has all these fields and requires no further processing)**

## How to use:

Change the settings in .env. You can compile this for Windows, Mac or Linux. Check the fhirmap.go file and make any desired changes. You should be able to figure out the mapping rules from this file. It reads data.csv file from the same folder. *(Will add commandline options soon)*

## Privacy and security:

This application does not encrypt the data. Use it only in a secure network. 

## Disclaimer:

This is an experimental application. Use it at your own risk.

