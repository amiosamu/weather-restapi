FROM golang:alpine as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

FROM golang:alpine as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /weather-api
WORKDIR /weather-api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags migrate -o /bin/weather-api ./cmd/app

FROM scratch
COPY --from=builder /weather-api/config /config
COPY --from=builder /bin/weather-api /weather-api

CMD ["/weather-api"]
