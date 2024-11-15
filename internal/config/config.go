// internal/config/config.go
package config

import (
	"fmt"
	"goconverter/internal/fetcher"
	"path"
	"strconv"
	"strings"

	"gopkg.in/ini.v1"
)

// ClashRule 表示一个规则配置
type ClashRule struct {
	Type      string
	Pararm    string
	Strategy  string
	NoResolve string
}

// ProxyGroup 表示一个代理组配置
type ProxyGroup struct {
	Name      string
	Type      string // select/url-test
	Proxies   []string
	URL       string // 用于 url-test
	Interval  int    // 用于 url-test
	Tolerance int    // 用于 url-test
}

// ClashConfig 存储完整的配置
type ClashConfig struct {
	RuleSets        []ClashRule
	ProxyGroups     []ProxyGroup
	EnableGenerator bool
	OverwriteRules  bool
}

func parseProxyGroup(value string) ProxyGroup {
	parts := strings.Split(value, "`")
	group := ProxyGroup{
		Name: parts[0],
	}

	if len(parts) > 1 {
		group.Type = parts[1]
	}

	// custom_proxy_group=Group_Name`url-test|fallback|load-balance`Rule_1`Rule_2`...`test_url`interval[,timeout][,tolerance]
	// custom_proxy_group=Group_Name`select`Rule_1`Rule_2`...
	if group.Type == "select" {
		for _, option := range parts[2:] {
			if option == "" {
				continue
			}
			group.Proxies = append(group.Proxies, option)
		}
	} else if group.Type == "url-test" || group.Type == "fallback" || group.Type == "load-balance" {
		parseTestOption := false
		for _, option := range parts[2:] {
			if option == "" {
				continue
			}
			if strings.HasPrefix(option, "http://") || strings.HasPrefix(option, "https://") {
				group.URL = option
				parseTestOption = true
				continue
			}
			if parseTestOption {
				testOptions := strings.Split(option, ",")
				if len(testOptions) > 0 {
					if num, err := strconv.Atoi(testOptions[0]); err == nil {
						group.Interval = num
					}
				}
				if len(testOptions) > 2 {
					if num, err := strconv.Atoi(testOptions[2]); err == nil {
						group.Tolerance = num
					}
				}
			}

			group.Proxies = append(group.Proxies, option)
		}
	}

	return group
}

func ParseConfig(content []byte) (*ClashConfig, error) {
	cfg, err := ini.LoadSources(ini.LoadOptions{
		AllowShadows:             true,
		Insensitive:              true,
		SpaceBeforeInlineComment: true,
	}, content)

	if err != nil {
		return nil, fmt.Errorf("failed to load config: %v", err)
	}

	section := cfg.Section("custom")
	config := &ClashConfig{}

	contentFetcher := fetcher.NewFetcher()

	// 解析 ruleset
	rulesetKeys := section.Key("ruleset").ValueWithShadows()
	for _, ruleStr := range rulesetKeys {
		parts := strings.SplitN(ruleStr, ",", 2)
		if len(parts) >= 2 {
			// ruleset=🎯 全球直连,[]GEOIP,CN
			// ruleset=🐟 漏网之鱼,[]FINAL

			//  - GEOIP,CN,🎯 全球直连
			//  - MATCH,🐟 漏网之鱼
			if after, found := strings.CutPrefix(parts[1], "[]"); found {
				typeAndParam := strings.SplitN(after, ",", 2)
				ruleType := typeAndParam[0]
				ruleParm := ""
				if len(typeAndParam) > 1 {
					ruleParm = typeAndParam[1]
				}
				rule := ClashRule{
					Strategy:  parts[0],
					Type:      ruleType,
					Pararm:    ruleParm,
					NoResolve: "",
				}
				config.RuleSets = append(config.RuleSets, rule)
			} else {
				contentUrl := "https://raw.githubusercontent.com/ACL4SSR/ACL4SSR/refs/heads/master/Clash" + strings.SplitAfterN(parts[1], "/Clash", 2)[1]
				listContent, err := contentFetcher.Fetch(contentUrl)
				if err != nil {
					continue
				}
				for _, line := range strings.Split(string(listContent), "\n") {
					if strings.HasPrefix(line, "#") || line == "" {
						continue
					}
					lineParam := strings.SplitN(line, ",", 3)
					param := ""
					if len(lineParam) > 1 {
						param = lineParam[1]
					}
					noReslove := ""
					if len(lineParam) > 2 {
						noReslove = lineParam[2]
					}
					rule := ClashRule{
						Strategy:  parts[0],
						Type:      lineParam[0],
						Pararm:    param,
						NoResolve: noReslove,
					}
					config.RuleSets = append(config.RuleSets, rule)
				}
			}
		}
	}

	// 解析 custom_proxy_group
	groupKeys := section.Key("custom_proxy_group").ValueWithShadows()
	for _, groupStr := range groupKeys {
		group := parseProxyGroup(groupStr)
		config.ProxyGroups = append(config.ProxyGroups, group)
	}

	// 解析其他设置
	config.EnableGenerator = section.Key("enable_rule_generator").MustBool(false)
	config.OverwriteRules = section.Key("overwrite_original_rules").MustBool(false)

	return config, nil
}

func getLastTwoPaths(urlPath string) string {
	// 使用 path.Clean 清理路径
	cleanPath := path.Clean(urlPath)
	parts := strings.Split(cleanPath, "/")

	if len(parts) <= 2 {
		return cleanPath
	}

	return "/" + strings.Join(parts[len(parts)-2:], "/")
}
