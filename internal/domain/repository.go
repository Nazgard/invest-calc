package domain

import "context"

type CalculatorRepository interface {
	SaveCalculation(ctx context.Context, result CalculationResult) (string, error)
	GetCalculation(ctx context.Context, id string) (*CalculationResult, error)
	ListCalculations(ctx context.Context) ([]CalculationResult, error)
}
