import { User } from "./user";
import { Match } from "./match";

export interface Match {
    id: number;
    user_a_id: number;
    user_b_id: number;
    game_id: string;
    created_at: string;
}

export interface MatchWithUser {
    match: Match;
    user: User;
    rank: string;
    role: string;
    avatar_url?: string;
}