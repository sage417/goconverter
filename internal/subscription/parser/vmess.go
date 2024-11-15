// internal/subscription/parser/vmess.go
package parser

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"goconverter/internal/subscription/model"
	"strings"
)

type VmessParser struct{}

func NewVmessParser() *VmessParser {
	return &VmessParser{}
}

type VmessLink struct {
	Version string `json:"v"`
	Name    string `json:"ps"`
	Address string `json:"add"`
	Port    int    `json:"port"`
	UUID    string `json:"id"`
	AlterId int    `json:"aid"`
	Network string `json:"net"`
	Type    string `json:"type"`
	TLS     string `json:"tls"`
}

func (p *VmessParser) Match(link string) bool {
	return strings.HasPrefix(link, "vmess://")
}

func (p *VmessParser) Parse(link string) (*model.Node, error) {
	link = strings.TrimPrefix(link, "vmess://")

	decoded, err := base64.RawURLEncoding.DecodeString(link)
	if err != nil {
		return nil, err
	}

	var vmessLink VmessLink
	if err := json.Unmarshal(decoded, &vmessLink); err != nil {
		return nil, err
	}

	settings := map[string]string{
		"uuid":    vmessLink.UUID,
		"alterId": fmt.Sprintf("%d", vmessLink.AlterId),
		"network": vmessLink.Network,
		"type":    vmessLink.Type,
		"tls":     vmessLink.TLS,
	}

	return &model.Node{
		Type:     "vmess",
		Name:     vmessLink.Name,
		Server:   vmessLink.Address,
		Port:     vmessLink.Port,
		Protocol: "vmess",
		Settings: settings,
	}, nil
}
