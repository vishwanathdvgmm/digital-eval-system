import React from "react";
import { useAuth } from "../../hooks/useAuth";

const TopbarUser: React.FC = () => {
  const { user, logout } = useAuth();

  return (
    <div className="flex items-center space-x-3">
      <div className="text-sm text-gray-700">
        <div className="font-medium">{user?.name || user?.user_id}</div>
        <div className="text-xs text-gray-500 capitalize">{user?.role}</div>
      </div>
      <button
        onClick={() => logout()}
        className="px-3 py-1 border rounded text-sm hover:bg-gray-100"
      >
        Logout
      </button>
    </div>
  );
};

export default TopbarUser;