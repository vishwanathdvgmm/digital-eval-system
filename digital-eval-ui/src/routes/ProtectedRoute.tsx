import { Navigate } from "react-router-dom";
import { useAuth } from "../hooks/useAuth";

export function ProtectedRoute({
  children,
  roles,
}: {
  children: JSX.Element;
  roles: string[];
}) {
  const { user, loading } = useAuth();

  if (loading) return null;

  if (!user) return <Navigate to="/login" replace />;

  if (!roles.includes(user.role)) return <Navigate to="/unauthorized" replace />;

  return children;
}