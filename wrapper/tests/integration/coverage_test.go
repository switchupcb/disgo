package integration_test

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/rs/zerolog"
	. "github.com/switchupcb/disgo/wrapper"
)

// Error Format Strings
const (
	errEmptyID    = "%v: expected non-empty %v ID"
	errNullObject = "%v: expected non-null %v object"
	errEmptySlice = "%v: expected non-empty %v slice"
)

// eventErrorGroup represents an error group for event handlers.
type eventErrorGroup struct {
	errors []error
	*sync.Mutex
}

// append appends an error the event error group.
func (e *eventErrorGroup) append(err error) {
	e.Lock()
	e.errors = append(e.errors, err)
	e.Unlock()
}

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
		t.Fatalf(errEmptyID, "GetCurrentBotApplicationInformation", "Application")
	}

	eventHandlerErrorGroup, err := initializeEventHandlers(bot)
	if err != nil {
		t.Fatal(fmt.Errorf("error setting up event handlers: %w", err))
	}

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
			return fmt.Errorf(errEmptySlice, "ListVoiceRegions", "Voice Regions")
		}

		return nil
	})

	eg.Go(func() error {
		if _, err := new(ListNitroStickerPacks).Send(bot); err != nil {
			return fmt.Errorf("ListNitroStickerPacks: %w", err)
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

	// wait until all required requests have been processed.
	select {
	case <-ctx.Done():
		t.Fatalf("%v", eg.Wait())
	default:
	}

	if err := eg.Wait(); err != nil {
		t.Fatalf("%v", err)
	}

	for _, err := range eventHandlerErrorGroup.errors {
		t.Errorf("%v", err)
	}

	// Disconnect the session from the Discord Gateway (WebSocket Connection).
	if err := bot.Sessions[0].Disconnect(); err != nil {
		t.Fatalf("%v", err)
	}

	// allow Discord to close the session.
	<-time.After(time.Second * 5)
}

// testCommands tests all endpoints that are dependent on a global command.
func testCommands(bot *Client) error {
	createGlobalApplicationCommand := CreateGlobalApplicationCommand{
		Name:        "main",
		Description: Pointer("A basic command"),
	}

	command, err := createGlobalApplicationCommand.Send(bot)
	if err != nil {
		return fmt.Errorf("CreateGlobalApplicationCommand: %w", err)
	}

	if command.ID == "" {
		return fmt.Errorf(errNullObject, "CreateGlobalApplicationCommand", "ApplicationCommand")
	}

	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		commands, err := new(GetGlobalApplicationCommands).Send(bot)
		if err != nil {
			return fmt.Errorf("GetGlobalApplicationCommands: %w", err)
		}

		if len(commands) == 0 {
			return fmt.Errorf(errEmptySlice, "GetGlobalApplicationCommands", "ApplicationCommand")
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
			return fmt.Errorf(errNullObject, "GetGlobalApplicationCommand", "ApplicationCommand")
		}

		return nil
	})

	eg.Go(func() error {
		editGlobalApplicationCommand := &EditGlobalApplicationCommand{
			CommandID:   command.ID,
			Name:        Pointer("notmain"),
			Description: Pointer("This is not a main global command."),
		}

		editedCommand, err := editGlobalApplicationCommand.Send(bot)
		if err != nil {
			return fmt.Errorf("EditGlobalApplicationCommand: %w", err)
		}

		if editedCommand.ID == "" {
			return fmt.Errorf(errEmptyID, "EditGlobalApplicationCommand", "ApplicationCommand")
		}

		return nil
	})

	// wait until all Application Command requests have been processed.
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
		return fmt.Errorf(errEmptyID, "GetGuild", "Guild")
	}

	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		return testGuildMember(bot, guild)
	})

	eg.Go(func() error {
		return testGuildAutoModeration(bot, guild)
	})

	eg.Go(func() error {
		return testGuildScheduledEvent(bot, guild)
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
			Query:   "D",
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

	// wait until all Guild requests have been processed.
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
	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		return testGuildRole(bot, guild)
	})

	eg.Go(func() error {
		getGuildMember := &GetGuildMember{
			GuildID: guild.ID,
			UserID:  bot.ApplicationID,
		}

		member, err := getGuildMember.Send(bot)
		if err != nil {
			return fmt.Errorf("GetGuildMember: %w", err)
		}

		if member == nil {
			return fmt.Errorf(errNullObject, "GetGuildMember", "GuildMember")
		}

		return nil
	})

	// wait until all Guild Member requests have been processed.
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

