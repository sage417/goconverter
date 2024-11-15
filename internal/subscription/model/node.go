// internal/subscription/model/node.go
package model

import (
	"encoding/json"
	"strings"
)

// NodeType 定义节点类型
type NodeType string

const (
	TypeSS     NodeType = "ss"
	TypeSSR    NodeType = "ssr"
	TypeVmess  NodeType = "vmess"
	TypeTrojan NodeType = "trojan"
)

// Node 定义统一的节点结构
type Node struct {
	Type       NodeType          `json:"type"`
	Name       string            `json:"name"`
	Server     string            `json:"server"`
	Port       int               `json:"port"`
	Password   string            `json:"password"`
	Cipher     string            `json:"cipher"`      // 加密方法
	UDP        bool              `json:"udp"`         // 是否支持UDP
	Plugin     string            `json:"plugin"`      // SS插件
	PluginOpts map[string]string `json:"plugin_opts"` // SS插件参数

	// SSR特定参数
	Protocol      string `json:"protocol"`       // SSR协议
	ProtocolParam string `json:"protocol_param"` // SSR协议参数
	Obfs          string `json:"obfs"`           // SSR混淆
	ObfsParam     string `json:"obfs_param"`     // SSR混淆参数

	// VMess特定参数
	UUID      string            `json:"uuid"`       // VMess UUID
	AlterID   int               `json:"alter_id"`   // VMess AlterID
	Network   string            `json:"network"`    // 传输协议
	TLS       bool              `json:"tls"`        // 是否启用TLS
	ALPN      []string          `json:"alpn"`       // ALPN
	SNI       string            `json:"sni"`        // TLS SNI
	WsPath    string            `json:"ws_path"`    // WebSocket路径
	WsHeaders map[string]string `json:"ws_headers"` // WebSocket请求头

	// Trojan特定参数
	AllowInsecure bool `json:"allow_insecure"` // 是否允许不安全TLS

	// 通用参数
	Group    string   `json:"group"` // 分组
	Tags     []string `json:"tags"`  // 标签
	Settings map[string]string
}

// ToClash 将节点转换为Clash配置
func (n *Node) ToClash() map[string]interface{} {
	proxy := make(map[string]interface{})
	proxy["name"] = n.Name
	proxy["server"] = n.Server
	proxy["port"] = n.Port

	switch n.Type {
	case TypeSS:
		proxy["type"] = "ss"
		proxy["cipher"] = n.Cipher
		proxy["password"] = n.Password
		if n.Plugin != "" {
			proxy["plugin"] = n.Plugin
			proxy["plugin-opts"] = n.PluginOpts
		}

	case TypeSSR:
		proxy["type"] = "ssr"
		proxy["cipher"] = n.Cipher
		proxy["password"] = n.Password
		proxy["protocol"] = n.Protocol
		proxy["protocol-param"] = n.ProtocolParam
		proxy["obfs"] = n.Obfs
		proxy["obfs-param"] = n.ObfsParam

	case TypeVmess:
		proxy["type"] = "vmess"
		proxy["uuid"] = n.UUID
		proxy["alterId"] = n.AlterID
		proxy["cipher"] = defaultIfEmpty(n.Cipher, "auto")
		if n.TLS {
			proxy["tls"] = true
			proxy["servername"] = defaultIfEmpty(n.SNI, n.Server)
			if len(n.ALPN) > 0 {
				proxy["alpn"] = n.ALPN
			}
		}
		if n.Network == "ws" {
			proxy["network"] = n.Network
			wsOpts := make(map[string]interface{})
			wsOpts["path"] = defaultIfEmpty(n.WsPath, "/")
			if len(n.WsHeaders) > 0 {
				wsOpts["headers"] = n.WsHeaders
			}
			proxy["ws-opts"] = wsOpts
		}

	case TypeTrojan:
		proxy["type"] = "trojan"
		proxy["password"] = n.Password

		proxy["sni"] = defaultIfEmpty(n.SNI, n.Server)
		proxy["skip-cert-verify"] = n.AllowInsecure
		if len(n.ALPN) > 0 {
			proxy["alpn"] = n.ALPN
		}

	}

	if n.UDP {
		proxy["udp"] = true
	}

	return proxy
}

// ToJSON 将节点转换为JSON字符串
func (n *Node) ToJSON() (string, error) {
	data, err := json.Marshal(n)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// Clone 深拷贝节点
func (n *Node) Clone() *Node {
	clone := &Node{}
	data, _ := json.Marshal(n)
	_ = json.Unmarshal(data, clone)
	return clone
}

// defaultIfEmpty 如果字符串为空则返回默认值
func defaultIfEmpty(str, def string) string {
	if strings.TrimSpace(str) == "" {
		return def
	}
	return str
}

// ValidateNode 验证节点配置的完整性
func ValidateNode(node *Node) error {
	if node.Server == "" {
		return NewValidationError("server address is required")
	}
	if node.Port <= 0 || node.Port > 65535 {
		return NewValidationError("invalid port number")
	}
	if node.Password == "" && node.Type != TypeVmess {
		return NewValidationError("password is required")
	}

	switch node.Type {
	case TypeSS:
		if node.Cipher == "" {
			return NewValidationError("cipher is required for shadowsocks")
		}
	case TypeSSR:
		if node.Protocol == "" {
			return NewValidationError("protocol is required for shadowsocksr")
		}
		if node.Obfs == "" {
			return NewValidationError("obfs is required for shadowsocksr")
		}
	case TypeVmess:
		if node.UUID == "" {
			return NewValidationError("uuid is required for vmess")
		}
		if node.AlterID < 0 {
			return NewValidationError("invalid alter_id for vmess")
		}
	}

	return nil
}

// ValidationError 验证错误
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func NewValidationError(message string) error {
	return &ValidationError{Message: message}
}
