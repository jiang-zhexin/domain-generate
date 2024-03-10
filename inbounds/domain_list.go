package inbounds

import (
	"bytes"
	"encoding/json"

	"Domain_Generate/data"
	"Domain_Generate/log"
)

func loadFromDomainList(inbound Inbound) (*data.DomainList, error) {
	input, err := getInput(inbound.Type, inbound.Src)
	if err != nil {
		return nil, err
	}
	dtl, err := parseDomainList(input)
	if err != nil {
		return nil, err
	}
	if inbound.AddTags != nil {
		dtl.AddTag(inbound.AddTags)
	}
	return dtl, nil
}

func parseDomainList(input []byte) (*data.DomainList, error) {
	dtl := data.NewDomainList()
	decoder := json.NewDecoder(bytes.NewReader(input))
	err := decoder.Decode(dtl)
	if err != nil {
		log.Warn("解析 domain list 错误")
		return nil, err
	}
	return dtl, nil
}
