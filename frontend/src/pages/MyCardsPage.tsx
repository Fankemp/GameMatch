import { useState, useEffect } from "react";
import { getMyCards, createCard, updateCard, deleteCard } from "../api/cards";
import { getRankColor } from "../utils/rankColors";
import { ROLE_ICONS } from "../App";
import type { GameCard, GameID, Rank, Role, CreateCardInput } from "../types";

const GAMES: { id: GameID; label: string }[] = [
    { id: "valorant", label: "Valorant" },
    { id: "cs2", label: "CS2" },
    { id: "lol", label: "LoL" },
    { id: "dota2", label: "Dota 2" },
];

const RANKS: Rank[] = ["Iron", "Bronze", "Silver", "Gold", "Platinum", "Diamond", "Ascendant", "Immortal", "Radiant"];
const ROLES: Role[] = ["Duelist", "Initiator", "Controller", "Sentinel"];

export default function MyCardsPage() {
    const [cards, setCards] = useState<GameCard[]>([]);
    const [loading, setLoading] = useState(true);
    const [showForm, setShowForm] = useState(false);
    const [form, setForm] = useState<CreateCardInput>({
        game_id: "valorant" as GameID,
        rank: "Gold" as Rank,
        role: "Duelist" as Role,
        description: "",
    });
    const [saving, setSaving] = useState(false);

    const loadCards = () => {
        setLoading(true);
        getMyCards()
            .then((data) => setCards(data ?? []))
            .catch(() => setCards([]))
            .finally(() => setLoading(false));
    };

    useEffect(() => { loadCards(); }, []);

    const handleCreate = async () => {
        if (!form.description.trim()) return;
        setSaving(true);
        try {
            await createCard(form);
            setShowForm(false);
            setForm({ game_id: "valorant" as GameID, rank: "Gold" as Rank, role: "Duelist" as Role, description: "" });
            loadCards();
        } catch (e) {
            console.error(e);
        } finally {
            setSaving(false);
        }
    };

    const handleToggle = async (card: GameCard) => {
        try {
            await updateCard(card.id, { is_active: !card.is_active } as any);
            loadCards();
        } catch (e) {
            console.error(e);
        }
    };

    const handleDelete = async (id: number) => {
        try {
            await deleteCard(id);
            loadCards();
        } catch (e) {
            console.error(e);
        }
    };

    if (loading) return <div className="page">Загрузка...</div>;

    return (
        <div className="page">
            <div className="page-title" style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                Мои карточки <span style={{ color: "var(--teal)", fontSize: 20 }}>{cards.length}</span>
                <button className="add-card-btn" onClick={() => setShowForm(!showForm)}>
                    {showForm ? "Отмена" : "+ Создать"}
                </button>
            </div>

            {showForm && (
                <div className="card-form">
                    <div className="field-row">
                        <div className="field">
                            <label>Игра</label>
                            <select value={form.game_id} onChange={(e) => setForm({ ...form, game_id: e.target.value as GameID })}>
                                {GAMES.map((g) => <option key={g.id} value={g.id}>{g.label}</option>)}
                            </select>
                        </div>
                        <div className="field">
                            <label>Ранг</label>
                            <select value={form.rank} onChange={(e) => setForm({ ...form, rank: e.target.value as Rank })}>
                                {RANKS.map((r) => <option key={r}>{r}</option>)}
                            </select>
                        </div>
                    </div>
                    <div className="field">
                        <label>Роль</label>
                        <select value={form.role} onChange={(e) => setForm({ ...form, role: e.target.value as Role })}>
                            {ROLES.map((r) => <option key={r}>{r}</option>)}
                        </select>
                    </div>
                    <div className="field">
                        <label>Описание</label>
                        <textarea
                            className="profile-bio"
                            placeholder="Ищу тимейта на ранкед..."
                            value={form.description}
                            onChange={(e) => setForm({ ...form, description: e.target.value })}
                            rows={3}
                        />
                    </div>
                    <button className="auth-btn" onClick={handleCreate} disabled={saving}>
                        {saving ? "Создаём..." : "Создать карточку"}
                    </button>
                </div>
            )}

            {cards.length === 0 && !showForm && (
                <div className="feed-status">У тебя пока нет карточек. Создай первую!</div>
            )}

            <div className="my-cards-list">
                {cards.map((card) => (
                    <div key={card.id} className={`my-card ${!card.is_active ? "inactive" : ""}`}>
                        <div className="my-card-header">
                            <span className="my-card-game">{GAMES.find((g) => g.id === card.game_id)?.label ?? card.game_id}</span>
                            <span className={`my-card-status ${card.is_active ? "active" : ""}`}>
                                {card.is_active ? "Активна" : "Неактивна"}
                            </span>
                        </div>
                        <div className="my-card-body">
                            <span className="my-card-rank" style={{ color: getRankColor(card.rank as Rank) }}>
                                {card.rank}
                            </span>
                            <span className="my-card-role">{ROLE_ICONS[card.role as Role]} {card.role}</span>
                        </div>
                        <p className="my-card-desc">{card.description}</p>
                        <div className="my-card-actions">
                            <button className="toggle-btn" onClick={() => handleToggle(card)}>
                                {card.is_active ? "Деактивировать" : "Активировать"}
                            </button>
                            <button className="delete-btn" onClick={() => handleDelete(card.id)}>
                                Удалить
                            </button>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
}
