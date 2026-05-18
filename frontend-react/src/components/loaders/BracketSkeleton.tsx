export default function BracketSkeleton() {
  return (
    <div className="animate-pulse space-y-4 py-4">
      <div className="h-5 w-40 rounded bg-zinc-700/50" />
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
        {Array.from({ length: 8 }).map((_, i) => (
          <div key={i} className="rounded-xl border border-zinc-800 bg-zinc-900/60 p-3 space-y-2">
            <div className="flex items-center gap-2">
              <div className="w-6 h-6 rounded-full bg-zinc-700/50" />
              <div className="h-3 w-20 rounded bg-zinc-700/50" />
              <div className="h-4 w-4 rounded bg-zinc-700/50 ml-auto" />
            </div>
            <div className="h-px bg-zinc-700/30" />
            <div className="flex items-center gap-2">
              <div className="w-6 h-6 rounded-full bg-zinc-700/50" />
              <div className="h-3 w-24 rounded bg-zinc-700/50" />
              <div className="h-4 w-4 rounded bg-zinc-700/50 ml-auto" />
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
