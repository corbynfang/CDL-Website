import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import type { Team, Player } from '../types';
import { teamApi } from '../services/api';
import TeamLogo from './TeamLogo';
import PlayerAvatar from './PlayerAvatar';

const TeamPlayers: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [team, setTeam] = useState<Team | null>(null);
  const [players, setPlayers] = useState<Player[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchTeamPlayers = async () => {
      if (!id) return;
      
      try {
        setLoading(true);
        const [teamData, playersData] = await Promise.all([
          teamApi.getTeam(parseInt(id)),
          teamApi.getTeamPlayers(parseInt(id))
        ]);
        setTeam(teamData);
        setPlayers(playersData);
      } catch (err) {
        setError('Failed to fetch team players');
        console.error('Error fetching team players:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchTeamPlayers();
  }, [id]);

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
        <Link to="/teams" className="btn-primary">
          Back to Teams
        </Link>
      </div>
    );
  }

  if (!team) {
    return (
      <div className="text-center py-8">
        <div className="text-gray-400 text-xl mb-4">Team not found</div>
        <Link to="/teams" className="btn-primary">
          Back to Teams
        </Link>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-white">{team.name} Roster</h1>
          <p className="text-gray-400">{players.length} Players</p>
        </div>
        <Link to={`/teams/${team.id}`} className="btn-secondary">
          Back to Team
        </Link>
      </div>

      {/* Team Info Card */}
      <div className="card">
        <div className="flex items-center space-x-4">
          <TeamLogo team={team} size="lg" />
          <div>
            <h2 className="text-xl font-semibold text-white">{team.name}</h2>
            <div className="flex items-center space-x-4 text-gray-400 text-sm">
              {team.city && <span>{team.city}</span>}
              <span>{team.is_active ? 'Active' : 'Inactive'}</span>
              {team.founded_date && (
                <span>Founded {new Date(team.founded_date).getFullYear()}</span>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Players Grid */}
      {players.length > 0 ? (
        <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
          {players.map((player) => (
            <div key={player.id} className="card hover:bg-gray-750 transition-all duration-300 transform hover:scale-105">
              {/* Player Header */}
              <div className="flex items-center space-x-4 mb-4">
                <PlayerAvatar player={player} size="lg" />
                <div className="flex-1">
                  <h3 className="text-xl font-bold text-white">{player.gamertag}</h3>
                  <div className="flex items-center space-x-2">
                    <span className="text-gray-400 text-sm">{player.role || 'Player'}</span>
                    <span className={`inline-block px-2 py-1 rounded-full text-xs font-medium ${
                      player.is_active ? 'bg-green-600 text-white' : 'bg-red-600 text-white'
                    }`}>
                      {player.is_active ? 'Active' : 'Inactive'}
                    </span>
                  </div>
                </div>
              </div>

              {/* Player Details */}
              <div className="space-y-3 mb-4">
                {player.first_name && player.last_name && (
                  <div className="flex justify-between">
                    <span className="text-gray-400">Name:</span>
                    <span className="text-white font-medium">
                      {player.first_name} {player.last_name}
                    </span>
                  </div>
                )}
                
                {player.country && (
                  <div className="flex justify-between">
                    <span className="text-gray-400">Country:</span>
                    <span className="text-white font-medium">{player.country}</span>
                  </div>
                )}

                {player.birthdate && (
                  <div className="flex justify-between">
                    <span className="text-gray-400">Birthdate:</span>
                    <span className="text-white font-medium">
                      {new Date(player.birthdate).toLocaleDateString()}
                    </span>
                  </div>
                )}
              </div>

              {/* Social Links */}
              <div className="flex space-x-2 mb-4">
                {player.twitter_handle && (
                  <a 
                    href={`https://twitter.com/${player.twitter_handle}`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="flex-1 bg-blue-600 hover:bg-blue-700 text-white text-center py-2 px-3 rounded text-sm font-medium transition-colors"
                  >
                    Twitter
                  </a>
                )}
                {player.liquipedia_url && (
                  <a 
                    href={player.liquipedia_url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="flex-1 bg-purple-600 hover:bg-purple-700 text-white text-center py-2 px-3 rounded text-sm font-medium transition-colors"
                  >
                    Liquipedia
                  </a>
                )}
              </div>

              {/* Action Buttons */}
              <div className="flex space-x-2">
                <Link 
                  to={`/players/${player.id}`}
                  className="btn-primary flex-1 text-center"
                >
                  View Profile
                </Link>
                <Link 
                  to={`/players/${player.id}/kd-stats`}
                  className="btn-secondary flex-1 text-center"
                >
                  KD Stats
                </Link>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <div className="card text-center py-12">
          <div className="text-gray-400 text-xl mb-4">No Players Found</div>
          <p className="text-gray-500 mb-6">This team currently has no active players assigned.</p>
          <Link to="/teams" className="btn-primary">
            Back to Teams
          </Link>
        </div>
      )}

      {/* Team Stats Summary */}
      {players.length > 0 && (
        <div className="card">
          <h3 className="text-xl font-semibold text-white mb-4">Roster Summary</h3>
          <div className="grid md:grid-cols-4 gap-4">
            <div className="text-center">
              <div className="text-2xl font-bold text-green-400">{players.length}</div>
              <div className="text-gray-400 text-sm">Total Players</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-blue-400">
                {players.filter(p => p.is_active).length}
              </div>
              <div className="text-gray-400 text-sm">Active Players</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-purple-400">
                {players.filter(p => p.role).length}
              </div>
              <div className="text-gray-400 text-sm">With Roles</div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-yellow-400">
                {players.filter(p => p.country).length}
              </div>
              <div className="text-gray-400 text-sm">With Country</div>
            </div>
          </div>
        </div>
      )}

      {/* Navigation */}
      <div className="flex justify-center space-x-4">
        <Link to="/teams" className="btn-secondary">
          All Teams
        </Link>
        <Link to={`/teams/${team.id}`} className="btn-primary">
          Team Details
        </Link>
      </div>
    </div>
  );
};

export default TeamPlayers; 