// testGuildRole tests all endpoints involving a guild role.
func testGuildRole(bot *Client, guild *Guild) error {
	createGuildRole := &CreateGuildRole{
		GuildID: guild.ID,
		Name:    Pointer("testing"),
		Hoist:   Pointer(true),
	}

	role, err := createGuildRole.Send(bot)
	if err != nil {
		return fmt.Errorf("CreateGuildRole: %w", err)
	}

	if role == nil {
		return fmt.Errorf(errNullObject, "CreateGuildRole", "Role")
	}

	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		getGuildRoles := &GetGuildRoles{GuildID: guild.ID}
		roles, err := getGuildRoles.Send(bot)
		if err != nil {
			return fmt.Errorf("GetGuildRoles: %w", err)
		}

		if len(roles) == 0 {
			return fmt.Errorf(errEmptySlice, "GetGuildRoles", "Role")
		}

		return nil
	})

	eg.Go(func() error {
		modifyGuildRole := &ModifyGuildRole{
			GuildID: guild.ID,
			RoleID:  role.ID,
			Name:    Pointer2("testing..."),
		}

		modifiedRole, err := modifyGuildRole.Send(bot)
		if err != nil {
			return fmt.Errorf("ModifyGuildRole: %w", err)
		}

		if modifiedRole.ID == "" {
			return fmt.Errorf(errEmptyID, "ModifyGuildRole", "Role")
		}

		return nil
	})

	eg.Go(func() error {
		addGuildMemberRole := &AddGuildMemberRole{
			GuildID: guild.ID,
			UserID:  bot.ApplicationID,
			RoleID:  role.ID,
		}

		if err := addGuildMemberRole.Send(bot); err != nil {
			return fmt.Errorf("AddGuildMemberRole: %w", err)
		}

		removeGuildMemberRole := &RemoveGuildMemberRole{
			GuildID: guild.ID,
			UserID:  bot.ApplicationID,
			RoleID:  role.ID,
		}

		if err := removeGuildMemberRole.Send(bot); err != nil {
			return fmt.Errorf("RemoveGuildMemberRole: %w", err)
		}

		return nil
	})

	// wait until Guild Role requests have been processed.
	select {
	case <-ctx.Done():
		return fmt.Errorf("%w", eg.Wait())
	default:
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("%w", eg.Wait())
	}

	deleteGuildRole := &DeleteGuildRole{
		GuildID: guild.ID,
		RoleID:  role.ID,
	}

	if err := deleteGuildRole.Send(bot); err != nil {
		return fmt.Errorf("DeleteGuildRole: %w", err)
	}

	return nil
}

