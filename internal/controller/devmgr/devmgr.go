package devmgr

import (
	"context"
	v1 "devmgr/api/devmgr/v1"
	filemgr "devmgr/api/filemgr/v1"
	"devmgr/internal/consts"
	"devmgr/internal/dao"
	"devmgr/internal/externalapi"
	"devmgr/internal/model"
	"devmgr/internal/model/do"
	"devmgr/internal/model/entity"
	"devmgr/internal/service"
	"devmgr/internal/utility"
	"encoding/json"
	"os"
	"strconv"
	"time"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/errors/gcode"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/guid"
)

type Controller struct {
	v1.UnimplementedDevmgrServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterDevmgrServer(s.Server, &Controller{})
}

func (*Controller) DevActivate(ctx context.Context, req *v1.DevActivateReq) (res *v1.DevActivateRes, err error) {
	var device []entity.Device
	err = dao.Device.Ctx(ctx).Where(do.Device{
		MachineCode:  req.MachineCode,
		DeviceSecret: req.DeviceSecret,
		Deleted:      false, //删除的无法激活
	}).Scan(&device)
	if err != nil {
		g.Log().Line().Error(ctx, "查询设备失败", err.Error())
		return nil, err
	}

	if len(device) != 1 {
		g.Log().Line().Error(ctx, "设备不存在或重复", len(device))
		return nil, gerror.NewCode(gcode.CodeInternalError, "设备不存在或重复", strconv.Itoa(len(device)))
	}

	if device[0].TenantId != 0 {
		g.Log().Line().Warning(ctx, "设备已被激活", len(device))
		if int(req.TenantId) != device[0].TenantId {
			return &v1.DevActivateRes{}, nil
		}
		return nil, gerror.NewCode(gcode.CodeInternalError, "设备已被激活")
	}

	//更新设备信息
	_, err = dao.Device.Ctx(ctx).Data(do.Device{
		TenantId: req.TenantId,
	}).WherePri(device[0].Id).Update()
	if err != nil {
		g.Log().Line().Error(ctx, "更新设备信息失败", err.Error())
		return nil, err
	}

	return &v1.DevActivateRes{}, nil
}

func (*Controller) DevAuth(ctx context.Context, req *v1.DevAuthReq) (res *v1.DevAuthRes, err error) {
	var devices []entity.Device

	err = dao.Device.Ctx(ctx).
		Where(do.Device{
			MachineCode:  req.MachineCode,
			DeviceSecret: req.DeviceSecret,
			Deleted:      false, //删除的无法认证
		}).
		Scan(&devices)

	if err != nil {
		g.Log().Line().Error(ctx, err.Error())
		return nil, err
	}
	if len(devices) != 1 {
		g.Log().Line().Error(ctx, "认证失败", len(devices))
		return nil, gerror.NewCode(gcode.CodeInternalError)
	}

	return &v1.DevAuthRes{DeviceKey: devices[0].DeviceKey}, nil
}

func (*Controller) ProductGetList(ctx context.Context, req *v1.ProductGetListReq) (res *v1.ProductGetListRes, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 查询数据库
	m := dao.Product.Ctx(ctx).Where(do.Product{Deleted: false})

	if req.ProductDesc != "" {
		m = m.Where("product_desc LIKE ?", "%"+req.ProductDesc+"%")
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}

	var products []*entity.Product
	err = m.Page(int(req.Page), int(req.PageSize)).Scan(&products)
	if err != nil {
		return nil, err
	}

	// 转换为响应结构
	res = &v1.ProductGetListRes{
		Total: int32(int64(total)),
	}
	for _, p := range products {
		res.List = append(res.List, &v1.ProductInfo{
			Id:          int32(p.Id),
			ProductKey:  p.ProductKey,
			ProductDesc: p.ProductDesc,
			Deleted:     p.Deleted,
		})
	}

	return
}

func (*Controller) ProductCreate(ctx context.Context, req *v1.ProductCreateReq) (res *v1.ProductCreateRes, err error) {
	_, err = dao.Product.Ctx(ctx).Data(do.Product{
		ProductKey:  req.ProductKey,
		ProductDesc: req.ProductDesc,
	}).Insert()

	////转为sql语句
	//sqlFun := func(ctx context.Context) error {
	//	_, err = dao.Product.Ctx(ctx).Data(do.Product{
	//		ProductKey:  req.ProductKey,
	//		ProductDesc: req.ProductDesc,
	//	}).Insert()
	//	return err
	//}
	//sql, err := gdb.ToSQL(ctx, sqlFun)
	//if err != nil {
	//	return nil, err
	//}

	if err != nil {
		return nil, err
	}

	return
}

func (*Controller) ProductDelete(ctx context.Context, req *v1.ProductDeleteReq) (res *v1.ProductDeleteRes, err error) {
	_, err = dao.Product.Ctx(ctx).WherePri(req.Ids).Data(do.Product{Deleted: true}).Update()

	//todo 通知云端代理断开和删除相关设备连接

	return
}

func (*Controller) ProductUpdate(ctx context.Context, req *v1.ProductUpdateReq) (res *v1.ProductUpdateRes, err error) {
	_, err = dao.Product.Ctx(ctx).Data(do.Product{
		ProductKey:  req.ProductKey,
		ProductDesc: req.ProductDesc,
	}).WherePri(req.Id).Update()
	return
}

