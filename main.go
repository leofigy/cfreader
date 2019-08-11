package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/leofigy/cfreader/reader"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage %s <cloudformation template>", os.Args[0])
		return
	}

	translator := &reader.SyntaxTransformer{}

	// reading input json file
	payload, err := ioutil.ReadFile(os.Args[1])

	if err != nil {
		log.Fatalf("%s", err)
		return
	}

	err = json.Unmarshal(payload, translator)

	if err != nil {
		log.Fatalf("%s", err)
	}

	fmt.Println(translator)
}
