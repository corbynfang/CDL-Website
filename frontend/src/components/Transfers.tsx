import { useApi } from "../hooks/useApi";
import { Link } from "react-router-dom";
import type { PlayerTransfer } from "../types";

interface TransfersResponse {
  transfers: PlayerTransfer[];
  count: number;
  timestamp?: number;
}

const Transfers = () => {
  const { data: response, loading, error } = useApi<TransfersResponse>(
    "/api/v1/transfers"
  );
  
  const transfers = response?.transfers || [];

  // Format date for display
  const formatDate = (dateString: string) => {
    try {
      const date = new Date(dateString);
      return date.toLocaleDateString("en-US", {
        year: "numeric",
        month: "long",
        day: "numeric",
      });
    } catch {
      return dateString;
    }
  };

  // Get transfer type badge color
  const getTransferTypeColor = (type: string) => {
    switch (type.toLowerCase()) {
      case "signing":
        return "bg-black/10 text-black";
      case "release":
        return "bg-[#6B7280]/20 text-[#6B7280]";
      case "trade":
        return "bg-black/10 text-black";
      default:
        return "bg-gray-200 text-black";
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
          <p className="text-[#6B7280]">Loading transfers...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
          <p className="text-[#555555]">Error: {error}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-white">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
        <h1 className="text-4xl font-bold mb-8 pb-4 text-black">TRANSFERS</h1>

        <div className="space-y-4">
          {transfers && transfers.length > 0 ? (
            transfers.map((transfer) => (
              <div
                key={transfer.id}
                className="bg-[#F4F4F5] rounded-2xl p-6 shadow-md shadow-[rgba(0,0,0,0.1)] hover:shadow-lg transition-all hover:scale-[1.01]"
              >
                <div className="flex flex-col md:flex-row md:items-center md:justify-between gap-4">
                  {/* Left: Player Info */}
                  <div className="flex items-center space-x-4 flex-1">
                    {transfer.player?.avatar_url && (
                      <img
                        src={transfer.player.avatar_url}
                        alt={transfer.player.gamertag}
                        className="w-16 h-16 rounded-xl object-cover shadow-sm"
                      />
                    )}
                    <div className="flex-1">
                      <div className="flex items-center gap-3 mb-2">
                        <Link
                          to={`/players/${transfer.player_id}`}
                          className="text-xl font-bold text-black hover:text-[#555555] transition-colors"
                        >
                          {transfer.player?.gamertag || "Unknown Player"}
                        </Link>
                        <span
                          className={`px-3 py-1 rounded-full text-xs font-medium ${getTransferTypeColor(
                            transfer.transfer_type
                          )}`}
                        >
                          {transfer.transfer_type.toUpperCase()}
                        </span>
                      </div>

                      {/* Transfer Path */}
                      <div className="flex items-center gap-2 text-sm">
                        <span className="text-[#6B7280]">
                          {transfer.from_team?.name || "Free Agent"}
                        </span>
                        <span className="text-black">→</span>
                        <span className="font-semibold text-black">
                          {transfer.to_team?.name || "Free Agent"}
                        </span>
                      </div>
                    </div>
                  </div>

                  {/* Right: Details */}
                  <div className="flex flex-col md:items-end gap-2 text-sm">
                    <div className="flex items-center gap-4 text-[#6B7280]">
                      {transfer.role && (
                        <span className="uppercase tracking-wider">
                          {transfer.role}
                        </span>
                      )}
                      {transfer.season && (
                        <span className="uppercase tracking-wider">
                          {transfer.season}
                        </span>
                      )}
                    </div>
                    <p className="text-[#6B7280]">
                      {formatDate(transfer.transfer_date)}
                    </p>
                  </div>
                </div>

                {/* Description if available */}
                {transfer.description && (
                  <div className="mt-4 pt-4 border-t border-gray-300">
                    <p className="text-sm text-[#555555]">
                      {transfer.description}
                    </p>
                  </div>
                )}
              </div>
            ))
          ) : (
            <div className="bg-[#F4F4F5] rounded-2xl p-8 text-center">
              <p className="text-[#6B7280]">No transfers available</p>
            </div>
          )}
        </div>

        <p className="mt-8 text-[#6B7280] text-sm">
          Total Transfers: {transfers?.length || 0}
        </p>
      </div>
    </div>
  );
};

export default Transfers;

