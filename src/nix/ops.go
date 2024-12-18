// SPDX-License-Identifier: MIT

package nix

import "errors"

func (c *Connection) QueryPathInfo(path string) (*PathInfo, error) {
	if err := c.writeUint64(uint64(QueryPathInfo)); err != nil {
		return nil, err
	}
	if err := c.writeString(path); err != nil {
		return nil, err
	}
	if err := c.processStderr(); err != nil {
		return nil, err
	}
	if c.version&0xff >= 17 {
		valid, err := c.readUint64()
		if err != nil {
			return nil, err
		}
		if valid != 1 {
			return nil, errors.New("invalid query path info")
		}
	}

	var err error
	info := &PathInfo{}

	if info.Deriver, err = c.readString(); err != nil {
		return nil, err
	}

	if info.NarHash, err = c.readString(); err != nil {
		return nil, err
	}

	if info.References, err = c.readStrings(); err != nil {
		return nil, err
	}

	if info.RegistrationTime, err = c.readUint64(); err != nil {
		return nil, err
	}

	if info.NarSize, err = c.readUint64(); err != nil {
		return nil, err
	}

	if c.version&0xff >= 16 {
		if info.Ultimate, err = c.readBool(); err != nil {
			return nil, err
		}
		if info.Sigs, err = c.readStrings(); err != nil {
			return nil, err
		}
		if info.Ca, err = c.readString(); err != nil {
			return nil, err
		}
	}

	return info, nil
}

func (c *Connection) GetAllReferences(path string) ([]string, error) {
	refs := make(map[string]struct{})
	queue := []string{path}

	for len(queue) > 0 {
		path, queue = queue[len(queue)-1], queue[:len(queue)-1]
		if _, ok := refs[path]; ok {
			continue
		}
		refs[path] = struct{}{}
		info, err := c.QueryPathInfo(path)
		if err != nil {
			return nil, err
		}
		for _, ref := range info.References {
			queue = append(queue, ref)
		}
	}

	var result []string
	for ref := range refs {
		result = append(result, ref)
	}

	return result, nil
}
