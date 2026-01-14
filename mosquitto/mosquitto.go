package mosquitto

import (
	"context"
	"pkg/errors"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MosquittoConfigEnv struct {
	Host     string `env:"MOSQUITTO_HOST"`
	User     string `env:"MOSQUITTO_USER"`
	Password string `env:"MOSQUITTO_PASSWORD"`
	ClientID string `env:"MOSQUITTO_CLIENT_ID"`
}

func NewClientMosquitto(_ context.Context, conf MosquittoConfigEnv) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions().
		AddBroker(conf.Host).
		SetClientID(conf.ClientID).
		SetUsername(conf.User).
		SetPassword(conf.Password).
		SetCleanSession(false).
		SetAutoReconnect(true).
		SetConnectRetry(true).
		SetConnectRetryInterval(5 * time.Second)

	client := mqtt.NewClient(opts)

	token := client.Connect()
	if !token.WaitTimeout(5 * time.Second) {
		return nil, errors.New("mosquitto timeout")
	}
	if token.Error() != nil {
		return nil, errors.Default.Wrap(token.Error())
	}

	return client, nil
}
