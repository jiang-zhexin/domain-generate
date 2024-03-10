package inbounds

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"Domain_Generate/data"
	"Domain_Generate/log"
)

type Inbound struct {
	Format  string     `json:"format"`
	Type    string     `json:"type"`
	Src     string     `json:"src"`
	AddTags *data.Tags `json:"addtags,omitempty"`
}

const (
	FormatDomainList = "domain-list"
	FormatGeosite    = "geosite"
	FormatList       = "list"
	FormatDnsmasq    = "dnsmasq"
	FormatURL        = "url"
	FormatHosts      = "hosts"
)
const (
	TypeLocal  = "local"
	TypeRemote = "remote"
)

func LoadAll(inbounds []Inbound) *data.DomainList {
	dtl := data.NewDomainList()
	var wg sync.WaitGroup
	newdl := make(chan *data.DomainList)
	for _, inbound := range inbounds {
		wg.Add(1)
		go workers(inbound, newdl, &wg)
	}
	go func(dl chan<- *data.DomainList, wg *sync.WaitGroup) {
		wg.Wait()
		close(dl)
	}(newdl, &wg)

	for dl := range newdl {
		dtl.Union(dl)
	}
	return dtl
}

func workers(inbound Inbound, dl chan<- *data.DomainList, wg *sync.WaitGroup) {
	defer wg.Done()
	newdl, err := load(inbound)
	if err != nil {
		log.Warn("inbound: %s 加载失败", inbound.Src)
		return
	}
	dl <- newdl
	log.Info("inbound: %s 加载成功", inbound.Src)
}

func load(inbound Inbound) (*data.DomainList, error) {
	switch inbound.Format {
	case FormatDomainList:
		return loadFromDomainList(inbound)
	case FormatGeosite:
		return loadFromGeosite(inbound)
	case FormatList:
		return loadFromList(inbound)
	case FormatDnsmasq:
		return loadFromDnsmasq(inbound)
	case FormatURL:
		return loadFromURL(inbound)
	case FormatHosts:
		return loadFromHosts(inbound)
	}
	log.Warn("未知输入格式: %s", inbound.Src)
	return nil, errors.New("unkown inbound.type")
}

func getInput(inboundType, src string) ([]byte, error) {
	switch inboundType {
	case TypeLocal:
		return local(src)
	case TypeRemote:
		return remote(src)
	}
	log.Warn("未知输入类型: %s", inboundType)
	return nil, errors.New("unkown inbound.type! ")
}

func remote(src string) ([]byte, error) {
	URL, err := url.Parse(src)
	if err != nil {
		return nil, err
	}
	client := http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := client.Get(URL.String())
	if err != nil {
		log.Warn("无法访问：%s", src)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Warn("无法下载 %s", src)
		return nil, errors.New("request error! ")
	}
	data, err := io.ReadAll(resp.Body)
	return data, err
}

func local(src string) ([]byte, error) {
	file, err := os.Open(src)
	if err != nil {
		log.Warn("无法打开文件 %s", src)
		return nil, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	return data, err
}
