import { Link } from "react-router-dom";
import { useApi } from "../hooks/useApi";

interface Team {
  id: number;
  name: string;
  abbreviation: string;
  logo_url: string;
}

const Teams = () => {
  const { data: teams, loading, error } = useApi<Team[]>("/api/v1/teams");

  if (loading) {
    return (
      <div className="min-h-screen bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-16">
          <p className="text-[#6B7280]">Loading teams...</p>
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
        <h1 className="text-4xl font-bold mb-8 pb-4 text-black">TEAMS</h1>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {teams?.map((team) => (
            <Link
              key={team.id}
              to={`/teams/${team.id}`}
              className="p-8 bg-[#F4F4F5] shadow-md shadow-[rgba(0,0,0,0.1)]"
            >
              <div className="flex items-center space-x-4 mb-4">
                {team.logo_url && (
                  <img
                    src={team.logo_url}
                    alt={team.name}
                    className="w-16 h-16 object-contain"
                  />
                )}
                <div>
                  <h2 className="text-xl font-bold text-black">{team.name}</h2>
                  <p className="text-[#6B7280]">{team.abbreviation}</p>
                </div>
              </div>
            </Link>
          ))}
        </div>

        <p className="mt-8 text-[#6B7280] text-sm">
          Total Teams: {teams?.length || 0}
        </p>
      </div>
    </div>
  );
};

export default Teams;
