// internal/subscription/parser/shadowsocks.go
package parser

import (
	"encoding/base64"
	"errors"
	"fmt"
	"goconverter/internal/subscription/model"
	"net/url"
	"strings"
)

type ShadowsocksParser struct{}

func NewShadowsocksParser() *ShadowsocksParser {
	return &ShadowsocksParser{}
}

func (p *ShadowsocksParser) Match(link string) bool {
	return strings.HasPrefix(link, "ss://")
}

func (p *ShadowsocksParser) Parse(link string) (*model.Node, error) {
	// 移除 "ss://" 前缀
	link = strings.TrimPrefix(link, "ss://")

	// 分离名称部分（如果存在）
	var name string
	if idx := strings.Index(link, "#"); idx != -1 {
		name = link[idx+1:]
		link = link[:idx]
		name, _ = url.QueryUnescape(name)
	}

	// 解析主体部分
	decodedBytes, err := base64.RawURLEncoding.DecodeString(link)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ss link: %v", err)
	}
	decoded := string(decodedBytes)

	// 分离方法和密码
	parts := strings.SplitN(decoded, "@", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid ss link format")
	}

	methodAndPass := strings.SplitN(parts[0], ":", 2)
	if len(methodAndPass) != 2 {
		return nil, errors.New("invalid method and password format")
	}

	serverAndPort := strings.SplitN(parts[1], ":", 2)
	if len(serverAndPort) != 2 {
		return nil, errors.New("invalid server and port format")
	}

	return &model.Node{
		Type:     "ss",
		Name:     name,
		Server:   serverAndPort[0],
		Port:     parseInt(serverAndPort[1]),
		Protocol: methodAndPass[0],
		Password: methodAndPass[1],
		Settings: make(map[string]string),
	}, nil
}