func (*Controller) DeviceGetList(ctx context.Context, req *v1.DeviceGetListReq) (res *v1.DeviceGetListRes, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	var listDevice []*entity.Device
	res = &v1.DeviceGetListRes{}
	m := dao.Device.Ctx(ctx).Where(do.Device{Deleted: false})

	if req.ProductId != 0 {
		m = m.Where("product_id = ?", req.ProductId)
	}

	if req.DeviceDesc != "" {
		m = m.Where("device_desc LIKE ?", "%"+req.DeviceDesc+"%")
	}

	if 0 != req.Status {
		m = m.Where("status = ?", req.Status)
	}

	total, err := m.Count()
	res.Total = int32(total)
	if err != nil {
		return nil, err
	}
	err = m.Page(int(req.Page), int(req.PageSize)).Scan(&listDevice)
	if err != nil {
		return nil, err
	}

	productDescMap := make(map[int]string)

	// 转换为响应结构
	for _, p := range listDevice {

		// 如果产品描述未缓存，则查询并缓存
		if _, exists := productDescMap[p.ProductId]; !exists {
			var product entity.Product
			err = dao.Product.Ctx(ctx).WherePri(p.ProductId).Scan(&product)
			if err != nil {
				return nil, gerror.Newf("查询产品信息失败: %v", err)
			}

			// 如果产品已被删除，则跳过
			if product.Deleted {
				continue
			}

			productDescMap[p.ProductId] = product.ProductDesc

		}
		// 使用缓存的产品描述
		productDescTemp := productDescMap[p.ProductId]

		res.List = append(res.List, &v1.DeviceInfo{
			Id:             int32(p.Id),
			ProductId:      int32(p.ProductId),
			ProductDesc:    productDescTemp,
			TenantId:       int32(p.TenantId),
			DeviceKey:      p.DeviceKey,
			DeviceDesc:     p.DeviceDesc,
			DeviceSecret:   p.DeviceSecret,
			MachineCode:    p.MachineCode,
			Status:         int32(p.Status),
			Version:        p.Version,
			Address:        p.Address,
			CreateTime:     utility.Gtime2Str(p.CreateTime),
			ActivateTime:   utility.Gtime2Str(p.ActivateTime),
			LastOnlineTime: utility.Gtime2Str(p.LastOnlineTime),
			Deleted:        p.Deleted,
		})
	}

	return
}

func (*Controller) DeviceCreate(ctx context.Context, req *v1.DeviceCreateReq) (res *v1.DeviceCreateRes, err error) {
	if 0 == req.ProductId {
		return nil, gerror.New("ProductId字段不能为空")
	}

	DeviceKey := guid.S()

	_, err = dao.Device.Ctx(ctx).Data(do.Device{
		ProductId:    req.ProductId,
		DeviceKey:    DeviceKey,
		DeviceDesc:   req.DeviceDesc,
		DeviceSecret: req.DeviceSecret,
	}).Insert()
	if err != nil {
		return nil, err
	}
	return
}

func (*Controller) DeviceDelete(ctx context.Context, req *v1.DeviceDeleteReq) (res *v1.DeviceDeleteRes, err error) {
	_, err = dao.Device.Ctx(ctx).WherePri(req.Ids).Data(do.Device{Deleted: true}).Update()
	if err != nil {
		return nil, err
	}
	return
}

func (*Controller) DeviceUpdate(ctx context.Context, req *v1.DeviceUpdateReq) (res *v1.DeviceUpdateRes, err error) {
	_, err = dao.Device.Ctx(ctx).Data(do.Device{
		DeviceDesc:   req.DeviceDesc,
		DeviceSecret: req.DeviceSecret,
	}).WherePri(req.Id).Update()
	if err != nil {
		return nil, err
	}
	return
}

func (*Controller) DeviceOverview(ctx context.Context, req *v1.DeviceOverviewReq) (res *v1.DeviceOverviewRes, err error) {
	//查询设备总数、在线数、激活数，Status为1或者2表示已经激活，为2表示在线
	var total, online, active int
	total, err = dao.Device.Ctx(ctx).Count(do.Device{Deleted: false})
	if err != nil {
		return nil, gerror.Newf("查询设备总数失败: %v", err)
	}
	active, err = dao.Device.Ctx(ctx).Where(do.Device{Deleted: false, Status: []int32{consts.DeviceStatusOffline, consts.DeviceStatusOnline}}).Count()
	if err != nil {
		return nil, gerror.Newf("查询激活设备数失败: %v", err)
	}
	online, err = dao.Device.Ctx(ctx).Where(do.Device{Deleted: false, Status: consts.DeviceStatusOnline}).Count()
	if err != nil {
		return nil, gerror.Newf("查询在线设备数失败: %v", err)
	}
	res = &v1.DeviceOverviewRes{
		Total:     int32(int64(total)),
		Online:    int32(int64(online)),
		Activated: int32(int64(active)),
	}

	return
}

