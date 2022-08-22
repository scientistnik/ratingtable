package telegram

import (
	"context"
	"fmt"
	"ratingtable/internal/app"
	"ratingtable/internal/app/domain"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Telegram struct {
	token string
}

func NewTelegram(token string) *Telegram {
	return &Telegram{token: token}
}

func (t Telegram) Launch(ctx context.Context, app *app.App) error {
	bot, err := tgbotapi.NewBotAPI(t.token)
	if err != nil {
		return fmt.Errorf("telegram error, %w", err)
	}

	bot.Debug = true

	//fmt.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for {
		select {
		case update, ok := <-updates:
			if !ok {
				return nil
			}

			if update.Message == nil {
				continue
			}

			// if !update.Message.IsCommand() {
			// 	continue
			// }

			//fmt.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			//msg.Text = "Hello"

			if update.Message.IsCommand() {
				switch update.Message.Command() {
				case "start":
					msg.Text = "Hello all"
				case "help":
					msg.Text = "I understand /start and /status."
				case "status":
					msg.Text = "I'm ok."
				case "party":
					msg.Text = "saved"
					fmt.Printf("[PARTY] %#v", update.Message)
					//update.Message.Text // "party"
					//app.TeamGet()
					app.AddParty("chess", []domain.TeamPoints{})
				default:
					msg.Text = "I don't know that command"
				}

			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
				msg.ReplyToMessageID = update.Message.MessageID
				//bot.Send(msg)
			}

			if len(msg.Text) > 0 {
				if _, err := bot.Send(msg); err != nil {
					fmt.Printf("telegram error: %#v\n", err)
				}
			}

		case <-ctx.Done():
			return nil
		}
	}
}
