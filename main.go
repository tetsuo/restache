package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/onur1/stache/parser"
	"github.com/onur1/stache/template"
)

func main() {
	if len(os.Args) > 1 {
		tpls, err := template.ParseGlob(os.Args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		data := make([]interface{}, len(tpls))
		for i, v := range tpls {
			data[i] = v.Serialize()
		}
		body, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(string(body))
	} else {
		if err := parser.Parse(os.Stdin, func(tree parser.Node) bool {
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
}
