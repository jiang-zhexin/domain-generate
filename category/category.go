package category

import (
	"Domain_Generate/data"
	"Domain_Generate/log"
	"Domain_Generate/outbounds"
)

type Rule struct {
	Category   string           `json:"category"`
	Over       int64            `json:"over,omitempty"`
	DomainType DomainType       `json:"domain type,omitempty"`
	TagWeight  map[string]int64 `json:"tag weight"`
}

type DomainType []string

func (dt DomainType) include(s string) bool {
	if dt == nil {
		return true
	}
	for _, t := range dt {
		if t == s {
			return true
		}
	}
	return false
}

func ParseRules(dl *data.DomainList, rules []Rule) outbounds.DomainRes {
	dr := make(outbounds.DomainRes)
	for _, rule := range rules {
		res := data.NewDomainList()
		log.Info("已创建域名集 %s", rule.Category)
		dr[rule.Category] = res
	}
	for domain, tags := range dl.Full {
		for _, rule := range rules {
			if rule.DomainType.include(data.DomainFull) && countWeight(*tags, rule.TagWeight) > rule.Over {
				dr[rule.Category].Full[domain] = tags
			}
		}
	}
	for domain, tags := range dl.Suffix {
		for _, rule := range rules {
			if rule.DomainType.include(data.DomainSuffix) && countWeight(*tags, rule.TagWeight) > rule.Over {
				dr[rule.Category].Suffix[domain] = tags
			}
		}
	}
	for domain, tags := range dl.Keyword {
		for _, rule := range rules {
			if rule.DomainType.include(data.DomainKeyword) && countWeight(*tags, rule.TagWeight) > rule.Over {
				dr[rule.Category].Keyword[domain] = tags
			}
		}
	}
	for domain, tags := range dl.Regexp {
		for _, rule := range rules {
			if rule.DomainType.include(data.DomainRegexp) && countWeight(*tags, rule.TagWeight) > rule.Over {
				dr[rule.Category].Regexp[domain] = tags
			}
		}
	}
	return dr
}

func countWeight(tags data.Tags, tagWeight map[string]int64) int64 {
	var count int64
	for tag := range tags.Value.Iter() {
		count += tagWeight[tag.(string)]
	}
	return count
}
