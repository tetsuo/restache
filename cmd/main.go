package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/onur1/stache/parser"
)

func main() {
	if err := parser.Parse(os.Stdin, func(tree *parser.TagNode) bool {
		b, err := json.Marshal(tree.Serialize())
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(b))
		return true
	}); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
