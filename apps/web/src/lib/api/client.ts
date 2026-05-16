import axios from 'axios'
import type {
  Incident, APIResponse, ListIncidentsParams,
  SLOStatus
} from '@/src/types'

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:9090',
  headers: { 'Content-Type': 'application/json' },
  timeout: 10000,
})

// ── Incidents ─────────────────────────────────────────────────────────
export const incidentsApi = {
  list: async (params?: ListIncidentsParams): Promise<APIResponse<Incident[]>> => {
    const { data } = await api.get('/api/v1/incidents', { params })
    return data
  },
  getById: async (id: string): Promise<APIResponse<Incident>> => {
    const { data } = await api.get(`/api/v1/incidents/${id}`)
    return data
  },
  create: async (payload: {
    title: string
    severity: string
    service?: string
    description?: string
  }): Promise<APIResponse<Incident>> => {
    const { data } = await api.post('/api/v1/incidents', payload)
    return data
  },
  updateStatus: async (id: string, status: string): Promise<APIResponse<Incident>> => {
    const { data } = await api.patch(`/api/v1/incidents/${id}/status`, { status })
    return data
  },
}

// ── SLOs ──────────────────────────────────────────────────────────────
export const slosApi = {
  list: async (): Promise<APIResponse<SLOStatus[]>> => {
    const { data } = await api.get('/api/v1/slos')
    return data
  },
  getById: async (id: string): Promise<APIResponse<SLOStatus>> => {
    const { data } = await api.get(`/api/v1/slos/${id}`)
    return data
  },
  create: async (payload: {
    name: string
    description?: string
    service: string
    slo_type: string
    target: number
    window_days: number
    promql_good?: string
    promql_total?: string
  }): Promise<APIResponse<SLOStatus>> => {
    const { data } = await api.post('/api/v1/slos', payload)
    return data
  },
}

export default api
