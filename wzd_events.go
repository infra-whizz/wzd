package wzd

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/infra-whizz/wzlib"
	wzlib_logger "github.com/infra-whizz/wzlib/logger"
	wzlib_transport "github.com/infra-whizz/wzlib/transport"
	wzlib_utils "github.com/infra-whizz/wzlib/utils"
	"github.com/nats-io/nats.go"
)

type WzDaemonEvents struct {
	daemon *WzDaemon
	wzlib_logger.WzLogger
}

func NewWzDaemonEvents(daemon *WzDaemon) *WzDaemonEvents {
	wde := new(WzDaemonEvents)
	wde.daemon = daemon
	return wde
}

func (wd *WzDaemonEvents) OnPublicEvent(m *nats.Msg) {
	wd.GetLogger().Debugln("received from public", len(m.Data), "bytes")
}

func (wd *WzDaemonEvents) OnControllerEvent(m *nats.Msg) {
	wd.GetLogger().Debugln("received from controller", len(m.Data), "bytes")
	envelope := wzlib_transport.NewWzEventMsgUtils().GetMessage(m.Data)
	switch envelope.Type {
	case wzlib_transport.MSGTYPE_PING:
		wd.daemon.GetPingService().Update(envelope)
	case wzlib_transport.MSGTYPE_REGISTRATION:
		funcRetPl, funcRetEx := envelope.Payload[wzlib_transport.PAYLOAD_FUNC_RET]
		if funcRetEx && funcRetPl != nil {
			status, statEx := funcRetPl.(map[string]interface{})["status"]
			if statEx {
				statusMsg := ""
				switch status.(int64) {
				case wzlib.CLIENT_STATUS_NEW:
					statusMsg = "new"
				case wzlib.CLIENT_STATUS_ACCEPTED:
					statusMsg = "accepted"
				case wzlib.CLIENT_STATUS_REJECTED:
					statusMsg = "rejected"
				default:
					statusMsg = "unknown"
				}
				wd.GetLogger().Debugf("Client is registered as '%s'", statusMsg)
			}

			rsaPubKey, rsaEx := envelope.Payload[wzlib_transport.PAYLOAD_RSA]
			if rsaEx {
				// Save RSA pub key into pki directory
				if rsaPKb, rsaPKt := rsaPubKey.([]byte); rsaPKt {
					rsaPubFp := wd.daemon.GetCryptoUtils().PEMKeyFingerprintFromBytes(rsaPKb)
					if rsaPubFp != wd.daemon.GetClusterFingerprint() {
						wd.GetLogger().Errorln("Discarded message due to RSA public key fingerprint do not match. This is a possible attack, BTW. :-)")
						wd.GetLogger().Debugln("Incoming fingerprint:", rsaPubFp)
						wd.GetLogger().Debugln("Expected fingerprint:", wd.daemon.GetClusterFingerprint())
					} else {
						// save PEM somewhere
						wd.GetLogger().Debug("Saving public PEM key")
						errcode, err := wd.daemon.SaveClusterPublicPEMKey(rsaPKb)
						if errcode != wzlib_utils.EX_TEMPFAIL && err != nil {
							wd.GetLogger().Errorf("Error saving cluster's public PEM: %s", err.Error())
						} else if err != nil {
							wd.GetLogger().Warningf("Skipping saving public PEM: %s", err.Error())
						} else {
							wd.GetLogger().Info("Public PEM file has been saved")
						}
					}
				} else {
					wd.GetLogger().Errorln("Discarded message due to a wrong type of RSA key payload")
				}
			} else {
				wd.GetLogger().Errorln("Discarded message due to no RSA public key has been received")
			}

		} else {
			wd.GetLogger().Debugln("Discarding response message: no function return slot has been defined")
		}
	default:
		wd.GetLogger().Debugln("Discaring message: no idea what type it is:", envelope.Type)
		spew.Dump(envelope)
	}
}
