module github.com/infra-whizz/wzd

go 1.13

require (
	github.com/antonfisher/nested-logrus-formatter v1.1.0 // indirect
	github.com/bramvdbogaerde/go-scp v0.0.0-20200518191442-5c8efdd1d925 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/elastic/go-windows v1.0.1 // indirect
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/infra-whizz/wzcmslib v0.0.0-20200629153348-2cd07ac6fe21
	github.com/infra-whizz/wzlib v0.0.0-20200629160001-37cd698f6d29
	github.com/isbm/go-nanoconf v0.0.0-20200623180822-caf90de1965e
	github.com/nats-io/jwt v1.0.1 // indirect
	github.com/nats-io/nats.go v1.10.0
	github.com/nats-io/nkeys v0.2.0 // indirect
	github.com/orcaman/concurrent-map v0.0.0-20190826125027-8c72a8bb44f6
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/procfs v0.1.3 // indirect
	github.com/shirou/gopsutil v2.20.5+incompatible // indirect
	github.com/sirupsen/logrus v1.6.0 // indirect
	github.com/urfave/cli/v2 v2.2.0
	github.com/vmihailenco/msgpack/v4 v4.3.12 // indirect
	go.starlark.net v0.0.0-20200619143648-50ca820fafb9 // indirect
	golang.org/x/net v0.0.0-20200625001655-4c5254603344 // indirect
	golang.org/x/sys v0.0.0-20200625212154-ddb9806d33ae // indirect
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
	howett.net/plist v0.0.0-20200419221736-3b63eb3a43b5 // indirect
)

replace github.com/infra-whizz/wzlib => ../wzlib

replace github.com/isbm/go-nanoconf => ../go-nanoconf

replace github.com/infra-whizz/wzcmslib => ../wzcmslib

replace github.com/infra-whizz/wzmodlib => ../wzmodlib
