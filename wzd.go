package wzd

import (
	"log"
	"time"

	wzd_transport "github.com/infra-whizz/wzd/transport"
	"github.com/nats-io/nats.go"
)

const (
	CHANNEL_PUBLIC     = "public"
	CHANNEL_CONTROLLER = "controller"
)

type WzDaemonStatus struct {
	Running bool
}

type WzDaemon struct {
	status                 *WzDaemonStatus
	transport              *wzd_transport.WzdPubSub
	publicSubscription     *nats.Subscription
	controllerSubscription *nats.Subscription
}

// Constructor
func NewWzDaemon() *WzDaemon {
	wd := new(WzDaemon)
	wd.status = &WzDaemonStatus{}
	wd.transport = wzd_transport.NewWizPubSub()

	return wd
}

func (wd *WzDaemon) onPublicEvent(m *nats.Msg) {
	log.Println("received from public", len(m.Data), "bytes")
}

func (wd *WzDaemon) onControllerEvent(m *nats.Msg) {
	log.Println("received from controller", len(m.Data), "bytes")
}

// GetTransport return transport object
func (wd *WzDaemon) GetTransport() *wzd_transport.WzdPubSub {
	return wd.transport
}

func (wd *WzDaemon) IsRunning() bool {
	return wd.status.Running
}

// Run the daemon
func (wd *WzDaemon) Run() *WzDaemon {
	if wd.IsRunning() {
		return wd
	}

	var err error
	wd.GetTransport().Start()
	wd.publicSubscription, err = wd.GetTransport().GetSubscriber().Subscribe(CHANNEL_PUBLIC, wd.onPublicEvent)
	if err != nil {
		log.Panicln("Unable to subscribe to a public channel:", err.Error())
	}
	wd.controllerSubscription, err = wd.GetTransport().GetSubscriber().Subscribe(CHANNEL_CONTROLLER, wd.onControllerEvent)
	if err != nil {
		log.Panicln("Unable to subscribe to a controller channel:", err.Error())
	}
	wd.status.Running = true

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
		panic("Drain error: " + err.Error())
	}
	wd.status.Running = false
}
