package bot

import (
	"fmt"
	"strings"
	"time"

	"gopkg.in/telebot.v3"
)

func (b *Bot) registerHandlers() {
	b.tgBot.Handle("/start", b.handleStart)
	b.tgBot.Handle("/help", b.handleHelp)
	b.tgBot.Handle("/meeting", b.handleMeetingInfo)
	b.tgBot.Handle(telebot.OnUserJoined, b.handleUserJoined)
	
	adminGroup := b.tgBot.Group()
	adminGroup.Use(b.requireAdmin)
	adminGroup.Handle("/admin", b.handleAdminHelp)
	adminGroup.Handle("/set_schedule", b.handleSetSchedule)
	adminGroup.Handle("/set_invitation", b.handleSetInvitation)
	adminGroup.Handle("/set_reminder", b.handleSetReminder)
	adminGroup.Handle("/set_link", b.handleSetStreamLink)
	adminGroup.Handle("/stats", b.handleStats)
	adminGroup.Handle("/participants", b.handleParticipantsList)
}

func (b *Bot) requireAdmin(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		for _, adminID := range b.cfg.Telegram.AdminIDs {
			if c.Sender().ID == adminID {
				return next(c)
			}
		}
		return c.Send("🚫 У вас нет прав администратора")
	}
}

func (b *Bot) handleStart(c telebot.Context) error {
	return c.Send("Добро пожаловать в группу Анонимных Наркоманов!")
}

func (b *Bot) handleHelp(c telebot.Context) error {
	return c.Send(`Доступные команды:
/start - Начало работы
/help - Эта справка
/meeting - Информация о встрече`)
}

func (b *Bot) handleMeetingInfo(c telebot.Context) error {
	if len(b.cfg.Meeting.Schedule) == 0 {
		return c.Send("Расписание не установлено")
	}
	
	next := b.cfg.Meeting.Schedule[0]
	msg := fmt.Sprintf("Следующая встреча: %s в %s", next.Day, next.Time)
	
	if b.cfg.Meeting.StreamLink != "" {
		msg += fmt.Sprintf("\nСсылка: %s", b.cfg.Meeting.StreamLink)
	}
	
	return c.Send(msg)
}

func (b *Bot) handleUserJoined(c telebot.Context) error {
	user := c.Sender()
	b.meeting.AddParticipant(user.ID, user.Username)
	
	text := strings.ReplaceAll(b.cfg.Meeting.InvitationText, 
		"{время}", b.cfg.Meeting.Schedule[0].Time)
	text = strings.ReplaceAll(text, "{ссылка}", b.cfg.Meeting.StreamLink)
	
	return c.Send(fmt.Sprintf("Добро пожаловать, %s!\n%s", user.FirstName, text))
}

// ... остальные обработчики ...
