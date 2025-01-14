// internal/subscription/parser/base.go
package parser

import (
	"errors"
	"fmt"
	"goconverter/internal/subscription/model"
	"log"
	"strings"

	"github.com/goccy/go-yaml"
)

// Parser 定义解析器接口
type Parser interface {
	Parse(link string) (*model.Node, error)
	Match(link string) bool
}

// ParseSubscription 解析整个订阅内容
func ParseSubscription(content string, format string) ([]*model.Node, error) {
	if content == "" {
		return nil, errors.New("empty subscription content")
	}
	nodes := make([]*model.Node, 0)

	if format == "line" {
		// 分割成单独的节点链接
		links := strings.Split(content, "\n")

		// 所有支持的解析器
		parsers := []Parser{
			NewShadowsocksParser(),
			NewShadowsocksRParser(),
			NewVmessParser(),
			NewTrojanParser(),
		}

		var parseErrors []error

		for _, link := range links {
			link = strings.TrimSpace(link)
			if link == "" {
				continue
			}

			node, err := parseLink(link, parsers)
			if err != nil {
				parseErrors = append(parseErrors, fmt.Errorf("parse link %s: %w", link, err))
				continue
			}
			if node != nil {
				nodes = append(nodes, node)
			}
		}

		if len(parseErrors) > 0 {
			return nodes, fmt.Errorf("some links failed to parse: %v", parseErrors)
		}

		return nodes, nil
	} else if format == "clashx" {
		clashConfig := ClashConfig{}

		err := yaml.Unmarshal([]byte(content), &clashConfig)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		parser := NewClashXParser()
		for _, proxy := range clashConfig.Proxies {
			node, _ := parser.Parse(proxy)
			if node != nil {
				node.AllowInsecure = proxy.SkipCertVerify
				node.UDP = true
				nodes = append(nodes, node)
			}
		}
	}
	return nodes, fmt.Errorf("unexpect format: %s", format)
}

// 将解析单个链接的逻辑抽取成独立函数
func parseLink(link string, parsers []Parser) (*model.Node, error) {
	for _, parser := range parsers {
		if parser.Match(link) {
			return parser.Parse(link)
		}
	}
	return nil, fmt.Errorf("no matching parser found for link: %s", link)
}
