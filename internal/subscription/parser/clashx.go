// internal/subscription/parser/clashx.go
package parser

import (
	"goconverter/internal/subscription/model"
	"strconv"
	"strings"
)

type ClashXParser struct{}

func NewClashXParser() *ClashXParser {
	return &ClashXParser{}
}

// ClashConfig Clash 配置
type ClashConfig struct {
	Proxies     []*Proxy     `yaml:"proxies"`      // 代理服务器列表
	ProxyGroups []ProxyGroup `yaml:"proxy-groups"` // 代理分组
	Rules       []string     `yaml:"rules"`        // 路由规则
	DNSSettings DNSConfig    `yaml:"dns"`          // DNS 配置
	AllowLan    bool         `yaml:"allow-lan"`    // 是否允许局域网连接
	Mode        string       `yaml:"mode"`         // 运行模式：rule/global/direct
	LogLevel    string       `yaml:"log-level"`    // 日志级别
	ExternalUI  string       `yaml:"external-ui"`  // 外部UI路径
	Secret      string       `yaml:"secret"`       // API密钥
}

// Proxy 定义单个代理服务器配置
type Proxy struct {
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

// ProxyGroup 定义代理分组配置
type ProxyGroup struct {
	Name    string   `yaml:"name"`          // 分组名称
	Type    string   `yaml:"type"`          // 分组类型：select/url-test/fallback/load-balance
	Proxies []string `yaml:"proxies"`       // 该组包含的代理列表
	Use     []string `yaml:"use,omitempty"` // 引用的代理集合

	// 自动测试相关配置
	URL       string `yaml:"url,omitempty"`       // 测试地址
	Interval  int32  `yaml:"interval,omitempty"`  // 测试间隔(秒)
	Tolerance uint16 `yaml:"tolerance,omitempty"` // 延迟容差
	Timeout   int32  `yaml:"timeout,omitempty"`   // 超时时间(秒)

	// 负载均衡配置
	Strategy string `yaml:"strategy,omitempty"` // 负载均衡策略
	Disable  bool   `yaml:"disable,omitempty"`  // 是否禁用
}

// DNSConfig 定义DNS配置
type DNSConfig struct {
	Enable         bool     `yaml:"enable"`             // 是否启用DNS服务器
	IPv6           bool     `yaml:"ipv6"`               // 是否解析IPv6地址
	NameServers    []string `yaml:"nameservers"`        // DNS服务器列表
	Fallback       []string `yaml:"fallback,omitempty"` // 备用DNS服务器
	FallbackFilter struct {
		GeoIP      bool     `yaml:"geoip"`
		IPCIDRList []string `yaml:"ipcidr"`
	} `yaml:"fallback-filter"`
	DefaultNameserver []string `yaml:"default-nameserver"` // 默认DNS服务器
	EnhancedMode      string   `yaml:"enhanced-mode"`      // 增强模式：fake-ip/redir-host
	FakeIPRange       string   `yaml:"fake-ip-range"`      // Fake IP地址范围
	Listen            string   `yaml:"listen"`             // 监听地址
}

func (p *ClashXParser) Match(link string) bool {
	return strings.HasPrefix(link, "vmess://")
}

func (p *ClashXParser) Parse(proxy *Proxy) (*model.Node, error) {

	settings := map[string]string{
		"uuid":             proxy.UUID,
		"alterId":          "",
		"network":          proxy.Network,
		"tls":              strconv.FormatBool(proxy.TLS),
		"skip-cert-verify": strconv.FormatBool(proxy.SkipCertVerify),
	}

	return &model.Node{
		Type:     model.NodeType(proxy.Type),
		Name:     proxy.Name,
		Server:   proxy.Server,
		Port:     proxy.Port,
		Password: proxy.Password,
		Protocol: "",
		Settings: settings,
		UDP:      proxy.UDP,
		ALPN:     proxy.Alpn,
	}, nil
}
