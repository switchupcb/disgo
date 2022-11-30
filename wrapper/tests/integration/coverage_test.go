package integration_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/rs/zerolog"
	. "github.com/switchupcb/disgo/wrapper"
)

// TestCoverage tests 100+ endpoints (requests) and respective events from the Discord API.
func TestCoverage(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	bot := &Client{
		Authentication: BotToken(os.Getenv("COVERAGE_TEST_TOKEN")),
		Config:         DefaultConfig(),
		Handlers:       new(Handlers),
		Sessions:       []*Session{NewSession()},
	}

	bot.Config.Request.Timeout = time.Second * time.Duration(3)

	// set the bot's Application ID.
	requestGetCurrentBotApplicationInformation := &GetCurrentBotApplicationInformation{}
	app, err := requestGetCurrentBotApplicationInformation.Send(bot)
	if err != nil {
		t.Fatal(fmt.Errorf("GetCurrentBotApplicationInformation: %w", err))
	}

	bot.ApplicationID = app.ID
	if app.ID == "" {
		t.Fatal("GetCurrentBotApplicationInformation: expected non-empty Application ID")
	}

	initializeEventHandlers(bot)

	// Connect the session to the Discord Gateway (WebSocket Connection).
	if err := bot.Sessions[0].Connect(bot); err != nil {
		t.Fatalf("can't open websocket session to Discord: %v", err)
	}

	eg, ctx := errgroup.WithContext(context.Background())

	// Call endpoints with no dependencies.
	eg.Go(func() error {
		request := &ListVoiceRegions{}
		regions, err := request.Send(bot)
		if err != nil {
			return fmt.Errorf("ListVoiceRegions: %w", err)
		}

		if len(regions) == 0 {
			return fmt.Errorf("ListVoiceRegions: expected non-empty Voice Regions Array")
		}

		return nil
	})

	// Call endpoints with one or more dependencies.
	eg.Go(func() error {
		return testCommands(bot)
	})

	eg.Go(func() error {
		return testGuild(bot)
	})

	eg.Go(func() error {
		return testChannel(bot)
	})

	eg.Go(func() error {
		return testMessage(bot)
	})

	// wait until all required requests have been processed.
	select {
	case <-ctx.Done():
		t.Fatalf("%v", eg.Wait())
	default:
	}

	if err := eg.Wait(); err != nil {
		t.Fatalf("%v", err)
	}

	// Disconnect the session from the Discord Gateway (WebSocket Connection).
	if err := bot.Sessions[0].Disconnect(); err != nil {
		t.Fatalf("%v", err)
	}

	// allow Discord to close the session.
	<-time.After(time.Second * 5)
}

// initializeEventHandlers initializes the event handlers necessary for this test.
func initializeEventHandlers(bot *Client) {

}

// testCommands tests all endpoints that are dependent on a global command.
func testCommands(bot *Client) error {
	createGlobalApplicationCommand := CreateGlobalApplicationCommand{
		Name:        "main",
		Description: "A basic command",
	}

	command, err := createGlobalApplicationCommand.Send(bot)
	if err != nil {
		return fmt.Errorf("CreateGlobalApplicationCommand: %w", err)
	}

	if command.ID == "" {
		return fmt.Errorf("CreateGlobalApplicationCommand: expected non-null ApplicationCommand object")
	}

	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		commands, err := new(GetGlobalApplicationCommands).Send(bot)
		if err != nil {
			return fmt.Errorf("GetGlobalApplicationCommands: %w", err)
		}

		if len(commands) == 0 {
			return fmt.Errorf("GetGlobalApplicationCommands: expected non-empty Global Application Command List")
		}

		return nil
	})

	eg.Go(func() error {
		getGlobalApplicationCommand := &GetGlobalApplicationCommand{CommandID: command.ID}
		got, err := getGlobalApplicationCommand.Send(bot)
		if err != nil {
			return fmt.Errorf("GetGlobalApplicationCommand: %w", err)
		}

		if got == nil {
			return fmt.Errorf("GetGlobalApplicationCommand: expected non-null ApplicationCommand object")
		}

		return nil
	})

	eg.Go(func() error {
		editGlobalApplicationCommand := &EditGlobalApplicationCommand{
			CommandID:   command.ID,
			Name:        "notmain",
			Description: "This is not a main global command.",
		}

		editedCommand, err := editGlobalApplicationCommand.Send(bot)
		if err != nil {
			return fmt.Errorf("EditGlobalApplicationCommand: %w", err)
		}

		if editedCommand.ID == "" {
			return fmt.Errorf("EditGlobalApplicationCommand: expected non-null ApplicationCommand object")
		}

		return nil
	})

	// wait until all requests have been processed.
	select {
	case <-ctx.Done():
		return fmt.Errorf("%w", eg.Wait())
	default:
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("%w", eg.Wait())
	}

	deleteGlobalApplicationCommand := &DeleteGlobalApplicationCommand{
		CommandID: command.ID,
	}

	if err := deleteGlobalApplicationCommand.Send(bot); err != nil {
		return fmt.Errorf("DeleteGlobalApplicationCommand: %w", err)
	}

	return nil
}

