package cmd

import "flag"

var (
	Config        *string
	PutAllDomains *string
)

func init() {
	Config = flag.String("c", "config.json", "set config file path")
	PutAllDomains = flag.String("putall", "", "put all domains")
}

func GetCmd() {
	flag.Parse()
}
