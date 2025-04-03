FROM golang:latest as builder
ARG CGO_ENABLED=0
WORKDIR /app

COPY go.mod ./
RUN go mod tidy
COPY . .

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g ./cmd/main.go -o ./docs

RUN go build ./cmd/main.go

FROM scratch
WORKDIR /bin
COPY --from=builder /app/main /bin
EXPOSE 8080
ENTRYPOINT ["/bin/main"]