func (*Controller) DeviceImport(ctx context.Context, req *v1.DeviceImportReq) (res *v1.DeviceImportRes, err error) {
	if req.ProductId <= 0 {
		return nil, gerror.New("ProductId字段不能为空")
	}

	//解析bytes File为excel，读取设备信息
	if req.File == nil {
		return nil, gerror.New("File字段不能为空")
	}

	// 解析Excel文件，获取设备信息
	devices, err := service.DevApi().ParseExcelFile(ctx, req.File)
	if err != nil {
		g.Log().Line().Error(ctx, "解析Excel文件失败")
		return nil, gerror.New("解析Excel文件失败")
	}

	g.Log().Debug(ctx, "解析Excel文件成功，设备数量:", len(devices))

	// 为所有设备添加产品ID和设备密钥
	for i := range devices {
		devices[i].ProductId = req.ProductId
		devices[i].DeviceKey = guid.S()
	}

	// 开启事务
	err = dao.Device.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := dao.Device.Ctx(ctx).TX(tx).Data(devices).Insert()
		if err != nil {
			return gerror.New("所有设备导入失败")
		}

		// 提交事务
		return nil
	})

	return &v1.DeviceImportRes{}, nil
}

func (*Controller) DeviceTemplateGet(ctx context.Context, req *v1.DeviceTemplateGetReq) (res *v1.DeviceTemplateGetRes, err error) {
	//将本目录下的device_template.xlsx文件字节流
	data, err := os.ReadFile("./resource/template/device_example.xlsx") // 读取文件为字节流
	if err != nil {
		g.Log().Line().Error(ctx, err.Error())
		return nil, err
	}

	return &v1.DeviceTemplateGetRes{
		Content: data,
	}, nil
}

func (*Controller) DeviceCommandList(ctx context.Context, req *v1.DeviceCommandListReq) (res *v1.DeviceCommandListRes, err error) {
	res = &v1.DeviceCommandListRes{}

	m := dao.CommandConfig.Ctx(ctx).Where("device_id=?", req.DeviceId)
	total, err := m.Count()
	res.Total = int32(int64(total))
	if err != nil {
		return nil, err
	}

	var commandConfigs []entity.CommandConfig
	err = m.Scan(&commandConfigs)
	if err != nil {
		g.Log().Line().Warning(ctx, err.Error())
		return nil, err
	}

	// 转换为响应结构
	for _, p := range commandConfigs {
		res.List = append(res.List, &v1.DeviceCommandListInfo{
			CommandId:   int32(p.Id),
			CommandName: p.CommandName,
			CommandDesc: p.CommandDesc,
		})
	}

	return res, nil
}

func (*Controller) DeviceCommandExecute(ctx context.Context, req *v1.DeviceCommandExecuteReq) (res *v1.DeviceCommandExecuteRes, err error) {
	if req.CommandId == 0 {
		return nil, gerror.New("CommandId字段不能为空")
	}

	//查询命令配置信息
	var commandCfg entity.CommandConfig
	err = dao.CommandConfig.Ctx(ctx).WherePri(req.CommandId).Scan(&commandCfg)
	if err != nil {
		g.Log().Line().Error(ctx, err.Error())
		return nil, err
	}

	//查询设备信息
	var device entity.Device
	err = dao.Device.Ctx(ctx).WherePri(commandCfg.DeviceId).Scan(&device)
	if err != nil {
		g.Log().Line().Error(ctx, err.Error())
		return nil, err
	}

	//设备必须在线才能下发命令
	if device.Status != int(consts.DeviceStatusOnline) {
		g.Log().Line().Error(ctx, "设备不在线，无法下发命令")
		return nil, gerror.New("设备不在线，无法下发命令")
	}

	sqlRes, err := dao.CommandRecord.Ctx(ctx).Data(do.CommandRecord{
		CommandConfigId: req.CommandId,
		Params:          req.Param,
		Status:          consts.Command_status_executing,
		TimeoutSeconds:  commandCfg.TimeoutSeconds,
	}).Insert()
	if err != nil {
		g.Log().Line().Error(ctx, err.Error())
		return nil, err
	}

	// 获取插入的记录ID
	execId, err := sqlRes.LastInsertId()
	if err != nil {
		g.Log().Line().Error(ctx, err.Error())
		return nil, err
	}

	//下发命令给设备
	msg := model.DeviceMessage{
		MessageType: consts.Edge_command_down,
		Content: model.DeviceMessageCommandSend{
			CommandId: int32(execId),
			Command:   commandCfg.CommandName,
			Param:     req.Param,
		},
	}
	msgData, err := json.Marshal(msg)
	if err != nil {
		g.Log().Line().Error(ctx, err.Error())
		return nil, err
	}
	g.Log().Line().Info(ctx, "下发命令给设备:", msgData)

	mqClient := service.MQClient()
	if mqClient == nil {
		g.Log().Line().Error(ctx, "MQClient未初始化")
		return nil, gerror.New("MQClient未初始化")
	}

	err = mqClient.Publish(consts.Kafka_topic_notify_down, []byte(device.DeviceKey), msgData)
	if err != nil {
		g.Log().Line().Error(ctx, err.Error())
		return nil, err
	}

	//记录日志
	_, _ = dao.CommandLog.Ctx(ctx).Data(do.CommandLog{CommandId: execId, Level: consts.Msg_level_info, Content: "已下发命令"}).Insert()

	return &v1.DeviceCommandExecuteRes{
		ExecuteId: int32(execId),
	}, nil
}

