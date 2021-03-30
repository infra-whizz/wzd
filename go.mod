module github.com/infra-whizz/wzd

go 1.13

require (
	github.com/StackExchange/wmi v0.0.0-20210224194228-fe8f1750fd46 // indirect
	github.com/antonfisher/nested-logrus-formatter v1.3.0 // indirect
	github.com/bramvdbogaerde/go-scp v0.0.0-20201229172121-7a6c0268fa67 // indirect
	github.com/davecgh/go-spew v1.1.1
	github.com/elastic/go-sysinfo v1.6.0 // indirect
	github.com/elastic/go-windows v1.0.1 // indirect
	github.com/go-ole/go-ole v1.2.5 // indirect
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/infra-whizz/wzcmslib v0.0.0-20210225171912-247b9a4ae031
	github.com/infra-whizz/wzlib v0.0.0-20210302184443-9c300336f5e1
	github.com/isbm/go-nanoconf v0.0.0-20200623180822-caf90de1965e
	github.com/karrick/godirwalk v1.16.1 // indirect
	github.com/magefile/mage v1.11.0 // indirect
	github.com/nats-io/jwt v1.2.2 // indirect
	github.com/nats-io/nats.go v1.10.0
	github.com/orcaman/concurrent-map v0.0.0-20210106121528-16402b402231
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/procfs v0.6.0 // indirect
	github.com/shirou/gopsutil v3.21.2+incompatible // indirect
	github.com/sirupsen/logrus v1.8.0
	github.com/tklauser/go-sysconf v0.3.4 // indirect
	github.com/urfave/cli/v2 v2.3.0
	github.com/vmihailenco/msgpack/v4 v4.3.12 // indirect
	github.com/vmihailenco/tagparser v0.1.2 // indirect
	go.starlark.net v0.0.0-20210305151048-6a590ae7f4eb // indirect
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83 // indirect
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110 // indirect
	golang.org/x/sys v0.0.0-20210305230114-8fe3ee5dd75b // indirect
	google.golang.org/appengine v1.6.7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	howett.net/plist v0.0.0-20201203080718-1454fab16a06 // indirect
)

replace github.com/infra-whizz/wzlib => ../wzlib

replace github.com/isbm/go-nanoconf => ../go-nanoconf

replace github.com/infra-whizz/wzcmslib => ../wzcmslib

replace github.com/infra-whizz/wzmodlib => ../wzmodlib
