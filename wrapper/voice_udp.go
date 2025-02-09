package wrapper

import (
	"encoding/binary"
	"fmt"
	"net"

	"golang.org/x/exp/slices"
)

// connectUDP connects to the Discord UDP Voice Server using the given Ready payload.
//
// https://discord.com/developers/docs/topics/voice-connections#establishing-a-voice-udp-connection
func (vc *VoiceConnection) connectUDP(r *VoiceReady) error {
	var err error

	// Open a UDP Connection to the provided IP and port.
	if vc.Connection, err = net.DialUDP(
		"udp",
		nil,
		&net.UDPAddr{
			IP:   net.IP(r.IP),
			Port: r.Port,
			Zone: "",
		},
	); err != nil {
		return fmt.Errorf("udp: %w", err)
	}

	// Perform an IP Discovery.
	// https://discord.com/developers/docs/topics/voice-connections#ip-discovery
	ipDiscoveryPacket := make([]byte, 74)
	binary.BigEndian.PutUint16(ipDiscoveryPacket, 1)                   // Type: 0x1 = request, 0x2 = response
	binary.BigEndian.PutUint16(ipDiscoveryPacket[2:4], 70)             // Message Length: 70
	binary.BigEndian.PutUint32(ipDiscoveryPacket[4:8], uint32(r.SSRC)) // SSRC
	vc.Connection.Write(ipDiscoveryPacket)

	ipDiscoveryPacket = make([]byte, 74)
	_, eaddr, err := vc.Connection.ReadFromUDP(ipDiscoveryPacket)
	if err != nil {
		return fmt.Errorf("udp: %w", err)
	}

	// Send the client's external IP and UDP Port to the Discord Voice WebSocket.
	// https://discord.com/developers/docs/topics/opcodes-and-status-codes#voice
	//
	// select a supported encryption mode (in order of priority).
	// https://discord.com/developers/docs/topics/voice-connections#establishing-a-voice-udp-connection-encryption-modes
	var mode string
	if slices.Contains(r.Modes, FlagVoiceEncryptionModeLite) {
		mode = FlagVoiceEncryptionModeLite
	} else if slices.Contains(r.Modes, FlagVoiceEncryptionModeSuffix) {
		mode = FlagVoiceEncryptionModeSuffix
	} else if slices.Contains(r.Modes, FlagVoiceEncryptionModeNormal) {
		mode = FlagVoiceEncryptionModeNormal
	} else {
		return fmt.Errorf("udp: supported mode is not available")
	}

	// TODO: also support modes on encryption and decryption

	// send an Opcode 1 Select Protocol Payload.
	selectProtocol := &SelectProtocol{
		Protocol: "udp",
		Data: SelectProtocolData{
			Address: eaddr.IP.String(),
			Port:    eaddr.Port,
			Mode:    mode,
		},
	}

	// Set up the event handler for the Opcode 4 Session Description event.
	if err := vc.Handle(FlagVoiceOpcodeNameSessionDescription, func(sd *SessionDescription) {
		// sd.SecretKey
	}); err != nil {
		return fmt.Errorf("udp: %w", err)
	}

	if err := selectProtocol.SendEvent(vc.VoiceSession); err != nil {
		return fmt.Errorf("udp: %w", err)
	}

	return nil
}

// Connection is established, create channel for external library to process voice
// https://discord.com/developers/docs/topics/voice-connections#encrypting-and-sending-voice
