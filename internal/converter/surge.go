// internal/subscription/converter/surge.go
package converter

import (
	"fmt"
	"goconverter/internal/subscription/model"
	"strings"
)

type SurgeConverter struct {
	info BaseInfo
}

func NewSurgeConverter(info BaseInfo) *SurgeConverter {
	return &SurgeConverter{
		info: info,
	}
}

func (s *SurgeConverter) Convert(nodes []*model.Node) (string, error) {
	var builder strings.Builder

	// 写入基础配置
	builder.WriteString("[General]\n")
	builder.WriteString("loglevel = notify\n")
	builder.WriteString("bypass-system = true\n")
	builder.WriteString("skip-proxy = 127.0.0.1,192.168.0.0/16,10.0.0.0/8,172.16.0.0/12,100.64.0.0/10,localhost,*.local,e.crashlytics.com,captive.apple.com,::ffff:0:0:0:0/1,::ffff:128:0:0:0/1\n")
	builder.WriteString("dns-server = system,114.114.114.114,8.8.8.8\n")
	builder.WriteString("allow-wifi-access = false\n\n")

	// 写入代理配置
	builder.WriteString("[Proxy]\n")
	builder.WriteString("DIRECT = direct\n")

	for _, node := range nodes {
		proxy, err := s.ConvertNode(node)
		if err != nil {
			return "", err
		}
		if proxyStr, ok := proxy.(string); ok {
			builder.WriteString(proxyStr + "\n")
		}
	}

	// 写入代理组
	builder.WriteString("\n[Proxy Group]\n")
	builder.WriteString("🚀 节点选择 = select,DIRECT")
	for _, node := range nodes {
		builder.WriteString("," + node.Name)
	}
	builder.WriteString("\n\n")

	// 写入规则
	builder.WriteString("[Rule]\n")
	for _, rule := range s.getRules() {
		builder.WriteString(rule + "\n")
	}

	return builder.String(), nil
}

func (s *SurgeConverter) ConvertNode(node *model.Node) (interface{}, error) {
	switch node.Type {
	case model.TypeSS:
		return fmt.Sprintf("%s = ss, %s, %d, encrypt-method=%s, password=%s%s",
			node.Name, node.Server, node.Port, node.Cipher, node.Password,
			s.getSSPluginOpts(node)), nil

	case model.TypeTrojan:
		return fmt.Sprintf("%s = trojan, %s, %d, password=%s, sni=%s, skip-cert-verify=%v",
			node.Name, node.Server, node.Port, node.Password,
			defaultIfEmpty(node.SNI, node.Server), node.AllowInsecure), nil

	case model.TypeVmess:
		return fmt.Sprintf("%s = vmess, %s, %d, username=%s, ws=%v, tls=%v%s",
			node.Name, node.Server, node.Port, node.UUID,
			node.Network == "ws", node.TLS,
			s.getVmessOpts(node)), nil

	default:
		return "", fmt.Errorf("unsupported node type: %s", node.Type)
	}
}

func (s *SurgeConverter) getSSPluginOpts(node *model.Node) string {
	if node.Plugin == "" {
		return ""
	}
	opts := []string{fmt.Sprintf(", plugin=%s", node.Plugin)}
	for k, v := range node.PluginOpts {
		opts = append(opts, fmt.Sprintf("%s=%s", k, v))
	}
	return strings.Join(opts, ",")
}

func (s *SurgeConverter) getVmessOpts(node *model.Node) string {
	var opts []string
	if node.Network == "ws" {
		if node.WsPath != "" {
			opts = append(opts, fmt.Sprintf("ws-path=%s", node.WsPath))
		}
		if len(node.WsHeaders) > 0 {
			for k, v := range node.WsHeaders {
				opts = append(opts, fmt.Sprintf("ws-headers=%s:%s", k, v))
			}
		}
	}
	if node.TLS {
		opts = append(opts, fmt.Sprintf("sni=%s", defaultIfEmpty(node.SNI, node.Server)))
	}
	if len(opts) > 0 {
		return ", " + strings.Join(opts, ", ")
	}
	return ""
}

func (s *SurgeConverter) getRules() []string {
	return []string{
		"DOMAIN-SUFFIX,google.com,🚀 节点选择",
		"DOMAIN-KEYWORD,google,🚀 节点选择",
		"DOMAIN-SUFFIX,ad.com,REJECT",
		"DOMAIN-KEYWORD,facebook,🚀 节点选择",
		"DOMAIN-KEYWORD,youtube,🚀 节点选择",
		"FINAL,🚀 节点选择",
	}
}

// 工具函数
func defaultIfEmpty(str, def string) string {
	if strings.TrimSpace(str) == "" {
		return def
	}
	return str
}
