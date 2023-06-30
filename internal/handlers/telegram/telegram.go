package telegram

import (
	"errors"
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/cost_control/internal/handlers/telegram/product"
	productRepos "github.com/cost_control/internal/repository/product"
	"github.com/cost_control/internal/service"
	"github.com/jdomzhang/goqr"
	"go.mongodb.org/mongo-driver/mongo"
	"image/jpeg"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	defaultTimeout  = 60
	inputDateFormat = "02.01.2006 15:04"
	telegramApiHost = "https://api.telegram.org"
	saveProduct     = "1"
)

var (
	begin       = time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Local)
	beginMonth  = begin.Format("02.01.2006")
	now         = time.Now().Format("02.01.2006")
	commandsKey = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Добавить новый продукт", "add_product"),
			tgbotapi.NewInlineKeyboardButtonData("Получить продукт", "get_product"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Получить продукты", fmt.Sprintf("get_products %s %s",
				beginMonth, now)),
			tgbotapi.NewInlineKeyboardButtonData("Удалить продукт", "delete_product"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Отчет", fmt.Sprintf("get_report %s %s", beginMonth, now)),
		),
	)
)

type BotHandler struct {
	productHandler product.ProductBotHandler
	Bot            tgbotapi.BotAPI
}

type InputData struct {
	ChatId    int64
	UserName  string
	Arguments string
}

func New(token string, db *mongo.Collection) (*BotHandler, error) {
	bot, err := tgbotapi.NewBotAPI(token)

	if err != nil {
		return nil, err
	}
	repos := productRepos.New(db)
	return &BotHandler{productHandler: *product.New(service.New(repos)), Bot: *bot}, err
}

