export interface User {
  id: number;
  user_id: string;      // e.g., examiner_1, evaluator_1, 4BD23AI104
  email: string;
  role: "examiner" | "evaluator" | "authority" | "student" | "admin";
  name: string;
  created_at: string;
  updated_at: string;
}

export interface LoginPayload {
  login: string;
  password: string;
}

export interface LoginResponse {
  access_token: string;
  token_type: string;
  user: User;
}