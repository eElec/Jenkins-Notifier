package utils

import (
	"io/ioutil"
	"log"
)

var iconMap map[string][]byte

func GetIcon(path string) []byte {
	if iconMap == nil {
		iconMap = make(map[string][]byte)
	}

	data, ok := iconMap[path]
	if ok {
		return data
	}

	data = convertIcon(path)
	iconMap[path] = data
	return data
}

func convertIcon(s string) []byte {
	b, err := ioutil.ReadFile(s)
	if err != nil {
		log.Print(err)
	}
	return b
}
