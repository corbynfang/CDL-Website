import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import type { Team, Player } from '../types';
import { teamApi } from '../services/api';
import TeamLogo from './TeamLogo';
import PlayerAvatar from './PlayerAvatar';

interface TeamWithPlayers extends Team {
  players?: Player[];
}

const Teams: React.FC = () => {
  const [teams, setTeams] = useState<TeamWithPlayers[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchTeamsWithPlayers = async () => {
      try {
        setLoading(true);
        const teamsData = await teamApi.getTeams();
        
        // Fetch players for each team
        const teamsWithPlayers = await Promise.all(
          teamsData.map(async (team) => {
            try {
              const players = await teamApi.getTeamPlayers(team.id);
              return { ...team, players };
            } catch (err) {
              console.error(`Failed to fetch players for team ${team.id}:`, err);
              return { ...team, players: [] };
            }
          })
        );
        
        setTeams(teamsWithPlayers);
      } catch (err) {
        setError('Failed to fetch teams');
        console.error('Error fetching teams:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchTeamsWithPlayers();
  }, []);

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-white"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-center py-8 px-4">
        <div className="text-red-500 text-lg sm:text-xl mb-4">{error}</div>
        <button
          onClick={() => window.location.reload()}
          className="btn-primary"
        >
          TRY AGAIN
        </button>
      </div>
    );
  }

  return (
    <div className="space-y-6 sm:space-y-8 px-4 sm:px-6 lg:px-8">
      {/* Hero Section */}
      <div className="text-center py-8 sm:py-12 lg:py-16 bg-black border border-gray-800 mb-6 sm:mb-8 rounded-lg">
        <h1 className="text-2xl sm:text-3xl lg:text-4xl xl:text-5xl font-black text-white mb-4 sm:mb-6 tracking-wider">
          CDL TEAMS
        </h1>
        <p className="text-sm sm:text-base lg:text-lg text-gray-400 max-w-4xl mx-auto leading-relaxed uppercase tracking-wider px-4 sm:px-6">
          PROFESSIONAL CALL OF DUTY LEAGUE ORGANIZATIONS COMPETING AT THE HIGHEST LEVEL.
          EXPLORE TEAM ROSTERS, PLAYER PROFILES, AND COMPREHENSIVE STATISTICS.
        </p>
        <div className="mt-4 sm:mt-6 lg:mt-8 text-gray-400 text-sm sm:text-base lg:text-lg uppercase tracking-wider">
          {teams.length} ACTIVE TEAMS
        </div>
      </div>

      {/* Teams Grid */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4 sm:gap-6 lg:gap-8">
        {teams.map((team) => (
          <div key={team.id} className="group">
            <div className="card hover:border-white transition-all duration-300 transform hover:scale-105 border border-gray-800 overflow-hidden h-full flex flex-col">
              {/* Team Header */}
              <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-4 sm:mb-6 p-4 sm:p-6">
                <div className="flex flex-col sm:flex-row sm:items-center space-y-3 sm:space-y-0 sm:space-x-3 lg:space-x-4 mb-3 sm:mb-0">
                  {/* Team Logo - Centered on mobile */}
                  <div className="flex justify-center sm:justify-start">
                    <TeamLogo team={team} size="lg" className="flex-shrink-0" />
                  </div>
                  <div className="text-center sm:text-left">
                    <h3 className="text-lg sm:text-xl lg:text-2xl font-bold text-white uppercase tracking-wider leading-tight">
                      {team.name}
                    </h3>
                    <p className="text-gray-400 uppercase tracking-wider text-sm sm:text-base mt-1">
                      {team.city}
                    </p>
                  </div>
                </div>
                <div className="text-center sm:text-right">
                  <span className="inline-block bg-white text-black px-3 sm:px-4 py-1 sm:py-2 rounded-none text-xs sm:text-sm font-bold uppercase tracking-wider">
                    {team.abbreviation}
                  </span>
                </div>
              </div>

              {/* Team Players */}
              <div className="px-4 sm:px-6 pb-4 sm:pb-6 flex-1 flex flex-col">
                <h4 className="text-base sm:text-lg font-semibold text-white mb-3 sm:mb-4 uppercase tracking-wider text-center sm:text-left">
                  ACTIVE ROSTER
                </h4>
                <div className="space-y-2 sm:space-y-3 flex-1">
                  {team.players && team.players.length > 0 ? (
                    team.players.map((player) => (
                      <div key={player.id} className="flex items-center space-x-2 sm:space-x-3 p-2 sm:p-3 bg-gray-750 rounded-lg hover:bg-gray-700 transition-colors duration-200">
                        <div className="flex-shrink-0">
                          <PlayerAvatar player={player} size="sm" className="w-8 h-8 sm:w-10 sm:h-10" />
                        </div>
                        <div className="flex-1 min-w-0">
                          <Link
                            to={`/players/${player.id}`}
                            className="text-white hover:text-red-500 font-medium transition-colors duration-200 text-sm sm:text-base truncate block"
                          >
                            {player.gamertag}
                          </Link>
                        </div>
                      </div>
                    ))
                  ) : (
                    <div className="text-center py-8">
                      <p className="text-gray-400 text-sm sm:text-base">No players available</p>
                    </div>
                  )}
                </div>

                {/* View Team Button */}
                <div className="mt-4 sm:mt-6 pt-4 border-t border-gray-700">
                  <Link
                    to={`/teams/${team.id}`}
                    className="btn-secondary w-full text-center text-sm sm:text-base py-2 sm:py-3 block"
                  >
                    VIEW TEAM
                  </Link>
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default Teams; 