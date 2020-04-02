package wzd

import (
	"github.com/infra-whizz/wzlib"
	wzlib_logger "github.com/infra-whizz/wzlib/logger"
	wzlib_transport "github.com/infra-whizz/wzlib/transport"
)

type WzDaemonBoot struct {
	pkiDir string
	daemon *WzDaemon
	wzlib_logger.WzLogger
}

func NewWzDaemonBoot(daemon *WzDaemon) *WzDaemonBoot {
	wdb := new(WzDaemonBoot)
	wdb.daemon = daemon
	return wdb
}

func (wd *WzDaemonBoot) waitForController() {
	ping := wd.daemon.GetPingService()
	seconds := 5
	for {
		latency, err := ping.Ping("controller node", seconds)
		if err != nil {
			wd.GetLogger().Warningln("Controller timeout. Trying again for", seconds, "seconds.")
		} else {
			wd.GetLogger().Warningln("Latency:", latency)
			break
		}
	}
}

// onClientBoot is when the client is booting up
func (wdb *WzDaemonBoot) onClientBoot() {
	if wdb.daemon.GetTraits() == nil {
		wdb.GetLogger().Panicln("Traits were not initalised")
	}
	// Check if this client has AES key
	rsa, aes := wdb.pkiVerify()
	wdb.GetLogger().Infoln("PKI check. RSA:", rsa)
	wdb.GetLogger().Infoln("PKI check. AES:", aes)

	wdb.waitForController()

	if !(rsa && aes) {
		// Generate key pair, if none
		if !rsa {
			if err := wdb.daemon.GetRSA().GenerateKeyPair(wdb.pkiDir); err != nil {
				wdb.GetLogger().Fatalln(err.Error()) // Game over, no PKI possible.
			}
		}

		// Request for AES token, using RSA PKI
		wdb.sendRegistrationRequest()
	} else {
		wdb.sendTraitsMap()
	}
}

// Verify if RSA and AES keys are present and loaded
func (wd *WzDaemonBoot) pkiVerify() (rsa bool, aes bool) {
	rsa, aes = false, false

	// Check RSA
	rsa = wd.daemon.GetRSA().LoadPEMKeyPair(wd.pkiDir) == nil

	// Check AES
	if !wd.daemon.GetAES().IsLoaded() {
		if err := wd.daemon.GetAES().LoadKey(wd.pkiDir); err != nil {
			wd.GetLogger().Errorln(err.Error())
		} else {
			aes = true
		}
	}
	return
}

/*
	Sends registration request message.
	The registration request is sent in three cases:

	1. Completely new client (PKI not yet completed)
	2. AES key was rotated
*/
func (wd *WzDaemonBoot) sendRegistrationRequest() {
	pem, err := wd.daemon.GetRSA().GetPublicPEMKey(wd.pkiDir)
	if err != nil {
		wd.GetLogger().Errorln(err.Error())
		return
	}
	envelope := wzlib_transport.NewWzMessage(wzlib_transport.MSGTYPE_REGISTRATION)
	envelope.Payload[wzlib_transport.PAYLOAD_RSA] = pem
	envelope.Payload[wzlib_transport.PAYLOAD_SYSTEM_ID] = wd.daemon.GetTraits().GetContainer().Get("uid")
	envelope.Payload[wzlib_transport.PAYLOAD_SYSTEM_FQDN] = wd.daemon.GetTraits().GetContainer().Get("hostname")

	msg, err := envelope.Serialise()
	if err != nil {
		wd.GetLogger().Errorln(err.Error())
		return
	}

	err = wd.daemon.GetTransport().GetPublisher().Publish(wzlib.CHANNEL_CLIENT, msg)
	if err != nil {
		wd.GetLogger().Errorln(err.Error())
		return
	}
	wd.GetLogger().Infoln("Sent registration request to", wzlib.CHANNEL_CLIENT)
}

// Sends traits map to the controller
func (wd *WzDaemonBoot) sendTraitsMap() {
}
