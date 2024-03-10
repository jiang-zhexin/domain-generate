package inbounds

import (
	"bufio"
	"bytes"
	"errors"
	"strings"

	"Domain_Generate/data"
	"Domain_Generate/log"
)

func loadFromList(inbound Inbound) (*data.DomainList, error) {
	if inbound.AddTags == nil {
		log.Warn("未指定 add tags")
		log.Info("已跳过 %s", inbound.Src)
		return nil, errors.New("not add tags")
	}
	input, err := getInput(inbound.Type, inbound.Src)
	if err != nil {
		return nil, err
	}
	return parseList(input, inbound.AddTags)
}

func parseList(input []byte, tags *data.Tags) (*data.DomainList, error) {
	list := data.NewDomainList()
	scanner := bufio.NewScanner(bytes.NewReader(input))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		line = removeComment(line)
		if len(line) == 0 {
			continue
		}
		list.Add(line, data.DomainSuffix, tags)
	}
	return list, nil
}

func removeComment(line string) string {
	idx := strings.Index(line, "#")
	if idx == -1 {
		return line
	}
	return strings.TrimSpace(line[:idx])
}
