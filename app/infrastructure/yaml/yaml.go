package yaml

import (
	"io/ioutil"
	"log"

	yml "gopkg.in/yaml.v2"
)

// ParseFromFile parse yaml from file path
func ParseFromFile(path string) (map[interface{}]interface{}, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return ParseFromString(string(data))
}

// ParseFromString parse yaml from string
func ParseFromString(data string) (map[interface{}]interface{}, error) {
	m := make(map[interface{}]interface{})

	err := yml.Unmarshal([]byte(data), &m)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return m, err
}
