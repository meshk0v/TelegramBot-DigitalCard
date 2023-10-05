package main

import (
	"log"
	"time"

	"github.com/fogleman/gg" // Импорт библиотеки для работы с изображениями
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	stepMap = make(map[int64]*StepData) // Карта для хранения данных о текущем шаге для каждого пользователя
)

type StepData struct {
	Step     int
	UserData UserData
}

type UserData struct {
	ChoiceDisegn   string
	FirstName    string
	LastName string
	PhoneNumber    string
	Login    string
	Photo    string
	City    string
	UTP    string
}

func main() {
	// Инициализация бота с помощью токена
    botToken := ""

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}

	// Настройка логгирования
	bot.Debug = true

	// Настройка обработчика команд и сообщений
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Обработчик сообщений от пользователя
	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := int64(update.Message.Chat.ID)

		if update.Message.IsCommand() && update.Message.Command() == "start" {
			// Начнем процесс создания визитки для пользователя
			startCreatingBusinessCard(bot, update.Message)
		} else {
			// Проверяем текущий шаг пользователя и маршрутизируем ввод
			stepData, ok := stepMap[chatID]
			if !ok {
				// Если объект StepData не найден, то пользователь начинает процесс заново
				startCreatingBusinessCard(bot, update.Message)
				continue
			}

			switch stepData.Step {
			case 1:
				handleFirstNameInput(bot, update.Message, stepData)
			case 2:
				handleLastNameInput(bot, update.Message, stepData)
			case 3:
				handlePhoneNumberInput(bot, update.Message, stepData)
			case 4:
				handleLoginInput(bot, update.Message, stepData)
			case 5:
				handlePhotoInput(bot, update.Message, stepData)
			case 6:
				handleCityInput(bot, update.Message, stepData)
			case 7:
				handleUTPInput(bot, update.Message, stepData)
			case 8:
				handleUsernameInput(bot, update.Message, stepData)
			default:
				// По умолчанию отправляем сообщение об ошибке
				reply := tgbotapi.NewMessage(chatID, "Неожиданный ввод. Пожалуйста, следуйте инструкциям.")
				bot.Send(reply)
			}
		}
	}
}

func startCreatingBusinessCard(bot *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	chatID := msg.Chat.ID

	// Создаем новый объект StepData для пользователя
	stepData := &StepData{
		Step:     1,
		UserData: UserData{},
	}

	// Сохраняем StepData в карту stepMap по chatID
	stepMap[chatID] = stepData

	// Отправляем сообщение с инструкциями и клавиатурой для выбора дизайна
	reply := tgbotapi.NewMessage(chatID, "Давайте начнем создание цифровой визитки.\n\nВыберите дизайн")
	reply.ReplyMarkup = createKeyboard() // Создаем клавиатуру с кнопками "1", "2", "3"
	bot.Send(reply)
	sendPhotos(bot, msg.Chat.ID)
}

func sendPhotos(bot *tgbotapi.BotAPI, chatID int64) {
	// Здесь можно отправить 3 фотографии пользователю
	photo1 := tgbotapi.NewPhotoUpload(chatID, "1.jpg")
	photo2 := tgbotapi.NewPhotoUpload(chatID, "2.jpg")
	photo3 := tgbotapi.NewPhotoUpload(chatID, "3.jpg")

	// Добавляем подписи к фотографиям
	photo1.Caption = "Дизайн 1"
	photo2.Caption = "Дизайн 2"
	photo3.Caption = "Дизайн 3"

	// Отправляем фотографии
	bot.Send(photo1)
	time.Sleep(1 * time.Second) // Задержка между отправкой фотографий
	bot.Send(photo2)
	time.Sleep(1 * time.Second)
	bot.Send(photo3)
}

func createKeyboard() tgbotapi.ReplyKeyboardMarkup {
	// Создаем клавиатуру с кнопками "1", "2", "3"
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("1"),
			tgbotapi.NewKeyboardButton("2"),
			tgbotapi.NewKeyboardButton("3"),
		),
	)
	return keyboard
}

func handleFirstNameInput(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, stepData *StepData) {
	chatID := msg.Chat.ID

	// Сохраняем введенное имя в объект UserData
	stepData.UserData.ChoiceDisegn = msg.Text

	// Переходим к следующему шагу
	stepData.Step++

	// Отправляем инструкции для следующего шага
	reply := tgbotapi.NewMessage(chatID, "Напишите ваше имя")
	bot.Send(reply)

}

func handleLastNameInput(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, stepData *StepData) {
	chatID := msg.Chat.ID

	// Сохраняем введенную фамилию в объект UserData
	stepData.UserData.FirstName = msg.Text

	// Переходим к следующему шагу
	stepData.Step++

	// Отправляем инструкции для следующего шага
	reply := tgbotapi.NewMessage(chatID, "Теперь напишите вашу фамилию")
	bot.Send(reply)
}

func handlePhoneNumberInput(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, stepData *StepData) {
	chatID := msg.Chat.ID

	// Сохраняем введенный номер телефона в объект UserData
	stepData.UserData.LastName = msg.Text

	// Переходим к следующему шагу
	stepData.Step++

	// Отправляем инструкции для следующего шага
	reply := tgbotapi.NewMessage(chatID, "Теперь напишите ваш номер телефона")
	bot.Send(reply)
}

