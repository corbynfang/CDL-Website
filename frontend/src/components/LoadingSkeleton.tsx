import React from 'react';

interface LoadingSkeletonProps {
  variant?: 'card' | 'list' | 'table' | 'profile';
  count?: number;
}

/**
 * Loading skeleton component for better UX during data fetching
 * 
 * @param variant - Type of skeleton to display
 * @param count - Number of skeleton items to show
 */
const LoadingSkeleton: React.FC<LoadingSkeletonProps> = ({ 
  variant = 'card', 
  count = 1 
}) => {
  const renderSkeleton = () => {
    switch (variant) {
      case 'card':
        return (
          <div className="card border border-gray-800">
            <div className="animate-pulse p-6">
              <div className="flex items-center space-x-4 mb-6">
                <div className="rounded-full bg-gray-700 h-16 w-16"></div>
                <div className="flex-1 space-y-3">
                  <div className="h-6 bg-gray-700 rounded w-3/4"></div>
                  <div className="h-4 bg-gray-700 rounded w-1/2"></div>
                </div>
              </div>
              <div className="space-y-3">
                <div className="h-4 bg-gray-700 rounded"></div>
                <div className="h-4 bg-gray-700 rounded"></div>
                <div className="h-4 bg-gray-700 rounded w-5/6"></div>
              </div>
              <div className="mt-6">
                <div className="h-10 bg-gray-700 rounded"></div>
              </div>
            </div>
          </div>
        );

      case 'list':
        return (
          <div className="border-b border-gray-800 py-4">
            <div className="animate-pulse flex items-center space-x-4">
              <div className="rounded-full bg-gray-700 h-12 w-12"></div>
              <div className="flex-1 space-y-2">
                <div className="h-4 bg-gray-700 rounded w-3/4"></div>
                <div className="h-3 bg-gray-700 rounded w-1/2"></div>
              </div>
            </div>
          </div>
        );

      case 'table':
        return (
          <tr className="border-b border-gray-800">
            <td className="py-4">
              <div className="animate-pulse">
                <div className="h-4 bg-gray-700 rounded w-full"></div>
              </div>
            </td>
            <td className="py-4">
              <div className="animate-pulse">
                <div className="h-4 bg-gray-700 rounded w-full"></div>
              </div>
            </td>
            <td className="py-4">
              <div className="animate-pulse">
                <div className="h-4 bg-gray-700 rounded w-full"></div>
              </div>
            </td>
          </tr>
        );

      case 'profile':
        return (
          <div className="space-y-8">
            {/* Hero Section */}
            <div className="card border border-gray-800">
              <div className="animate-pulse p-8">
                <div className="flex flex-col md:flex-row items-center gap-6 mb-6">
                  <div className="rounded-full bg-gray-700 h-32 w-32"></div>
                  <div className="flex-1 space-y-4 w-full">
                    <div className="h-8 bg-gray-700 rounded w-1/2"></div>
                    <div className="h-6 bg-gray-700 rounded w-1/3"></div>
                    <div className="h-6 bg-gray-700 rounded w-1/4"></div>
                  </div>
                </div>
              </div>
            </div>

            {/* Stats Grid */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              {[1, 2, 3].map((i) => (
                <div key={i} className="card border border-gray-800">
                  <div className="animate-pulse p-6">
                    <div className="h-4 bg-gray-700 rounded w-1/2 mb-4"></div>
                    <div className="h-8 bg-gray-700 rounded w-3/4"></div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        );

      default:
        return null;
    }
  };

  return (
    <>
      {Array.from({ length: count }).map((_, index) => (
        <React.Fragment key={index}>
          {renderSkeleton()}
        </React.Fragment>
      ))}
    </>
  );
};

export default LoadingSkeleton;

/**
 * Simple spinner component for inline loading states
 */
export const Spinner: React.FC<{ size?: 'sm' | 'md' | 'lg' }> = ({ size = 'md' }) => {
  const sizeClasses = {
    sm: 'h-4 w-4',
    md: 'h-8 w-8',
    lg: 'h-12 w-12',
  };

  return (
    <div className="flex justify-center items-center">
      <div className={`animate-spin rounded-none border-b-2 border-white ${sizeClasses[size]}`}></div>
    </div>
  );
};

/**
 * Error display component with retry button
 */
interface ErrorDisplayProps {
  message: string;
  onRetry?: () => void;
}

export const ErrorDisplay: React.FC<ErrorDisplayProps> = ({ message, onRetry }) => {
  return (
    <div className="text-center py-12 px-4">
      <div className="text-red-500 text-xl mb-6 flex items-center justify-center gap-3">
        <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        {message}
      </div>
      {onRetry && (
        <button
          onClick={onRetry}
          className="btn-primary px-6 py-3 text-base"
        >
          TRY AGAIN
        </button>
      )}
    </div>
  );
};

