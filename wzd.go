package wzd

import (
	"os"
	"time"

	wzd_events "github.com/infra-whizz/wzd/events"
	wzlib "github.com/infra-whizz/wzlib"
	wzlib_crypto "github.com/infra-whizz/wzlib/crypto"
	wzlib_logger "github.com/infra-whizz/wzlib/logger"
	wzlib_sockets "github.com/infra-whizz/wzlib/sockets"
	wzlib_traits "github.com/infra-whizz/wzlib/traits"
	wzlib_traits_attributes "github.com/infra-whizz/wzlib/traits/attributes"
	wzlib_transport "github.com/infra-whizz/wzlib/transport"
	"github.com/nats-io/nats.go"
)

type WzDaemonStatus struct {
	Running bool
}

type WzDaemon struct {
	events                 *WzDaemonEvents
	boot                   *WzDaemonBoot
	status                 *WzDaemonStatus
	unixSock               *wzlib_sockets.WzLocalSocketCommunicator
	transport              *wzlib_transport.WzdPubSub
	traits                 *wzlib_traits.WzTraits
	publicSubscription     *nats.Subscription
	controllerSubscription *nats.Subscription
	aes                    *wzlib_crypto.WzAES
	rsa                    *wzlib_crypto.WzRSA

	// Services
	pingService *wzd_events.WzPingEvent

	wzlib_logger.WzLogger
}

// Constructor
func NewWzDaemon() *WzDaemon {
	wd := new(WzDaemon)
	wd.status = &WzDaemonStatus{}
	wd.unixSock = wzlib_sockets.NewWzLocalSocketCommunicator("/tmp/wzd.sock")
	wd.boot = NewWzDaemonBoot(wd)
	wd.events = NewWzDaemonEvents(wd)
	wd.transport = wzlib_transport.NewWizPubSub()
	wd.aes = wzlib_crypto.NewWzAES()
	wd.rsa = wzlib_crypto.NewWzRSA()

	return wd
}

// Sets PKI directory
func (wd *WzDaemon) SetPkiDirectory(dirname string) *WzDaemon {
	wd.boot.pkiDir = dirname
	return wd
}

// SetTraitsFile initialises traits instance in the boot sub-object
func (wd *WzDaemon) SetTraitsFile(fpath string) *WzDaemon {
	wd.traits = wzlib_traits.NewTraits(fpath)
	wd.traits.LoadAttribute(wzlib_traits_attributes.NewSysInfo())
	wd.traits.LoadAttribute(wzlib_traits_attributes.NewCPUInfo())
	wd.traits.LoadAttribute(wzlib_traits_attributes.NewDiskInfo())
	wd.traits.LoadAttribute(wzlib_traits_attributes.NewNetInfo())
	wd.traits.Save()

	return wd
}

func (wd *WzDaemon) GetPingService() *wzd_events.WzPingEvent {
	if wd.pingService == nil {
		uid := wd.GetTraits().GetContainer().Get("uid").(string)
		wd.pingService = wzd_events.NewWzPingEvent(wd.transport, uid)
	}
	return wd.pingService
}

// GetAES returns AES utility API
func (wd *WzDaemon) GetAES() *wzlib_crypto.WzAES {
	return wd.aes
}

// GetRSA returns RSA utility API
func (wd *WzDaemon) GetRSA() *wzlib_crypto.WzRSA {
	return wd.rsa
}

// GetTransport return transport object
func (wd *WzDaemon) GetTransport() *wzlib_transport.WzdPubSub {
	return wd.transport
}

// GetTraits returns initialised traits
func (wd *WzDaemon) GetTraits() *wzlib_traits.WzTraits {
	return wd.traits
}

func (wd *WzDaemon) IsRunning() bool {
	return wd.status.Running
}

// Run the daemon
func (wd *WzDaemon) Run() *WzDaemon {
	if wd.IsRunning() {
		return wd
	}

	err := wd.unixSock.Bind()
	if err != nil {
		panic(err)
	}
	if wd.unixSock.IsClient() {
		wd.GetLogger().Errorln("Another instance is running already!")
		os.Exit(1)
	}

	wd.GetTransport().Start()
	wd.publicSubscription, err = wd.GetTransport().GetSubscriber().Subscribe(wzlib.CHANNEL_PUBLIC, wd.events.OnPublicEvent)
	if err != nil {
		wd.GetLogger().Panicln("Unable to subscribe to a public channel:", err.Error())
	}
	wd.controllerSubscription, err = wd.GetTransport().GetSubscriber().Subscribe(wzlib.CHANNEL_CONTROLLER, wd.events.OnControllerEvent)
	if err != nil {
		wd.GetLogger().Panicln("Unable to subscribe to a controller channel:", err.Error())
	}
	wd.status.Running = true

	wd.boot.onClientBoot()

	return wd
}

func (wd *WzDaemon) AppLoop() {
	for {
		time.Sleep(10 * time.Second)
	}
}

// Stop wzd
func (wd *WzDaemon) Stop() {
	if err := wd.GetTransport().GetSubscriber().Drain(); err != nil {
		wd.GetLogger().Panicln("Drain error: " + err.Error())
	}
	wd.status.Running = false
}
