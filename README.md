# Case Study

## Installation

Clone this repo

```
git clone git@github.com:albertsundjaja/frankie.git
```

init vendors

```
cd frankie
go mod vendor
```

## Running the project

### Run with docker

make sure you have docker and docker-compose installed

to spin up the container simply type the command

```
docker-compose up -d
```

### Run without docker

To run this project without docker

```
go run main.go
```

To run Ginkgo testing suite 

```
go test ./..

OR

ginkgo ./...
```

To check coverage on testing

```
go test -cover

OR

ginkgo -r -cover
```

### Accessing the endpoint

Use a REST API testing app such as Postman

the server is running on port `8888`

example request

```
POST http://localhost:8888/isgood
Content-Type: application/json

Body
[
  {
    "checkType": "BIOMETRIC",
    "activityType": "SIGNUP",
    "checkSessionKey": "string1",
    "activityData": [
      {
        "kvpKey": "ip.address",
        "kvpValue": "1.23.45.123",
        "kvpType": "general.string"
      }
    ]
  },
  {
    "checkType": "BIOMETRIC",
    "activityType": "SIGNUP",
    "checkSessionKey": "string2",
    "activityData": [
      {
        "kvpKey": "ip.address4",
        "kvpValue": "1.1.1.1",
        "kvpType": "general.string"
      }
    ]
  }
]

```