// testGuild tests all endpoints dependent on a guild.
func testGuild(bot *Client) error {
	getGuild := &GetGuild{GuildID: os.Getenv("COVERAGE_TEST_GUILD")}

	guild, err := getGuild.Send(bot)
	if err != nil {
		return fmt.Errorf("GetGuild: %w", err)
	}

	if guild.ID == "" {
		return fmt.Errorf("GetGuild: expected non-empty Guild ID.")
	}

	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		return testGuildMember(bot, guild)
	})

	eg.Go(func() error {
		return testGuildScheduledEvent(bot, guild)
	})

	eg.Go(func() error {
		/*
			ListAutoModerationRulesForGuild
			CreateAutoModerationRule
		*/

		return nil
	})

	eg.Go(func() error {
		getGuildPreview := &GetGuildPreview{GuildID: guild.ID}
		if _, err := getGuildPreview.Send(bot); err != nil {
			return fmt.Errorf("GetGuildPreview: %w", err)
		}

		return nil
	})

	eg.Go(func() error {
		getGuildChannels := &GetGuildChannels{GuildID: guild.ID}
		if _, err := getGuildChannels.Send(bot); err != nil {
			return fmt.Errorf("GetGuildChannels: %w", err)
		}

		return nil
	})

	eg.Go(func() error {
		listActiveGuildThreads := &ListActiveGuildThreads{GuildID: guild.ID}
		if _, err := listActiveGuildThreads.Send(bot); err != nil {
			return fmt.Errorf("ListActiveGuildThreads: %w", err)
		}

		return nil
	})

	eg.Go(func() error {
		searchGuildMembers := &SearchGuildMembers{
			GuildID: guild.ID,
			Query:   Pointer("D"),
		}

		if _, err := searchGuildMembers.Send(bot); err != nil {
			return fmt.Errorf("SearchGuildMembers: %w", err)
		}

		return nil
	})

	eg.Go(func() error {
		getGuildVoiceRegions := &GetGuildVoiceRegions{GuildID: guild.ID}
		if _, err := getGuildVoiceRegions.Send(bot); err != nil {
			return fmt.Errorf("GetGuildVoiceRegions: %w", err)
		}

		return nil
	})

	eg.Go(func() error {
		getGuildInvites := &GetGuildInvites{GuildID: guild.ID}
		if _, err := getGuildInvites.Send(bot); err != nil {
			return fmt.Errorf("GetGuildInvites: %w", err)
		}

		return nil
	})

	eg.Go(func() error {
		getGuildWidgetSettings := &GetGuildWidgetSettings{GuildID: guild.ID}
		if _, err := getGuildWidgetSettings.Send(bot); err != nil {
			return fmt.Errorf("GetGuildWidgetSettings: %w", err)
		}

		return nil
	})

	// test all endpoints involving a guild.
	eg.Go(func() error {
		getGuildAuditLog := &GetGuildAuditLog{GuildID: guild.ID}
		if _, err := getGuildAuditLog.Send(bot); err != nil {
			return fmt.Errorf("GetGuildAuditLog: %w", err)
		}

		return nil
	})
	eg.Go(func() error {
		listGuildEmojis := &ListGuildEmojis{GuildID: guild.ID}
		if _, err := listGuildEmojis.Send(bot); err != nil {
			return fmt.Errorf("ListGuildEmojis: %w", err)
		}

		return nil
	})

	eg.Go(func() error {
		getGuildTemplates := &GetGuildTemplates{GuildID: guild.ID}
		if _, err := getGuildTemplates.Send(bot); err != nil {
			return fmt.Errorf("GetGuildTemplates: %w", err)
		}

		return nil
	})

	eg.Go(func() error {
		listGuildStickers := &ListGuildStickers{GuildID: guild.ID}
		if _, err := listGuildStickers.Send(bot); err != nil {
			return fmt.Errorf("ListGuildStickers: %w", err)
		}

		return nil
	})

	// wait until all requests have been processed.
	select {
	case <-ctx.Done():
		return fmt.Errorf("%w", eg.Wait())
	default:
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("%w", eg.Wait())
	}

	return nil
}

