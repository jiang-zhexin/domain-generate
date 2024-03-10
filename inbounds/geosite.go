package inbounds

import (
	"strings"

	"Domain_Generate/data"
	"Domain_Generate/log"

	"github.com/v2fly/v2ray-core/v5/app/router/routercommon"
	"google.golang.org/protobuf/proto"
)

func loadFromGeosite(inbound Inbound) (*data.DomainList, error) {
	input, err := getInput(inbound.Type, inbound.Src)
	if err != nil {
		return nil, err
	}
	dtl, err := parseGeosite(input)
	if err != nil {
		log.Warn("解析 geosite 错误")
		return nil, err
	}
	if inbound.AddTags != nil {
		dtl.AddTag(inbound.AddTags)
	}
	return dtl, nil
}

func parseGeosite(vGeositeData []byte) (*data.DomainList, error) {
	vGeositeList := routercommon.GeoSiteList{}
	err := proto.Unmarshal(vGeositeData, &vGeositeList)
	if err != nil {
		return nil, err
	}
	dtl := data.NewDomainList()
	for _, vGeositeEntry := range vGeositeList.Entry {
		code := strings.ToLower(vGeositeEntry.CountryCode)
		for _, domain := range vGeositeEntry.Domain {
			tag := data.NewTags(code)
			if len(domain.Attribute) > 0 {
				for _, attribute := range domain.Attribute {
					tag.Add(attribute.Key)
				}
			}
			switch domain.Type {
			case routercommon.Domain_RootDomain:
				dtl.Add(domain.Value, data.DomainSuffix, tag)
			case routercommon.Domain_Full:
				dtl.Add(domain.Value, data.DomainFull, tag)
			case routercommon.Domain_Plain:
				dtl.Add(domain.Value, data.DomainKeyword, tag)
			case routercommon.Domain_Regex:
				dtl.Add(domain.Value, data.DomainRegexp, tag)
			}
		}

	}
	return dtl, nil
}
