export interface OneTimeContrib {
  id: number
  amount: string
  month: string
  year: string
}

export interface PeriodicContrib {
  id: number
  amount: string
  period: 'monthly' | 'quarterly' | 'yearly'
  startMonth: string
}

export interface FormParams {
  initialCapital: string
  annualReturn: string
  investmentYears: string
  investmentMonths: string
  taxRate: string
  inflationRate: string
  reinvestEnabled: boolean
  reinvestPeriod: 'daily' | 'monthly' | 'quarterly' | 'semi_yearly' | 'yearly'
  oneTimeContribs: OneTimeContrib[]
  periodicContribs: PeriodicContrib[]
}

export interface MonthResult {
  month: number
  startBalance: number
  contributions: number
  returnEarned: number
  taxPaid: number
  endBalance: number
  endBalanceReal: number
  cumulativeReturn: number
  cumulativeContrib: number
}

export interface CalculationResult {
  params: {
    initialCapital: number
    annualReturn: number
    investmentMonths: number
    taxRate: number
    inflationRate: number
    reinvestEnabled: boolean
    reinvestPeriod: string
  }
  monthlyResults: MonthResult[]
  finalBalance: number
  finalBalanceReal: number
  totalContrib: number
  totalReturn: number
  totalTax: number
  netProfit: number
  netProfitReal: number
}

export interface ApiRequest {
  initial_capital: number
  annual_return: number
  investment_years: number
  investment_months: number
  tax_rate: number
  inflation_rate: number
  reinvest_enabled: boolean
  reinvest_period: string
  onetime_contribs: { amount: number; month: number; year: number }[]
  periodic_contribs: { amount: number; period: string; start_month: number }[]
}
