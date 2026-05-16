export type Severity = 'critical' | 'high' | 'medium' | 'low' | 'info'
export type IncidentStatus = 'open' | 'investigating' | 'identified' | 'monitoring' | 'resolved'
export type SLOType = 'availability' | 'latency' | 'error_rate' | 'throughput'
export type SLOHealthStatus = 'healthy' | 'degraded' | 'warning' | 'critical' | 'unknown'

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

export interface SLO {
  id: string
  name: string
  description: string
  cluster_id?: string
  service: string
  slo_type: SLOType
  target: number
  window_days: number
  promql_good: string
  promql_total: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface SLOSnapshot {
  id: string
  slo_id: string
  compliance: number
  error_budget_total: number
  error_budget_used: number
  error_budget_pct: number
  burn_rate_1h: number
  burn_rate_6h: number
  burn_rate_24h: number
  snapshot_at: string
}

export interface SLOStatus {
  slo: SLO
  snapshot?: SLOSnapshot
  status: SLOHealthStatus
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
