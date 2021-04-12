package ps_event_bus

/*
	Подписчик на событие
*/
type EventBusSubscriber struct {
	/*
		ID пользователя. Уникален только в паре с Type
	*/
	ID uint `json:"id"`
	/*
		customer - клиенты
		user - пользователи админ панели
	*/
	Type string `json:"type"`
}
