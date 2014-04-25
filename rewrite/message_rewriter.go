package rewrite

import "github.com/wolfeidau/mqtt"

type MsgRewriter struct {
	CredentialsRewriter CredentialsRewriter
	IngressRewriter     TopicRewriter
	EgressRewriter      TopicRewriter
}

func CreatMsgRewriter(credentialsRewriter CredentialsRewriter, ingressRewriter TopicRewriter, egressRewriter TopicRewriter) *MsgRewriter {
	return &MsgRewriter{
		CredentialsRewriter: credentialsRewriter,
		IngressRewriter:     ingressRewriter,
		EgressRewriter:      egressRewriter,
	}
}

func (mr *MsgRewriter) RewriteIngress(msg mqtt.Message) mqtt.Message {

	//	log.Printf("[ingress] msg: %s %v", reflect.TypeOf(msg), msg)

	switch msg := msg.(type) {
	case *mqtt.Connect:
		msg = mr.CredentialsRewriter.RewriteCredentials(msg)
	case *mqtt.Publish:
		msg.TopicName = mr.IngressRewriter.RewriteTopicName(msg.TopicName)
	case *mqtt.Subscribe:
		msg.Topics = mr.IngressRewriter.RewriteTopics(msg.Topics)
	case *mqtt.Unsubscribe:
		msg.Topics = mr.IngressRewriter.RenameTopicNames(msg.Topics)
	}
	return msg
}

func (mr *MsgRewriter) RewriteEgress(msg mqtt.Message) mqtt.Message {

	//	log.Printf("[egress] msg: %s %v", reflect.TypeOf(msg), msg)

	switch msg := msg.(type) {
	case *mqtt.Publish:
		msg.TopicName = mr.EgressRewriter.RewriteTopicName(msg.TopicName)
	}
	return msg
}
