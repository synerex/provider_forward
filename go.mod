module forward-provider

go 1.14

require (
	github.com/golang/protobuf v1.4.0 // indirect
	github.com/shirou/gopsutil v2.20.3+incompatible // indirect
	github.com/stretchr/testify v1.5.1 // indirect
	github.com/synerex/synerex_api v0.3.1
	github.com/synerex/synerex_sxutil v0.4.9
	golang.org/x/net v0.0.0-20200421231249-e086a090c8fd // indirect
	golang.org/x/sys v0.0.0-20200420163511-1957bb5e6d1f // indirect
	google.golang.org/genproto v0.0.0-20200420144010-e5e8543f8aeb // indirect
	google.golang.org/grpc v1.29.0 // indirect
)

//replace github.com/synerex/synerex_sxutil => ../../sxutil
