package payloads

import "net/url"

type BuildpackList struct{}

func (d *BuildpackList) SupportedKeys() []string {
	return []string{"order_by"}
}

func (d *BuildpackList) DecodeFromURLValues(values url.Values) error {
	return nil
}
