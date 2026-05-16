'use client'

import { Target } from 'lucide-react'
import { Card, CardContent, CardHeader } from '@/src/components/ui/card'
import { SLOStatusBadge } from './slo-status-badge'
import { ErrorBudgetBar } from './error-budget-bar'
import { BurnRateIndicator } from './burn-rate-indicator'
import type { SLOStatus } from '@/src/types'

interface SLOCardProps {
  sloStatus: SLOStatus
}

export function SLOCard({ sloStatus }: SLOCardProps) {
  const { slo, snapshot, status } = sloStatus
  const targetPct = (slo.target * 100).toFixed(2)

  return (
    <Card className="hover:shadow-md transition-shadow">
      <CardHeader className="pb-3">
        <div className="flex items-start justify-between gap-2">
          <div className="flex items-center gap-2">
            <div className="p-1.5 bg-blue-50 rounded">
              <Target className="h-4 w-4 text-blue-600" />
            </div>
            <div>
              <h3 className="text-sm font-semibold text-gray-900 leading-tight">
                {slo.name}
              </h3>
              <p className="text-xs text-gray-500 mt-0.5">{slo.service}</p>
            </div>
          </div>
          <SLOStatusBadge status={status} />
        </div>
      </CardHeader>

      <CardContent className="space-y-4">
        {/* SLO Details */}
        <div className="grid grid-cols-2 gap-2 text-xs">
          <div className="bg-gray-50 rounded p-2">
            <p className="text-gray-500">Target</p>
            <p className="font-semibold text-gray-900">{targetPct}%</p>
          </div>
          <div className="bg-gray-50 rounded p-2">
            <p className="text-gray-500">Window</p>
            <p className="font-semibold text-gray-900">{slo.window_days}d</p>
          </div>
          <div className="bg-gray-50 rounded p-2">
            <p className="text-gray-500">Type</p>
            <p className="font-semibold text-gray-900 capitalize">
              {slo.slo_type.replace('_', ' ')}
            </p>
          </div>
          <div className="bg-gray-50 rounded p-2">
            <p className="text-gray-500">Compliance</p>
            <p className="font-semibold text-gray-900">
              {snapshot
                ? `${(snapshot.compliance * 100).toFixed(3)}%`
                : '—'}
            </p>
          </div>
        </div>

        {/* Error Budget Bar */}
        {snapshot ? (
          <ErrorBudgetBar
            pct={snapshot.error_budget_pct}
            total={snapshot.error_budget_total}
            used={snapshot.error_budget_used}
          />
        ) : (
          <div className="space-y-1.5">
            <div className="flex justify-between text-xs text-gray-500">
              <span>Error Budget</span>
              <span>No data yet</span>
            </div>
            <div className="w-full bg-gray-100 rounded-full h-2.5" />
          </div>
        )}

        {/* Burn Rates */}
        {snapshot && (
          <div className="space-y-1 pt-1 border-t border-gray-100">
            <BurnRateIndicator rate={snapshot.burn_rate_1h}  window="1h" />
            <BurnRateIndicator rate={snapshot.burn_rate_6h}  window="6h" />
            <BurnRateIndicator rate={snapshot.burn_rate_24h} window="24h" />
          </div>
        )}
      </CardContent>
    </Card>
  )
}
