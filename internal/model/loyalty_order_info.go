package model

type LoyaltyOrderInfo struct {
	Number  string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float32 `json:"accrual,omitempty"`
}
