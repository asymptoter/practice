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
COPY --from=builder /go/src/github.com/asymptoter/practice-backend/app/server/server.crt /app/server/server.crt
COPY --from=builder /go/src/github.com/asymptoter/practice-backend/app/server/server.key /app/server/server.key
COPY --from=builder /go/src/github.com/asymptoter/practice-backend/config /config
COPY --from=builder /go/src/github.com/asymptoter/practice-backend/creds.json /creds.json
#COPY --from=builder /go/src/github.com/asymptoter/practice-backend/scripts/wait-for-it.sh /scripts/wait-for-it.sh
WORKDIR /app/server
CMD ["./main"]
EXPOSE 80
