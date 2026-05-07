package reports

import (
	"time"

	"jijin/backend/internal/alerts"
	"jijin/backend/internal/marketdata"
)

type StockDetail struct {
	Market         string
	Symbol         string
	LatestPrice    *float64
	LatestDataTime *time.Time
	LatestAlert    *alerts.Event
	RiskLevel      alerts.RiskLevel
	Snapshots      []marketdata.Snapshot
	Alerts         []alerts.Event
}

func BuildStockDetail(market string, symbol string, snapshots []marketdata.Snapshot, events []alerts.Event) StockDetail {
	detail := StockDetail{Market: market, Symbol: symbol, RiskLevel: alerts.RiskLow}
	for _, snapshot := range snapshots {
		if snapshot.Market != market || snapshot.Symbol != symbol {
			continue
		}
		detail.Snapshots = append(detail.Snapshots, snapshot)
		price := snapshot.Price
		dataTime := snapshot.DataTime
		if detail.LatestDataTime == nil || dataTime.After(*detail.LatestDataTime) {
			detail.LatestPrice = &price
			detail.LatestDataTime = &dataTime
		}
	}
	for _, event := range events {
		if event.Market != market || event.Symbol != symbol {
			continue
		}
		detail.Alerts = append(detail.Alerts, event)
		if detail.LatestAlert == nil || event.CreatedAt.After(detail.LatestAlert.CreatedAt) {
			copy := event
			detail.LatestAlert = &copy
			detail.RiskLevel = event.RiskLevel
		}
	}
	return detail
}

type DailyReport struct {
	ID                string
	UserID            string
	Date              string
	Summary           string
	RiskPoints        []string
	NeedsConfirmation []string
	DataTime          time.Time
	CreatedAt         time.Time
}

func GenerateDailyReport(id string, userID string, events []alerts.Event, now time.Time) DailyReport {
	report := DailyReport{
		ID:        id,
		UserID:    userID,
		Date:      now.Format("2006-01-02"),
		DataTime:  now,
		CreatedAt: now,
	}
	if len(events) == 0 {
		report.Summary = "今日没有触发新的监控提醒。"
		report.NeedsConfirmation = []string{"继续观察股票池和持仓数据。"}
		return report
	}
	report.Summary = "今日存在需要关注的监控提醒。"
	for _, event := range events {
		if event.UserID != userID {
			continue
		}
		report.RiskPoints = append(report.RiskPoints, event.Market+":"+event.Symbol+" "+event.Summary)
		report.NeedsConfirmation = append(report.NeedsConfirmation, "人工确认 "+event.Market+":"+event.Symbol+" 的 "+string(event.Signal)+" 提醒。")
	}
	if len(report.RiskPoints) == 0 {
		report.Summary = "今日没有触发新的监控提醒。"
	}
	return report
}

type Repository struct {
	reports map[string]DailyReport
}

func NewRepository() *Repository {
	return &Repository{reports: make(map[string]DailyReport)}
}

func (r *Repository) Save(report DailyReport) {
	r.reports[report.ID] = report
}

func (r *Repository) FindByID(id string) (DailyReport, bool) {
	report, ok := r.reports[id]
	return report, ok
}
