package metric

import (
	"encoding/json"

	"github.com/mrlovelybob/gps_mid_common/pb"
	"google.golang.org/protobuf/types/known/structpb"
)

// NewRawMetric 按指标 ID 构造一条 pb.RawMetric。
//
// dictID 请使用 dict 包中的具名常量（如 dict.MetricBattery）。
// 本函数不校验字典：未知 ID 同样会构造，Id 取 dictID 本身，由调用方按需用 dict.Lookup 校验。
// 共享模块不依赖任何业务项目的日志组件。
func NewRawMetric(dictID int, value any) *pb.RawMetric {
	return &pb.RawMetric{
		Id:    int32(dictID),
		Value: toValue(value),
	}
}

// toValue 把 Go 值转换为 google.protobuf.Value。
// number/string/bool 直接映射；[]WifiInfo 与 []CellTower 转为 ListValue；
// 其余类型退化为 JSON 字符串，无法序列化时为 NullValue。
func toValue(v any) *structpb.Value {
	if v == nil {
		return structpb.NewNullValue()
	}
	switch val := v.(type) {
	case string:
		return structpb.NewStringValue(val)
	case int:
		return structpb.NewNumberValue(float64(val))
	case int32:
		return structpb.NewNumberValue(float64(val))
	case int64:
		return structpb.NewNumberValue(float64(val))
	case float32:
		return structpb.NewNumberValue(float64(val))
	case float64:
		return structpb.NewNumberValue(val)
	case bool:
		return structpb.NewBoolValue(val)
	case []WifiInfo:
		items := make([]*structpb.Value, len(val))
		for i, w := range val {
			items[i] = structpb.NewStructValue(&structpb.Struct{
				Fields: map[string]*structpb.Value{
					"mac":  structpb.NewStringValue(w.Mac),
					"rssi": structpb.NewNumberValue(float64(w.RSSI)),
				},
			})
		}
		return structpb.NewListValue(&structpb.ListValue{Values: items})
	case []CellTower:
		items := make([]*structpb.Value, len(val))
		for i, c := range val {
			items[i] = structpb.NewStructValue(&structpb.Struct{
				Fields: map[string]*structpb.Value{
					"mcc":    structpb.NewNumberValue(float64(c.MCC)),
					"mnc":    structpb.NewNumberValue(float64(c.MNC)),
					"lac":    structpb.NewNumberValue(float64(c.LAC)),
					"cellId": structpb.NewNumberValue(float64(c.CellID)),
					"rssi":   structpb.NewNumberValue(float64(c.RSSI)),
				},
			})
		}
		return structpb.NewListValue(&structpb.ListValue{Values: items})
	default:
		if b, err := json.Marshal(val); err == nil {
			return structpb.NewStringValue(string(b))
		}
		return structpb.NewNullValue()
	}
}
