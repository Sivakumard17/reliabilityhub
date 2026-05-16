'use client'

import { formatRelativeTime } from '@/src/lib/utils'
import { SeverityBadge } from './severity-badge'
import { StatusBadge } from './status-badge'
import type { Incident } from '@/src/types'

interface IncidentRowProps {
  incident: Incident
  onStatusUpdate: (id: string, status: string) => void
}

export function IncidentRow({ incident, onStatusUpdate }: IncidentRowProps) {
  return (
    <tr className="hover:bg-gray-50 transition-colors">
      <td className="px-6 py-4">
        <div className="flex flex-col">
          <span className="text-sm font-medium text-gray-900">{incident.title}</span>
          {incident.service && (
            <span className="text-xs text-gray-500 mt-0.5">{incident.service}</span>
          )}
        </div>
      </td>
      <td className="px-6 py-4">
        <SeverityBadge severity={incident.severity} />
      </td>
      <td className="px-6 py-4">
        <StatusBadge status={incident.status} />
      </td>
      <td className="px-6 py-4 text-sm text-gray-500">
        {formatRelativeTime(incident.started_at)}
      </td>
      <td className="px-6 py-4">
        {incident.status !== 'resolved' && (
          <select
            className="text-xs border border-gray-200 rounded px-2 py-1 bg-white"
            defaultValue=""
            onChange={(e) => {
              if (e.target.value) onStatusUpdate(incident.id, e.target.value)
            }}
          >
            <option value="" disabled>Update...</option>
            <option value="investigating">Investigating</option>
            <option value="identified">Identified</option>
            <option value="monitoring">Monitoring</option>
            <option value="resolved">Resolved</option>
          </select>
        )}
      </td>
    </tr>
  )
}
