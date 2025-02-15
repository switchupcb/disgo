package wrapper

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"slices"
	"strconv"
)

// connectUDP connects to the Discord UDP Voice Server using the given Ready payload.
//
// https://discord.com/developers/docs/topics/voice-connections#establishing-a-voice-udp-connection
func (vc *VoiceChannelConnection) connectUDP(r *VoiceReady) error {
	var err error

	// Open a UDP Connection to the provided IP and port.
	address := r.IP + ":" + strconv.Itoa(r.Port)
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return fmt.Errorf("udp: %w", err)
	}

	if vc.Connection, err = net.DialUDP("udp", nil, udpAddr); err != nil {
		return fmt.Errorf("udp: %w", err)
	}

	// Perform an IP Discovery.
	// https://discord.com/developers/docs/topics/voice-connections#ip-discovery
	ipDiscoveryPacket := make([]byte, 74)                              //nolint:mnd // TODO add numbers as dasgo flags
	binary.BigEndian.PutUint16(ipDiscoveryPacket, 1)                   //nolint:nolintlint //nolint:mnd // Type: 0x1 = request, 0x2 = response
	binary.BigEndian.PutUint16(ipDiscoveryPacket[2:4], 70)             //nolint:mnd // Message Length: 70
	binary.BigEndian.PutUint32(ipDiscoveryPacket[4:8], uint32(r.SSRC)) //nolint:gosec // disable G115 // SSRC
	if _, err := vc.Connection.Write(ipDiscoveryPacket); err != nil {
		return fmt.Errorf("udp: %w", err)
	}

	ipDiscoveryPacket = make([]byte, 74) //nolint:mnd
	_, externalAddr, err := vc.Connection.ReadFromUDP(ipDiscoveryPacket)
	if err != nil {
		return fmt.Errorf("udp: %w", err)
	}

	// Send the client's external IP and UDP Port to the Discord Voice WebSocket.
	//
	// select a supported encryption mode (in order of priority).
	// https://discord.com/developers/docs/topics/voice-connections#transport-encryption-and-sending-voice
	var mode string
	switch {
	case slices.Contains(r.Modes, FlagVoiceEncryptionModeAES256):
		mode = FlagVoiceEncryptionModeAES256
	case slices.Contains(r.Modes, FlagVoiceEncryptionModeXChaCha20):
		mode = FlagVoiceEncryptionModeXChaCha20
	default:
		return errors.New("udp: supported mode is not available")
	}

	// send an Opcode 1 Select Protocol Payload.
	selectProtocol := &SelectProtocol{
		Protocol: "udp",
		Data: SelectProtocolData{
			Address: externalAddr.IP.String(),
			Port:    externalAddr.Port,
			Mode:    mode,
		},
	}

	if err := selectProtocol.SendEvent(vc.VoiceSession); err != nil {
		return fmt.Errorf("udp: %w", err)
	}

	// TODO: DAVE
	// https://discord.com/developers/docs/topics/voice-connections#endtoend-encryption-dave-protocol

	// TODO: Connection is established, create routine for external library to process Voice Connection Data
	// https://discord.com/developers/docs/topics/voice-connections#encrypting-and-sending-voice
	// go func()
	// VoicePacket...

	return nil
}