// testGuildMember tests all endpoints involving a Guild Member.
func testGuildMember(bot *Client, guild *Guild) error {
	// get the User ID of the bot.
	user, err := new(GetCurrentUser).Send(bot)
	if err != nil {
		return fmt.Errorf("Guild.GetCurrentUser: %w", err)
	}

	if user.ID == "" {
		return fmt.Errorf("GetCurrentUser: expected non-empty Guild ID.")
	}

	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		getGuildMember := &GetGuildMember{
			GuildID: guild.ID,
			UserID:  user.ID,
		}

		member, err := getGuildMember.Send(bot)
		if err != nil {
			return fmt.Errorf("GetGuildMember: %w", err)
		}

		if member == nil {
			return fmt.Errorf("GetGuildMember: expected non-null GuildMember object")
		}

		return nil
	})

	// test all endpoints involving a Guild Role.
	eg.Go(func() error {
		createGuildRole := &CreateGuildRole{
			GuildID: guild.ID,
			Name:    "testing",
			Hoist:   Pointer(true),
		}

		role, err := createGuildRole.Send(bot)
		if err != nil {
			return fmt.Errorf("CreateGuildRole: %w", err)
		}

		if role == nil {
			return fmt.Errorf("CreateGuildRole: expected non-null Role object")
		}

		rEG, rCTX := errgroup.WithContext(context.Background())

		rEG.Go(func() error {
			getGuildRoles := &GetGuildRoles{GuildID: guild.ID}
			roles, err := getGuildRoles.Send(bot)
			if err != nil {
				return fmt.Errorf("GetGuildRoles: %w", err)
			}

			if len(roles) == 0 {
				return fmt.Errorf("GetGuildRoles: expected non-empty roles slice")
			}

			return nil
		})

		rEG.Go(func() error {
			modifyGuildRole := &ModifyGuildRole{
				GuildID: guild.ID,
				RoleID:  role.ID,
				Name:    Pointer("testing..."),
			}

			if _, err := modifyGuildRole.Send(bot); err != nil {
				return fmt.Errorf("ModifyGuildRole: %w", err)
			}

			return nil
		})

		rEG.Go(func() error {
			addGuildMemberRole := &AddGuildMemberRole{
				GuildID: guild.ID,
				UserID:  user.ID,
				RoleID:  role.ID,
			}

			if err := addGuildMemberRole.Send(bot); err != nil {
				return fmt.Errorf("AddGuildMemberRole: %w", err)
			}

			removeGuildMemberRole := &RemoveGuildMemberRole{
				GuildID: guild.ID,
				UserID:  user.ID,
				RoleID:  role.ID,
			}

			if err := removeGuildMemberRole.Send(bot); err != nil {
				return fmt.Errorf("RemoveGuildMemberRole: %w", err)
			}

			return nil
		})

		// wait until role requests have been processed.
		select {
		case <-rCTX.Done():
			return fmt.Errorf("%w", rEG.Wait())
		default:
		}

		if err := rEG.Wait(); err != nil {
			return fmt.Errorf("%w", rEG.Wait())
		}

		deleteGuildRole := &DeleteGuildRole{
			GuildID: guild.ID,
			RoleID:  role.ID,
		}

		if err := deleteGuildRole.Send(bot); err != nil {
			return fmt.Errorf("DeleteGuildRole: %w", err)
		}

		return nil
	})

	// wait until all requests have been processed.
	select {
	case <-ctx.Done():
		return fmt.Errorf("%w", eg.Wait())
	default:
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("%w", eg.Wait())
	}

	return nil
}

