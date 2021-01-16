FROM golang:alpine AS build-env

WORKDIR /app

RUN mkdir files

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /go/bin/main

####################
### Second stage ###
####################

FROM alpine:3.7

EXPOSE 8080

WORKDIR /app

RUN apk add python3 py3-pip
RUN python3 -m pip install semgrep

COPY --from=build-env /app /app
COPY --from=build-env /go/bin/main /go/bin/main

CMD ["/go/bin/main"]