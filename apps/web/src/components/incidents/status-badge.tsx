import { cn, getStatusColor } from '@/src/lib/utils'
import type { IncidentStatus } from '@/src/types'

export function StatusBadge({ status }: { status: IncidentStatus }) {
  return (
    <span className={cn(
      'inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium',
      getStatusColor(status)
    )}>
      <span className="w-1.5 h-1.5 rounded-full bg-current" />
      {status.charAt(0).toUpperCase() + status.slice(1)}
    </span>
  )
}
