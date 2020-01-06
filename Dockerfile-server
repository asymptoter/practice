FROM golang:latest AS builder
ENV GO111MODULE on
WORKDIR /src/github.com/asymptoter/geochallenge-backend
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
WORKDIR ./app/server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main .

FROM scratch
COPY --from=builder /src/github.com/asymptoter/geochallenge-backend/app/server/main .
COPY ./config ./config
CMD ["./main"]
EXPOSE 8080
