import React from "react";
import { Navigate } from "react-router-dom";
import { useAuth } from "../hooks/useAuth";
import AuthorityDashboard from "../pages/dashboard/AuthorityDashboard";
import EvaluatorDashboard from "../pages/dashboard/EvaluatorDashboard";
import ExaminerDashboard from "../pages/dashboard/ExaminerDashboard";
import StudentDashboard from "../pages/dashboard/StudentDashboard";
import AdminDashboard from "@/pages/dashboard/AdminDashboard";

const DashboardRoute: React.FC = () => {
  const { user } = useAuth();

  if (!user) return <Navigate to="/auth/login" replace />;

  switch (user.role) {
    case "authority":
      return <AuthorityDashboard />;
    case "evaluator":
      return <EvaluatorDashboard />;
    case "examiner":
      return <ExaminerDashboard />;
    case "student":
      return <StudentDashboard />;
    case "admin":
      return <AdminDashboard />;
    default:
      return <div>Unknown role</div>;
  }
};

export default DashboardRoute;