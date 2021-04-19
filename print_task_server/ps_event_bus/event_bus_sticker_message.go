package ps_event_bus

import "encoding/json"

type StickerMessageUser struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type StickerMessage struct {
	CustomerId         string             `json:"customer_id"`
	CustomerTitle      string             `json:"customer_title"`
	CustomerLogo       string             `json:"customer_logo"`
	RegionCustomerId   string             `json:"region_customer_id"`
	RegionCustomerLogo string             `json:"region_customer_logo"`
	RegionName         string             `json:"region_name"`
	SupplierLogo       string             `json:"supplier_logo"`
	SupplierPriceLogo  string             `json:"supplier_price_logo"`
	Oem                string             `json:"oem"`
	Ean                string             `json:"ean"`
	MakeName           string             `json:"make_name"`
	DetailName         string             `json:"detail_name"`
	Quantity           int                `json:"qnt"`
	QuantityAccept     int                `json:"qnt_accept"`
	QuantityIncome     int                `json:"qnt_income"`
	DeliveryType       string             `json:"delivery_type"`
	Url                string             `json:"url"`
	Site               string             `json:"site"`
	User               StickerMessageUser `json:"user"`
}

var testStickerMessage = StickerMessage{
	CustomerId:         "123",
	CustomerTitle:      `ОАО "Атмосфера"`,
	CustomerLogo:       "ASFE",
	RegionCustomerId:   "342",
	RegionCustomerLogo: "ASFE-CLI",
	RegionName:         "Курск",
	SupplierLogo:       "KORE",
	SupplierPriceLogo:  "KORE-MSK",
	Oem:                "C110",
	Ean:                "ean-code",
	MakeName:           "DOLZ",
	DetailName:         "Помпа",
	Quantity:           3,
	QuantityAccept:     2,
	QuantityIncome:     1,
	DeliveryType:       "Контейнер",
	Url:                "url",
	Site:               "kuzparts.ru",
	User: StickerMessageUser{
		Id:   1,
		Name: "Иванов",
	},
}

func GetTestStickerMessage(clientId string) EventBusMessage {

	data, _ := json.Marshal(&testStickerMessage)
	return EventBusMessage{
		CommandType:  "sticker",
		CommandParam: "",
		StickerType:  "current_sticker",
		ClientId:     clientId,
		Data:         data,
		IsBroadcast:  false,
	}
}
