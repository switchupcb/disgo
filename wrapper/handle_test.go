package wrapper

import (
	"fmt"
	"testing"
)

// TestHandle tests the Handle function for adding and removing event handlers.
func TestHandle(t *testing.T) {
	bot := &Client{
		Config:   DefaultConfig(),
		Handlers: new(Handlers),
	}

	// add event correctly.
	err := bot.Handle(FlagGatewayEventNameReady, func(r *Ready) {})
	if err != nil || len(bot.Handlers.Ready) != 1 {
		t.Fatalf("(%v): got %v, wanted %v", "AddReady", err, nil)
	}

	// test for automatic intent calculation.
	if bot.Config.Intents != 0 {
		t.Fatalf("(automatic intent calculation): got %v, wanted %v", bot.Config.Intents, 0)
	}

	// add event incorrectly.
	err = bot.Handle(FlagGatewayEventNameHello, func(r *Ready) {})
	if err == nil || len(bot.Handlers.Ready) != 1 || len(bot.Handlers.Hello) != 0 {
		t.Fatalf("(%v): got %v, wanted %v", "AddIncorrectHello", err, fmt.Errorf("event handler for %s was not added.", FlagGatewayEventNameHello))
	}

	// test for automatic intent calculation.
	if bot.Config.Intents != 0 {
		t.Fatalf("(automatic intent calculation): got %v, wanted %v", bot.Config.Intents, 0)
	}

	// add event correctly.
	err = bot.Handle(FlagGatewayEventNameChannelCreate, func(r *ChannelCreate) {})
	if err != nil || len(bot.Handlers.ChannelCreate) != 1 {
		t.Fatalf("(%v): got %v, wanted %v", "AddChannelCreate", err, nil)
	}

	// test for automatic intent calculation.
	if bot.Config.Intents != FlagIntentGUILDS {
		t.Fatalf("(automatic intent calculation): got %v, wanted %v", bot.Config.Intents, FlagIntentGUILDS)
	}

	// add similar event correctly.
	err = bot.Handle(FlagGatewayEventNameChannelUpdate, func(r *ChannelUpdate) {})
	if err != nil || len(bot.Handlers.ChannelUpdate) != 1 {
		t.Fatalf("(%v): got %v, wanted %v", "AddChannelUpdate", err, nil)
	}

	// test for automatic intent calculation.
	if bot.Config.Intents != FlagIntentGUILDS {
		t.Fatalf("(automatic intent calculation): got %v, wanted %v", bot.Config.Intents, FlagIntentGUILDS)
	}

	// remove event correctly.
	err = bot.Remove(FlagGatewayEventNameChannelUpdate, 0)
	if err != nil || len(bot.Handlers.ChannelUpdate) != 0 {
		t.Fatalf("(%v): got %v, wanted %v", "RemoveChannelUpdate", err, nil)
	}

	// remove event incorrectly (index out of bounds).
	err = bot.Remove(FlagGatewayEventNameChannelUpdate, 0)
	if err == nil || len(bot.Handlers.ChannelUpdate) != 0 {
		t.Fatalf(
			"(%v): got %v, wanted %v", "RemoveChannelUpdate",
			err,
			fmt.Errorf(errRemoveInvalidEventHandler, FlagGatewayEventNameChannelUpdate, 0),
		)
	}
}
