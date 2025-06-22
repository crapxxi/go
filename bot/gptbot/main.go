package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

type User struct {
	telegram_id int64
	name        string
	level       string
	target_lang string
	progress    string
}
type Word struct {
	Text          string `json:"text"`
	Describe      string `json:"describe"`
	Translation   string `json:"translation"`
	Example       string `json:"example"`
	Example_ru    string `json:"example_ru"`
	Transcription string `json:"transcription"`
}
type Question struct {
	Question       string `json:"Question"`
	VariantA       string `json:"VariantA"`
	VariantB       string `json:"VariantB"`
	VariantC       string `json:"VariantC"`
	VariantD       string `json:"VariantD"`
	CorrectVariant string `json:"CorrectVariant"`
}

func handleError(err error) {
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func main() {

	db, err := wdb()
	if err != nil {
		log.Print("Something witth db")
	}
	handleError(err)
	defer db.Close()

	godotenv.Load()
	botToken := os.Getenv("TGBOTAPI")
	log.Print(botToken)
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Print("Something witth tg")
	}
	handleError(err)
	bot.Debug = true

	log.Printf(bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	issettings := false
	isdialogue := false
	for update := range updates {
		if update.CallbackQuery != nil {
			data := update.CallbackQuery.Data
			user := update.CallbackQuery
			if data[0] == data[1] {
				edit := tgbotapi.NewEditMessageText(user.Message.Chat.ID, user.Message.MessageID, "Правильно!")
				bot.Send(edit)
			} else {
				edit := tgbotapi.NewEditMessageText(user.Message.Chat.ID, user.Message.MessageID, "Не правильно!")
				bot.Send(edit)
			}
		}
		if update.Message != nil {
			if update.Message.Text == "/start" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет, я языковой помощник! Ваш уровень автоматический выбран как начинающий, а язык для обучения выбран как английский, в настройках вы сможете изменить ее!")
				msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton("Учить"),
						tgbotapi.NewKeyboardButton("Тест"),
						tgbotapi.NewKeyboardButton("Диалог"),
						tgbotapi.NewKeyboardButton("Настройки"),
					),
				)
				userdata := User{
					telegram_id: update.Message.Chat.ID,
					name:        update.Message.From.UserName,
					level:       "A1",
					target_lang: "Английский",
					progress:    "",
				}
				addData(userdata, db)
				bot.Send(msg)
			} else if update.Message.Text == "Настройки" {
				issettings = true
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите каков ваш уровень: ")
				msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton("A1"),
						tgbotapi.NewKeyboardButton("A2"),
						tgbotapi.NewKeyboardButton("B1"),
						tgbotapi.NewKeyboardButton("B2"),
						tgbotapi.NewKeyboardButton("C1"),
						tgbotapi.NewKeyboardButton("C2"),
					),
				)
				bot.Send(msg)
			} else if issettings && (update.Message.Text[0] == 'A' || update.Message.Text[0] == 'B' || update.Message.Text[0] == 'C') {
				updateData(update.Message.Chat.ID, update.Message.Text, "level", db)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите какой язык хотите изучать: ")
				msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton("Английский"),
						tgbotapi.NewKeyboardButton("Немецкий"),
						tgbotapi.NewKeyboardButton("Китайский"),
					),
				)
				bot.Send(msg)
			} else if issettings && (update.Message.Text == "Английский" || update.Message.Text == "Немецкий" || update.Message.Text == "Китайский") {
				updateData(update.Message.Chat.ID, update.Message.Text, "", db)
				err := deleteprogress(update.Message.Chat.ID, db)
				handleError(err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Отлично! Сохранено!")
				msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton("Учить"),
						tgbotapi.NewKeyboardButton("Тест"),
						tgbotapi.NewKeyboardButton("Диалог"),
						tgbotapi.NewKeyboardButton("Настройки"),
					),
				)
				bot.Send(msg)
			} else if update.Message.Text == "Учить" {
				wordsprogress, err := getProgress(db, update.Message.Chat.ID)
				handleError(err)
				level, err := getLevel(db, update.Message.Chat.ID)
				handleError(err)
				lang, err := getLanguage(db, update.Message.Chat.ID)
				handleError(err)
				prompt := fmt.Sprintf(`Дай 5 слов на какую то тему уровня "%s" в JSON-формате. На "%s" языке.Пример для использования так же должен быть на этом языке! Добавь произношения транскрипцию целого предложение в поле transcription и подчеркни произношение этого слово в транскрипции. Слова которые были изучены "%s" ты должен их опять не давать. У каждого слова должны быть поля: "text","describe", "translation", "example","example_ru", "transcription". Без лишнего текста.`, level, lang, wordsprogress)
				response, err := gpt(prompt)
				handleError(err)
				re := regexp.MustCompile(`\[[\s\S]*\]`)
				cleaned := re.FindString(response)
				var words []Word
				err = json.Unmarshal([]byte(cleaned), &words)
				handleError(err)
				i := 1
				for _, word := range words {
					err = updateData(update.Message.Chat.ID, word.Text, "progress", db)
					handleError(err)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, strconv.Itoa(i)+". Слово: \b"+word.Text+"\nОбьяснение слова: "+word.Describe+"\nПример для использования и перевод на русском:\n\b"+word.Example+"\n\b"+word.Example_ru+"\nЗначение слова:\b"+word.Translation+"\nТранскрипция:\b"+word.Transcription)
					msg.ParseMode = "HTML"
					bot.Send(msg)
					i++
				}
			} else if update.Message.Text == "Тест" {
				wordsprogress, err := getProgress(db, update.Message.Chat.ID)
				handleError(err)
				level, err := getLevel(db, update.Message.Chat.ID)
				handleError(err)
				lang, err := getLanguage(db, update.Message.Chat.ID)
				handleError(err)
				prompt := fmt.Sprintf(`Дай тест из 10 вопросов в JSON-формате. Уровень пользователя "%s". Слова по которым берется тест:"%s", тест должен быть на языке "%s".В поле CorrectVariant должен указываться правильный вариант "A,B,C,D" У каждого вопроса должны быть поля: "question", "VariantA", "VariantB","VariantC", "VariantD", "CorrectVariant". Без лишнего текста.`, level, wordsprogress, lang)
				response, err := gpt(prompt)
				handleError(err)
				re := regexp.MustCompile(`\[[\s\S]*\]`)
				cleaned := re.FindString(response)
				var tests []Question
				err = json.Unmarshal([]byte(cleaned), &tests)
				handleError(err)
				for _, test := range tests {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вопрос: "+test.Question+"\nA."+test.VariantA+"\nB."+test.VariantB+"\nC."+test.VariantC+"\nD."+test.VariantD)
					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("A", "A"+test.CorrectVariant),
							tgbotapi.NewInlineKeyboardButtonData("B", "B"+test.CorrectVariant),
							tgbotapi.NewInlineKeyboardButtonData("C", "C"+test.CorrectVariant),
							tgbotapi.NewInlineKeyboardButtonData("D", "D"+test.CorrectVariant),
						),
					)
					bot.Send(msg)
				}
			} else if update.Message.Text == "Диалог" {
				isdialogue = true
				level, err := getLevel(db, update.Message.Chat.ID)
				handleError(err)
				lang, err := getLanguage(db, update.Message.Chat.ID)
				handleError(err)
				prompt := fmt.Sprintf(`Ты сейчас будешь общаться с пользователем и вести диалог на "%s" языке. У пользователя уровень "%s". Сейчас без лишнего текста начни диалог`, lang, level)
				response, err := gpt(prompt)
				handleError(err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
				bot.Send(msg)
			} else if isdialogue {
				if update.Message.Text == "Остановить диалог" {
					isdialogue = false
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Диалог остановлен!")
					msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
						tgbotapi.NewKeyboardButtonRow(
							tgbotapi.NewKeyboardButton("Учить"),
							tgbotapi.NewKeyboardButton("Тест"),
							tgbotapi.NewKeyboardButton("Диалог"),
							tgbotapi.NewKeyboardButton("Настройки"),
						),
					)
					bot.Send(msg)
				} else {
					lang, err := getLanguage(db, update.Message.Chat.ID)
					handleError(err)
					response, err := gpt("Ты ведешь диалог с пользователем, и ты должен не отвечать на вопросы типо дай код или дай формулу. Ты должен вести себя как человек и разговаривать на " + lang + " языке. На вопросы не касающиеся обучение языка лучше отвечать типо я незнаю. Без лишних слов веди диалог. Вот ответ пользователя: " + update.Message.Text)
					handleError(err)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
					msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
						tgbotapi.NewKeyboardButtonRow(
							tgbotapi.NewKeyboardButton("Остановить диалог"),
						),
					)
					bot.Send(msg)
				}
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Я вас не понял!")
				msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton("Учить"),
						tgbotapi.NewKeyboardButton("Тест"),
						tgbotapi.NewKeyboardButton("Диалог"),
						tgbotapi.NewKeyboardButton("Настройки"),
					),
				)
				bot.Send(msg)
			}
		}
	}
}
