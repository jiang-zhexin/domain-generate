package inbounds

import (
	"bufio"
	"bytes"
	"errors"
	"net"
	"net/url"
	"strings"

	"Domain_Generate/data"
	"Domain_Generate/log"
)

func loadFromURL(inbound Inbound) (*data.DomainList, error) {
	if inbound.AddTags == nil {
		log.Warn("未指定 add tags")
		log.Info("已跳过 %s", inbound.Src)
		return nil, errors.New("not add tags")
	}
	input, err := getInput(inbound.Type, inbound.Src)
	if err != nil {
		return nil, err
	}
	return parseURL(input, inbound.AddTags)
}

func parseURL(input []byte, tags *data.Tags) (*data.DomainList, error) {
	list := data.NewDomainList()
	scanner := bufio.NewScanner(bytes.NewReader(input))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		line = removeComment(line)
		if len(line) == 0 {
			continue
		}
		line, err := getHost(line)
		if err != nil {
			continue
		}
		list.Add(line, data.DomainFull, tags)
	}
	return list, nil
}

func getHost(urlStr string) (string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		log.Warn("输入值不是URL: %s", urlStr)
		return "", err
	}
	host := u.Hostname()

	if net.ParseIP(host) != nil {
		log.Warn("输入URL: %s host部分不是域名", urlStr)
		return "", errors.New("host is not domain")
	}
	return host, nil
}
