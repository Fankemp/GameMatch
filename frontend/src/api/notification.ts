import { client } from "./client";
import type { Notification } from "../types";

export const getNotifications = () =>
    client.get<Notification[]>("/notifications").then((r) => r.data);

export const markRead = (id: number) =>
    client.put(`/notifications/${id}/read`).then((r) => r.data);