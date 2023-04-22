package egts

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
)

const (
	// DefaultHeaderLen is default header length for EGTS protocol.
	DefaultHeaderLen       = 11
	allowedFirstByte1 byte = 0x01 // 1 - version
)

// Packet structure describes of the EGTS packet.
// The Transport Layer Protocol header consists of the following fields: PRV, PRF, PR, CMP, ENA, RTE, HL, HE, FDL,
// PID, PT, PRA, RCA, TTL, HCS. The Service Support Layer protocol is represented by the SFRD field,
// the checksum of the Service Support Layer field is contained in the SFRCS field.
type Packet struct {
	// ProtocolVersion (PRV) parameter contains the value 0x01.
	// The value of this parameter is incremented every time you change the header structure.
	ProtocolVersion byte `json:"PRV"`
	// SecurityKeyID (SKID) parameter defines the identifier of the key used for encryption.
	SecurityKeyID byte `json:"SKID"`
	// Prefix (PRF) parameter defines the Transport layer header prefix and contains the value 00.
	Prefix string `json:"PRF"`
	// Route (RTE) field determines the need for further routing of this packet to the
	// remote hardware and software complex, as well as the presence of optional parameters PRA, RCA, TTL,
	// necessary for routing this packet. If the field is 1, then routing is required and the PRA, RCA,
	// TTL fields are present in the packet. This field sets the Dispatcher of the hardware and software complex on
	// which the packet was generated,
	// or the subscriber terminal that generated the packet for sending to the hardware and software complex,
	// in case it is set to the parameter "HOME_DISPATCHER_ID",
	// defining the address of the hardware and software complex on which this subscriber terminal is registered.
	Route string `json:"RTE"`
	// EncryptionAlg (ENA) field specifies the algorithm code used to encrypt data from the SFRD field.
	// If the field is 00, the data in the SFRD field is not encrypted.
	EncryptionAlg string `json:"ENA"`
	// Compression (CMP) field determines whether data from the SFRD field is compressed.
	// If the field is set to 1, the data in the SFRD field is considered compressed.
	Compression string `json:"CMP"`
	// Priority (PR) field determines the routing priority of this packet and can take the following values:
	// 00 - highest,
	// 01 - high,
	// 10 - average,
	// 11 - low.
	// When a packet is received,
	// Dispatcher routes a packet with a higher priority faster than packets with a lower priority.
	Priority string `json:"PR"`
	// HeaderLength (HL) field is the length of the Transport Layer header in bytes,
	// including the checksum byte (HCS fields).
	HeaderLength byte `json:"HL"`
	// HeaderEncoding (HE) field defines the encoding method used for the next part of the Transport Layer Header
	// following this parameter.
	HeaderEncoding byte `json:"HE"`
	// FrameDataLength (FDL) field specifies the size in bytes of the SFRD data field containing Service Level Support
	// protocol information.
	FrameDataLength uint16 `json:"FDL"`
	// PacketIdentifier (PID) field contains the Transport Layer packet number,
	// increasing by 1 with each new packet sent on the sender side.
	// The values in this field change according to cyclic counter rules in the range from 0 to 65535,
	// i.e. when the value 65535 is reached the next value is 0.
	PacketIdentifier uint16 `json:"PID"`
	// PacketType (PT) field is the packet type of the Transport Layer. The PT field can take the following values.
	// 0 - EGTS_PT_RESPONSE (acknowledgement on Transport Layer packet);
	// 1 - EGTS_PT_APPDATA (packet containing Service Level Support Protocol data);
	// 2 - EGTS_PT_SIGNED_APPDATA (packet containing Service Level Support Protocol data with digital signature).
	PacketType byte `json:"PT"`
	// PeerAddress (PRA) field is the address of the hardware and software complex on which this packet
	// generated. This address is unique within the network and is used to create a confirmation packet on the
	// receiving side.
	PeerAddress uint16 `json:"PRA"`
	// RecipientAddress (RCA) field - the address of the hardware and software complex,
	// for which the package is intended. At this address the identification of the package belonging to a
	// particular hardware and software complex and its routing when using the intermediate hardware and software
	// complexes.
	RecipientAddress uint16 `json:"RCA"`
	// ThTimeToLive (TTL) field is the packet lifetime when routing between hardware and software complexes.
	// Using this parameter prevents the packet from looping during retransmission in systems with a complex
	// topology of address points. The TTL is initially set by the hardware and software complex that generated the
	// packet. The TTL value is set equal to the maximum allowed number of hardware and software complexes between
	// the sending and receiving hardware and software complexes.
	// The TTL value decreases by one when a packet is transmitted through each hardware and software complex,
	// and the Transport Layer Header checksum is recalculated.
	// When this parameter reaches a value of 0 and when it is found necessary to further route the packet is
	// destroyed and the corresponding code PC_TTLEXPIRED is generated.
	TimeToLive byte `json:"TTL"`
	// HeaderCheckSum (HCS) field is the checksum of the Transport Layer Header (
	// from the "PRV" field to the "HCS" field, not including the "HCS" field).
	// The CRC-8 algorithm is applied to all bytes of the specified sequence to calculate the value of the HCS field.
	HeaderCheckSum byte `json:"HCS"`
	// ServicesFrameData (SFRD) field is a packet type-dependent data structure containing Service Level Support
	// Protocol information.
	ServicesFrameData BinaryData `json:"SFRD"`
	// ServicesFrameDataCheckSum (SFRCS) field is the checksum of the Service Support Protocol level field.
	// To calculate the checksum on data from the SFRD field,
	// the CRC-16 algorithm is used. This field is present only if the SFRD field is present.
	ServicesFrameDataCheckSum uint16 `json:"SFRCS"`
	// ErrorCode contains result of decode package.
	ErrorCode uint8 `json:"-"`
}

