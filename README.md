# Domain Generate
本项目供需要自定义域名集的用户用于快速、方便地生成域名集

本项目不对域名来源负责

## How to use
运行 `domain_generate -c config.json`

# config 字段说明
> 在此处查看[示例配置](/blob/main/config.json)

## inbounds
`inbounds` 字段定义了域名来源

### format
`format` 字段定义了域名来源的格式，目前支持以下格式：`domain-list` `geosite` `list` `dnsmasq`
> 1. domain-list 是本项目的源格式，在[这里](#domain-list)有详细说明
> 2. geosite 是 domain-list-community 项目的格式，请前往[项目仓库](https://github.com/v2fly/domain-list-community)查看格式说明
> 3. list 是一个简单的域名列表，每行一个域名，Domain Generate 会忽略空行和 `#` 开头的行，并将所有域名解释为 `suffix` 类型
> 4. dnsmasq 是满足 dnsmasq 转发规则的配置文件，提取域名部分得到，例：`server=/0-6.com/114.114.114.114`
> 5. url 是一个 url 列表提取 host 部分得到，例：`udp://tracker.opentrackr.org:1337/announce`
> 6. hosts 是一个 hosts 文件，提取域名部分得到，例：`0.0.0.0 1oavsblobprodcus350.blob.core.windows.net`

###### geosite 输入时的行为
1. 保留 `full` `suffix` `keyword` `regexp`
2. 将 `domain` 所在文件及 `include` 该域名的所有文件名，以及其拥有的 `attr` 解释为 `tag`
> 注意 cn 分类和 @cn 会并成同一个 tag

### type
`type` 字段定义了域名来源，为：`local` `remote`

### addtags
`addtags` 字段会为该域名来源的所有域名添加 tags
> 1. 在 `domain-list` `geosite` 中可省略
> 2. 在 `list` `dnsmasq`中为**必填**项

### src
`src` 字段定义了域名来源的位置，值为 `path` 或者 `URL`

## rules
`rules` 字段定义了域名的分类

### category
`category` 字段定义了域名分类的名字

### over
`over` 字段定义了大于该权重的域名会被匹配

### domain type
`domain type` 可选填 `full` `suffix` `keyword` `regexp`，留空表示匹配所有

### tag weight
`tag weight` 字段定义了该分类下，tag 对应的权重

## outbounds
`outbounds` 定义了域名的输出格式

### format
`format` 字段定义了域名输出的格式，目前支持以下格式：`domain-list` `geosite` `clash` `rule-set`
> 1. domain-list 是本项目的源格式，在[这里](#domain-list)有详细说明
> 2. geosite 是 domain-list-community 项目的格式，请前往[项目仓库](https://github.com/v2fly/domain-list-community)查看格式说明
> 3. clash 是 clash 域名匹配格式，请前往[此处](https://wiki.metacubex.one/config/rule-providers/content/#__tabbed_1_2)查看格式说明
> 4. rule-set 是 sing-box 的路由规则格式，请前往[此处](https://sing-box.sagernet.org/zh/configuration/rule-set/source-format/)查看格式说明

##### geosite 输出时的行为
每个 `include` 为一个类别，转化为 geosite 格式

> v2fly 中使用方法：geosite:$include-name

##### clash 输出时的行为
> 仅实现 `classical.text` 格式

每个 `include` 的分类为一个文件，以 `$category.txt` 为文件名，转化为 clash 格式，`regexp` 自动忽略

##### rule-set 输出时的行为
每个 `include` 的分类为一个文件，以 `$category.srs` 为文件名，转化为 rule-set 格式

### path
`path` 字段定义了输出文件夹或输出文件(geosite时)

### include
`include` 字段定义了输出时包含的分类

> 分类名称应存在于 rules.category
> 留空表示包含所有分类


# domain list
`domain list` 是本项目的源格式，基本格式如下

```json
{
    "full": {
        "domain1": [
            "tag1",
            "tag2"
        ],
        "domain2": "tag3" //只有一个 tag 时，可省略 []
    },
    "suffix": {},
    "keyword": {},
    "regexp": {}
}
```