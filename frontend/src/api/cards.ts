import { client } from "./client";
import type { GameCard, CreateCardInput } from "../types";

// POST /api/v1/cards  — TODO: бэк пустой
export const createCard = (data: CreateCardInput) =>
    client.post<GameCard>("/cards", data).then((r) => r.data);

// GET /api/v1/cards
export const getMyCards = () =>
    client.get<GameCard[]>("/cards").then((r) => r.data);

// PUT /api/v1/cards/:id
export const updateCard = (id: number, data: Partial<CreateCardInput>) =>
    client.put<GameCard>(`/cards/${id}`, data).then((r) => r.data);

// DELETE /api/v1/cards/:id
export const deleteCard = (id: number) =>
    client.delete(`/cards/${id}`);