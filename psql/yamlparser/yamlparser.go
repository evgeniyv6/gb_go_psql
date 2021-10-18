package yamlparser

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var data = `
UVS_SCENARIO_ONE:
  DB:
    psql_url: http://foo.bar
  d: [3, 4]
`

func YamlParser() {
	var (
		domain string = "SIGMA"
		env_key string = "psi"

	)
	m := make(map[interface{}]interface{})
	ymlFile, err := ioutil.ReadFile("/opt/yevgen/githome/publicgit/gb/gb_go_psql/psql/yamlparser/f.yaml")
	if err != nil {
		log.Fatalf("Read file error %v", err)
	}
	err = yaml.Unmarshal(ymlFile, &m)
	if err != nil {
		fmt.Printf("yaml parser err{%s}\n", err)
	}
	fmt.Println(m)
	fmt.Println(m["UVS_SCENARIO_ONE"].(map[interface{}]interface{})[domain].(map[interface{}]interface{})["DB"].(map[interface{}]interface{})[env_key])
}



