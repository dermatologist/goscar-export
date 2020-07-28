https_proxy='socks5://127.0.0.1:9090' go mod tidy

java -jar -Xmx2g org.hl7.fhir.validator.jar data.json -version 3.0

NO_HMR=1 meteor run --port 8085


## MongoDB

* show dbs
* use fhir
* db.getCollectionNames()
* db.observations.find().pretty()
* db.patients.find().pretty()
* db.patients.drop()
* db.patients.count()


* show collections
* show tables

## 

```
https_proxy='socks5://127.0.0.1:9090' go run cmd/fhirpost.go 
https_proxy='socks5://127.0.0.1:9090' go get github.com/stamblerre/gocode
https_proxy='socks5://127.0.0.1:9090' go mod tidy

```

## "but does not contain package"

* Files have been deleted during purge
* remove github folder in GOROOT and go mod tidy again