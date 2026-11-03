// Package dict 定义设备上报数据的指标字典。
//
// 指标 ID 是 GPS-TcpServer-1 与 GPS-Handler-Service-1 之间约定的协议常量：
// 发送端按 ID 构造 pb.RawMetric，接收端按 ID 解析。请统一使用本包的具名常量，
// 不要在业务代码里硬编码数字。
package dict

// DictItem 记录单个指标的定义。
type DictItem struct {
	ID       int    // 指标 ID
	Name     string // 指标名称
	Category string // 分类
	DataType string // 数据类型
	Unit     string // 单位
	Var      string // 对应的变量名
	Desc     string // 描述
}

// 基础类指标 ID。
const (
	MetricBattery         = 10001 // 电量
	MetricVoltage         = 10002 // 电池电压
	MetricSpeed           = 10003 // 速度
	MetricLongitude       = 10004 // WGS84 坐标系经度
	MetricLatitude        = 10005 // WGS84 坐标系纬度
	MetricAltitude        = 10006 // 海拔高度
	MetricDirection       = 10007 // 方向
	MetricChargeStatus    = 10008 // 充电状态
	MetricHasGps          = 10009 // GPS 定位
	MetricHasBeidou       = 10010 // 北斗定位
	MetricHasLocation     = 10011 // 定位状态
	MetricAudioMode       = 10012 // 声控模式
	MetricStepCount       = 10013 // 步数
	MetricSatelliteCount  = 10014 // 卫星数量
	MetricDeviceState     = 10015 // 设备状态
	MetricHardwareVersion = 10016 // 硬件版本
	MetricFirmwareVersion = 10017 // 固件版本
	MetricICCID           = 10018 // ICCID
)

// 音频类指标 ID。
const (
	MetricClientAppID = 10200 // 客户端唯一 ID
)

// 告警与事件类指标 ID。
const (
	MetricAlarmSos      = 20001 // 紧急告警
	MetricAlarmSpeed    = 20002 // 超速告警
	MetricAlarmLowPower = 20003 // 低电告警
	MetricAlarmShake    = 20004 // 震动告警
	MetricAlarmLight    = 20005 // 光感告警
	MetricAlarmMove     = 20006 // 移动告警
	MetricAlarmSide     = 20007 // 侧翻告警
	MetricAlarmFall     = 20008 // 跌倒告警
)

// 网络类指标 ID。
const (
	MetricWifiCount      = 30001 // WiFi 数量
	MetricCellCount      = 30002 // 基站数量
	MetricSignalStrength = 30003 // 信号强度
	MetricWifiInfos      = 30004 // WiFi 数据
	MetricCellTowers     = 30005 // 基站数据
)

