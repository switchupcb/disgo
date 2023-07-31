package voice_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	. "github.com/switchupcb/disgo/wrapper"
)

func TestConnectVoice(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.TraceLevel)

	bot := &Client{
		Authentication: BotToken(os.Getenv("TOKEN")),
		Config:         DefaultConfig(),
		Handlers:       new(Handlers),
		Sessions:       NewSessionManager(),
	}

	bot.Config.Gateway.EnableIntent(FlagIntentGUILD_VOICE_STATES)

	s := NewSession()

	// connect to the Discord Gateway (WebSocket Session).
	if err := s.Connect(bot); err != nil {
		t.Fatalf("%v", err)
	}

	vc := &VoiceConnection{
		State: GatewayVoiceStateUpdate{
			GuildID:   os.Getenv("COVERAGE_TEST_GUILD"),
			ChannelID: Pointer(os.Getenv("COVERAGE_TEST_VOICE_CHANNEL")),
			SelfMute:  false,
			SelfDeaf:  false,
		},
		Session:      s,
		VoiceSession: nil,
		Connection:   nil,
		Handlers:     nil,
	}

	// connect to a Discord Voice Channel.
	if err := vc.Connect(bot); err != nil {
		if sErr := s.Disconnect(); err != nil {
			t.Fatalf("%v", fmt.Errorf("session: %q\nvoice session: %q", sErr, err))
		}

		t.Fatalf("%v", err)
	}

	time.Sleep(time.Second * 20)

	if err := s.Disconnect(); err != nil {
		t.Fatalf("%v", err)
	}

	time.Sleep(time.Second * 10)
}
