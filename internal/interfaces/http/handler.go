package http

import (
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"

	"invest-calc/internal/domain"
	"invest-calc/internal/usecases"
)

type Handler struct {
	calc *usecases.Calculator
}

func NewHandler(calc *usecases.Calculator) (*Handler, error) {
	return &Handler{calc: calc}, nil
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", h.serveIndex)
	mux.HandleFunc("POST /api/calculate", h.calculate)
	staticFS, _ := fs.Sub(staticFiles, "static")
	mux.Handle("GET /assets/", http.FileServer(http.FS(staticFS)))
}

func (h *Handler) serveIndex(w http.ResponseWriter, r *http.Request) {
	data, err := staticFiles.ReadFile("static/index.html")
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}

type apiRequest struct {
	InitialCapital   float64 `json:"initial_capital"`
	AnnualReturn     float64 `json:"annual_return"`
	InvestmentYears  int     `json:"investment_years"`
	InvestmentMonths int     `json:"investment_months"`
	TaxRate          float64 `json:"tax_rate"`
	InflationRate    float64 `json:"inflation_rate"`
	ReinvestEnabled  bool    `json:"reinvest_enabled"`
	ReinvestPeriod   string  `json:"reinvest_period"`
	OneTimeContribs  []struct {
		Amount float64 `json:"amount"`
		Month  int     `json:"month"`
		Year   int     `json:"year"`
	} `json:"onetime_contribs"`
	PeriodicContribs []struct {
		Amount     float64 `json:"amount"`
		Period     string  `json:"period"`
		StartMonth int     `json:"start_month"`
	} `json:"periodic_contribs"`
}

func (h *Handler) calculate(w http.ResponseWriter, r *http.Request) {
	var req apiRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}

	investmentMonths := req.InvestmentYears*12 + req.InvestmentMonths

	var contributions []domain.Contribution

	for _, c := range req.OneTimeContribs {
		if c.Amount > 0 && c.Month > 0 && c.Year > 0 {
			contributions = append(contributions, domain.Contribution{
				Type:   domain.OneTime,
				Amount: c.Amount,
				Month:  c.Month + (c.Year-1)*12,
				Year:   c.Year,
			})
		}
	}

	for _, c := range req.PeriodicContribs {
		if c.Amount > 0 && c.Period != "" {
			contributions = append(contributions, domain.Contribution{
				Type:   domain.Periodic,
				Amount: c.Amount,
				Period: domain.Period(c.Period),
				Month:  c.StartMonth,
			})
		}
	}

	reinvestPeriod := req.ReinvestPeriod
	if req.ReinvestEnabled && reinvestPeriod == "" {
		reinvestPeriod = "monthly"
	}

	params := domain.InvestmentParams{
		InitialCapital:   req.InitialCapital,
		AnnualReturn:     req.AnnualReturn,
		InvestmentMonths: investmentMonths,
		TaxRate:          req.TaxRate,
		InflationRate:    req.InflationRate,
		Contributions:    contributions,
		ReinvestEnabled:  req.ReinvestEnabled,
		ReinvestPeriod:   domain.Period(reinvestPeriod),
	}

	result := h.calc.Calculate(r.Context(), params)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

//go:embed static/*
var staticFiles embed.FS
