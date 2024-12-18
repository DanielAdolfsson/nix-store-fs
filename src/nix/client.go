package nix

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
)

type OpCode int

const (
	ClientVersion uint64 = (1 << 8) | 38
	ClientMagic   uint64 = 0x6e697863
	ServerMagic   uint64 = 0x6478696f

	QueryPathInfo OpCode = 26

	StderrNext          uint64 = 0x6f6c6d67
	StderrRead          uint64 = 0x64617461 // data needed from source
	StderrWrite         uint64 = 0x64617416 // data for sink
	StderrLast          uint64 = 0x616c7473
	StderrError         uint64 = 0x63787470
	StderrStartActivity uint64 = 0x53545254
	StderrStopActivity  uint64 = 0x53544f50
	StderrResult        uint64 = 0x52534c54
)

type Connection struct {
	conn    net.Conn
	version int
}

type PathInfo struct {
	Deriver          string
	NarSize          uint64
	NarHash          string
	References       []string
	RegistrationTime uint64
	Ultimate         bool
	Sigs             []string
	Ca               string
}

func Connect(path string) (*Connection, error) {
	conn, err := net.Dial("unix", path)
	if err != nil {
		return nil, err
	}

	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()

	var buf [32]byte
	binary.LittleEndian.PutUint64(buf[:], ClientMagic)
	binary.LittleEndian.PutUint64(buf[8:], ClientVersion)
	if _, err := conn.Write(buf[:16]); err != nil {
		return nil, err
	}

	if _, err := io.ReadFull(conn, buf[:16]); err != nil {
		return nil, err
	}

	magic := binary.LittleEndian.Uint64(buf[:])
	if magic != ServerMagic {
		return nil, fmt.Errorf("invalid nix magic number: %d", magic)
	}

	version := binary.LittleEndian.Uint64(buf[8:])

	version = min(version, ClientVersion)
	connection := &Connection{conn: conn, version: int(version)}

	if version&0xff >= 38 {
		if err := connection.writeStrings([]string{}); err != nil {
			return nil, err
		}
		features, err := connection.readStrings()
		if err != nil {
			return nil, err
		}
		_ = features
	}

	if connection.version&0xff >= 14 {
		// Obsolete CPU affinity.
		if err := connection.writeUint64(0); err != nil {
			return nil, err
		}
	}

	if connection.version&0xff >= 11 {
		// Obsolete reserveSpace.
		if err := connection.writeUint64(0); err != nil {
			return nil, err
		}
	}

	if connection.version&0xff >= 33 {
		daemonVersion, err := connection.readString()
		if err != nil {
			return nil, err
		}
		_ = daemonVersion
	}

	if connection.version&0xff >= 35 {
		trusted, err := connection.readUint64()
		if err != nil {
			return nil, err
		}
		_ = trusted
		// 0 = undefined
		// 1 = trusted
		// 2 = not trusted
	}

	if err := connection.processStderr(); err != nil {
		return nil, err
	}

	// Clear this variable so we're not closing it during cleanup.
	conn = nil

	return connection, nil
}

func (c *Connection) processStderr() error {
	for {
		msg, err := c.readUint64()
		if err != nil {
			return err
		}

		// TODO: Handle all Stderr* messages.

		switch msg {
		case StderrLast:
			return nil

		default:
			return errors.New("unexpected stderr message")
		}
	}
}

func (c *Connection) Close() error {
	if c.conn == nil {
		return nil
	}
	err := c.conn.Close()
	c.conn = nil
	return err
}
