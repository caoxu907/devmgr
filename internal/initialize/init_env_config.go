package initialize

import (
	"context"

	"github.com/gogf/gf/v2/encoding/gyaml"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/genv"
	"github.com/gogf/gf/v2/text/gregex"
)

func InitEnvConfig(Ctx context.Context) {
	config, err := g.Cfg().Data(Ctx)
	if err != nil {
		g.Log().Fatal(Ctx, err)
	}
	envData := getEnvConfig(config)
	yamlData, err := gyaml.Encode(envData)
	if err != nil {
		g.Log().Fatal(Ctx, err)
	}
	adapter, err := gcfg.NewAdapterContent(string(yamlData))
	if err != nil {
		g.Log().Fatal(Ctx, err)
	}
	g.Cfg().SetAdapter(adapter)
	g.Log().Line().Info(Ctx, "env config init success")
}

func getEnvConfig(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{})
		for key, value := range v {
			result[key] = getEnvConfig(value)
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = getEnvConfig(item)
		}
		return result
	case string:
		return getEnvValue(v)
	default:
		return v
	}
}

func getEnvValue(s string) string {
	pattern := `\$\{([^}:|]+)[:|]([^}]+)\}`
	value, err := gregex.ReplaceStringFuncMatch(pattern, s, func(match []string) string {
		if len(match) < 3 {
			return match[0]
		}
		key := match[1]
		def := match[2]
		value := genv.Get(key, def).String()
		return value
	})
	if err != nil {
		return s
	}
	return value
}
