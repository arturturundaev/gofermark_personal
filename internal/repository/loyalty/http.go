package loyalty

import (
	"encoding/json"
	"go.uber.org/zap"
	"gofermark_personal/internal/model"
	"io"
	"net/http"
)

type LoyaltyHTTPRepository struct {
	address string
	http    *http.Client
	logger  *zap.Logger
}

func NewLoyaltyHTTPRepository(address string, http *http.Client, logger *zap.Logger) *LoyaltyHTTPRepository {
	return &LoyaltyHTTPRepository{address: address, http: http, logger: logger}
}

func (repository LoyaltyHTTPRepository) GetOrderInfo(orderNumber string) (*model.LoyaltyOrderInfo, error) {
	fullURL := repository.address + "/api/orders/" + orderNumber
	repository.logger.Info(fullURL)

	resp, err := repository.http.Get(fullURL)

	if err != nil {
		repository.logger.Error(err.Error())

		return nil, err
	}
	defer resp.Body.Close()

	var loyaltyOrderInfo model.LoyaltyOrderInfo
	if err := json.NewDecoder(resp.Body).Decode(&loyaltyOrderInfo); err != nil {
		repository.logger.Error("failed to decode loyalty order", zap.String("error", err.Error()))
		body, _ := io.ReadAll(resp.Body)
		repository.logger.Error("response", zap.String("body", string(body)))
	}

	return &loyaltyOrderInfo, nil
}
