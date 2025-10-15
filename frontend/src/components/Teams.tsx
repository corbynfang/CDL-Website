import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import type { Team, Player } from '../types';
import { useApi } from '../hooks/useApi';
import { teamApi } from '../services/api';
import TeamLogo from './TeamLogo';
import LoadingSkeleton, { ErrorDisplay } from './LoadingSkeleton';

interface TeamWithPlayers extends Team {
  players?: Player[];
}

const Teams: React.FC = () => {
  const [teamsWithPlayers, setTeamsWithPlayers] = useState<TeamWithPlayers[]>([]);
  const [loadingPlayers, setLoadingPlayers] = useState(false);

  const { data: teams, loading: teamsLoading, error, refetch } = useApi<Team[]>(
    '/api/v1/teams',
    { retries: 3, retryDelay: 1000 }
  );

  useEffect(() => {
    if (!teams || teams.length === 0) return;

    const fetchPlayersForTeams = async () => {
      setLoadingPlayers(true);
      
      const teamsData = await Promise.all(
        teams.map(async (team) => {
          try {
            const players = await teamApi.getTeamPlayers(team.id);
            return { ...team, players };
          } catch (err) {
            return { ...team, players: [] };
          }
        })
      );
      
      setTeamsWithPlayers(teamsData);
      setLoadingPlayers(false);
    };

    fetchPlayersForTeams();
  }, [teams]);

  if (teamsLoading || loadingPlayers) {
    return <LoadingSkeleton variant="card" count={8} />;
  }

  if (error) {
    return <ErrorDisplay message={error} onRetry={refetch} />;
  }

  return (
    <div className="max-w-4xl mx-auto px-4 sm:px-6 py-8">
      <div className="mb-12">
        <h1 className="text-3xl sm:text-4xl font-bold text-white mb-3">teams</h1>
        <p className="text-lg" style={{ color: '#a3a3a3' }}>
          {teamsWithPlayers.length} cdl organizations
        </p>
      </div>

      <div className="space-y-2">
        {teamsWithPlayers.map((team) => (
          <Link
            key={team.id}
            to={`/teams/${team.id}`}
            className="group block"
          >
            <div className="flex items-center justify-between p-4 rounded-lg transition-all duration-200 hover:bg-card">
              <div className="flex items-center space-x-4 flex-1 min-w-0">
                <TeamLogo team={team} size="md" />
                <div className="flex-1 min-w-0">
                  <div className="flex items-center gap-2">
                    <h3 className="text-white font-semibold">{team.name}</h3>
                    <span className="text-xs px-2 py-1 rounded" style={{ backgroundColor: '#1a1a1a', color: '#a3a3a3' }}>
                      {team.abbreviation}
                    </span>
                  </div>
                  <p className="text-sm" style={{ color: '#737373' }}>
                    {team.players && team.players.length > 0 
                      ? `${team.players.map(p => p.gamertag).join(', ')}`
                      : 'no roster'
                    }
                  </p>
                </div>
              </div>
              <svg className="w-5 h-5 text-muted group-hover:translate-x-1 transition-transform flex-shrink-0" style={{ color: '#737373' }} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
              </svg>
            </div>
          </Link>
        ))}
      </div>
    </div>
  );
};

export default Teams; 