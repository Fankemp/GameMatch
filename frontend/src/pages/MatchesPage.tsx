import { useState, useEffect } from "react";
import { getMatches } from "../api/matches";
import type { MatchWithUser } from "../types";
import { getRankColor } from "../utils/rankColors";
import { ROLE_ICONS } from "../App"; // или вынеси в константы

function MatchesPage() {
    const [matches, setMatches] = useState<MatchWithUser[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        getMatches()
            .then(setMatches)
            .catch(console.error)
            .finally(() => setLoading(false));
    }, []);

    if (loading) return <div className="page">Загрузка...</div>;

    return (
        <div className="page">
            <div className="page-title">Матчи <span style={{ color: "var(--teal)", fontSize: 20 }}>{matches.length}</span></div>
            {matches.map(m => (
                <div key={m.match.id} className="match-card">
                    <div className="match-card-av"><img src={m.avatar_url || `https://api.dicebear.com/7.x/avataaars/svg?seed=${m.user.username}`} alt={m.user.username} /></div>
                    <div style={{ flex: 1 }}>
                        <div className="mc-name">{m.user.username}</div>
                        <div className="mc-meta">
                            <span style={{ color: getRankColor(m.rank as any) }}>● {m.rank}</span> · {ROLE_ICONS[m.role as any]} {m.role}
                        </div>
                    </div>
                    <div style={{ textAlign: "right", display: "flex", flexDirection: "column", alignItems: "flex-end", gap: 8 }}>
                        <div className="mc-time">{new Date(m.match.created_at).toLocaleDateString()}</div>
                        {m.user.discord && <button className="discord-btn" onClick={() => navigator.clipboard.writeText(m.user.discord!)}>Discord →</button>}
                    </div>
                </div>
            ))}
        </div>
    );
}