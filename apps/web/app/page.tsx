'use client'

import { useState, useEffect } from 'react'
import Link from 'next/link'
import { AlertTriangle, Activity, Shield, Server } from 'lucide-react'
import { incidentsApi } from '@/src/lib/api/client'
import { Card, CardContent } from '@/src/components/ui/card'

export default function DashboardPage() {
  const [stats, setStats] = useState({ total: 0, open: 0, critical: 0 })
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    incidentsApi.list({ per_page: 100 }).then(res => {
      const data = res.data ?? []
      setStats({
        total:    res.total ?? 0,
        open:     data.filter(i => i.status === 'open').length,
        critical: data.filter(i => i.severity === 'critical').length,
      })
    }).finally(() => setLoading(false))
  }, [])

  const statCards = [
    { label: 'Total Incidents', value: stats.total,    icon: AlertTriangle, color: 'text-red-500',    bg: 'bg-red-50',    href: '/incidents' },
    { label: 'Open',            value: stats.open,     icon: Activity,      color: 'text-orange-500', bg: 'bg-orange-50', href: '/incidents' },
    { label: 'Critical',        value: stats.critical, icon: Shield,        color: 'text-purple-500', bg: 'bg-purple-50', href: '/incidents' },
    { label: 'Clusters',        value: 1,              icon: Server,        color: 'text-blue-500',   bg: 'bg-blue-50',   href: '/clusters' },
  ]

  return (
    <div className="p-6 space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-gray-900">ReliabilityHub</h1>
        <p className="text-sm text-gray-500 mt-1">SRE Control Plane for Kubernetes</p>
      </div>

      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        {statCards.map(({ label, value, icon: Icon, color, bg, href }) => (
          <Link key={label} href={href}>
            <Card className="hover:shadow-md transition-shadow cursor-pointer">
              <CardContent className="py-4">
                <div className={`inline-flex p-2 rounded-lg ${bg} mb-3`}>
                  <Icon className={`h-5 w-5 ${color}`} />
                </div>
                <p className="text-2xl font-bold text-gray-900">
                  {loading ? '...' : value}
                </p>
                <p className="text-sm text-gray-500 mt-1">{label}</p>
              </CardContent>
            </Card>
          </Link>
        ))}
      </div>

      <Card>
        <CardContent>
          <h2 className="text-sm font-semibold text-gray-700 mb-3">Quick Links</h2>
          <div className="flex gap-3">
            <Link href="/incidents" className="text-sm px-4 py-2 bg-red-50 text-red-700 rounded-lg hover:bg-red-100">
              View Incidents
            </Link>
            <Link href="/slos" className="text-sm px-4 py-2 bg-blue-50 text-blue-700 rounded-lg hover:bg-blue-100">
              SLO Dashboard
            </Link>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
