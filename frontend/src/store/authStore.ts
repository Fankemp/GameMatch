import { create } from "zustand";
import { getMe } from "../api/auth";
import type { User } from "../types";
import { getToken, setToken, removeToken } from "../utils/token";

interface AuthState {
    token: string | null;
    user: User | null;
    isAuthed: boolean;
    isLoading: boolean;

    setAuth: (token: string, user: User) => void;
    logout: () => void;
    fetchMe: () => Promise<void>;
}

export const useAuthStore = create<AuthState>((set) => ({
    token: getToken(),
    user: null,
    isAuthed: !!getToken(),
    isLoading: false,

    setAuth: (token, user) => {
        setToken(token);
        set({ token, user, isAuthed: true });
    },

    logout: () => {
        removeToken();
        set({ token: null, user: null, isAuthed: false });
    },

    fetchMe: async () => {
        set({ isLoading: true });
        try {
            const user = await getMe();
            set({ user, isAuthed: true });
        } catch {
            removeToken();
            set({ token: null, user: null, isAuthed: false });
        } finally {
            set({ isLoading: false });
        }
    },
}));