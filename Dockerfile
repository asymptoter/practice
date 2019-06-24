FROM golang:latest
ENV GO111MODULE on
RUN go get github.com/pilu/fresh
COPY . $GOPATH/src/github.com/asymptoter/practice/app/server
WORKDIR $GOPATH/src/github.com/asymptoter/practice/app/server
CMD fresh -c fresh.conf main.go
EXPOSE 8080
