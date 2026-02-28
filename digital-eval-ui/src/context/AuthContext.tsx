import { createContext, useContext } from "react";
import { User } from "../types/auth";

export interface AuthState {
  user: User | null;
  loading: boolean;
  login: (loginId: string, password: string) => Promise<User>;
  logout: () => Promise<void>;
}

export const AuthContext = createContext<AuthState>({
  user: null,
  loading: true,
  login: async () => {
    throw new Error("Context not provided");
  },
  logout: async () => { },
});

export const useAuth = () => useContext(AuthContext);