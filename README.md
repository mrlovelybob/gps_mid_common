# gps_mid_common

GPS 设备中台共享协议模块：把（发送端）与（接收端）之间约定的 **protobuf 载荷、数据指标字典、指标值构造/解析工具** 收敛到一处，避免两个项目各自维护重复副本导致协议漂移。

- 仓库远端：`https://github.com/mrlovelybob/gps_mid_common.git`（分支 `main`）
- 当前状态：共享模块已可独立构建、测试通过；**两个原项目尚未切换到本模块**
- License：MIT

## 目录

- [1. 模块定位与数据流](#1-模块定位与数据流)
- [2. 当前目录结构](#2-当前目录结构)
- [3. 环境与依赖](#3-环境与依赖)
- [4. 编译 pb（macOS）](#4-编译-pbmacos)
- [5. 使用说明](#5-使用说明)
- [6. 验证与测试](#6-验证与测试)
- [7. 发布与引用](#7-发布与引用)
- [8. 原项目迁移对照表](#8-原项目迁移对照表)
- [9. 设计约定与兼容性](#9-设计约定与兼容性)

---

## 1. 模块定位与数据流

```
设备上报
  │
  ▼
GPS-TcpServer-1 ── metric.NewRawMetric(dict.MetricXxx, value) ─┐
                                                                ├─► pb.DevicePayload
                                                                │   └─ pb.RawMetric{Id, Value}
                                                                ▼
                                                     proto.Marshal → Kafka
                                                                │
                                                                ▼
GPS-Handler-Service-1 ← proto.Unmarshal ←──────────────────────┘
  │
  └─ metric.NewMetricsIndex(payload) → idx.Int / Float64 / String / ...
     metric.ParseWifiInfos(...) / metric.ParseCellTowers(...)
```

**抽取范围**：proto 契约 + 指标字典 + 与 pb 指标值直接相关的构造/解析工具。

**不在本模块内**（属于各项目业务层，继续留在原项目）：`Keep6Decimals`、`GpsDataExtra`、`ExtractGPSData`、`ExtractDeviceState`、各项目 `dto`、Kafka 收发、LBS 适配、数据库写入。

---

## 2. 当前目录结构

```
gps_mid_common/
├── .gitignore                      # 忽略 .DS_Store、/.idea
├── LICENSE                          # MIT，Copyright (c) 2026 mrlovelybob
├── README.md                        # 本文档
├── go.mod                           # module 声明 + 两个直接依赖
├── go.sum                           # 依赖校验和
├── proto/
│   └── device_transport.proto       # 唯一 proto 源：DevicePayload / RawMetric
├── pb/
│   └── device_transport.pb.go       # protoc-gen-go 生成码，禁止手改
├── dict/
│   ├── dict.go                      # DictItem、DeviceDataDict、32 个 Metric* 常量、Lookup
│   └── dict_test.go                 # 字典完整性 / 未知 ID 测试
├── metric/
│   ├── types.go                     # WifiInfo、CellTower（数组指标的元素类型）
│   ├── build.go                     # NewRawMetric、toValue（发送端构造）
│   ├── build_test.go                # 值类型转换、未知 ID 回归、JSON 兜底
│   ├── index.go                     # MetricsIndex、NewMetricsIndex、解析工具（接收端）
│   └── index_test.go                # 边界条件、取值方法、构造→解析往返
└── scripts/
    └── gen_pb.sh                    # macOS 一键重新生成 pb.go（可执行）
```

> `.idea/` 是本机 GoLand/IDE 配置，已被 `.gitignore` 忽略，不属于仓库内容。

### 结构描述

| 路径 | 职责 | 说明 |
| --- | --- | --- |
| `proto/` | 协议源 | 只有一份 `.proto`。`option go_package` 决定生成码的 Go 包路径，必须与 `go.mod` 的 module 前缀一致 |
| `pb/` | 生成码 | 由脚本产出，含 `DevicePayload`、`RawMetric` 及反射描述符。**不要手工编辑** |
| `dict/` | 指标字典 | 32 个指标定义 + 同名具名常量，业务代码统一引用常量，禁止硬编码 `10001` 这类数字 |
| `metric/` | 指标值工具 | 构造侧（Go 值 → `*pb.RawMetric`）与解析侧（`*pb.DevicePayload` → 按 ID 索引）对称放置 |
| `scripts/` | 构建脚本 | 只服务于「改了 proto 之后重新生成 pb.go」这一件事 |
| `go.mod` / `go.sum` | 模块定义 | Go 1.27；两个直接依赖，无间接业务依赖 |

### 包依赖关系

```
dict   ── 独立（不依赖任何其他包）
pb     ── 独立（仅依赖 google.golang.org/protobuf）
metric ── 依赖 pb（不依赖 dict，指标 ID 由调用方传入）

业务项目 ──► dict + metric + pb
```

`metric` 刻意不 import `dict`：字典增删不会破坏库的编译，调用方用 `dict.MetricXxx` 传入 ID 即可。

---

## 3. 环境与依赖

| 项 | 版本 / 值 | 用途 |
| --- | --- | --- |
| Go | `1.27`（`go.mod`） | 编译运行 |
| `google.golang.org/protobuf` | `v1.36.12` | 运行时：`pb.go`、`structpb` |
| `github.com/spf13/cast` | `v1.10.0` | `MetricsIndex` 的 `Int/Float64/String/Bool` 类型转换 |
| macOS | Apple Silicon / Intel 均可 | 编译 pb 的脚本面向 macOS |
| `protoc` | `libprotoc 35.1`（Homebrew） | **仅生成时**需要 |
| `protoc-gen-go` | 当前为 `v1.36.11` | **仅生成时**需要 |

> **版本提示**：`pb/device_transport.pb.go` 由 `protoc-gen-go v1.36.11` 生成，而 `go.mod` 的运行时库是 `v1.36.12`。两者兼容，实测 `go build ./...`、`go test ./...` 通过。若希望完全对齐，把 `protoc-gen-go` 升级到 `@v1.36.12` 后重新生成，生成码的版本注释会变化，属正常。

外部依赖只有两个，与两个原项目所用版本一致，不会引入版本冲突。

---

## 4. 编译 pb（macOS）

`pb/device_transport.pb.go` 是生成物。只有在**修改 `proto/device_transport.proto`** 或**替换 GitHub owner** 时才需要重新生成。

### 4.1 一次性环境准备（macOS / zsh）

```bash
# 1) protoc（含 google/protobuf/*.proto 等 well-known types）
brew install protobuf

# 2) Go 代码生成插件
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.12

# 3) 确保 $(go env GOPATH)/bin 在 PATH 里（protoc-gen-go 装在这里）
if ! grep -q 'go env GOPATH' ~/.zshrc; then
  echo 'export PATH="$(go env GOPATH)/bin:$PATH"' >> ~/.zshrc
fi
source ~/.zshrc

# 4) 校验
protoc --version              # libprotoc 35.1
protoc-gen-go --version       # protoc-gen-go v1.36.12
which protoc                  # /opt/homebrew/bin/protoc（Apple Silicon）
which protoc-gen-go           # /Users/<you>/go/bin/protoc-gen-go
```

Intel Mac 的 Homebrew 前缀是 `/usr/local`，脚本用 `brew --prefix` 自动适配，无需改代码。

### 4.2 一键脚本（推荐）

仓库已自带 `scripts/gen_pb.sh`，首次赋予执行权限后即可在任意目录运行：

```bash
cd /path/to/gps_mid_common
chmod +x scripts/gen_pb.sh    # 仅首次需要
./scripts/gen_pb.sh
```

脚本内容（macOS / bash 3.2 兼容，与仓库内文件一致）：

```bash
#!/usr/bin/env bash
# macOS 专用：由 proto/device_transport.proto 重新生成 pb/device_transport.pb.go。
# 用法：在任意目录执行 ./scripts/gen_pb.sh（脚本会自动切到仓库根目录）。
set -euo pipefail

cd "$(dirname "$0")/.."

if ! command -v protoc >/dev/null 2>&1; then
  echo "[gen_pb] 缺少 protoc，请先执行：brew install protobuf" >&2
  exit 1
fi
if ! command -v protoc-gen-go >/dev/null 2>&1; then
  echo "[gen_pb] 缺少 protoc-gen-go，请先执行：go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.12" >&2
  exit 1
fi

# module 路径从 go.mod 读取，替换 GitHub owner 后无需改动本脚本。
MODULE="$(awk '/^module[[:space:]]/ { print $2; exit }' go.mod)"
if [[ -z "$MODULE" ]]; then
  echo "[gen_pb] 无法从 go.mod 读取 module 路径" >&2
  exit 1
fi

# include 搜索路径：proto/ 放业务 proto；Homebrew 的 include 放 google/protobuf/*.proto。
# 数组至少含一个元素，兼容 macOS 自带的 bash 3.2（set -u 下空数组会报 unbound variable）。
INCLUDE_ARGS=(-I proto)
if command -v brew >/dev/null 2>&1; then
  protobuf_include="$(brew --prefix)/include"
  if [[ -d "$protobuf_include/google/protobuf" ]]; then
    INCLUDE_ARGS+=(-I "$protobuf_include")
  fi
fi

protoc "${INCLUDE_ARGS[@]}" \
  --go_out=. \
  --go_opt=module="$MODULE" \
  device_transport.proto

gofmt -w pb/device_transport.pb.go

echo "[gen_pb] 已生成 pb/device_transport.pb.go"
echo "[gen_pb] module        = $MODULE"
echo "[gen_pb] $(protoc --version)"
echo "[gen_pb] $(protoc-gen-go --version)"
```

脚本做了四件事：校验工具链 → 从 `go.mod` 读 module 路径 → 调 `protoc` → `gofmt` 生成码并打印版本。**不自动安装任何软件**，缺工具时只给提示。

### 4.3 等价的手工命令

不想用脚本时，在仓库根目录执行：

```bash
protoc -I proto -I "$(brew --prefix)/include" \
  --go_out=. \
  --go_opt=module=github.com/your-org/gps_mid_common \
  device_transport.proto

gofmt -w pb/device_transport.pb.go
go build ./...
```

> 手工命令里的 `--go_opt=module=` 必须与 `go.mod` 的 `module` 完全一致；脚本会自动读取，无需手填。

### 4.4 参数说明

| 参数 | 含义 |
| --- | --- |
| `-I proto` | 业务 proto 的搜索根目录，输入文件因此写作 `device_transport.proto`（不带 `proto/` 前缀） |
| `-I "$(brew --prefix)/include"` | 提供 `google/protobuf/struct.proto` 等 well-known types；Homebrew 的 protoc 通常已内置，加上更稳妥 |
| `--go_out=.` | 生成码输出到当前目录 |
| `--go_opt=module=<module>` | 从生成路径中剥离 module 前缀，最终落到 `pb/device_transport.pb.go`；这个值也会写进 descriptor 的 `go_package` |
| `device_transport.proto` | 输入文件，必须在某个 `-I` 目录下可找到 |

### 4.5 修改 proto 的注意事项

1. **不要手改 `pb/device_transport.pb.go`**，一律改 proto 后跑脚本。
2. **字段编号、`json_name`、message 名不可变更**：Kafka 上已有二进制消息依赖这些，`proto.DevicePayload` / `proto.RawMetric` 的全名也保持不变。
3. 新增字段只能追加新编号；删除字段时保留编号占位（`reserved`）。
4. `option go_package` 必须跟随 `go.mod` 的 module；改 owner 后重新生成一次。
5. 改完执行第 6 节的验证命令。

### 4.6 常见报错

| 报错 | 原因 / 处理 |
| --- | --- |
| `protoc-gen-go: program not found or is not executable` | `$(go env GOPATH)/bin` 不在 PATH，按 4.1 第 3 步处理 |
| `google/protobuf/struct.proto: File not found` | 缺 WKT include 路径，加上 `-I "$(brew --prefix)/include"` |
| `device_transport.proto: File does not reside within any path specified using --proto_path` | 输入写成了 `proto/device_transport.proto`；配合 `-I proto` 时应只写文件名 |
| 生成码出现 `github.com/your-org/...` | `go.mod` 的 module 还是占位符，见第 7 节 |

---

## 5. 使用说明

### 5.1 引入模块

```bash
# 正式引用（发布 tag 后）
go get github.com/mrlovelybob/gps_mid_common@v0.1.0

# 本地联调：在业务项目的 go.mod 里加
# replace github.com/mrlovelybob/gps_mid_common => ../gps_mid_common
```

三个包按职责单独 import：

| 包 | import | 什么时候用 |
| --- | --- | --- |
| `pb` | `github.com/<owner>/gps_mid_common/pb` | 需要 `*pb.DevicePayload` / `*pb.RawMetric` 类型时 |
| `dict` | `github.com/<owner>/gps_mid_common/dict` | 引用指标 ID 常量、查字典 |
| `metric` | `github.com/<owner>/gps_mid_common/metric` | 构造指标、按 ID 取值、解析数组指标 |

### 5.2 API 速查

**`dict`**

| 符号 | 说明 |
| --- | --- |
| `Metric*` | 32 个无类型常量（如 `MetricBattery = 10001`），可直接传给 `metric` 的函数 |
| `DictItem` | `{ID, Name, Category, DataType, Unit, Var, Desc}` |
| `DeviceDataDict` | `map[int]DictItem`，**只读**，请勿运行时写入 |
| `Lookup(id int) (DictItem, bool)` | 查字典，未知 ID 返回 `false` |

指标 ID 分段：

| 分类 | 范围 | 数量 |
| --- | --- | --- |
| 基础 / 状态 | `10001`–`10018` | 18 |
| 音频 | `10200` | 1 |
| 告警与事件 | `20001`–`20008` | 8 |
| 网络 | `30001`–`30005` | 5 |
| 合计 | | **32** |

**`metric`**

| 符号 | 说明 |
| --- | --- |
| `NewRawMetric(dictID int, value any) *pb.RawMetric` | 构造单条指标，不校验字典 |
| `NewMetricsIndex(payload *pb.DevicePayload) MetricsIndex` | 由载荷建索引，nil 安全 |
| `MetricsIndex.Value/Exists/Raw/Int/Int64/Float64/String/Bool` | 按 ID 取值，缺失返回零值 |
| `ParseWifiInfos(v *structpb.Value) []WifiInfo` | 解析指标 `30004` |
| `ParseCellTowers(v *structpb.Value) []CellTower` | 解析指标 `30005` |
| `WifiInfo` / `CellTower` | 数组指标的元素结构（字段名 `RSSI` 全大写） |

### 5.3 发送端（TcpServer 侧）

```go
package main

import (
    "github.com/your-org/gps_mid_common/dict"
    "github.com/your-org/gps_mid_common/metric"
    "github.com/your-org/gps_mid_common/pb"

    "google.golang.org/protobuf/proto"
)

func buildPayload() ([]byte, error) {
    payload := &pb.DevicePayload{
        ClientId:  "device-001",
        FactoryId: 1,
        Ts:        1730000000,
        Method:    "location",
        Metrics: []*pb.RawMetric{
            metric.NewRawMetric(dict.MetricDeviceState, 1), // 心跳
            metric.NewRawMetric(dict.MetricBattery, 88),
            metric.NewRawMetric(dict.MetricLongitude, 116.397128),
            metric.NewRawMetric(dict.MetricLatitude, 39.908722),
            metric.NewRawMetric(dict.MetricFirmwareVersion, "v1.2.3"),
            metric.NewRawMetric(dict.MetricWifiInfos, []metric.WifiInfo{
                {Mac: "aa:bb:cc:dd:ee:ff", RSSI: -55},
            }),
            metric.NewRawMetric(dict.MetricCellTowers, []metric.CellTower{
                {MCC: 460, MNC: 0, LAC: 1234, CellID: 987654321, RSSI: -80},
            }),
        },
    }
    return proto.Marshal(payload) // 写入 Kafka
}
```

### 5.4 接收端（Handler 侧）

```go
package main

import (
    "fmt"

    "github.com/your-org/gps_mid_common/dict"
    "github.com/your-org/gps_mid_common/metric"
    "github.com/your-org/gps_mid_common/pb"

    "google.golang.org/protobuf/proto"
)

func handle(raw []byte) error {
    var payload pb.DevicePayload
    if err := proto.Unmarshal(raw, &payload); err != nil {
        return err
    }

    idx := metric.NewMetricsIndex(&payload)

    battery   := idx.Int(dict.MetricBattery)
    longitude := idx.Float64(dict.MetricLongitude)
    version   := idx.String(dict.MetricFirmwareVersion)
    hasGps    := idx.Bool(dict.MetricHasGps)
    wifis     := metric.ParseWifiInfos(idx.Value(dict.MetricWifiInfos))
    towers    := metric.ParseCellTowers(idx.Value(dict.MetricCellTowers))

    item, ok := dict.Lookup(dict.MetricBattery) // 查指标元信息
    if ok {
        fmt.Println(item.Name, item.Unit)
    }
    fmt.Println(battery, longitude, version, hasGps, len(wifis), len(towers))
    return nil
}
```

> 完整可运行的「构造 → Marshal → Unmarshal → 索引 → 解析」往返示例见 `metric/index_test.go` 的 `TestBuildThenParseRoundTrip`。

### 5.5 值类型映射表

`NewRawMetric` 把 Go 值转成 `google.protobuf.Value`，接收端对应取法：

| 传给 `NewRawMetric` 的值 | 生成的 `structpb.Value` | 接收端取法 |
| --- | --- | --- |
| `nil` | `NullValue` | `idx.Raw(id) == nil` |
| `string` | `StringValue` | `idx.String(id)` |
| `int` / `int32` / `int64` / `float32` / `float64` | `NumberValue`（统一为 `float64`） | `idx.Int(id)` / `Int64` / `Float64` |
| `bool` | `BoolValue` | `idx.Bool(id)` |
| `[]metric.WifiInfo` | `ListValue` of `Struct{mac, rssi}` | `metric.ParseWifiInfos(idx.Value(id))` |
| `[]metric.CellTower` | `ListValue` of `Struct{mcc,mnc,lac,cellId,rssi}` | `metric.ParseCellTowers(idx.Value(id))` |
| 其他类型 | JSON 字符串，序列化失败则为 `NullValue` | `idx.String(id)` |

### 5.6 使用约束

- **缺失 ID 不报错**：`Int/Float64/String/Bool` 返回对应零值，`Value` 返回 `nil`，`Exists` 返回 `false`，不会 panic。
- **`NewRawMetric` 不校验字典**：传入未知 `dictID` 仍会构造，`Id` 取传入值本身；需要校验用 `dict.Lookup`。
- **`DeviceDataDict` 只读**：并发读安全，不要在运行时写入或替换该 map。
- **数值精度**：protobuf 的 `number` 是 `double`，`int64` 超过 `2^53` 会丢精度，`CellTower.CellID` 等大整数需留意。
- **字段名**：`WifiInfo.RSSI`、`CellTower.RSSI` 为全大写缩写；迁移时从项目内 `dto`（`Rssi`）转换要显式赋值。

---

## 6. 验证与测试

```bash
cd gps_mid_common

gofmt -l .          # 应无输出
go vet ./...
go build ./...
go test ./...       # 全量
go test ./... -v    # 看用例明细
```

当前共 **11 个测试函数**（含 9 个表驱动子用例）：

| 测试文件 | 覆盖点 |
| --- | --- |
| `dict/dict_test.go` | 32 个常量与字典表一一对应、`ID` 与 key 一致、`Lookup` 未知 ID |
| `metric/build_test.go` | 9 种值类型转换、未知 ID 回归（曾误写 `Id=0`）、非受支持类型 JSON 兜底 |
| `metric/index_test.go` | nil payload、nil 指标项跳过、缺失 ID 零值、各取值方法与解包、数组解析边界、**构造→索引→解析往返** |

实测结果：`go build ./...` 通过，`go test ./...` 两个包全 PASS。

---

## 7. 发布与引用

### ⚠️ 当前存在一处不一致（发布前必须处理）

| 项 | 当前值 |
| --- | --- |
| Git remote | `https://github.com/mrlovelybob/gps_mid_common.git` |
| `go.mod` module | `github.com/your-org/gps_mid_common` ← **占位符，与远端不一致** |
| 源码 import / proto `go_package` / pb.go descriptor | 同样是 `github.com/your-org/...` |

module 路径与仓库地址不一致时，别人执行 `go get github.com/mrlovelybob/gps_mid_common` 会失败。发布前把 `your-org` 全部替换为真实 owner：

```bash
# macOS / BSD sed 需要 -i ''；排除生成的 pb.go，避免破坏 descriptor 里的长度前缀
find . -type f \( -name '*.go' -o -name 'go.mod' -o -name '*.proto' -o -name 'README.md' \) \
  ! -name '*.pb.go' \
  -exec sed -i '' 's#github.com/your-org/gps_mid_common#github.com/mrlovelybob/gps_mid_common#g' {} +

# 按新 module 重新生成 pb.go（脚本自动读 go.mod，descriptor 里的 go_package 一并更新）
./scripts/gen_pb.sh

go build ./... && go test ./...
```

**不要直接对 `pb/device_transport.pb.go` 做文本替换**：`go_package` 被编码进 `rawDesc` 且带长度前缀，owner 长度变化会破坏描述符。必须重新生成。

需要同步的位置共 4 类：`go.mod` 的 `module`、`proto/device_transport.proto` 的 `option go_package`、`metric/` 下 4 个文件的 import（及本文档示例）、以及重新生成的 `pb.go`。

### 发布步骤

```bash
git add -A
git commit -m "chore: align module path with GitHub owner"
git tag v0.1.0
git push origin main
git push origin v0.1.0
```

业务项目引用：

```bash
go get github.com/mrlovelybob/gps_mid_common@v0.1.0
```

---

