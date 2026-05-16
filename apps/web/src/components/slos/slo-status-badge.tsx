import { cn } from '@/src/lib/utils'
import type { SLOHealthStatus } from '@/src/types'

const statusConfig: Record<SLOHealthStatus, { label: string; className: string; dot: string }> = {
  healthy:  { label: 'Healthy',  className: 'bg-green-100 text-green-700',  dot: 'bg-green-500'  },
  degraded: { label: 'Degraded', className: 'bg-yellow-100 text-yellow-700', dot: 'bg-yellow-500' },
  warning:  { label: 'Warning',  className: 'bg-orange-100 text-orange-700', dot: 'bg-orange-500' },
  critical: { label: 'Critical', className: 'bg-red-100 text-red-700',      dot: 'bg-red-500'    },
  unknown:  { label: 'Unknown',  className: 'bg-gray-100 text-gray-600',    dot: 'bg-gray-400'   },
}

export function SLOStatusBadge({ status }: { status: SLOHealthStatus }) {
  const cfg = statusConfig[status] ?? statusConfig.unknown
  return (
    <span className={cn(
      'inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold',
      cfg.className
    )}>
      <span className={cn('w-1.5 h-1.5 rounded-full', cfg.dot)} />
      {cfg.label}
    </span>
  )
}
