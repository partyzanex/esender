package domain

import (
	"context"
	"log"
	"time"

	"github.com/getsentry/raven-go"
	"github.com/pkg/errors"
)

type AgentConfig struct {
	Interval time.Duration
	Pause    time.Duration
	Limit    int
	Status   EmailStatus
}

type agent struct {
	EmailStorage
	EmailSender

	interval *time.Ticker
}

func (agent *agent) Process() error {
	agentConfig := agent.AgentConfig()

	emails, err := agent.Search(context.Background(), &Filter{
		Status: agentConfig.Status,
		Limit:  agentConfig.Limit,
	})
	if err != nil {
		return err
	}

	for _, email := range emails {
		_, err := agent.Send(*email)
		if err != nil {
			agent.log(errors.Wrap(err, "sending email failed"))
			email.Status = StatusError
			errMsg := err.Error()
			email.Error = &errMsg
		} else {
			DTSent := time.Now()
			email.DTSent = &DTSent
			email.Status = StatusSent
		}

		_, err = agent.Update(context.Background(), *email)
		if err != nil {
			agent.log(errors.Wrap(err, "updating email failed"))
		}
	}

	return nil
}

func (agent *agent) Run() {
	for range agent.interval.C {
		err := agent.Process()
		if err != nil {
			agent.log(err)
		}
	}
}

func (agent *agent) Stop() {
	agent.interval.Stop()
}

func (agent) log(err error) {
	log.Println("Agent error:", err)
	raven.CaptureError(err, nil)
}

func NewAgent(storage EmailStorage, sender EmailSender) *agent {
	return &agent{
		EmailStorage: storage,
		EmailSender:  sender,
		interval:     time.NewTicker(sender.AgentConfig().Interval),
	}
}
