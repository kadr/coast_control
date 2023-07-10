package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/cost_control/config"
	"github.com/cost_control/internal/handlers/telegram/product"
	productRepos "github.com/cost_control/internal/repository/product"
	userRepos "github.com/cost_control/internal/repository/user"
	product2 "github.com/cost_control/internal/service/product"
	"github.com/cost_control/internal/service/user"
	"github.com/cost_control/pkg/jwt"
	"github.com/cost_control/pkg/logger"
	"github.com/cost_control/pkg/password_hasher"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"net/mail"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	defaultTimeout  = 60
	inputDateFormat = "02.01.2006 15:04"
	//telegramApiHost = "https://api.telegram.org"
	saveProduct = "1"
)

var (
	begin       = time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Local)
	beginMonth  = begin.Format("02.01.2006")
	now         = time.Now().Format("02.01.2006")
	nowCmd      = time.Now().Format("02.01.2006T15:04:05 +0400 +04")
	commandsKey = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Добавить новый продукт", "add_product"),
			tgbotapi.NewInlineKeyboardButtonData("Получить продукты", fmt.Sprintf("get_products %s %s",
				beginMonth, nowCmd)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отчет", fmt.Sprintf("get_report %s %s", beginMonth, now)),
		),
	)
)

const productCollection = "product"
const userCollection = "user"

type BotHandler struct {
	productHandler product.ProductBotHandler
	userService    user.UserService
	bot            tgbotapi.BotAPI
	log            logger.ILogger
	cfg            *config.Config
	jwtManager     *jwt.Token
	hasherManager  *password_hasher.PasswordHasher
}

type InputData struct {
	ChatId  int64
	User    *tgbotapi.User
	Text    string
	Command string
}

func New(
	token string,
	db *mongo.Database,
	log logger.ILogger,
	cfg *config.Config,
) (*BotHandler, error) {
	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		return nil, err
	}
	productRepo := productRepos.New(db.Collection(productCollection))
	userRepo := userRepos.New(db.Collection(userCollection))
	return &BotHandler{
		productHandler: *product.New(product2.New(productRepo)),
		bot:            *bot,
		log:            log,
		cfg:            cfg,
		jwtManager:     &jwt.Token{},
		userService:    *user.New(userRepo),
		hasherManager:  password_hasher.New(),
	}, err
}

func (b BotHandler) Start(wg *sync.WaitGroup, updateTimeout *int, offset *int) error {
	defer wg.Done()
	b.log.Info("Start telegram bot Api")
	inputData := InputData{}
	config := initConfig(updateTimeout, offset)

	b.getUpdates(config, &inputData)
	return nil
}

func (b BotHandler) AddProduct(productDto product.CreateProductDTO) (string, error) {
	_, err := b.productHandler.Create(productDto)
	if err != nil {
		return "", err
	}

	return "Продукт добавлен", nil
}

func (b BotHandler) GetProducts(data *InputData) ([]product.GetProductDTO, error) {
	products, err := b.productHandler.Get(data.Text)
	if err != nil {
		b.log.Info(err)
		return nil, err
	}
	var getProducts []product.GetProductDTO
	for _, _product := range products {
		getProducts = append(getProducts, product.GetProductDTO{
			Name:        _product.Name,
			Price:       _product.Price,
			Description: _product.Description,
			BuyAt:       _product.BuyAt,
			User:        _product.User,
		})
	}

	return getProducts, nil

}

func (b BotHandler) GetReport(data *InputData) (map[string]float32, error) {
	// Реализовать группировку по пользователям
	report, err := b.productHandler.Report(data.Text)
	if err != nil {
		return nil, err
	}
	return report, nil
}

