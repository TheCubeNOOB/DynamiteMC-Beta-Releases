package server

import (
	"github.com/aimjel/minecraft/packet"
	"github.com/dynamitemc/dynamite/server/network/handlers"
	"time"
)

func (p *PlayerController) HandlePackets() error {
	ticker := time.NewTicker(25 * time.Second)
	for {
		select {
		case <-ticker.C:
			p.Keepalive()
		default:
		}

		packt, err := p.conn.ReadPacket()
		if err != nil {
			return err
		}

		switch pk := packt.(type) {
		case *packet.PlayerCommandServer:
			handlers.PlayerCommand(p, pk.ActionID)
		case *packet.ChatMessageServer:
			handlers.ChatMessagePacket(p, pk.Message)
		case *packet.ChatCommandServer:
			handlers.ChatCommandPacket(p, p.Server.commandGraph, pk.Command)
		case *packet.ClientSettings:
			handlers.ClientSettings(p, pk)
		case *packet.PlayerPosition, *packet.PlayerPositionRotation, *packet.PlayerRotation:
			handlers.PlayerMovement(p, p.player, pk)
		case *packet.PlayerActionServer:
			handlers.PlayerAction(p, pk)
		case *packet.InteractServer:
			handlers.Interact(p, pk)
		case *packet.SwingArmServer:
			handlers.SwingArm(p, pk.Hand)
		case *packet.CommandSuggestionsRequest:
			handlers.CommandSuggestionsRequest(pk.TransactionId, pk.Text, p.Server.commandGraph, p)
		case *packet.ClientCommandServer:
			handlers.ClientCommand(p, p.player, pk.ActionID)
		case *packet.PlayerAbilitiesServer:
			handlers.PlayerAbilities(p.player, pk.Flags)
		}
	}
}
