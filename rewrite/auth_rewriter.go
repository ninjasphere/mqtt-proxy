package rewrite

import (
	"fmt"
	"log"

	"github.com/wolfeidau/mqtt"
)

type CredentialsRewriter interface {
	RewriteCredentials(msg *mqtt.Connect) *mqtt.Connect
}

type CredentialsReplaceRewriter struct {
	User   string
	Pass   string
	UserId string
	MqttId string
}

func NewCredentialsReplaceRewriter(user string, pass string, uid string, mqttId string) *CredentialsReplaceRewriter {
	return &CredentialsReplaceRewriter{
		User:   user,
		Pass:   pass,
		UserId: uid,
		MqttId: mqttId,
	}
}

func (crr *CredentialsReplaceRewriter) RewriteCredentials(msg *mqtt.Connect) *mqtt.Connect {

	if crr.User != "" {
		msg.UsernameFlag = true
		msg.Username = crr.User
	}

	if crr.Pass != "" {
		msg.PasswordFlag = true
		msg.Password = crr.Pass
	}

	msg.ClientId = fmt.Sprintf("%s", crr.MqttId)

	log.Printf("[creds] connecting ClientId %s", msg.ClientId)

	return msg
}