func (b BotHandler) getUpdates(config tgbotapi.UpdateConfig, inputData *InputData) {
	updates, err := b.bot.GetUpdatesChan(config)
	if err != nil {
		b.log.Fatal(err)
	}
	productDTO := &product.CreateProductDTO{}
	readyToSave := false
	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}
		if update.Message != nil {
			inputData.ChatId = update.Message.Chat.ID
			inputData.User = update.Message.From
		} else if update.CallbackQuery != nil {
			inputData.ChatId = update.CallbackQuery.Message.Chat.ID
			inputData.User = update.CallbackQuery.Message.From
		}
		msg := tgbotapi.NewMessage(inputData.ChatId, "")
		msg.ReplyMarkup = commandsKey
		err := b.getCommand(update, inputData)
		if err != nil {
			b.log.Error(err)
			msg.Text = err.Error()
			_, err = b.bot.Send(msg)
			if err != nil {
				b.log.Error(err)
			}
			continue
		}
		err = b.checkAccess(inputData)
		if err != nil {
			b.log.Error(err)
			msg.Text = err.Error()
			_, err = b.bot.Send(msg)
			if err != nil {
				b.log.Error(err)
			}
			continue
		}

		executeMsg, err := b.executeCommand(inputData, productDTO, &readyToSave)
		if err != nil {
			b.log.Error(err)
			msg.Text = err.Error()
			_, err = b.bot.Send(msg)
			if err != nil {
				b.log.Error(err)
			}
			continue
		}
		if readyToSave {
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Сохранить", "add_product"),
					tgbotapi.NewInlineKeyboardButtonData("Отмена", "cancel_add_product"),
				),
			)
		}
		msg.Text = executeMsg
		_, err = b.bot.Send(msg)
		if err != nil {
			b.log.Error(err)
			continue
		}
	}
}

func (b BotHandler) executeCommand(inputData *InputData, productDTO *product.CreateProductDTO, readyToSave *bool,
) (string, error) {
	switch inputData.Command {
	case "start":
		text := "Привет, тебя приветствует система контроля расходов."
		err := b.startMessage(inputData)
		if err != nil {
			return "", err
		}
		return text, nil
	case "add_product":
		var result map[string]string
		var err error
		var text string
		if !*readyToSave {
			result, err = b.prepareAddProductData(inputData, productDTO)
			if err != nil {
				return "", err
			}
		}
		text = result["text"]
		if save, ok := result["save"]; ok && save == saveProduct {
			*readyToSave = true
			return text, nil
		}
		if *readyToSave {
			text, err = b.AddProduct(*productDTO)
			if err != nil {
				return "", err
			}
		}

		return text, nil
	case "get_products":
		products, err := b.GetProducts(inputData)
		if err != nil {
			return "", err
		}
		var text strings.Builder
		text.WriteString(fmt.Sprintf("Период: %s - %s\n", beginMonth, now))
		for _, _product := range products {
			text.WriteString(fmt.Sprintf("Название: %s\nЦена: %.2f руб.\nОписание: %s\n"+
				"Дата покупки: %s\nПользователь: %s\n\n", _product.Name, _product.Price, _product.Description,
				_product.BuyAt.Format(inputDateFormat), _product.User))
		}
		return text.String(), nil
	case "get_report":
		report, err := b.GetReport(inputData)
		if err != nil {
			return "", err
		}
		sumByUser := strings.Builder{}
		for key, sum := range report {
			if key != "sum" {
				sumByUser.WriteString(fmt.Sprintf("%s: %.2f руб.\n", key, sum))
			}
		}
		return fmt.Sprintf("Период: %s - %s\n%sИтоговая сумма: %.2f руб.", beginMonth, now,
			sumByUser.String(), report["sum"]), nil
	case "cancel_add_product":
		inputData = nil
		productDTO = &product.CreateProductDTO{}
		*readyToSave = false
	default:
		return "", errors.New("Не понял!")
	}

	return "", nil
}

