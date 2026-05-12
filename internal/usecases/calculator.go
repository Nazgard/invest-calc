package usecases

import (
	"context"
	"math"

	"invest-calc/internal/domain"
)

type Calculator struct {
	repo domain.CalculatorRepository
}

func NewCalculator(repo domain.CalculatorRepository) *Calculator {
	return &Calculator{repo: repo}
}

func (c *Calculator) Calculate(ctx context.Context, params domain.InvestmentParams) domain.CalculationResult {
	monthlyRate := params.AnnualReturn / 100 / 12
	monthlyInflation := params.InflationRate / 100 / 12
	dailyRate := params.AnnualReturn / 100 / 365

	daysInMonth := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}

	principal := params.InitialCapital
	cumulativeContrib := params.InitialCapital
	cumulativeReturn := 0.0
	totalTax := 0.0
	accumulatedIncome := 0.0
	accumulatedReturn := 0.0
	yearlyAccumulatedReturn := 0.0

	var results []domain.MonthResult

	for month := 1; month <= params.InvestmentMonths; month++ {
		monthContrib := c.getContributionsForMonth(params.Contributions, month)
		principal += monthContrib
		cumulativeContrib += monthContrib

		var returnEarned float64
		var taxPaid float64

		if params.ReinvestEnabled {
			balance := principal + accumulatedIncome

			if params.ReinvestPeriod == domain.Daily {
				days := daysInMonth[(month-1)%12]
				returnEarned = 0.0
				for d := 0; d < days; d++ {
					dailyReturn := balance * dailyRate
					balance += dailyReturn
					returnEarned += dailyReturn
				}
				cumulativeReturn += returnEarned
				accumulatedReturn += returnEarned
				yearlyAccumulatedReturn += returnEarned
				accumulatedIncome = balance - principal
			} else {
				returnEarned = balance * monthlyRate
				cumulativeReturn += returnEarned
				accumulatedReturn += returnEarned
				yearlyAccumulatedReturn += returnEarned

				shouldReinvest := false
				switch params.ReinvestPeriod {
				case domain.Monthly:
					shouldReinvest = true
				case domain.Quarterly:
					shouldReinvest = (month % 3 == 0)
				case domain.SemiYearly:
					shouldReinvest = (month % 6 == 0)
				case domain.Yearly:
					shouldReinvest = (month % 12 == 0)
				}

				if shouldReinvest && accumulatedReturn > 0 {
					accumulatedIncome += accumulatedReturn
					accumulatedReturn = 0
				}
			}
		} else {
			returnEarned = principal * monthlyRate
			cumulativeReturn += returnEarned
			yearlyAccumulatedReturn += returnEarned
			accumulatedIncome += returnEarned
		}

		isLastMonth := (month == params.InvestmentMonths)
		isYearEnd := (month%12 == 0) || isLastMonth

		if isYearEnd && yearlyAccumulatedReturn > 0 {
			taxPaid = yearlyAccumulatedReturn * params.TaxRate / 100
			accumulatedIncome -= taxPaid
			totalTax += taxPaid
			yearlyAccumulatedReturn = 0
		}

		balance := principal + accumulatedIncome
		inflationFactor := math.Pow(1+monthlyInflation, float64(month))
		endBalanceReal := balance / inflationFactor

		results = append(results, domain.MonthResult{
			Month:             month,
			StartBalance:      principal + accumulatedIncome - monthContrib,
			ReturnEarned:      returnEarned,
			Contributions:     monthContrib,
			TaxPaid:           taxPaid,
			EndBalance:        balance,
			EndBalanceReal:    endBalanceReal,
			CumulativeReturn:  cumulativeReturn,
			CumulativeContrib: cumulativeContrib,
		})
	}

	finalBalance := principal + accumulatedIncome
	if len(results) > 0 {
		finalBalance = results[len(results)-1].EndBalance
	}

	finalBalanceReal := finalBalance / math.Pow(1+monthlyInflation, float64(params.InvestmentMonths))
	totalContrib := cumulativeContrib
	netProfit := finalBalance - totalContrib
	netProfitReal := finalBalanceReal - totalContrib

	result := domain.CalculationResult{
		Params:           params,
		MonthlyResults:   results,
		FinalBalance:     finalBalance,
		FinalBalanceReal: finalBalanceReal,
		TotalContrib:     totalContrib,
		TotalReturn:      cumulativeReturn,
		TotalTax:         totalTax,
		NetProfit:        netProfit,
		NetProfitReal:    netProfitReal,
	}

	if c.repo != nil {
		// Saving removed as history feature is disabled
	}

	return result
}

func (c *Calculator) getContributionsForMonth(contributions []domain.Contribution, month int) float64 {
	total := 0.0
	for _, contrib := range contributions {
		switch contrib.Type {
		case domain.OneTime:
			if contrib.Month == month {
				total += contrib.Amount
			}
		case domain.Periodic:
			if month >= contrib.Month {
				switch contrib.Period {
				case domain.Monthly:
					total += contrib.Amount
				case domain.Quarterly:
					if (month-contrib.Month)%3 == 0 {
						total += contrib.Amount
					}
				case domain.Yearly:
					if (month-contrib.Month)%12 == 0 {
						total += contrib.Amount
					}
				}
			}
		}
	}
	return total
}

func (c *Calculator) SaveResult(ctx context.Context, result domain.CalculationResult) (string, error) {
	return c.repo.SaveCalculation(ctx, result)
}

func (c *Calculator) GetResult(ctx context.Context, id string) (*domain.CalculationResult, error) {
	return c.repo.GetCalculation(ctx, id)
}

func (c *Calculator) ListResults(ctx context.Context) ([]domain.CalculationResult, error) {
	return c.repo.ListCalculations(ctx)
}
