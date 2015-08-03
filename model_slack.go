package main

import ()

const (
	SlackApiUrl = "https://slack.com/api/chat.postMessage"
)

type Slack struct {
	ApiKey  string `gorethink:"api_key" json:"api_key"`
	Channel string `gorethink:"channel" json:"channel"`
}

// TODO: take the alert and make a nice Slack notification with some sort of
// colored status
func (s *Slack) SendToSlack(deploymentName, message string, status float64) {

}
