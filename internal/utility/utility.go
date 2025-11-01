package utility

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

func Gtime2Str(_gtime *gtime.Time) string {
	if _gtime != nil && !_gtime.IsZero() {
		// 明确设置为中国/北京时区 (UTC+8)
		cst := time.FixedZone("CST", 8*3600)
		return _gtime.ToLocation(cst).Format("Y-m-d H:i:s")
	}
	return ""
}

// parseContent 将 content 序列化后反序列化到 out（out 必须是指针）
func ParseContent(ctx context.Context, content interface{}, out interface{}, errInfo string) (err error) {
	defer func() {
		if err != nil {
			g.Log().Line().Error(ctx, errInfo)
		}
	}()

	contentBytes, err := json.Marshal(content)
	if err != nil {
		g.Log().Line().Error(ctx, err.Error())
		return gerror.New("序列化设备消息内容失败")
	}
	if err = json.Unmarshal(contentBytes, out); err != nil {
		g.Log().Line().Error(ctx, "解析设备消息JSON失败:", err.Error())
		return gerror.New("解析设备消息JSON失败")
	}
	return nil
}