func (*Controller) DeviceCommandLogList(ctx context.Context, req *v1.DeviceCommandLogListReq) (res *v1.DeviceCommandLogListRes, err error) {
	if req.ExecuteId == 0 {
		g.Log().Line().Error(ctx, "ExecuteId字段不能为空")
		return nil, gerror.New("ExecuteId字段不能为空")
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	m := dao.CommandLog.Ctx(ctx).Where(do.CommandLog{CommandId: req.ExecuteId})

	total, err := m.Count()
	if err != nil {
		return nil, err
	}

	// 查询数据库结构
	var cmdLogs []entity.CommandLog
	err = m.Page(int(req.Page), int(req.PageSize)).Scan(&cmdLogs)
	if err != nil {
		return nil, err
	}

	var cmdRecord entity.CommandRecord
	err = dao.CommandRecord.Ctx(ctx).WherePri(req.ExecuteId).Scan(&cmdRecord)
	if err != nil {
		g.Log().Line().Error(ctx, err.Error())
		return nil, err
	}

	res = &v1.DeviceCommandLogListRes{
		Total:  int32(int64(total)),
		Status: int32(cmdRecord.Status),
	}

	// 转换为响应结构
	for _, p := range cmdLogs {
		res.List = append(res.List, &v1.CommandLog{
			Timestamp: utility.Gtime2Str(p.Timestamp),
			Content:   p.Content,
		})
	}

	return res, err
}

func (*Controller) PackageList(ctx context.Context, req *v1.PackageListReq) (res *v1.PackageListRes, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 查询数据库
	m := dao.UpgradePackage.Ctx(ctx)
	if req.PackageName != "" {
		m = m.Where("package_name LIKE ?", "%"+req.PackageName+"%")
	}
	if req.PackageDesc != "" {
		m = m.Where("package_desc LIKE ?", "%"+req.PackageDesc+"%")
	}
	if req.PackageVersion != "" {
		m = m.Where("version LIKE ?", "%"+req.PackageVersion+"%")
	}
	if req.ProductId != 0 {
		m = m.Where(do.UpgradePackage{ProductId: req.ProductId})
	}
	total, err := m.Count()
	if err != nil {
		return nil, err
	}
	var objs []*entity.UpgradePackage
	err = m.Page(int(req.Page), int(req.PageSize)).Scan(&objs)
	if err != nil {
		return nil, err
	}

	var productDescMap = make(map[int]string)

	// 转换为响应结构
	res = &v1.PackageListRes{
		Total: int32(int64(total)),
	}
	for _, p := range objs {

		// 如果产品描述未缓存，则查询并缓存
		if _, exists := productDescMap[p.ProductId]; !exists {
			var product entity.Product
			err = dao.Product.Ctx(ctx).WherePri(p.ProductId).Scan(&product)
			if err != nil {
				return nil, gerror.Newf("查询产品信息失败: %v", err)
			}

			// 如果产品已被删除，则跳过
			if product.Deleted {
				continue
			}

			productDescMap[p.ProductId] = product.ProductDesc

		}
		// 使用缓存的产品描述
		productDescTemp := productDescMap[p.ProductId]

		// p.AvailableVersions 是 string，内容为 JSON 数组，如 '["v1", "v2", "v3"]'
		var versions []string
		err := json.Unmarshal([]byte(p.AvailableVersions), &versions)
		if err != nil {
			// 处理错误
		}

		res.List = append(res.List, &v1.PackageInfo{
			Id:                 int32(p.Id),
			ProductId:          int32(p.ProductId),
			ProductDesc:        productDescTemp,
			PackageName:        p.PackageName,
			PackageDesc:        p.PackageDesc,
			PackageVersion:     p.Version,
			SupportAllVersions: p.SupportAllVersions,
			AvailableVersions:  versions,
			PackagePath:        p.PackagePath,
			PackageSize:        int32(p.PackageSize),
			Signature:          p.PackageHash,
			SignAlgorithm:      "md5",
			CreateTime:         utility.Gtime2Str(p.CreateTime),
			UpdateTime:         utility.Gtime2Str(p.CreateTime),
			PackageStatus:      int32(p.PackageStatus),
		})
	}

	return res, nil
}

func (*Controller) DeviceVersionList(ctx context.Context, req *v1.DeviceVersionListReq) (res *v1.DeviceVersionListRes, err error) {
	if req.ProductId == 0 {
		g.Log().Line().Error(ctx, "ProductId字段不能为空")
		return nil, gerror.New("ProductId字段不能为空")
	}
	// 查询数据库
	m := dao.Device.Ctx(ctx).Where(do.Device{ProductId: req.ProductId}).Group("version").Order("version DESC")
	var jsonVersions []string
	err = m.Fields("version").Scan(&jsonVersions)
	if err != nil {
		g.Log().Line().Error(ctx, "查询设备版本失败:", err.Error())
		return nil, err
	}

	// 解析JSON格式的版本号
	versions := make([]string, 0, len(jsonVersions))
	for _, jsonVersion := range jsonVersions {
		// 跳过空版本
		if jsonVersion == "" {
			continue
		}

		// 尝试解析JSON
		var versionObj struct {
			Version string `json:"version"`
		}
		if err := json.Unmarshal([]byte(jsonVersion), &versionObj); err == nil && versionObj.Version != "" {
			versions = append(versions, versionObj.Version)
		} else {
			// 如果不是JSON或解析失败，直接使用原始字符串
			versions = append(versions, jsonVersion)
		}
	}

	return &v1.DeviceVersionListRes{
		Total: int32(len(versions)),
		List:  versions,
	}, nil
}

func (*Controller) PackageCreate(ctx context.Context, req *v1.PackageCreateReq) (res *v1.PackageCreateRes, err error) {
	if req.ProductId == 0 {
		g.Log().Line().Warning(ctx, "ProductId字段不能为空")
		return nil, gerror.New("ProductId字段不能为空")
	}

	jsonStr, err := json.Marshal(req.AvailableVersions)
	if err != nil {
		g.Log().Line().Error(ctx, "序列化失败:", err.Error())
		return nil, err
	}

	PackagePath := strconv.FormatInt(int64(req.ProductId), 10) + "/" + req.PackageVersion + "/" + req.PackageFileName
	data := do.UpgradePackage{
		ProductId:          req.ProductId,
		Version:            req.PackageVersion,
		PackageName:        req.PackageName,
		PackageDesc:        req.PackageDesc,
		PackagePath:        PackagePath,
		SupportAllVersions: req.SupportAllVersions,
		AvailableVersions:  string(jsonStr),
		PackageStatus:      consts.Ota_package_status_draft,
	}

	// 检查是否已存在相同ProductId和Version且状态为draft的记录，如果存在则更新，否则插入
	result, err := dao.UpgradePackage.Ctx(ctx).Where(do.UpgradePackage{
		ProductId:     req.ProductId,
		Version:       req.PackageVersion,
		PackageStatus: consts.Ota_package_status_draft,
	}).Data(data).Update()

	if err != nil {
		g.Log().Line().Error(ctx, "更新升级包记录失败:", err.Error())
		return nil, err
	}

	// 检查是否有记录被更新
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		g.Log().Line().Error(ctx, "获取受影响行数失败:", err.Error())
		return nil, err
	}

	// 如果没有记录被更新，则尝试插入
	if rowsAffected == 0 {
		_, err = dao.UpgradePackage.Ctx(ctx).Data(data).Insert()
		if err != nil {
			// 如果插入失败，可能是因为其他线程刚好插入了记录
			// 再尝试一次更新
			_, err = dao.UpgradePackage.Ctx(ctx).Where(do.UpgradePackage{
				ProductId:     req.ProductId,
				Version:       req.PackageVersion,
				PackageStatus: consts.Ota_package_status_draft,
			}).Data(data).Update()
			if err != nil {
				g.Log().Line().Error(ctx, "创建或更新升级包记录失败:", err.Error())
				return nil, err
			}
		}
	}

	//调用文件服务，获取上传URL
	fileClient, err := externalapi.GetFilemgrClient()
	if err != nil {
		return nil, err
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, consts.Grpc_Timeout*time.Second)
	defer cancel()
	presignDownloadRes, err := fileClient.GeneratePresignedUploadURL(timeoutCtx, &filemgr.PresignUploadReq{
		Bucket:   consts.File_bucket,
		Key:      PackagePath,
		Expiry:   consts.File_upload_expiry,
		Intranet: false,
	})

	if err != nil {
		g.Log().Line().Error(ctx, "获取上传URL失败:", err.Error())
		return nil, err
	}

	return &v1.PackageCreateRes{
		Url:    presignDownloadRes.Url,
		Fields: presignDownloadRes.Fields,
	}, nil
}

