import { useApi } from "../hooks/useApi";
import { Link } from "react-router-dom";
import type { PlayerTransfer } from "../types";

interface TransfersResponse {
  transfers: PlayerTransfer[];
  count: number;
  timestamp?: number;
}

const transferBadge = (type: string) => {
  switch (type.toLowerCase()) {
    case "signing":
      return "text-green-400 bg-green-400/10";
    case "release":
      return "text-[#737373] bg-white/5";
    case "trade":
      return "text-blue-400 bg-blue-400/10";
    default:
      return "text-[#a3a3a3] bg-white/5";
  }
};

const formatDate = (dateString: string) => {
  try {
    return new Date(dateString).toLocaleDateString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
    });
  } catch {
    return dateString;
  }
};

const Transfers = () => {
  const { data: response, loading, error } = useApi<TransfersResponse>("/api/v1/transfers");

  const transfers = response?.transfers || [];

  if (loading) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <p className="text-[#737373] text-sm">Loading transfers...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <p className="text-[#737373] text-sm">Error: {error}</p>
      </div>
    );
  }

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div className="mb-8">
        <p className="text-xs uppercase tracking-widest text-[#737373] mb-2">Roster Moves</p>
        <h1 className="font-grotesk text-3xl font-bold text-white">TRANSFERS</h1>
      </div>

      <div className="space-y-2">
        {transfers.length > 0 ? (
          transfers.map((transfer) => (
            <div
              key={transfer.id}
              className="bg-[#111111] border border-[#1a1a1a] p-5 hover:border-[#2a2a2a] transition-colors"
            >
              <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
                <div className="flex items-center gap-4">
                  {transfer.player?.avatar_url && (
                    <img
                      src={transfer.player.avatar_url}
                      alt={transfer.player.gamertag}
                      className="w-12 h-12 object-cover rounded-full opacity-90 flex-shrink-0"
                    />
                  )}
                  <div>
                    <div className="flex items-center gap-2 mb-1.5">
                      <Link
                        to={`/players/${transfer.player_id}`}
                        className="font-grotesk font-semibold text-white hover:text-[#a3a3a3] transition-colors text-sm"
                      >
                        {transfer.player?.gamertag || "Unknown Player"}
                      </Link>
                      <span
                        className={`px-2 py-0.5 text-xs font-medium uppercase tracking-wider ${transferBadge(transfer.transfer_type)}`}
                      >
                        {transfer.transfer_type}
                      </span>
                    </div>

                    <div className="flex items-center gap-2 text-xs">
                      <span className="text-[#737373]">
                        {transfer.from_team?.name || "Free Agent"}
                      </span>
                      <span className="text-[#737373]">→</span>
                      <span className="text-[#a3a3a3] font-medium">
                        {transfer.to_team?.name || "Free Agent"}
                      </span>
                    </div>
                  </div>
                </div>

                <div className="flex items-center gap-4 text-xs text-[#737373] md:flex-col md:items-end md:gap-1">
                  {transfer.role && (
                    <span className="uppercase tracking-wider">{transfer.role}</span>
                  )}
                  {transfer.season && (
                    <span className="uppercase tracking-wider">{transfer.season}</span>
                  )}
                  <span>{formatDate(transfer.transfer_date)}</span>
                </div>
              </div>

              {transfer.description && (
                <div className="mt-4 pt-4 border-t border-[#1a1a1a]">
                  <p className="text-xs text-[#737373] leading-relaxed">
                    {transfer.description}
                  </p>
                </div>
              )}
            </div>
          ))
        ) : (
          <div className="border border-[#1a1a1a] p-16 text-center">
            <p className="text-[#737373] text-sm">No transfers available</p>
          </div>
        )}
      </div>

      <p className="mt-4 text-[#737373] text-xs">
        {transfers.length} transfers
      </p>
    </div>
  );
};

export default Transfers;
