package usecases

import (
	"math"
	"testing"

	"invest-calc/internal/domain"
)

const epsilon = 0.01

func almostEqual(a, b, eps float64) bool {
	return math.Abs(a-b) < eps
}

func TestCalculate_SimpleInterest_NoReinvest(t *testing.T) {
	calc := NewCalculator(nil)

	params := domain.InvestmentParams{
		InitialCapital:   100000,
		AnnualReturn:     12,
		InvestmentMonths: 12,
		TaxRate:          0,
		InflationRate:    0,
		ReinvestEnabled:  false,
	}

	result := calc.Calculate(nil, params)

	expectedFinal := 100000 * (1 + 0.12/12*12)
	if !almostEqual(result.FinalBalance, expectedFinal, epsilon) {
		t.Errorf("FinalBalance = %.2f, want %.2f", result.FinalBalance, expectedFinal)
	}

	if !almostEqual(result.TotalContrib, 100000, epsilon) {
		t.Errorf("TotalContrib = %.2f, want 100000.00", result.TotalContrib)
	}

	expectedReturn := 100000 * 0.12
	if !almostEqual(result.TotalReturn, expectedReturn, epsilon) {
		t.Errorf("TotalReturn = %.2f, want %.2f", result.TotalReturn, expectedReturn)
	}
}

func TestCalculate_CompoundInterest_MonthlyReinvest(t *testing.T) {
	calc := NewCalculator(nil)

	params := domain.InvestmentParams{
		InitialCapital:   100000,
		AnnualReturn:     12,
		InvestmentMonths: 12,
		TaxRate:          0,
		InflationRate:    0,
		ReinvestEnabled:  true,
		ReinvestPeriod:   domain.Monthly,
	}

	result := calc.Calculate(nil, params)

	expectedFinal := 100000 * math.Pow(1+0.12/12, 12)
	if !almostEqual(result.FinalBalance, expectedFinal, epsilon) {
		t.Errorf("FinalBalance = %.2f, want %.2f", result.FinalBalance, expectedFinal)
	}
}

func TestCalculate_CompoundInterest_YearlyReinvest(t *testing.T) {
	calc := NewCalculator(nil)

	params := domain.InvestmentParams{
		InitialCapital:   100000,
		AnnualReturn:     10,
		InvestmentMonths: 24,
		TaxRate:          0,
		InflationRate:    0,
		ReinvestEnabled:  true,
		ReinvestPeriod:   domain.Yearly,
	}

	result := calc.Calculate(nil, params)

	expectedFinal := 100000 * math.Pow(1+0.10, 2)
	if !almostEqual(result.FinalBalance, expectedFinal, epsilon) {
		t.Errorf("FinalBalance = %.2f, want %.2f", result.FinalBalance, expectedFinal)
	}
}

func TestCalculate_CompoundInterest_QuarterlyReinvest(t *testing.T) {
	calc := NewCalculator(nil)

	params := domain.InvestmentParams{
		InitialCapital:   100000,
		AnnualReturn:     12,
		InvestmentMonths: 12,
		TaxRate:          0,
		InflationRate:    0,
		ReinvestEnabled:  true,
		ReinvestPeriod:   domain.Quarterly,
	}

	result := calc.Calculate(nil, params)

	expectedFinal := 100000 * math.Pow(1+0.12/4, 4)
	if !almostEqual(result.FinalBalance, expectedFinal, epsilon) {
		t.Errorf("FinalBalance = %.2f, want %.2f", result.FinalBalance, expectedFinal)
	}
}

func TestCalculate_CompoundInterest_SemiYearlyReinvest(t *testing.T) {
	calc := NewCalculator(nil)

	params := domain.InvestmentParams{
		InitialCapital:   100000,
		AnnualReturn:     10,
		InvestmentMonths: 12,
		TaxRate:          0,
		InflationRate:    0,
		ReinvestEnabled:  true,
		ReinvestPeriod:   domain.SemiYearly,
	}

	result := calc.Calculate(nil, params)

	expectedFinal := 100000 * math.Pow(1+0.10/2, 2)
	if !almostEqual(result.FinalBalance, expectedFinal, epsilon) {
		t.Errorf("FinalBalance = %.2f, want %.2f", result.FinalBalance, expectedFinal)
	}
}

