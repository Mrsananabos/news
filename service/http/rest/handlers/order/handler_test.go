package order

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"log"
	"net/http/httptest"
	"orderService/internal/models"
	"orderService/internal/service/mocks"
	"testing"
	"time"
)

var uid uuid.UUID
var dateCreated time.Time

func init() {
	var err error
	uid, err = uuid.Parse("1e9ad4fb-2615-46f9-9458-20b59253086b")
	if err != nil {
		log.Fatalf("Error while parse uuid: %v", err)
	}

	dateCreatedStr := "2021-11-26T06:22:19Z"
	dateCreated, err = time.Parse(time.RFC3339, dateCreatedStr)
	if err != nil {
		log.Fatalf("Error while parse date: %v", err)
	}
}

func TestHandler_GetOrderById(t *testing.T) {
	orderView := models.OrderView{DeliveryService: "meest", DateCreated: dateCreated, Delivery: models.DeliveryView{
		Name:    "Test Testov",
		Phone:   "+9720000000",
		Zip:     "2639809",
		City:    "Moscow",
		Address: "Prospekt Mira 15",
		Region:  "Moscow",
		Email:   "test@gmail.com",
	},
		Payment: models.PaymentView{
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			DeliveryCost: 1500,
			GoodsTotal:   317,
		},
		Items: []models.ItemView{{
			Name:       "Mascaras",
			TotalPrice: 317,
			Brand:      "Vivienne Sabo",
		},
		}}

	orderViewResponse := `{
    "DeliveryService": "meest",
    "DateCreated": "2021-11-26T06:22:19Z",
    "Delivery": {
        "Name": "Test Testov",
        "Phone": "+9720000000",
        "Zip": "2639809",
        "City": "Moscow",
        "Address": "Prospekt Mira 15",
        "Region": "Moscow",
        "Email": "test@gmail.com"
    },
    "Payment": {
        "Currency": "USD",
        "Provider": "wbpay",
        "Amount": 1817,
        "DeliveryCost": 1500,
        "GoodsTotal": 317
    },
    "Items": [
        {
            "Name": "Mascaras",
            "TotalPrice": 317,
            "Brand": "Vivienne Sabo"
        }
    ]
}`

	t.Run("Success", func(t *testing.T) {
		c := gomock.NewController(t)
		defer c.Finish()

		mockOrderService := new(mocks.IOrderService)
		mockOrderService.On("GetById", uid).Return(orderView, nil)

		handler := NewHandler(mockOrderService)
		g := gin.New()
		g.GET("/order/:uid", handler.GetOrderById)

		h := httptest.NewRecorder()
		r := httptest.NewRequest("GET", fmt.Sprintf("/order/%s", uid.String()), nil)

		g.ServeHTTP(h, r)

		assert.Equal(t, 200, h.Code)
		assert.JSONEq(t, orderViewResponse, h.Body.String())
		mockOrderService.AssertCalled(t, "GetById", uid)
	})

	t.Run("UidIsNotUUIDType", func(t *testing.T) {
		c := gomock.NewController(t)
		defer c.Finish()

		mockOrderService := new(mocks.IOrderService)

		handler := NewHandler(mockOrderService)
		g := gin.New()
		g.GET("/order/:uid", handler.GetOrderById)

		h := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/order/1", nil)

		g.ServeHTTP(h, r)

		assert.Equal(t, 400, h.Code)
		assert.JSONEq(t, `{"error":"uid is not UUID format"}`, h.Body.String())
		mockOrderService.AssertNotCalled(t, "GetById")
	})

	t.Run("OrderNotFound", func(t *testing.T) {
		mockOrderService := new(mocks.IOrderService)
		mockOrderService.On("GetById", uid).Return(models.OrderView{}, fmt.Errorf("record not found"))
		handler := NewHandler(mockOrderService)

		c := gomock.NewController(t)
		defer c.Finish()

		g := gin.New()
		g.GET("/order/:uid", handler.GetOrderById)

		h := httptest.NewRecorder()
		r := httptest.NewRequest("GET", fmt.Sprintf("/order/%s", uid.String()), nil)

		g.ServeHTTP(h, r)

		assert.Equal(t, 404, h.Code)
		assert.JSONEq(t, `{"error":"record not found"}`, h.Body.String())
		mockOrderService.AssertCalled(t, "GetById", uid)
	})
}
