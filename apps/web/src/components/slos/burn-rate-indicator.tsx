import { cn } from '@/src/lib/utils'

interface BurnRateProps {
  rate: number
  window: string
}

function getBurnRateColor(rate: number): string {
  if (rate > 14.4) return 'text-red-600 font-bold'
  if (rate > 6)    return 'text-orange-600 font-semibold'
  if (rate > 1)    return 'text-yellow-600'
  return 'text-green-600'
}

export function BurnRateIndicator({ rate, window }: BurnRateProps) {
  if (!rate && rate !== 0) return null

  return (
    <div className="flex justify-between items-center text-xs">
      <span className="text-gray-500">Burn {window}</span>
      <span className={cn('font-mono', getBurnRateColor(rate))}>
        {rate === 0 ? '—' : `${rate.toFixed(2)}x`}
      </span>
    </div>
  )
}
