export function formatMoney(value: number): string {
  const negative = value < 0
  const abs = Math.abs(value)
  const rounded = Math.round(abs * 100) / 100
  const intPart = Math.floor(rounded)
  const fracPart = Math.round((rounded - intPart) * 100)
  const intStr = intPart.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ' ')
  const result = intStr + '.' + fracPart.toString().padStart(2, '0')
  return negative ? '-' + result : result
}

export function formatInput(value: string): string {
  const cleaned = value.replace(/[^\d.]/g, '')
  if (!cleaned) return ''
  const parts = cleaned.split('.')
  parts[0] = parts[0].replace(/\B(?=(\d{3})+(?!\d))/g, ' ')
  return parts.join('.')
}

export function parseMoney(value: string): string {
  return value.replace(/\s/g, '')
}

export function defaultFormParams() {
  return {
    initialCapital: '100 000',
    annualReturn: '10',
    investmentYears: '5',
    investmentMonths: '0',
    taxRate: '13',
    inflationRate: '7',
    reinvestEnabled: false,
    reinvestPeriod: 'monthly' as const,
    oneTimeContribs: [] as { id: number; amount: string; month: string; year: string }[],
    periodicContribs: [] as { id: number; amount: string; period: 'monthly' | 'quarterly' | 'yearly'; startMonth: string }[]
  }
}

let nextId = 1
export function genId(): number {
  return nextId++
}
