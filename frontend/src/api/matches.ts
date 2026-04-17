import { client } from "./client";
import type { MatchWithUser } from "../types";

export const getMatches = () =>
    client.get<MatchWithUser[]>("/matches").then((r) => r.data);
