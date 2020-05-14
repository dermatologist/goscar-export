https_proxy='socks5://127.0.0.1:9090' go mod tidy

java -jar -Xmx2g org.hl7.fhir.validator.jar data.json -version 3.0