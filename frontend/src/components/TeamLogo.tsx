import React from 'react';

interface TeamLogoProps {
  team: {
    id: number;
    name: string;
    abbreviation: string;
    logo_url?: string;
  };
  size?: 'sm' | 'md' | 'lg' | 'xl';
  className?: string;
}

const TeamLogo: React.FC<TeamLogoProps> = ({ team, size = 'md', className = '' }) => {
  const sizeClasses = {
    sm: 'w-8 h-8 text-xs',
    md: 'w-12 h-12 text-sm',
    lg: 'w-16 h-16 text-lg',
    xl: 'w-24 h-24 text-2xl'
  };

  const sizeClassesImg = {
    sm: 'w-8 h-8',
    md: 'w-12 h-12',
    lg: 'w-16 h-16',
    xl: 'w-24 h-24'
  };

  // If we have a logo URL, use it
  if (team.logo_url) {
    return (
      <div className={`relative ${sizeClassesImg[size]} ${className}`}>
        <img
          src={team.logo_url}
          alt={`${team.name} logo`}
          className={`w-full h-full object-contain object-center`}
          onError={(e) => {
            // Fallback to placeholder if image fails to load
            const target = e.target as HTMLImageElement;
            target.style.display = 'none';
            const parent = target.parentElement;
            if (parent) {
              const fallback = document.createElement('div');
              fallback.className = `bg-gradient-to-br from-blue-500 to-purple-600 rounded-full flex items-center justify-center text-white font-bold ${sizeClasses[size]} ${className} w-full h-full`;
              fallback.textContent = team.abbreviation;
              parent.appendChild(fallback);
            }
          }}
          loading="lazy"
        />
      </div>
    );
  }

  // Fallback to placeholder
  return (
    <div className={`bg-gradient-to-br from-blue-500 to-purple-600 rounded-full flex items-center justify-center text-white font-bold ${sizeClasses[size]} ${className} flex-shrink-0`}>
      {team.abbreviation}
    </div>
  );
};

export default TeamLogo; 