// testGuildAutoModeration tests all endpoints involving AutoModeration.
func testGuildAutoModeration(bot *Client, guild *Guild) error {
	createAutoModerationRule := &CreateAutoModerationRule{
		GuildID:     guild.ID,
		Name:        "Test",
		EventType:   FlagEventTypeMESSAGE_SEND,
		TriggerType: FlagTriggerTypeKEYWORD,
		TriggerMetadata: &TriggerMetadata{
			KeywordFilter:     []string{"!@#$%^&*"},
			RegexPatterns:     []Flag{},
			Presets:           []Flag{},
			AllowList:         []string{},
			MentionTotalLimit: 0,
		},
		Actions: []*AutoModerationAction{
			{
				Type:     FlagActionTypeBLOCK_MESSAGE,
				Metadata: nil,
			},
		},
		Enabled:        Pointer(false),
		ExemptRoles:    nil,
		ExemptChannels: nil,
	}

	rule, err := createAutoModerationRule.Send(bot)
	if err != nil {
		return fmt.Errorf("CreateAutoModerationRule: %w", err)
	}

	if rule == nil {
		return fmt.Errorf(errNullObject, "CreateAutoModerationRule", "AutoModerationRule")
	}

	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		listAutoModerationRules := &ListAutoModerationRulesForGuild{GuildID: guild.ID}
		rules, err := listAutoModerationRules.Send(bot)
		if err != nil {
			return fmt.Errorf("ListAutoModerationRulesForGuild: %w", err)
		}

		if len(rules) == 0 {
			return fmt.Errorf(errEmptySlice, "ListAutoModerationRulesForGuild", "AutoModerationRule")
		}

		return nil
	})

	eg.Go(func() error {
		getAutoModerationRule := &GetAutoModerationRule{
			GuildID:              guild.ID,
			AutoModerationRuleID: rule.ID,
		}

		got, err := getAutoModerationRule.Send(bot)
		if err != nil {
			return fmt.Errorf("GetAutoModerationRule: %w", err)
		}

		if got == nil {
			return fmt.Errorf(errNullObject, "GetAutoModerationRule", "AutoModerationRule")
		}

		return nil
	})

	eg.Go(func() error {
		modifyAutoModerationRule := &ModifyAutoModerationRule{
			GuildID:              guild.ID,
			AutoModerationRuleID: rule.ID,
			Name:                 Pointer("Testing..."),
		}

		modifiedRule, err := modifyAutoModerationRule.Send(bot)
		if err != nil {
			return fmt.Errorf("ModifyAutoModerationRule: %w", err)
		}

		if modifiedRule.ID == "" {
			return fmt.Errorf(errEmptyID, "ModifyAutoModerationRule", "AutoModerationRule")
		}

		return nil
	})

	// wait until all Guild Scheduled Event requests have been processed.
	select {
	case <-ctx.Done():
		return fmt.Errorf("%w", eg.Wait())
	default:
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("%w", eg.Wait())
	}

	deleteAutoModerationRule := &DeleteAutoModerationRule{
		GuildID:              guild.ID,
		AutoModerationRuleID: rule.ID,
	}

	if err := deleteAutoModerationRule.Send(bot); err != nil {
		return fmt.Errorf("DeleteAutoModerationRule: %w", err)
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
		ScheduledStartTime: tomorrow,
		ScheduledEndTime:   &overmorrow,
		Description:        Pointer("A test event."),
		EntityType:         Pointer(FlagGuildScheduledEventEntityTypeEXTERNAL),
		Image:              nil,
	}

	event, err := createGuildScheduledEvent.Send(bot)
	if err != nil {
		return fmt.Errorf("CreateGuildScheduledEvent: %w", err)
	}

	if event == nil {
		return fmt.Errorf(errNullObject, "CreateGuildScheduledEvent", "GuildScheduledEvent")
	}

	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		listScheduledEvents := &ListScheduledEventsforGuild{GuildID: guild.ID}
		events, err := listScheduledEvents.Send(bot)
		if err != nil {
			return fmt.Errorf("ListScheduledEventsforGuild: %w", err)
		}

		if len(events) == 0 {
			return fmt.Errorf(errEmptySlice, "ListScheduledEventsforGuild", "GuildScheduledEvent")
		}

		return nil
	})

	eg.Go(func() error {
		getGuildScheduledEvent := &GetGuildScheduledEvent{
			GuildID:               guild.ID,
			GuildScheduledEventID: event.ID,
		}

		got, err := getGuildScheduledEvent.Send(bot)
		if err != nil {
			return fmt.Errorf("GetGuildScheduledEvent: %w", err)
		}

		if got == nil {
			return fmt.Errorf(errNullObject, "GetGuildScheduledEvent", "GuildScheduledEvent")
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

		modifiedEvent, err := modifyGuildScheduledEvent.Send(bot)
		if err != nil {
			return fmt.Errorf("ModifyGuildScheduledEvent: %w", err)
		}

		if modifiedEvent.ID == "" {
			return fmt.Errorf(errEmptyID, "ModifyGuildScheduledEvent", "GuildScheduledEvent")
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

	// wait until all Guild Scheduled Event requests have been processed.
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
	createGuildChannel := &CreateGuildChannel{
		GuildID:                    os.Getenv("COVERAGE_TEST_GUILD"),
		Name:                       "Test",
		Type:                       Pointer2(FlagChannelTypeGUILD_TEXT),
		Topic:                      nil,
		Bitrate:                    nil,
		UserLimit:                  nil,
		RateLimitPerUser:           nil,
		Position:                   nil,
		PermissionOverwrites:       nil,
		ParentID:                   Pointer2(os.Getenv("COVERAGE_TEST_CATEGORY")),
		NSFW:                       nil,
		RTCRegion:                  nil,
		VideoQualityMode:           nil,
		DefaultAutoArchiveDuration: nil,
		DefaultReactionEmoji:       nil,
		AvailableTags:              nil,
		DefaultSortOrder:           nil,
	}

	channel, err := createGuildChannel.Send(bot)
	if err != nil {
		return fmt.Errorf("Channel.CreateGuildChannel: %w", err)
	}

	if channel.ID == "" {
		return fmt.Errorf(errEmptyID, "Channel.CreateGuildChannel", "Channel")
	}

	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		return testVoiceChannel(bot)
	})

	eg.Go(func() error {
		return testStageInstance(bot)
	})

	eg.Go(func() error {
		return testThread(bot, channel)
	})

	eg.Go(func() error {
		return testMessage(bot, channel)
	})

	eg.Go(func() error {
		getChannel := &GetChannel{ChannelID: channel.ID}
		got, err := getChannel.Send(bot)
		if err != nil {
			return fmt.Errorf("GetChannel: %w", err)
		}

		if got == nil {
			return fmt.Errorf(errNullObject, "GetChannel", "Channel")
		}

		return nil
	})

	eg.Go(func() error {
		modifyChannel := &ModifyChannelGuild{
			ChannelID: channel.ID,
			Name:      Pointer("Testing"),
			Type:      Pointer(FlagChannelTypeGUILD_TEXT),
			ParentID:  Pointer2(os.Getenv("COVERAGE_TEST_CATEGORY")),
		}

		modifiedChannel, err := modifyChannel.Send(bot)
		if err != nil {
			return fmt.Errorf("Channel.ModifyChannelGuild: %w", err)
		}

		if modifiedChannel.ID == "" {
			return fmt.Errorf(errEmptyID, "Channel.ModifyChannelGuild", "Channel")
		}

		return nil
	})

	// wait until all Channel requests have been processed.
	select {
	case <-ctx.Done():
		return fmt.Errorf("%w", eg.Wait())
	default:
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("%w", eg.Wait())
	}

	deleteChannel := &DeleteCloseChannel{ChannelID: channel.ID}
	if _, err := deleteChannel.Send(bot); err != nil {
		return fmt.Errorf("Channel.DeleteCloseChannel: %w", err)
	}

	return nil
}

// testVoiceChannel tests all endpoints involving a voice channel.
func testVoiceChannel(bot *Client) error {
	createVoiceChannel := &CreateGuildChannel{
		GuildID:                    os.Getenv("COVERAGE_TEST_GUILD"),
		Name:                       "Test",
		Type:                       Pointer2(FlagChannelTypeGUILD_VOICE),
		Topic:                      nil,
		Bitrate:                    nil,
		UserLimit:                  nil,
		RateLimitPerUser:           nil,
		Position:                   nil,
		PermissionOverwrites:       nil,
		ParentID:                   Pointer2(os.Getenv("COVERAGE_TEST_CATEGORY")),
		NSFW:                       nil,
		RTCRegion:                  nil,
		VideoQualityMode:           nil,
		DefaultAutoArchiveDuration: nil,
		DefaultReactionEmoji:       nil,
		AvailableTags:              nil,
		DefaultSortOrder:           nil,
	}

	channel, err := createVoiceChannel.Send(bot)
	if err != nil {
		return fmt.Errorf("VoiceChannel.CreateGuildChannel: %w", err)
	}

	if channel.ID == "" {
		return fmt.Errorf("VoiceChannel.CreateGuildChannel: expected non-empty Channel ID.")
	}

	/*
		JoinChannel (Gateway Voice)
		ModifyCurrentUserVoiceState (Mute)
		ModifyUserVoiceState (Unmute)
	*/

	deleteChannel := &DeleteCloseChannel{ChannelID: channel.ID}
	if _, err := deleteChannel.Send(bot); err != nil {
		return fmt.Errorf("Voice.DeleteCloseChannel: %w", err)
	}

	return nil
}

// testThread tests all endpoints involving a thread.
func testThread(bot *Client, channel *Channel) error {
	startThread := &StartThreadwithoutMessage{
		ChannelID:           channel.ID,
		Name:                "Test",
		AutoArchiveDuration: nil,
		Type:                Pointer(FlagChannelTypePRIVATE_THREAD),
		Invitable:           nil,
		RateLimitPerUser:    nil,
	}

	thread, err := startThread.Send(bot)
	if err != nil {
		return fmt.Errorf("StartThreadwithoutMessage: %w", err)
	}

	if thread.ID == "" {
		return fmt.Errorf(errEmptyID, "StartThreadwithoutMessage", "Channel")
	}

	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		joinThread := &JoinThread{ChannelID: thread.ID}
		if err := joinThread.Send(bot); err != nil {
			return fmt.Errorf("JoinThread: %w", err)
		}

		getThreadMember := &GetThreadMember{
			ChannelID: thread.ID,
			UserID:    bot.ApplicationID,
		}

		member, err := getThreadMember.Send(bot)
		if err != nil {
			return fmt.Errorf("GetThreadMember: %w", err)
		}

		if member == nil {
			return fmt.Errorf(errNullObject, "GetThreadMember", "ThreadMember")
		}

		leaveThread := &LeaveThread{ChannelID: thread.ID}
		if err := leaveThread.Send(bot); err != nil {
			return fmt.Errorf("LeaveThread: %w", err)
		}

		return nil
	})

	eg.Go(func() error {
		listPublicArchivedThreads := &ListPublicArchivedThreads{ChannelID: channel.ID}
		if _, err := listPublicArchivedThreads.Send(bot); err != nil {
			return fmt.Errorf("ListPublicArchivedThreads: %w", err)
		}

		return nil
	})

	eg.Go(func() error {
		listPrivateArchivedThreads := &ListPrivateArchivedThreads{ChannelID: channel.ID}
		if _, err := listPrivateArchivedThreads.Send(bot); err != nil {
			return fmt.Errorf("ListPrivateArchivedThreads: %w", err)
		}

		return nil
	})

	eg.Go(func() error {
		listJoinedPrivateArchivedThreads := &ListJoinedPrivateArchivedThreads{ChannelID: channel.ID}
		if _, err := listJoinedPrivateArchivedThreads.Send(bot); err != nil {
			return fmt.Errorf("ListJoinedPrivateArchivedThreads: %w", err)
		}

		return nil
	})

	// wait until all Thread requests have been processed.
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

