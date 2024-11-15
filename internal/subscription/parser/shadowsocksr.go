// internal/subscription/parser/shadowsocksr.go
package parser

import (
	"encoding/base64"
	"fmt"
	"goconverter/internal/subscription/model"
	"strings"
)

type ShadowsocksRParser struct{}

func NewShadowsocksRParser() *ShadowsocksRParser {
	return &ShadowsocksRParser{}
}

func (p *ShadowsocksRParser) Match(link string) bool {
	return strings.HasPrefix(link, "ssr://")
}

func (p *ShadowsocksRParser) Parse(link string) (*model.Node, error) {
	// 移除 "ssr://" 前缀
	link = strings.TrimPrefix(link, "ssr://")

	// Base64 解码
	decodedBytes, err := base64.RawURLEncoding.DecodeString(link)
	if err != nil {
		// 尝试标准Base64解码
		decodedBytes, err = base64.StdEncoding.DecodeString(link)
		if err != nil {
			return nil, fmt.Errorf("failed to decode ssr link: %v", err)
		}
	}

	decoded := string(decodedBytes)

	// 分离主要部分和参数部分
	parts := strings.SplitN(decoded, "/?", 2)
	if len(parts) < 1 {
		return nil, fmt.Errorf("invalid ssr link format")
	}

	// 解析主要部分
	mainParts := strings.Split(parts[0], ":")
	if len(mainParts) < 6 {
		return nil, fmt.Errorf("invalid ssr link main parts")
	}

	node := &model.Node{
		Type:     model.TypeSSR,
		Server:   mainParts[0],
		Port:     parseInt(mainParts[1]),
		Protocol: mainParts[2],
		Cipher:   mainParts[3],
		Obfs:     mainParts[4],
	}

	// 解码密码
	password, err := base64URLDecode(mainParts[5])
	if err != nil {
		return nil, fmt.Errorf("failed to decode password: %v", err)
	}
	node.Password = password

	// 解析参数部分
	if len(parts) > 1 {
		params := parts[1]
		values, err := parseSSRParams(params)
		if err != nil {
			return nil, fmt.Errorf("failed to parse ssr params: %v", err)
		}

		if obfsParam, ok := values["obfsparam"]; ok {
			node.ObfsParam = obfsParam
		}
		if protocolParam, ok := values["protoparam"]; ok {
			node.ProtocolParam = protocolParam
		}
		if remarks, ok := values["remarks"]; ok {
			node.Name = remarks
		}
		if group, ok := values["group"]; ok {
			node.Group = group
		}
	}

	if node.Name == "" {
		node.Name = fmt.Sprintf("SSR_%s_%d", node.Server, node.Port)
	}

	return node, nil
}

// 解析SSR参数
func parseSSRParams(params string) (map[string]string, error) {
	result := make(map[string]string)
	pairs := strings.Split(params, "&")

	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value, err := base64URLDecode(parts[1])
		if err != nil {
			continue
		}

		result[key] = value
	}

	return result, nil
}

// Base64 URL 解码
func base64URLDecode(s string) (string, error) {
	// 添加补齐的等号
	padding := 4 - len(s)%4
	if padding != 4 {
		s += strings.Repeat("=", padding)
	}

	// 尝试Raw解码
	bytes, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		// 尝试标准解码
		bytes, err = base64.URLEncoding.DecodeString(s)
		if err != nil {
			return "", err
		}
	}

	return string(bytes), nil
}
