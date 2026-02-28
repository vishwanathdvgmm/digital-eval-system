import React, { useState } from "react";
import Card from "../../components/Card";
import Button from "../../components/Button";

const UserManagement: React.FC = () => {
  const [users, setUsers] = useState([
    { id: "1", name: "Admin User", email: "admin@example.com", role: "admin" },
    { id: "2", name: "John Doe", email: "john@example.com", role: "authority" },
    { id: "3", name: "Jane Smith", email: "jane@example.com", role: "examiner" },
  ]);

  const handleDelete = (id: string) => {
    if (confirm("Are you sure you want to delete this user?")) {
      setUsers(users.filter(u => u.id !== id));
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold text-slate-800">User Management</h1>
        <Button variant="primary">Add User</Button>
      </div>

      <Card>
        <div className="overflow-x-auto">
          <table className="w-full text-left">
            <thead>
              <tr className="border-b border-slate-200">
                <th className="pb-3 font-semibold text-slate-600">Name</th>
                <th className="pb-3 font-semibold text-slate-600">Email</th>
                <th className="pb-3 font-semibold text-slate-600">Role</th>
                <th className="pb-3 font-semibold text-slate-600 text-right">Actions</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {users.map((user) => (
                <tr key={user.id} className="group hover:bg-slate-50 transition-colors">
                  <td className="py-3">{user.name}</td>
                  <td className="py-3 text-slate-500">{user.email}</td>
                  <td className="py-3">
                    <span className="px-2 py-1 rounded-full text-xs font-medium bg-indigo-100 text-indigo-700 uppercase">
                      {user.role}
                    </span>
                  </td>
                  <td className="py-3 text-right">
                    <button 
                      onClick={() => handleDelete(user.id)}
                      className="text-red-500 hover:text-red-700 text-sm font-medium"
                    >
                      Remove
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </Card>
    </div>
  );
};

export default UserManagement;
