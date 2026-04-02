import { client } from "./client";
import type { Profile, CreateProfileInput } from "../types";

export const createProfile = (data: CreateProfileInput) =>
    client.post<Profile>("/profile", data).then((r) => r.data);

export const getProfile = () =>
    client.get<Profile>("/profile").then((r) => r.data);

export const updateProfile = (data: CreateProfileInput) =>
    client.put<Profile>("/profile", data).then((r) => r.data);