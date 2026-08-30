// 与后端 internal/controller/restapi/v1 及 internal/entity 对应的类型

export interface User {
  id: string
  username: string
  email: string
  created_at: string
  updated_at: string
}

export interface TokenResponse {
  token: string
}

export type TaskStatus = 'todo' | 'in_progress' | 'done'

export interface Task {
  id: string
  user_id: string
  title: string
  description: string
  status: TaskStatus
  created_at: string
  updated_at: string
}

export interface TaskList {
  tasks: Task[]
  total: number
}

export interface Translation {
  source: string
  destination: string
  original: string
  translation: string
}

// 对应后端 internal/entity/translation.history.go（GET /translation/history 返回 { history: [...] }）
export interface TranslationHistory {
  history: Translation[]
}

export interface ErrorResponse {
  error: string
}

// 后端统一响应信封 {code, message, data}（对应 internal/controller/restapi/v1/response/envelope.go）
export interface Envelope<T> {
  code: number
  message: string
  data: T
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
}

export interface LoginRequest {
  email: string
  password: string
}

export interface CreateTaskRequest {
  title: string
  description?: string
}

export interface UpdateTaskRequest {
  title: string
  description?: string
}

export interface TransitionTaskRequest {
  status: TaskStatus
}

export interface TranslateRequest {
  source: string
  destination: string
  original: string
}
