'use client'

import { useState, useEffect } from 'react'
import { AlertTriangle, RefreshCw } from 'lucide-react'
import { incidentsApi } from '@/src/lib/api/client'
import { IncidentRow } from '@/src/components/incidents/incident-row'
import { Card, CardHeader, CardContent } from '@/src/components/ui/card'
import type { Incident, IncidentStatus, Severity } from '@/src/types'

export default function IncidentsPage() {
  const [incidents, setIncidents] = useState<Incident[]>([])
  const [total, setTotal]         = useState(0)
  const [loading, setLoading]     = useState(true)
  const [error, setError]         = useState<string | null>(null)
  const [statusFilter, setStatus] = useState<IncidentStatus | ''>('')
  const [severityFilter, setSeverity] = useState<Severity | ''>('')

  const fetchIncidents = async () => {
    setLoading(true)
    setError(null)
    try {
      const res = await incidentsApi.list({
        status:   statusFilter   || undefined,
        severity: severityFilter || undefined,
      })
      setIncidents(res.data ?? [])
      setTotal(res.total ?? 0)
    } catch {
      setError('Failed to fetch incidents. Is the API running on port 9090?')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => { fetchIncidents() }, [statusFilter, severityFilter])

  const handleStatusUpdate = async (id: string, status: string) => {
    try {
      await incidentsApi.updateStatus(id, status)
      fetchIncidents()
    } catch {
      alert('Failed to update status')
    }
  }

  const openCount     = incidents.filter(i => i.status === 'open').length
  const criticalCount = incidents.filter(i => i.severity === 'critical').length

  return (
    <div className="p-6 space-y-6">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <AlertTriangle className="h-6 w-6 text-red-500" />
          <h1 className="text-2xl font-bold text-gray-900">Incidents</h1>
          <span className="bg-gray-100 text-gray-700 text-sm px-2.5 py-0.5 rounded-full font-medium">
            {total} total
          </span>
        </div>
        <button
          onClick={fetchIncidents}
          className="flex items-center gap-2 px-4 py-2 text-sm bg-white border border-gray-200 rounded-lg hover:bg-gray-50"
        >
          <RefreshCw className="h-4 w-4" />
          Refresh
        </button>
      </div>

      <div className="grid grid-cols-3 gap-4">
        <Card><CardContent className="py-3">
          <p className="text-sm text-gray-500">Open</p>
          <p className="text-2xl font-bold text-red-600">{openCount}</p>
        </CardContent></Card>
        <Card><CardContent className="py-3">
          <p className="text-sm text-gray-500">Critical</p>
          <p className="text-2xl font-bold text-orange-600">{criticalCount}</p>
        </CardContent></Card>
        <Card><CardContent className="py-3">
          <p className="text-sm text-gray-500">Total</p>
          <p className="text-2xl font-bold text-gray-900">{total}</p>
        </CardContent></Card>
      </div>

      <div className="flex gap-3">
        <select
          className="text-sm border border-gray-200 rounded-lg px-3 py-2 bg-white"
          value={statusFilter}
          onChange={e => setStatus(e.target.value as IncidentStatus | '')}
        >
          <option value="">All Statuses</option>
          <option value="open">Open</option>
          <option value="investigating">Investigating</option>
          <option value="resolved">Resolved</option>
        </select>
        <select
          className="text-sm border border-gray-200 rounded-lg px-3 py-2 bg-white"
          value={severityFilter}
          onChange={e => setSeverity(e.target.value as Severity | '')}
        >
          <option value="">All Severities</option>
          <option value="critical">Critical</option>
          <option value="high">High</option>
          <option value="medium">Medium</option>
          <option value="low">Low</option>
        </select>
      </div>

      <Card>
        <CardHeader>
          <h2 className="text-sm font-semibold text-gray-700">Active Incidents</h2>
        </CardHeader>
        {error ? (
          <CardContent><p className="text-sm text-red-500">{error}</p></CardContent>
        ) : loading ? (
          <CardContent><p className="text-sm text-gray-500 animate-pulse">Loading...</p></CardContent>
        ) : incidents.length === 0 ? (
          <CardContent><p className="text-sm text-gray-500">No incidents found.</p></CardContent>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-gray-200 bg-gray-50">
                  <th className="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Title</th>
                  <th className="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Severity</th>
                  <th className="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Status</th>
                  <th className="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Started</th>
                  <th className="px-6 py-3 text-left text-xs font-semibold text-gray-500 uppercase">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-100">
                {incidents.map(incident => (
                  <IncidentRow
                    key={incident.id}
                    incident={incident}
                    onStatusUpdate={handleStatusUpdate}
                  />
                ))}
              </tbody>
            </table>
          </div>
        )}
      </Card>
    </div>
  )
}