func (b BotHandler) startMessage(inputData *InputData) error {
	err := b.checkAccess(inputData)
	if err != nil {
		b.log.Error(err)
		return err
	}

	return nil
}
func (b BotHandler) checkAccess(inputData *InputData) error {
	//	TODO: Создать табличку session в которой будут храниться токен и данные о пользователе
	//	затем из этой таблички брать токен по логину пользователя в телеграме
	//	полученый токен проверить на валидность, и если ве нормально, вернуть true, в противном случае false
	if strings.Contains(inputData.Text, ":") && strings.Contains(inputData.Text, "@") {
		split := strings.Split(inputData.Text, ";")
		if len(split) == 2 {
			email, password := split[0], split[1]
			dbUser, err := b.getUser(email)
			if err != nil {
				return err
			}
			if !b.hasherManager.Verify(dbUser.Password, password) {
				b.log.Error(err)
				return errors.New("Не корректный email или пароль")
			}
			token, err := b.jwtManager.Generate(dbUser.Email, 60, b.cfg.SignedKey)
			if err != nil {
				return err
			}
			err = b.updateOrCreateSession(dbUser.Id, token)
			if err != nil {
				return err
			}
		}
		return errors.New("Переданы не корректные данные для получения токена.")
	}
	return errors.New("Для работы в системе, нужно авторизоваться. Введите email и пароль через : email:password")
}

func (b BotHandler) getUser(email string) (user.UserServiceOutput, error) {
	_, err := mail.ParseAddress(email)
	if err != nil {
		b.log.Error(err)
		return user.UserServiceOutput{},
			errors.New("не удалось распарсить email адрес, возможно был введен не корректный адрес.")
	}
	findUser, err := b.userService.GetByEmail(context.Background(), email)
	if err != nil {
		b.log.Error(err)
		return user.UserServiceOutput{}, errors.New("Пользователь с таким email не найден")
	}

	return findUser, nil
}

func (b BotHandler) getCommand(update tgbotapi.Update, inputData *InputData) error {
	var err error
	switch {
	case update.CallbackQuery != nil:
		b.prepareCallbackQuery(update.CallbackQuery, inputData)
	//case update.Message.Photo != nil:
	//	err = b.preparePhoto(update.Message.Photo, inputData)
	//	if err != nil {
	//		b.log.Error(err)
	//		return err
	//	}
	case update.Message != nil:
		err = b.prepareMessage(update.Message, inputData)
		if err != nil {
			b.log.Error(err)
			return err
		}
	}

	return nil
}

func (b BotHandler) prepareCallbackQuery(callbackQuery *tgbotapi.CallbackQuery, inputData *InputData) {
	var splitStr []string
	sep := " "
	if strings.Contains(callbackQuery.Data, "\n") {
		splitStr = strings.Split(callbackQuery.Data, "\n")
		sep = "\n"
	} else {
		splitStr = strings.Split(callbackQuery.Data, " ")
	}
	inputData.Command = strings.TrimSpace(strings.TrimSpace(splitStr[0]))
	if len(splitStr) > 1 {
		inputData.Text = strings.TrimSpace(strings.Join(splitStr[1:], sep))
	}

}

func (b BotHandler) prepareMessage(message *tgbotapi.Message, inputData *InputData) error {
	inputData.Command = "add_product"
	if message.IsCommand() {
		return errors.New("Команды передаются только через кнопки.")
	} else {
		if strings.Contains(message.Text, ":") && strings.Contains(message.Text, "@") {
			inputData.Command = "start"
			inputData.Text = message.Text
		}
	}
	return nil
}

//func (b BotHandler) preparePhoto(photo *[]tgbotapi.PhotoSize, inputData *InputData) error {
//	var err error
//	inputData.Text, err = b.getQrCodeText(photo)
//	if err != nil {
//		b.log.Error(err)
//		return err
//	}
//	inputData.Command = "add_product"
//
//	return nil
//}

