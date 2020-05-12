FROM golang:1.14.2-buster

## install git
RUN apt-get update && apt-get install -y git

## add key to known hosts
RUN mkdir /root/.ssh
RUN touch /root/.ssh/known_hosts
RUN ssh-keyscan github.com >> /root/.ssh/known_hosts

## clone the repo
RUN git clone https://github.com/albertsundjaja/frankie.git

## init vendors
WORKDIR ./frankie
RUN go mod vendor

EXPOSE 8888

## run the server
CMD ["go", "run", "main.go"]