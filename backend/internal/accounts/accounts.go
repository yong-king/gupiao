package accounts

import (
	"errors"
	"strings"

	"jijin/backend/internal/holdings"
	"jijin/backend/internal/marketdata"
	"jijin/backend/internal/settings"
)

type Config struct {
	ID          string
	UserID      string
	Alias       string
	RefreshMode settings.RefreshMode
	ReadOnly    bool
	Metadata    map[string]string
}

func (c Config) Validate() error {
	if c.ID == "" || c.UserID == "" || strings.TrimSpace(c.Alias) == "" {
		return errors.New("account id, user id and alias are required")
	}
	if !c.ReadOnly {
		return errors.New("account monitoring must be read-only")
	}
	if c.RefreshMode != settings.RefreshModeManual && c.RefreshMode != settings.RefreshModeConservative && c.RefreshMode != settings.RefreshModeStandard {
		return errors.New("unsupported refresh mode")
	}
	for key := range c.Metadata {
		normalized := strings.ToLower(key)
		if strings.Contains(normalized, "password") || strings.Contains(normalized, "secret") || strings.Contains(normalized, "token") {
			return errors.New("sensitive account metadata is not allowed")
		}
	}
	return nil
}

type PositionRisk struct {
	Market          string
	Symbol          string
	Value           float64
	UnrealizedPnL   float64
	Concentration   float64
	RiskObservation string
}

type PortfolioRisk struct {
	TotalValue  float64
	Positions   []PositionRisk
	MissingData []string
}

func CalculateRisk(items []holdings.Holding, snapshots []marketdata.Snapshot, concentrationThreshold float64) PortfolioRisk {
	priceByKey := map[string]float64{}
	for _, snapshot := range snapshots {
		priceByKey[snapshot.Market+":"+snapshot.Symbol] = snapshot.Price
	}

	var result PortfolioRisk
	values := make([]float64, len(items))
	for i, item := range items {
		key := item.Market + ":" + item.Symbol
		price, ok := priceByKey[key]
		if !ok {
			result.MissingData = append(result.MissingData, key)
			continue
		}
		value := item.Quantity * price
		values[i] = value
		result.TotalValue += value
	}

	for i, item := range items {
		key := item.Market + ":" + item.Symbol
		price, ok := priceByKey[key]
		if !ok {
			continue
		}
		value := values[i]
		concentration := 0.0
		if result.TotalValue > 0 {
			concentration = value / result.TotalValue
		}
		observation := "继续观察"
		if concentration > concentrationThreshold {
			observation = "持仓集中度偏高，建议人工确认风险"
		}
		result.Positions = append(result.Positions, PositionRisk{
			Market:          item.Market,
			Symbol:          item.Symbol,
			Value:           value,
			UnrealizedPnL:   (price - item.CostBasis) * item.Quantity,
			Concentration:   concentration,
			RiskObservation: observation,
		})
	}
	return result
}
