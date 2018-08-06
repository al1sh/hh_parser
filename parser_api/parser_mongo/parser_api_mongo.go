package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	// start := time.Now()
	db := PrepareDB()

	// user := "user2"
	// cities := []int64{1, 2, 3}

	// InitUser(db, user)
	// // SetString(db, user, "search", "python")
	// // SetString(db, user, "experience", "0")
	// // SetArray(db, user, "cities", cities)

	// SetElement(db, user, "search", "go")
	// SetElement(db, user, "experience", int64(0))
	// SetArray(db, user, "cities", cities)

	// RetrieveUser(db)

	// elapsed := time.Since(start)
	// fmt.Printf("Search took %s", elapsed)

	bot, err := tgbotapi.NewBotAPI("664176668:AAEGlhy3pLIJEhDO6NAjfiurw8HAb04lQ3g")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		input := strings.Fields(update.Message.Text)
		var reply tgbotapi.MessageConfig
		var message string

		user := strconv.Itoa(int(update.Message.Chat.ID))
		help := `Данный бот предназначен для получения уведомлений о новых вакансиях на сайте hh.ru 
			Чтобы включить уведомления введите необходимые данные с помощью команд
			/setJob *вакансия* 
			/setExperience *(0,1, или 3)*
			/setCities *код городов или регионов через запятую*

			После заполнения параметров используйте /on и /off для включения/отключения уведомлений
			`

		switch input[0] {

		case "/start":
			result := InitUser(db, user)
			if result {
				fmt.Println("User Initiated")
				message = help
			} else {
				message = "Для запуска введите /on"
			}

		case "/setJob":
			search := strings.Join(input[1:], "+")
			SetElement(db, user, SearchField, search)
			message = "Вакансия для поиска успешно обновлена"

		case "/setExperience":

			status := true

			if exp, err := strconv.Atoi(input[1]); err == nil {
				if exp >= 0 && exp < 4 {
					SetElement(db, user, ExperienceField, int64(exp))

				} else {
					status = false
				}

			} else {
				status = false
			}

			if status {
				message = "Опыт вакансий установлен на " + input[1]
			} else {
				message = `Ошибка. Опыт может быть задан только числами 0, 1 и 3, что соответствует минимальному количеству лет 
					Пример: /setExperience 1`
			}

		case "/setCities":
			cities := strings.Split(input[1], ",")
			var citiesInt []int64
			status := true

			if len(input) == 2 {
				for _, i := range cities {
					if num, err := strconv.Atoi(i); err == nil {
						if num > 0 {
							citiesInt = append(citiesInt, int64(num))
						} else {
							status = false
						}
					} else {
						status = false
					}
				}

			} else {
				status = false
			}

			if status {
				SetArray(db, user, CitiesField, citiesInt)
				message = "Города успешно выбраны"
			} else {
				message = `Ошибка. Выберите коды регионов через запятую, без пробелов
						Пример: /setCities 1,2,3 или /setCities 1`
			}

		case "/on":
			status := true

			userinfo, err := GetUser(db, user)
			if err == nil {
				if userinfo.Search != "" {
					SetElement(db, user, ActiveField, true)
				}
			} else {
				status = false
			}

			if status {
				message = fmt.Sprintf(`Вы подписались на уведомления при добавлении вакансий %v с опытом %v в указанных городах `,
					userinfo.Search, userinfo.Experience)
			} else {
				message = "Ошибка чтения данных пользователя. "
			}

		case "/help":
			message = help

		case "/off":
			SetElement(db, user, ActiveField, false)
			message = "Уведомления отключены"

		default:
			message = "Ваша команда не распознана. Для инструкций по использованию введите /help"

		}

		// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text+" boiii")
		// msg.ReplyToMessageID = update.Message.MessageID

		// reply.ReplyToMessageID = update.Message.MessageID

		reply = tgbotapi.NewMessage(update.Message.Chat.ID, message)
		bot.Send(reply)
	}
}
