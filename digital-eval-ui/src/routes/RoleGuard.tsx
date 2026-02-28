import { Navigate } from "react-router-dom";
import { useAuth } from "../hooks/useAuth";

export const RoleGuard = ({
  allowedRoles,
  children,
}: {
  allowedRoles: string[];
  children: JSX.Element;
}) => {
  const { user } = useAuth();

  if (!user) return <Navigate to="/login" replace />;
  if (user.role && !allowedRoles.includes(user.role)) return <Navigate to="/unauthorized" replace />;

  return children;
};