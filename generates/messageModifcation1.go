package generates

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

func MessageModificationV1() {
	files, err := ioutil.ReadDir("./messages/")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fileContent, err := os.ReadFile("./messages/" + f.Name())
		if err != nil {
			return
		}
		var rgx = regexp.MustCompile(`^([\s\S]*func)(.+)(NOA[\s\S]*func \(([^)]+)[\s\S]*)`)
		rs := rgx.FindStringSubmatch(string(fileContent))
		if len(rs) == 5 {
			newContent := fmt.Sprintf("%s (%s) Get%s", rs[1], rs[4], rs[3])
			_ = os.WriteFile("./messages/"+f.Name(), []byte(newContent), 0644)
		}
	}
}

func MessageModificationV2() {
	files, err := ioutil.ReadDir("./messages/")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fileContent, err := os.ReadFile("./messages/" + f.Name())
		if err != nil {
			return
		}
		var rgx = regexp.MustCompile(`/^[\s\S]*NOA\(instance uint\)([^{]+)([\s\S]*)/gm`)
		rs := rgx.FindStringSubmatch(string(fileContent))
		if len(rs) == 3 {
			newContent := fmt.Sprintf("%s Message %s", rs[1], rs[2])
			_ = os.WriteFile("./messages/"+f.Name(), []byte(newContent), 0644)
		}
	}
}
