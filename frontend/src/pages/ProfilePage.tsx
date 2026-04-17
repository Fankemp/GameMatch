import { useState, useEffect } from "react";
import { useAuthStore } from "../store/authStore";
import { getProfile, createProfile, updateProfile } from "../api/profile";
import type { Profile } from "../types";

const REGIONS = ["EU Central", "EU West", "NA East", "NA West", "CIS", "Asia"];
const LANGS = ["RU", "EN", "DE", "KZ", "TR", "PL"];

export default function ProfilePage() {
    const { user, fetchMe } = useAuthStore();
    const [profile, setProfile] = useState<Profile | null>(null);
    const [bio, setBio] = useState("");
    const [avatarUrl, setAvatarUrl] = useState("");
    const [loading, setLoading] = useState(true);
    const [saving, setSaving] = useState(false);
    const [msg, setMsg] = useState("");

    useEffect(() => {
        getProfile()
            .then((p) => {
                setProfile(p);
                setBio(p.bio ?? "");
                setAvatarUrl(p.avatar_url ?? "");
            })
            .catch(() => {})
            .finally(() => setLoading(false));
    }, []);

    const handleSave = async () => {
        setSaving(true);
        setMsg("");
        try {
            const data = { bio, avatar_url: avatarUrl };
            if (profile) {
                const updated = await updateProfile(data);
                setProfile(updated);
            } else {
                const created = await createProfile(data);
                setProfile(created);
            }
            setMsg("Сохранено!");
        } catch {
            setMsg("Ошибка сохранения");
        } finally {
            setSaving(false);
        }
    };

    if (loading) return <div className="page">Загрузка...</div>;

    return (
        <div className="page">
            <div className="page-title">Профиль</div>

            <div className="profile-card">
                <div className="profile-avatar">
                    <img
                        src={avatarUrl || `https://api.dicebear.com/7.x/avataaars/svg?seed=${user?.username}`}
                        alt="avatar"
                    />
                </div>
                <div className="profile-info">
                    <div className="profile-name">{user?.username}</div>
                    <div className="profile-meta">
                        {user?.age} лет · {user?.language} · {user?.region}
                    </div>
                    {user?.discord && <div className="profile-contact">Discord: {user.discord}</div>}
                    {user?.telegram && <div className="profile-contact">Telegram: {user.telegram}</div>}
                </div>
            </div>

            <div className="profile-form">
                <div className="field">
                    <label>Аватар URL</label>
                    <input
                        placeholder="https://..."
                        value={avatarUrl}
                        onChange={(e) => setAvatarUrl(e.target.value)}
                    />
                </div>
                <div className="field">
                    <label>О себе</label>
                    <textarea
                        className="profile-bio"
                        placeholder="Расскажи о себе..."
                        value={bio}
                        onChange={(e) => setBio(e.target.value)}
                        rows={4}
                    />
                </div>
                <button className="auth-btn" onClick={handleSave} disabled={saving}>
                    {saving ? "Сохраняем..." : profile ? "Обновить профиль" : "Создать профиль"}
                </button>
                {msg && <div className={`profile-msg ${msg.includes("Ошибка") ? "err" : ""}`}>{msg}</div>}
            </div>
        </div>
    );
}
