// internal/subscription/converter/clash.go
package converter

import (
	"fmt"
	"goconverter/internal/config"
	"goconverter/internal/subscription/model"
	"slices"
	"strings"

	"github.com/goccy/go-yaml"
)

type ClashConverter struct {
	info *BaseInfo
}

func NewClashConverter(info *BaseInfo) *ClashConverter {
	return &ClashConverter{
		info: info,
	}
}

type ClashConfig struct {
	Port               int                      `yaml:"port"`
	SocksPort          int                      `yaml:"socks-port"`
	AllowLan           bool                     `yaml:"allow-lan"`
	Mode               string                   `yaml:"mode"`
	LogLevel           string                   `yaml:"log-level"`
	ExternalController string                   `yaml:"external-controller"`
	Secret             string                   `yaml:"secret"`
	DNS                map[string]interface{}   `yaml:"dns"`
	Proxies            []map[string]interface{} `yaml:"proxies"`
	ProxyGroups        []*ProxyGroup            `yaml:"proxy-groups"`
	Rules              []string                 `yaml:"rules"`
}

// ClashProxy 定义单个代理服务器配置
type ClashProxy struct {
	Name     string `yaml:"name"`             // 代理名称
	Type     string `yaml:"type"`             // 代理类型：ss/ssr/vmess/trojan/http
	Server   string `yaml:"server"`           // 服务器地址
	Port     int    `yaml:"port"`             // 端口号
	Password string `yaml:"password"`         // 密码
	UUID     string `yaml:"uuid,omitempty"`   // UUID(VMess)
	Cipher   string `yaml:"cipher,omitempty"` // 加密方式
	UDP      bool   `yaml:"udp,omitempty"`    // 是否启用UDP

	// TLS相关配置
	TLS            bool     `yaml:"tls,omitempty"`    // 是否启用TLS
	SkipCertVerify bool     `yaml:"skip-cert-verify"` // 是否跳过证书验证
	Alpn           []string `yaml:"alpn,omitempty"`   // ALPN配置
	SNI            string   `yaml:"sni,omitempty"`    // SNI配置

	// 传输层配置
	Network   string            `yaml:"network,omitempty"`    // 传输协议：ws/h2/grpc
	WsPath    string            `yaml:"ws-path,omitempty"`    // WebSocket路径
	WsHeaders map[string]string `yaml:"ws-headers,omitempty"` // WebSocket请求头

	// 插件配置
	Plugin     string                 `yaml:"plugin,omitempty"`      // 插件名称
	PluginOpts map[string]interface{} `yaml:"plugin-opts,omitempty"` // 插件配置
}

type ProxyGroup struct {
	Name      string   `yaml:"name"`
	Type      string   `yaml:"type"`                // select/url-test
	URL       string   `yaml:"url,omitempty"`       // 用于 url-test
	Interval  int      `yaml:"interval,omitempty"`  // 用于 url-test
	Tolerance int      `yaml:"tolerance,omitempty"` // 用于 url-test
	Proxies   []string `yaml:"proxies"`
}

func (c *ClashConverter) Convert(nodes []*model.Node, clashConfig *config.ClashConfig) (string, error) {
	config := &ClashConfig{
		Port:               7890,
		SocksPort:          7891,
		AllowLan:           false,
		Mode:               "rule",
		LogLevel:           "info",
		ExternalController: "127.0.0.1:9090",
		DNS: map[string]interface{}{
			"enable":     true,
			"ipv6":       false,
			"nameserver": []string{"114.114.114.114", "8.8.8.8"},
		},
		Proxies:     make([]map[string]interface{}, 0, len(nodes)),
		ProxyGroups: make([]*ProxyGroup, 0),
		Rules:       make([]string, 0),
	}

	// 转换所有节点
	nodeNames := make([]string, 0, len(nodes))
	for _, node := range nodes {
		proxy, err := c.ConvertNode(node)
		if err != nil {
			return "", err
		}
		if proxyMap, ok := proxy.(map[string]interface{}); ok {
			config.Proxies = append(config.Proxies, proxyMap)
			nodeNames = append(nodeNames, node.Name)
		}
	}

	// 添加代理组
	for _, configProxyGroup := range clashConfig.ProxyGroups {
		proxyGroup := &ProxyGroup{
			Name:      configProxyGroup.Name,
			Type:      configProxyGroup.Type,
			URL:       configProxyGroup.URL,
			Interval:  configProxyGroup.Interval,
			Tolerance: configProxyGroup.Tolerance,
			Proxies:   make([]string, 0),
		}

		for _, name := range configProxyGroup.Proxies {
			if after, found := strings.CutPrefix(name, "[]"); found {
				proxyGroup.Proxies = append(proxyGroup.Proxies, after)
			} else {
				for _, nodeName := range nodeNames {
					proxyGroup.Proxies = append(proxyGroup.Proxies, nodeName)
				}
			}
		}

		config.ProxyGroups = append(config.ProxyGroups, proxyGroup)
	}

	// 添加规则
	for _, rule := range c.getRules(clashConfig) {
		config.Rules = append(config.Rules, rule)
	}

	// 转换为YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return "", fmt.Errorf("failed to marshal clash config: %v", err)
	}

	return string(data), nil
}

func (c *ClashConverter) ConvertNode(node *model.Node) (interface{}, error) {
	return node.ToClash(), nil
}

func (c *ClashConverter) getRules(clashConfig *config.ClashConfig) []string {
	rules := make([]string, 0)
	for _, ruleset := range clashConfig.RuleSets {
		if ruleset.Type == "FINAL" {
			ruleset.Type = "MATCH"
		}
		// if ruleset.Type == "USER-AGENT" {
		// 	continue
		// }
		if !slices.Contains([]string{"DOMAIN", "DOMAIN-SUFFIX", "DOMAIN-KEYWORD",
			"GEOIP", "IP-CIDR", "IP-CIDR6", "SRC-IP-CIDR", "SRC-PORT", "DST-PORT",
			"PROCESS-NAME", "PROCESS-PATH", "IPSET", "RULE-SET", "SCRIPT", "MATCH"}, ruleset.Type) {
			continue
		}
		rulesetslice := make([]string, 0)
		rulesetslice = append(rulesetslice, ruleset.Type)
		if ruleset.Pararm != "" {
			rulesetslice = append(rulesetslice, ruleset.Pararm)
		}
		rulesetslice = append(rulesetslice, ruleset.Strategy)
		if ruleset.NoResolve != "" {
			rulesetslice = append(rulesetslice, ruleset.NoResolve)
		}
		rules = append(rules, strings.Join(rulesetslice, ","))
	}

	return rules
}