// DeviceDataDict 静态字典表，以指标 ID 为键。
//
// 只读使用：请勿在运行时写入或替换该 map，多协程并发读取是安全的。
var DeviceDataDict = map[int]DictItem{
	// 基础类
	MetricBattery:         {ID: MetricBattery, Name: "电量", Category: "基础", DataType: "int32", Unit: "%", Var: "battery", Desc: "电池电量百分比"},
	MetricVoltage:         {ID: MetricVoltage, Name: "电池电压", Category: "基础", DataType: "float", Unit: "V", Var: "voltage", Desc: "电池电压"},
	MetricSpeed:           {ID: MetricSpeed, Name: "速度", Category: "基础", DataType: "float", Unit: "km/h", Var: "speed", Desc: "设备当前行驶速度"},
	MetricLongitude:       {ID: MetricLongitude, Name: "WGS84 坐标系经度", Category: "基础", DataType: "double", Unit: "", Var: "longitude", Desc: "WGS84 坐标系经度"},
	MetricLatitude:        {ID: MetricLatitude, Name: "WGS84 坐标系纬度", Category: "基础", DataType: "double", Unit: "", Var: "latitude", Desc: "WGS84 坐标系纬度"},
	MetricAltitude:        {ID: MetricAltitude, Name: "海拔高度", Category: "基础", DataType: "float", Unit: "m", Var: "altitude", Desc: "海拔高度"},
	MetricDirection:       {ID: MetricDirection, Name: "方向", Category: "基础", DataType: "int32", Unit: "", Var: "direction", Desc: "运动方向/航向角"},
	MetricChargeStatus:    {ID: MetricChargeStatus, Name: "充电状态", Category: "基础", DataType: "int32", Unit: "", Var: "charge_status", Desc: "(0:未充电, 1:充电中)"},
	MetricHasGps:          {ID: MetricHasGps, Name: "GPS定位", Category: "基础", DataType: "int32", Unit: "", Var: "has_gps", Desc: "(0:不是, 1:是)"},
	MetricHasBeidou:       {ID: MetricHasBeidou, Name: "北斗定位", Category: "基础", DataType: "int32", Unit: "", Var: "has_beidou", Desc: "(0:不是, 1:是)"},
	MetricHasLocation:     {ID: MetricHasLocation, Name: "定位状态", Category: "基础", DataType: "int32", Unit: "", Var: "has_location", Desc: "定位状态 (0:未定位, 1...)"},
	MetricAudioMode:       {ID: MetricAudioMode, Name: "声控模式", Category: "基础", DataType: "int32", Unit: "", Var: "audio_model", Desc: "声控模式"},
	MetricStepCount:       {ID: MetricStepCount, Name: "步数", Category: "基础", DataType: "int", Unit: "", Var: "step_num", Desc: "步数"},
	MetricSatelliteCount:  {ID: MetricSatelliteCount, Name: "卫星数量", Category: "基础", DataType: "int", Unit: "", Var: "gps_num", Desc: "卫星数量"},
	MetricDeviceState:     {ID: MetricDeviceState, Name: "设备状态", Category: "基础", DataType: "int32", Unit: "", Var: "device_state", Desc: "设备状态 1=心跳, 2=注册, 3=离线, 4=关机, 99=鉴权"},
	MetricHardwareVersion: {ID: MetricHardwareVersion, Name: "硬件版本", Category: "状态", DataType: "string", Unit: "", Var: "hardware_version", Desc: "硬件版本"},
	MetricFirmwareVersion: {ID: MetricFirmwareVersion, Name: "固件版本", Category: "状态", DataType: "string", Unit: "", Var: "firmware_version", Desc: "固件版本"},
	MetricICCID:           {ID: MetricICCID, Name: "ICCID", Category: "状态", DataType: "string", Unit: "", Var: "iccid", Desc: "ICCID"},

	// 音频类
	MetricClientAppID: {ID: MetricClientAppID, Name: "客户端唯一ID", Category: "音频", DataType: "string", Unit: "", Var: "clientAppId", Desc: "用户设备通知App拨打音频"},

	// 告警与事件类
	MetricAlarmSos:      {ID: MetricAlarmSos, Name: "紧急告警", Category: "告警与事件", DataType: "int32", Unit: "", Var: "alarm_sos", Desc: "SOS 紧急求助告警"},
	MetricAlarmSpeed:    {ID: MetricAlarmSpeed, Name: "超速告警", Category: "告警与事件", DataType: "int32", Unit: "", Var: "alarm_speed", Desc: "超速告警"},
	MetricAlarmLowPower: {ID: MetricAlarmLowPower, Name: "低电告警", Category: "告警与事件", DataType: "int32", Unit: "", Var: "alarm_low_power", Desc: "低电量告警"},
	MetricAlarmShake:    {ID: MetricAlarmShake, Name: "震动告警", Category: "告警与事件", DataType: "int32", Unit: "", Var: "alarm_shake", Desc: "震动告警"},
	MetricAlarmLight:    {ID: MetricAlarmLight, Name: "光感告警", Category: "告警与事件", DataType: "int32", Unit: "", Var: "alarm_light", Desc: "光感告警"},
	MetricAlarmMove:     {ID: MetricAlarmMove, Name: "移动告警", Category: "告警与事件", DataType: "int32", Unit: "", Var: "alarm_move", Desc: "移动告警"},
	MetricAlarmSide:     {ID: MetricAlarmSide, Name: "侧翻告警", Category: "告警与事件", DataType: "int32", Unit: "", Var: "alarm_side", Desc: "侧翻告警"},
	MetricAlarmFall:     {ID: MetricAlarmFall, Name: "跌倒告警", Category: "告警与事件", DataType: "int32", Unit: "", Var: "alarm_fall", Desc: "跌倒告警"},

	// 网络类
	MetricWifiCount:      {ID: MetricWifiCount, Name: "WiFi数量", Category: "网络", DataType: "int32", Unit: "", Var: "wifi_count", Desc: "WiFi数量"},
	MetricCellCount:      {ID: MetricCellCount, Name: "基站数量", Category: "网络", DataType: "int32", Unit: "", Var: "cell_count", Desc: "基站数量"},
	MetricSignalStrength: {ID: MetricSignalStrength, Name: "信号强度", Category: "网络", DataType: "int32", Unit: "", Var: "4g_rssi", Desc: "信号强度百分比"},
	MetricWifiInfos:      {ID: MetricWifiInfos, Name: "wifiInfos", Category: "网络", DataType: "Array", Unit: "", Var: "wifi_ssid", Desc: "WiFi数据 (mac,rssi)"},
	MetricCellTowers:     {ID: MetricCellTowers, Name: "cellTowers", Category: "网络", DataType: "Array", Unit: "", Var: "cell_towers", Desc: "基站数据"},
}

// Lookup 按指标 ID 查询字典项，未知 ID 返回 false。
func Lookup(id int) (DictItem, bool) {
	item, ok := DeviceDataDict[id]
	return item, ok
}
