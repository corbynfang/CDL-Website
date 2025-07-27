import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import type { Team, Player, TeamTournamentStats } from '../types';
import { teamApi } from '../services/api';
import TeamLogo from './TeamLogo';
import PlayerAvatar from './PlayerAvatar';

const TeamDetails: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [team, setTeam] = useState<Team | null>(null);
  const [players, setPlayers] = useState<Player[]>([]);
  const [stats, setStats] = useState<TeamTournamentStats[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchTeamData = async () => {
      if (!id) return;
      
      try {
        setLoading(true);
        const [teamData, playersData, statsData] = await Promise.all([
          teamApi.getTeam(parseInt(id)),
          teamApi.getTeamPlayers(parseInt(id)),
          teamApi.getTeamStats(parseInt(id))
        ]);
        setTeam(teamData);
        setPlayers(playersData);
        setStats(statsData);
      } catch (err) {
        setError('Failed to fetch team data');
        console.error('Error fetching team data:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchTeamData();
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
      {/* Team Header */}
      <div className="bg-gradient-to-r from-blue-900 to-purple-900 rounded-lg p-8">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-6">
            {/* Team Logo */}
            <TeamLogo team={team} size="xl" />
            <div>
              <h1 className="text-4xl font-bold text-white mb-2">{team.name}</h1>
              <div className="flex items-center space-x-4 text-gray-300">
                {team.city && (
                  <div className="flex items-center">
                    <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                    </svg>
                    {team.city}
                  </div>
                )}
                <div className="flex items-center">
                  <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  {team.is_active ? 'Active' : 'Inactive'}
                </div>
                {team.founded_date && (
                  <div className="flex items-center">
                    <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                    </svg>
                    Founded {new Date(team.founded_date).getFullYear()}
                  </div>
                )}
              </div>
            </div>
          </div>
          <div className="text-right">
            <span className="inline-block bg-blue-600 text-white px-4 py-2 rounded-full text-lg font-medium">
              {team.abbreviation}
            </span>
          </div>
        </div>
      </div>

      {/* Team Colors */}
      {(team.primary_color || team.secondary_color) && (
        <div className="card">
          <h2 className="text-xl font-semibold text-white mb-4">Team Colors</h2>
          <div className="flex space-x-4">
            {team.primary_color && (
              <div className="text-center">
                <div 
                  className="w-16 h-16 rounded-full border-4 border-gray-600 mb-2"
                  style={{ backgroundColor: team.primary_color }}
                ></div>
                <span className="text-gray-400 text-sm">Primary</span>
              </div>
            )}
            {team.secondary_color && (
              <div className="text-center">
                <div 
                  className="w-16 h-16 rounded-full border-4 border-gray-600 mb-2"
                  style={{ backgroundColor: team.secondary_color }}
                ></div>
                <span className="text-gray-400 text-sm">Secondary</span>
              </div>
            )}
          </div>
        </div>
      )}

      {/* Team Stats Overview */}
      <div className="grid md:grid-cols-4 gap-6">
        <div className="card text-center">
          <div className="text-3xl font-bold text-green-400">
            {players.length}
          </div>
          <div className="text-gray-400 text-sm">Active Players</div>
        </div>
        <div className="card text-center">
          <div className="text-3xl font-bold text-blue-400">
            {stats.length}
          </div>
          <div className="text-gray-400 text-sm">Tournaments</div>
        </div>
        <div className="card text-center">
          <div className="text-3xl font-bold text-purple-400">
            {stats.reduce((total, stat) => total + stat.matches_played, 0)}
          </div>
          <div className="text-gray-400 text-sm">Matches Played</div>
        </div>
        <div className="card text-center">
          <div className="text-3xl font-bold text-yellow-400">
            {stats.reduce((total, stat) => total + stat.matches_won, 0)}
          </div>
          <div className="text-gray-400 text-sm">Matches Won</div>
        </div>
      </div>

      {/* Roster Section */}
      <div className="card">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-2xl font-semibold text-white">Current Roster</h2>
          <span className="text-gray-400">{players.length} Players</span>
        </div>
        
        {players.length > 0 ? (
          <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
            {players.map((player) => (
              <div key={player.id} className="bg-gray-800 rounded-lg p-4 hover:bg-gray-750 transition-colors">
                <div className="flex items-center space-x-3 mb-3">
                  <PlayerAvatar player={player} size="lg" />
                  <div>
                    <div className="text-white font-semibold">{player.gamertag}</div>
                    <div className="text-gray-400 text-sm">{player.role || 'Player'}</div>
                  </div>
                </div>
                
                <div className="space-y-1 text-sm">
                  {player.first_name && player.last_name && (
                    <div className="text-gray-300">
                      {player.first_name} {player.last_name}
                    </div>
                  )}
                  {player.country && (
                    <div className="text-gray-400">
                      Country: {player.country}
                    </div>
                  )}
                  <div className="text-gray-400">
                    Status: {player.is_active ? 'Active' : 'Inactive'}
                  </div>
                </div>
                
                <div className="mt-4">
                  <Link 
                    to={`/players/${player.id}`}
                    className="btn-primary w-full text-center text-sm"
                  >
                    View Profile
                  </Link>
                </div>
              </div>
            ))}
          </div>
        ) : (
          <div className="text-center py-8">
            <div className="text-gray-400 text-lg mb-2">No players assigned</div>
            <p className="text-gray-500">This team currently has no active players.</p>
          </div>
        )}
      </div>

      {/* Tournament Performance */}
      {stats.length > 0 && (
        <div className="card">
          <h2 className="text-2xl font-semibold text-white mb-6">Tournament Performance</h2>
          <div className="overflow-x-auto">
            <table className="w-full">
              <thead>
                <tr className="border-b border-gray-700">
                  <th className="text-left py-3 px-4 text-gray-300">Tournament</th>
                  <th className="text-right py-3 px-4 text-gray-300">Placement</th>
                  <th className="text-right py-3 px-4 text-gray-300">Matches</th>
                  <th className="text-right py-3 px-4 text-gray-300">W-L</th>
                  <th className="text-right py-3 px-4 text-gray-300">Maps</th>
                  <th className="text-right py-3 px-4 text-gray-300">Prize Money</th>
                </tr>
              </thead>
              <tbody>
                {stats.map((stat) => (
                  <tr key={stat.id} className="border-b border-gray-800">
                    <td className="py-3 px-4 text-white font-medium">
                      {stat.tournament?.name || 'Unknown Tournament'}
                    </td>
                    <td className="py-3 px-4 text-right text-gray-300">
                      {stat.placement ? `${stat.placement}${getOrdinalSuffix(stat.placement)}` : 'N/A'}
                    </td>
                    <td className="py-3 px-4 text-right text-gray-300">
                      {stat.matches_played}
                    </td>
                    <td className="py-3 px-4 text-right text-gray-300">
                      <span className="text-green-400">{stat.matches_won}</span>
                      <span className="text-gray-500 mx-1">-</span>
                      <span className="text-red-400">{stat.matches_lost}</span>
                    </td>
                    <td className="py-3 px-4 text-right text-gray-300">
                      <span className="text-green-400">{stat.maps_won}</span>
                      <span className="text-gray-500 mx-1">-</span>
                      <span className="text-red-400">{stat.maps_lost}</span>
                    </td>
                    <td className="py-3 px-4 text-right text-gray-300">
                      ${stat.prize_money.toLocaleString()}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {/* Action Buttons */}
      <div className="flex justify-center space-x-4">
        <Link to="/teams" className="btn-secondary">
          Back to Teams
        </Link>
        <Link to={`/teams/${team.id}/players`} className="btn-primary">
          View Full Roster
        </Link>
      </div>
    </div>
  );
};

// Helper function to get ordinal suffix
const getOrdinalSuffix = (num: number): string => {
  const j = num % 10;
  const k = num % 100;
  if (j === 1 && k !== 11) {
    return "st";
  }
  if (j === 2 && k !== 12) {
    return "nd";
  }
  if (j === 3 && k !== 13) {
    return "rd";
  }
  return "th";
};

export default TeamDetails; 