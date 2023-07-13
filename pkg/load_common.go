package pkg

import (
	"encoding/json"

	"separa/common/utils"

	"github.com/chainreactors/parsers"
	"github.com/chainreactors/parsers/iutils"
)

var (
	NameMap = utils.NameMap
	PortMap = utils.PortMap
	TagMap  = utils.TagMap
	//WorkFlowMap    map[string][]*Workflow
	Extractor      []*parsers.Extractor
	Extractors     = make(parsers.Extractors)
	ExtractRegexps = map[string][]*parsers.Extractor{}
)

type PortFinger struct {
	Name  string   `json:"name"`
	Ports []string `json:"ports"`
	Tags  []string `json:"tags"`
}

func LoadPortConfig() {
	var portfingers []PortFinger
	err := json.Unmarshal(LoadConfig("port"), &portfingers)

	if err != nil {
		iutils.Fatal("port config load FAIL!, " + err.Error())
	}
	for _, v := range portfingers {
		v.Ports = utils.ParsePorts(v.Ports)
		utils.NameMap.Append(v.Name, v.Ports...)
		for _, t := range v.Tags {
			utils.TagMap.Append(t, v.Ports...)
		}
		for _, p := range v.Ports {
			utils.PortMap.Append(p, v.Name)
		}
	}
}

func LoadExtractor() {
	err := json.Unmarshal(LoadConfig("extract"), &Extractor)
	if err != nil {
		iutils.Fatal("extract config load FAIL!, " + err.Error())
	}

	for _, extract := range Extractor {
		extract.Compile()

		ExtractRegexps[extract.Name] = []*parsers.Extractor{extract}
		for _, tag := range extract.Tags {
			if _, ok := ExtractRegexps[tag]; !ok {
				ExtractRegexps[tag] = []*parsers.Extractor{extract}
			} else {
				ExtractRegexps[tag] = append(ExtractRegexps[tag], extract)
			}
		}
	}
}
