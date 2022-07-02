package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
)

type gatewayHandlerGuildUpdate struct{}

func (h *gatewayHandlerGuildUpdate) EventType() gateway.EventType {
	return gateway.EventTypeGuildUpdate
}

func (h *gatewayHandlerGuildUpdate) New() any {
	return &discord.Guild{}
}

func (h *gatewayHandlerGuildUpdate) HandleGatewayEvent(client bot.Client, sequenceNumber int, shardID int, v any) {
	guild := *v.(*discord.Guild)

	oldGuild, _ := client.Caches().Guilds().Get(guild.ID)
	client.Caches().Guilds().Put(guild.ID, guild)

	client.EventManager().DispatchEvent(&events.GuildUpdate{
		GenericGuild: &events.GenericGuild{
			GenericEvent: events.NewGenericEvent(client, sequenceNumber, shardID),
			Guild:        guild,
		},
		OldGuild: oldGuild,
	})

}
