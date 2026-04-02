import { client } from "./client";
import type { SwipeAction, Match } from "../types";

//TODO: бэк пустой
export const swipe = (targetCardId: number, action: SwipeAction) =>
    client
        .post<{ matched: boolean; match?: Match }>("/swipes", {
            target_card_id: targetCardId,
            action,
        })
        .then((r) => r.data);