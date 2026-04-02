import { client } from "./client";
import type { FeedCard, FeedFilters } from "../types";

// TODO: бэк пустой
export const getFeed = (gameId: string, filters?: FeedFilters) =>
    client
        .get<FeedCard[]>(`/feed/${gameId}`, { params: filters })
        .then((r) => r.data);