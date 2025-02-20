package integration_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	. "github.com/switchupcb/disgo"
)

func TestConnectVoice(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)

	bot := &Client{
		Authentication: BotToken(os.Getenv("TOKEN")),
		Config:         DefaultConfig(),
		Handlers:       new(Handlers),
		VoiceHandlers:  new(VoiceHandlers),
		Sessions:       NewSessionManager(),
		ApplicationID:  os.Getenv("APPID"),
	}

	bot.Config.Gateway.EnableIntent(FlagIntentGUILD_VOICE_STATES)

	s := NewSession()

	// connect to the Discord Gateway (WebSocket Session).
	if err := s.Connect(bot); err != nil {
		t.Fatalf("%v", err)
	}

	vc := &VoiceChannelConnection{
		State: GatewayVoiceStateUpdate{
			GuildID:   os.Getenv("COVERAGE_TEST_GUILD"),
			ChannelID: Pointer(os.Getenv("COVERAGE_TEST_VOICE_CHANNEL")),
			SelfMute:  false,
			SelfDeaf:  false,
		},
		GatewaySession: s,
		VoiceSession:   nil,
		Connection:     nil,
	}

	// connect to a Discord Voice Channel.
	if err := vc.Connect(bot); err != nil {
		if sErr := s.Disconnect(); sErr != nil {
			t.Fatalf("%v", fmt.Errorf("session: %q\nvoice session: %q", sErr, err))
		}

		t.Fatalf("%v", err)
	}

	<-time.After(time.Second * 5)

	// disconnect from a Discord Voice Channel.
	if err := s.Disconnect(); err != nil {
		t.Fatalf("%v", err)
	}

	// allow Discord to close the session.
	<-time.After(time.Second * 5)
}
