package wrapper

// TODO: note that Group DM requests are not being used

var (
	endpoints = map[string][]string{
		"GetGlobalApplicationCommands":           {"CreateGlobalApplicationCommand"},
		"CreateGlobalApplicationCommand":         {"CreateGlobalApplicationCommand"},
		"GetGlobalApplicationCommand":            {"CreateGlobalApplicationCommand"},
		"EditGlobalApplicationCommand":           {"CreateGlobalApplicationCommand"},
		"DeleteGlobalApplicationCommand":         {"CreateGlobalApplicationCommand"},
		"BulkOverwriteGlobalApplicationCommands": {"CreateGlobalApplicationCommand"},
		"GetGuildApplicationCommands":            {},
		"CreateGuildApplicationCommand":          {},
		"GetGuildApplicationCommand":             {"CreateGuildApplicationCommand"},
		"EditGuildApplicationCommand":            {"CreateGuildApplicationCommand"},
		"DeleteGuildApplicationCommand":          {"CreateGuildApplicationCommand"},
		"BulkOverwriteGuildApplicationCommands":  {"CreateGuildApplicationCommand"},
		"GetGuildApplicationCommandPermissions":  {},
		"GetApplicationCommandPermissions":       {},
		"EditApplicationCommandPermissions":      {},
		"BatchEditApplicationCommandPermissions": {},
		"CreateInteractionResponse":              {},
		"GetOriginalInteractionResponse":         {"CreateInteractionResponse"},
		"EditOriginalInteractionResponse":        {"CreateInteractionResponse"},
		"DeleteOriginalInteractionResponse":      {"CreateInteractionResponse"},
		"CreateFollowupMessage":                  {},
		"GetFollowupMessage":                     {"CreateFollowupMessage"},
		"EditFollowupMessage":                    {"CreateFollowupMessage"},
		"DeleteFollowupMessage":                  {"CreateFollowupMessage"},
		"GetGuildAuditLog":                       {"CreateGuild"},
		"ListAutoModerationRulesForGuild":        {"CreateAutoModerationRule"},
		"GetAutoModerationRule":                  {"CreateAutoModerationRule"},
		"CreateAutoModerationRule":               {"CreateGuild"},
		"ModifyAutoModerationRule":               {"CreateAutoModerationRule"},
		"DeleteAutoModerationRule":               {"CreateAutoModerationRule"},
		"GetChannel":                             {"CreateGuild"},
		"ModifyChannel":                          {"CreateGuild"},
		"ModifyChannelGuild":                     {"CreateGuild"},
		"ModifyChannelThread":                    {},
		"DeleteCloseChannel":                     {},
		"GetChannelMessages":                     {},
		"GetChannelMessage":                      {},
		"CreateMessage":                          {},
		"CrosspostMessage":                       {},
		"CreateReaction":                         {},
		"DeleteOwnReaction":                      {"CreateReaction"},
		"DeleteUserReaction":                     {"CreateReaction"},
		"GetReactions":                           {"CreateReaction"},
		"DeleteAllReactions":                     {"CreateReaction"},
		"DeleteAllReactionsforEmoji":             {},
		"EditMessage":                            {"CreateMessage"},
		"DeleteMessage":                          {"CreateMessage"},
		"BulkDeleteMessages":                     {"CreateMessage"},
		"EditChannelPermissions":                 {},
		"GetChannelInvites":                      {},
		"CreateChannelInvite":                    {},
		"DeleteChannelPermission":                {},
		"FollowNewsChannel":                      {},
		"TriggerTypingIndicator":                 {},
		"GetPinnedMessages":                      {},
		"PinMessage":                             {"CreateMessage"},
		"UnpinMessage":                           {"CreateMessage"},
		"StartThreadfromMessage":                 {"CreateMessage"},
		"StartThreadwithoutMessage":              {},
		"StartThreadinForumChannel":              {},
		"JoinThread":                             {"StartThreadfromMessage"},
		"AddThreadMember":                        {"StartThreadfromMessage"},
		"LeaveThread":                            {"StartThreadfromMessage"},
		"RemoveThreadMember":                     {"StartThreadfromMessage"},
		"GetThreadMember":                        {"StartThreadfromMessage"},
		"ListThreadMembers":                      {"StartThreadfromMessage"},
		"ListPublicArchivedThreads":              {"StartThreadfromMessage"},
		"ListPrivateArchivedThreads":             {"StartThreadfromMessage"},
		"ListJoinedPrivateArchivedThreads":       {"StartThreadfromMessage"},
		"ListGuildEmojis":                        {},
		"GetGuildEmoji":                          {},
		"CreateGuildEmoji":                       {},
		"ModifyGuildEmoji":                       {},
		"DeleteGuildEmoji":                       {},
		"CreateGuild":                            {},
		"GetGuild":                               {"CreateGuild"},
		"GetGuildPreview":                        {"CreateGuild"},
		"ModifyGuild":                            {"CreateGuild"},
		"DeleteGuild":                            {"CreateGuild"},
		"GetGuildChannels":                       {"CreateGuild"},
		"CreateGuildChannel":                     {"CreateGuild"},
		"ModifyGuildChannelPositions":            {},
		"ListActiveGuildThreads":                 {"CreateGuild"},
		"GetGuildMember":                         {"CreateGuild"},
		"ListGuildMembers":                       {"CreateGuild"},
		"SearchGuildMembers":                     {"CreateGuild"},
		"AddGuildMember":                         {"CreateGuild"},
		"ModifyGuildMember":                      {"CreateGuild"},
		"ModifyCurrentMember":                    {"AddGuildMember"},
		"AddGuildMemberRole":                     {"AddGuildMember"},
		"RemoveGuildMemberRole":                  {"AddGuildMember"},
		"RemoveGuildMember":                      {"AddGuildMember"},
		"GetGuildBans":                           {"CreateGuildBan"},
		"GetGuildBan":                            {"CreateGuildBan"},
		"CreateGuildBan":                         {"CreateGuild"},
		"RemoveGuildBan":                         {"CreateGuildBan"},
		"GetGuildRoles":                          {"CreateGuild"},
		"CreateGuildRole":                        {"CreateGuild"},
		"ModifyGuildRolePositions":               {"CreateGuild"},
		"ModifyGuildRole":                        {"CreateGuildRole"},
		"DeleteGuildRole":                        {"CreateGuildRole"},
		"GetGuildPruneCount":                     {"CreateGuild"},
		"BeginGuildPrune":                        {"CreateGuild", "AddGuildMember"},
		"GetGuildVoiceRegions":                   {"CreateGuild"},
		"GetGuildInvites":                        {"CreateGuild"},
		"GetGuildIntegrations":                   {"CreateGuild"},
		"DeleteGuildIntegration":                 {"CreateGuild"},
		"GetGuildWidgetSettings":                 {"CreateGuild"},
		"ModifyGuildWidget":                      {"CreateGuild"},
		"GetGuildWidget":                         {"CreateGuild"},
		"GetGuildVanityURL":                      {"CreateGuild"},
		"GetGuildWidgetImage":                    {"CreateGuild"},
		"GetGuildWelcomeScreen":                  {"CreateGuild"},
		"ModifyGuildWelcomeScreen":               {"CreateGuild"},
		"ModifyCurrentUserVoiceState":            {},
		"ModifyUserVoiceState":                   {},
		"ListScheduledEventsforGuild":            {"CreateGuildScheduledEvent"},
		"CreateGuildScheduledEvent":              {"CreateGuild"},
		"GetGuildScheduledEvent":                 {"CreateGuildScheduledEvent"},
		"ModifyGuildScheduledEvent":              {"CreateGuildScheduledEvent"},
		"DeleteGuildScheduledEvent":              {"CreateGuildScheduledEvent"},
		"GetGuildScheduledEventUsers":            {"CreateGuildScheduledEvent"},
		"GetGuildTemplate":                       {"CreateGuildTemplate"},
		"CreateGuildfromGuildTemplate":           {},
		"GetGuildTemplates":                      {},
		"CreateGuildTemplate":                    {},
		"SyncGuildTemplate":                      {"CreateGuildTemplate"},
		"ModifyGuildTemplate":                    {"CreateGuildTemplate"},
		"DeleteGuildTemplate":                    {"CreateGuildTemplate"},
		"GetInvite":                              {},
		"DeleteInvite":                           {},
		"CreateStageInstance":                    {},
		"GetStageInstance":                       {"CreateStageInstance"},
		"ModifyStageInstance":                    {"CreateStageInstance"},
		"DeleteStageInstance":                    {"CreateStageInstance"},
		"GetSticker":                             {},
		"ListNitroStickerPacks":                  {},
		"ListGuildStickers":                      {},
		"GetGuildSticker":                        {"CreateGuildSticker"},
		"CreateGuildSticker":                     {},
		"ModifyGuildSticker":                     {"CreateGuildSticker"},
		"DeleteGuildSticker":                     {"CreateGuildSticker"},
		"GetCurrentUser":                         {},
		"GetUser":                                {},
		"ModifyCurrentUser":                      {},
		"GetCurrentUserGuilds":                   {},
		"GetCurrentUserGuildMember":              {},
		"LeaveGuild":                             {"CreateGuild"},
		"GetUserConnections":                     {"GetUser"},
		"ListVoiceRegions":                       {},
		"CreateWebhook":                          {},
		"GetChannelWebhooks":                     {},
		"GetGuildWebhooks":                       {"CreateGuild"},
		"GetWebhook":                             {"CreateWebhook"},
		"GetWebhookwithToken":                    {"CreateWebhook"},
		"ModifyWebhook":                          {"CreateWebhook"},
		"ModifyWebhookwithToken":                 {"CreateWebhook"},
		"DeleteWebhook":                          {"CreateWebhook"},
		"DeleteWebhookwithToken":                 {"CreateWebhook"},
		"ExecuteWebhook":                         {"CreateWebhook"},
		"ExecuteSlackCompatibleWebhook":          {"CreateWebhook"},
		"ExecuteGitHubCompatibleWebhook":         {"CreateWebhook"},
		"GetWebhookMessage":                      {"CreateWebhook"},
		"EditWebhookMessage":                     {"CreateWebhook"},
		"DeleteWebhookMessage":                   {"CreateWebhook"},
		"GetGateway":                             {},
		"GetGatewayBot":                          {},
		"GetCurrentBotApplicationInformation":    {},
		"GetCurrentAuthorizationInformation":     {},
	}
)

