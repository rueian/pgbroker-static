package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

func NewSettings(path string) (*Settings, error) {
	s := &Settings{path: path}
	if err := s.load(path); err != nil {
		return nil, err
	}

	return s, nil
}

type YML struct {
	Databases map[string]struct {
		Address string `yaml:"address"`
	} `yaml:"databases"`
}

type Settings struct {
	mu   sync.Mutex
	path string
	yml  YML
}

func (s *Settings) GetAddress(database string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if m, ok := s.yml.Databases[database]; ok {
		return m.Address
	}
	return ""
}

func (s *Settings) Watch() {
	for {
		time.Sleep(3 * time.Second)
		if err := s.load(s.path); err != nil {
			log.Printf("Fail to reload config from %s\n", s.path)
		}
	}
}

func (s *Settings) load(path string) error {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	yml := YML{}
	if err = yaml.Unmarshal(bs, &yml); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.yml = yml
	return nil
}
