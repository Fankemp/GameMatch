import { create } from "zustand";
import { getFeed } from "../api/feed";
import { swipe } from "../api/swipes";
import type { FeedCard, FeedFilters, GameID, Match } from "../types";
// временные моки
const MOCK_CARDS: FeedCard[] = [
    { id:1, user_id:10, game_id:"valorant", rank:"Gold", role:"Duelist", description:"Ищу дуо", is_active:true, created_at:"", updated_at:"", username:"AlexFPS", age:23, language:"RU", region:"EU Central", avatar_url:"https://api.dicebear.com/7.x/avataaars/svg?seed=Alex" },
    { id:2, user_id:11, game_id:"valorant", rank:"Platinum", role:"Controller", description:"Looking for duo", is_active:true, created_at:"", updated_at:"", username:"NightOwl", age:19, language:"EN", region:"EU West", avatar_url:"https://api.dicebear.com/7.x/avataaars/svg?seed=Night" },
];

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
        } catch {
            // fallback на моки, бэк ис пока анреди
            set({ cards: MOCK_CARDS });
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
        } catch {
            // имитация матча для демо тоже пока бэк анреди
            if (action === "like" && Math.random() < 0.4) {
                const fakeMatch: Match = {
                    id: Date.now(),
                    user_a_id: 0,
                    user_b_id: card.user_id,
                    game_id: card.game_id,
                    created_at: new Date().toISOString(),
                };
                set({ newMatch: fakeMatch });
            }
        }
    },

    clearMatch: () => set({ newMatch: null }),
}));