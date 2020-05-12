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
