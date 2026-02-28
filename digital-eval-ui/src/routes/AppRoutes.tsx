import { Routes, Route, Navigate } from "react-router-dom";

import Login from "../pages/auth/Login";

import AuthorityLayout from "../layouts/AuthorityLayout";
import ExaminerLayout from "../layouts/ExaminerLayout";
import EvaluatorLayout from "../layouts/EvaluatorLayout";
import StudentLayout from "../layouts/StudentLayout";
import AdminLayout from "../layouts/AdminLayout";

import AuthorityDashboard from "../pages/dashboard/AuthorityDashboard";
import ExaminerDashboard from "../pages/dashboard/ExaminerDashboard";
import EvaluatorDashboard from "../pages/dashboard/EvaluatorDashboard";
import StudentDashboard from "../pages/dashboard/StudentDashboard";
import AdminDashboard from "@/pages/dashboard/AdminDashboard";

import PendingRequests from "../pages/authority/PendingRequests";
import ApproveRequest from "../pages/authority/ApproveRequest";
import ReleaseResults from "../pages/authority/ReleaseResults";

import UploadScripts from "../pages/examiner/UploadScripts";
import UploadHistory from "../pages/examiner/UploadHistory";

import AssignedScripts from "../pages/evaluator/AssignedScripts";
import RequestEvaluation from "../pages/evaluator/RequestEvaluation";
import RequestHistory from "../pages/evaluator/RequestHistory";
import EvaluationPage from "../pages/evaluator/EvaluationPage";

import StudentResults from "../pages/student/StudentResults";

import UserManagement from "../pages/admin/UserManagement";

import DashboardRoute from "./DashboardRoute";
import { PrivateRoute } from "./PrivateRoute";
import { RoleGuard } from "./RoleGuard";

const AppRoutes = () => {
    return (
        <Routes>

            {/* Root Redirect */}
            <Route path="/" element={<Navigate to="/login" replace />} />

            {/* Login */}
            <Route path="/login" element={<Login />} />

            {/* Dashboard */}
            <Route path="/dashboard">
                <Route index element={
                    <PrivateRoute>
                        <DashboardRoute />
                    </PrivateRoute>
                } />

                {/* AUTHORITY */}
                <Route
                    path="authority"
                    element={
                        <PrivateRoute>
                            <RoleGuard allowedRoles={["authority"]}>
                                <AuthorityLayout />
                            </RoleGuard>
                        </PrivateRoute>
                    }
                >
                    <Route index element={<AuthorityDashboard />} />
                    <Route path="requests" element={<PendingRequests />} />
                    <Route path="approve" element={<ApproveRequest />} />
                    <Route path="release" element={<ReleaseResults />} />
                </Route>

                {/* EXAMINER */}
                <Route
                    path="examiner"
                    element={
                        <PrivateRoute>
                            <RoleGuard allowedRoles={["examiner"]}>
                                <ExaminerLayout />
                            </RoleGuard>
                        </PrivateRoute>
                    }
                >
                    <Route index element={<ExaminerDashboard />} />
                    <Route path="upload" element={<UploadScripts />} />
                    <Route path="history" element={<UploadHistory />} />
                </Route>

                {/* EVALUATOR */}
                <Route
                    path="evaluator"
                    element={
                        <PrivateRoute>
                            <RoleGuard allowedRoles={["evaluator"]}>
                                <EvaluatorLayout />
                            </RoleGuard>
                        </PrivateRoute>
                    }
                >
                    <Route index element={<EvaluatorDashboard />} />
                    <Route path="assigned" element={<AssignedScripts />} />
                    <Route path="request" element={<RequestEvaluation />} />
                    <Route path="requests" element={<RequestHistory />} />
                    <Route path="evaluate/:scriptId" element={<EvaluationPage />} />
                </Route>

                {/* STUDENT */}
                <Route
                    path="student"
                    element={
                        <PrivateRoute>
                            <RoleGuard allowedRoles={["student"]}>
                                <StudentLayout />
                            </RoleGuard>
                        </PrivateRoute>
                    }
                >
                    <Route index element={<StudentDashboard />} />
                    <Route path="results" element={<StudentResults />} />
                </Route>

                {/* ADMIN */}
                <Route
                    path="admin"
                    element={
                        <PrivateRoute>
                            <RoleGuard allowedRoles={["admin"]}>
                                <AdminLayout />
                            </RoleGuard>
                        </PrivateRoute>
                    }
                >
                    <Route index element={<AdminDashboard />} />
                    <Route path="users" element={<UserManagement />} />
                </Route>
            </Route>

            {/* 404 */}
            <Route path="*" element={<Navigate to="/login" replace />} />

        </Routes>
    );
};

export default AppRoutes;