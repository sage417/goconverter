// internal/subscription/converter/base.go
package converter

import (
	"goconverter/internal/subscription/model"
)

// Converter 定义转换器接口
type Converter interface {
	// Convert 将节点列表转换为目标格式的配置字符串
	Convert(nodes []*model.Node) (string, error)
	// ConvertNode 转换单个节点配置
	ConvertNode(node *model.Node) (interface{}, error)
}

// BaseInfo 基础转换信息
type BaseInfo struct {
	Name        string            // 配置名称
	Description string            // 配置描述
	Author      string            // 作者信息
	Tags        []string          // 标签
	Rules       map[string]string // 规则配置
}
