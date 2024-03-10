package data

import (
	"encoding/json"
	"os"
	"path/filepath"

	"Domain_Generate/log"
)

type DomainList struct {
	Full    DomainTags `json:"full,omitempty"`
	Suffix  DomainTags `json:"suffix,omitempty"`
	Keyword DomainTags `json:"keyword,omitempty"`
	Regexp  DomainTags `json:"regexp,omitempty"`
}
type DomainTags map[string]*Tags

const (
	DomainFull    = "full"
	DomainSuffix  = "suffix"
	DomainKeyword = "keyword"
	DomainRegexp  = "regexp"
)

func NewDomainList() *DomainList {
	return &DomainList{
		Full:    make(DomainTags),
		Suffix:  make(DomainTags),
		Keyword: make(DomainTags),
		Regexp:  make(DomainTags),
	}
}

func Load(path string) (*DomainList, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var dl DomainList

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&dl)
	if err != nil {
		return nil, err
	}
	return &dl, nil
}

func (dl *DomainList) Save(path string) error {
	jsonData, err := json.Marshal(dl)
	if err != nil {
		log.Warn("序列化错误：%s", path)
		return err
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Warn("无法创建目录：%s", dir)
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		log.Warn("无法创建文件：%s", path)
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}
	return nil
}

func (dl *DomainList) Union(otherDomainList *DomainList) {
	if otherDomainList == nil {
		return
	}
	for domain, tags := range otherDomainList.Full {
		dl.Full[domain] = dl.Full[domain].Union(tags)
	}
	for domain, tags := range otherDomainList.Suffix {
		dl.Suffix[domain] = dl.Suffix[domain].Union(tags)
	}
	for domain, tags := range otherDomainList.Keyword {
		dl.Keyword[domain] = dl.Keyword[domain].Union(tags)
	}
	for domain, tags := range otherDomainList.Regexp {
		dl.Regexp[domain] = dl.Regexp[domain].Union(tags)
	}
}

func (dl *DomainList) Add(domain string, domainType string, tags *Tags) {
	if tags == nil {
		return
	}
	switch domainType {
	case DomainFull:
		dl.Full[domain] = dl.Full[domain].Union(tags)
	case DomainSuffix:
		dl.Suffix[domain] = dl.Suffix[domain].Union(tags)
	case DomainKeyword:
		dl.Keyword[domain] = dl.Keyword[domain].Union(tags)
	case DomainRegexp:
		dl.Regexp[domain] = dl.Regexp[domain].Union(tags)
	}
}

func (dl *DomainList) AddTag(otherTags *Tags) {
	if otherTags == nil {
		return
	}
	for domain := range dl.Full {
		dl.Full[domain].Union(otherTags)
	}
	for domain := range dl.Suffix {
		dl.Suffix[domain].Union(otherTags)
	}
	for domain := range dl.Keyword {
		dl.Keyword[domain].Union(otherTags)
	}
	for domain := range dl.Regexp {
		dl.Regexp[domain].Union(otherTags)
	}
}
