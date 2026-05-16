export type Severity = 'critical' | 'high' | 'medium' | 'low' | 'info'
export type IncidentStatus = 'open' | 'investigating' | 'identified' | 'monitoring' | 'resolved'

export interface Incident {
  id: string
  title: string
  description: string
  severity: Severity
  status: IncidentStatus
  cluster_id?: string
  service: string
  alert_name: string
  labels: Record<string, string>
  annotations: Record<string, string>
  started_at: string
  resolved_at?: string
  created_at: string
  updated_at: string
}

export interface APIResponse<T> {
  data: T
  error?: string
  total?: number
  page?: number
  per_page?: number
}

export interface ListIncidentsParams {
  status?: IncidentStatus
  severity?: Severity
  service?: string
  page?: number
  per_page?: number
}
