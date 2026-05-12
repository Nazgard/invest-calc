import { useState, useMemo } from 'react'
import type { CalculationResult, MonthResult } from '../types'
import { formatMoney } from '../utils'

interface PeriodData {
  label: string
  initialCapital: number
  contributions: number
  income: number
  endBalance: number
}

function groupByPeriod(monthlyResults: MonthResult[], totalMonths: number): PeriodData[] {
  const groupByYear = totalMonths > 12
  const originalCapital = monthlyResults.length > 0 ? monthlyResults[0].startBalance : 0

  const groups: { label: string; months: MonthResult[] }[] = []

  if (groupByYear) {
    const yearGroups: Map<number, MonthResult[]> = new Map()
    monthlyResults.forEach(row => {
      const year = Math.ceil(row.month / 12)
      if (!yearGroups.has(year)) yearGroups.set(year, [])
      yearGroups.get(year)!.push(row)
    })
    yearGroups.forEach((months, year) => {
      groups.push({ label: `${year} ${year === 1 ? 'год' : year < 5 ? 'года' : 'лет'}`, months })
    })
  } else {
    monthlyResults.forEach(row => {
      groups.push({ label: `Мес ${row.month}`, months: [row] })
    })
  }

  return groups.map(g => {
    const last = g.months[g.months.length - 1]
    const totalContribs = monthlyResults
      .filter(m => m.month <= last.month)
      .reduce((s, m) => s + m.contributions, 0)
    const totalIncome = last.endBalance - originalCapital - totalContribs

    return {
      label: g.label,
      initialCapital: originalCapital,
      contributions: totalContribs,
      income: totalIncome,
      endBalance: last.endBalance
    }
  })
}

interface ResultsProps {
  result: CalculationResult
}

