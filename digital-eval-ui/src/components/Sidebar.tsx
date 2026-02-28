import React from "react";
import { Link, useLocation } from "react-router-dom";
import { useAuth } from "../hooks/useAuth";

const Sidebar: React.FC = () => {
  const { user } = useAuth();
  const location = useLocation();

  if (!user) return null;

  const role = user.role;

  const menu: Record<string, { label: string; to: string }[]> = {
    examiner: [
      { label: "Upload Scripts", to: "/dashboard/examiner/upload" },
      { label: "Upload History", to: "/dashboard/examiner/history" },
    ],
    evaluator: [
      { label: "Assigned Scripts", to: "/dashboard/evaluator/assigned" },
      { label: "Request Evaluation", to: "/dashboard/evaluator/request" },
      { label: "Request History", to: "/dashboard/evaluator/requests" },
    ],
    authority: [
      { label: "Pending Requests", to: "/dashboard/authority/requests" },
      { label: "Approve Requests", to: "/dashboard/authority/approve" },
      { label: "Release Results", to: "/dashboard/authority/release" },
    ],
    student: [
      { label: "View Results", to: "/dashboard/student/results" },
    ],
    admin: [
      { label: "Dashboard", to: "/dashboard/admin" },
      { label: "User Management", to: "/dashboard/admin/users" },
    ],
  };

  const items = menu[role] ?? [];

  return (
    <aside className="w-72 bg-slate-900 text-white h-screen flex flex-col shadow-2xl">
      <div className="p-6 border-b border-slate-800">
        <h1 className="text-2xl font-bold bg-gradient-to-r from-indigo-400 to-cyan-400 bg-clip-text text-transparent">
          Digital Eval
        </h1>
        <p className="text-xs text-slate-400 mt-1 uppercase tracking-wider">{role} Portal</p>
      </div>

      <nav className="flex-1 p-4 space-y-2 overflow-y-auto">
        {items.map((item) => {
          const isActive = location.pathname === item.to;
          return (
            <Link
              key={item.to}
              to={item.to}
              className={`block px-4 py-3 rounded-xl text-sm font-medium transition-all duration-200 ${
                isActive
                  ? "bg-indigo-600 text-white shadow-lg shadow-indigo-900/50 translate-x-1"
                  : "text-slate-300 hover:bg-slate-800 hover:text-white hover:translate-x-1"
              }`}
            >
              {item.label}
            </Link>
          );
        })}
      </nav>

      <div className="p-4 border-t border-slate-800">
        <div className="flex items-center gap-3 px-4 py-3 rounded-xl bg-slate-800/50">
          <div className="w-8 h-8 rounded-full bg-gradient-to-tr from-indigo-500 to-purple-500 flex items-center justify-center text-xs font-bold">
            {user.user_id.charAt(0).toUpperCase()}
          </div>
          <div className="overflow-hidden">
            <p className="text-sm font-medium truncate">{user.user_id}</p>
            <p className="text-xs text-slate-400 truncate">{user.email}</p>
          </div>
        </div>
      </div>
    </aside>
  );
};

export default Sidebar;