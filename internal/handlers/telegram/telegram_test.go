package telegram

import (
	"github.com/cost_control/internal/handlers/telegram/product"
	productRepos "github.com/cost_control/internal/repository/product"
	"github.com/cost_control/internal/service"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

var fakeStorage *productRepos.ProductFake
var telegramHandler *product.ProductBotHandler
var botHandler BotHandler

func TestMain(m *testing.M) {
	// Write code here to run before tests
	fakeStorage = productRepos.NewFakeDb()
	botHandler = BotHandler{productHandler: *product.New(service.New(fakeStorage))}
	// Run tests
	exitVal := m.Run()

	// Write code here to run after tests

	// Exit with exit value from tests
	os.Exit(exitVal)
}

func TestPrepareDataForAddProduct(t *testing.T) {
	t.Log("Инициализируем обработчик для бота и входные данные")
	inputData := InputData{}
	productDto := &product.CreateProductDTO{}
	t.Log("Отправляем запрос на добавление продукта")
	result, err := prepareAddProductData(inputData, productDto)
	require.Equal(t, "Введите название продукта. Название должно начинаться со знака +", result["text"], err)
	t.Log("В ответ на приглашение отправить название продукта, мы кладем в InputData.Arguments какое нибудь " +
		"название")
	inputData.Arguments = "+ Яблоки"
	result, err = prepareAddProductData(inputData, productDto)
	require.Equal(t, "Введите цену продукта. Цена должна начинаться со знака +", result["text"], err)

	t.Log("В ответ на приглашение отправить цену продукта, мы кладем в InputData.Arguments какую нибудь" +
		"цену")
	inputData.Arguments = "+ 123"
	result, err = prepareAddProductData(inputData, productDto)
	require.Equal(t, "Введите дату покупки продукта (не обязательно). "+
		"Дата должна начинаться со знака +", result["text"], err)

	t.Log("В ответ на приглашение отправить дату покупки, мы кладем в InputData.Arguments какую нибудь " +
		"дату")
	inputData.Arguments = "+ 22.06.2023 14:05"
	result, err = prepareAddProductData(inputData, productDto)
	require.Contains(t, result["text"], "Сохранить?", err)
	require.Equal(t, saveProduct, result["save"], err)
}

func TestAddProduct(t *testing.T) {
	date := time.Now().Local()
	dto := product.CreateProductDTO{
		Name:        "Манго сушеное",
		Price:       564.25,
		BuyAt:       &date,
		Description: "Манго сушеное без добавления сахара",
		User:        "kadr86",
	}
	result, err := botHandler.AddProduct(dto)
	require.Nilf(t, err, "Не должно быть ошибки", err)
	require.Equalf(t, "Продукт добавлен", result, "", err)
}

func TestGetProducts(t *testing.T) {
	data := InputData{
		Arguments: "",
	}
	products, err := botHandler.GetProducts(data)
	require.Nilf(t, err, "Не должно быть ошибки", err)
	require.GreaterOrEqualf(t, len(products), 1, "Должно быть больше 0 продуктов.")
}

func TestGetReport(t *testing.T) {
	data := InputData{
		Arguments: "",
	}
	report, err := botHandler.GetReport(data)
	sum := report["sum"]
	require.Nilf(t, err, "Не должно быть ошибки", err)
	require.NotZerof(t, sum, "Сумма должна быть больше 0.")
}

func TestDeleteProduct(t *testing.T) {
	data := InputData{
		Arguments: "99646b1f-181f-43f7-bf5d-25fc27883293",
	}
	err := botHandler.DeleteProduct(data)
	require.Nilf(t, err, "Не должно быть ошибки", err)

}
