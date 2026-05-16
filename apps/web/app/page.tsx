'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { AlertTriangle, Activity, Target, Server } from 'lucide-react'
import { incidentsApi } from '@/src/lib/api/client'
import { slosApi } from '@/src/lib/api/client'
import { Card, CardContent } from '@/src/components/ui/card'

export const dynamic = 'force-dynamic'

export default function DashboardPage() {
  const [incidentStats, setIncidentStats] = useState({ total: 0, open: 0, critical: 0 })
  const [sloStats, setSloStats]           = useState({ total: 0, healthy: 0, atRisk: 0 })
  const [loading, setLoading]             = useState(true)

  useEffect(() => {
    Promise.all([
      incidentsApi.list({ per_page: 100 }),
      slosApi.list(),
    ]).then(([incRes, sloRes]) => {
      const incidents = incRes.data ?? []
      setIncidentStats({
        total:    incRes.total ?? 0,
        open:     incidents.filter(i => i.status === 'open').length,
        critical: incidents.filter(i => i.severity === 'critical').length,
      })

      const slos = sloRes.data ?? []
      setSloStats({
        total:   sloRes.total ?? 0,
        healthy: slos.filter(s => s.status === 'healthy').length,
        atRisk:  slos.filter(s => ['critical','warning','degraded'].includes(s.status)).length,
      })
    }).finally(() => setLoading(false))
  }, [])

  const statCards = [
    {
      label: 'Total Incidents',
      value: incidentStats.total,
      sub: `${incidentStats.open} open`,
      icon: AlertTriangle,
      color: 'text-red-500',
      bg: 'bg-red-50',
      href: '/incidents',
    },
    {
      label: 'Active SLOs',
      value: sloStats.total,
      sub: `${sloStats.healthy} healthy`,
      icon: Target,
      color: 'text-blue-500',
      bg: 'bg-blue-50',
      href: '/slos',
    },
    {
      label: 'Critical Incidents',
      value: incidentStats.critical,
      sub: 'last 30 days',
      icon: Activity,
      color: 'text-orange-500',
      bg: 'bg-orange-50',
      href: '/incidents',
    },
    {
      label: 'SLOs at Risk',
      value: sloStats.atRisk,
      sub: 'burning budget',
      icon: Server,
      color: 'text-purple-500',
      bg: 'bg-purple-50',
      href: '/slos',
    },
  ]

  return (
    <div className="p-6 space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">ReliabilityHub</h1>
        <p className="text-sm text-gray-500 mt-1">SRE Control Plane for Kubernetes</p>
      </div>

      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        {statCards.map(({ label, value, sub, icon: Icon, color, bg, href }) => (
          <Link key={label} href={href}>
            <Card className="hover:shadow-md transition-shadow cursor-pointer">
              <CardContent className="py-4">
                <div className={`inline-flex p-2 rounded-lg ${bg} mb-3`}>
                  <Icon className={`h-5 w-5 ${color}`} />
                </div>
                <p className="text-2xl font-bold text-gray-900">
                  {loading ? '...' : value}
                </p>
                <p className="text-sm text-gray-700 mt-0.5">{label}</p>
                <p className="text-xs text-gray-400 mt-0.5">{sub}</p>
              </CardContent>
            </Card>
          </Link>
        ))}
      </div>

      <div className="grid grid-cols-2 gap-4">
        <Card>
          <CardContent className="py-4">
            <h2 className="text-sm font-semibold text-gray-700 mb-3">Quick Actions</h2>
            <div className="flex flex-col gap-2">
              <Link href="/incidents"
                className="text-sm px-4 py-2 bg-red-50 text-red-700 rounded-lg hover:bg-red-100 transition-colors">
                View Incidents →
              </Link>
              <Link href="/slos"
                className="text-sm px-4 py-2 bg-blue-50 text-blue-700 rounded-lg hover:bg-blue-100 transition-colors">
                View SLO Dashboard →
              </Link>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="py-4">
            <h2 className="text-sm font-semibold text-gray-700 mb-3">System Status</h2>
            <div className="space-y-2">
              <div className="flex justify-between text-sm">
                <span className="text-gray-500">API</span>
                <span className="text-green-600 font-medium">● Healthy</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-gray-500">Database</span>
                <span className="text-green-600 font-medium">● Healthy</span>
              </div>
              <div className="flex justify-between text-sm">
                <span className="text-gray-500">Cluster</span>
                <span className="text-green-600 font-medium">● Running</span>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
