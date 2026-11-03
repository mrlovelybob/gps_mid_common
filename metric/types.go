// Package metric 提供 pb 载荷中指标值的构造与解析工具。
//
// 发送端用 NewRawMetric 把 Go 值转成 pb.RawMetric，
// 接收端用 NewMetricsIndex 把 pb.DevicePayload.Metrics 转成按 ID 取值的索引。
package metric

// WifiInfo WiFi 信息，对应指标 30004 的数组元素。
type WifiInfo struct {
	Mac  string `json:"mac"`  // 格式：xx:xx:xx:xx:xx:xx
	RSSI int    `json:"rssi"` // 信号强度（dBm）
}

// CellTower 基站信息，对应指标 30005 的数组元素。
type CellTower struct {
	MCC    int   `json:"mcc"`    // 移动国家代码
	MNC    int   `json:"mnc"`    // 移动网络代码
	LAC    int   `json:"lac"`    // 位置区码
	CellID int64 `json:"cellId"` // 小区 ID
	RSSI   int   `json:"rssi"`   // 信号强度（dBm）
}
