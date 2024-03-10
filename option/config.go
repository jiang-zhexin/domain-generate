package option

import (
	"encoding/json"
	"os"

	"Domain_Generate/category"
	"Domain_Generate/inbounds"
	"Domain_Generate/log"
	"Domain_Generate/outbounds"
)

type Option struct {
	Inbounds  []inbounds.Inbound   `json:"inbounds"`
	Rules     []category.Rule      `jaon:"rules"`
	Outbounds []outbounds.Outbound `json:"outbounds"`
}

func LoadConfig(path string) (*Option, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Error("无法打开文件 %s", path)
		return nil, err
	}
	defer file.Close()
	var o Option

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&o)
	if err != nil {
		log.Error("无法解析配置 %s", path)
		return nil, err
	}
	return &o, nil
}
