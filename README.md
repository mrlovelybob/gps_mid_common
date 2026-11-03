## 使用示例

### 发送端（TcpServer 侧）

```go
import (
    "github.com/your-org/gps_mid_common/dict"
    "github.com/your-org/gps_mid_common/metric"
    "github.com/your-org/gps_mid_common/pb"
)

payload := &pb.DevicePayload{
    ClientId: "device-001",
    Metrics: []*pb.RawMetric{
        metric.NewRawMetric(dict.MetricDeviceState, 1),   // 心跳
        metric.NewRawMetric(dict.MetricBattery, 88),
        metric.NewRawMetric(dict.MetricLongitude, 116.397128),
        metric.NewRawMetric(dict.MetricWifiInfos, []metric.WifiInfo{
            {Mac: "aa:bb:cc:dd:ee:ff", RSSI: -55},
        }),
    },
}
// proto.Marshal(payload) 后写入 Kafka
```

### 接收端（Handler 侧）

```go
idx := metric.NewMetricsIndex(payload)

battery  := idx.Int(dict.MetricBattery)
longitude := idx.Float64(dict.MetricLongitude)
version   := idx.String(dict.MetricFirmwareVersion)
wifis     := metric.ParseWifiInfos(idx.Value(dict.MetricWifiInfos))
towers    := metric.ParseCellTowers(idx.Value(dict.MetricCellTowers))

if _, ok := dict.Lookup(someID); !ok {
    // 未知指标 ID
}
```

`idx` 下标用 `int32`，具名常量为无类型常量，可直接传入；缺失 ID 返回对应类型零值，不会 panic。

---

## 验证

```bash
go mod tidy
go build ./...
go test ./...
```
