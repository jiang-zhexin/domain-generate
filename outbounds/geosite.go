package outbounds

import (
	"os"
	"sort"

	"Domain_Generate/data"
	"Domain_Generate/log"

	router "github.com/v2fly/v2ray-core/v5/app/router/routercommon"
	"google.golang.org/protobuf/proto"
)

func (res DomainRes) geositeSave(outbound Outbound) error {
	protoList := new(router.GeoSiteList)
	if outbound.Include == nil || len(outbound.Include) == 0 {
		for category, dtl := range res {
			site := toGeosite(dtl, category)
			protoList.Entry = append(protoList.Entry, site)
		}
	} else {
		for _, category := range outbound.Include {
			site := toGeosite(res[category], category)
			protoList.Entry = append(protoList.Entry, site)
		}
	}
	sort.SliceStable(protoList.Entry, func(i, j int) bool {
		return protoList.Entry[i].CountryCode < protoList.Entry[j].CountryCode
	})

	protoBytes, err := proto.Marshal(protoList)
	if err != nil {
		log.Warn("geosite 保存错误...跳过保存：%s", outbound.Path)
		return err
	}
	return os.WriteFile(outbound.Path, protoBytes, 0644)
}

func toGeosite(dtl *data.DomainList, cc string) *router.GeoSite {
	site := &router.GeoSite{
		CountryCode: cc,
	}
	if l := len(dtl.Full) + len(dtl.Suffix) + len(dtl.Keyword) + len(dtl.Regexp); l > 0 {
		site.Domain = make([]*router.Domain, 0, l)
		for domain := range dtl.Full {
			site.Domain = append(site.Domain, &router.Domain{
				Type:  router.Domain_Full,
				Value: domain,
			})
		}
		for domain := range dtl.Suffix {
			site.Domain = append(site.Domain, &router.Domain{
				Type:  router.Domain_RootDomain,
				Value: domain,
			})
		}
		for domain := range dtl.Keyword {
			site.Domain = append(site.Domain, &router.Domain{
				Type:  router.Domain_Plain,
				Value: domain,
			})
		}
		for domain := range dtl.Regexp {
			site.Domain = append(site.Domain, &router.Domain{
				Type:  router.Domain_Regex,
				Value: domain,
			})
		}
	}
	return site
}
