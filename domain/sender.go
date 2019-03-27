package domain

type EmailSender interface {
	Send(email Email) error
	Name() string

	AgentConfig() AgentConfig
}

type Senders struct {
	list map[string]EmailSender
	def  EmailSender
}

func (senders *Senders) Add(sender EmailSender) {
	if senders.list == nil {
		senders.list = make(map[string]EmailSender)
	}

	senders.list[sender.Name()] = sender
	if senders.def == nil {
		senders.def = sender
	}
}

func (senders *Senders) SetDefault(sender EmailSender) {
	senders.Add(sender)
	senders.def = sender
}

func (senders *Senders) Get(name string) (EmailSender, bool) {
	sender, ok := senders.list[name]
	if !ok {
		sender = senders.def
	}

	return sender, ok
}

func (senders *Senders) All() []EmailSender {
	result := make([]EmailSender, len(senders.list))
	i := 0

	for _, sender := range senders.list {
		result[i] = sender
		i++
	}

	return result
}
