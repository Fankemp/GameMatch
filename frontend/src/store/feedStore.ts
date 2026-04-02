import { create } from "zustand";
import { getFeed } from "../api/feed";
import { swipe } from "../api/swipes";
import type { FeedCard, FeedFilters, GameID, Match } from "../types";


interface FeedState {
    gameId: GameID;
    cards: FeedCard[];
    filters: FeedFilters;
    isLoading: boolean;
    error: string | null;
    newMatch: Match | null;

    setGame: (gameId: GameID) => void;
    setFilters: (f: FeedFilters) => void;
    loadFeed: () => Promise<void>;
    swipe: (card: FeedCard, action: "like" | "dislike") => Promise<void>;
    clearMatch: () => void;
}

export const useFeedStore = create<FeedState>((set, get) => ({
    gameId: "valorant",
    cards: [],
    filters: {},
    isLoading: false,
    error: null,
    newMatch: null,

    setGame: (gameId) => {
        set({ gameId, cards: [] });
        get().loadFeed();
    },

    setFilters: (filters) => set({ filters }),

    loadFeed: async () => {
        set({ isLoading: true, error: null });
        try {
            const { gameId, filters } = get();
            const cards = await getFeed(gameId, filters);
            set({ cards });
        } finally {
            set({ isLoading: false });
        }
    },

    swipe: async (card, action) => {
        set((s) => ({ cards: s.cards.filter((c) => c.id !== card.id) }));
        try {
            const result = await swipe(card.id, action);
            if (result.matched && result.match) {
                set({ newMatch: result.match });
            }
        } catch (err) {
            console.error("Swipe failed", err);
        }
    },

    clearMatch: () => set({ newMatch: null }),
}));