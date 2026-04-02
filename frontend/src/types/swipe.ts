export type SwipeAction = "like" | "dislike";

export interface Swipe {
    id: number;
    user_id: number;
    target_card_id: number;
    action: SwipeAction;
    created_at: string;
}