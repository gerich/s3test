# upload text file
POST http://{{host}}:8080/api/v1/upload
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoiZm9vIn0.6ee7NeQqcFzmPymq-5N3h52P98_JCOKwsrwye_SVfZA
[MultipartFormData]
file: file,my-text.txt;
HTTP 200
[Asserts]
header "Content-Type" contains "application/json; charset=utf-8"
jsonpath "$.hash" == "4e398d35fb30b30b2228ce5c5e2f4f34"


# upload image
POST http://{{host}}:8080/api/v1/upload
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoiZm9vIn0.6ee7NeQqcFzmPymq-5N3h52P98_JCOKwsrwye_SVfZA
[MultipartFormData]
file: file,my-image.png;
HTTP 200
[Asserts]
header "Content-Type" contains "application/json; charset=utf-8"
jsonpath "$.hash" == "1e02a85bee78984d20942bf213f4c461"


# list files
GET http://{{host}}:8080/api/v1/list
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoiZm9vIn0.6ee7NeQqcFzmPymq-5N3h52P98_JCOKwsrwye_SVfZA
HTTP 200
[Asserts]
header "Content-Type" contains "application/json; charset=utf-8"
jsonpath "$.user" == "foo"
jsonpath "$.files[?(@.name == 'my-image.png')].name" exists
jsonpath "$.files[?(@.name == 'my-text.txt')].name" exists

# download image
GET http://{{host}}:8080/api/v1/download/my-image.png
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoiZm9vIn0.6ee7NeQqcFzmPymq-5N3h52P98_JCOKwsrwye_SVfZA
HTTP 200
[Asserts]
md5 == hex,1e02a85bee78984d20942bf213f4c461;


# download text
GET http://{{host}}:8080/api/v1/download/my-text.txt
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoiZm9vIn0.6ee7NeQqcFzmPymq-5N3h52P98_JCOKwsrwye_SVfZA
HTTP 200
[Asserts]
md5 == hex,4e398d35fb30b30b2228ce5c5e2f4f34;
