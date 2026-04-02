import { client } from "./client";
import type { Match } from "../types";

// TODO: бэк пустой
export const getMatches = () =>
    client.get<Match[]>("/matches").then((r) => r.data);