package domain

import "time"

type ContributionType string

const (
	OneTime     ContributionType = "one_time"
	Periodic    ContributionType = "periodic"
)

type Period string

const (
	Daily       Period = "daily"
	Monthly     Period = "monthly"
	Quarterly   Period = "quarterly"
	SemiYearly  Period = "semi_yearly"
	Yearly      Period = "yearly"
)

type Contribution struct {
	Type      ContributionType `json:"type"`
	Amount    float64          `json:"amount"`
	Period    Period           `json:"period"`
	Month     int              `json:"month"`
	Year      int              `json:"year"`
	CreatedAt time.Time        `json:"created_at"`
}

type InvestmentParams struct {
	InitialCapital   float64          `json:"initialCapital"`
	AnnualReturn     float64          `json:"annualReturn"`
	InvestmentMonths int              `json:"investmentMonths"`
	TaxRate          float64          `json:"taxRate"`
	InflationRate    float64          `json:"inflationRate"`
	Contributions    []Contribution   `json:"contributions"`
	ReinvestEnabled  bool             `json:"reinvestEnabled"`
	ReinvestPeriod   Period           `json:"reinvestPeriod"`
}

type MonthResult struct {
	Month             int     `json:"month"`
	StartBalance      float64 `json:"startBalance"`
	ReturnEarned      float64 `json:"returnEarned"`
	Contributions     float64 `json:"contributions"`
	TaxPaid           float64 `json:"taxPaid"`
	EndBalance        float64 `json:"endBalance"`
	EndBalanceReal    float64 `json:"endBalanceReal"`
	CumulativeReturn  float64 `json:"cumulativeReturn"`
	CumulativeContrib float64 `json:"cumulativeContrib"`
}

type CalculationResult struct {
	ID               string         `json:"id"`
	Params           InvestmentParams `json:"params"`
	MonthlyResults   []MonthResult  `json:"monthlyResults"`
	FinalBalance     float64        `json:"finalBalance"`
	FinalBalanceReal float64        `json:"finalBalanceReal"`
	TotalContrib     float64        `json:"totalContrib"`
	TotalReturn      float64        `json:"totalReturn"`
	TotalTax         float64        `json:"totalTax"`
	NetProfit        float64        `json:"netProfit"`
	NetProfitReal    float64        `json:"netProfitReal"`
}
