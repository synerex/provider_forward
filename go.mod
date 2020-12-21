module forward-provider

go 1.14

require (
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/shirou/gopsutil v3.20.11+incompatible // indirect
	github.com/synerex/synerex_api v0.4.2
	github.com/synerex/synerex_nodeapi v0.5.4 // indirect
	github.com/synerex/synerex_proto v0.1.10 // indirect
	github.com/synerex/synerex_sxutil v0.6.2
	golang.org/x/net v0.0.0-20201216054612-986b41b23924 // indirect
	golang.org/x/sys v0.0.0-20201218084310-7d0127a74742 // indirect
	golang.org/x/text v0.3.4 // indirect
	google.golang.org/genproto v0.0.0-20201214200347-8c77b98c765d // indirect
	google.golang.org/grpc v1.34.0 // indirect
)

//replace github.com/synerex/synerex_sxutil => ../../sxutil
