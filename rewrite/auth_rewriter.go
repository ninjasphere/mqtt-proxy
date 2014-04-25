package rewrite

import (
	"fmt"

	"github.com/wolfeidau/mqtt"
)

type CredentialsRewriter interface {
	RewriteCredentials(msg *mqtt.Connect) *mqtt.Connect
}

type CredentialsReplaceRewriter struct {
	User   string
	Pass   string
	UserId uint
}

func NewCredentialsReplaceRewriter(user string, pass string, uid uint) *CredentialsReplaceRewriter {
	return &CredentialsReplaceRewriter{
		User:   user,
		Pass:   pass,
		UserId: uid,
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

	msg.ClientId = fmt.Sprintf("%d-%s", crr.UserId, msg.ClientId)

	//	log.Printf("[creds] %v", msg)

	return msg
}
