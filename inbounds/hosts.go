package inbounds

import (
	"bufio"
	"bytes"
	"errors"
	"strings"

	"Domain_Generate/data"
	"Domain_Generate/log"
)

func loadFromHosts(inbound Inbound) (*data.DomainList, error) {
	if inbound.AddTags == nil {
		log.Warn("未指定 add tags")
		log.Info("已跳过 %s", inbound.Src)
		return nil, errors.New("not add tags")
	}
	input, err := getInput(inbound.Type, inbound.Src)
	if err != nil {
		return nil, err
	}
	return parseHosts(input, inbound.AddTags)
}

func parseHosts(input []byte, tags *data.Tags) (*data.DomainList, error) {
	list := data.NewDomainList()
	scanner := bufio.NewScanner(bytes.NewReader(input))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		line = removeComment(line)
		if len(line) == 0 {
			continue
		}
		ds := strings.Fields(line)
		if len(ds) < 2 {
			log.Warn("无法解析 hosts 行: %s", line)
			continue
		}
		for _, d := range ds[1:] {
			if d[0] == '*' {
				list.Add(d[2:], data.DomainSuffix, tags)
			} else {
				list.Add(d, data.DomainFull, tags)
			}
		}
	}
	return list, nil
}
