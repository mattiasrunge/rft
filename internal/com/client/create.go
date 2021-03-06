package client

import (
	"errors"
	"time"

	"barf/internal/com/channel"
	"barf/internal/com/protocol"
	"barf/internal/op"
)

// CreateOperation sends a create request
func CreateOperation(opType op.OperationType, args op.OperationArgs) (*op.Operation, error) {
	message := protocol.NewRequestCreateMessage(opType, args)
	err := channel.Broadcast(message)

	if err != nil {
		return nil, err
	}

	c1 := make(chan *protocol.ResponseCreate)

	onResponse := func(response *protocol.ResponseCreate) {
		c1 <- response
	}

	bus.SubscribeOnce(string(message.RequestCreate.ID), onResponse)

	select {
	case res := <-c1:
		if res.Result == protocol.ResponseOk {
			return res.Operation, nil
		}

		return nil, errors.New(string(res.Message))
	case <-time.After(10 * time.Second):
		bus.Unsubscribe(string(message.RequestCreate.ID), onResponse)
		return nil, errors.New("timeout")
	}
}
