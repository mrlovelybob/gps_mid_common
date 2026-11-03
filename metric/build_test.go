package metric

import (
	"testing"

	"github.com/your-org/gps_mid_common/dict"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestNewRawMetricValueKinds(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  *structpb.Value
	}{
		{"nil", nil, structpb.NewNullValue()},
		{"string", "3G", structpb.NewStringValue("3G")},
		{"int", 88, structpb.NewNumberValue(88)},
		{"int32", int32(12), structpb.NewNumberValue(12)},
		{"int64", int64(1700000000), structpb.NewNumberValue(1700000000)},
		{"float32", float32(3.75), structpb.NewNumberValue(3.75)},
		{"float64", 116.397128, structpb.NewNumberValue(116.397128)},
		{"bool", true, structpb.NewBoolValue(true)},
		{"fallback_json", struct{ A int }{A: 1}, structpb.NewStringValue(`{"A":1}`)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRawMetric(dict.MetricBattery, tt.value)
			if got.GetId() != dict.MetricBattery {
				t.Errorf("Id = %d, want %d", got.GetId(), dict.MetricBattery)
			}
			if !proto.Equal(got.GetValue(), tt.want) {
				t.Errorf("Value = %v, want %v", got.GetValue(), tt.want)
			}
		})
	}
}

// 回归：原实现在字典查不到 ID 时取零值 item.ID，会把 Id 写成 0。
// 共享模块应保留调用方传入的 dictID。
func TestNewRawMetricUnknownDictIDKeepsID(t *testing.T) {
	const unknownID = 99999
	got := NewRawMetric(unknownID, 1)
	if got.GetId() != unknownID {
		t.Fatalf("未知 ID 时 Id = %d, want %d", got.GetId(), unknownID)
	}
	if _, ok := dict.Lookup(unknownID); ok {
		t.Fatal("测试用例使用了已存在的字典 ID，请换一个")
	}
}

func TestNewRawMetricUnsupportedTypeFallsBackToJSON(t *testing.T) {
	got := NewRawMetric(dict.MetricWifiInfos, map[string]int{"a": 1})
	if _, ok := got.GetValue().GetKind().(*structpb.Value_StringValue); !ok {
		t.Fatalf("未知类型应退化为 JSON 字符串，实际为 %T", got.GetValue().GetKind())
	}
	if got.GetValue().GetStringValue() != `{"a":1}` {
		t.Errorf("JSON 兜底结果 = %q", got.GetValue().GetStringValue())
	}
}
