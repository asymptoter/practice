FROM golang:latest AS builder
ENV GO111MODULE on
WORKDIR /Users/ayuang/go/src/github.com/asymptoter/practice-backend
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
WORKDIR ./app/server
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main .

#FROM scratch
#COPY --from=builder /src/github.com/asymptoter/practice-backend/app/server/main .
#COPY ./config ./config
#COPY ./scripts/wait-for-it.sh ./scripts/wait-for-it.sh
CMD go run main.go
EXPOSE 8080
