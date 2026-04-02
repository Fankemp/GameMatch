export type GameID = "valorant" | "cs2" | "lol" | "dota2";
export type Rank =
    | "Iron" | "Bronze" | "Silver" | "Gold"
    | "Platinum" | "Diamond" | "Ascendant" | "Immortal" | "Radiant";
export type Role = "Duelist" | "Initiator" | "Controller" | "Sentinel";

export interface GameCard {
    id: number;
    user_id: number;
    game_id: GameID;
    rank: Rank;
    role: Role;
    description: string;
    is_active: boolean;
    created_at: string;
    updated_at: string;
}

export interface FeedCard extends GameCard {
    username: string;
    age: number;
    language: string;
    region: string;
    discord?: string;
    telegram?: string;
    avatar_url?: string;
    bio?: string;
}

export interface CreateCardInput {
    game_id: GameID;
    rank: Rank;
    role: Role;
    description: string;
}

export interface FeedFilters {
    rank_min?: Rank;
    rank_max?: Rank;
    role?: Role;
    language?: string;
    region?: string;
    age_min?: number;
    age_max?: number;
}