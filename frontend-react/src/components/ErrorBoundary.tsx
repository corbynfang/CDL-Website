import { Component, type ReactNode } from "react";

interface Props {
  children: ReactNode;
}

interface State {
  crashed: boolean;
}

class ErrorBoundary extends Component<Props, State> {
  state: State = { crashed: false };

  static getDerivedStateFromError(): State {
    return { crashed: true };
  }

  componentDidCatch(error: Error) {
    console.error("App crashed:", error);
  }

  render() {
    if (this.state.crashed) {
      return (
        <div className="min-h-screen flex items-center justify-center bg-[#0a0a0a]">
          <div className="text-center px-4">
            <p className="text-xs uppercase tracking-widest text-[#737373] mb-4">Error</p>
            <h1 className="font-grotesk text-2xl font-bold text-white mb-4">
              Something went wrong
            </h1>
            <p className="text-[#737373] text-sm mb-8">
              The page crashed. Try refreshing.
            </p>
            <button
              onClick={() => window.location.reload()}
              className="text-[#737373] hover:text-white text-sm transition-colors border border-[#1a1a1a] px-4 py-2"
            >
              Refresh
            </button>
          </div>
        </div>
      );
    }
    return this.props.children;
  }
}

export default ErrorBoundary;
