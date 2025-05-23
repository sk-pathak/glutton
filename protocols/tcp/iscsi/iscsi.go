package iscsi

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type IscsiMsg struct {
	Opcode  uint8
	Flags   uint8
	TaskTag uint32
	Data    uint32
	CID     uint32
	LUN     uint64
}

func ParseISCSIMessage(buffer []byte) (IscsiMsg, IscsiMsg, []byte, error) {
	if len(buffer) < 48 {
		return IscsiMsg{}, IscsiMsg{}, nil, fmt.Errorf("Incomplete iSCSI message")
	}
	msg := IscsiMsg{}
	r := bytes.NewReader(buffer)
	if err := binary.Read(r, binary.BigEndian, &msg); err != nil {
		return IscsiMsg{}, IscsiMsg{}, nil, fmt.Errorf("error reading iSCSI message: %w", err)
	}

	var res IscsiMsg
	switch msg.Opcode {
	case 0x03:
		res = IscsiMsg{
			Opcode:  0x23, // Login response
			Flags:   0x00,
			TaskTag: msg.TaskTag,
			Data:    0,
			CID:     msg.CID,
			LUN:     msg.LUN,
		}
	case 0x01: //Initiator SCSI Command
		res = IscsiMsg{
			Opcode:  0x21, // Target SCSI response
			Flags:   0x00,
			TaskTag: msg.TaskTag,
			Data:    8, //Can vary
			CID:     msg.CID,
			LUN:     msg.LUN,
		}

	case 0x06: // Logout Request
		res = IscsiMsg{
			Opcode:  0x26, // Logout Response
			Flags:   0x00,
			TaskTag: msg.TaskTag,
			Data:    0,
			CID:     msg.CID,
			LUN:     msg.LUN,
		}
	default:
		res = IscsiMsg{
			Opcode:  0x20, // No Operation response
			Flags:   0x00,
			TaskTag: msg.TaskTag,
			Data:    0x00000001, // A generic error data
			CID:     msg.CID,
			LUN:     msg.LUN,
		}
	}

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.BigEndian, res); err != nil {
		return IscsiMsg{}, IscsiMsg{}, nil, fmt.Errorf("failed to write response: %w", err)
	}
	return msg, res, buf.Bytes(), nil
}
