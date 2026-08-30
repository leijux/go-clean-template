import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useState,
  type ReactNode,
} from 'react'
import {
  authApi,
  clearToken,
  getToken,
  setToken,
  ApiError,
} from '../api/client'
import type { User } from '../api/types'

interface AuthContextValue {
  user: User | null
  loading: boolean
  login: (email: string, password: string) => Promise<User>
  register: (username: string, email: string, password: string) => Promise<User>
  logout: () => void
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined)

async function postLogin(token: string): Promise<User> {
  setToken(token)
  // 401 时清 token 并抛出，避免把过期的 user 缓存在页面里
  try {
    return await authApi.profile()
  } catch (err) {
    if (err instanceof ApiError && err.status === 401) {
      clearToken()
      throw err
    }
    throw err
  }
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState<boolean>(() => Boolean(getToken()))

  useEffect(() => {
    if (!getToken()) {
      setLoading(false)
      return
    }
    authApi
      .profile()
      .then(setUser)
      .catch(() => clearToken())
      .finally(() => setLoading(false))
  }, [])

  const login = useCallback(async (email: string, password: string) => {
    const { token } = await authApi.login({ email, password })
    const u = await postLogin(token)
    setUser(u)
    return u
  }, [])

  const register = useCallback(
    async (username: string, email: string, password: string) => {
      // 注册接口不返回 token，注册成功后直接触发登录
      await authApi.register({ username, email, password })
      const { token } = await authApi.login({ email, password })
      const u = await postLogin(token)
      setUser(u)
      return u
    },
    [],
  )

  const logout = useCallback(() => {
    clearToken()
    setUser(null)
  }, [])

  return (
    <AuthContext.Provider value={{ user, loading, login, register, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext)
  if (!ctx) {
    throw new Error('useAuth must be used within AuthProvider')
  }
  return ctx
}
