package shard_test

import (
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/switchupcb/disgo"
	. "github.com/switchupcb/disgo/shard"
)

// TestReconnect tests Connect(), Disconnect(), and Reconnect() of the Shard Manager.
func TestReconnect(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	bot := &disgo.Client{
		Authentication: disgo.BotToken(os.Getenv("TOKEN")),
		Config:         disgo.DefaultConfig(),
		Handlers:       &disgo.Handlers{},
		Sessions:       disgo.NewSessionManager(),
	}

	bot.Config.Gateway.ShardManager = new(InstanceShardManager)
	bot.Config.Gateway.ShardManager.SetNumShards(2)

	s := bot.Config.Gateway.ShardManager

	// connect to the Discord Gateway (WebSocket Session).
	if err := s.Connect(bot); err != nil {
		t.Fatalf("%v", err)
	}

	if len(s.GetSessions()) != 2 {
		t.Fatalf("expected 2 sessions but got %d", len(s.GetSessions()))
	}

	time.Sleep(time.Second)

	// reconnect to the Discord Gateway (WebSocket Session).
	if err := s.Reconnect(bot); err != nil {
		t.Fatalf("%v", err)
	}

	if len(s.GetSessions()) != 2 {
		t.Fatalf("expected 2 sessions but got %d", len(s.GetSessions()))
	}

	time.Sleep(time.Second)

	// disconnect from the Discord Gateway (WebSocket Connection).
	if err := s.Disconnect(); err != nil {
		t.Fatalf("%v", err)
	}

	if len(s.GetSessions()) != 0 {
		t.Fatalf("expected 0 sessions but got %d", len(s.GetSessions()))
	}
}
