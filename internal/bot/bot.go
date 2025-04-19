package bot

import (
	"log"
	"sync"
	"time"

	"github.com/korsakov-kuzjma/chatan/internal/config"
	"github.com/korsakov-kuzjma/chatan/internal/meeting"
	"github.com/korsakov-kuzjma/chatan/internal/modules/questionnaire"
	"github.com/korsakov-kuzjma/chatan/internal/modules/statistics"
	"gopkg.in/telebot.v3"
)

type Bot struct {
	cfg         *config.Config
	tgBot       *telebot.Bot
	meeting     *meeting.Service
	stats       *statistics.Service
	question    *questionnaire.Service
	mu          sync.RWMutex
	scheduleJob *time.Timer
}

func NewBot(cfg *config.Config) (*Bot, error) {
	pref := telebot.Settings{
		Token:  cfg.Telegram.Token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	tgBot, err := telebot.NewBot(pref)
	if err != nil {
		return nil, err
	}

	b := &Bot{
		cfg:      cfg,
		tgBot:    tgBot,
		meeting:  meeting.NewService(),
		stats:    statistics.NewService(),
		question: questionnaire.NewService(),
	}

	b.registerHandlers()
	b.scheduleMeetings()

	return b, nil
}

func (b *Bot) Start() {
	log.Println("Starting bot...")
	b.tgBot.Start()
}

func (b *Bot) Stop() {
	if b.scheduleJob != nil {
		b.scheduleJob.Stop()
	}
	b.tgBot.Stop()
}