func findOrder(endpoints map[string][]string) []string {

	var numEndpoints int = len(endpoints)
	// pcmap represents a map of prerequisite courses to courses.
	pcmap := make([][]string, numEndpoints)

	// indegree represents the amount of prerequisites required for a course.
	indegree := make([]int, numEndpoints)

	// bucket represents a [course, prerequisite] ordered pair.
	for _, bucket := range endpoints {
		// map prerequisites to d.
		pcmap[bucket[1][0]] = append(pcmap[bucket[1][0]], bucket[0])

		// bucket[0] represents a course (from 0 < numCourses).
		// add to its prerequisite count.
		indegree[bucket[0][0]]++
	}

	// queue represents a first-in first-out data structure.
	//
	// queue can't have more entries than the number of courses,
	// so initialize a map of length 0 with capacity = numCourses.
	queue := make([]string, len(endpoints))

	// fill the queue with courses that have no prerequisites.
	for course, prerequisiteCount := range indegree {
		if prerequisiteCount == 0 {
			queue = append(queue, course)
		}
	}

	// output represents the returned result.
	output := make([]string, 0, len(endpoints))

	// add the courses with no prerequisites.
	for len(queue) > 0 {
		// select the first entry in the queue (of courses with no prereqs).
		current := queue[0]

		// remove the first entry out the queue (of courses with no prereqs).
		queue = queue[1:]

		// add the entry to the output.
		output = append(output, current)
		numEndpoints--

		// check courses that depend on the selected entry to the queue
		for _, course := range pcmap[current] {
			indegree[course]--

			// when the course's prerequisite count is 0,
			// add it to the queue.
			if indegree[course] == 0 {
				queue = append(queue, course)
			}
		}
	}

	// numCourses is not 0 when a cycle occurs (i.e [1,0],[0,1])
	if numEndpoints != 0 {
		return []string{}
	}

	return output
}