// SecretKey is interface for secret key.
type SecretKey interface {
	Decode([]byte) ([]byte, error)
	Encode(data []byte) ([]byte, error)
}

// Options is struct for options of decode/encode operations.
type Options struct {
	Secret SecretKey
}

// Decode parses the set of bytes into the packet structure.
func (p *Packet) Decode(content []byte, opt ...func(*Options)) error {
	options := &Options{}
	for _, o := range opt {
		o(options)
	}

	secretKey := options.Secret

	var (
		err   error
		flags byte
	)
	buf := bytes.NewReader(content)
	if p.ProtocolVersion, err = buf.ReadByte(); err != nil {
		p.ErrorCode = EgtsPcIncHeaderform
		return fmt.Errorf("failed to retrieve the protocol version: %w", err)
	}

	if p.SecurityKeyID, err = buf.ReadByte(); err != nil {
		p.ErrorCode = EgtsPcIncHeaderform
		return fmt.Errorf("failed to get a security key identifier: %w", err)
	}

	// parse flags
	if flags, err = buf.ReadByte(); err != nil {
		p.ErrorCode = EgtsPcIncHeaderform
		return fmt.Errorf("failed to read flags: %w", err)
	}
	flagBits := fmt.Sprintf("%08b", flags)
	p.Prefix = flagBits[:2]         // flags << 7, flags << 6
	p.Route = flagBits[2:3]         // flags << 5
	p.EncryptionAlg = flagBits[3:5] // flags << 4, flags << 3
	p.Compression = flagBits[5:6]   // flags << 2
	p.Priority = flagBits[6:]       // flags << 1, flags << 0

	isEncrypted := p.EncryptionAlg != "00"

	if p.HeaderLength, err = buf.ReadByte(); err != nil {
		p.ErrorCode = EgtsPcIncHeaderform
		return fmt.Errorf("failed to get header length: %w", err)
	}

	if p.HeaderEncoding, err = buf.ReadByte(); err != nil {
		p.ErrorCode = EgtsPcIncHeaderform
		return fmt.Errorf("failed to get encoding method: %w", err)
	}

	tmpIntBuf := make([]byte, 2)
	if _, err = buf.Read(tmpIntBuf); err != nil {
		p.ErrorCode = EgtsPcIncHeaderform
		return fmt.Errorf("failed to get the length of the data section: %w", err)
	}
	p.FrameDataLength = binary.LittleEndian.Uint16(tmpIntBuf)

	if _, err = buf.Read(tmpIntBuf); err != nil {
		p.ErrorCode = EgtsPcIncHeaderform
		return fmt.Errorf("failed to retrieve package identifierа: %w", err)
	}
	p.PacketIdentifier = binary.LittleEndian.Uint16(tmpIntBuf)

	if p.PacketType, err = buf.ReadByte(); err != nil {
		p.ErrorCode = EgtsPcIncHeaderform
		return fmt.Errorf("failed to get package type: %w", err)
	}

	if p.Route == "1" {
		if _, err = buf.Read(tmpIntBuf); err != nil {
			p.ErrorCode = EgtsPcIncHeaderform
			return fmt.Errorf("failed to get the sender's apk address: %w", err)
		}
		p.PeerAddress = binary.LittleEndian.Uint16(tmpIntBuf)

		if _, err = buf.Read(tmpIntBuf); err != nil {
			p.ErrorCode = EgtsPcIncHeaderform
			return fmt.Errorf("failed to get the recipient's apk address: %w", err)
		}
		p.RecipientAddress = binary.LittleEndian.Uint16(tmpIntBuf)

		if p.TimeToLive, err = buf.ReadByte(); err != nil {
			p.ErrorCode = EgtsPcIncHeaderform
			return fmt.Errorf("failed to get TTL of a packet: %w", err)
		}
	}

	if p.HeaderCheckSum, err = buf.ReadByte(); err != nil {
		p.ErrorCode = EgtsPcIncHeaderform
		return fmt.Errorf("failed to get header crc: %w", err)
	}

	if p.HeaderCheckSum != CRC8(content[:p.HeaderLength-1]) {
		p.ErrorCode = EgtsPcHeaderCrcError
		return fmt.Errorf("incorrect checksom of header: %d", p.HeaderCheckSum)
	}

	dataFrameBytes := make([]byte, p.FrameDataLength)
	if _, err = buf.Read(dataFrameBytes); err != nil {
		p.ErrorCode = EgtsPcIncDataform
		return fmt.Errorf("failed to read packet body: %w", err)
	}
	switch p.PacketType {
	case PtAppdataPacket:
		p.ServicesFrameData = &ServiceDataSet{}
	case PtResponsePacket:
		p.ServicesFrameData = &PtResponse{}
	default:
		p.ErrorCode = EgtsPcUnsType
		return fmt.Errorf("unknown package type: %d", p.PacketType)
	}

	if isEncrypted {
		if secretKey == nil {
			p.ErrorCode = EgtsPcDecryptError
			return ErrSecretKey
		}
		dataFrameBytes, err = secretKey.Decode(dataFrameBytes)
		if err != nil {
			p.ErrorCode = EgtsPcDecryptError
			return fmt.Errorf("failed to decrypt packet body: %w", err)
		}
	}

	if err = p.ServicesFrameData.Decode(dataFrameBytes); err != nil {
		p.ErrorCode = EgtsPcDecryptError
		return fmt.Errorf("failed to decode packet body: %w", err)
	}

	crcBytes := make([]byte, 2)
	if _, err = buf.Read(crcBytes); err != nil {
		p.ErrorCode = EgtsPcDecryptError
		return fmt.Errorf("failed to read the CRC16 of the packet: %w", err)
	}
	p.ServicesFrameDataCheckSum = binary.LittleEndian.Uint16(crcBytes)

	if p.ServicesFrameDataCheckSum != CRC16(content[p.HeaderLength:uint16(p.HeaderLength)+p.FrameDataLength]) {
		p.ErrorCode = EgtsPcHeaderCrcError
		return fmt.Errorf("incorrect checksom of body packer: %d", p.ServicesFrameDataCheckSum)
	}
	p.ErrorCode = EgtsPcOk
	return nil
}

