import axios from 'axios'
import type { Incident, APIResponse, ListIncidentsParams } from '@/src/types'

const api = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:9090',
  headers: { 'Content-Type': 'application/json' },
  timeout: 10000,
})

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

export default api
