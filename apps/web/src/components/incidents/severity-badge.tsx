import { Badge } from '@/src/components/ui/badge'
import { getSeverityColor } from '@/src/lib/utils'
import type { Severity } from '@/src/types'

export function SeverityBadge({ severity }: { severity: Severity }) {
  return (
    <Badge className={getSeverityColor(severity)}>
      {severity.toUpperCase()}
    </Badge>
  )
}
