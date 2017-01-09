package smtp

import (
	"math/rand"
	"net"
	"time"

	"github.com/ian-kent/linkio"
	"github.com/ian-kent/service.go/log"
)

var _ ChaosMonkey = &Jim{}

// Jim is a chaos monkey
type Jim struct {
	JimConfig
}

// Accept implements ChaosMonkey.Accept
func (j *Jim) Accept(context string, conn net.Conn) bool {
	if rand.Float64() > j.AcceptChance {
		log.DebugC(context, "jim: rejecting connection", nil)
		return false
	}
	log.DebugC(context, "jim: allowing connection", nil)
	return true
}

// LinkSpeed implements ChaosMonkey.LinkSpeed
func (j *Jim) LinkSpeed(context string) *linkio.Throughput {
	rand.Seed(time.Now().Unix())
	if rand.Float64() < j.LinkSpeedAffect {
		lsDiff := j.LinkSpeedMax - j.LinkSpeedMin
		lsAffect := j.LinkSpeedMin + (lsDiff * rand.Float64())
		f := linkio.Throughput(lsAffect) * linkio.BytePerSecond
		log.DebugC(context, "jim: restricting throughput", log.Data{"throughput": f})
		return &f
	}
	log.DebugC(context, "jim: allowing unrestricted throughput", nil)
	return nil
}

// ValidRCPT implements ChaosMonkey.ValidRCPT
func (j *Jim) ValidRCPT(context string, rcpt string) bool {
	if rand.Float64() < j.RejectRecipientChance {
		log.DebugC(context, "jim: rejecting recipient", log.Data{"recipient": rcpt})
		return false
	}
	log.DebugC(context, "jim: accepting recipient", log.Data{"recipient": rcpt})
	return true
}

// ValidMAIL implements ChaosMonkey.ValidMAIL
func (j *Jim) ValidMAIL(context string, mail string) bool {
	if rand.Float64() < j.RejectSenderChance {
		log.DebugC(context, "jim: rejecting sender", log.Data{"sender": mail})
		return false
	}
	log.DebugC(context, "jim: accepting sender", log.Data{"sender": mail})
	return true
}

// ValidAUTH implements ChaosMonkey.ValidAUTH
func (j *Jim) ValidAUTH(context string, mechanism string, args ...string) bool {
	if rand.Float64() < j.RejectAuthChance {
		log.DebugC(context, "jim: rejecting authentication", log.Data{"mechanism": mechanism, "args": args})
		return false
	}
	log.DebugC(context, "jim: accepting authentication", log.Data{"mechanism": mechanism, "args": args})
	return true
}

// Disconnect implements ChaosMonkey.Disconnect
func (j *Jim) Disconnect(context string) bool {
	if rand.Float64() < j.DisconnectChance {
		log.DebugC(context, "jim: being nasty, kicking them off", nil)
		return true
	}
	log.DebugC(context, "jim: being nice, letting them stay", nil)
	return false
}
