package loyalty

import (
	"encoding/json"
	"go.uber.org/zap"
	"gofermark_personal/internal/model"
	"io"
	"net/http"
)

type LoyaltyHttpRepository struct {
	address string
	http    *http.Client
	logger  *zap.Logger
}

func NewLoyaltyHttpRepository(address string, http *http.Client, logger *zap.Logger) *LoyaltyHttpRepository {
	return &LoyaltyHttpRepository{address: address, http: http, logger: logger}
}

func (repository LoyaltyHttpRepository) GetOrderInfo(orderNumber string) (*model.LoyaltyOrderInfo, error) {
	fullUrl := repository.address + "/api/orders/" + orderNumber
	repository.logger.Info(fullUrl)

	resp, err := repository.http.Get(fullUrl)
	defer resp.Body.Close()

	if err != nil {
		repository.logger.Error(err.Error())

		return nil, err
	}

	var loyaltyOrderInfo model.LoyaltyOrderInfo
	if err := json.NewDecoder(resp.Body).Decode(&loyaltyOrderInfo); err != nil {
		repository.logger.Error("failed to decode loyalty order", zap.String("error", err.Error()))
		body, _ := io.ReadAll(resp.Body)
		repository.logger.Error("response", zap.String("body", string(body)))
	}

	return &loyaltyOrderInfo, nil
}
