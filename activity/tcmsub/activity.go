package tcmsub

import (
        "github.com/TIBCOSoftware/flogo-lib/core/activity"
        "github.com/TIBCOSoftware/flogo-lib/logger"
        "github.com/jvanderl/tib-eftl"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
)

// log is the default package logger
var log = logger.GetLogger("activity-dabbate-tcmsub")

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new AppActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MyActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval - Receives a message from TIBCO eFTL/TCM
// TODO
func (a *MyActivity) Eval(context activity.Context) (done bool, err error) {

	// Get the activity data from the context
	wsURL, _ := context.GetInput("url").(string)
	wsAuthKey, _ := context.GetInput("authkey").(string)
	wsClientID, _ := context.GetInput("clientid").(string)
	wsDestinationName, _ := context.GetInput("destinationname").(string)
	wsDestinationMatch, _ := context.GetInput("destinationmatch").(string)
	wsdurable, _ := context.GetInput("durable").(bool)

	// wsDestinationValue, _ := context.GetInput("destinationvalue").(string)
	//wsMessageName, _ := context.GetInput("messagename").(string)
	//wsMessageValue, _ := context.GetInput("messagevalue").(string)
	//	wsCert, _ := context.GetInput("certificate").(string)
	wsCert := ""

	var tlsConfig *tls.Config

	if wsCert != "" {
		// TLS configuration uses CA certificate from a PEM file to
		// authenticate the server certificate when using wss:// for
		// a secure connection
		caCert, err := base64.StdEncoding.DecodeString(wsCert)
		if err != nil {
			log.Errorf("unable to decode certificate: %s", err)
			return false, err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig = &tls.Config{
			RootCAs: caCertPool,
		}
	} else {
		// TLS configuration accepts all server certificates
		// when using wss:// for a secure connection
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	// channel for receiving connection errors
	errChan := make(chan error, 1)

	// set connection options
	opts := &eftl.Options{
		ClientID:  wsClientID,
		Username:  "",
		Password:  wsAuthKey,
		TLSConfig: tlsConfig,
	}

	// connect to the server
	conn, err := eftl.Connect(wsURL, opts, errChan)
	if err != nil {
		context.SetOutput("result", "ERR_CONNECT_HOST")
		log.Errorf("Error connecing to TCM server: [%s]", err)
		return false, err
	}

	// close the connection when done
	defer conn.Disconnect()

	// channel for receiving subscription response
	subChan := make(chan *eftl.Subscription, 1)

	// channel for receiving published messages
	msgChan := make(chan eftl.Message, 1000)

	// destName := handler.GetSetting("destinationname") // replaced by wsDestinationName
	// destMatch := handler.GetSetting("destinationmatch") // replaced by wsDestinationMatch
	// create the message content matcher
	matcher := ""
	if wsDestinationMatch == "*" {
		matcher = fmt.Sprintf("{\"%s\":true}", wsDestinationName)
	} else {
		matcher = fmt.Sprintf("{\"%s\":\"%s\"}", wsDestinationName, wsDestinationMatch)
	}
	//durable = wsdurable // was ... handler.GetSetting("durable"))
	if wsdurable {
		durablename, _ := context.GetInput("durablename").(string)
		log.Infof("Subscribing to destination: %s:%s, durable name:%s", wsDestinationName, wsDestinationMatch, durablename)
		conn.SubscribeAsync(matcher, durablename, msgChan, subChan)
	} else {
		log.Infof("Subscribing to destination: %s:%s", wsDestinationName, wsDestinationMatch)
		conn.SubscribeAsync(matcher, "", msgChan, subChan)
	}

	for {
		select {
		case sub := <-subChan:
			if sub.Error != nil {
				log.Infof("subscribe operation failed: %s", sub.Error)
				return false, sub.Error
			}
			log.Infof("subscribed with matcher %s", sub.Matcher)
		case msg := <-msgChan:
			log.Infof("received message: %s", msg)
			// see if we can find a matching handler

			//destName := handler.GetSetting("destinationname")
			//destMatch := handler.GetSetting("destinationmatch")
			//msgName := handler.GetSetting("messagename")
			msgName := context.GetInput("messagename").(string)
			if (msg[wsDestinationName].(string) == wsDestinationMatch) || (msg[wsDestinationName] != nil && wsDestinationMatch == "*") {
				destination := wsDestinationName + "_" + wsDestinationMatch
				message := msg[msgName].(string)
				log.Debugf("Received message", msgName)
				//actionId, found := t.destinationToActionId[destination]
				// if found {
				//log.Debugf("About to run action for Id [%s]", actionId)
				//t.RunAction(actionId, message, destination)
				//} else {
				//	log.Debug("actionId not found")
				//}
			}

		case err := <-errChan:
			log.Infof("connection error: %s", err)
			return false, err
		}
	}

	//TODO

	// channel for receiving publish completions
	//	compChan := make(chan *eftl.Completion, 1000)
	//	// publish the message
	//	if wsDestinationValue != "" {
	//		conn.PublishAsync(eftl.Message{
	//			wsDestinationName: wsDestinationValue,
	//			wsMessageName:     wsMessageValue,
	//		}, compChan)
	//	} else {
	//		conn.PublishAsync(eftl.Message{
	//			wsMessageName: wsMessageValue,
	//		}, compChan)
	//	}
	//
	//	for {
	//		select {
	//		case comp := <-compChan:
	//			if comp.Error != nil {
	//				log.Errorf("Error while sending message to wsHost: [%s]", comp.Error)
	//				context.SetOutput("result", "ERR_SEND_MESSAGE")
	//				return false, comp.Error
	//			}
	//			log.Debugf("published message: %s", comp.Message)
	//			context.SetOutput("result", "OK")
	//			return true, nil
	//		case err := <-errChan:
	//			log.Errorf("connection error: %s", err)
	//			context.SetOutput("result", "ERR_CONNECT_HOST")
	//			return false, err
	//		}
	//	}

}
