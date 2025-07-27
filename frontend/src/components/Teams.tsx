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
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-500"></div>
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
          Try Again
        </button>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      {/* Hero Section */}
      <div className="text-center py-16 bg-gradient-to-r from-blue-900 to-purple-900 rounded-lg mb-8">
        <h1 className="text-5xl font-bold text-white mb-6">CDL Teams</h1>
        <p className="text-xl text-gray-300 max-w-3xl mx-auto leading-relaxed">
          Professional Call of Duty League organizations competing at the highest level.
          Explore team rosters, player profiles, and comprehensive statistics.
        </p>
        <div className="mt-8 text-gray-400 text-lg">
          {teams.length} Active Teams
        </div>
      </div>

      {/* Teams Grid */}
      <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-8">
        {teams.map((team) => (
          <div key={team.id} className="group">
            <div className="card hover:bg-gray-750 transition-all duration-300 transform hover:scale-105 border border-gray-700 hover:border-blue-500 overflow-hidden">
              {/* Team Header */}
              <div className="flex items-center justify-between mb-6 p-6">
                <div className="flex items-center space-x-4">
                  {/* Team Logo */}
                  <TeamLogo team={team} size="xl" />
                  <div>
                    <h3 className="text-2xl font-bold text-white">{team.name}</h3>
                    <p className="text-gray-400">{team.city}</p>
                  </div>
                </div>
                <div className="text-right">
                  <span className="inline-block bg-blue-600 text-white px-4 py-2 rounded-full text-sm font-bold">
                    {team.abbreviation}
                  </span>
                </div>
              </div>

              {/* Team Stats */}
              <div className="grid grid-cols-3 gap-4 mb-6">
                <div className="text-center">
                  <div className="text-2xl font-bold text-green-400">
                    {team.players?.length || 0}
                  </div>
                  <div className="text-gray-400 text-sm">Players</div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold text-blue-400">
                    {team.is_active ? 'Active' : 'Inactive'}
                  </div>
                  <div className="text-gray-400 text-sm">Status</div>
                </div>
                <div className="text-center">
                  <div className="text-2xl font-bold text-purple-400">
                    {team.founded_date ? new Date(team.founded_date).getFullYear() : 'N/A'}
                  </div>
                  <div className="text-gray-400 text-sm">Founded</div>
                </div>
              </div>

              {/* Team Colors */}
              {team.primary_color && (
                <div className="mb-6">
                  <h4 className="text-white font-semibold mb-2">Team Colors</h4>
                  <div className="flex space-x-2">
                    <div 
                      className="w-8 h-8 rounded-full border-2 border-gray-600"
                      style={{ backgroundColor: team.primary_color }}
                    ></div>
                    {team.secondary_color && (
                      <div 
                        className="w-8 h-8 rounded-full border-2 border-gray-600"
                        style={{ backgroundColor: team.secondary_color }}
                      ></div>
                    )}
                  </div>
                </div>
              )}

              {/* Players Section */}
              <div className="mb-6">
                <h4 className="text-white font-semibold mb-3">Current Roster</h4>
                {team.players && team.players.length > 0 ? (
                  <div className="space-y-2">
                    {team.players.slice(0, 4).map((player) => (
                      <div key={player.id} className="flex items-center justify-between p-2 bg-gray-800 rounded">
                        <div className="flex items-center space-x-3">
                          <PlayerAvatar player={player} size="sm" />
                          <div>
                            <div className="text-white font-medium">{player.gamertag}</div>
                            <div className="text-gray-400 text-xs">{player.role || 'Player'}</div>
                          </div>
                        </div>
                        <Link 
                          to={`/players/${player.id}`}
                          className="text-blue-400 hover:text-blue-300 text-sm"
                        >
                          View
                        </Link>
                      </div>
                    ))}
                    {team.players.length > 4 && (
                      <div className="text-center py-2">
                        <span className="text-gray-400 text-sm">
                          +{team.players.length - 4} more players
                        </span>
                      </div>
                    )}
                  </div>
                ) : (
                  <div className="text-center py-4">
                    <div className="text-gray-400 text-sm">No players assigned</div>
                  </div>
                )}
              </div>

              {/* Action Buttons */}
              <div className="flex space-x-2">
                <Link
                  to={`/teams/${team.id}`}
                  className="btn-primary flex-1 text-center"
                >
                  Team Details
                </Link>
                <Link
                  to={`/teams/${team.id}/players`}
                  className="btn-secondary flex-1 text-center"
                >
                  Full Roster
                </Link>
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Call to Action */}
      <div className="text-center py-12 bg-gradient-to-r from-gray-800 to-gray-900 rounded-lg">
        <h2 className="text-3xl font-bold text-white mb-4">Follow Your Favorite Team</h2>
        <p className="text-gray-300 mb-6 max-w-2xl mx-auto">
          Stay updated with the latest news, match results, and player statistics from your favorite CDL teams.
        </p>
        <div className="flex justify-center space-x-4">
          <Link to="/players" className="btn-primary">
            View All Players
          </Link>
          <Link to="/" className="btn-secondary">
            Back to Home
          </Link>
        </div>
      </div>

      {teams.length === 0 && (
        <div className="text-center py-12">
          <div className="text-gray-400 text-xl mb-4">No teams found</div>
          <p className="text-gray-500">There are currently no teams available.</p>
        </div>
      )}
    </div>
  );
};

export default Teams; 