'use client'

import { useState, useEffect } from 'react'
import { Target, RefreshCw, Plus } from 'lucide-react'
import { slosApi } from '@/src/lib/api/client'
import { SLOCard } from '@/src/components/slos/slo-card'
import { Card, CardContent } from '@/src/components/ui/card'
import type { SLOStatus, SLOHealthStatus } from '@/src/types'

export const dynamic = 'force-dynamic'

export default function SLOsPage() {
  const [slos, setSlos]       = useState<SLOStatus[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError]     = useState<string | null>(null)
  const [total, setTotal]     = useState(0)

  const fetchSLOs = async () => {
    setLoading(true)
    setError(null)
    try {
      const res = await slosApi.list()
      setSlos(res.data ?? [])
      setTotal(res.total ?? 0)
    } catch {
      setError('Failed to fetch SLOs. Is the API running?')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchSLOs() }, [])

  const healthyCount  = slos.filter(s => s.status === 'healthy').length
  const criticalCount = slos.filter(s => s.status === 'critical' || s.status === 'warning').length
  const unknownCount  = slos.filter(s => s.status === 'unknown').length

  return (
    <div className="p-6 space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Target className="h-6 w-6 text-blue-500" />
          <h1 className="text-2xl font-bold text-gray-900">SLOs</h1>
          <span className="bg-gray-100 text-gray-700 text-sm px-2.5 py-0.5 rounded-full font-medium">
            {total} total
          </span>
        </div>
        <button
          onClick={fetchSLOs}
          className="flex items-center gap-2 px-4 py-2 text-sm bg-white border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors"
        >
          <RefreshCw className="h-4 w-4" />
          Refresh
        </button>
      </div>

      {/* Summary Stats */}
      <div className="grid grid-cols-4 gap-4">
        <Card>
          <CardContent className="py-3">
            <p className="text-sm text-gray-500">Total SLOs</p>
            <p className="text-2xl font-bold text-gray-900">{total}</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="py-3">
            <p className="text-sm text-gray-500">Healthy</p>
            <p className="text-2xl font-bold text-green-600">{healthyCount}</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="py-3">
            <p className="text-sm text-gray-500">At Risk</p>
            <p className="text-2xl font-bold text-orange-600">{criticalCount}</p>
          </CardContent>
        </Card>
        <Card>
          <CardContent className="py-3">
            <p className="text-sm text-gray-500">Unknown</p>
            <p className="text-2xl font-bold text-gray-500">{unknownCount}</p>
          </CardContent>
        </Card>
      </div>

      {/* SLO Cards Grid */}
      {error ? (
        <Card>
          <CardContent>
            <p className="text-sm text-red-500">{error}</p>
          </CardContent>
        </Card>
      ) : loading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {[1, 2].map(i => (
            <Card key={i} className="animate-pulse">
              <CardContent className="h-48" />
            </Card>
          ))}
        </div>
      ) : slos.length === 0 ? (
        <Card>
          <CardContent className="py-12 text-center">
            <Target className="h-12 w-12 text-gray-300 mx-auto mb-3" />
            <p className="text-gray-500 font-medium">No SLOs defined yet</p>
            <p className="text-sm text-gray-400 mt-1">
              Create SLOs via the API to track service reliability
            </p>
          </CardContent>
        </Card>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {slos.map(sloStatus => (
            <SLOCard
              key={sloStatus.slo.id}
              sloStatus={sloStatus}
            />
          ))}
        </div>
      )}

      {/* Quick Create via API hint */}
      <Card className="border-dashed border-2 border-gray-200">
        <CardContent className="py-4">
          <div className="flex items-center gap-3">
            <Plus className="h-5 w-5 text-gray-400" />
            <div>
              <p className="text-sm font-medium text-gray-700">Add a new SLO</p>
              <p className="text-xs text-gray-400 font-mono mt-0.5">
                POST /api/v1/slos
              </p>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