func TestCalculate_CompoundInterest_DailyReinvest(t *testing.T) {
	calc := NewCalculator(nil)

	params := domain.InvestmentParams{
		InitialCapital:   100000,
		AnnualReturn:     10,
		InvestmentMonths: 12,
		TaxRate:          0,
		InflationRate:    0,
		ReinvestEnabled:  true,
		ReinvestPeriod:   domain.Daily,
	}

	result := calc.Calculate(nil, params)

	expectedFinal := 100000 * math.Pow(1+0.10/365, 365)
	if !almostEqual(result.FinalBalance, expectedFinal, 1.0) {
		t.Errorf("FinalBalance = %.2f, want %.2f", result.FinalBalance, expectedFinal)
	}
}

func TestCalculate_TaxCalculation(t *testing.T) {
	calc := NewCalculator(nil)

	params := domain.InvestmentParams{
		InitialCapital:   100000,
		AnnualReturn:     12,
		InvestmentMonths: 12,
		TaxRate:          13,
		InflationRate:    0,
		ReinvestEnabled:  false,
	}

	result := calc.Calculate(nil, params)

	expectedReturnBeforeTax := 100000 * 0.12
	expectedTax := expectedReturnBeforeTax * 0.13
	expectedFinal := 100000 + expectedReturnBeforeTax - expectedTax

	if !almostEqual(result.TotalTax, expectedTax, epsilon) {
		t.Errorf("TotalTax = %.2f, want %.2f", result.TotalTax, expectedTax)
	}

	if !almostEqual(result.FinalBalance, expectedFinal, epsilon) {
		t.Errorf("FinalBalance = %.2f, want %.2f", result.FinalBalance, expectedFinal)
	}
}

func TestCalculate_InflationAdjustment(t *testing.T) {
	calc := NewCalculator(nil)

	params := domain.InvestmentParams{
		InitialCapital:   100000,
		AnnualReturn:     10,
		InvestmentMonths: 12,
		TaxRate:          0,
		InflationRate:    5,
		ReinvestEnabled:  false,
	}

	result := calc.Calculate(nil, params)

	monthlyInflation := 0.05 / 12
	inflationFactor := math.Pow(1+monthlyInflation, 12)
	expectedRealBalance := result.FinalBalance / inflationFactor

	if !almostEqual(result.FinalBalanceReal, expectedRealBalance, epsilon) {
		t.Errorf("FinalBalanceReal = %.2f, want %.2f", result.FinalBalanceReal, expectedRealBalance)
	}
}

func TestCalculate_OneTimeContribution(t *testing.T) {
	calc := NewCalculator(nil)

	params := domain.InvestmentParams{
		InitialCapital: 100000,
		AnnualReturn:   12,
		InvestmentMonths: 12,
		TaxRate:        0,
		InflationRate:  0,
		ReinvestEnabled: false,
		Contributions: []domain.Contribution{
			{
				Type:   domain.OneTime,
				Amount: 50000,
				Month:  1,
			},
		},
	}

	result := calc.Calculate(nil, params)

	expectedTotalContrib := 150000.0
	if !almostEqual(result.TotalContrib, expectedTotalContrib, epsilon) {
		t.Errorf("TotalContrib = %.2f, want %.2f", result.TotalContrib, expectedTotalContrib)
	}

	expectedReturn := 150000 * 0.12
	if !almostEqual(result.TotalReturn, expectedReturn, epsilon) {
		t.Errorf("TotalReturn = %.2f, want %.2f", result.TotalReturn, expectedReturn)
	}
}