func (*Controller) PackageConfirm(ctx context.Context, req *v1.PackageConfirmReq) (res *v1.PackageConfirmRes, err error) {
	var pkg entity.UpgradePackage
	err = dao.UpgradePackage.Ctx(ctx).Where(do.UpgradePackage{
		ProductId: req.ProductId,
		Version:   req.PackageVersion,
	}).Scan(&pkg)
	if err != nil {
		g.Log().Line().Error(ctx, "查询升级包失败:", err.Error())
		return nil, gerror.Newf("查询升级包失败:", err.Error())
	}

	//调用文件服务，获取文件元数据
	fileClient, err := externalapi.GetFilemgrClient()
	if err != nil {
		return nil, err
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, consts.Grpc_Timeout*time.Second)
	defer cancel()
	getMetaRes, err := fileClient.GetFileMetadata(timeoutCtx, &filemgr.GetFileMetadataReq{
		Bucket: consts.File_bucket,
		Key:    pkg.PackagePath,
	})
	if err != nil {
		g.Log().Line().Error(ctx, "获取文件元数据失败:", err.Error())
		return nil, gerror.Newf("获取文件元数据失败:", err.Error())
	}

	_, err = dao.UpgradePackage.Ctx(ctx).Data(do.UpgradePackage{
		PackageSize:   getMetaRes.Size,
		PackageHash:   getMetaRes.Etag,
		PackageStatus: consts.Ota_package_status_release,
	}).WherePri(pkg.Id).Update()
	if err != nil {
		g.Log().Line().Error(ctx, "更新升级包信息失败:", err.Error())
		return nil, err
	}

	return &v1.PackageConfirmRes{}, nil
}