//func (b BotHandler) getQrCodeText(photo *[]tgbotapi.PhotoSize) (string, error) {
//	resp, err := b.bot.GetFile(tgbotapi.FileConfig{photo[1].FileID})
//	if err != nil {
//		return "", err
//	}
//	url := fmt.Sprintf("%s/file/bot%s/%s", telegramApiHost, b.bot.Token, resp.FilePath)
//	r, err := http.Get(url) //загружаем изображение с сервера telegram
//	if err != nil {
//		b.log.Errorf("Не удалось загрузить изображение с сервера. %v\n", err)
//		return "", err
//	}
//	defer r.Body.Close()
//	img, err := jpeg.Decode(r.Body) //конвертируем io.Reader в image.Image
//	if err != nil {
//		b.log.Infof("Не смог декодировать изображение: %v\n", err)
//		return "", err
//	}
//	qrCodes, err := goqr.Recognize(img)
//	if err != nil {
//		b.log.Infof("Recognize failed: %v\n", err)
//		return "", err
//	}
//	for _, qrCode := range qrCodes {
//		return string(qrCode.Payload), nil
//	}
//
//	return "", nil
//
//}

func (b BotHandler) prepareAddProductData(data *InputData, productDto *product.CreateProductDTO) (map[string]string, error) {
	var err error
	result := make(map[string]string)
	if productDto.Name == "" {
		if strings.Contains(data.Text, "+") {
			data.Text = strings.Replace(data.Text, "+", "", 1)
			data.Text = strings.TrimSpace(data.Text)
			productDto.Name = data.Text
			productDto.Description = data.Text
			result["text"] = "Введите цену продукта. Цена должна начинаться со знака +"
		} else {
			result["text"] = "Введите название продукта. Название должно начинаться со знака +"
		}
		return result, err
	} else if productDto.Price == 0 {
		if strings.Contains(data.Text, "+") {
			data.Text = strings.Replace(data.Text, "+", "", 1)
			data.Text = strings.TrimSpace(data.Text)
			price, err := strconv.ParseFloat(data.Text, 32)
			if err != nil {
				return result, err
			}
			productDto.Price = float32(price)
			result["text"] = "Введите дату покупки продукта (не обязательно). Дата должна начинаться со знака +"
		} else {
			result["text"] = "Введите цену продукта. Цена должна начинаться со знака +"
		}
		return result, err
	} else if productDto.BuyAt.IsZero() {
		if strings.Contains(data.Text, "+") {
			data.Text = strings.Replace(data.Text, "+", "", 1)
			data.Text = strings.TrimSpace(data.Text)
			productDto.BuyAt, err = time.ParseInLocation(inputDateFormat, data.Text, time.Local)
			if err != nil {
				return result, err
			}
			result["text"] = fmt.Sprintf("Название: %s\nЦена: %.2f\nДата покупки: %s\nСохранить?",
				productDto.Name, productDto.Price, productDto.BuyAt.Format(inputDateFormat))
			result["save"] = saveProduct
			productDto.User = data.User.UserName
		} else {
			result["text"] = "Введите дату покупки продукта (не обязательно). Дата должна начинаться со знака +"
		}
	}

	return result, err
}

func (b BotHandler) updateOrCreateSession(userId, token string) error {
	//TODO: реализовать сессии
	session, err := sessionService.GetByUserId(userId)
	if err != nil {
		return err
	}
	if len(session) == 0 {
		createSessionDTO := sessionService.CreateSession(userId, token)
		err = b.sessionService.Save(createSessionDTO)
	} else {
		updateSessionDTO := sessionService.UpdateSession(session.Id, token)
		err = b.sessionService.Update(updateSessionDTO)
	}
	return nil
}

func initConfig(updateTimeout *int, offset *int) tgbotapi.UpdateConfig {
	updateConfig := tgbotapi.NewUpdate(0)
	if offset != nil {
		updateConfig = tgbotapi.NewUpdate(*offset)
	}
	updateConfig.Timeout = defaultTimeout
	if updateTimeout != nil {
		updateConfig = tgbotapi.NewUpdate(*updateTimeout)
	}

	return updateConfig
}
