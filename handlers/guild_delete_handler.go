package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
)

type gatewayHandlerGuildDelete struct{}

func (h *gatewayHandlerGuildDelete) EventType() gateway.EventType {
	return gateway.EventTypeGuildDelete
}

func (h *gatewayHandlerGuildDelete) New() any {
	return &discord.UnavailableGuild{}
}

func (h *gatewayHandlerGuildDelete) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	unavailableGuild := *v.(*discord.UnavailableGuild)

	guild, _ := client.Caches().Guilds().Remove(unavailableGuild.ID)
	client.Caches().VoiceStates().RemoveAll(unavailableGuild.ID)
	client.Caches().Presences().RemoveAll(unavailableGuild.ID)
	client.Caches().ThreadMembers().RemoveIf(func(_ snowflake.ID, threadMember discord.ThreadMember) bool {
		// TODO: figure out a better way to remove thread members from cache via guild id without requiring cached GuildThreads
		if thread, ok := client.Caches().Channels().GetGuildThread(threadMember.ThreadID); ok {
			return thread.GuildID() == unavailableGuild.ID
		}
		return false
	})
	client.Caches().Channels().RemoveIf(func(channel discord.Channel) bool {
		if guildChannel, ok := channel.(discord.GuildChannel); ok {
			return guildChannel.GuildID() == unavailableGuild.ID
		}
		return false
	})
	client.Caches().Emojis().RemoveAll(unavailableGuild.ID)
	client.Caches().Stickers().RemoveAll(unavailableGuild.ID)
	client.Caches().Roles().RemoveAll(unavailableGuild.ID)
	client.Caches().StageInstances().RemoveAll(unavailableGuild.ID)

	client.Caches().Messages().RemoveIf(func(channelID snowflake.ID, message discord.Message) bool {
		return message.GuildID != nil && *message.GuildID == unavailableGuild.ID
	})

	if unavailableGuild.Unavailable {
		client.Caches().Guilds().SetUnavailable(unavailableGuild.ID)
	}

	genericGuildEvent := &events.GenericGuild{
		GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
		GuildID:      unavailableGuild.ID,
		Guild:        guild,
	}

	if unavailableGuild.Unavailable {
		client.EventManager().DispatchEvent(&events.GuildUnavailable{
			GenericGuild: genericGuildEvent,
		})
	} else {
		client.EventManager().DispatchEvent(&events.GuildLeave{
			GenericGuild: genericGuildEvent,
		})
	}
}
