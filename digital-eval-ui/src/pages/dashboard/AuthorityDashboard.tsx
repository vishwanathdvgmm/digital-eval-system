import React from "react";
import { Link } from "react-router-dom";
import Card from "../../components/Card";

const AuthorityDashboard: React.FC = () => {
  return (
    <div className="p-6 space-y-6">
      <div className="bg-blue-600 text-white p-6 rounded-lg shadow-md">
        <h1 className="text-3xl font-bold">Welcome, Authority</h1>
        <p className="mt-2 opacity-90">Manage the evaluation process and release results.</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <Link to="/dashboard/authority/requests" className="block">
          <Card className="h-full hover:shadow-md transition-shadow cursor-pointer border-l-4 border-l-yellow-500">
            <h2 className="text-xl font-semibold mb-2">Pending Requests</h2>
            <p className="text-gray-600">Review and approve evaluator access requests.</p>
          </Card>
        </Link>

        <Link to="/dashboard/authority/approve" className="block">
          <Card className="h-full hover:shadow-md transition-shadow cursor-pointer border-l-4 border-l-green-500">
            <h2 className="text-xl font-semibold mb-2">Approve Actions</h2>
            <p className="text-gray-600">Finalize approvals and manage permissions.</p>
          </Card>
        </Link>

        <Link to="/dashboard/authority/release" className="block">
          <Card className="h-full hover:shadow-md transition-shadow cursor-pointer border-l-4 border-l-purple-500">
            <h2 className="text-xl font-semibold mb-2">Release Results</h2>
            <p className="text-gray-600">Publish semester results for students.</p>
          </Card>
        </Link>
      </div>
    </div>
  );
};

export default AuthorityDashboard;