// cmd/converter/main.go
package main

import (
	"flag"
	"fmt"
	"goconverter/internal/config"
	"goconverter/internal/converter"
	"goconverter/internal/fetcher"
	"goconverter/internal/subscription/parser"
	"log"
	"os"
)

func main() {
	// 定义命令行参数
	subscriptionURL := flag.String("url", "", "订阅地址URL")
	configURL := flag.String("config", "https://raw.githubusercontent.com/ACL4SSR/ACL4SSR/refs/heads/master/Clash/config/ACL4SSR.ini", "配置文件URL")
	outputFile := flag.String("output", "", "输出文件路径(可选)")
	// targetFormat := flag.String("target", "clash", "目标格式(clash/surge/quantumult)")

	flag.Parse()

	if *subscriptionURL == "" {
		log.Fatal("订阅地址不能为空")
	}

	contentFetcher := fetcher.NewFetcher()
	configBytes, err := contentFetcher.Fetch(*configURL)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	cfg, err := config.ParseConfig(configBytes)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	subscriptionBytes, err := contentFetcher.Fetch(*subscriptionURL)
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	nodes, err := parser.ParseSubscription(string(subscriptionBytes), "clashx")

	conv := converter.NewClashConverter(&converter.BaseInfo{})

	result, err := conv.Convert(nodes, cfg)
	if err != nil {
		log.Fatalf("转换失败: %v", err)
	}

	// 输出结果
	if *outputFile != "" {
		err = os.WriteFile(*outputFile, []byte(result), 0644)
		if err != nil {
			log.Fatalf("写入文件失败: %v", err)
		}
		fmt.Printf("已保存到文件: %s\n", *outputFile)
	} else {
		fmt.Println(result)
	}
}
