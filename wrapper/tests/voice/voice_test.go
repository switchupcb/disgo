package voice_test

import (
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

	s := NewSession()

	voiceChannel := GatewayVoiceStateUpdate{
		GuildID:   os.Getenv("COVERAGE_TEST_GUILD"),
		ChannelID: Pointer(os.Getenv("COVERAGE_TEST_VOICE_CHANNEL")),
		SelfMute:  false,
		SelfDeaf:  false,
	}

	// connect to the Discord Gateway (WebSocket Session).
	if err := s.ConnectVoice(bot, voiceChannel); err != nil {
		t.Fatalf("%v", err)
	}

	time.Sleep(time.Second * 10)

	if err := s.Disconnect(); err != nil {
		t.Fatalf("%v", err)
	}

	time.Sleep(time.Second * 10)
}
