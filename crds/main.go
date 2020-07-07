package main

import (
	"io/ioutil"
	"log"
	"os"
)

func main() {
	path := os.Args[1]
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("yamlFile.Get err	#%v ", err)
	}
	r := NewResourceDefinition(yamlFile)
	r.GenerateNodeJS(os.Stdout)
}