// testStageInstance tests all endpoints involving a stage instance.
func testStageInstance(bot *Client) error {
	createStageChannel := &CreateGuildChannel{
		GuildID:                    os.Getenv("COVERAGE_TEST_GUILD"),
		Name:                       "Test",
		Type:                       Pointer2(FlagChannelTypeGUILD_STAGE_VOICE),
		Topic:                      nil,
		Bitrate:                    nil,
		UserLimit:                  nil,
		RateLimitPerUser:           nil,
		Position:                   nil,
		PermissionOverwrites:       nil,
		ParentID:                   Pointer2(os.Getenv("COVERAGE_TEST_CATEGORY")),
		NSFW:                       nil,
		RTCRegion:                  nil,
		VideoQualityMode:           nil,
		DefaultAutoArchiveDuration: nil,
		DefaultReactionEmoji:       nil,
		AvailableTags:              nil,
		DefaultSortOrder:           nil,
	}

	channel, err := createStageChannel.Send(bot)
	if err != nil {
		return fmt.Errorf("StageInstance.CreateGuildChannel: %w", err)
	}

	if channel.ID == "" {
		return fmt.Errorf(errEmptyID, "StageInstance.CreateGuildChannel", "Channel")
	}

	createStageInstance := &CreateStageInstance{
		ChannelID:    channel.ID,
		Topic:        "Test",
		PrivacyLevel: Pointer(FlagStageInstancePrivacyLevelGUILD_ONLY),
	}

	stage, err := createStageInstance.Send(bot)
	if err != nil {
		return fmt.Errorf("CreateStageInstance: %w", err)
	}

	if stage == nil {
		return fmt.Errorf(errNullObject, "CreateStageInstance", "StageInstance")
	}

	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		getStageInstance := &GetStageInstance{ChannelID: channel.ID}
		got, err := getStageInstance.Send(bot)
		if err != nil {
			return fmt.Errorf("GetStageInstance: %w", err)
		}

		if got == nil {
			return fmt.Errorf(errNullObject, "GetStageInstance", "StageInstance")
		}

		return nil
	})

	eg.Go(func() error {
		modifyStageInstance := &ModifyStageInstance{
			ChannelID: channel.ID,
			Topic:     Pointer("Testing"),
		}

		modifiedStage, err := modifyStageInstance.Send(bot)
		if err != nil {
			return fmt.Errorf("ModifyStageInstance: %w", err)
		}

		if modifiedStage.ID == "" {
			return fmt.Errorf(errEmptyID, "ModifyStageInstance", "StageInstance")
		}

		return nil
	})

	// wait until all Stage Instance requests have been processed.
	select {
	case <-ctx.Done():
		return fmt.Errorf("%w", eg.Wait())
	default:
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("%w", eg.Wait())
	}

	deleteStageInstance := &DeleteStageInstance{ChannelID: channel.ID}
	if err := deleteStageInstance.Send(bot); err != nil {
		return fmt.Errorf("DeleteStageInstance: %w", err)
	}

	deleteChannel := &DeleteCloseChannel{ChannelID: channel.ID}
	if _, err := deleteChannel.Send(bot); err != nil {
		return fmt.Errorf("StageInstance.DeleteCloseChannel: %w", err)
	}

	return nil
}

