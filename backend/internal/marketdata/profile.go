package marketdata

import (
	"context"
	"strings"
)

type CompanyProfile struct {
	Market         string   `json:"market"`
	Symbol         string   `json:"symbol"`
	Name           string   `json:"name"`
	Sector         string   `json:"sector"`
	Products       []string `json:"products"`
	Business       string   `json:"business"`
	DataSource     string   `json:"data_source"`
	Analysis       string   `json:"analysis"`
	Recommendation string   `json:"recommendation"`
	Disclaimer     string   `json:"disclaimer"`
}

type ProfileProvider interface {
	FetchProfile(ctx context.Context, market string, symbol string) (CompanyProfile, error)
}

func ProfileFromSnapshots(market string, symbol string, snapshots []Snapshot) CompanyProfile {
	profile := EnrichKnownProfile(CompanyProfile{
		Market:     strings.ToUpper(strings.TrimSpace(market)),
		Symbol:     strings.ToUpper(strings.TrimSpace(symbol)),
		DataSource: "snapshot",
	})
	for _, snapshot := range snapshots {
		if snapshot.Market == profile.Market && snapshot.Symbol == profile.Symbol && snapshot.Name != "" {
			profile.Name = snapshot.Name
		}
	}
	if profile.Name == "" {
		profile.Name = profile.Symbol
	}
	return AnalyzeProfile(profile, snapshots)
}

func EnrichKnownProfile(profile CompanyProfile) CompanyProfile {
	key := strings.ToUpper(strings.TrimSpace(profile.Symbol))
	known := map[string]CompanyProfile{
		"AAPL":  {Sector: "Consumer Electronics", Products: []string{"iPhone", "Mac", "iPad", "Wearables", "Services"}, Business: "Apple sells consumer devices, software, and services."},
		"MSFT":  {Sector: "Software and Cloud", Products: []string{"Azure", "Microsoft 365", "Windows", "GitHub", "Xbox"}, Business: "Microsoft sells cloud infrastructure, productivity software, operating systems, and gaming services."},
		"NVDA":  {Sector: "Semiconductors", Products: []string{"GPU", "AI accelerators", "Networking", "CUDA"}, Business: "NVIDIA designs accelerated computing chips, systems, and software used in AI, graphics, and data centers."},
		"TSLA":  {Sector: "Automotive and Energy", Products: []string{"Electric vehicles", "Battery storage", "Charging", "Energy generation"}, Business: "Tesla sells electric vehicles, charging services, and energy products."},
		"AMZN":  {Sector: "E-commerce and Cloud", Products: []string{"Amazon Marketplace", "AWS", "Prime", "Advertising"}, Business: "Amazon operates online commerce, cloud infrastructure, subscriptions, logistics, and advertising businesses."},
		"GOOGL": {Sector: "Internet Services", Products: []string{"Search", "YouTube", "Google Cloud", "Android", "Ads"}, Business: "Alphabet sells digital advertising, cloud services, subscriptions, devices, and platform services."},
		"META":  {Sector: "Internet Services", Products: []string{"Facebook", "Instagram", "WhatsApp", "Reality Labs", "Ads"}, Business: "Meta operates social apps, advertising systems, messaging, and metaverse hardware/software initiatives."},
		"0700":  {Sector: "Internet Services", Products: []string{"WeChat", "Online games", "FinTech", "Cloud", "Advertising"}, Business: "Tencent operates social platforms, games, digital content, fintech, cloud, and advertising businesses."},
	}
	if item, ok := known[key]; ok {
		if profile.Sector == "" {
			profile.Sector = item.Sector
		}
		if len(profile.Products) == 0 {
			profile.Products = item.Products
		}
		if profile.Business == "" {
			profile.Business = item.Business
		}
	}
	if profile.Sector == "" {
		profile.Sector = "Unknown"
	}
	if len(profile.Products) == 0 {
		profile.Products = []string{"待补充产品和业务信息"}
	}
	if profile.Business == "" {
		profile.Business = "待补充公司业务描述。"
	}
	return profile
}

func AnalyzeProfile(profile CompanyProfile, snapshots []Snapshot) CompanyProfile {
	latest := latestSnapshot(profile.Market, profile.Symbol, snapshots)
	signal := "observe"
	analysis := profile.Name + " belongs to " + profile.Sector + ". Business: " + profile.Business
	if latest != nil {
		if latest.ChangePercent >= 3 {
			signal = "risk_watch"
			analysis += " Latest collected price moved sharply upward; review valuation and news before acting."
		} else if latest.ChangePercent <= -3 {
			signal = "risk_watch"
			analysis += " Latest collected price moved sharply downward; review fundamentals, news, and position risk before acting."
		} else {
			analysis += " Latest collected price movement is moderate based on available snapshots."
		}
	}
	profile.Analysis = analysis
	profile.Recommendation = signal
	profile.Disclaimer = "仅用于监控和研究，不构成投资建议或自动交易指令。"
	return profile
}

func latestSnapshot(market string, symbol string, snapshots []Snapshot) *Snapshot {
	var latest *Snapshot
	for i := range snapshots {
		snapshot := &snapshots[i]
		if snapshot.Market != market || snapshot.Symbol != symbol {
			continue
		}
		if latest == nil || snapshot.DataTime.After(latest.DataTime) {
			latest = snapshot
		}
	}
	return latest
}
