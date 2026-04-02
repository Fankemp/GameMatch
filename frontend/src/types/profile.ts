export interface Profile {
    id: number;
    user_id: number;
    bio?: string;
    avatar_url?: string;
    created_at: string;
    updated_at: string;
}

export interface CreateProfileInput {
    bio?: string;
    avatar_url?: string;
}