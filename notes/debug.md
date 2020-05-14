Error @ Bundle.entry[0].resource.ofType(Composition).date (line 1, col316) : @value cannot be empty
  Error @ Bundle.entry[0].resource.ofType(Composition).title (line 1, col361) : @value cannot be empty
  Error @ Bundle.entry[1].resource.author[1].reference (line 1, col327) : Can't find "MyName" in the bundle (Composition.author)
  Error @ Bundle.entry[1].resource.subject.reference (line 1, col305) : Can't find "MyName" in the bundle (Composition.subject)
  Error @ Bundle.entry[0].resource.type (line 1, col273) : Object must have some content
  Error @ Bundle.entry[2].resource.code (line 1, col757) : Object must have some content
  Error @ Bundle.entry[4].resource.code (line 1, col1130) : Object must have some content
  Error @ Bundle.entry[6].resource.code (line 1, col1529) : Object must have some content
  Error @ Bundle.entry[2].resource.ofType(Observation).id (line 1, col585) : id value "MyClinic-116-200-08:48.0-fdid" is not valid
  Error @ Bundle.entry[4].resource.ofType(Observation).id (line 1, col955) : id value "MyClinic-116-201-08:48.0-cocaine" is not valid
  Error @ Bundle.entry[6].resource.ofType(Observation).id (line 1, col1341) : id value "MyClinic-116-203-09:25.0-other_drug_of_choice" is not valid
  Error @ Bundle : A document must have a date {0}
  Error @ Bundle.entry[0].resource.ofType(Composition).subject : Relative Reference appears inside Bundle whose entry is missing a fullUrl
  Error @ Bundle.entry[0].resource.ofType(Composition).author[0] : Relative Reference appears inside Bundle whose entry is missing a fullUrl
  Error @ Bundle.entry[2].resource.ofType(Observation).subject : Relative Reference appears inside Bundle whose entry is missing a fullUrl
  Error @ Bundle.entry[4].resource.ofType(Observation).subject : Relative Reference appears inside Bundle whose entry is missing a fullUrl
  Error @ Bundle.entry[6].resource.ofType(Observation).subject : Relative Reference appears inside Bundle whose entry is missing a fullUrl
  Error @ Bundle.entry[0] (line 1, col116) : Bundle entry missing fullUrl
  Error @ Bundle.entry[2] (line 1, col536) : Bundle entry missing fullUrl
  Error @ Bundle.entry[4] (line 1, col903) : Bundle entry missing fullUrl
  Error @ Bundle.entry[6] (line 1, col1276) : Bundle entry missing fullUrl
  Error @ Bundle.entry[0].resource.ofType(Composition).date (line 1, col316) : ele-1: All FHIR elements must have a @value or children [hasValue() | (children().count() > id.count())]
  Error @ Bundle.entry[0].resource.ofType(Composition).title (line 1, col361) : ele-1: All FHIR elements must have a @value or children [hasValue() | (children().count() > id.count())]


```
{
   "identifier":{
      "system":"http://canehealth.com/goscar",
      "value":"other_drug_of_choice"
   },
   "type":"document",
   "entry":[
      {
         "resource":{
            "id":"a9575356-90c1-436a-b2a4-602ba2d6ca8c",
            "identifier":{
               "system":"http://canehealth.com/goscar",
               "value":"MyClinic"
            },
            "status":"final",
            "type":{

            },
            "subject":{
               "reference":"MyName"
            },
            "date":"",
            "author":[
               {
                  "reference":"MyName"
               }
            ],
            "title":"",
            "resourceType":"Composition"
         }
      },
      {
         "resource":{
            "id":"MyClinic-100",
            "identifier":[
               {
                  "system":"http://canehealth.com/goscar",
                  "value":"MyClinic-MyName"
               }
            ],
            "resourceType":"Patient"
         }
      },
      {
         "resource":{
            "id":"MyClinic-116-200-08:48.0-fdid",
            "identifier":[
               {
                  "system":"http://canehealth.com/goscar",
                  "value":"fdid"
               }
            ],
            "subject":{
               "reference":"MyClinic-100"
            },
            "resourceType":"Observation",
            "status":"registered",
            "code":{

            }
         }
      },
      {
         "resource":{
            "id":"MyClinic-100",
            "identifier":[
               {
                  "system":"http://canehealth.com/goscar",
                  "value":"MyClinic-MyName"
               }
            ],
            "resourceType":"Patient"
         }
      },
      {
         "resource":{
            "id":"MyClinic-116-201-08:48.0-cocaine",
            "identifier":[
               {
                  "system":"http://canehealth.com/goscar",
                  "value":"cocaine"
               }
            ],
            "subject":{
               "reference":"MyClinic-100"
            },
            "resourceType":"Observation",
            "status":"registered",
            "code":{

            }
         }
      },
      {
         "resource":{
            "id":"MyClinic-101",
            "identifier":[
               {
                  "system":"http://canehealth.com/goscar",
                  "value":"MyClinic-MyName"
               }
            ],
            "resourceType":"Patient"
         }
      },
      {
         "resource":{
            "id":"MyClinic-116-203-09:25.0-other_drug_of_choice",
            "identifier":[
               {
                  "system":"http://canehealth.com/goscar",
                  "value":"other_drug_of_choice"
               }
            ],
            "subject":{
               "reference":"MyClinic-101"
            },
            "resourceType":"Observation",
            "status":"registered",
            "code":{

            }
         }
      }
   ],
   "resourceType":"Bundle"
}
```