func TestCalculate_PeriodicMonthlyContribution(t *testing.T) {
	calc := NewCalculator(nil)

	params := domain.InvestmentParams{
		InitialCapital: 100000,
		AnnualReturn:   12,
		InvestmentMonths: 12,
		TaxRate:        0,
		InflationRate:  0,
		ReinvestEnabled: false,
		Contributions: []domain.Contribution{
			{
				Type:   domain.Periodic,
				Amount: 10000,
				Period: domain.Monthly,
				Month:  1,
			},
		},
	}

	result := calc.Calculate(nil, params)

	expectedContrib := 100000.0 + 10000.0*12
	if !almostEqual(result.TotalContrib, expectedContrib, epsilon) {
		t.Errorf("TotalContrib = %.2f, want %.2f", result.TotalContrib, expectedContrib)
	}
}

func TestCalculate_PeriodicQuarterlyContribution(t *testing.T) {
	calc := NewCalculator(nil)

	params := domain.InvestmentParams{
		InitialCapital: 100000,
		AnnualReturn:   12,
		InvestmentMonths: 12,
		TaxRate:        0,
		InflationRate:  0,
		ReinvestEnabled: false,
		Contributions: []domain.Contribution{
			{
				Type:   domain.Periodic,
				Amount: 30000,
				Period: domain.Quarterly,
				Month:  1,
			},
		},
	}

	result := calc.Calculate(nil, params)

	expectedContrib := 100000.0 + 30000.0*4
	if !almostEqual(result.TotalContrib, expectedContrib, epsilon) {
		t.Errorf("TotalContrib = %.2f, want %.2f", result.TotalContrib, expectedContrib)
	}
}

func TestCalculate_PeriodicYearlyContribution(t *testing.T) {
	calc := NewCalculator(nil)

	params := domain.InvestmentParams{
		InitialCapital: 100000,
		AnnualReturn:   12,
		InvestmentMonths: 24,
		TaxRate:        0,
		InflationRate:  0,
		ReinvestEnabled: false,
		Contributions: []domain.Contribution{
			{
				Type:   domain.Periodic,
				Amount: 100000,
				Period: domain.Yearly,
				Month:  1,
			},
		},
	}

	result := calc.Calculate(nil, params)

	expectedContrib := 100000.0 + 100000.0*2
	if !almostEqual(result.TotalContrib, expectedContrib, epsilon) {
		t.Errorf("TotalContrib = %.2f, want %.2f", result.TotalContrib, expectedContrib)
	}
}

func TestCalculate_MonthlyResultsCount(t *testing.T) {
	calc := NewCalculator(nil)

	params := domain.InvestmentParams{
		InitialCapital:   100000,
		AnnualReturn:     10,
		InvestmentMonths: 36,
		TaxRate:          0,
		InflationRate:    0,
		ReinvestEnabled:  true,
		ReinvestPeriod:   domain.Monthly,
	}

	result := calc.Calculate(nil, params)

	if len(result.MonthlyResults) != 36 {
		t.Errorf("MonthlyResults count = %d, want 36", len(result.MonthlyResults))
	}
}

func TestCalculate_NetProfit(t *testing.T) {
	calc := NewCalculator(nil)

	params := domain.InvestmentParams{
		InitialCapital:   100000,
		AnnualReturn:     12,
		InvestmentMonths: 12,
		TaxRate:          13,
		InflationRate:    0,
		ReinvestEnabled:  false,
	}

	result := calc.Calculate(nil, params)

	expectedNetProfit := result.FinalBalance - result.TotalContrib
	if !almostEqual(result.NetProfit, expectedNetProfit, epsilon) {
		t.Errorf("NetProfit = %.2f, want %.2f", result.NetProfit, expectedNetProfit)
	}
}

func TestCalculate_CompoundVsSimple(t *testing.T) {
	calc := NewCalculator(nil)

	compoundParams := domain.InvestmentParams{
		InitialCapital:   100000,
		AnnualReturn:     12,
		InvestmentMonths: 12,
		TaxRate:          0,
		InflationRate:    0,
		ReinvestEnabled:  true,
		ReinvestPeriod:   domain.Monthly,
	}

	simpleParams := domain.InvestmentParams{
		InitialCapital:   100000,
		AnnualReturn:     12,
		InvestmentMonths: 12,
		TaxRate:          0,
		InflationRate:    0,
		ReinvestEnabled:  false,
	}

	compoundResult := calc.Calculate(nil, compoundParams)
	simpleResult := calc.Calculate(nil, simpleParams)

	if compoundResult.FinalBalance <= simpleResult.FinalBalance {
		t.Errorf("Compound interest FinalBalance (%.2f) should be greater than simple interest (%.2f)",
			compoundResult.FinalBalance, simpleResult.FinalBalance)
	}
}

