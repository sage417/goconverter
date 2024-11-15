// internal/subscription/parser/trojan.go
package parser

import (
	"goconverter/internal/subscription/model"
	"net/url"
	"strings"
)

type TrojanParser struct{}

func NewTrojanParser() *TrojanParser {
	return &TrojanParser{}
}

func (p *TrojanParser) Match(link string) bool {
	return strings.HasPrefix(link, "trojan://")
}

func (p *TrojanParser) Parse(link string) (*model.Node, error) {
	u, err := url.Parse(link)
	if err != nil {
		return nil, err
	}

	password := u.User.Username()
	host := u.Host
	name, _ := url.QueryUnescape(strings.TrimPrefix(u.Fragment, "#"))

	// 解析查询参数
	query := u.Query()
	settings := make(map[string]string)
	for k, v := range query {
		if len(v) > 0 {
			settings[k] = v[0]
		}
	}

	// 分离端口
	hostParts := strings.Split(host, ":")
	server := hostParts[0]
	port := 443
	if len(hostParts) > 1 {
		port = parseInt(hostParts[1])
	}

	return &model.Node{
		Type:     "trojan",
		Name:     name,
		Server:   server,
		Port:     port,
		Protocol: "trojan",
		Password: password,
		Settings: settings,
	}, nil
}
