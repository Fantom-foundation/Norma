module github.com/Fantom-foundation/Norma

go 1.21

toolchain go1.21.6

require (
	github.com/Fantom-foundation/go-opera v1.1.3-rc.5
	github.com/docker/go-connections v0.4.0
	github.com/ethereum/go-ethereum v1.11.6
	github.com/jupp0r/go-priority-queue v0.0.0-20160601094913-ab1073853bde
	github.com/tyler-smith/go-bip32 v1.0.0
	github.com/tyler-smith/go-bip39 v1.1.0
	github.com/urfave/cli/v2 v2.25.7
	go.uber.org/mock v0.4.0
	golang.org/x/exp v0.0.0-20231110203233-9a3e6036ecaa
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/docker/docker v24.0.5+incompatible
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/holiman/uint256 v1.2.4 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/shirou/gopsutil v3.21.4-0.20210419000835-c7a38de76ee5+incompatible // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
)

require (
	github.com/FactomProject/basen v0.0.0-20150613233007-fe3947df716e // indirect
	github.com/FactomProject/btcutilecc v0.0.0-20130527213604-d3a63a5752ec // indirect
	github.com/Fantom-foundation/Substate v0.0.0-20230224090651-4c8c024214f4 // indirect
	github.com/Fantom-foundation/lachesis-base v0.0.0-20240116072301-a75735c4ef00 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/VictoriaMetrics/fastcache v1.12.2 // indirect
	github.com/btcsuite/btcd v0.20.1-beta // indirect
	github.com/cakturk/go-netstat v0.0.0-20200220111822-e5b49efee7a5 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/deckarep/golang-set v1.8.0 // indirect
	github.com/docker/distribution v2.8.2+incompatible // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v0.0.5-0.20220116011046-fa5810519dcb // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/hashicorp/golang-lru v0.5.5-0.20210104140557-80c98217689d // indirect
	github.com/holiman/bloomfilter/v2 v2.0.3 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/mattn/go-sqlite3 v1.14.17 // indirect
	github.com/moby/term v0.0.0-20221205130635-1aeaba878587 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/tsdb v0.10.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/rjeczalik/notify v0.9.2 // indirect
	github.com/stretchr/testify v1.9.0 // indirect
	github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7 // indirect
	golang.org/x/crypto v0.23.0 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
	gotest.tools/v3 v3.4.0 // indirect
)

replace github.com/Fantom-foundation/go-opera v1.1.3-rc.5 => github.com/Fantom-foundation/Sonic v1.2.1-d

replace github.com/ethereum/go-ethereum => github.com/Fantom-foundation/go-ethereum-substate v1.1.1-0.20240523185630-b058188a043a

replace github.com/golang/mock v1.3.1 => go.uber.org/mock v0.4.0
