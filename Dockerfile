FROM golang:latest AS builder
ENV GO111MODULE on
WORKDIR /go/src/github.com/asymptoter/practice-backend
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
WORKDIR ./app/server
RUN CGO_ENABLED=0 GOOS=linux go build -o main

FROM alpine:latest
COPY --from=builder /go/src/github.com/asymptoter/practice-backend/app/server/main /app/server/main
COPY --from=builder /go/src/github.com/asymptoter/practice-backend/config/config_ci.yml /config/config.yml
COPY --from=builder /go/src/github.com/asymptoter/practice-backend/creds.json /creds.json
WORKDIR /app/server
CMD ["./main"]
EXPOSE 80