func (*Controller) PackageDelete(ctx context.Context, req *v1.PackageDeleteReq) (res *v1.PackageDeleteRes, err error) {
	_, err = dao.UpgradePackage.Ctx(ctx).WherePri(req.Ids).Where(do.UpgradePackage{PackageStatus: consts.Ota_package_status_draft}).Delete()
	//todo 要不要删除文件存储中的文件？
	return
}

// UpgradeBatchCreate 简化操作，升级包下只能存在一个“进行中”的升级批次
func (*Controller) UpgradeBatchCreate(ctx context.Context, req *v1.UpgradeBatchCreateReq) (res *v1.UpgradeBatchCreateRes, err error) {
	//判断升级包是否存在且状态为release
	var pkg entity.UpgradePackage
	err = dao.UpgradePackage.Ctx(ctx).Where(do.UpgradePackage{
		Id: req.PackageId,
	}).Scan(&pkg)
	if err != nil {
		g.Log().Line().Error(ctx, "查询升级包失败:", err.Error())
		return nil, err
	}
	if pkg.Id == 0 {
		g.Log().Line().Error(ctx, "升级包不存在")
		return nil, gerror.New("升级包不存在")
	}
	if pkg.PackageStatus != int(consts.Ota_package_status_release) {
		g.Log().Line().Error(ctx, "升级包状态不合法，必须为已发布状态")
		return nil, gerror.New("升级包状态不合法，必须为已发布状态")
	}

	jsonInt64, err := json.Marshal(req.ScopeDevices)
	if err != nil {
		g.Log().Line().Error(ctx, "序列化失败:", err.Error())
		return nil, err
	}

	// 升级包下不能同时存在"进行中"或"待处理"的升级批次
	count, err := dao.UpgradeBatch.Ctx(ctx).Where(do.UpgradeBatch{
		PackageId: req.PackageId,
	}).WhereIn("batchstatus", []int32{
		consts.Ota_batch_status_executing,
		consts.Ota_batch_status_pending,
	}).Count()
	if err != nil {
		g.Log().Line().Error(ctx, "查询升级批次失败:", err.Error())
		return nil, err
	}
	if count > 0 {
		g.Log().Line().Error(ctx, "该升级包下已存在进行中或待处理的升级批次")
		return nil, gerror.New("该升级包下已存在进行中或待处理的升级批次")
	}

	//创建升级批次
	objData := do.UpgradeBatch{
		PackageId:            req.PackageId,
		ScopeType:            req.ScopeType,
		ScopeDevices:         jsonInt64,
		UserConfirm:          req.UserConfirm,
		RetryInterval:        req.RetryInterval,
		MaxRetryCount:        req.MaxRetryCount,
		OverrideExistingTask: req.OverrideTask,
		UpgradeMode:          req.UpgradeMode,
		StrategyType:         req.StrategyType,
		GrayPercentage:       req.GrayPercentage,
		GraySuccessThreshold: req.GraySuccessThreshold,
		MaxConcurrent:        consts.Upgrade_max_concurrent, //最大升级并发数
		GroupInterval:        consts.Upgrade_group_interval, //组间基础延迟（秒）
		RandomDelay:          consts.Upgrade_random_delay,   //组内随机最大延迟 （秒）
		Batchstatus:          consts.Ota_batch_status_pending,
	}

	_, err = dao.UpgradeBatch.Ctx(ctx).Data(objData).Insert()
	if err != nil {
		g.Log().Line().Error(ctx, "创建升级批次失败:", err.Error())
		return nil, err
	}

	// 定时任务去触发升级

	return &v1.UpgradeBatchCreateRes{}, nil
}

func (*Controller) UpgradeBatchList(ctx context.Context, req *v1.UpgradeBatchListReq) (res *v1.UpgradeBatchListRes, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 查询数据库
	m := dao.UpgradeBatch.Ctx(ctx).Where(do.UpgradeBatch{PackageId: req.PackageId})
	total, err := m.Count()
	if err != nil {
		return nil, err
	}

	var upgradeBatches []*entity.UpgradeBatch
	err = m.Page(int(req.Page), int(req.PageSize)).Scan(&upgradeBatches)
	if err != nil {
		return nil, err
	}

	// 转换为响应结构
	res = &v1.UpgradeBatchListRes{
		Total: int32(total),
	}
	for _, p := range upgradeBatches {
		res.List = append(res.List, &v1.UpgradeBatchInfo{
			Id:             int32(p.Id),
			Status:         int32(p.Batchstatus),
			TotalDevices:   int32(p.TotalDevices),
			SuccessDevices: int32(p.SuccessDevices),
			FailureDevices: int32(p.FailureDevices),
			PendingDevices: int32(p.PendingDevices),
			StartTime:      utility.Gtime2Str(p.StartTime),
			CompleteTime:   utility.Gtime2Str(p.CompleteTime),
			ScopeType:      int32(p.ScopeType),
			UpgradeMode:    int32(p.UpgradeMode),
			StrategyType:   int32(p.StrategyType),
			UserConfirm:    p.UserConfirm,
			RetryInterval:  int32(p.RetryInterval),
			MaxRetryCount:  int32(p.MaxRetryCount),
			OverrideTask:   p.OverrideExistingTask,
		})
	}

	return
}