// testMessage tests all endpoints involving a message.
func testMessage(bot *Client, channel *Channel) error {
	createMessage := &CreateMessage{
		ChannelID: channel.ID,
		Content:   Pointer("Test."),
	}

	message, err := createMessage.Send(bot)
	if err != nil {
		return fmt.Errorf("CreateMessage: %w", err)
	}

	if message.ID == "" {
		return fmt.Errorf(errEmptyID, "CreateMessage", "Message")
	}

	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		return testReaction(bot, channel, message)
	})

	eg.Go(func() error {
		pinMessage := &PinMessage{
			ChannelID: channel.ID,
			MessageID: message.ID,
		}

		if err := pinMessage.Send(bot); err != nil {
			return fmt.Errorf("PinMessage: %w", err)
		}

		getPinnedMessages := &GetPinnedMessages{ChannelID: channel.ID}
		messages, err := getPinnedMessages.Send(bot)
		if err != nil {
			return fmt.Errorf("GetPinnedMessages: %w", err)
		}

		if len(messages) == 0 {
			return fmt.Errorf(errEmptySlice, "GetPinnedMessages", "Message")
		}

		unpinMessage := &UnpinMessage{
			ChannelID: channel.ID,
			MessageID: message.ID,
		}

		if err := unpinMessage.Send(bot); err != nil {
			return fmt.Errorf("UnpinMessage: %w", err)
		}

		return nil
	})

	eg.Go(func() error {
		getChannelMessages := &GetChannelMessages{ChannelID: channel.ID}
		messages, err := getChannelMessages.Send(bot)
		if err != nil {
			return fmt.Errorf("GetChannelMessages: %w", err)
		}

		if len(messages) == 0 {
			return fmt.Errorf(errEmptySlice, "GetChannelMessages", "Message")
		}

		return nil
	})

	eg.Go(func() error {
		getChannelMessage := &GetChannelMessage{
			ChannelID: channel.ID,
			MessageID: message.ID,
		}

		got, err := getChannelMessage.Send(bot)
		if err != nil {
			return fmt.Errorf("GetChannelMessage: %w", err)
		}

		if got == nil {
			return fmt.Errorf(errNullObject, "GetChannelMessage", "Message")
		}

		return nil
	})

	eg.Go(func() error {
		editMessage := &EditMessage{
			ChannelID: channel.ID,
			MessageID: message.ID,
			Content:   Pointer2("Testing..."),
		}

		if _, err := editMessage.Send(bot); err != nil {
			return fmt.Errorf("EditMessage: %w", err)
		}

		return nil
	})

	// wait until all Message requests have been processed.
	select {
	case <-ctx.Done():
		return fmt.Errorf("%w", eg.Wait())
	default:
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("%w", eg.Wait())
	}

	deleteMessage := &DeleteMessage{
		ChannelID: channel.ID,
		MessageID: message.ID,
	}

	if err := deleteMessage.Send(bot); err != nil {
		return fmt.Errorf("DeleteMessage: %w", err)
	}

	return nil
}

