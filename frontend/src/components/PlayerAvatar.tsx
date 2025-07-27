import React from 'react';

interface PlayerAvatarProps {
  player: {
    id: number;
    gamertag: string;
    first_name?: string;
    last_name?: string;
    avatar_url?: string;
  };
  size?: 'sm' | 'md' | 'lg' | 'xl';
  className?: string;
}

const PlayerAvatar: React.FC<PlayerAvatarProps> = ({ player, size = 'md', className = '' }) => {
  const sizeClasses = {
    sm: 'w-8 h-8 text-xs',
    md: 'w-12 h-12 text-sm',
    lg: 'w-16 h-16 text-lg',
    xl: 'w-24 h-24 text-xl'
  };

  const sizeClassesImg = {
    sm: 'w-8 h-8',
    md: 'w-12 h-12',
    lg: 'w-16 h-16',
    xl: 'w-24 h-24'
  };

  // If we have an avatar URL, use it
  if (player.avatar_url) {
    return (
      <img
        src={player.avatar_url}
        alt={`${player.gamertag} avatar`}
        className={`${sizeClassesImg[size]} rounded-full object-cover ${className}`}
        onError={(e) => {
          // Fallback to placeholder if image fails to load
          const target = e.target as HTMLImageElement;
          target.style.display = 'none';
          const parent = target.parentElement;
          if (parent) {
            const fallback = document.createElement('div');
            fallback.className = `bg-gradient-to-br from-green-500 to-blue-500 rounded-full flex items-center justify-center text-white font-bold ${sizeClasses[size]} ${className}`;
            fallback.textContent = player.gamertag.charAt(0).toUpperCase();
            parent.appendChild(fallback);
          }
        }}
      />
    );
  }

  // Fallback to placeholder
  return (
    <div className={`bg-gradient-to-br from-green-500 to-blue-500 rounded-full flex items-center justify-center text-white font-bold ${sizeClasses[size]} ${className}`}>
      {player.gamertag.charAt(0).toUpperCase()}
    </div>
  );
};

export default PlayerAvatar; 