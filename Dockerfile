FROM golang:latest as builder
ARG CGO_ENABLED=0
WORKDIR /app

COPY go.mod ./
RUN go mod tidy
COPY . .

RUN go build ./cmd/main.go

FROM scratch
WORKDIR /bin
COPY --from=builder /app/main /bin
EXPOSE 1337
ENTRYPOINT ["/bin/main"]