func handleLoginInput(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, stepData *StepData) {
	chatID := msg.Chat.ID

	// Сохраняем введенный номер телефона в объект UserData
	stepData.UserData.PhoneNumber = msg.Text

	// Переходим к следующему шагу
	stepData.Step++

	// Отправляем инструкции для следующего шага
	reply := tgbotapi.NewMessage(chatID, "Теперь напишите ваш логин в телеграм")
	bot.Send(reply)
}

func handlePhotoInput(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, stepData *StepData) {
	chatID := msg.Chat.ID

	// Сохраняем введенный номер телефона в объект UserData
	stepData.UserData.Login = msg.Text

	// Переходим к следующему шагу
	stepData.Step++

	// Отправляем инструкции для следующего шага
	reply := tgbotapi.NewMessage(chatID, "Отправьте фото, чтобы пропустить введите -")
	bot.Send(reply)
}

func handleCityInput(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, stepData *StepData) {
	chatID := msg.Chat.ID

	// Сохраняем введенный номер телефона в объект UserData
	stepData.UserData.Photo = msg.Text

	// Переходим к следующему шагу
	stepData.Step++

	// Отправляем инструкции для следующего шага
	reply := tgbotapi.NewMessage(chatID, "Отправьте город, чтобы пропустить введите -")
	bot.Send(reply)
}

func handleUTPInput(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, stepData *StepData) {
	chatID := msg.Chat.ID

	// Сохраняем введенный номер телефона в объект UserData
	stepData.UserData.City = msg.Text

	// Переходим к следующему шагу
	stepData.Step++

	// Отправляем инструкции для следующего шага
	reply := tgbotapi.NewMessage(chatID, "Напишите УТП")
	bot.Send(reply)
}

func handleUsernameInput(bot *tgbotapi.BotAPI, msg *tgbotapi.Message, stepData *StepData) {
	chatID := msg.Chat.ID

	// Сохраняем введенный логин в объект UserData
	stepData.UserData.UTP = msg.Text

	// Пользователь завершил процесс создания визитки, вы можете сделать здесь что-то еще

	// Очищаем StepData из карты stepMap
	delete(stepMap, chatID)

	// Наложение данных пользователя на фото
	mergedPhotoPath := mergeDataWithPhoto("1.jpg", stepData.UserData)
	if mergedPhotoPath == "" {
		reply := tgbotapi.NewMessage(chatID, "Произошла ошибка при наложении данных на фото.")
		bot.Send(reply)
		return
	}

	// Отправляем фото с данными обратно пользователю
	photo := tgbotapi.NewPhotoUpload(chatID, mergedPhotoPath)
	bot.Send(photo)

	// Отправляем сообщение о завершении процесса
	reply := tgbotapi.NewMessage(chatID, "Создание визитки завершено. Ваши данные:\n"+
		"Дизайн: "+stepData.UserData.ChoiceDisegn+"\n"+
		"Имя: "+stepData.UserData.FirstName+"\n"+
		"Фамилия: "+stepData.UserData.LastName+"\n"+
		"Номер телефона: "+stepData.UserData.PhoneNumber+"\n"+
		"Логин: "+stepData.UserData.Login+"\n"+
		"Город: "+stepData.UserData.City+"\n"+
		"УТП: "+stepData.UserData.UTP)
	bot.Send(reply)
}
// Функция для наложения данных пользователя на фото дизайна 1
func mergeDataWithPhoto(designFileName string, userData UserData) string {
	// Загружаем изображение дизайна
	im, err := gg.LoadImage(designFileName)
	if err != nil {
		log.Println("Ошибка при загрузке изображения дизайна:", err)
		return ""
	}

	// Создаем контекст для рисования на изображении
	dc := gg.NewContextForImage(im)

	// Устанавливаем цвет текста и размер шрифта
	dc.SetRGB(0, 0, 0) // черный цвет
	err = dc.LoadFontFace("/home/romeo/Desktop/AllMyProject/chatgpt_arsen/arial.ttf", 150) // Укажите путь к файлу шрифта и размер
	if err != nil {
		log.Println("Ошибка при загрузке шрифта:", err)
		return ""
	}

	// Налагаем данные пользователя на фото
	// Разместить имя и фамилию в одной строке
	dc.DrawStringAnchored(userData.FirstName+" "+userData.LastName, 300, 500, 0.1, 0.1)
	dc.DrawStringAnchored(userData.PhoneNumber, 300, 750, 0.1, 0.1)
	dc.DrawStringAnchored(userData.Login, 300, 950, 0.1, 0.1)
	dc.DrawStringAnchored(userData.UTP, 300, 1150, 0.1, 0.1)

	// Сохраняем результат в новый файл
	resultFileName := "result_design1.jpg" // Укажите имя файла для сохранения
	err = dc.SavePNG(resultFileName)
	if err != nil {
		log.Println("Ошибка при сохранении изображения:", err)
		return ""
	}

	return resultFileName
}
