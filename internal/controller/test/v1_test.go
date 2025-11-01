package test

import (
	"testing"

	"github.com/gogf/gf/contrib/registry/nacos/v2"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gctx"
)

func register() (client *gclient.Client) {
	gsvc.SetRegistry(nacos.New("127.0.0.1:18848"))
	client = g.Client()
	client.SetDiscovery(gsvc.GetRegistry())
	return
}

type CodeResp struct {
	Code int `json:"code"`
}

//func AnalysisRes(res string, t *testing.T) {
//	var cr CodeResp
//	err := json.Unmarshal([]byte(res), &cr)
//	if err != nil {
//		t.Fatal("解析返回值失败:", err)
//	}
//
//	if 200 != cr.Code {
//		t.Fatal("执行失败, 错误码:", cr.Code)
//	} else {
//		t.Logf("执行成功")
//	}
//
//	var obj interface{}
//	err = json.Unmarshal([]byte(res), &obj)
//	if err != nil {
//		fmt.Println("解析失败:", err)
//		return
//	}
//	pretty, _ := json.MarshalIndent(obj, "", "  ")
//
//	g.Log().Line().Info(context.Background(), pretty)
//
//}

func TestControllerV1_ProductCreate(t *testing.T) {
	client := register()
	if client == nil {
		t.Fatal("client is nil")
		return
	}
	var ctx = gctx.New()

	result, err := client.Get(ctx, "http://devmgr-test/test/template")
	if err != nil {
		panic(err)
	}
	gfile.PutBytes("/home/shr/temp/cover.json", result.ReadAll())

}

func TestControllerV1_FileUpload(t *testing.T) {
	client := register()
	if client == nil {
		t.Fatal("client is nil")
		return
	}
	var ctx = gctx.New()

	path := "../../../resource/template/device_example.xlsx"
	result, err := client.Post(ctx, "http://devmgr-test/test/import2",
		"file=@file:"+path,
		"productId=2",
	)
	if err != nil {
		t.Fatal("失败:", err)
	}
	defer result.Close()
	t.Logf(result.ReadAllString())

}
