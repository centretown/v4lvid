package ha

import "encoding/json"

type Subscription struct {
	consumer Consumer
	run      func(consumer Consumer)
	Enabled  bool
}

func NewSubcription(consumer Consumer, run func(Consumer)) *Subscription {
	sub := &Subscription{
		consumer: consumer,
		run:      run,
		Enabled:  true,
	}
	return sub
}

func (sub *Subscription) Consume(newState *Entity[json.RawMessage]) {
	sub.consumer.Copy(newState)
	sub.run(sub.consumer)
}
