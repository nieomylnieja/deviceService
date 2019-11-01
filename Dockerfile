# Device service app with GO
FROM golang:1.13.3

WORKDIR /go/src/deviceService
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["deviceService"]
