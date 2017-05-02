package mqtt2

import (
	"context"
	"encoding/json"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
//	"github.com/eclipse/paho.mqtt.golang"
	"testing"
)

const testConfig string = `{
  "name": "mqtt2",
  "settings": {
    "broker": "tcp://127.0.0.1:1883",
    "id": "flogoEngine",
    "user": "",
    "password": "",
    "store": "",
    "qos": "0",
    "cleansess": "false"
  },
  "endpoints": [
    {
      "flowURI": "local://testFlow",
      "settings": {
        "topic": "flogo/#"
      }
    }
  ]
}`

type TestRunner struct {
}

// Run implements action.Runner.Run
func (tr *TestRunner) Run(context context.Context, action action.Action, uri string, options interface{}) (code int, data interface{}, err error) {
	log.Debugf("Ran Action: %v", uri)
	return 0, nil, nil
}

func TestInit(t *testing.T) {
	config := trigger.Config{}
	json.Unmarshal([]byte(testConfig), config)


	f := &MQTT2Factory{}
	tgr := f.New(&config)

	runner := &TestRunner{}

	tgr.Init(runner)

}

func TestEndpoint(t *testing.T) {

	config := trigger.Config{}
	json.Unmarshal([]byte(testConfig), &config)
	// New  factory
	f := &MQTT2Factory{}
	tgr := f.New(&config)

	runner := &TestRunner{}

	tgr.Init(runner)

	tgr.Start()
	defer tgr.Stop()

/*	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	opts.SetClientID("flogoEngine")
	opts.SetUsername("")
	opts.SetPassword("")
	opts.SetCleanSession(false)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	log.Debug("---- doing publish ----")
	token := client.Publish("flogo/test_start", 0, false, "Test message payload!")
	token.Wait()

	client.Disconnect(250)
	log.Debug("Sample Publisher Disconnected") */
}