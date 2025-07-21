import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import type { Team } from '../types';
import { teamApi } from '../services/api';

const Teams: React.FC = () => {
  const [teams, setTeams] = useState<Team[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchTeams = async () => {
      try {
        setLoading(true);
        const data = await teamApi.getTeams();
        setTeams(data);
      } catch (err) {
        setError('Failed to fetch teams');
        console.error('Error fetching teams:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchTeams();
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
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-3xl font-bold text-white">CDL Teams</h1>
        <div className="text-gray-400">
          {teams.length} Active Teams
        </div>
      </div>

      <div className="grid md:grid-cols-2 lg:grid-cols-3 gap-6">
        {teams.map((team) => (
          <div key={team.id} className="card hover:bg-gray-750 transition-colors duration-200">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-xl font-semibold text-white">{team.name}</h3>
              <span className="text-sm text-gray-400 bg-gray-700 px-2 py-1 rounded">
                {team.abbreviation}
              </span>
            </div>

            <div className="space-y-2 mb-4">
              {team.city && (
                <div className="flex items-center text-gray-300">
                  <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                  </svg>
                  {team.city}
                </div>
              )}



              <div className="flex items-center text-gray-300">
                <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                {team.is_active ? 'Active' : 'Inactive'}
              </div>
            </div>

            <div className="flex space-x-2">
              <Link
                to={`/teams/${team.id}`}
                className="btn-primary flex-1 text-center"
              >
                View Details
              </Link>
              <Link
                to={`/teams/${team.id}/players`}
                className="btn-secondary flex-1 text-center"
              >
                Players
              </Link>
            </div>
          </div>
        ))}
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