// Encode encodes the string into a byte slice.
func (p *Packet) Encode(opt ...func(*Options)) ([]byte, error) {
	var (
		result []byte
		err    error
		flags  uint64
	)

	options := &Options{}
	for _, o := range opt {
		o(options)
	}

	secretKey := options.Secret

	buf := new(bytes.Buffer)

	if err = buf.WriteByte(p.ProtocolVersion); err != nil {
		return result, fmt.Errorf("failed to record protocol version: %w", err)
	}
	if err = buf.WriteByte(p.SecurityKeyID); err != nil {
		return result, fmt.Errorf("failed to write security key ID: %w", err)
	}

	// compile flags
	flagsBits := p.Prefix + p.Route + p.EncryptionAlg + p.Compression + p.Priority
	if flags, err = strconv.ParseUint(flagsBits, 2, 8); err != nil {
		return result, fmt.Errorf("failed to generate a flag byte: %w", err)
	}

	if err = buf.WriteByte(uint8(flags)); err != nil {
		return result, fmt.Errorf("failed to write flags: %w", err)
	}

	if p.HeaderLength == 0 {
		p.HeaderLength = DefaultHeaderLen
		if p.Route == "1" {
			p.HeaderLength += 5
		}
	}

	if err = buf.WriteByte(p.HeaderLength); err != nil {
		return result, fmt.Errorf("failed to write header length: %w", err)
	}

	if err = buf.WriteByte(p.HeaderEncoding); err != nil {
		return result, fmt.Errorf("failed to write encoding method: %w", err)
	}

	var sfrd []byte
	if p.ServicesFrameData != nil { //nolint:nestif
		sfrd, err = p.ServicesFrameData.Encode()
		if err != nil {
			return result, fmt.Errorf("failed to encode services frame data: %w", err)
		}

		if p.EncryptionAlg != "00" {
			if secretKey == nil {
				return result, ErrSecretKey
			}
			sfrd, err = secretKey.Encode(sfrd)
			if err != nil {
				return result, fmt.Errorf("failed to encrypt services frame data: %w", err)
			}
		}
	}
	p.FrameDataLength = uint16(len(sfrd))
	if err = binary.Write(buf, binary.LittleEndian, p.FrameDataLength); err != nil {
		return result, fmt.Errorf("failed to write the length of the data section: %w", err)
	}

	if err = binary.Write(buf, binary.LittleEndian, p.PacketIdentifier); err != nil {
		return result, fmt.Errorf("failed to write packet identifier: %w", err)
	}

	if err = buf.WriteByte(p.PacketType); err != nil {
		return result, fmt.Errorf("failed to write packet identifier: %w", err)
	}

	if p.Route == "1" {
		if err = binary.Write(buf, binary.LittleEndian, p.PeerAddress); err != nil {
			return result, fmt.Errorf("failed to write the sender's apk address: %w", err)
		}

		if err = binary.Write(buf, binary.LittleEndian, p.RecipientAddress); err != nil {
			return result, fmt.Errorf("failed to write the recipient's apk address: %w", err)
		}

		if err = buf.WriteByte(p.TimeToLive); err != nil {
			return result, fmt.Errorf("failed to write TTL packet: %w", err)
		}
	}

	buf.WriteByte(CRC8(buf.Bytes()))

	if p.FrameDataLength > 0 {
		buf.Write(sfrd)
		if err = binary.Write(buf, binary.LittleEndian, CRC16(sfrd)); err != nil {
			return result, fmt.Errorf("failed to write the CRC16 of the packet: %w", err)
		}
	}

	result = buf.Bytes()
	return result, nil
}

