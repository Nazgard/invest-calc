import { useState, useCallback } from 'react'
import type { FormParams, OneTimeContrib, PeriodicContrib } from '../types'
import { formatInput, parseMoney } from '../utils'
import ContribSection from './ContribSection'

interface InvestFormProps {
  onSubmit: (params: FormParams) => void
  loading: boolean
}

const defaultParams = (): FormParams => ({
  initialCapital: '100 000',
  annualReturn: '10',
  investmentYears: '5',
  investmentMonths: '0',
  taxRate: '13',
  inflationRate: '7',
  reinvestEnabled: false,
  reinvestPeriod: 'monthly',
  oneTimeContribs: [],
  periodicContribs: []
})

export default function InvestForm({ onSubmit, loading }: InvestFormProps) {
  const [form, setForm] = useState<FormParams>(defaultParams)

  const update = useCallback((field: keyof FormParams, value: unknown) => {
    setForm(prev => ({ ...prev, [field]: value }))
  }, [])

  const handleMoneyInput = useCallback((field: keyof FormParams, raw: string) => {
    const cleaned = raw.replace(/[^\d.]/g, '')
    const parts = cleaned.split('.')
    if (parts.length > 2) parts.length = 2
    parts[0] = parts[0].replace(/\B(?=(\d{3})+(?!\d))/g, ' ')
    setForm(prev => ({ ...prev, [field]: parts.join('.') }))
  }, [])

  const handleSubmit = useCallback((e: React.FormEvent) => {
    e.preventDefault()
    onSubmit(form)
  }, [form, onSubmit])

  return (
    <form onSubmit={handleSubmit}>
      <div className="card">
        <h2>Параметры инвестирования</h2>

        <div className="form-row">
          <div className="form-group">
            <label>Стартовый капитал</label>
            <div className="input-with-icon">
              <input
                type="text"
                value={form.initialCapital}
                onChange={e => handleMoneyInput('initialCapital', e.target.value)}
                className="money-input"
              />
              <span className="input-icon">&#8381;</span>
            </div>
          </div>
          <div className="form-group">
            <label>Годовая доходность</label>
            <div className="input-with-icon">
              <input
                type="number"
                value={form.annualReturn}
                onChange={e => update('annualReturn', e.target.value)}
                step="any"
                placeholder="10"
              />
              <span className="input-icon">%</span>
            </div>
          </div>
        </div>

        <div className="form-row">
          <div className="form-group">
            <label>Срок (лет)</label>
            <input
              type="number"
              value={form.investmentYears}
              onChange={e => update('investmentYears', e.target.value)}
              min={0} max={50}
            />
          </div>
          <div className="form-group">
            <label>Срок (месяцев)</label>
            <input
              type="number"
              value={form.investmentMonths}
              onChange={e => update('investmentMonths', e.target.value)}
              min={0} max={11}
            />
          </div>
        </div>

        <div className="divider" />

        <div className="form-row">
          <div className="form-group">
            <label>Налог на доход</label>
            <div className="input-with-icon">
              <input
                type="number"
                value={form.taxRate}
                onChange={e => update('taxRate', e.target.value)}
                step="any"
                placeholder="13"
              />
              <span className="input-icon">%</span>
            </div>
          </div>
          <div className="form-group">
            <label>Инфляция</label>
            <div className="input-with-icon">
              <input
                type="number"
                value={form.inflationRate}
                onChange={e => update('inflationRate', e.target.value)}
                step="any"
                placeholder="7"
              />
              <span className="input-icon">%</span>
            </div>
          </div>
        </div>

        <div className="divider" />

        <div className="form-row">
          <div className="form-group">
            <div className="checkbox-group">
              <input
                type="checkbox"
                id="reinvest_enabled"
                checked={form.reinvestEnabled}
                onChange={e => update('reinvestEnabled', e.target.checked)}
              />
              <label htmlFor="reinvest_enabled">Реинвестировать доход</label>
            </div>
          </div>
        </div>

        {form.reinvestEnabled && (
          <div>
            <label style={{ fontSize: '0.82rem', color: '#a0a0a0', fontWeight: 500 }}>Период реинвестирования</label>
            <div className="reinvest-options">
              {([
                ['daily', 'Каждый день'],
                ['monthly', 'Раз в месяц'],
                ['quarterly', 'Раз в квартал'],
                ['semi_yearly', 'Раз в полгода'],
                ['yearly', 'Раз в год']
              ] as const).map(([val, label]) => (
                <label className="reinvest-option" key={val}>
                  <input
                    type="radio"
                    name="reinvest_period"
                    value={val}
                    checked={form.reinvestPeriod === val}
                    onChange={() => update('reinvestPeriod', val)}
                  />
                  {label}
                </label>
              ))}
            </div>
          </div>
        )}
      </div>

      <div className="card">
        <h2>Дополнительные вложения</h2>

        <ContribSection
          title="Единичные вложения"
          contribs={form.oneTimeContribs}
          type="onetime"
          onChange={c => update('oneTimeContribs', c as OneTimeContrib[])}
        />

        <ContribSection
          title="Периодические вложения"
          contribs={form.periodicContribs}
          type="periodic"
          onChange={c => update('periodicContribs', c as PeriodicContrib[])}
        />
      </div>

      <button type="submit" className="btn btn-calc" disabled={loading}>
        {loading ? 'Расчет...' : 'Рассчитать'}
      </button>
    </form>
  )
}
