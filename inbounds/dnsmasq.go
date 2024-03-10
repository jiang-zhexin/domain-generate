package inbounds

import (
	"bufio"
	"bytes"
	"errors"
	"strings"

	"Domain_Generate/data"
	"Domain_Generate/log"
)

func loadFromDnsmasq(inbound Inbound) (*data.DomainList, error) {
	if inbound.AddTags == nil {
		log.Warn("未指定 add tags")
		log.Info("已跳过 %s", inbound.Src)
		return nil, errors.New("not add tags")
	}
	input, err := getInput(inbound.Type, inbound.Src)
	if err != nil {
		return nil, err
	}
	return parseDnsmasq(input, inbound.AddTags)
}

func parseDnsmasq(input []byte, tags *data.Tags) (*data.DomainList, error) {
	list := data.NewDomainList()
	scanner := bufio.NewScanner(bytes.NewReader(input))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		line = removeServer(line)
		if len(line) == 0 {
			continue
		}
		list.Add(line, data.DomainFull, tags)
	}
	return list, nil
}

func removeServer(line string) string {
	s := strings.Split(line, "/")
	if len(s) != 3 {
		return ""
	}
	return s[1]
}