// MarshalJSON translates the package into json. Use it to get simple text representation of the package content.
func (p *Packet) MarshalJSON() ([]byte, error) {
	return json.Marshal(p) //nolint:wrapcheck
}

// Response prepares response for incoming packet.
func (p *Packet) Response() ([]byte, error) {
	var (
		resultCode []byte
		err        error
	)

	dataSet := RecordDataSet{}
	serviceType := UndefinedService
	if p.PacketType == PtAppdataPacket && p.ServicesFrameData != nil {
		for _, record := range *p.ServicesFrameData.(*ServiceDataSet) { //nolint:forcetypeassert
			r := record
			data := RecordData{
				SubrecordType:   SrRecordResponseType,
				SubrecordLength: 3,
				SubrecordData: &SrResponse{
					ConfirmedRecordNumber: r.RecordNumber,
					RecordStatus:          EgtsPcOk,
				},
			}
			dataSet = append(dataSet, data)
			serviceType = r.SourceServiceType

			for _, subRec := range r.RecordDataSet {
				switch subRec.SubrecordType {
				case SrTermIdentityType:
					resultCode, err = p.prepareSRResultCode() // ToDo move to sub record level code?
					if err != nil {
						return nil, fmt.Errorf("failed to prepare result code: %w", err)
					}
				case SrAuthInfoType:
					resultCode, err = p.prepareSRResultCode() // ToDo move to sub record level code?
					if err != nil {
						return nil, fmt.Errorf("failed to prepare result code: %w", err)
					}
				}
			}
		}
	}

	respSection := PtResponse{
		ResponsePacketID: p.PacketIdentifier,
		ProcessingResult: p.ErrorCode,
	}

	if len(dataSet) > 0 {
		respSection.SDR = &ServiceDataSet{
			ServiceDataRecord{
				RecordLength:             dataSet.Length(),
				RecordNumber:             nextRecordNumber(),
				SourceServiceOnDevice:    "0",
				RecipientServiceOnDevice: "0",
				Group:                    "1",
				RecordProcessingPriority: "00",
				TimeFieldExists:          "0",
				EventIDFieldExists:       "0",
				ObjectIDFieldExists:      "0",
				SourceServiceType:        serviceType,
				RecipientServiceType:     serviceType,
				RecordDataSet:            dataSet,
			},
		}
	}

	respPkg := Packet{
		ProtocolVersion:   p.ProtocolVersion,
		SecurityKeyID:     p.SecurityKeyID,
		Prefix:            "00",
		Route:             "0",
		EncryptionAlg:     "00",
		Compression:       "0",
		Priority:          "00",
		HeaderLength:      DefaultHeaderLen,
		HeaderEncoding:    0,
		FrameDataLength:   respSection.Length(),
		PacketIdentifier:  nextPacketIdentifier(),
		PacketType:        PtResponsePacket,
		ServicesFrameData: &respSection,
	}

	respBytes, err := respPkg.Encode()
	if err != nil {
		return nil, fmt.Errorf("failed to encode response package: %w", err)
	}
	return append(respBytes, resultCode...), nil
}

