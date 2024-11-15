// internal/config/config_test.go
package config

import (
	"os"
	"testing"
)

func TestParseConfig(t *testing.T) {

	testConfig, err := os.ReadFile("D:\\Projects\\goconverter\\test\\data\\ACL4SSR.ini")
	if err != nil {
		t.Fatalf("读取配置文件失败: %v", err)
	}

	cfg, err := ParseConfig(testConfig)
	if err != nil {
		t.Fatalf("解析配置失败: %v", err)
	}

	t.Logf("Rulesets size: %d", len(cfg.RuleSets))
	t.Logf("ProxyGroups size: %d", len(cfg.ProxyGroups))

	// fmt.Printf("%+v\n", cfg)

}
