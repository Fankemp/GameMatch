import type { Rank } from "../types";

export const RANK_COLORS: Record<Rank, string> = {
    Iron: "#8b9bb4",
    Bronze: "#cd7f32",
    Silver: "#c0c0c0",
    Gold: "#ffd700",
    Platinum: "#4fc3b0",
    Diamond: "#a96cdb",
    Ascendant: "#00d4aa",
    Immortal: "#ff4d6d",
    Radiant: "#ffe066",
};

export const getRankColor = (rank: Rank) => RANK_COLORS[rank] ?? "#888";