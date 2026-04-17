import axios, { AxiosError } from "axios";
import type { ApiError } from "../types";

const BASE_URL = import.meta.env.VITE_API_URL ?? "http://localhost:8080";

export const client = axios.create({
    baseURL: `${BASE_URL}/api/v1`,
    headers: { "Content-Type": "application/json" },
});

// JWT to every request
client.interceptors.request.use((config) => {
    const token = localStorage.getItem("token");
    if (token) config.headers.Authorization = `Bearer ${token}`;
    return config;
});

// Normalise error shape
client.interceptors.response.use(
    (res) => res,
    (err: AxiosError<ApiError>) => {
        if (err.response?.status === 401) {
            localStorage.removeItem("token");
            // Don't redirect, let auth store handle it
        }
        return Promise.reject(err);
    }
);

export default client;