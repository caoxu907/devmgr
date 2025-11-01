package main

import (
	_ "devmgr/internal/logic"

	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"go.lanniu.top/nebula-cloud/go-scaffold/sflake"
	"go.lanniu.top/nebula-cloud/go-scaffold/sframe"

	"github.com/gogf/gf/v2/os/gctx"

	"devmgr/internal/cmd"
)

func main() {
	sframe.Init()
	sflake.Init()
	cmd.Main.Run(gctx.GetInitCtx())
}
