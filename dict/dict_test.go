package dict

import (
	"slices"
	"testing"
)

// allMetricIDs 列出全部指标 ID 常量，用于防止常量与字典表脱节。
var allMetricIDs = []int{
	MetricBattery, MetricVoltage, MetricSpeed, MetricLongitude, MetricLatitude,
	MetricAltitude, MetricDirection, MetricChargeStatus, MetricHasGps, MetricHasBeidou,
	MetricHasLocation, MetricAudioMode, MetricStepCount, MetricSatelliteCount,
	MetricDeviceState, MetricHardwareVersion, MetricFirmwareVersion, MetricICCID,
	MetricClientAppID,
	MetricAlarmSos, MetricAlarmSpeed, MetricAlarmLowPower, MetricAlarmShake,
	MetricAlarmLight, MetricAlarmMove, MetricAlarmSide, MetricAlarmFall,
	MetricWifiCount, MetricCellCount, MetricSignalStrength, MetricWifiInfos, MetricCellTowers,
}

func TestDeviceDataDictMatchesConstants(t *testing.T) {
	if len(DeviceDataDict) != len(allMetricIDs) {
		t.Fatalf("字典项数 %d 与常量个数 %d 不一致", len(DeviceDataDict), len(allMetricIDs))
	}
	for _, id := range allMetricIDs {
		item, ok := Lookup(id)
		if !ok {
			t.Errorf("常量 %d 在字典表中缺失", id)
			continue
		}
		if item.ID != id {
			t.Errorf("字典项 %d 的 ID 字段为 %d，应与 key 一致", id, item.ID)
		}
		if item.Name == "" || item.Var == "" {
			t.Errorf("字典项 %d 的 Name/Var 为空", id)
		}
	}
	for id := range DeviceDataDict {
		if !slices.Contains(allMetricIDs, id) {
			t.Errorf("字典项 %d 没有对应的具名常量", id)
		}
	}
}

func TestLookupUnknownID(t *testing.T) {
	if _, ok := Lookup(99999); ok {
		t.Fatal("未知指标 ID 不应命中字典")
	}
}
