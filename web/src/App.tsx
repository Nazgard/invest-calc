import { useState, useCallback } from 'react'
import './index.css'
import type { FormParams, ApiRequest, CalculationResult } from './types'
import { parseMoney } from './utils'
import InvestForm from './components/InvestForm'
import Results from './components/Results'
import { calculate } from './api'

export default function App() {
  const [result, setResult] = useState<CalculationResult | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleSubmit = useCallback(async (form: FormParams) => {
    setLoading(true)
    setError(null)

    const toNum = (s: string) => parseFloat(parseMoney(s)) || 0

    const req: ApiRequest = {
      initial_capital: toNum(form.initialCapital),
      annual_return: toNum(form.annualReturn),
      investment_years: parseInt(form.investmentYears) || 0,
      investment_months: parseInt(form.investmentMonths) || 0,
      tax_rate: toNum(form.taxRate),
      inflation_rate: toNum(form.inflationRate),
      reinvest_enabled: form.reinvestEnabled,
      reinvest_period: form.reinvestPeriod,
      onetime_contribs: form.oneTimeContribs
        .filter(c => toNum(c.amount) > 0 && parseInt(c.month) > 0 && parseInt(c.year) > 0)
        .map(c => ({ amount: toNum(c.amount), month: parseInt(c.month), year: parseInt(c.year) })),
      periodic_contribs: form.periodicContribs
        .filter(c => toNum(c.amount) > 0 && c.period)
        .map(c => ({ amount: toNum(c.amount), period: c.period, start_month: parseInt(c.startMonth) || 1 }))
    }

    try {
      const res = await calculate(req)
      setResult(res)
      setTimeout(() => {
        document.getElementById('results')?.scrollIntoView({ behavior: 'smooth' })
      }, 100)
    } catch {
      setError('Ошибка при расчете. Попробуйте снова.')
    } finally {
      setLoading(false)
    }
  }, [])

  return (
    <div className="container">
      <h1>Инвестиционный калькулятор</h1>
      <InvestForm onSubmit={handleSubmit} loading={loading} />
      {error && <div className="error">{error}</div>}
      {result && (
        <div id="results" style={{ marginTop: 16 }}>
          <Results result={result} />
        </div>
      )}
    </div>
  )
}
