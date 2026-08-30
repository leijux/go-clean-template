import { useEffect, useState, type FormEvent } from 'react'
import { ApiError, translationApi } from '@/api/client'
import type { Translation } from '@/api/types'
import { Alert, AlertTitle } from '@/components/ui/alert'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Empty, EmptyDescription, EmptyTitle } from '@/components/ui/empty'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Separator } from '@/components/ui/separator'
import { Skeleton } from '@/components/ui/skeleton'
import { Spinner } from '@/components/ui/spinner'
import { Textarea } from '@/components/ui/textarea'

export function TranslationPage() {
  const [source, setSource] = useState('auto')
  const [destination, setDestination] = useState('en')
  const [original, setOriginal] = useState('')
  const [result, setResult] = useState<Translation | null>(null)
  const [history, setHistory] = useState<Translation[]>([])
  const [error, setError] = useState<string | null>(null)
  const [submitting, setSubmitting] = useState(false)
  const [loadingHistory, setLoadingHistory] = useState(true)

  async function loadHistory() {
    setLoadingHistory(true)
    try {
      const data = await translationApi.history()
      setHistory(data.history)
    } catch (err) {
      setError(err instanceof ApiError ? err.message : '加载历史失败')
    } finally {
      setLoadingHistory(false)
    }
  }

  useEffect(() => {
    void loadHistory()
  }, [])

  async function onSubmit(e: FormEvent) {
    e.preventDefault()
    if (!original.trim()) return
    setError(null)
    setSubmitting(true)
    try {
      const t = await translationApi.translate({
        source: source.trim(),
        destination: destination.trim(),
        original: original.trim(),
      })
      setResult(t)
      void loadHistory()
    } catch (err) {
      setError(err instanceof ApiError ? err.message : '翻译失败')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <div className="flex flex-col gap-4">
      <h1 className="text-2xl font-semibold">翻译</h1>

      {error && (
        <Alert variant="destructive">
          <AlertTitle>{error}</AlertTitle>
        </Alert>
      )}

      <Card>
        <CardHeader>
          <CardTitle>翻译文本</CardTitle>
        </CardHeader>
        <CardContent>
          <form onSubmit={onSubmit} className="flex flex-col gap-4">
            <div className="flex flex-col gap-3 sm:flex-row">
              <div className="flex flex-1 flex-col gap-2">
                <Label htmlFor="source">源语言</Label>
                <Input
                  id="source"
                  value={source}
                  onChange={(e) => setSource(e.target.value)}
                  placeholder="auto"
                />
              </div>
              <div className="flex flex-1 flex-col gap-2">
                <Label htmlFor="dest">目标语言</Label>
                <Input
                  id="dest"
                  value={destination}
                  onChange={(e) => setDestination(e.target.value)}
                  placeholder="en"
                />
              </div>
            </div>
            <div className="flex flex-col gap-2">
              <Label htmlFor="original">原文</Label>
              <Textarea
                id="original"
                value={original}
                onChange={(e) => setOriginal(e.target.value)}
                placeholder="输入要翻译的文本"
                className="min-h-28"
              />
            </div>
            <Button type="submit" disabled={submitting} className="w-fit">
              {submitting && <Spinner data-icon="inline-start" />}
              {submitting ? '翻译中…' : '翻译'}
            </Button>
          </form>
        </CardContent>
      </Card>

      {result && (
        <Card>
          <CardHeader>
            <CardTitle>翻译结果</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <Badge variant="secondary">{result.source}</Badge>
              <span>→</span>
              <Badge variant="secondary">{result.destination}</Badge>
            </div>
            <Separator className="my-3" />
            <p className="text-base">{result.translation}</p>
          </CardContent>
        </Card>
      )}

      <Card>
        <CardHeader>
          <CardTitle>历史记录</CardTitle>
        </CardHeader>
        <CardContent>
          {loadingHistory ? (
            <div className="flex flex-col gap-3">
              {Array.from({ length: 2 }).map((_, i) => (
                <Skeleton key={i} className="h-16 w-full rounded-lg" />
              ))}
            </div>
          ) : history.length === 0 ? (
            <Empty className="py-10">
              <EmptyTitle>暂无翻译记录</EmptyTitle>
              <EmptyDescription>翻译的内容会显示在这里。</EmptyDescription>
            </Empty>
          ) : (
            <ul className="flex flex-col divide-y divide-border">
              {history.map((item, i) => (
                <li key={i} className="flex flex-col gap-1 py-3 first:pt-0 last:pb-0">
                  <span className="text-xs text-muted-foreground">
                    {item.source} → {item.destination}
                  </span>
                  <span className="text-sm">{item.original}</span>
                  <span className="text-sm font-medium text-primary">
                    {item.translation}
                  </span>
                </li>
              ))}
            </ul>
          )}
        </CardContent>
      </Card>
    </div>
  )
}
