package metric

import (
	"github.com/spf13/cast"
	"github.com/your-org/gps_mid_common/pb"
	"google.golang.org/protobuf/types/known/structpb"
)

// MetricsIndex 以指标 ID 为键的索引，便于按 ID 快速取值。
type MetricsIndex map[int32]*structpb.Value

// NewMetricsIndex 从 pb.DevicePayload.Metrics 构建指标索引。
// payload 为 nil 或 metrics 为空时返回空索引，不会 panic。
func NewMetricsIndex(payload *pb.DevicePayload) MetricsIndex {
	metrics := payload.GetMetrics()
	idx := make(MetricsIndex, len(metrics))
	for _, m := range metrics {
		if m == nil {
			continue
		}
		idx[m.GetId()] = m.GetValue()
	}
	return idx
}

// Value 返回指定指标 ID 的原始 *structpb.Value，不存在时返回 nil。
func (idx MetricsIndex) Value(id int32) *structpb.Value {
	if idx == nil {
		return nil
	}
	return idx[id]
}

// Exists 判断指定指标 ID 是否存在。
func (idx MetricsIndex) Exists(id int32) bool {
	_, ok := idx[id]
	return ok
}

// Raw 返回指定指标 ID 解包后的原生类型（number/string/bool/nil；
// list 与 struct 返回原始 *structpb.Value，供 ParseWifiInfos / ParseCellTowers 使用）。
func (idx MetricsIndex) Raw(id int32) any {
	return unwrapValue(idx.Value(id))
}

// Int 按指标 ID 取 int 值，不存在或类型不符时返回 0。
func (idx MetricsIndex) Int(id int32) int {
	return cast.ToInt(idx.Raw(id))
}

// Int64 按指标 ID 取 int64 值，不存在或类型不符时返回 0。
func (idx MetricsIndex) Int64(id int32) int64 {
	return cast.ToInt64(idx.Raw(id))
}

// Float64 按指标 ID 取 float64 值，不存在或类型不符时返回 0。
func (idx MetricsIndex) Float64(id int32) float64 {
	return cast.ToFloat64(idx.Raw(id))
}

// String 按指标 ID 取 string 值，不存在或类型不符时返回 ""。
func (idx MetricsIndex) String(id int32) string {
	return cast.ToString(idx.Raw(id))
}

// Bool 按指标 ID 取 bool 值，不存在或类型不符时返回 false。
func (idx MetricsIndex) Bool(id int32) bool {
	return cast.ToBool(idx.Raw(id))
}

// ParseWifiInfos 把指标 30004 的 ListValue 解析为 []WifiInfo。
func ParseWifiInfos(v *structpb.Value) []WifiInfo {
	list := v.GetListValue()
	if list == nil {
		return nil
	}
	values := list.GetValues()
	result := make([]WifiInfo, 0, len(values))
	for _, item := range values {
		st := item.GetStructValue()
		if st == nil {
			continue
		}
		fields := st.GetFields()
		result = append(result, WifiInfo{
			Mac:  getStructStringField(fields, "mac"),
			RSSI: getStructIntField(fields, "rssi"),
		})
	}
	return result
}

// ParseCellTowers 把指标 30005 的 ListValue 解析为 []CellTower。
func ParseCellTowers(v *structpb.Value) []CellTower {
	list := v.GetListValue()
	if list == nil {
		return nil
	}
	values := list.GetValues()
	result := make([]CellTower, 0, len(values))
	for _, item := range values {
		st := item.GetStructValue()
		if st == nil {
			continue
		}
		fields := st.GetFields()
		result = append(result, CellTower{
			MCC:    getStructIntField(fields, "mcc"),
			MNC:    getStructIntField(fields, "mnc"),
			LAC:    getStructIntField(fields, "lac"),
			CellID: getStructInt64Field(fields, "cellId"),
			RSSI:   getStructIntField(fields, "rssi"),
		})
	}
	return result
}

// unwrapValue 把 *structpb.Value 展开为 Go 原生类型，供 cast 使用。
func unwrapValue(v *structpb.Value) any {
	if v == nil {
		return nil
	}
	switch v.GetKind().(type) {
	case *structpb.Value_NumberValue:
		return v.GetNumberValue()
	case *structpb.Value_StringValue:
		return v.GetStringValue()
	case *structpb.Value_BoolValue:
		return v.GetBoolValue()
	case *structpb.Value_NullValue:
		return nil
	case *structpb.Value_ListValue, *structpb.Value_StructValue:
		return v
	default:
		return nil
	}
}

// getStructStringField 从 structpb.Fields 中取 string 值。
func getStructStringField(fields map[string]*structpb.Value, key string) string {
	if v, ok := fields[key]; ok {
		return v.GetStringValue()
	}
	return ""
}

// getStructIntField 从 structpb.Fields 中取 int 值（protobuf number 以 float64 存储）。
func getStructIntField(fields map[string]*structpb.Value, key string) int {
	if v, ok := fields[key]; ok {
		return cast.ToInt(v.GetNumberValue())
	}
	return 0
}

// getStructInt64Field 从 structpb.Fields 中取 int64 值（protobuf number 以 float64 存储）。
func getStructInt64Field(fields map[string]*structpb.Value, key string) int64 {
	if v, ok := fields[key]; ok {
		return cast.ToInt64(v.GetNumberValue())
	}
	return 0
}
