package bot

import (
	"fmt"
	"strings"
	"time"

	"github.com/korsakov-kuzjma/chatan/internal/config"
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
	b.stats.RecordAttendance(user.ID, user.Username, time.Now(), false)

	text := strings.ReplaceAll(b.cfg.Meeting.InvitationText,
		"{время}", b.cfg.Meeting.Schedule[0].Time)
	text = strings.ReplaceAll(text, "{ссылка}", b.cfg.Meeting.StreamLink)

	return c.Send(fmt.Sprintf("Добро пожаловать, %s!\n%s", user.FirstName, text))
}

func (b *Bot) handleAdminHelp(c telebot.Context) error {
	return c.Send(`Команды администратора:
/admin - Эта справка
/set_schedule <день время> - Установить расписание
/set_invitation <текст> - Изменить приглашение
/set_reminder <время> - Установить напоминание
/set_link <url> - Установить ссылку на трансляцию
/stats - Статистика посещений
/participants - Список участников`)
}

func (b *Bot) handleSetSchedule(c telebot.Context) error {
	args := strings.Fields(c.Message().Payload)
	if len(args) < 2 {
		return c.Send("Используйте: /set_schedule <день> <время>\nПример: /set_schedule Tuesday 19:00")
	}

	b.cfg.Meeting.Schedule = []config.MeetingSchedule{
		{Day: args[0], Time: args[1]},
	}

	if err := config.Save(b.cfg); err != nil {
		return c.Send("Ошибка сохранения: " + err.Error())
	}

	b.scheduleMeetings()
	return c.Send("Расписание обновлено!")
}

func (b *Bot) handleSetInvitation(c telebot.Context) error {
	newText := c.Message().Payload
	if newText == "" {
		return c.Send("Укажите текст приглашения")
	}

	b.cfg.Meeting.InvitationText = newText
	if err := config.Save(b.cfg); err != nil {
		return c.Send("Ошибка сохранения: " + err.Error())
	}

	return c.Send("Текст приглашения обновлен!")
}

func (b *Bot) handleSetReminder(c telebot.Context) error {
	duration, err := time.ParseDuration(c.Message().Payload)
	if err != nil {
		return c.Send("Неверный формат. Пример: 1h или 30m")
	}

	b.cfg.Meeting.ReminderBefore = duration
	if err := config.Save(b.cfg); err != nil {
		return c.Send("Ошибка сохранения: " + err.Error())
	}

	b.scheduleMeetings()
	return c.Send(fmt.Sprintf("Напоминание установлено за %s до встречи", duration))
}

func (b *Bot) handleSetStreamLink(c telebot.Context) error {
	newLink := c.Message().Payload
	if newLink == "" {
		return c.Send("Укажите ссылку на трансляцию")
	}

	b.cfg.Meeting.StreamLink = newLink
	if err := config.Save(b.cfg); err != nil {
		return c.Send("Ошибка сохранения: " + err.Error())
	}

	return c.Send("Ссылка на трансляцию обновлена!")
}

func (b *Bot) handleStats(c telebot.Context) error {
	stats := b.stats.GetGeneralStats()
	if len(stats) == 0 {
		return c.Send("Статистика пока отсутствует")
	}

	msg := "📊 Статистика посещений:\n"
	for user, count := range stats {
		msg += fmt.Sprintf("- %s: %d встреч\n", user, count)
	}

	return c.Send(msg)
}

func (b *Bot) handleParticipantsList(c telebot.Context) error {
	participants := b.meeting.GetParticipants()
	if len(participants) == 0 {
		return c.Send("Нет зарегистрированных участников")
	}

	msg := "👥 Список участников:\n"
	for _, p := range participants {
		msg += fmt.Sprintf("- @%s (с %s)\n",
			p.Username, p.JoinedAt.Format("02.01.2006"))
	}

	return c.Send(msg)
}

func (b *Bot) scheduleMeetings() {
	if b.scheduleJob != nil {
		b.scheduleJob.Stop()
	}

	if len(b.cfg.Meeting.Schedule) == 0 {
		return
	}

	// Здесь будет логика планирования встреч
	// Например, отправка напоминаний
}
