import type { ApiRequest, CalculationResult } from './types'

export async function calculate(params: ApiRequest): Promise<CalculationResult> {
  const res = await fetch('/api/calculate', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(params)
  })
  if (!res.ok) throw new Error('Calculation failed')
  return res.json()
}
