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
			return fmt.Errorf("ListVoiceRegions: expected non-empty Voice Regions slice")
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
		Description: Pointer("A basic command"),
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
			Name:        Pointer("notmain"),
			Description: Pointer("This is not a main global command."),
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
			return fmt.Errorf("GetGuildMember: expected non-null GuildMember object")
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
		return fmt.Errorf("CreateGuildRole: expected non-null Role object")
	}

	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
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

	eg.Go(func() error {
		modifyGuildRole := &ModifyGuildRole{
			GuildID: guild.ID,
			RoleID:  role.ID,
			Name:    Pointer2("testing..."),
		}

		if _, err := modifyGuildRole.Send(bot); err != nil {
			return fmt.Errorf("ModifyGuildRole: %w", err)
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
		return fmt.Errorf("Channel.CreateGuildChannel: expected non-empty Channel ID.")
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
			return fmt.Errorf("GetChannel: expected non-null Channel object")
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

		if _, err := modifyChannel.Send(bot); err != nil {
			return fmt.Errorf("ModifyChannelGuild: %w", err)
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
		return fmt.Errorf("StartThreadwithoutMessage: expected non-empty Channel ID.")
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
			return fmt.Errorf("GetThreadMember: expected non-null ThreadMember object")
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
		return fmt.Errorf("StageInstance.CreateGuildChannel: expected non-empty Channel ID.")
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
		return fmt.Errorf("CreateStageInstance: expected non-null StageInstance object")
	}

	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(func() error {
		getStageInstance := &GetStageInstance{ChannelID: channel.ID}
		got, err := getStageInstance.Send(bot)
		if err != nil {
			return fmt.Errorf("GetStageInstance: %w", err)
		}

		if got == nil {
			return fmt.Errorf("GetStageInstance: expected non-null StageInstance object")
		}

		return nil
	})

	eg.Go(func() error {
		modifyStageInstance := &ModifyStageInstance{
			ChannelID: channel.ID,
			Topic:     Pointer("Testing"),
		}

		if _, err := modifyStageInstance.Send(bot); err != nil {
			return fmt.Errorf("ModifyStageInstance: %w", err)
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
		return fmt.Errorf("CreateMessage: expected non-empty Message ID.")
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
			return fmt.Errorf("GetPinnedMessages: expected non-empty Message slice")
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
			return fmt.Errorf("GetChannelMessages: expected non-empty Message slice")
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
			return fmt.Errorf("GetChannelMessage: expected non-null Message object")
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
			return fmt.Errorf("GetReactions: expected non-empty User slice")
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
