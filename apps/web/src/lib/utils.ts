import { clsx, type ClassValue } from 'clsx'
import { twMerge } from 'tailwind-merge'
import { formatDistanceToNow, format } from 'date-fns'
import type { Severity, IncidentStatus } from '@/src/types'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function formatRelativeTime(dateString: string): string {
  return formatDistanceToNow(new Date(dateString), { addSuffix: true })
}

export function formatDateTime(dateString: string): string {
  return format(new Date(dateString), 'MMM dd, yyyy HH:mm')
}

export function getSeverityColor(severity: Severity): string {
  const colors = {
    critical: 'bg-red-100 text-red-800 border-red-200',
    high:     'bg-orange-100 text-orange-800 border-orange-200',
    medium:   'bg-yellow-100 text-yellow-800 border-yellow-200',
    low:      'bg-blue-100 text-blue-800 border-blue-200',
    info:     'bg-gray-100 text-gray-800 border-gray-200',
  }
  return colors[severity] ?? colors.info
}

export function getStatusColor(status: IncidentStatus): string {
  const colors = {
    open:          'bg-red-100 text-red-700',
    investigating: 'bg-orange-100 text-orange-700',
    identified:    'bg-yellow-100 text-yellow-700',
    monitoring:    'bg-blue-100 text-blue-700',
    resolved:      'bg-green-100 text-green-700',
  }
  return colors[status] ?? 'bg-gray-100 text-gray-700'
}
