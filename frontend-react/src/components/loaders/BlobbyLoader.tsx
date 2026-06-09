interface BlobbyLoaderProps {
  label?: string;
  color?: string;
}

export default function BlobbyLoader({
  label = "Loading...",
  color = "#78BE20",
}: BlobbyLoaderProps) {
  return (
    <div className="flex flex-col items-center justify-center gap-6 py-20">
      <div className="relative w-20 h-20">
        <div
          className="absolute inset-0"
          style={{
            background: `radial-gradient(circle at 40% 40%, ${color}cc, ${color}44)`,
            animation:
              "blob-morph 4s ease-in-out infinite, blob-pulse 2s ease-in-out infinite",
            boxShadow: `0 0 32px ${color}55`,
          }}
        />
      </div>
      <p className="text-sm text-zinc-400 tracking-widest uppercase">{label}</p>
    </div>
  );
}
