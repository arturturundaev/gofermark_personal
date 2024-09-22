package service

import "gofermark_personal/internal/model"

type LoyalityRepository interface {
	GetOrderInfo(orderNumber string) (*model.LoyaltyOrderInfo, error)
}
