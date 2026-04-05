import { useState, useEffect } from "react";
import { getMatches } from "../api/matches";
import type { MatchWithUser } from "../types";
import { getRankColor } from "../utils/rankColors";
import { ROLE_ICONS } from "../App";

function MatchesPage() {
    const [matches, setMatches] = useState<MatchWithUser[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");

    useEffect(() => {
        getMatches()
            .then((data) => setMatches(data ?? []))
            .catch(() => {
                setMatches([]);
                setError("Не удалось загрузить матчи");
            })
            .finally(() => setLoading(false));
    }, []);

    if (loading) return <div className="page">Загрузка...</div>;

    return (
        <div className="page">
            <div className="page-title">Матчи <span style={{ color: "var(--teal)", fontSize: 20 }}>{matches.length}</span></div>
            {error && <div style={{ color: "var(--accent)", fontSize: 14, marginBottom: 12 }}>{error}</div>}
            {matches.length === 0 && !error && <div className="feed-status">Пока нет матчей</div>}
            {matches.map(m => (
                <div key={m.match?.id ?? Math.random()} className="match-card">
                    <div className="match-card-av">
                        <img src={m.avatar_url || `https://api.dicebear.com/7.x/avataaars/svg?seed=${m.user?.username ?? "anon"}`} alt={m.user?.username ?? ""} />
                    </div>
                    <div style={{ flex: 1 }}>
                        <div className="mc-name">{m.user?.username ?? "Unknown"}</div>
                        <div className="mc-meta">
                            <span style={{ color: getRankColor(m.rank as any) }}>● {m.rank ?? "?"}</span> · {ROLE_ICONS[m.role as keyof typeof ROLE_ICONS] ?? ""} {m.role ?? ""}
                        </div>
                    </div>
                    <div style={{ textAlign: "right", display: "flex", flexDirection: "column", alignItems: "flex-end", gap: 8 }}>
                        <div className="mc-time">{m.match?.created_at ? new Date(m.match.created_at).toLocaleDateString() : ""}</div>
                        {m.user?.discord && <button className="discord-btn" onClick={() => navigator.clipboard.writeText(m.user.discord!)}>Discord →</button>}
                    </div>
                </div>
            ))}
        </div>
    );
}

export default MatchesPage;
