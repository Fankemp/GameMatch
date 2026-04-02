export interface Notification {
    id: number;
    user_id: number;
    type: string;
    content: string;
    read: boolean;
    created_at: string;
}