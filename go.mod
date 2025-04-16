module git.tcp.direct/kayos/common

go 1.21
toolchain go1.24.1

require (
	golang.org/x/crypto v0.35.0
	nullprogram.com/x/rng v1.1.0
)

require golang.org/x/sys v0.30.0 // indirect

retract (
	v0.9.8 // nil error push
	v0.9.1 // premature (race condition)
	v0.9.0 // premature
	v0.0.0-20220210125455-40e3d2190a52
)
