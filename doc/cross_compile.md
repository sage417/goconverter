### How to

execute in powershell

```
$ENV:GOOS="linux"
$ENV:GOARCH="arm"
$ENV:GOARM="arm64"
```


Architecture	Status	GOARM value	GOARCH value
- ARMv4 and below	not supported	n/a	n/a
- ARMv5	supported	GOARM=5	GOARCH=arm
- ARMv6	supported	GOARM=6	GOARCH=arm
- ARMv7	supported	GOARM=7	GOARCH=arm
- ARMv8	supported	n/a	GOARCH=arm64
- ARMv9	supported	n/a	GOARCH=arm64

https://go.dev/wiki/GoArm

https://www.ohyee.cc/post/note_compile_go_to_openwrt

### Example
```
$ENV:GOOS="linux" ; $ENV:GOARCH="arm64"; go build -ldflags="-s -w" -o clash_meta

$ENV:GOOS="linux" ; $ENV:GOARCH="arm64"; go build -ldflags="-s -w" -trimpath  -o clash_meta
 
$ENV:GOOS="linux" ; $ENV:GOARCH="arm64"; $ENV:CGO_ENABLED="0"; go build -ldflags="-s -w" --trimpath -o clash_meta; ~\Downloads\upx-5.1.0-win64\upx.exe clash_meta

$ENV:GOOS="linux"; $ENV:GOARCH="arm64"; $ENV:CGO_ENABLED="0"; go build -ldflags="-s -w" -gcflags="all=-B" --trimpath -o clash_meta; ~\Downloads\upx-5.1.0-win64\upx.exe -5 clash_meta

$ENV:GOOS="windows" ; $ENV:GOARCH="amd64";go build -ldflags="-s -w" --trimpath -o verge-mihomo.exe; ~\Downloads\upx-5.1.0-win64\upx.exe verge-mihomo.exe;

$ENV:GOOS="linux" ; $ENV:GOARCH="mipsle"; $env:GOMIPS="softfloat"; $ENV:CGO_ENABLED="0"; go build -x -v -ldflags="-s -w" --trimpath -o clash_meta; ~\Downloads\upx-5.1.0-win64\upx.exe clash_meta
```


### rust compile
```
export CC_aarch64_unknown_linux_musl=clang
export AR_aarch64_unknown_linux_musl=llvm-ar
export CARGO_TARGET_AARCH64_UNKNOWN_LINUX_MUSL_RUSTFLAGS="-Clink-self-contained=yes -Clinker=rust-lld"
cargo build --release -target = aarch64-unknown-linux-musl

CGO_ENABLED=0
```

### go pprof

```
go tool pprof http://192.168.31.1:9090/debug/pprof/allocs
go tool pprof http://192.168.31.1:9090/debug/pprof/heap
go tool pprof http://192.168.31.1:9090/debug/pprof/profiles?seconds=120
```