import React from "react";
import { Link } from "react-router-dom";
import Card from "../../components/Card";

const ExaminerDashboard: React.FC = () => {
  return (
    <div className="p-6 space-y-6">
      <div className="bg-indigo-600 text-white p-6 rounded-lg shadow-md">
        <h1 className="text-3xl font-bold">Welcome, Examiner</h1>
        <p className="mt-2 opacity-90">Upload answer scripts and track your submissions.</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Link to="/dashboard/examiner/upload" className="block">
          <Card className="h-full hover:shadow-md transition-shadow cursor-pointer border-l-4 border-l-blue-500">
            <h2 className="text-xl font-semibold mb-2">Upload Scripts</h2>
            <p className="text-gray-600">Upload new answer scripts for evaluation.</p>
          </Card>
        </Link>

        <Link to="/dashboard/examiner/history" className="block">
          <Card className="h-full hover:shadow-md transition-shadow cursor-pointer border-l-4 border-l-gray-500">
            <h2 className="text-xl font-semibold mb-2">Upload History</h2>
            <p className="text-gray-600">View status of previously uploaded scripts.</p>
          </Card>
        </Link>
      </div>
    </div>
  );
};

export default ExaminerDashboard;