import { useCallback, useEffect, useState } from 'react'
import { ApiError, tasksApi } from '@/api/client'
import type { Task, TaskStatus } from '@/api/types'
import { Alert, AlertTitle } from '@/components/ui/alert'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Empty, EmptyDescription, EmptyTitle } from '@/components/ui/empty'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Separator } from '@/components/ui/separator'
import { Skeleton } from '@/components/ui/skeleton'
import { Spinner } from '@/components/ui/spinner'
import { Textarea } from '@/components/ui/textarea'

const STATUS_LABEL: Record<TaskStatus, string> = {
  todo: '待办',
  in_progress: '进行中',
  done: '已完成',
}

const NEXT_STATUS: Record<TaskStatus, TaskStatus | null> = {
  todo: 'in_progress',
  in_progress: 'done',
  done: null,
}

// 精简的中文标签用于按钮文案
const nextLabel = (status: TaskStatus): string => {
  const next = NEXT_STATUS[status]
  return next ? STATUS_LABEL[next] : ''
}

interface TaskForm {
  title: string
  description: string
}

const emptyForm: TaskForm = { title: '', description: '' }

export function TasksPage() {
  const [tasks, setTasks] = useState<Task[]>([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [actionError, setActionError] = useState<string | null>(null)
  const [saving, setSaving] = useState(false)

  const [createOpen, setCreateOpen] = useState(false)
  const [createForm, setCreateForm] = useState<TaskForm>(emptyForm)

  const [editing, setEditing] = useState<Task | null>(null)
  const [editForm, setEditForm] = useState<TaskForm>(emptyForm)

  const [deleting, setDeleting] = useState<string | null>(null)

  async function load() {
    setLoading(true)
    setError(null)
    try {
      const data = await tasksApi.list({ limit: 100 })
      setTasks(data.tasks)
      setTotal(data.total)
    } catch (err) {
      setError(err instanceof ApiError ? err.message : '加载任务失败')
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    void load()
  }, [])

  function handleActionError(err: unknown) {
    setActionError(err instanceof ApiError ? err.message : '操作失败')
  }

  async function submitCreate() {
    const title = createForm.title.trim()
    if (!title) {
      setActionError('标题不能为空')
      return
    }
    setActionError(null)
    setSaving(true)
    try {
      await tasksApi.create({ title, description: createForm.description.trim() })
      setCreateForm(emptyForm)
      setCreateOpen(false)
      void load()
    } catch (err) {
      handleActionError(err)
    } finally {
      setSaving(false)
    }
  }

  function openEdit(task: Task) {
    setEditing(task)
    setEditForm({ title: task.title, description: task.description })
  }

  async function submitEdit() {
    if (!editing) return
    const title = editForm.title.trim()
    if (!title) {
      setActionError('标题不能为空')
      return
    }
    setActionError(null)
    setSaving(true)
    try {
      await tasksApi.update(editing.id, {
        title,
        description: editForm.description.trim(),
      })
      setEditing(null)
      void load()
    } catch (err) {
      handleActionError(err)
    } finally {
      setSaving(false)
    }
  }

  async function transition(task: Task, status: TaskStatus) {
    setActionError(null)
    try {
      await tasksApi.transition(task.id, { status })
      void load()
    } catch (err) {
      handleActionError(err)
    }
  }

  const deleteTask = useCallback(async (task: Task) => {
    if (!window.confirm(`确认删除任务「${task.title}」？`)) return
    setDeleting(task.id)
    setActionError(null)
    try {
      await tasksApi.remove(task.id)
      void load()
    } catch (err) {
      handleActionError(err)
    } finally {
      setDeleting(null)
    }
  }, [])

  return (
    <div className="flex flex-col gap-4">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold">我的任务</h1>
        <Button onClick={() => setCreateOpen(true)}>
          新建任务
        </Button>
      </div>

      {actionError && (
        <Alert variant="destructive">
          <AlertTitle>{actionError}</AlertTitle>
        </Alert>
      )}

      {loading ? (
        <div className="flex flex-col gap-3">
          {Array.from({ length: 3 }).map((_, i) => (
            <Skeleton key={i} className="h-24 w-full rounded-xl" />
          ))}
        </div>
      ) : error ? (
        <Alert variant="destructive">
          <AlertTitle>{error}</AlertTitle>
        </Alert>
      ) : total === 0 ? (
        <Empty className="py-16">
          <EmptyTitle>暂无任务</EmptyTitle>
          <EmptyDescription>
            点击右上角「新建任务」创建第一个任务。
          </EmptyDescription>
        </Empty>
      ) : (
        <ul className="flex flex-col gap-3">
          {tasks.map((task) => (
            <li key={task.id}>
              <Card>
                <CardHeader>
                  <div className="flex items-start justify-between gap-4">
                    <CardTitle className="text-base">{task.title}</CardTitle>
                    <Badge variant="secondary">{STATUS_LABEL[task.status]}</Badge>
                  </div>
                </CardHeader>
                <CardContent className="flex flex-col gap-3">
                  {task.description && (
                    <p className="whitespace-pre-wrap text-sm text-muted-foreground">
                      {task.description}
                    </p>
                  )}
                  <p className="text-xs text-muted-foreground">
                    创建于 {new Date(task.created_at).toLocaleString()}
                  </p>
                  <Separator />
                  <div className="flex flex-wrap gap-2">
                    {nextLabel(task.status) && (
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() =>
                          transition(task, NEXT_STATUS[task.status]!)
                        }
                      >
                        推进为「{nextLabel(task.status)}」
                      </Button>
                    )}
                    {task.status === 'in_progress' && (
                      <Button
                        variant="outline"
                        size="sm"
                        onClick={() => transition(task, 'todo')}
                      >
                        回退为「待办」
                      </Button>
                    )}
                    <Button variant="ghost" size="sm" onClick={() => openEdit(task)}>
                      编辑
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      className="text-destructive hover:text-destructive"
                      disabled={deleting === task.id}
                      onClick={() => void deleteTask(task)}
                    >
                      {deleting === task.id && <Spinner data-icon="inline-start" />}
                      删除
                    </Button>
                  </div>
                </CardContent>
              </Card>
            </li>
          ))}
        </ul>
      )}

      {/* 新建任务 */}
      <Dialog open={createOpen} onOpenChange={setCreateOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>新建任务</DialogTitle>
            <DialogDescription>填写任务的标题和可选描述。</DialogDescription>
          </DialogHeader>
          <div className="flex flex-col gap-4">
            <div className="flex flex-col gap-2">
              <Label htmlFor="create-title">标题</Label>
              <Input
                id="create-title"
                value={createForm.title}
                maxLength={255}
                onChange={(e) =>
                  setCreateForm((f) => ({ ...f, title: e.target.value }))
                }
                placeholder="任务标题"
              />
            </div>
            <div className="flex flex-col gap-2">
              <Label htmlFor="create-desc">描述（可选）</Label>
              <Textarea
                id="create-desc"
                value={createForm.description}
                maxLength={1000}
                onChange={(e) =>
                  setCreateForm((f) => ({ ...f, description: e.target.value }))
                }
                placeholder="任务描述"
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setCreateOpen(false)}>
              取消
            </Button>
            <Button onClick={submitCreate} disabled={saving}>
              {saving && <Spinner data-icon="inline-start" />}
              创建
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      {/* 编辑任务 */}
      <Dialog open={Boolean(editing)} onOpenChange={(o) => !o && setEditing(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>编辑任务</DialogTitle>
            <DialogDescription>修改任务的标题和描述。</DialogDescription>
          </DialogHeader>
          <div className="flex flex-col gap-4">
            <div className="flex flex-col gap-2">
              <Label htmlFor="edit-title">标题</Label>
              <Input
                id="edit-title"
                value={editForm.title}
                maxLength={255}
                onChange={(e) =>
                  setEditForm((f) => ({ ...f, title: e.target.value }))
                }
                placeholder="任务标题"
              />
            </div>
            <div className="flex flex-col gap-2">
              <Label htmlFor="edit-desc">描述（可选）</Label>
              <Textarea
                id="edit-desc"
                value={editForm.description}
                maxLength={1000}
                onChange={(e) =>
                  setEditForm((f) => ({ ...f, description: e.target.value }))
                }
                placeholder="任务描述"
              />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setEditing(null)}>
              取消
            </Button>
            <Button onClick={submitEdit} disabled={saving}>
              {saving && <Spinner data-icon="inline-start" />}
              保存
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
