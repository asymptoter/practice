package script

import (
	"bufio"
	"log"
	"os"
	"strings"
)

func toUpper(s string) string {
	res := strings.ToUpper(s[0:1])
	flag := false
	for i := 1; i < len(s); i++ {
		if flag {
			res = res + string(s[i])
			flag = false
		} else {
			res = res + strings.ToLower(string(s[i]))
		}
		if s[i] == ' ' {
			flag = true
		}
	}
	return res
}

func main() {
	file, err := os.Open("iso3166-1")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	f, err := os.Create("mapping.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		country := toUpper(scanner.Text()[:48])
		country = strings.Replace(country, " And ", " and ", -1)
		country = strings.Replace(country, " Of", " of", -1)
		b := strings.ToLower(scanner.Text()[48:50])
		f.WriteString("[\"" + b + "\", \"" + country + "\"],\n")
		//f.WriteString("\"" + country + "\", \n")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
