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
        <div className="animate-spin rounded-none h-12 w-12 border-b-2 border-white"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-center py-8">
        <div className="text-red-500 text-xl mb-4">{error}</div>
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
    <div className="space-y-8">
      {/* Hero Section */}
      <div className="text-center py-16 bg-black border border-gray-800 mb-8">
        <h1 className="text-heading text-white mb-6">CDL TEAMS</h1>
        <p className="text-subheading text-gray-400 max-w-3xl mx-auto leading-relaxed uppercase tracking-wider">
          PROFESSIONAL CALL OF DUTY LEAGUE ORGANIZATIONS COMPETING AT THE HIGHEST LEVEL.
          EXPLORE TEAM ROSTERS, PLAYER PROFILES, AND COMPREHENSIVE STATISTICS.
        </p>
        <div className="mt-8 text-gray-400 text-lg uppercase tracking-wider">
          {teams.length} ACTIVE TEAMS
        </div>
      </div>

      {/* Teams Grid */}
      <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
        {teams.map((team) => (
          <div key={team.id} className="group">
            <div className="card hover:border-white transition-all duration-300 transform hover:scale-105 border border-gray-800 overflow-hidden">
              {/* Team Header */}
              <div className="flex items-center justify-between mb-6 p-6">
                <div className="flex items-center space-x-4">
                  {/* Team Logo */}
                  <TeamLogo team={team} size="xl" />
                  <div>
                    <h3 className="text-2xl font-bold text-white uppercase tracking-wider">{team.name}</h3>
                    <p className="text-gray-400 uppercase tracking-wider">{team.city}</p>
                  </div>
                </div>
                <div className="text-right">
                  <span className="inline-block bg-white text-black px-4 py-2 rounded-none text-sm font-bold uppercase tracking-wider">
                    {team.abbreviation}
                  </span>
                </div>
              </div>

              {/* Team Players */}
              <div className="px-6 pb-6">
                <h4 className="text-lg font-semibold text-white mb-4 uppercase tracking-wider">ACTIVE ROSTER</h4>
                <div className="space-y-3">
                  {team.players && team.players.length > 0 ? (
                    team.players.map((player) => (
                      <div key={player.id} className="flex items-center space-x-3">
                        <PlayerAvatar player={player} size="sm" />
                        <div className="flex-1">
                          <Link
                            to={`/players/${player.id}`}
                            className="text-white hover:text-red-500 font-medium transition-colors duration-200"
                          >
                            {player.gamertag}
                          </Link>
                        </div>
                      </div>
                    ))
                  ) : (
                    <p className="text-gray-400 text-sm">No players available</p>
                  )}
                </div>

                {/* View Team Button */}
                <div className="mt-6">
                  <Link
                    to={`/teams/${team.id}`}
                    className="btn-secondary w-full text-center"
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