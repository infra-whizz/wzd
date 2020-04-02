package wzd_events

import (
	"fmt"
	"time"

	"github.com/infra-whizz/wzlib"
	wzlib_logger "github.com/infra-whizz/wzlib/logger"
	wzlib_transport "github.com/infra-whizz/wzlib/transport"
	cmap "github.com/orcaman/concurrent-map"
)

type WzPingStat struct {
	Ticks     int64
	Responded bool
}

// WzPingEvent is a class for pinger
type WzPingEvent struct {
	uid       string
	pings     cmap.ConcurrentMap
	transport *wzlib_transport.WzdPubSub
}

// NewWzPingEvent creates a new ping event
func NewWzPingEvent(transport *wzlib_transport.WzdPubSub, uid string) *WzPingEvent {
	pe := new(WzPingEvent)
	pe.uid = uid
	pe.transport = transport
	pe.pings = cmap.New()

	return pe
}

// Update ping event container by message handler
func (pe *WzPingEvent) Update(msg *wzlib_transport.WzGenericMessage) {
	uid := msg.Payload[wzlib_transport.PAYLOAD_SYSTEM_ID].(string)
	if pe.uid != uid {
		return
	}

	pingId, ok := msg.Payload[wzlib_transport.PAYLOAD_PING_ID]
	if !ok {
		log.Println("Ping message contains no 'ping.id' section!")
	} else {
		pingStatItf, ok := pe.pings.Get(pingId.(string))
		pingStat := pingStatItf.(*WzPingStat)
		if !ok {
			log.Println("Unable to find ping ID for", pingId)
		} else {
			pingStat.Ticks = time.Now().Unix() - pingStat.Ticks
			pingStat.Responded = true
		}
	}
}

// Ping any channel, as long as subscriber supports this feature.
func (pe *WzPingEvent) Ping(descr string, seconds int) (int64, error) {
	return pe.waitForResponse(pe.ping(wzlib.CHANNEL_CLIENT), descr, seconds)
}

// Pings the channel on the MQ network
func (pe *WzPingEvent) ping(channel string) string {
	pingId := wzlib.MakeJid()
	pe.pings.Set(pingId, &WzPingStat{
		Ticks:     time.Now().Unix(),
		Responded: false})

	envelope := wzlib_transport.NewWzMessage(wzlib_transport.MSGTYPE_PING)
	envelope.Payload[wzlib_transport.PAYLOAD_SYSTEM_ID] = pe.uid
	envelope.Payload[wzlib_transport.PAYLOAD_PING_ID] = pingId

	msg, _ := envelope.Serialise()
	if err := pe.transport.GetPublisher().Publish(channel, msg); err != nil {
		log.Println("Unable to ping controller:", err.Error())
	}

	return pingId
}

// waitForResponse from the ping. If seconds are less then 1, then wait 1 second.
func (pe *WzPingEvent) waitForResponse(pingId string, descr string, seconds int) (int64, error) {
	var latency int64 = -1
	var err error = nil
	if seconds < 1 {
		seconds = 1
	}

	cycles := 1000 * seconds
	cycle := 0
	for {
		pingStatItf, ok := pe.pings.Get(pingId)
		pingStat := pingStatItf.(*WzPingStat)
		if !ok {
			err = fmt.Errorf("Unknown pingId: %s", pingId)
			break
		}
		if pingStat.Responded {
			latency = pingStat.Ticks
			pe.pings.Remove(pingId)
			log.Println("Ping latency:", latency)
			break
		}
		time.Sleep(time.Millisecond)
		cycle++

		if cycle > cycles {
			err = fmt.Errorf("Ping timeout")
			break
		}
	}

	return latency, err
}
