import type {
  CreateTaskRequest,
  Envelope,
  LoginRequest,
  RegisterRequest,
  Task,
  TaskList,
  TaskStatus,
  TokenResponse,
  TranslateRequest,
  Translation,
  TranslationHistory,
  TransitionTaskRequest,
  UpdateTaskRequest,
  User,
} from './types'

const TOKEN_KEY = 'access_token'
const BASE = '/v1'

export class ApiError extends Error {
  status: number

  constructor(message: string, status: number) {
    super(message)
    this.name = 'ApiError'
    this.status = status
  }
}

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

export function setToken(token: string): void {
  localStorage.setItem(TOKEN_KEY, token)
}

export function clearToken(): void {
  localStorage.removeItem(TOKEN_KEY)
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string> | undefined),
  }

  const token = getToken()
  if (token) {
    headers.Authorization = `Bearer ${token}`
  }

  const res = await fetch(`${BASE}${path}`, { ...options, headers })

  const isJson = res.headers.get('content-type')?.includes('application/json')

  // 后端统一信封 {code, message, data}。非 2xx 时 data 为 null，错误信息在 message。
  if (!res.ok) {
    const body = isJson ? await res.json() : null
    const message =
      body && typeof body === 'object' && 'message' in body
        ? (body as Envelope<unknown>).message
        : `Request failed (${res.status})`
    throw new ApiError(message, res.status)
  }

  if (!isJson) {
    return undefined as T
  }

  const envelope = (await res.json()) as Envelope<T>
  return envelope.data
}

// ---- Auth ----
export const authApi = {
  register: (body: RegisterRequest) =>
    request<User>('/auth/register', { method: 'POST', body: JSON.stringify(body) }),
  login: (body: LoginRequest) =>
    request<TokenResponse>('/auth/login', { method: 'POST', body: JSON.stringify(body) }),
  profile: () => request<User>('/user/profile'),
}

// ---- Tasks ----
export const tasksApi = {
  list: (params: { status?: TaskStatus; limit?: number; offset?: number } = {}) => {
    const qs = new URLSearchParams()
    if (params.status) qs.set('status', params.status)
    if (params.limit != null) qs.set('limit', String(params.limit))
    if (params.offset != null) qs.set('offset', String(params.offset))
    const q = qs.toString()
    return request<TaskList>(`/tasks${q ? `?${q}` : ''}`)
  },
  get: (id: string) => request<Task>(`/tasks/${id}`),
  create: (body: CreateTaskRequest) =>
    request<Task>('/tasks', { method: 'POST', body: JSON.stringify(body) }),
  update: (id: string, body: UpdateTaskRequest) =>
    request<Task>(`/tasks/${id}`, { method: 'PUT', body: JSON.stringify(body) }),
  transition: (id: string, status: TransitionTaskRequest) =>
    request<Task>(`/tasks/${id}/status`, { method: 'PATCH', body: JSON.stringify(status) }),
  remove: (id: string) => request<void>(`/tasks/${id}`, { method: 'DELETE' }),
}

// ---- Translation ----
export const translationApi = {
  history: () => request<TranslationHistory>('/translation/history'),
  translate: (body: TranslateRequest) =>
    request<Translation>('/translation/do-translate', { method: 'POST', body: JSON.stringify(body) }),
}
