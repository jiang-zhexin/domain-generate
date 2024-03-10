package main

import (
	"Domain_Generate/category"
	"Domain_Generate/cmd"
	"Domain_Generate/inbounds"
	"Domain_Generate/option"
)

func main() {
	cmd.GetCmd()
	config, err := option.LoadConfig(*cmd.Config)
	if err != nil {
		return
	}

	domainList := inbounds.LoadAll(config.Inbounds)
	if *cmd.PutAllDomains != "" {
		domainList.Save(*cmd.PutAllDomains)
	}

	result := category.ParseRules(domainList, config.Rules)
	result.Save(config.Outbounds)
}
