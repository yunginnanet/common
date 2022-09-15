module git.tcp.direct/kayos/common

go 1.19

require (
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.27.0
	golang.org/x/crypto v0.0.0-20220817201139-bc19a97f63c8
	inet.af/netaddr v0.0.0-20220811202034-502d2d690317
	nullprogram.com/x/rng v1.1.0
)

require (
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	go4.org/intern v0.0.0-20211027215823-ae77deb06f29 // indirect
	go4.org/unsafe/assume-no-moving-gc v0.0.0-20220617031537-928513b29760 // indirect
	golang.org/x/sys v0.0.0-20210927094055-39ccf1dd6fa6 // indirect
)

retract (
	v0.0.0-20220210125455-40e3d2190a52
)
