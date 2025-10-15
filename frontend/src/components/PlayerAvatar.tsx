import React from 'react';

interface PlayerAvatarProps {
  player: {
    id: number;
    gamertag: string;
    first_name?: string;
    last_name?: string;
    avatar_url?: string;
  };
  size?: 'sm' | 'md' | 'lg' | 'xl' | '2xl';
  className?: string;
}

const PlayerAvatar: React.FC<PlayerAvatarProps> = ({ player, size = 'md', className = '' }) => {
  const sizeClasses = {
    sm: 'w-8 h-8 text-xs',
    md: 'w-12 h-12 text-sm',
    lg: 'w-16 h-16 text-lg',
    xl: 'w-24 h-24 text-xl',
    '2xl': 'w-32 h-32 text-2xl'
  };

  const sizeClassesImg = {
    sm: 'w-8 h-8',
    md: 'w-12 h-12',
    lg: 'w-16 h-16',
    xl: 'w-24 h-24',
    '2xl': 'w-32 h-32'
  };

  // Fallback image for players without avatars
  const FALLBACK_AVATAR = '/assets/avatars/Unknown.webp';

  // Determine which avatar to use
  const avatarUrl = player.avatar_url || FALLBACK_AVATAR;

  return (
    <div className={`relative ${sizeClassesImg[size]} ${className}`}>
      <img
        src={avatarUrl}
        alt={`${player.gamertag} avatar`}
        className={`w-full h-full rounded-full object-cover object-center`}
        onError={(e) => {
          // If the avatar fails to load, use the Unknown.webp fallback
          const target = e.target as HTMLImageElement;
          if (target.src !== FALLBACK_AVATAR) {
            target.src = FALLBACK_AVATAR;
          } else {
            // If even Unknown.webp fails, show initials
            target.style.display = 'none';
            const parent = target.parentElement;
            if (parent) {
              const fallback = document.createElement('div');
              fallback.className = `bg-gradient-to-br from-gray-600 to-gray-800 rounded-full flex items-center justify-center text-white font-bold ${sizeClasses[size]} w-full h-full`;
              fallback.textContent = player.gamertag.charAt(0).toUpperCase();
              parent.appendChild(fallback);
            }
          }
        }}
        loading="lazy"
      />
    </div>
  );
};

export default PlayerAvatar; 