package metric

import (
	"reflect"
	"testing"

	"github.com/mrlovelybob/gps_mid_common/dict"
	"github.com/mrlovelybob/gps_mid_common/pb"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestNewMetricsIndexNilPayload(t *testing.T) {
	idx := NewMetricsIndex(nil)
	if len(idx) != 0 {
		t.Fatalf("nil payload 应得到空索引，实际 %d 项", len(idx))
	}
	if idx.Exists(dict.MetricBattery) || idx.Value(dict.MetricBattery) != nil {
		t.Fatal("空索引不应命中任何指标")
	}
}

func TestNewMetricsIndexSkipsNilMetric(t *testing.T) {
	idx := NewMetricsIndex(&pb.DevicePayload{
		Metrics: []*pb.RawMetric{
			nil,
			NewRawMetric(dict.MetricBattery, 66),
		},
	})
	if len(idx) != 1 || idx.Int(dict.MetricBattery) != 66 {
		t.Fatalf("nil 指标应被跳过，实际索引 %v", idx)
	}
}

func TestMetricsIndexMissingID(t *testing.T) {
	var idx MetricsIndex = nil
	if idx.Exists(dict.MetricBattery) {
		t.Fatal("nil 索引不应命中")
	}
	if idx.Int(dict.MetricBattery) != 0 ||
		idx.Int64(dict.MetricBattery) != 0 ||
		idx.Float64(dict.MetricBattery) != 0 ||
		idx.String(dict.MetricBattery) != "" ||
		idx.Bool(dict.MetricBattery) ||
		idx.Raw(dict.MetricBattery) != nil {
		t.Fatal("缺失 ID 应返回各类型零值")
	}
}

func TestMetricsIndexAccessorsAndUnwrap(t *testing.T) {
	idx := MetricsIndex{
		dict.MetricBattery:         structpb.NewNumberValue(88),
		dict.MetricLongitude:       structpb.NewNumberValue(116.397128),
		dict.MetricFirmwareVersion: structpb.NewStringValue("v1.2.3"),
		dict.MetricHasGps:          structpb.NewBoolValue(true),
		dict.MetricICCID:           structpb.NewNullValue(),
		dict.MetricWifiInfos:       structpb.NewListValue(&structpb.ListValue{}),
	}
	if !idx.Exists(dict.MetricBattery) {
		t.Fatal("已写入的指标应命中")
	}
	if idx.Int(dict.MetricBattery) != 88 {
		t.Errorf("Int = %d, want 88", idx.Int(dict.MetricBattery))
	}
	if idx.Float64(dict.MetricLongitude) != 116.397128 {
		t.Errorf("Float64 = %v, want 116.397128", idx.Float64(dict.MetricLongitude))
	}
	if idx.String(dict.MetricFirmwareVersion) != "v1.2.3" {
		t.Errorf("String = %q", idx.String(dict.MetricFirmwareVersion))
	}
	if !idx.Bool(dict.MetricHasGps) {
		t.Error("Bool 应为 true")
	}
	if idx.Raw(dict.MetricICCID) != nil {
		t.Error("NullValue 应解包为 nil")
	}
	if _, ok := idx.Raw(dict.MetricWifiInfos).(*structpb.Value); !ok {
		t.Error("ListValue 应保留原始 *structpb.Value 供 Parse 系列使用")
	}
}

func TestParseWifiInfosAndCellTowersEdgeCases(t *testing.T) {
	if got := ParseWifiInfos(nil); got != nil {
		t.Errorf("nil 输入应返回 nil，实际 %v", got)
	}
	if got := ParseCellTowers(nil); got != nil {
		t.Errorf("nil 输入应返回 nil，实际 %v", got)
	}
	if got := ParseWifiInfos(structpb.NewStringValue("not-a-list")); got != nil {
		t.Errorf("非 ListValue 应返回 nil，实际 %v", got)
	}
	// 非 struct 元素应被跳过，不影响其余元素解析
	mixed := structpb.NewListValue(&structpb.ListValue{Values: []*structpb.Value{
		structpb.NewStringValue("skip"),
		structpb.NewStructValue(&structpb.Struct{Fields: map[string]*structpb.Value{
			"mac":  structpb.NewStringValue("aa:bb:cc:dd:ee:ff"),
			"rssi": structpb.NewNumberValue(-55),
		}}),
	}})
	if got := ParseWifiInfos(mixed); len(got) != 1 || got[0].RSSI != -55 {
		t.Errorf("混合元素解析结果 = %v", got)
	}
}

// 核心链路：发送端构造 pb.RawMetric → 接收端按 ID 索引并解析数组指标。
func TestBuildThenParseRoundTrip(t *testing.T) {
	wifis := []WifiInfo{
		{Mac: "aa:bb:cc:dd:ee:ff", RSSI: -55},
		{Mac: "11:22:33:44:55:66", RSSI: -70},
	}
	towers := []CellTower{
		{MCC: 460, MNC: 0, LAC: 1234, CellID: 987654321, RSSI: -80},
	}

	payload := &pb.DevicePayload{
		ClientId: "device-001",
		Metrics: []*pb.RawMetric{
			NewRawMetric(dict.MetricBattery, 88),
			NewRawMetric(dict.MetricLongitude, 116.397128),
			NewRawMetric(dict.MetricFirmwareVersion, "v1.2.3"),
			NewRawMetric(dict.MetricHasGps, true),
			NewRawMetric(dict.MetricWifiInfos, wifis),
			NewRawMetric(dict.MetricCellTowers, towers),
		},
	}

	idx := NewMetricsIndex(payload)
	if idx.Int(dict.MetricBattery) != 88 {
		t.Errorf("电量 = %d, want 88", idx.Int(dict.MetricBattery))
	}
	if idx.Float64(dict.MetricLongitude) != 116.397128 {
		t.Errorf("经度 = %v, want 116.397128", idx.Float64(dict.MetricLongitude))
	}
	if idx.String(dict.MetricFirmwareVersion) != "v1.2.3" {
		t.Errorf("固件版本 = %q", idx.String(dict.MetricFirmwareVersion))
	}
	if !idx.Bool(dict.MetricHasGps) {
		t.Error("GPS 定位标志应为 true")
	}
	if got := ParseWifiInfos(idx.Value(dict.MetricWifiInfos)); !reflect.DeepEqual(got, wifis) {
		t.Errorf("WiFi 往返不一致: got %v, want %v", got, wifis)
	}
	if got := ParseCellTowers(idx.Value(dict.MetricCellTowers)); !reflect.DeepEqual(got, towers) {
		t.Errorf("基站往返不一致: got %v, want %v", got, towers)
	}
}
