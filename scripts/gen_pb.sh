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