func (*Controller) UpgradeRecordList(ctx context.Context, req *v1.UpgradeRecordListReq) (res *v1.UpgradeRecordIRes, err error) {
	if req.BatchId == 0 {
		g.Log().Line().Error(ctx, "BatchId字段不能为空")
		return nil, gerror.New("BatchId字段不能为空")
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 查询数据库
	m := dao.UpgradeRecord.Ctx(ctx).Where(do.UpgradeRecord{BatchId: req.BatchId})
	if req.UpgradeStatus != 0 {
		m = m.Where(do.UpgradeRecord{UpgradeStatus: req.UpgradeStatus})
	}

	if req.DeviceStatus != 0 {
		//查询设备状态
		var devices []entity.Device
		err = dao.Device.Ctx(ctx).Where(do.Device{Status: int(req.DeviceStatus)}).Fields("id").Scan(&devices)
		if err != nil {
			g.Log().Line().Error(ctx, "查询设备状态失败:", err.Error())
			return nil, err
		}
		if len(devices) == 0 {
			//没有符合条件的设备，直接返回空结果
			return &v1.UpgradeRecordIRes{
				Total: 0,
				List:  []*v1.UpgradeRecordInfo{},
			}, nil
		}
		var deviceIds []int
		for _, d := range devices {
			deviceIds = append(deviceIds, d.Id)
		}
		m = m.WhereIn("device_id", deviceIds)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}

	var UpgradeRecords []*entity.UpgradeRecord
	err = m.Page(int(req.Page), int(req.PageSize)).Scan(&UpgradeRecords)
	if err != nil {
		return nil, err
	}

	// 转换为响应结构
	res = &v1.UpgradeRecordIRes{
		Total: int32(total),
	}
	for _, p := range UpgradeRecords {
		//查询设备信息
		var device entity.Device
		err = dao.Device.Ctx(ctx).WherePri(p.DeviceId).Scan(&device)
		if err != nil {
			g.Log().Line().Error(ctx, "查询设备信息失败:", err.Error())
			return nil, err
		}
		res.List = append(res.List, &v1.UpgradeRecordInfo{
			Id:            int32(p.Id),
			DeviceId:      int32(p.DeviceId),
			DeviceKey:     device.DeviceKey,
			DeviceDesc:    device.DeviceDesc,
			DeviceStatus:  int32(device.Status),
			Version:       device.Version,
			RetryTimes:    int32(p.RetryTimes),
			StartTime:     utility.Gtime2Str(p.StartTime),
			Duration:      int32(p.Duration),
			UpgradeStatus: int32(p.UpgradeStatus),
		})
	}

	return
}

func (*Controller) UpgradeRecordOverview(ctx context.Context, req *v1.UpgradeRecordOverviewReq) (res *v1.UpgradeRecordOverviewRes, err error) {
	var records []entity.UpgradeRecord
	err = dao.UpgradeRecord.Ctx(ctx).Where(do.UpgradeRecord{BatchId: int(req.BatchId)}).Fields("upgrade_status").Scan(&records)
	if err != nil {
		g.Log().Line().Error(ctx, "查询升级记录失败:", err.Error())
		return nil, err
	}
	overview := make(map[int32]int32)
	for _, r := range records {
		overview[int32(r.UpgradeStatus)]++
	}

	res = &v1.UpgradeRecordOverviewRes{}

	for status, count := range overview {
		res.List = append(res.List, &v1.OverviewInfo{
			Status: status,
			Count:  count,
		})
	}

	res.List = append(res.List, &v1.OverviewInfo{
		Status: 0,
		Count:  int32(len(records)),
	})

	return res, nil
}

func (*Controller) UpgradeRecordLog(ctx context.Context, req *v1.UpgradeRecordLogReq) (res *v1.UpgradeRecordLogRes, err error) {
	if req.RecordId == 0 {
		g.Log().Line().Error(ctx, "RecordId字段不能为空")
		return nil, gerror.New("RecordId字段不能为空")
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 查询数据库
	m := dao.UpgradeLog.Ctx(ctx).Where(do.UpgradeLog{RecordId: req.RecordId})
	total, err := m.Count()
	if err != nil {
		return nil, err
	}

	var upgradeLogs []entity.UpgradeLog
	err = m.Page(int(req.Page), int(req.PageSize)).Scan(&upgradeLogs)
	if err != nil {
		return nil, err
	}

	// 转换为响应结构
	res = &v1.UpgradeRecordLogRes{
		Total: int32(total),
	}
	for _, p := range upgradeLogs {
		res.List = append(res.List, &v1.UpgradeRecordLogInfo{
			Time:    utility.Gtime2Str(p.Timestamp),
			Content: p.Content,
		})
	}

	return

}

func (*Controller) DeviceUpgradeCancel(ctx context.Context, req *v1.DeviceUpgradeCancelReq) (res *v1.DeviceUpgradeCancelRes, err error) {
	//查询设备信息
	var records []entity.UpgradeRecord
	err = dao.UpgradeRecord.Ctx(ctx).WhereIn("id", req.RecordIds).Scan(&records)
	if err != nil {
		g.Log().Line().Error(ctx, err.Error())
		return nil, gerror.New("查询升级记录失败")
	}

	res = &v1.DeviceUpgradeCancelRes{}

	for _, record := range records {
		var device entity.Device
		err = dao.Device.Ctx(ctx).WherePri(record.DeviceId).Scan(&device)
		if err != nil {
			g.Log().Line().Error(ctx, err.Error())
			res.DenyList = append(res.DenyList, int32(record.Id))
		}

		//没完成（5-成功，6-失败，7-已取消 之外）的才能取消升级
		if record.UpgradeStatus == int(consts.Ota_record_status_success) ||
			record.UpgradeStatus == int(consts.Ota_record_status_failed) ||
			record.UpgradeStatus == int(consts.Ota_record_status_canceled) {
			g.Log().Line().Warning(ctx, "该设备当前升级状态无法取消")
			res.DenyList = append(res.DenyList, int32(record.Id))
		}

		//下发命令给设备
		msg := model.DeviceMessage{
			MessageType: consts.Edge_ota_down,
			Content: model.DeviceMessageOtaSend{
				MessageType: consts.Ota_command_type_cancel,
				Id:          int32(record.Id),
			},
		}
		msgData, err := json.Marshal(msg)
		if err != nil {
			g.Log().Line().Error(ctx, err.Error())
			res.DenyList = append(res.DenyList, int32(record.Id))
		}
		g.Log().Line().Info(ctx, "下发命令给设备:", msgData)

		mqClient := service.MQClient()
		if mqClient == nil {
			g.Log().Line().Error(ctx, "MQClient未初始化")
			return nil, gerror.New("MQClient未初始化")
		}

		err = mqClient.Publish(consts.Kafka_topic_notify_down, []byte(device.DeviceKey), msgData)
		if err != nil {
			g.Log().Line().Error(ctx, err.Error())
			res.DenyList = append(res.DenyList, int32(record.Id))
		}

		res.AllowList = append(res.AllowList, int32(record.Id))
	}

	return
}

func (*Controller) DeviceUpgradeRetry(ctx context.Context, req *v1.DeviceUpgradeRetryReq) (res *v1.DeviceUpgradeRetryRes, err error) {
	// 查询升级记录信息
	var records []entity.UpgradeRecord
	err = dao.UpgradeRecord.Ctx(ctx).WhereIn("id", req.RecordIds).Scan(&records)
	if err != nil {
		g.Log().Line().Error(ctx, err.Error())
		return nil, err
	}

	res = &v1.DeviceUpgradeRetryRes{}

	for _, record := range records {
		// 查询设备信息以检查设备是否在线
		var device entity.Device
		err = dao.Device.Ctx(ctx).WherePri(record.DeviceId).Scan(&device)
		if err != nil {
			g.Log().Line().Error(ctx, "查询设备信息失败:", err.Error())
			res.DenyList = append(res.DenyList, int32(record.Id))
		}

		// 检查设备状态
		if device.Status != int(consts.DeviceStatusOnline) {
			errMsg := "设备不在线，无法重试升级"
			g.Log().Line().Error(ctx, errMsg)

			// 记录重试失败日志
			_, _ = dao.UpgradeLog.Ctx(ctx).Data(do.UpgradeLog{
				RecordId: record.Id,
				Content:  errMsg,
				Level:    consts.Msg_level_error,
			}).Insert()

			res.DenyList = append(res.DenyList, int32(record.Id))
		}

		// 重置升级记录状态
		_, err = dao.UpgradeRecord.Ctx(ctx).WherePri(record.Id).Data(do.UpgradeRecord{
			UpgradeStatus: consts.Ota_record_status_to_be_pushed, // 重置为待推送
			RetryTimes:    0,
			StartTime:     nil,
			Duration:      0,
			IsActive:      true, // 确保记录是活跃的
		}).Update()
		if err != nil {
			g.Log().Line().Error(ctx, "更新升级记录状态失败:", err.Error())
			return nil, gerror.Newf("更新升级记录状态失败: %v", err)
		}

		// 更新批次状态为执行中
		_, err = dao.UpgradeBatch.Ctx(ctx).WherePri(record.BatchId).Data(do.UpgradeBatch{
			Batchstatus: consts.Ota_batch_status_executing,
		}).Update()
		if err != nil {
			g.Log().Line().Error(ctx, "更新批次状态失败:", err.Error())
			return nil, gerror.Newf("更新批次状态失败: %v", err)
		}

		// 记录重试日志
		_, err = dao.UpgradeLog.Ctx(ctx).Data(do.UpgradeLog{
			RecordId: record.Id,
			Content:  "手动触发升级重试",
			Level:    consts.Msg_level_info,
		}).Insert()
		if err != nil {
			g.Log().Line().Error(ctx, "创建升级日志失败:", err.Error())
			// 不影响主流程，继续执行
		}

		res.AllowList = append(res.AllowList, int32(record.Id))
	}

	return &v1.DeviceUpgradeRetryRes{}, nil

}
