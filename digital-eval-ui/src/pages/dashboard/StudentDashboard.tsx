import React from "react";
import { Link } from "react-router-dom";
import Card from "../../components/Card";

const StudentDashboard: React.FC = () => {
  return (
    <div className="p-6 space-y-6">
      <div className="bg-purple-600 text-white p-6 rounded-lg shadow-md">
        <h1 className="text-3xl font-bold">Welcome, Student</h1>
        <p className="mt-2 opacity-90">Access your academic results and performance reports.</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Link to="/dashboard/student/results" className="block">
          <Card className="h-full hover:shadow-md transition-shadow cursor-pointer border-l-4 border-l-blue-500">
            <h2 className="text-xl font-semibold mb-2">View Results</h2>
            <p className="text-gray-600">Check your latest semester results.</p>
          </Card>
        </Link>
      </div>
    </div>
  );
};

export default StudentDashboard;