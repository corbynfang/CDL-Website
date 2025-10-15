import React, { useState } from 'react';

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
  const [imageError, setImageError] = useState(false);
  const [fallbackError, setFallbackError] = useState(false);

  const sizeClasses = {
    sm: 'w-8 h-8 text-xs',
    md: 'w-12 h-12 text-sm',
    lg: 'w-16 h-16 text-lg',
    xl: 'w-24 h-24 text-xl',
    '2xl': 'w-32 h-32 text-2xl'
  };

  const FALLBACK_AVATAR = '/assets/avatars/Unknown.webp';
  const avatarUrl = player.avatar_url || FALLBACK_AVATAR;

  const handleImageError = () => {
    if (!imageError) {
      setImageError(true);
    } else if (!fallbackError) {
      setFallbackError(true);
    }
  };

  if (fallbackError || (!avatarUrl && imageError)) {
    return (
      <div className={`bg-gradient-to-br from-gray-600 to-gray-800 rounded-full flex items-center justify-center text-white font-bold ${sizeClasses[size]} ${className} flex-shrink-0`}>
        {player.gamertag.charAt(0).toUpperCase()}
      </div>
    );
  }

  const displayUrl = imageError && avatarUrl !== FALLBACK_AVATAR ? FALLBACK_AVATAR : avatarUrl;

  return (
    <div className={`relative ${sizeClasses[size]} ${className} flex-shrink-0`}>
      <img
        src={displayUrl}
        alt={`${player.gamertag} avatar`}
        className="w-full h-full rounded-full object-cover object-center"
        onError={handleImageError}
        loading="lazy"
      />
    </div>
  );
};

export default PlayerAvatar; 