// testGuildScheduledEvent tests all endpoints involving a scheduled event.
func testGuildScheduledEvent(bot *Client, guild *Guild) error {
	day := time.Hour * time.Duration(24)
	tomorrow := time.Now().Add(day)
	overmorrow := tomorrow.Add(day)

	createGuildScheduledEvent := &CreateGuildScheduledEvent{
		GuildID:   guild.ID,
		ChannelID: nil,
		EntityMetadata: &GuildScheduledEventEntityMetadata{
			Location: "Test",
		},
		Name:               "Test Event",
		PrivacyLevel:       FlagGuildScheduledEventPrivacyLevelGUILD_ONLY,
		ScheduledStartTime: tomorrow.Format(TimestampFormatISO8601),
		ScheduledEndTime:   overmorrow.Format(TimestampFormatISO8601),
		Description:        Pointer("A test event."),
		EntityType:         FlagGuildScheduledEventEntityTypeEXTERNAL,
		Image:              nil,
	}

	event, err := createGuildScheduledEvent.Send(bot)
	if err != nil {
		return fmt.Errorf("CreateGuildScheduledEvent: %w", err)
	}

	if event == nil {
		return fmt.Errorf("CreateGuildScheduledEvent: expected non-null GuildScheduledEvent object")
	}

	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		listScheduledEvents := &ListScheduledEventsforGuild{GuildID: guild.ID}
		events, err := listScheduledEvents.Send(bot)
		if err != nil {
			return fmt.Errorf("ListScheduledEventsforGuild: %w", err)
		}

		if len(events) == 0 {
			return fmt.Errorf("ListScheduledEventsforGuild: expected non-empty GuildScheduledEvent slice")
		}

		return nil
	})

	eg.Go(func() error {
		getGuildScheduledEvent := &GetGuildScheduledEvent{
			GuildID:               guild.ID,
			GuildScheduledEventID: event.ID,
		}

		if _, err := getGuildScheduledEvent.Send(bot); err != nil {
			return fmt.Errorf("GetGuildScheduledEvent: %w", err)
		}

		return nil
	})

	eg.Go(func() error {
		modifyGuildScheduledEvent := &ModifyGuildScheduledEvent{
			GuildID:               guild.ID,
			GuildScheduledEventID: event.ID,
			ChannelID:             nil,
			Name:                  Pointer("Test Event Modified"),
		}

		event, err := modifyGuildScheduledEvent.Send(bot)
		if err != nil {
			return fmt.Errorf("ModifyGuildScheduledEvent: %w", err)
		}

		if event == nil {
			return fmt.Errorf("ModifyGuildScheduledEvent: expected non-null GuildScheduledEvent object")
		}

		return nil
	})

	eg.Go(func() error {
		getGuildScheduledEventUsers := &GetGuildScheduledEventUsers{
			GuildID:               guild.ID,
			GuildScheduledEventID: event.ID,
		}

		_, err := getGuildScheduledEventUsers.Send(bot)
		if err != nil {
			return fmt.Errorf("GetGuildScheduledEventUsers: %w", err)
		}

		return nil
	})

	// wait until all requests have been processed.
	select {
	case <-ctx.Done():
		return fmt.Errorf("%w", eg.Wait())
	default:
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("%w", eg.Wait())
	}

	deleteGuildScheduledEvent := &DeleteGuildScheduledEvent{
		GuildID:               guild.ID,
		GuildScheduledEventID: event.ID,
	}

	if err := deleteGuildScheduledEvent.Send(bot); err != nil {
		return fmt.Errorf("DeleteGuildScheduledEvent: %w", err)
	}

	return nil
}

// testChannel tests all endpoints dependent on a channel.
func testChannel(bot *Client) error {
	/*
		CreateGuildChannel
		GetChannel
		ModifyChannelGuild
		ModifyChannel
		DeleteCloseChannel
		EditChannelPermissions
		DeleteChannelPermission
		CreateChannelInvite
		GetChannelInvites

		StartThreadfromMessage
		ListPublicArchivedThreads
		ListPrivateArchivedThreads
		ListJoinedPrivateArchivedThreads

		ModifyCurrentUserVoiceState
		ModifyUserVoiceState

		CreateStageInstance
		ModifyStageInstance
		DeleteStageInstance
		GetStageInstance
	*/

	return nil
}

// testMessage tests all endpoints dependent on a message.
func testMessage(bot *Client) error {
	/*
		CreateMessage
		GetChannelMessages
	*/

	return nil
}
