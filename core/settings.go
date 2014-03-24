package core

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
)

type Settings struct {
	io.Writer         `json:"-"`
	Location          string                 `json:"-"`
	DownloadDirectory string                 `json:"download_directory"`
	HosterCredentials map[string]Credentials `json:"hoster_credentials"`
}

func (s *Settings) Save() (err error) {
	m := json.NewEncoder(s.Writer)
	return m.Encode(s)
}

func (s *Settings) String() string {
	b, err := json.Marshal(s)
	if err != nil {
		return ""
	}
	return string(b)

}

type SettingsStore struct {
	io.ReadWriter
}

func homeDir() (s string, err error) {
	usr, err := user.Current()
	if err != nil {
		return
	}
	s = usr.HomeDir
	return
}

// TODO: Make me… better.
func LoadSettings() (s *Settings, err error) {
	s = &Settings{}
	dir, err := homeDir()
	if err != nil {
		return nil, fmt.Errorf("Could not load home dir")
	}
	s.Location = dir + "/.go-download"

	if err != nil {
		return
	}
	store, err := newSettingsStore(s)
	if err != nil {
		return nil, err
	}
	s.Writer = store
	_, err = os.Stat(s.Location)
	if err != nil {
		if os.IsNotExist(err) {
			// No existing Settings
			err = defaultSettings(s)
			if err != nil {
				return
			}
			s.Save()
			err = nil
		} else {
			// Can't stat file
			return
		}
	} else {
		// Existing Settings
		content, err := ioutil.ReadFile(s.Location)
		if err != nil {
			//	Can't load file
			return nil, err
		}
		//	Can load file
		if string(content) == "" {
			err = defaultSettings(s)
			if err != nil {
				return s, err
			}
			s.Save()
			err = nil

		} else {
			err = json.Unmarshal(content, s)
			if err != nil {
				return s, err
			}
		}
	}

	return
}
func defaultSettings(s *Settings) (err error) {
	dir, err := homeDir()
	if err != nil {
		return
	}
	s.DownloadDirectory = dir + "/Downloads"
	s.HosterCredentials = make(map[string]Credentials)
	return
}

func newSettingsStore(s *Settings) (str *SettingsStore, err error) {
	file, err := os.OpenFile(
		s.Location,
		os.O_RDWR|os.O_CREATE,
		0644,
	)
	if err != nil {
		return
	}
	str = &SettingsStore{}
	str.ReadWriter = file
	return
}