// testReaction tests all endpoints involving a reaction.
func testReaction(bot *Client, channel *Channel, message *Message) error {
	emoji := "✅"
	createReaction := &CreateReaction{
		ChannelID: channel.ID,
		MessageID: message.ID,
		Emoji:     emoji,
	}

	if err := createReaction.Send(bot); err != nil {
		return fmt.Errorf("CreateReaction: %w", err)
	}

	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		getReactions := &GetReactions{
			ChannelID: channel.ID,
			MessageID: message.ID,
			Emoji:     emoji,
		}

		users, err := getReactions.Send(bot)
		if err != nil {
			return fmt.Errorf("GetReactions: %w", err)
		}

		if len(users) == 0 {
			return fmt.Errorf(errEmptySlice, "GetReactions", "User")
		}

		return nil
	})

	// wait until all Reaction requests have been processed.
	select {
	case <-ctx.Done():
		return fmt.Errorf("%w", eg.Wait())
	default:
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("%w", eg.Wait())
	}

	deleteReaction := &DeleteAllReactions{
		ChannelID: channel.ID,
		MessageID: message.ID,
	}

	if err := deleteReaction.Send(bot); err != nil {
		return fmt.Errorf("DeleteAllReactions: %w", err)
	}

	return nil
}

// initializeEventHandlers initializes the event handlers necessary for this test.
func initializeEventHandlers(bot *Client) (*eventErrorGroup, error) {
	eg := new(eventErrorGroup)
	category := "Handler"

	// Connection to WebSocket.
	if err := bot.Handle(FlagGatewayEventNameGuildCreate, func(e *GuildCreate) {
		if e.ID == "" {
			eg.append(fmt.Errorf(errEmptyID, category, "GuildCreate"))
		}
	}); err != nil {
		return nil, err
	}

	// CreateGuildRole
	if err := bot.Handle(FlagGatewayEventNameGuildRoleCreate, func(e *GuildRoleCreate) {
		if e.Role.ID == "" {
			eg.append(fmt.Errorf(errEmptyID, category, "GuildRoleCreate"))
		}
	}); err != nil {
		return nil, err
	}

	// ModifyGuildRole
	if err := bot.Handle(FlagGatewayEventNameGuildRoleUpdate, func(e *GuildRoleUpdate) {
		if e.Role.ID == "" {
			eg.append(fmt.Errorf(errEmptyID, category, "GuildRoleUpdate"))
		}
	}); err != nil {
		return nil, err
	}

	// DeleteGuildRole
	if err := bot.Handle(FlagGatewayEventNameGuildRoleDelete, func(e *GuildRoleDelete) {
		if e.RoleID == "" {
			eg.append(fmt.Errorf(errEmptyID, category, "GuildRoleCreate"))
		}
	}); err != nil {
		return nil, err
	}

	// CreateAutoModerationRule
	if err := bot.Handle(FlagGatewayEventNameAutoModerationRuleCreate, func(e *AutoModerationRuleCreate) {
		if e.ID == "" {
			eg.append(fmt.Errorf(errEmptyID, category, "AutoModerationRuleCreate"))
		}
	}); err != nil {
		return nil, err
	}

	// ModifyAutoModerationRule
	if err := bot.Handle(FlagGatewayEventNameAutoModerationRuleUpdate, func(e *AutoModerationRuleUpdate) {
		if e.ID == "" {
			eg.append(fmt.Errorf(errEmptyID, category, "AutoModerationRuleUpdate"))
		}
	}); err != nil {
		return nil, err
	}

	// DeleteAutoModerationRule
	if err := bot.Handle(FlagGatewayEventNameAutoModerationRuleDelete, func(e *AutoModerationRuleDelete) {
		if e.ID == "" {
			eg.append(fmt.Errorf(errEmptyID, category, "AutoModerationRuleDelete"))
		}
	}); err != nil {
		return nil, err
	}

	// CreateGuildScheduledEvent
	if err := bot.Handle(FlagGatewayEventNameGuildScheduledEventCreate, func(e *GuildScheduledEventCreate) {
		if e.ID == "" {
			eg.append(fmt.Errorf(errEmptyID, category, "GuildScheduledEventCreate"))
		}
	}); err != nil {
		return nil, err
	}

	// ModifyGuildScheduledEvent
	if err := bot.Handle(FlagGatewayEventNameGuildScheduledEventUpdate, func(e *GuildScheduledEventUpdate) {
		if e.ID == "" {
			eg.append(fmt.Errorf(errEmptyID, category, "GuildScheduledEventUpdate"))
		}
	}); err != nil {
		return nil, err
	}

	// DeleteGuildScheduledEvent
	if err := bot.Handle(FlagGatewayEventNameGuildScheduledEventDelete, func(e *GuildScheduledEventDelete) {
		if e.ID == "" {
			eg.append(fmt.Errorf(errEmptyID, category, "GuildScheduledEventDelete"))
		}
	}); err != nil {
		return nil, err
	}

	// CreateGuildChannel
	if err := bot.Handle(FlagGatewayEventNameChannelCreate, func(e *ChannelCreate) {
		if e.ID == "" {
			eg.append(fmt.Errorf(errEmptyID, category, "ChannelCreate"))
		}
	}); err != nil {
		return nil, err
	}

	// ModifyChannel
	if err := bot.Handle(FlagGatewayEventNameChannelUpdate, func(e *ChannelUpdate) {
		if e.ID == "" {
			eg.append(fmt.Errorf(errEmptyID, category, "ChannelUpdate"))
		}
	}); err != nil {
		return nil, err
	}

	// DeleteCloseChannel
	if err := bot.Handle(FlagGatewayEventNameChannelDelete, func(e *ChannelDelete) {
		if e.ID == "" {
			eg.append(fmt.Errorf(errEmptyID, category, "ChannelDelete"))
		}
	}); err != nil {
		return nil, err
	}

	if err := bot.Handle(FlagGatewayEventNameThreadDelete, func(e *ThreadDelete) {
		if e.ID == "" {
			eg.append(fmt.Errorf(errEmptyID, category, "ThreadDelete"))
		}
	}); err != nil {
		return nil, err
	}

	// StartThreadwithoutMessage
	if err := bot.Handle(FlagGatewayEventNameThreadCreate, func(e *ThreadCreate) {
		if e.ID == "" {
			eg.append(fmt.Errorf(errEmptyID, category, "ThreadCreate"))
		}
	}); err != nil {
		return nil, err
	}

	// JoinThread, LeaveThread
	if err := bot.Handle(FlagGatewayEventNameThreadMembersUpdate, func(e *ThreadMembersUpdate) {
		if e.ID == "" {
			eg.append(fmt.Errorf(errEmptyID, category, "ThreadMembersUpdate"))
		}
	}); err != nil {
		return nil, err
	}

	// CreateStageInstance
	if err := bot.Handle(FlagGatewayEventNameStageInstanceCreate, func(e *StageInstanceCreate) {
		if e.ID == "" {
			eg.append(fmt.Errorf(errEmptyID, category, "StageInstanceCreate"))
		}
	}); err != nil {
		return nil, err
	}

	// ModifyStageInstance
	if err := bot.Handle(FlagGatewayEventNameStageInstanceUpdate, func(e *StageInstanceUpdate) {
		if e.ID == "" {
			eg.append(fmt.Errorf(errEmptyID, category, "StageInstanceUpdate"))
		}
	}); err != nil {
		return nil, err
	}

	// DeleteStageInstance
	if err := bot.Handle(FlagGatewayEventNameStageInstanceDelete, func(e *StageInstanceDelete) {
		if e.ID == "" {
			eg.append(fmt.Errorf(errEmptyID, category, "StageInstanceDelete"))
		}
	}); err != nil {
		return nil, err
	}

	// CreateMessage
	if err := bot.Handle(FlagGatewayEventNameMessageCreate, func(e *MessageCreate) {
		if e.ID == "" {
			eg.append(fmt.Errorf(errEmptyID, category, "MessageCreate"))
		}
	}); err != nil {
		return nil, err
	}

	// EditMessage
	if err := bot.Handle(FlagGatewayEventNameMessageUpdate, func(e *MessageUpdate) {
		if e.ID == "" {
			eg.append(fmt.Errorf(errEmptyID, category, "MessageUpdate"))
		}
	}); err != nil {
		return nil, err
	}

	// DeleteMessage
	if err := bot.Handle(FlagGatewayEventNameMessageDelete, func(e *MessageDelete) {
		if e.MessageID == "" {
			eg.append(fmt.Errorf(errEmptyID, category, "MessageDelete"))
		}
	}); err != nil {
		return nil, err
	}

	// PinMessage, UnpinMessage
	if err := bot.Handle(FlagGatewayEventNameChannelPinsUpdate, func(e *ChannelPinsUpdate) {
		if e.ChannelID == "" {
			eg.append(fmt.Errorf(errEmptyID, category, "ChannelPinsUpdate"))
		}
	}); err != nil {
		return nil, err
	}

	// CreateReaction
	if err := bot.Handle(FlagGatewayEventNameMessageReactionAdd, func(e *MessageReactionAdd) {
		if e.ChannelID == "" {
			eg.append(fmt.Errorf(errEmptyID, category, "MessageReactionAdd"))
		}
	}); err != nil {
		return nil, err
	}

	// DeleteAllReactions
	if err := bot.Handle(FlagGatewayEventNameMessageReactionRemoveAll, func(e *MessageReactionRemoveAll) {
		if e.ChannelID == "" {
			eg.append(fmt.Errorf(errEmptyID, category, "MessageReactionRemoveAll"))
		}
	}); err != nil {
		return nil, err
	}

	return eg, nil
}
