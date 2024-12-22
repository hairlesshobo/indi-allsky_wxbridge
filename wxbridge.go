package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"indi-allsky_wxbridge/model"
)

var (
	wxClient     mqtt.Client
	allskyClient mqtt.Client
	config       *model.Config
)

func main() {
	if len(os.Args) < 2 {
		slog.Error("the config file must be provided as the first and only argument!")
		os.Exit(1)
	}

	configFile := os.Args[1]

	readConfig(configFile)

	wxOptions := brokerOptions(config.WxHost, config.WxPort, config.WxClientID, config.WxUsername, config.WxPassword)
	wxClient = mqtt.NewClient(wxOptions)
	if token := wxClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	allskyOptions := brokerOptions(config.AllskyHost, config.AllskyPort, config.AllskyClientID, config.AllskyUsername, config.AllskyPassword)
	allskyClient = mqtt.NewClient(allskyOptions)
	if token := allskyClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	subscribe(wxClient, config.WxTopic)

	c := make(chan os.Signal, 5)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	slog.Info("Waiting for messages...")
	for {
		if len(c) > 0 {
			// received sigterm handle it
			connectionLostHandler(nil, nil)
			break
		}

		time.Sleep(100 * time.Millisecond)
	}
}

func publishAllsky(name string, value float64) {
	token := allskyClient.Publish(config.AllskyTopic+"/"+name, 0, false, fmt.Sprintf("%1.3f", value))
	if token.Wait() && token.Error() != nil {
		slog.Error(fmt.Sprintf("Failed to send %s message to allsky broker: %s", name, token.Error().Error()))
	}
}

func publishToAllsky(payload *model.WeewxLoopPayload) {
	if payload.OutTempF != "" {
		tempF, err := strconv.ParseFloat(payload.OutTempF, 64)
		if err != nil {
			slog.Error("failed to parse weewx temperature: " + err.Error())
		} else {
			tempC := (tempF - 32) / 1.8
			publishAllsky("temperature", tempC)
		}
	}

	if payload.DewpointF != "" {
		dewpointF, err := strconv.ParseFloat(payload.DewpointF, 64)
		if err != nil {
			slog.Error("failed to parse weewx dewpoint: " + err.Error())
		} else {
			dewpointC := (dewpointF - 32) / 1.8
			publishAllsky("dewpoint", dewpointC)
		}
	}

	if payload.OutHumidity != "" {
		outHumidity, err := strconv.ParseFloat(payload.OutHumidity, 64)
		if err != nil {
			slog.Error("failed to parse weewx humidity: " + err.Error())
		} else {
			publishAllsky("humidity", outHumidity)
		}
	}

	if payload.BarometerInHg != "" {
		barometerInHg, err := strconv.ParseFloat(payload.BarometerInHg, 64)
		if err != nil {
			slog.Error("failed to parse weewx barometer: " + err.Error())
		} else {
			pressure := barometerInHg * 33.86389
			publishAllsky("pressure", pressure)
		}
	}

	if payload.WindDir != "" {
		windDir, err := strconv.ParseFloat(payload.WindDir, 64)
		if err != nil {
			slog.Error("failed to parse weewx winddir: " + err.Error())
		} else {
			publishAllsky("winddir", windDir)
		}
	}
}

func subscribe(client mqtt.Client, topic string) {
	token := wxClient.Subscribe(topic, 1, nil)
	token.Wait()
	slog.Info(fmt.Sprintf("subscribed to %s topic: %s\n", getBrokerName(client), topic))
}

func brokerOptions(host string, port int, clientId string, username string, password string) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", host, port))
	opts.SetClientID(clientId)

	if username != "" {
		opts.SetUsername(username)
		opts.SetPassword(password)
	}

	opts.SetDefaultPublishHandler(messagePubHandler)

	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectionLostHandler

	return opts
}

func messagePubHandler(client mqtt.Client, msg mqtt.Message) {
	slog.Info(fmt.Sprintf("received %s message: %s ", getBrokerName(client), msg.Payload()))

	if client == wxClient {
		var payload model.WeewxLoopPayload
		err := json.Unmarshal(msg.Payload(), &payload)
		if err != nil {
			slog.Error("error unmarshaling wx payload: " + err.Error())
			os.Exit(1)
		}

		publishToAllsky(&payload)
	}
}

func connectHandler(client mqtt.Client) {
	slog.Info(fmt.Sprintf("%s broker connected", getBrokerName(client)))
}

func connectionLostHandler(client mqtt.Client, err error) {
	if err != nil {
		slog.Error("%s broker connection lost: %v", getBrokerName(client), err)
	} else {
		slog.Info("shutting down")
	}

	if wxClient != nil {
		wxClient.Disconnect(200)
	}

	if allskyClient != nil {
		allskyClient.Disconnect(200)
	}

	if err != nil {
		os.Exit(1)
	}
}

func readConfig(configFile string) {
	content, err := os.ReadFile(configFile)
	if err != nil {
		slog.Error("error reading config file: " + err.Error())
		os.Exit(1)
	}

	// Now let's unmarshall the data into `config`
	err = json.Unmarshal(content, &config)
	if err != nil {
		slog.Error("error unmarshaling config: " + err.Error())
		os.Exit(1)
	}
}

func getBrokerName(client mqtt.Client) string {
	if client == allskyClient {
		return "allsky"
	}

	if client == wxClient {
		return "wx"
	}

	return "UNKNOWN"
}
