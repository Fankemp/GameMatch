import { client } from "./client";
import type { User, LoginInput, RegisterInput } from "../types";

interface AuthResponse {
    token: string;
    user: User;
}

// POST /api/v1/auth/register
export const register = (data: RegisterInput) =>
    client.post<AuthResponse>("/auth/register", data).then((r) => r.data);

// POST /api/v1/auth/login
export const login = (data: LoginInput) =>
    client.post<AuthResponse>("/auth/login", data).then((r) => r.data);

// GET /api/v1/auth/me  — пока нет JWT
export const getMe = () =>
    client.get<User>("/auth/me").then((r) => r.data);