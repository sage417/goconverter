// internal/config/config.go
package config

import (
	"fmt"
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

type RuleProvider struct {
	Name     string
	Type     string
	Url      string
	Path     string
	Interval int
	Proxy    string
	Behavior string
	Format   string
}

// ProxyGroup 表示一个代理组配置
type ProxyGroup struct {
	Name      string
	Type      string // select/url-test
	Proxies   []string
	URL       string // 用于 url-test
	Interval  int    // 用于 url-test
	Timeout   int    // 用于 url-test
	Tolerance int    // 用于 url-test
}

// ClashConfig 存储完整的配置
type ClashConfig struct {
	RuleSets        []ClashRule
	ProxyGroups     []ProxyGroup
	RuleProviders   []RuleProvider
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
		interval_idx := len(parts) - 2 - 1

		for idx, option := range parts[2:] {
			if option == "" {
				continue
			}

			if idx < interval_idx-1 {
				group.Proxies = append(group.Proxies, option)
				continue
			}

			if idx == interval_idx-1 {
				group.URL = option
				continue
			}

			if idx == interval_idx {
				testOptions := strings.Split(option, ",")
				if len(testOptions) > 0 {
					if num, err := strconv.Atoi(testOptions[0]); err == nil {
						group.Interval = num
					}
				}
				// if len(testOptions) > 1 {
				// 	if num, err := strconv.Atoi(testOptions[1]); err == nil {
				// 		group.Timeout = num
				// 	}
				// }
				if len(testOptions) > 2 {
					if num, err := strconv.Atoi(testOptions[2]); err == nil {
						group.Tolerance = num
					}
				}
				break
			}
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

	// 解析 ruleset
	rulesetKeys := section.Key("ruleset").ValueWithShadows()
	for _, ruleStr := range rulesetKeys {
		parts := strings.SplitN(ruleStr, ",", 4)
		if len(parts) < 2 {
			continue
		}
		// inline
		// ruleset=🎯 全球直连,[]GEOIP,CN
		// ruleset=🐟 漏网之鱼,[]FINAL

		//  - GEOIP,CN,🎯 全球直连
		//  - MATCH,🐟 漏网之鱼
		if ruleType, found := strings.CutPrefix(parts[1], "[]"); found {
			var ruleParm, noResolve string
			if len(parts) > 2 {
				ruleParm = parts[2]
			}

			if len(parts) > 3 {
				noResolve = parts[3]
			}
			rule := ClashRule{
				Strategy:  parts[0],
				Type:      ruleType,
				Pararm:    ruleParm,
				NoResolve: noResolve,
			}
			config.RuleSets = append(config.RuleSets, rule)
		} else {
			// remote or local file
			var contentUrl, behavior string
			// convert to online rule
			interval := 3600 * 24

			if after, found := strings.CutPrefix(parts[1], "clash-"); found {
				behavior, contentUrl, _ = strings.Cut(after, ":")
			} else {
				behavior = "classical"
				contentUrl = parts[1]
			}

			if behavior == "classic" {
				behavior = "classical"
			}

			if strings.HasPrefix(contentUrl, "rules/ACL4SSR/Clash/") {
				contentUrl = "https://raw.githubusercontent.com/ACL4SSR/ACL4SSR/refs/heads/master/Clash" +
					strings.SplitAfterN(parts[1], "/Clash", 2)[1]
			}
			// use cdn to fetch content
			contentUrl = strings.Replace(contentUrl, "https://raw.githubusercontent.com/ACL4SSR/ACL4SSR/refs/heads/master",
				"https://cdn.jsdelivr.net/gh/ACL4SSR/ACL4SSR@master", 1)

			fileName := path.Base(contentUrl)
			name, format, _ := strings.Cut(fileName, ".")

			ruleProvider := RuleProvider{
				Name:     name,
				Type:     "http",
				Url:      contentUrl,
				Interval: interval,
				Behavior: behavior,
				Format:   format,
			}

			config.RuleProviders = append(config.RuleProviders, ruleProvider)
			config.RuleSets = append(config.RuleSets, ClashRule{
				Strategy:  parts[0],
				Type:      "RULE-SET",
				Pararm:    name,
				NoResolve: "",
			})
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