// prepareSRResultCode prepares result code (SR_Result_Code) for incoming packet.
func (p *Packet) prepareSRResultCode() ([]byte, error) {
	data := RecordDataSet{
		RecordData{
			SubrecordType:   SrResultCodeType,
			SubrecordLength: 1,
			SubrecordData: &SrResultCode{
				ResultCode: p.ErrorCode,
			},
		},
	}

	sfrd := ServiceDataSet{
		ServiceDataRecord{
			RecordLength:             data.Length(),
			RecordNumber:             nextRecordNumber(),
			SourceServiceOnDevice:    "0",
			RecipientServiceOnDevice: "0",
			Group:                    "1",
			RecordProcessingPriority: "00",
			TimeFieldExists:          "0",
			EventIDFieldExists:       "0",
			ObjectIDFieldExists:      "0", // return object ID?
			SourceServiceType:        AuthService,
			RecipientServiceType:     AuthService,
			RecordDataSet:            data,
		},
	}

	resp := Packet{
		ProtocolVersion:   p.ProtocolVersion,
		SecurityKeyID:     p.SecurityKeyID,
		Prefix:            "00",
		Route:             "0",
		EncryptionAlg:     "00",
		Compression:       "0",
		Priority:          "00",
		HeaderLength:      DefaultHeaderLen,
		HeaderEncoding:    0,
		FrameDataLength:   sfrd.Length(),
		PacketIdentifier:  nextPacketIdentifier(),
		PacketType:        PtResponsePacket,
		ServicesFrameData: &sfrd,
	}

	return resp.Encode()
}
