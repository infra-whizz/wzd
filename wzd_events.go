package wzd

import (
	"log"

	"github.com/davecgh/go-spew/spew"
	wzlib_transport "github.com/infra-whizz/wzlib/transport"
	"github.com/nats-io/nats.go"
)

type WzDaemonEvents struct {
	daemon *WzDaemon
}

func NewWzDaemonEvents(daemon *WzDaemon) *WzDaemonEvents {
	wde := new(WzDaemonEvents)
	wde.daemon = daemon
	return wde
}

func (wd *WzDaemonEvents) OnPublicEvent(m *nats.Msg) {
	log.Println("received from public", len(m.Data), "bytes")
}

func (wd *WzDaemonEvents) OnControllerEvent(m *nats.Msg) {
	log.Println("received from controller", len(m.Data), "bytes")
	envelope := wzlib_transport.NewWzEventMsgUtils().GetMessage(m.Data)
	switch envelope.Type {
	case wzlib_transport.MSGTYPE_PING:
		wd.daemon.GetPingService().Update(envelope)
	default:
		log.Println("Discaring message: no idea what type it is:", envelope.Type)
		spew.Dump(envelope)
	}
}
