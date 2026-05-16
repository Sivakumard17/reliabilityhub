'use client'

import { cn } from '@/src/lib/utils'

interface ErrorBudgetBarProps {
  pct: number        // 0-1 remaining
  total: number      // total minutes
  used: number       // used minutes
}

export function ErrorBudgetBar({ pct, total, used }: ErrorBudgetBarProps) {
  const remaining = Math.max(0, Math.min(1, pct))
  const pctDisplay = (remaining * 100).toFixed(1)

  const barColor = remaining > 0.5
    ? 'bg-green-500'
    : remaining > 0.25
    ? 'bg-yellow-500'
    : 'bg-red-500'

  const formatMinutes = (mins: number) => {
    if (mins >= 60 * 24) return `${(mins / 60 / 24).toFixed(1)}d`
    if (mins >= 60)      return `${(mins / 60).toFixed(1)}h`
    return `${mins.toFixed(0)}m`
  }

  return (
    <div className="space-y-1.5">
      <div className="flex justify-between text-xs text-gray-500">
        <span>Error Budget</span>
        <span className="font-medium">{pctDisplay}% remaining</span>
      </div>
      <div className="w-full bg-gray-100 rounded-full h-2.5">
        <div
          className={cn('h-2.5 rounded-full transition-all', barColor)}
          style={{ width: `${remaining * 100}%` }}
        />
      </div>
      <div className="flex justify-between text-xs text-gray-400">
        <span>Used: {formatMinutes(used)}</span>
        <span>Total: {formatMinutes(total)}</span>
      </div>
    </div>
  )
}
