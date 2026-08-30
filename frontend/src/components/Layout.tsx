import { NavLink, Outlet } from 'react-router-dom'
import { useAuth } from '../auth/AuthContext'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'

function NavItem({ to, end, children }: { to: string; end?: boolean; children: React.ReactNode }) {
  return (
    <NavLink to={to} end={end}>
      {({ isActive }) => (
        <Button
          variant={isActive ? 'secondary' : 'ghost'}
          className="h-9"
        >
          {children}
        </Button>
      )}
    </NavLink>
  )
}

export function Layout() {
  const { user, logout } = useAuth()

  return (
    <div className="min-h-dvh flex flex-col">
      <header className="sticky top-0 z-10 border-b bg-background">
        <div className="mx-auto flex h-14 w-full max-w-4xl items-center gap-4 px-4">
          <span className="font-semibold">Go Clean Template</span>
          <nav className="flex flex-1 items-center gap-1">
            <NavItem to="/" end>
              任务
            </NavItem>
            <NavItem to="/translation">翻译</NavItem>
          </nav>
          <Separator orientation="vertical" className="h-6" />
          <span className="text-sm text-muted-foreground">{user?.username}</span>
          <Button variant="outline" size="sm" onClick={logout}>
            退出
          </Button>
        </div>
      </header>
      <main className="mx-auto w-full max-w-4xl flex-1 px-4 py-6">
        <Outlet />
      </main>
    </div>
  )
}
