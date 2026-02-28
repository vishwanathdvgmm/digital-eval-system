import React from "react";
import { Link } from "react-router-dom";
import Card from "../../components/Card";

const EvaluatorDashboard: React.FC = () => {
  return (
    <div className="p-6 space-y-6">
      <div className="bg-green-600 text-white p-6 rounded-lg shadow-md">
        <h1 className="text-3xl font-bold">Welcome, Evaluator</h1>
        <p className="mt-2 opacity-90">Manage your assigned scripts and evaluation requests.</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <Link to="/dashboard/evaluator/assigned" className="block">
          <Card className="h-full hover:shadow-md transition-shadow cursor-pointer border-l-4 border-l-blue-500">
            <h2 className="text-xl font-semibold mb-2">Assigned Scripts</h2>
            <p className="text-gray-600">View and evaluate scripts assigned to you.</p>
          </Card>
        </Link>

        <Link to="/dashboard/evaluator/request" className="block">
          <Card className="h-full hover:shadow-md transition-shadow cursor-pointer border-l-4 border-l-yellow-500">
            <h2 className="text-xl font-semibold mb-2">Request Evaluation</h2>
            <p className="text-gray-600">Request new scripts for evaluation.</p>
          </Card>
        </Link>

        <Link to="/dashboard/evaluator/requests" className="block">
          <Card className="h-full hover:shadow-md transition-shadow cursor-pointer border-l-4 border-l-gray-500">
            <h2 className="text-xl font-semibold mb-2">Request History</h2>
            <p className="text-gray-600">Track status of your evaluation requests.</p>
          </Card>
        </Link>
      </div>
    </div>
  );
};

export default EvaluatorDashboard;