func TestCalculate_ZeroInitialCapital(t *testing.T) {
	calc := NewCalculator(nil)

	params := domain.InvestmentParams{
		InitialCapital: 0,
		AnnualReturn:   12,
		InvestmentMonths: 12,
		TaxRate:        0,
		InflationRate:  0,
		ReinvestEnabled: false,
		Contributions: []domain.Contribution{
			{
				Type:   domain.Periodic,
				Amount: 10000,
				Period: domain.Monthly,
				Month:  1,
			},
		},
	}

	result := calc.Calculate(nil, params)

	if result.FinalBalance <= 0 {
		t.Errorf("FinalBalance should be positive, got %.2f", result.FinalBalance)
	}
}

func TestCalculate_LongTermCompoundInterest(t *testing.T) {
	calc := NewCalculator(nil)

	params := domain.InvestmentParams{
		InitialCapital:   100000,
		AnnualReturn:     10,
		InvestmentMonths: 120,
		TaxRate:          0,
		InflationRate:    0,
		ReinvestEnabled:  true,
		ReinvestPeriod:   domain.Monthly,
	}

	result := calc.Calculate(nil, params)

	expectedFinal := 100000 * math.Pow(1+0.10/12, 120)
	if !almostEqual(result.FinalBalance, expectedFinal, 10.0) {
		t.Errorf("FinalBalance = %.2f, want %.2f", result.FinalBalance, expectedFinal)
	}
}

func TestCalculate_MonthlyBalanceProgression(t *testing.T) {
	calc := NewCalculator(nil)

	params := domain.InvestmentParams{
		InitialCapital:   100000,
		AnnualReturn:     12,
		InvestmentMonths: 3,
		TaxRate:          0,
		InflationRate:    0,
		ReinvestEnabled:  true,
		ReinvestPeriod:   domain.Monthly,
	}

	result := calc.Calculate(nil, params)

	expectedMonth1 := 100000 * (1 + 0.12/12)
	expectedMonth2 := expectedMonth1 * (1 + 0.12/12)
	expectedMonth3 := expectedMonth2 * (1 + 0.12/12)

	if !almostEqual(result.MonthlyResults[0].EndBalance, expectedMonth1, epsilon) {
		t.Errorf("Month 1 EndBalance = %.2f, want %.2f", result.MonthlyResults[0].EndBalance, expectedMonth1)
	}
	if !almostEqual(result.MonthlyResults[1].EndBalance, expectedMonth2, epsilon) {
		t.Errorf("Month 2 EndBalance = %.2f, want %.2f", result.MonthlyResults[1].EndBalance, expectedMonth2)
	}
	if !almostEqual(result.MonthlyResults[2].EndBalance, expectedMonth3, epsilon) {
		t.Errorf("Month 3 EndBalance = %.2f, want %.2f", result.MonthlyResults[2].EndBalance, expectedMonth3)
	}
}

func TestCalculate_TaxAppliedAnnually(t *testing.T) {
	calc := NewCalculator(nil)

	params := domain.InvestmentParams{
		InitialCapital:   100000,
		AnnualReturn:     12,
		InvestmentMonths: 24,
		TaxRate:          13,
		InflationRate:    0,
		ReinvestEnabled:  false,
	}

	result := calc.Calculate(nil, params)

	taxMonths := 0
	for _, mr := range result.MonthlyResults {
		if mr.TaxPaid > 0 {
			taxMonths++
		}
	}

	if taxMonths != 2 {
		t.Errorf("Tax should be applied 2 times (at month 12 and 24), but was applied %d times", taxMonths)
	}
}