export default function Results({ result }: ResultsProps) {
  const [activeTab, setActiveTab] = useState<'table' | 'chart' | 'formulas'>('chart')

  const periods = useMemo(
    () => groupByPeriod(result.monthlyResults, result.params.investmentMonths),
    [result]
  )

  const maxBalance = Math.max(...periods.map(p => p.endBalance), 1)
  const maxHeight = 200

  return (
    <div>
      <div className="card">
        <h2>Результаты расчета</h2>
        <div className="summary">
          <div className="summary-item">
            <div className="label">Итоговая сумма</div>
            <div className="value">{formatMoney(result.finalBalance)}</div>
          </div>
          <div className="summary-item">
            <div className="label">С учетом инфляции</div>
            <div className="value">{formatMoney(result.finalBalanceReal)}</div>
          </div>
          <div className="summary-item">
            <div className="label">Всего вложено</div>
            <div className="value">{formatMoney(result.totalContrib)}</div>
          </div>
          <div className="summary-item">
            <div className="label">Чистый доход</div>
            <div className="value positive">{formatMoney(result.netProfit)}</div>
          </div>
          <div className="summary-item">
            <div className="label">Реальный доход</div>
            <div className={result.netProfitReal < 0 ? 'value negative' : 'value positive'}>
              {formatMoney(result.netProfitReal)}
            </div>
          </div>
          <div className="summary-item">
            <div className="label">Налоги</div>
            <div className="value negative">{formatMoney(result.totalTax)}</div>
          </div>
        </div>
      </div>

      <div className="card">
        <h2>Структура капитала по периодам</h2>

        <div className="chart-legend" style={{ display: 'flex', gap: 20, justifyContent: 'center', marginBottom: 16, flexWrap: 'wrap' }}>
          <div style={{ display: 'flex', alignItems: 'center', gap: 6, fontSize: '0.82rem', color: '#ccc' }}>
            <div style={{ width: 14, height: 14, borderRadius: 3, background: '#3b82f6' }} />
            Начальный капитал
          </div>
          <div style={{ display: 'flex', alignItems: 'center', gap: 6, fontSize: '0.82rem', color: '#ccc' }}>
            <div style={{ width: 14, height: 14, borderRadius: 3, background: '#10b981' }} />
            Вложения
          </div>
          <div style={{ display: 'flex', alignItems: 'center', gap: 6, fontSize: '0.82rem', color: '#ccc' }}>
            <div style={{ width: 14, height: 14, borderRadius: 3, background: '#f59e0b' }} />
            Доход
          </div>
        </div>

        <div className="chart-scroll" style={{ overflowX: 'auto', paddingBottom: 8 }}>
          <div style={{ display: 'flex', alignItems: 'flex-end', gap: 8, minWidth: Math.max(periods.length * 80, 400), padding: '0 10px', height: 260 }}>
            {periods.map(p => {
              const initialH = Math.max((p.initialCapital / maxBalance) * maxHeight, 0)
              const contribH = Math.max((p.contributions / maxBalance) * maxHeight, 0)
              const incomeH = Math.max((p.income / maxBalance) * maxHeight, 0)

              return (
                <div key={p.label} style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', flex: 1, minWidth: 60 }}>
                  <div style={{ fontSize: '0.72rem', color: '#a0a0a0', marginBottom: 4, textAlign: 'center' }}>
                    {formatMoney(p.endBalance)}
                  </div>
                  <div style={{ display: 'flex', flexDirection: 'column', width: '100%', maxWidth: 70 }}>
                    {p.initialCapital > 0 && (
                      <div style={{ height: initialH, background: '#3b82f6', borderRadius: '4px 4px 0 0' }} />
                    )}
                    <div style={{ height: contribH, background: '#10b981' }} />
                    <div style={{
                      height: incomeH,
                      background: p.income < 0 ? '#ef4444' : '#f59e0b',
                      borderRadius: '0 0 4px 4px'
                    }} />
                  </div>
                  <div style={{ fontSize: '0.75rem', color: '#a0a0a0', marginTop: 8, textAlign: 'center', whiteSpace: 'nowrap' }}>
                    {p.label}
                  </div>
                </div>
              )
            })}
          </div>
        </div>
      </div>

      <div className="card">
        <div className="tabs">
          <div className={`tab${activeTab === 'chart' ? ' active' : ''}`} onClick={() => setActiveTab('chart')}>График</div>
          <div className={`tab${activeTab === 'table' ? ' active' : ''}`} onClick={() => setActiveTab('table')}>Таблица</div>
          <div className={`tab${activeTab === 'formulas' ? ' active' : ''}`} onClick={() => setActiveTab('formulas')}>Формулы</div>
        </div>

        <div className={`tab-content${activeTab === 'chart' ? ' active' : ''}`}>
          <div className="chart-scroll" style={{ overflowX: 'auto' }}>
            <div style={{ display: 'flex', gap: 16, minWidth: Math.max(periods.length * 200, 400) }}>
              {periods.map(p => (
                <div key={p.label} style={{ flex: 1, minWidth: 180, background: '#252525', borderRadius: 8, padding: 14, border: '1px solid #333' }}>
                  <div style={{ fontSize: '0.9rem', fontWeight: 600, color: '#e0e0e0', marginBottom: 10 }}>{p.label}</div>
                  <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                    {p.initialCapital > 0 && (
                      <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '0.82rem' }}>
                        <span style={{ color: '#a0a0a0' }}>Начальный капитал</span>
                        <span style={{ color: '#3b82f6', fontWeight: 500 }}>{formatMoney(p.initialCapital)}</span>
                      </div>
                    )}
                    <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '0.82rem' }}>
                      <span style={{ color: '#a0a0a0' }}>Вложения</span>
                      <span style={{ color: '#10b981', fontWeight: 500 }}>{formatMoney(p.contributions)}</span>
                    </div>
                    <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '0.82rem' }}>
                      <span style={{ color: '#a0a0a0' }}>Доход</span>
                      <span style={{ color: p.income < 0 ? '#ef4444' : '#f59e0b', fontWeight: 500 }}>{formatMoney(p.income)}</span>
                    </div>
                    <div style={{ height: 1, background: '#333', margin: '4px 0' }} />
                    <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: '0.85rem' }}>
                      <span style={{ color: '#e0e0e0', fontWeight: 600 }}>Итого</span>
                      <span style={{ color: '#ffffff', fontWeight: 700 }}>{formatMoney(p.endBalance)}</span>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>

        <div className={`tab-content${activeTab === 'table' ? ' active' : ''}`}>
          <div className="table-wrap">
            <table>
              <thead>
                <tr>
                  <th>Месяц</th>
                  <th>Начальный баланс</th>
                  <th>Вложения</th>
                  <th>Доход</th>
                  <th>Налог</th>
                  <th>Конечный баланс</th>
                  <th>Реальный баланс</th>
                </tr>
              </thead>
              <tbody>
                {result.monthlyResults.map(row => (
                  <tr key={row.month}>
                    <td>{row.month}</td>
                    <td>{formatMoney(row.startBalance)}</td>
                    <td>{formatMoney(row.contributions)}</td>
                    <td>{formatMoney(row.returnEarned)}</td>
                    <td>{formatMoney(row.taxPaid)}</td>
                    <td>{formatMoney(row.endBalance)}</td>
                    <td>{formatMoney(row.endBalanceReal)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        <div className={`tab-content${activeTab === 'formulas' ? ' active' : ''}`}>
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 16 }}>
            <div style={{ background: '#1e3a5f', border: '1px solid #3b82f6', borderRadius: 8, padding: 14 }}>
              <div style={{ fontSize: '0.88rem', fontWeight: 600, color: '#60a5fa', marginBottom: 8 }}>С реинвестированием</div>
              <div style={{ fontFamily: 'monospace', fontSize: '0.85rem', background: '#2a2a2a', padding: '8px 10px', borderRadius: 6, margin: '6px 0', color: '#e0e0e0' }}>A = P &times; (1 + r/n)^(n&times;t)</div>
              <div style={{ fontSize: '0.8rem', color: '#a0a0a0', marginTop: 4 }}>P &mdash; начальная сумма, r &mdash; ставка, n &mdash; периодов в год, t &mdash; лет</div>
            </div>
            <div style={{ background: '#1a3d2e', border: '1px solid #10b981', borderRadius: 8, padding: 14 }}>
              <div style={{ fontSize: '0.88rem', fontWeight: 600, color: '#34d399', marginBottom: 8 }}>Без реинвестирования</div>
              <div style={{ fontFamily: 'monospace', fontSize: '0.85rem', background: '#2a2a2a', padding: '8px 10px', borderRadius: 6, margin: '6px 0', color: '#e0e0e0' }}>A = P &times; (1 + r &times; t)</div>
              <div style={{ fontSize: '0.8rem', color: '#a0a0a0', marginTop: 4 }}>P &mdash; начальная сумма, r &mdash; ставка, t &mdash; лет</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