func (b BotHandler) Start(updateTimeout *int, offset *int) error {
	updateConfig := tgbotapi.NewUpdate(0)
	if offset != nil {
		updateConfig = tgbotapi.NewUpdate(*offset)
	}
	updateConfig.Timeout = defaultTimeout
	if updateTimeout != nil {
		updateConfig = tgbotapi.NewUpdate(*updateTimeout)
	}

	updates, err := b.Bot.GetUpdatesChan(updateConfig)
	if err != nil {
		return err
	}
	inputData := InputData{}
	productDTO := &product.CreateProductDTO{}
	readyToSave := false
	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}
		var command string
		switch {
		case update.CallbackQuery != nil:
			var splitStr []string
			sep := " "
			if strings.Contains(update.CallbackQuery.Data, "\n") {
				splitStr = strings.Split(update.CallbackQuery.Data, "\n")
				sep = "\n"
			} else {
				splitStr = strings.Split(update.CallbackQuery.Data, " ")
			}
			command = strings.TrimSpace(strings.TrimSpace(splitStr[0]))
			inputData.ChatId = update.CallbackQuery.Message.Chat.ID
			inputData.UserName = update.CallbackQuery.Message.From.UserName
			if len(splitStr) > 1 {
				inputData.Arguments = strings.TrimSpace(strings.Join(splitStr[1:], sep))
			}
		case update.Message != nil:
			inputData.ChatId = update.Message.Chat.ID
			inputData.UserName = update.Message.From.UserName
			command = update.Message.Command()
			inputData.Arguments = update.Message.CommandArguments()
			if !update.Message.IsCommand() && strings.Contains(update.Message.Text, "+") {
				command = "add_product"
				inputData.Arguments = update.Message.Text
			}
			if update.Message.Photo != nil {
				photo := *update.Message.Photo
				inputData.Arguments, err = getQrCodeText(photo, b)
				if err != nil {
					log.Print(err)
					return err
				}
				command = "add_product"
			}
		}
		msg := tgbotapi.NewMessage(inputData.ChatId, "")
		switch command {
		case "start":
			msg.Text = "Привет, тебя приветствует система контроля расходов. Выбери что ты хочешь сделать."
			msg.ReplyMarkup = commandsKey
			_, err = b.Bot.Send(msg)
			if err != nil {
				log.Print(err)
			}
		case "add_product":
			var result map[string]string
			if !readyToSave {
				result, err = prepareAddProductData(inputData, productDTO)
			}
			if err != nil {
				msg.Text = err.Error()
			} else {
				msg.Text = result["text"]
				if save, ok := result["save"]; ok && save == saveProduct {
					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
						tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData("Сохранить", "add_product"),
							tgbotapi.NewInlineKeyboardButtonData("Отмена", "cancel_add_product"),
						),
					)
					readyToSave = true
					_, err = b.Bot.Send(msg)
					if err != nil {
						log.Print(err)
					}
					continue
				}
			}
			if readyToSave {
				productDTO.User = inputData.UserName
				msg.Text, err = b.AddProduct(*productDTO)
				msg.ReplyMarkup = commandsKey
				if err != nil {
					log.Print(err)
					msg.Text = err.Error()
				}
			}
			_, err = b.Bot.Send(msg)
			if err != nil {
				log.Print(err)
			}
		case "get_products":
			products, err := b.GetProducts(inputData)
			if err != nil {
				msg.Text = err.Error()
				_, err = b.Bot.Send(msg)
				if err != nil {
					log.Print(err)
				}
				continue
			}
			var text strings.Builder
			for _, _product := range products {
				text.WriteString(fmt.Sprintf("ID: %s\nНазвание: %s\nЦена: %.2f руб.\nОписание: %s\nДата покупки: %s\n"+
					"Пользователь: %s\n\n", _product.Id, _product.Name, _product.Price, _product.Description,
					_product.BuyAt.Format(inputDateFormat), _product.User))
			}
			msg.Text = text.String()
			_, err = b.Bot.Send(msg)
			if err != nil {
				log.Print(err)
			}
		case "get_product":
			_product, err := b.GetProduct(inputData)
			if err != nil {
				msg.Text = err.Error()
				_, err = b.Bot.Send(msg)
				if err != nil {
					log.Print(err)
				}
				continue
			}
			msg.Text = fmt.Sprintf("ID: %s\nНазвание: %s\nЦена: %.2f руб.\nОписание: %s\nДата покупки: %s\nПользователь: %s",
				inputData.Arguments, _product.Name, _product.Price, _product.Description, _product.BuyAt.Format(inputDateFormat),
				_product.User)
			_, err = b.Bot.Send(msg)
			if err != nil {
				log.Print(err)
			}
		case "get_report":
			report, err := b.GetReport(inputData)
			if err != nil {
				msg.Text = err.Error()
				_, err = b.Bot.Send(msg)
				if err != nil {
					log.Print(err)
				}
				continue
			}
			sumByUser := strings.Builder{}
			for key, sum := range report {
				if key != "sum" {
					sumByUser.WriteString(fmt.Sprintf("%s: %.2f руб.\n", key, sum))
				}
			}
			msg.Text = fmt.Sprintf("%sИтоговая сумма: %.2f руб.", sumByUser.String(), report["sum"])
			_, err = b.Bot.Send(msg)
			if err != nil {
				log.Print(err)
			}
		case "delete_product":
			err = b.DeleteProduct(inputData)
			if err != nil {
				msg.Text = err.Error()
				_, err = b.Bot.Send(msg)
				if err != nil {
					log.Print(err)
				}
				continue
			}
			msg.Text = "Продукт был удален."
			_, err = b.Bot.Send(msg)
			if err != nil {
				log.Print(err)
			}
		case "cancel_add_product":
			inputData = InputData{}
			productDTO = &product.CreateProductDTO{}
			readyToSave = false
		default:
			msg.Text = "Не понял!"
			_, err = b.Bot.Send(msg)
			if err != nil {
				log.Print(err)
			}
		}

	}
	return nil
}

func (b BotHandler) AddProduct(productDto product.CreateProductDTO) (string, error) {
	_, err := b.productHandler.Create(productDto)
	if err != nil {
		log.Print(err)
		return "", err
	}

	return "Продукт добавлен", nil
}

func (b BotHandler) GetProduct(data InputData) (product.GetProductDTO, error) {
	var err error
	if len(data.Arguments) == 0 {
		return product.GetProductDTO{}, errors.New("Не переданы идентификатор продукта.")
	}
	_product, err := b.productHandler.GetById(data.Arguments)
	if err != nil {
		return product.GetProductDTO{}, errors.New(fmt.Sprintf("Не удалось выполнить запрос. %v", err))
	}
	return product.GetProductDTO{
		Id:          _product.Id,
		Name:        _product.Name,
		Price:       _product.Price,
		BuyAt:       &_product.BuyAt,
		Description: _product.Description,
		User:        _product.User,
	}, nil
}

