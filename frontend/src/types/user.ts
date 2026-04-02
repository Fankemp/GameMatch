export interface User {
    id: number;
    username: string;
    email: string;
    age: number;
    language: string;
    discord?: string;
    telegram?: string;
    region: string;
    created_at: string;
    updated_at: string;
}

export interface LoginInput {
    email: string;
    password: string;
}

export interface RegisterInput {
    username: string;
    email: string;
    password: string;
    age: number;
    language: string;
    region: string;
    discord?: string;
    telegram?: string;
}