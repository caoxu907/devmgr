package tool

import (
	"fmt"
	"net"
)

// GetIPsInSubnet 获取所有网卡中属于指定网段的IP地址
// subnet: CIDR格式的网段，如 "10.0.2.0/24", "172.18.0.0/16"
// 返回: IP地址字符串切片和错误
func GetIPsInSubnet(subnet string) ([]string, error) {
	// 解析网段
	_, ipNet, err := net.ParseCIDR(subnet)
	if err != nil {
		return nil, fmt.Errorf("解析网段失败: %v", err)
	}

	var result []string

	// 获取所有网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("获取网络接口失败: %v", err)
	}

	// 遍历所有网络接口
	for _, iface := range interfaces {
		// 跳过未启用或loopback接口
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		// 获取接口的所有地址
		addrs, err := iface.Addrs()
		if err != nil {
			continue // 跳过获取地址失败的接口
		}

		// 检查每个地址
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				// 检查IP是否在指定网段内
				if ipNet.Contains(ipnet.IP) {
					result = append(result, ipnet.IP.String())
				}
			}
		}
	}

	return result, nil
}

// GetFirstIPInSubnet 获取指定网段中的第一个IP地址（常用于单个IP场景）
func GetFirstIPInSubnet(subnet string) (string, error) {
	ips, err := GetIPsInSubnet(subnet)
	if err != nil {
		return "", err
	}

	if len(ips) == 0 {
		return "", fmt.Errorf("在网段 %s 中未找到任何IP地址", subnet)
	}

	return ips[0], nil
}
