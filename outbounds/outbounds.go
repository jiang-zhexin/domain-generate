package outbounds

import (
	"path/filepath"

	"Domain_Generate/data"
	"Domain_Generate/log"
)

type Outbound struct {
	Format  string   `json:"format"`
	Path    string   `json:"path"`
	Include []string `json:"include,omitempty"`
}

type DomainRes map[string]*data.DomainList

const (
	FormatDomainList = "domain-list"
	FormatGeosite    = "geosite"
	FormatClash      = "clash"
	FormatRuleSet    = "rule-set"
)

func (res DomainRes) Save(outbounds []Outbound) {
	for _, outbound := range outbounds {
		switch outbound.Format {
		case FormatDomainList:
			res.domainListSave(outbound)
		case FormatGeosite:
			res.geositeSave(outbound)
		case FormatClash:
			res.clashSave(outbound)
		case FormatRuleSet:
			res.rulesetSave(outbound)
		default:
			log.Warn("未知输出格式: %s", outbound.Format)
			continue
		}
		log.Info("保存 outbound.%s 到 %s 成功", outbound.Format, outbound.Path)
	}
}

func (res DomainRes) domainListSave(outbound Outbound) error {
	if outbound.Include == nil || len(outbound.Include) == 0 {
		for category, dtl := range res {
			err := dtl.Save(filepath.Join(outbound.Path, category+".json"))
			if err != nil {
				log.Warn("跳过保存：%s", outbound.Path)
				return err
			}
		}
	} else {
		for _, category := range outbound.Include {
			err := res[category].Save(filepath.Join(outbound.Path, category+".json"))
			if err != nil {
				log.Warn("跳过保存：%s", outbound.Path)
				return err
			}
		}
	}
	return nil
}

func (res DomainRes) clashSave(outbound Outbound) error {
	if outbound.Include == nil || len(outbound.Include) == 0 {
		for category, dtl := range res {
			err := save2clash(dtl, filepath.Join(outbound.Path, category+".txt"))
			if err != nil {
				log.Warn("跳过保存：%s", outbound.Path)
				return err
			}
		}
	} else {
		for _, category := range outbound.Include {
			err := save2clash(res[category], filepath.Join(outbound.Path, category+".txt"))
			if err != nil {
				log.Warn("跳过保存：%s", outbound.Path)
				return err
			}
		}
	}
	return nil
}

func (res DomainRes) rulesetSave(outbound Outbound) error {
	if outbound.Include == nil || len(outbound.Include) == 0 {
		for category, dtl := range res {
			err := save2ruleset(dtl, filepath.Join(outbound.Path, category+".srs"))
			if err != nil {
				log.Warn("跳过保存：%s", outbound.Path)
				return err
			}
		}
	} else {
		for _, category := range outbound.Include {
			err := save2ruleset(res[category], filepath.Join(outbound.Path, category+".srs"))
			if err != nil {
				log.Warn("跳过保存：%s", outbound.Path)
				return err
			}
		}
	}
	return nil
}
