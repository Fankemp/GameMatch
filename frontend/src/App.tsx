import { useEffect, useState } from "react";
import { useAuthStore } from "./store/authStore";
import { useFeedStore } from "./store/feedStore";
import { LoginPage, RegisterPage } from "./pages/AuthPages";
import MatchesPage from "./pages/MatchesPage";
import MyCardsPage from "./pages/MyCardsPage";
import ProfilePage from "./pages/ProfilePage";
import { getRankColor } from "./utils/rankColors";
import { LANG_FLAGS } from "./utils/langFlags";
import type { FeedCard, GameID, Role } from "./types";
import "./App.css";

export const ROLE_ICONS: Record<Role, string> = {
  Duelist: "⚔️",
  Initiator: "🔍",
  Controller: "💨",
  Sentinel: "🛡️",
};

const GAMES: { id: GameID; label: string; icon: string }[] = [
  { id: "valorant", label: "Valorant", icon: "🎯" },
  { id: "cs2", label: "CS2", icon: "🔫" },
  { id: "lol", label: "LoL", icon: "🧙" },
  { id: "dota2", label: "Dota 2", icon: "⚡" },
];

type Page = "feed" | "matches" | "cards" | "profile";

function App() {
  const { isAuthed, isLoading, user, fetchMe, logout } = useAuthStore();
  const [authMode, setAuthMode] = useState<"login" | "register">("login");
  const [page, setPage] = useState<Page>("feed");

  useEffect(() => {
    if (isAuthed && !user) fetchMe();
  }, [isAuthed, user, fetchMe]);

  if (isLoading) {
    return (
      <div className="loader-wrap">
        <div className="loader-text">
          Game<span>Match</span>
        </div>
      </div>
    );
  }

  if (!isAuthed) {
    return authMode === "login" ? (
      <LoginPage onSwitch={() => setAuthMode("register")} />
    ) : (
      <RegisterPage onSwitch={() => setAuthMode("login")} />
    );
  }

  return (
    <div className="app">
      <header className="topbar">
        <div className="topbar-logo">
          Game<span>Match</span>
        </div>
        <nav className="topbar-nav">
          <button className={page === "feed" ? "active" : ""} onClick={() => setPage("feed")}>
            Feed
          </button>
          <button className={page === "cards" ? "active" : ""} onClick={() => setPage("cards")}>
            Cards
          </button>
          <button className={page === "matches" ? "active" : ""} onClick={() => setPage("matches")}>
            Matches
          </button>
          <button className={page === "profile" ? "active" : ""} onClick={() => setPage("profile")}>
            Profile
          </button>
        </nav>
        <div className="topbar-user">
          <span>{user?.username}</span>
          <button className="logout-btn" onClick={logout}>
            Выйти
          </button>
        </div>
      </header>

      <main className="main">
        {page === "feed" && <FeedPage />}
        {page === "cards" && <MyCardsPage />}
        {page === "matches" && <MatchesPage />}
        {page === "profile" && <ProfilePage />}
      </main>
    </div>
  );
}

function FeedPage() {
  const { gameId, cards, isLoading, setGame, swipe, newMatch, clearMatch, loadFeed } = useFeedStore();
  const [animCard, setAnimCard] = useState<{ id: number; dir: "left" | "right" } | null>(null);

  useEffect(() => {
    loadFeed();
  }, [loadFeed]);

  const handleSwipe = (card: FeedCard, action: "like" | "dislike") => {
    setAnimCard({ id: card.id, dir: action === "like" ? "right" : "left" });
    setTimeout(() => {
      swipe(card, action);
      setAnimCard(null);
    }, 300);
  };

  const topCard = cards[0];

  return (
    <div className="feed-page">
      <div className="game-tabs">
        {GAMES.map((g) => (
          <button
            key={g.id}
            className={`game-tab ${gameId === g.id ? "active" : ""}`}
            onClick={() => setGame(g.id)}
          >
            {g.icon} {g.label}
          </button>
        ))}
      </div>

      <div className="card-stack">
        {isLoading && <div className="feed-status">Загрузка...</div>}
        {!isLoading && !topCard && <div className="feed-status">Карточки закончились</div>}
        {topCard && (
          <div
            className={`swipe-card ${animCard?.id === topCard.id ? (animCard.dir === "left" ? "anim-left" : "anim-right") : ""}`}
          >
            <div className="card-header">
              <img
                className="card-avatar"
                src={topCard.avatar_url || `https://api.dicebear.com/7.x/avataaars/svg?seed=${topCard.username}`}
                alt={topCard.username}
              />
              <div>
                <div className="card-name">
                  {topCard.username}, {topCard.age}
                </div>
                <div className="card-meta">
                  {LANG_FLAGS[topCard.language] || ""} {topCard.region}
                </div>
              </div>
            </div>
            <div className="card-body">
              <div className="card-rank" style={{ color: getRankColor(topCard.rank) }}>
                {topCard.rank}
              </div>
              <div className="card-role">
                {ROLE_ICONS[topCard.role]} {topCard.role}
              </div>
              <p className="card-desc">{topCard.description}</p>
            </div>
            <div className="card-actions">
              <button className="btn-dislike" onClick={() => handleSwipe(topCard, "dislike")}>
                ✕
              </button>
              <button className="btn-like" onClick={() => handleSwipe(topCard, "like")}>
                ♥
              </button>
            </div>
          </div>
        )}
      </div>

      {newMatch && (
        <div className="match-overlay" onClick={clearMatch}>
          <div className="match-popup">
            <div className="match-title">It's a Match!</div>
            <p>Теперь вы можете связаться друг с другом</p>
            <button className="auth-btn" onClick={clearMatch}>
              Продолжить
            </button>
          </div>
        </div>
      )}
    </div>
  );
}

export default App;
