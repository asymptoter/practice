package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("vim-go")
	http.Client{
		Transport: http.Transport{},
	}
}
