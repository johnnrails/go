package services

type Tavern struct {
	OrderService   *OrderService
	BillingService interface{}
}

func NewTavern(os *OrderService, bs interface{}) *Tavern {
	return &Tavern{
		OrderService:   os,
		BillingService: bs,
	}
}
