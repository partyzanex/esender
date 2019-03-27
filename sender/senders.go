package sender

import (
	"github.com/partyzanex/esender/domain"
	"github.com/partyzanex/esender/sender/smtp"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

type Config map[string]map[string]interface{}

func Create(config []interface{}) (*domain.Senders, error) {
	senders := &domain.Senders{}

	for _, senderMap := range config {
		senderMap := senderMap.(map[interface{}]interface{})
		name, ok := senderMap["name"].(string)
		if !ok {
			return nil, errors.New("sender name not found in config")
		}

		interval, ok := senderMap["interval"]
		if !ok {
			return nil, errors.New("agent interval is undefined")
		}
		pause, ok := senderMap["pause"]
		if !ok {
			return nil, errors.New("agent pause is undefined")
		}
		limit, ok := senderMap["limit"]
		if !ok {
			return nil, errors.New("agent limit is undefined")
		}
		status, ok := senderMap["status"].(string)
		if !ok {
			return nil, errors.New("agent status is undefined")
		}
		agentCfg := domain.AgentConfig{
			Interval: cast.ToDuration(interval),
			Pause:    cast.ToDuration(pause),
			Limit:    cast.ToInt(limit),
			Status:   domain.EmailStatus(status),
		}

		switch name {
		case "smtp":
			host, ok := senderMap["host"].(string)
			if !ok {
				return nil, errors.New("smtp host is undefined")
			}
			port, ok := senderMap["port"]
			if !ok {
				return nil, errors.New("smtp port is undefined")
			}
			tls, ok := senderMap["tls"]
			if !ok {
				return nil, errors.New("smtp tls is undefined")
			}
			username, ok := senderMap["username"].(string)
			if !ok {
				return nil, errors.New("smtp username is undefined")
			}
			password, ok := senderMap["password"].(string)
			if !ok {
				return nil, errors.New("smtp password is undefined")
			}
			isDefault, ok := senderMap["default"]

			sender := smtp.New(smtp.Config{
				AgentConfig: agentCfg,
				Host:        host,
				Port:        uint16(cast.ToInt(port)),
				TLS:         cast.ToBool(tls),
				UserName:    username,
				Password:    password,
			})

			senders.Add(sender)

			if ok && cast.ToBool(isDefault) {
				senders.SetDefault(sender)
			}
		default:
			return nil, errors.New("Unknown sender " + name)
		}
	}

	return senders, nil
}
