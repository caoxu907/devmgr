package devmgr

import (
	v1 "devmgr/api/devmgr/v1"
	"encoding/json"
	"os"
	"testing"
	"time"

	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/contrib/registry/nacos/v2"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func PrintJSON(t *testing.T, v interface{}) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Logf("succeeded: %v", v)
		return
	}
	t.Logf("%s", data)
}

func register() (client v1.DevmgrClient) {
	grpcx.Resolver.Register(nacos.New("127.0.0.1:18848"))
	var conn = grpcx.Client.MustNewGrpcClientConn("devmgr", grpcx.Balancer.WithRandom())
	client = v1.NewDevmgrClient(conn)
	return
}

func TestRedis_key(t *testing.T) {
	var ctx = gctx.New()

	// set key
	_, err := g.Redis().Set(ctx, "key_test", "value_test")
	if err != nil {
		g.Log().Fatal(ctx, err)
	}
	value, err := g.Redis().Get(ctx, "key_test")
	if err != nil {
		g.Log().Fatal(ctx, err)
	}
	g.Log().Line().Info(ctx, value.String())
}

func TestRedis_hash(t *testing.T) {
	var ctx = gctx.New()

	// set hash map
	key := "key_test_hash"
	_, err := g.Redis().HSet(ctx, key, g.Map{
		"id":   11,
		"name": "gdy",
	})
	if err != nil {
		g.Log().Fatal(ctx, err)
	}

	// retrieve hash map
	value, err := g.Redis().HGetAll(ctx, key)
	if err != nil {
		g.Log().Fatal(ctx, err)
	}
	g.Log().Line().Info(ctx, value.Map())

	// scan to struct
	type User struct {
		Id   uint64
		Name string
	}
	var user *User
	if err = value.Scan(&user); err != nil {
		g.Log().Fatal(ctx, err)
	}
	g.Dump(user)
}

func TestControllerV1_DevAuth(t *testing.T) {
	var ctx = gctx.New()
	client := register()

	req := v1.DevAuthReq{
		MachineCode:  "device_key_1",
		DeviceSecret: "device_secret_1",
	}
	res, err := client.DevAuth(ctx, &req)
	if err != nil {
		t.Fatal("failed:", err)
	}

	PrintJSON(t, res)
}

func TestControllerV1_ProductGetList(t *testing.T) {
	var ctx = gctx.New()
	client := register()

	res, err := client.ProductGetList(ctx, &v1.ProductGetListReq{})
	if err != nil {
		t.Fatal("failed:", err)
	}

	PrintJSON(t, res)
}

func TestControllerV1_ProductCreate(t *testing.T) {
	var ctx = gctx.New()
	client := register()

	res, err := client.ProductCreate(ctx, &v1.ProductCreateReq{
		ProductKey:  "product_key_1113",
		ProductDesc: "product_desc_111",
	})
	if err != nil {
		t.Fatal("failed:", err)
	}

	PrintJSON(t, res)
}

func TestControllerV1_DeviceGetList(t *testing.T) {
	var ctx = gctx.New()
	client := register()

	res, err := client.DeviceGetList(ctx, &v1.DeviceGetListReq{
		ProductId: 0,
		PageSize:  50,
		Status:    1,
	})
	if err != nil {
		t.Fatal("failed:", err)
	}

	PrintJSON(t, res)
}

func TestControllerV1_DeviceCreate(t *testing.T) {
	var ctx = gctx.New()
	client := register()

	res, err := client.DeviceCreate(ctx, &v1.DeviceCreateReq{
		ProductId:    2,
		DeviceDesc:   "product_" + time.Now().Format("20060102150405"),
		DeviceSecret: "product_" + time.Now().Format("20060102150405"),
	})
	if err != nil {
		t.Fatal("failed:", err)
	}

	PrintJSON(t, res)
}

func TestControllerV1_DeviceOverview(t *testing.T) {
	var ctx = gctx.New()
	client := register()

	res, err := client.DeviceOverview(ctx, &v1.DeviceOverviewReq{})
	if err != nil {
		t.Fatal("failed:", err)
	}

	PrintJSON(t, res)
}

func TestControllerV1_DeviceTemplateGet(t *testing.T) {
	var ctx = gctx.New()
	client := register()

	res, err := client.DeviceTemplateGet(ctx, &v1.DeviceTemplateGetReq{})
	if err != nil {
		t.Fatal("failed:", err)
	}

	//PrintJSON(t, res)

	err = os.WriteFile("/home/shr/temp/output.xlsx", res.Content, 0644)
	if err != nil {
		t.Fatal("保存文件失败:", err)
	}
}
