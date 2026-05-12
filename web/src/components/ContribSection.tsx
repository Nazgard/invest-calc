import { useCallback } from 'react'
import type { OneTimeContrib, PeriodicContrib } from '../types'
import { formatInput, genId } from '../utils'

interface ContribSectionProps {
  title: string
  contribs: OneTimeContrib[] | PeriodicContrib[]
  type: 'onetime' | 'periodic'
  onChange: (contribs: OneTimeContrib[] | PeriodicContrib[]) => void
}

export default function ContribSection({ title, contribs, type, onChange }: ContribSectionProps) {
  const addContrib = useCallback(() => {
    if (type === 'onetime') {
      const newContrib: OneTimeContrib = { id: genId(), amount: '', month: '', year: '' }
      onChange([...(contribs as OneTimeContrib[]), newContrib])
    } else {
      const newContrib: PeriodicContrib = { id: genId(), amount: '', period: 'monthly', startMonth: '' }
      onChange([...(contribs as PeriodicContrib[]), newContrib])
    }
  }, [type, contribs, onChange])

  const removeContrib = useCallback((id: number) => {
    if (type === 'onetime') {
      onChange((contribs as OneTimeContrib[]).filter(c => c.id !== id))
    } else {
      onChange((contribs as PeriodicContrib[]).filter(c => c.id !== id))
    }
  }, [type, contribs, onChange])

  const updateContrib = useCallback((id: number, field: string, value: string) => {
    if (type === 'onetime') {
      const items = contribs as OneTimeContrib[]
      onChange(items.map(c => c.id === id ? { ...c, [field]: value } : c))
    } else {
      const items = contribs as PeriodicContrib[]
      onChange(items.map(c => c.id === id ? { ...c, [field]: value } : c))
    }
  }, [type, contribs, onChange])

  return (
    <div className="contrib-section">
      <div className="contrib-section-header">
        <span className="contrib-section-title">{title}</span>
        <button type="button" className="btn btn-sm btn-outline" onClick={addContrib}>+ Добавить</button>
      </div>
      {contribs.map(c => (
        <div className="contrib-row" key={c.id}>
          <div className="form-group">
            <label>Сумма</label>
            <div className="input-with-icon">
              <input
                type="text"
                value={c.amount}
                onChange={e => updateContrib(c.id, 'amount', formatInput(e.target.value))}
                className="money-input"
              />
              <span className="input-icon">&#8381;</span>
            </div>
          </div>
          {type === 'onetime' ? (
            <>
              <div className="form-group small">
                <label>Месяц</label>
                <input
                  type="number"
                  value={(c as OneTimeContrib).month}
                  onChange={e => updateContrib(c.id, 'month', e.target.value)}
                  min={1} max={12}
                />
              </div>
              <div className="form-group small">
                <label>Год</label>
                <input
                  type="number"
                  value={(c as OneTimeContrib).year}
                  onChange={e => updateContrib(c.id, 'year', e.target.value)}
                  min={1}
                />
              </div>
            </>
          ) : (
            <>
              <div className="form-group">
                <label>Период</label>
                <select
                  value={(c as PeriodicContrib).period}
                  onChange={e => updateContrib(c.id, 'period', e.target.value)}
                >
                  <option value="monthly">Ежемесячно</option>
                  <option value="quarterly">Ежеквартально</option>
                  <option value="yearly">Ежегодно</option>
                </select>
              </div>
              <div className="form-group small">
                <label>С месяца</label>
                <input
                  type="number"
                  value={(c as PeriodicContrib).startMonth}
                  onChange={e => updateContrib(c.id, 'startMonth', e.target.value)}
                  min={1}
                />
              </div>
            </>
          )}
          <div className="form-group" style={{ minWidth: 40, flex: 0 }}>
            <label>&nbsp;</label>
            <button type="button" className="btn btn-sm btn-del" onClick={() => removeContrib(c.id)}>&#215;</button>
          </div>
        </div>
      ))}
    </div>
  )
}
