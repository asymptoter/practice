package main

import (
	"fmt"
	"log"
	"os/exec"
)

func main() {
	res, err := exec.Command("pwd").Output()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(res))
}
