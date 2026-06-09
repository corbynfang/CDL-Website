export default function EventCardSkeleton() {
  return (
    <div className="rounded-2xl border border-zinc-800 bg-zinc-900/60 overflow-hidden animate-pulse">
      <div className="h-2 bg-zinc-700/50 w-full" />
      <div className="p-5 space-y-3">
        <div className="flex items-center gap-2">
          <div className="h-3 w-16 rounded bg-zinc-700/50" />
          <div className="h-3 w-10 rounded bg-zinc-700/50" />
        </div>
        <div className="h-5 w-3/4 rounded bg-zinc-700/50" />
        <div className="h-4 w-1/2 rounded bg-zinc-700/50" />
        <div className="flex gap-2 pt-1">
          <div className="h-3 w-20 rounded bg-zinc-700/50" />
          <div className="h-3 w-16 rounded bg-zinc-700/50" />
        </div>
      </div>
    </div>
  );
}
