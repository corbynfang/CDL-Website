export default function MatchCardSkeleton() {
  return (
    <div className="rounded-xl border border-zinc-800 bg-zinc-900/60 p-4 animate-pulse">
      <div className="flex items-center justify-between gap-4">
        <div className="flex items-center gap-3 flex-1">
          <div className="w-8 h-8 rounded-full bg-zinc-700/50" />
          <div className="h-4 w-28 rounded bg-zinc-700/50" />
        </div>
        <div className="flex items-center gap-2">
          <div className="h-6 w-6 rounded bg-zinc-700/50" />
          <div className="h-4 w-4 rounded bg-zinc-700/50" />
          <div className="h-6 w-6 rounded bg-zinc-700/50" />
        </div>
        <div className="flex items-center gap-3 flex-1 justify-end">
          <div className="h-4 w-28 rounded bg-zinc-700/50" />
          <div className="w-8 h-8 rounded-full bg-zinc-700/50" />
        </div>
      </div>
    </div>
  )
}
