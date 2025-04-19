package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Telegram struct {
		Token    string  `yaml:"token"`
		ChatID   int64   `yaml:"chat_id"`
		AdminIDs []int64 `yaml:"admin_ids"`
	} `yaml:"telegram"`
	
	Meeting struct {
		Schedule       []MeetingSchedule `yaml:"schedule"`
		InvitationText string            `yaml:"invitation_text"`
		ReminderBefore time.Duration     `yaml:"reminder_before"`
		StreamLink     string            `yaml:"stream_link"`
	} `yaml:"meeting"`
	
	Modules []string `yaml:"modules"`
}

type MeetingSchedule struct {
	Day  string `yaml:"day"`
	Time string `yaml:"time"`
}

func Load(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func Save(cfg *Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile("config.yaml", data, 0644)
}
