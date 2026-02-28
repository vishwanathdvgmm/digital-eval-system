import React from "react";
import { useAuth } from "../hooks/useAuth";

const Navbar: React.FC = () => {
  const { user, logout } = useAuth();

  return (
    <header className="sticky top-0 z-50 w-full bg-white/80 backdrop-blur-md border-b border-slate-200 px-8 py-4 flex items-center justify-between">
      <div>
        <h2 className="text-lg font-semibold text-slate-800 tracking-tight">
          {user?.role.toUpperCase()} <span className="text-slate-400 font-normal">Dashboard</span>
        </h2>
      </div>

      <div className="flex items-center gap-6">
        <div className="text-right hidden sm:block">
          <p className="text-sm font-medium text-slate-700">{user?.user_id}</p>
          <p className="text-xs text-slate-500">{user?.email}</p>
        </div>
        
        <button
          onClick={logout}
          className="px-4 py-2 rounded-lg text-sm font-medium text-red-600 bg-red-50 hover:bg-red-100 hover:text-red-700 transition-colors"
        >
          Logout
        </button>
      </div>
    </header>
  );
};

export default Navbar;


