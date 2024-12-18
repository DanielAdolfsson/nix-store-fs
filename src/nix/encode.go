package nix

import (
	"encoding/binary"
	"io"
)

func (c *Connection) readBool() (bool, error) {
	var buf [8]byte
	if _, err := io.ReadFull(c.conn, buf[:]); err != nil {
		return false, err
	}
	return binary.LittleEndian.Uint64(buf[:]) == 1, nil
}

func (c *Connection) readUint64() (uint64, error) {
	var buf [8]byte
	if _, err := io.ReadFull(c.conn, buf[:]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(buf[:]), nil
}

func (c *Connection) readString() (string, error) {
	length, err := c.readUint64()
	if err != nil {
		return "", err
	}
	// Length rounded up to nearest multiple of 8.
	buf := make([]byte, (length+7)&^7)
	if _, err := io.ReadFull(c.conn, buf); err != nil {
		return "", err
	}
	return string(buf[:length]), nil
}

func (c *Connection) readStrings() ([]string, error) {
	count, err := c.readUint64()
	if err != nil {
		return nil, err
	}
	result := make([]string, count)
	for i := uint64(0); i < count; i++ {
		str, err := c.readString()
		if err != nil {
			return nil, err
		}
		result[i] = str
	}
	return result, nil
}

func (c *Connection) writeUint64(value uint64) error {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], uint64(value))
	_, err := c.conn.Write(buf[:])
	return err
}

func (c *Connection) writeString(value string) error {
	length := len(value)
	// 8 + (length rounded up to nearest multiple of 8).
	buf := make([]byte, 8+(length+7)&^7)
	copy(buf[8:], value)
	binary.LittleEndian.PutUint64(buf[:], uint64(length))
	_, err := c.conn.Write(buf[:])
	return err
}

func (c *Connection) writeStrings(value []string) error {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], uint64(len(value)))
	if _, err := c.conn.Write(buf[:]); err != nil {
		return err
	}
	for _, str := range value {
		if err := c.writeString(str); err != nil {
			return err
		}
	}
	return nil
}