func (b BotHandler) GetProducts(data InputData) ([]product.GetProductDTO, error) {
	products, err := b.productHandler.Get(data.Arguments)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	var getProducts []product.GetProductDTO
	for _, _product := range products {
		getProducts = append(getProducts, product.GetProductDTO{
			Name:        _product.Name,
			Price:       _product.Price,
			Description: _product.Description,
			BuyAt:       &_product.BuyAt,
		})
	}

	return getProducts, nil

}

func (b BotHandler) DeleteProduct(data InputData) error {
	err := b.productHandler.Delete(data.Arguments)
	if err != nil {
		return errors.New(fmt.Sprintf("Ну удалось удалить пролукт с id: %s. %v", data.Arguments, err))
	}

	return nil
}

func (b BotHandler) GetReport(data InputData) (map[string]float32, error) {
	// Реализовать группировку по пользователям
	report, err := b.productHandler.Report(data.Arguments)
	if err != nil {
		return nil, err
	}
	return report, nil
}

func getQrCodeText(photo []tgbotapi.PhotoSize, b BotHandler) (string, error) {
	resp, err := b.Bot.GetFile(tgbotapi.FileConfig{photo[1].FileID})
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("%s/file/bot%s/%s", telegramApiHost, b.Bot.Token, resp.FilePath)
	r, err := http.Get(url) //загружаем изображение с сервера telegram
	if err != nil {
		log.Printf("Не удалось загрузить изображение с сервера. %v\n", err)
		return "", err
	}
	defer r.Body.Close()
	img, err := jpeg.Decode(r.Body) //конвертируем io.Reader в image.Image
	if err != nil {
		log.Printf("Не смог декодировать изображение: %v\n", err)
		return "", err
	}
	qrCodes, err := goqr.Recognize(img)
	if err != nil {
		log.Printf("Recognize failed: %v\n", err)
		return "", err
	}
	for _, qrCode := range qrCodes {
		return string(qrCode.Payload), nil
	}

	return "", nil

}

func prepareAddProductData(data InputData, productDto *product.CreateProductDTO) (map[string]string, error) {
	var err error
	result := make(map[string]string)
	if productDto.Name == "" {
		if strings.Contains(data.Arguments, "+") {
			data.Arguments = strings.Replace(data.Arguments, "+", "", 1)
			data.Arguments = strings.TrimSpace(data.Arguments)
			productDto.Name = data.Arguments
			productDto.Description = data.Arguments
			result["text"] = "Введите цену продукта. Цена должна начинаться со знака +"
		} else {
			result["text"] = "Введите название продукта. Название должно начинаться со знака +"
		}
		return result, err
	} else if productDto.Price == 0 {
		if strings.Contains(data.Arguments, "+") {
			data.Arguments = strings.Replace(data.Arguments, "+", "", 1)
			data.Arguments = strings.TrimSpace(data.Arguments)
			price, err := strconv.ParseFloat(data.Arguments, 32)
			if err != nil {
				return result, err
			}
			productDto.Price = float32(price)
			result["text"] = "Введите дату покупки продукта (не обязательно). Дата должна начинаться со знака +"
		} else {
			result["text"] = "Введите цену продукта. Цена должна начинаться со знака +"
		}
		return result, err
	} else if productDto.BuyAt == nil {
		if strings.Contains(data.Arguments, "+") {
			data.Arguments = strings.Replace(data.Arguments, "+", "", 1)
			data.Arguments = strings.TrimSpace(data.Arguments)
			date, err := time.ParseInLocation(inputDateFormat, data.Arguments, time.Local)
			if err != nil {
				log.Print(err)
				return result, err
			}
			productDto.BuyAt = &date
			result["text"] = fmt.Sprintf("Название: %s\nЦена: %.2f\nДата покупки: %s\nСохранить?",
				productDto.Name, productDto.Price, productDto.BuyAt.Format(inputDateFormat))
			result["save"] = saveProduct
		} else {
			result["text"] = "Введите дату покупки продукта (не обязательно). Дата должна начинаться со знака +"
		}
	}

	return result, err
}