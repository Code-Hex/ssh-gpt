package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/gliderlabs/ssh"
	"github.com/tailscale/hujson"
)

func SFTPHandler(serverName, background string) ssh.Handler {
	return func(s ssh.Session) {
		ctx := s.Context()
		simulator, err := NewSimulator(SFTPMode, serverName, background)
		if err != nil {
			log.Fatal(err)
		}

		for {
			packet, err := recvPacket(s)
			if err != nil {
				ErrLog("recvPacket err: %v", err)
				s.Close()
				return
			}

			ary, err := marshalBytesJSON(packet)
			if err != nil {
				ErrLog("marshalBytesJSON err: %v", err)
				s.Close()
				return
			}

			RecvLog("recv(%s): %s", fxp(packet[4]), ary)

			result, err := simulator.Simulate(ctx, ary)
			if err != nil {
				ErrLog("Simulate err:", err)
				s.Close()
				return
			}

			resp, err := unmarshalJSONBytes(result)
			if err != nil {
				ErrLog("unmarshalJSONBytes err: %v, %q", err, result)
				s.Close()
				return
			}
			SendLog("simulate(%s): %s", fxp(resp[4]), result)

			length := len(resp) - 4 // subtract the uint32(length) from the start
			binary.BigEndian.PutUint32(resp[:4], uint32(length))

			if _, err := s.Write(resp); err != nil {
				ErrLog("unmarshalJSONBytes err: %v", err)
				s.Close()
				return
			}
		}
	}
}

func marshalBytesJSON(b []byte) (string, error) {
	tmp := make([]uint, len(b))
	for i, v := range b {
		tmp[i] = uint(v)
	}
	data, err := json.Marshal(tmp)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func unmarshalJSONBytes(jsonStr string) ([]byte, error) {
	ast, err := hujson.Parse([]byte(jsonStr))
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWCC: %w", err)
	}
	ast.Standardize()

	var data []byte
	if err := json.Unmarshal(ast.Pack(), &data); err != nil {
		return nil, err
	}
	return data, nil
}

func unmarshalUint32(b []byte) (uint32, []byte) {
	v := uint32(b[3]) | uint32(b[2])<<8 | uint32(b[1])<<16 | uint32(b[0])<<24
	return v, b[4:]
}

const maxMsgLength = 256 * 1024

func recvPacket(r io.Reader) ([]byte, error) {
	b := make([]byte, 4)
	if _, err := io.ReadFull(r, b); err != nil {
		return nil, err
	}
	length, _ := unmarshalUint32(b)
	if length > maxMsgLength {
		log.Printf("recv packet %d bytes too long", length)
		return nil, errors.New("long packet")
	}
	if length == 0 {
		log.Printf("recv packet of 0 bytes too short")
		return nil, errors.New("short packet")
	}

	bb := make([]byte, length)
	if _, err := io.ReadFull(r, bb[:length]); err != nil {
		log.Printf("recv packet %d bytes: err %v", length, err)
		return nil, err
	}
	return append(b, bb...), nil
}

func fxp(f byte) string {
	switch f {
	case 1:
		return "SSH_FXP_INIT"
	case 2:
		return "SSH_FXP_VERSION"
	case 3:
		return "SSH_FXP_OPEN"
	case 4:
		return "SSH_FXP_CLOSE"
	case 5:
		return "SSH_FXP_READ"
	case 6:
		return "SSH_FXP_WRITE"
	case 7:
		return "SSH_FXP_LSTAT"
	case 8:
		return "SSH_FXP_FSTAT"
	case 9:
		return "SSH_FXP_SETSTAT"
	case 10:
		return "SSH_FXP_FSETSTAT"
	case 11:
		return "SSH_FXP_OPENDIR"
	case 12:
		return "SSH_FXP_READDIR"
	case 13:
		return "SSH_FXP_REMOVE"
	case 14:
		return "SSH_FXP_MKDIR"
	case 15:
		return "SSH_FXP_RMDIR"
	case 16:
		return "SSH_FXP_REALPATH"
	case 17:
		return "SSH_FXP_STAT"
	case 18:
		return "SSH_FXP_RENAME"
	case 19:
		return "SSH_FXP_READLINK"
	case 20:
		return "SSH_FXP_SYMLINK"
	case 101:
		return "SSH_FXP_STATUS"
	case 102:
		return "SSH_FXP_HANDLE"
	case 103:
		return "SSH_FXP_DATA"
	case 104:
		return "SSH_FXP_NAME"
	case 105:
		return "SSH_FXP_ATTRS"
	case 200:
		return "SSH_FXP_EXTENDED"
	case 201:
		return "SSH_FXP_EXTENDED_REPLY"
	default:
		return "Unknown"
	}
}
