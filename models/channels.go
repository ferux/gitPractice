package models

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

// ReceiverChannels contains registered RTSP channels to cameras
type ReceiverChannels struct {
	EnabledChannels []string               `json:"enabled-channels"`
	Channels        map[string]RTSPChannel `json:"channels"`
}

// RTSPChannel contains link to rtsp channel, auth token and channel id.
type RTSPChannel struct {
	RTSPURL string `json:"rtsp-url"`
	CID     string `json:"cid"`
	Token   string `json:"token"`
}

// ParseJSON parses info to ReceiverChannels struct
func ParseJSON(path string) (rc ReceiverChannels, err error) {
	var file *os.File

	file, err = os.Open(path)
	if err != nil {
		return rc, errors.Wrap(err, "can't open file")
	}
	defer func() {
		if errClose := file.Close(); errClose != nil {
			panic(errClose)
		}

	}()

	err = errors.Wrap(json.NewDecoder(file).Decode(&rc), "can't unmarshal")
	return rc, err
}

// Validate checks
func (rc ReceiverChannels) Validate() (seenSlice []string, unseenSlice []string, err error) {
	// if len(rc.Channels) != len(rc.EnabledChannels) {
	// 	return false, errors.Errorf("channel len: %d enabled len: %d", len(rc.Channels), len(rc.EnabledChannels))
	// }
	seenMap := map[string]bool{}
	unseenSlice = []string{}
	seenSlice = []string{}
	for k := range rc.Channels {
		seenMap[k] = true
	}

	for _, c := range rc.EnabledChannels {
		if _, ok := seenMap[c]; !ok {
			unseenSlice = append(unseenSlice, c)
			continue
		}
		seenSlice = append(seenSlice, c)
	}

	return seenSlice, unseenSlice, nil
}

// DeleteEnabled shows all other channels.
func (rc ReceiverChannels) DeleteEnabled(outpath string) error {
	seen := map[string]bool{}
	for _, c := range rc.EnabledChannels {
		seen[c] = true
	}

	rcDsbld := ReceiverChannels{}
	for k, v := range rc.Channels {
		if !seen[k] {
			rcDsbld.Channels[k] = v
		}
	}

	file, err := os.Create(outpath)
	if err != nil {
		return err
	}
	defer func() {
		if errClose := file.Close(); errClose != nil {
			panic(errClose)
		}
	}()

	return json.NewEncoder(file).Encode(&rcDsbld)
}
