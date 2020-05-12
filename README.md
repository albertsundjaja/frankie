# Case Study

## Installation

Clone this repo

```
git clone git@github.com:albertsundjaja/frankie.git
```

init vendors

```
go mod vendor
```

## Running the project

To run this project

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
