import { Link } from "react-router-dom";

const NotFound = () => (
  <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-32 text-center">
    <p className="text-xs uppercase tracking-widest text-[#737373] mb-4">404</p>
    <h1 className="font-grotesk text-3xl font-bold text-white mb-6">
      Page not found
    </h1>
    <Link
      to="/"
      className="text-[#737373] hover:text-white text-sm transition-colors"
    >
      ← Back to home
    </Link>
  </div>
);

export default NotFound;
