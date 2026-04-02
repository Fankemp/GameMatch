import { useState, useRef, useEffect } from "react";
import { useAuthStore } from "./store/authStore";
import { useFeedStore } from "./store/feedStore";
import { LoginPage, RegisterPage } from "./pages/AuthPages";
import { createProfile } from "./api/profile";
import { getMatches } from "./api/matches";
import type { FeedCard, GameID, Match, Rank, Role } from "./types";
import { getRankColor } from "./utils/rankColors";
import { LANG_FLAGS } from "./utils/langFlags";

// Constants
const RANKS: Rank[] = ["Iron","Bronze","Silver","Gold","Platinum","Diamond","Ascendant","Immortal","Radiant"];
const ROLES: Role[] = ["Duelist","Initiator","Controller","Sentinel"];
const ROLE_ICONS: Record<Role,string> = { Duelist:"⚔️", Initiator:"🧲", Controller:"🌫️", Sentinel:"🛡️" };
const GAMES: { id: GameID; label: string; icon: string }[] = [
    { id:"valorant", label:"Valorant", icon:"🔺" },
    { id:"cs2",      label:"CS2",      icon:"💣" },
    { id:"lol",      label:"LoL",      icon:"⚔" },
    { id:"dota2",    label:"Dota 2",   icon:"🐉" },
];
const LANGS = ["RU","EN","DE","KZ","TR","PL"];
const REGIONS = ["EU Central","EU West","NA East","NA West","CIS","Asia"];

// CSS
const APP_CSS = `
@import url('https://fonts.googleapis.com/css2?family=Rajdhani:wght@400;500;600;700&family=Barlow:wght@300;400;500;600&family=Barlow+Condensed:wght@400;600;700&display=swap');
*,*::before,*::after{box-sizing:border-box;margin:0;padding:0}
:root{
  --bg:#0a0b0e; --surface:#111318; --surface2:#181b22;
  --border:#1e2230; --border2:#252a38;
  --accent:#ff4655; --accent2:#ff6b78;
  --teal:#00d4aa; --gold:#ffd700;
  --text:#e8eaf0; --text2:#8b91a8; --text3:#4f5570;
  --glow:0 0 20px rgba(255,70,85,.35); --glow-teal:0 0 20px rgba(0,212,170,.35);
  font-family:'Barlow',sans-serif;
}
body{background:var(--bg);color:var(--text);min-height:100vh;overflow-x:hidden}
::-webkit-scrollbar{width:4px}
::-webkit-scrollbar-track{background:var(--surface)}
::-webkit-scrollbar-thumb{background:var(--border2);border-radius:2px}
.app{display:flex;flex-direction:column;min-height:100vh}
/* NAV */
.nav{display:flex;align-items:center;gap:12px;padding:0 24px;height:60px;background:var(--surface);border-bottom:1px solid var(--border);position:sticky;top:0;z-index:100}
.nav-logo{font-family:'Rajdhani',sans-serif;font-size:22px;font-weight:700;letter-spacing:1px;cursor:pointer}
.nav-logo span{color:var(--accent)}
.nav-game-select{display:flex;align-items:center;gap:8px;background:var(--surface2);border:1px solid var(--border2);border-radius:8px;padding:7px 14px;cursor:pointer;font-family:'Barlow Condensed',sans-serif;font-size:15px;font-weight:600;color:var(--text);transition:.2s;position:relative}
.nav-game-select:hover{border-color:var(--accent)}
.nav-spacer{flex:1}
.nav-btn{display:flex;align-items:center;gap:7px;background:var(--surface2);border:1px solid var(--border2);border-radius:8px;padding:7px 14px;cursor:pointer;font-size:13px;font-weight:500;color:var(--text2);transition:.2s}
.nav-btn:hover,.nav-btn.active{color:var(--text);border-color:var(--accent);color:var(--accent)}
.nav-badge{background:var(--accent);color:#fff;border-radius:10px;font-size:10px;font-weight:700;padding:1px 5px}
.nav-avatar{width:34px;height:34px;border-radius:50%;border:2px solid var(--accent);overflow:hidden;cursor:pointer;background:var(--surface2)}
.nav-avatar img{width:100%;height:100%;object-fit:cover}
/* DROPDOWN */
.dropdown{position:absolute;top:calc(100% + 8px);left:0;background:var(--surface);border:1px solid var(--border2);border-radius:10px;padding:6px;min-width:160px;z-index:200;box-shadow:0 8px 32px rgba(0,0,0,.5)}
.dropdown-opt{display:flex;align-items:center;gap:10px;padding:9px 12px;border-radius:7px;cursor:pointer;font-size:14px;font-weight:500;color:var(--text2);transition:.15s}
.dropdown-opt:hover{background:var(--surface2);color:var(--text)}
.dropdown-opt.active{color:var(--accent)}
/* LAYOUT */
.main{display:flex;flex:1;min-height:0}
/* SIDEBAR */
.sidebar{width:240px;min-width:240px;background:var(--surface);border-right:1px solid var(--border);padding:20px 14px;display:flex;flex-direction:column;gap:5px}
.sidebar-section{font-family:'Barlow Condensed',sans-serif;font-size:11px;font-weight:600;letter-spacing:2px;color:var(--text3);text-transform:uppercase;padding:0 8px;margin:10px 0 4px}
.sidebar-btn{display:flex;align-items:center;gap:10px;padding:9px 12px;border-radius:8px;cursor:pointer;font-size:14px;font-weight:500;color:var(--text2);transition:.18s;border:1px solid transparent}
.sidebar-btn:hover{background:var(--surface2);color:var(--text)}
.sidebar-btn.active{background:rgba(255,70,85,.1);border-color:rgba(255,70,85,.2);color:var(--accent)}
/* FILTER */
.filter-panel{background:var(--surface2);border:1px solid var(--border2);border-radius:10px;padding:16px;margin-top:10px;display:flex;flex-direction:column;gap:14px}
.filter-label{font-family:'Barlow Condensed',sans-serif;font-size:11px;font-weight:600;letter-spacing:1.5px;color:var(--text3);text-transform:uppercase;margin-bottom:6px}
.filter-select{width:100%;background:var(--surface);border:1px solid var(--border2);color:var(--text);border-radius:6px;padding:7px 10px;font-size:13px;font-family:'Barlow',sans-serif;cursor:pointer;outline:none}
.filter-select:focus{border-color:var(--accent)}
.chips{display:flex;flex-wrap:wrap;gap:5px}
.chip{padding:4px 10px;border-radius:20px;font-size:12px;font-weight:500;border:1px solid var(--border2);color:var(--text2);cursor:pointer;transition:.15s;background:var(--surface)}
.chip:hover{border-color:var(--text3);color:var(--text)}
.chip.on{background:rgba(255,70,85,.15);border-color:var(--accent);color:var(--accent)}
.chip.teal.on{background:rgba(0,212,170,.15);border-color:var(--teal);color:var(--teal)}
.range-row{display:flex;align-items:center;gap:8px}
.range-input{flex:1;background:var(--surface);border:1px solid var(--border2);color:var(--text);border-radius:6px;padding:6px 10px;font-size:13px;outline:none}
.range-input:focus{border-color:var(--accent)}
.apply-btn{background:var(--accent);color:#fff;border:none;border-radius:8px;padding:9px;font-size:13px;font-weight:600;cursor:pointer;font-family:'Barlow Condensed',sans-serif;letter-spacing:.5px;transition:.18s}
.apply-btn:hover{background:var(--accent2);box-shadow:var(--glow)}
/* FEED */
.feed-area{flex:1;display:flex;flex-direction:column;align-items:center;padding:32px 20px;gap:24px;overflow-y:auto}
.card-stack{position:relative;width:360px;height:520px}
.p-card{position:absolute;top:0;left:0;width:360px;height:520px;background:var(--surface);border:1px solid var(--border2);border-radius:18px;overflow:hidden;cursor:grab;user-select:none;transform-origin:bottom center;transition:box-shadow .2s}
.p-card:active{cursor:grabbing}
.p-card.behind{transform:scale(.94) translateY(16px);z-index:0;filter:brightness(.6);pointer-events:none}
.p-card.front{z-index:2}
.p-card-img{width:100%;height:320px;object-fit:cover;background:linear-gradient(135deg,#1a1f2e,#0f1117);display:flex;align-items:center;justify-content:center}
.p-card-img img{width:100%;height:100%;object-fit:cover}
.p-card-overlay{position:absolute;top:0;left:0;right:0;height:320px;background:linear-gradient(to bottom,transparent 40%,rgba(10,11,14,.95) 100%);pointer-events:none}
.p-card-body{padding:20px 22px}
.p-card-name{font-family:'Rajdhani',sans-serif;font-size:26px;font-weight:700;display:flex;align-items:center;gap:8px}
.p-card-mic{margin-left:auto;display:flex;align-items:center;gap:4px;font-size:11px;color:var(--teal);font-weight:600}
.rank-badge{display:inline-flex;align-items:center;gap:6px;font-family:'Barlow Condensed',sans-serif;font-size:14px;font-weight:700;letter-spacing:.5px;padding:4px 10px;border-radius:6px;margin-top:6px;border:1px solid}
.p-card-bio{margin-top:10px;font-size:13px;color:var(--text2);line-height:1.5;display:-webkit-box;-webkit-line-clamp:2;-webkit-box-orient:vertical;overflow:hidden}
.tags{display:flex;gap:6px;margin-top:10px;flex-wrap:wrap}
.tag{font-size:11px;padding:3px 9px;border-radius:4px;font-weight:600;border:1px solid var(--border2);color:var(--text2);font-family:'Barlow Condensed',sans-serif}
/* SWIPE LABELS */
.sw-like{position:absolute;top:24px;left:20px;z-index:10;background:rgba(0,212,170,.18);border:3px solid var(--teal);color:var(--teal);border-radius:10px;padding:6px 16px;font-family:'Rajdhani',sans-serif;font-size:28px;font-weight:700;letter-spacing:2px;opacity:0;transform:rotate(-12deg);pointer-events:none;transition:.1s}
.sw-nope{position:absolute;top:24px;right:20px;z-index:10;background:rgba(255,70,85,.18);border:3px solid var(--accent);color:var(--accent);border-radius:10px;padding:6px 16px;font-family:'Rajdhani',sans-serif;font-size:28px;font-weight:700;letter-spacing:2px;opacity:0;transform:rotate(12deg);pointer-events:none;transition:.1s}
/* ACTIONS */
.action-row{display:flex;align-items:center;gap:20px}
.btn-dis{width:62px;height:62px;border-radius:50%;border:2px solid var(--accent);background:rgba(255,70,85,.1);color:var(--accent);font-size:24px;display:grid;place-items:center;cursor:pointer;transition:.18s}
.btn-dis:hover{background:rgba(255,70,85,.22);box-shadow:var(--glow);transform:scale(1.07)}
.btn-lik{width:62px;height:62px;border-radius:50%;border:2px solid var(--teal);background:rgba(0,212,170,.1);color:var(--teal);font-size:24px;display:grid;place-items:center;cursor:pointer;transition:.18s}
.btn-lik:hover{background:rgba(0,212,170,.22);box-shadow:var(--glow-teal);transform:scale(1.07)}
.btn-skip{flex:1;background:transparent;border:1px solid var(--border2);color:var(--text3);font-size:12px;border-radius:8px;padding:8px;cursor:pointer;transition:.15s}
.btn-skip:hover{color:var(--text2)}
/* MATCH POPUP */
.overlay{position:fixed;inset:0;background:rgba(0,0,0,.8);display:flex;align-items:center;justify-content:center;z-index:200;animation:fadeIn .3s}
.match-box{background:var(--surface);border:1px solid var(--border2);border-radius:20px;padding:48px 40px;text-align:center;max-width:380px;animation:popIn .35s cubic-bezier(.175,.885,.32,1.275);position:relative;overflow:hidden}
.match-box::before{content:'';position:absolute;inset:0;background:radial-gradient(circle at 50% 0%,rgba(0,212,170,.1) 0%,transparent 70%);pointer-events:none}
.match-avs{display:flex;justify-content:center;align-items:center;gap:4px;margin-bottom:24px}
.match-av{width:80px;height:80px;border-radius:50%;overflow:hidden;border:3px solid var(--teal)}
.match-av img{width:100%;height:100%;object-fit:cover}
.match-title{font-family:'Rajdhani',sans-serif;font-size:36px;font-weight:700;color:var(--teal);letter-spacing:2px;margin-bottom:8px}
.match-sub{color:var(--text2);font-size:14px;margin-bottom:28px}
.match-contact{background:rgba(0,212,170,.08);border:1px solid rgba(0,212,170,.3);border-radius:10px;padding:14px;margin-bottom:20px}
.match-contact-lbl{font-size:11px;color:var(--teal);font-weight:600;letter-spacing:1px;margin-bottom:4px}
.match-contact-val{font-size:16px;font-weight:600;color:var(--text)}
.match-close{width:100%;background:var(--surface2);border:1px solid var(--border2);color:var(--text2);border-radius:10px;padding:12px;font-size:14px;cursor:pointer;transition:.18s}
.match-close:hover{color:var(--text);border-color:var(--teal)}
/* PAGES */
.page{flex:1;padding:32px;overflow-y:auto;max-width:700px}
.page-title{font-family:'Rajdhani',sans-serif;font-size:28px;font-weight:700;letter-spacing:1px;margin-bottom:24px}
/* MATCHES */
.match-card{background:var(--surface);border:1px solid var(--border2);border-radius:12px;padding:16px 20px;display:flex;align-items:center;gap:16px;transition:.18s;cursor:pointer;margin-bottom:10px}
.match-card:hover{border-color:var(--teal);box-shadow:0 0 0 1px rgba(0,212,170,.2)}
.match-card-av{width:52px;height:52px;border-radius:50%;overflow:hidden;border:2px solid var(--teal);flex-shrink:0}
.match-card-av img{width:100%;height:100%;object-fit:cover}
.mc-name{font-family:'Rajdhani',sans-serif;font-size:18px;font-weight:700}
.mc-meta{font-size:12px;color:var(--text2);margin-top:2px}
.mc-time{font-size:11px;color:var(--text3)}
.discord-btn{background:rgba(0,212,170,.1);border:1px solid rgba(0,212,170,.3);color:var(--teal);border-radius:8px;padding:8px 16px;font-size:12px;font-weight:600;cursor:pointer;transition:.15s}
.discord-btn:hover{background:rgba(0,212,170,.2)}
/* PROFILE */
.prof-header{display:flex;gap:24px;align-items:flex-start;background:var(--surface);border:1px solid var(--border2);border-radius:16px;padding:28px;margin-bottom:20px}
.prof-av{width:80px;height:80px;border-radius:50%;overflow:hidden;border:3px solid var(--accent);flex-shrink:0}
.prof-av img{width:100%;height:100%;object-fit:cover}
.prof-name{font-family:'Rajdhani',sans-serif;font-size:26px;font-weight:700}
.prof-sub{font-size:13px;color:var(--text2);margin-top:3px}
.prof-edit{margin-left:auto;background:var(--surface2);border:1px solid var(--border2);color:var(--text2);border-radius:8px;padding:8px 16px;font-size:13px;cursor:pointer;transition:.15s}
.prof-edit:hover{color:var(--text);border-color:var(--accent)}
.section-card{background:var(--surface);border:1px solid var(--border2);border-radius:14px;padding:22px;margin-bottom:16px}
.section-head{display:flex;align-items:center;justify-content:space-between;margin-bottom:16px}
.section-heading{font-family:'Barlow Condensed',sans-serif;font-size:13px;font-weight:600;letter-spacing:1.5px;color:var(--text3);text-transform:uppercase}
.add-btn{background:rgba(255,70,85,.1);border:1px solid rgba(255,70,85,.3);color:var(--accent);border-radius:6px;padding:5px 12px;font-size:12px;font-weight:600;cursor:pointer;transition:.15s}
.add-btn:hover{background:rgba(255,70,85,.2)}
.stat-row{display:flex;gap:16px}
.stat-box{flex:1;background:var(--surface2);border-radius:10px;padding:16px;text-align:center;border:1px solid var(--border)}
.stat-num{font-family:'Rajdhani',sans-serif;font-size:32px;font-weight:700;color:var(--accent)}
.stat-lbl{font-size:11px;color:var(--text3);margin-top:2px}
/* EMPTY */
.empty{display:flex;flex-direction:column;align-items:center;gap:12px;padding:60px 24px;color:var(--text3);text-align:center}
.empty-icon{font-size:56px;opacity:.3}
.empty-title{font-family:'Rajdhani',sans-serif;font-size:22px;font-weight:700;color:var(--text2)}
/* BOTTOM NAV */
.bottom-nav{display:none}
/* ANIMS */
@keyframes fadeIn{from{opacity:0}to{opacity:1}}
@keyframes popIn{from{opacity:0;transform:scale(.8)}to{opacity:1;transform:scale(1)}}
.anim-left{animation:swLeft .35s ease forwards}
.anim-right{animation:swRight .35s ease forwards}
@keyframes swLeft{to{transform:translateX(-160%) rotate(-20deg);opacity:0}}
@keyframes swRight{to{transform:translateX(160%) rotate(20deg);opacity:0}}
@media(max-width:768px){
  .sidebar{display:none}
  .bottom-nav{display:flex;justify-content:space-around;background:var(--surface);border-top:1px solid var(--border);padding:8px 0 12px;position:sticky;bottom:0;z-index:50}
  .bn-btn{display:flex;flex-direction:column;align-items:center;gap:3px;font-size:10px;color:var(--text3);cursor:pointer;padding:4px 14px;transition:.15s}
  .bn-btn.active{color:var(--accent)}
  .bn-icon{font-size:20px}
  .feed-area{padding:16px 12px}
  .card-stack,.p-card{width:calc(100vw - 24px);max-width:360px}
}
`;

// Helper
const rc = (rank: Rank) => getRankColor(rank);

// PlayerCard
function PlayerCard({ card, isFront, animClass, onDragStart }: {
    card: FeedCard; isFront: boolean; animClass: string;
    onDragStart?: (e: React.MouseEvent | React.TouchEvent) => void;
}) {
    const likeId = `like-${card.id}`;
    const nopeId = `nope-${card.id}`;
    return (
        <div
            className={`p-card ${isFront ? "front" : "behind"} ${animClass}`}
            onMouseDown={isFront ? onDragStart : undefined}
            onTouchStart={isFront ? onDragStart : undefined}
        >
            <div className="sw-like" id={likeId}>LIKE</div>
            <div className="sw-nope" id={nopeId}>NOPE</div>
            <div className="p-card-img">
                <img src={card.avatar_url ?? `https://api.dicebear.com/7.x/avataaars/svg?seed=${card.username}`} alt={card.username} draggable={false} />
            </div>
            <div className="p-card-overlay" />
            <div className="p-card-body">
                <div className="p-card-name">
                    {card.username}, {card.age}
                    <span>{LANG_FLAGS[card.language] ?? "🌐"}</span>
                    <span className="p-card-mic">🎙 МИК</span>
                </div>
                <div className="rank-badge" style={{ color: rc(card.rank as Rank), borderColor: rc(card.rank as Rank) + "55", background: rc(card.rank as Rank) + "18" }}>
                    ● {card.rank}
                </div>
                <p className="p-card-bio">{card.description || card.bio}</p>
                <div className="tags">
                    <span className="tag">{ROLE_ICONS[card.role as Role]} {card.role}</span>
                    <span className="tag">📍 {card.region}</span>
                    <span className="tag">{card.language}</span>
                </div>
            </div>
        </div>
    );
}

//MatchPopup
function MatchPopup({ matchedCard, onClose }: { matchedCard: FeedCard; onClose: () => void }) {
    const { user } = useAuthStore();
    return (
        <div className="overlay" onClick={onClose}>
            <div className="match-box" onClick={e => e.stopPropagation()}>
                <div className="match-avs">
                    <div className="match-av">
                        <img src={`https://api.dicebear.com/7.x/avataaars/svg?seed=${user?.username}`} alt="me" />
                    </div>
                    <span style={{ fontSize: 28, margin: "0 4px" }}>💚</span>
                    <div className="match-av">
                        <img src={matchedCard.avatar_url ?? `https://api.dicebear.com/7.x/avataaars/svg?seed=${matchedCard.username}`} alt={matchedCard.username} />
                    </div>
                </div>
                <div className="match-title">IT'S A MATCH!</div>
                <p className="match-sub">Вы и <b>{matchedCard.username}</b> понравились друг другу</p>
                {matchedCard.discord && (
                    <div className="match-contact">
                        <div className="match-contact-lbl">DISCORD</div>
                        <div className="match-contact-val">{matchedCard.discord}</div>
                    </div>
                )}
                {matchedCard.telegram && (
                    <div className="match-contact">
                        <div className="match-contact-lbl">TELEGRAM</div>
                        <div className="match-contact-val">{matchedCard.telegram}</div>
                    </div>
                )}
                <button className="match-close" onClick={onClose}>Продолжить поиск</button>
            </div>
        </div>
    );
}

// FeedPage
function FeedPage() {
    const { cards, isLoading, newMatch, swipe, clearMatch, loadFeed, setFilters } = useFeedStore();
    const [animDir, setAnimDir] = useState<"left" | "right" | null>(null);
    const [matchedCard, setMatchedCard] = useState<FeedCard | null>(null);
    const [localFilters, setLocalFilters] = useState({ rankMin: "Iron", rankMax: "Radiant", roles: [] as string[], langs: [] as string[] });
    const cardRef = useRef<HTMLDivElement>(null);

    useEffect(() => { if (cards.length === 0) loadFeed(); }, []);

    const doSwipe = (card: FeedCard, dir: "left" | "right") => {
        setAnimDir(dir);
        if (dir === "right") setMatchedCard(card);
        setTimeout(() => {
            setAnimDir(null);
            swipe(card, dir === "right" ? "like" : "dislike");
        }, 350);
    };

    const handleDragStart = (card: FeedCard) => (e: React.MouseEvent | React.TouchEvent) => {
        e.preventDefault();
        const startX = "touches" in e ? e.touches[0].clientX : e.clientX;
        const el = cardRef.current;
        let delta = 0;

        const move = (ev: MouseEvent | TouchEvent) => {
            const x = ("touches" in ev ? ev.touches[0].clientX : ev.clientX) - startX;
            delta = x;
            if (el) el.style.transform = `translateX(${x}px) rotate(${x * 0.06}deg)`;
            const likeEl = document.getElementById(`like-${card.id}`);
            const nopeEl = document.getElementById(`nope-${card.id}`);
            if (likeEl) likeEl.style.opacity = String(Math.max(0, x / 80));
            if (nopeEl) nopeEl.style.opacity = String(Math.max(0, -x / 80));
        };
        const up = () => {
            window.removeEventListener("mousemove", move);
            window.removeEventListener("touchmove", move);
            window.removeEventListener("mouseup", up);
            window.removeEventListener("touchend", up);
            if (el) el.style.transform = "";
            const likeEl = document.getElementById(`like-${card.id}`);
            const nopeEl = document.getElementById(`nope-${card.id}`);
            if (likeEl) likeEl.style.opacity = "0";
            if (nopeEl) nopeEl.style.opacity = "0";
            if (Math.abs(delta) > 90) doSwipe(card, delta > 0 ? "right" : "left");
        };
        window.addEventListener("mousemove", move);
        window.addEventListener("touchmove", move);
        window.addEventListener("mouseup", up);
        window.addEventListener("touchend", up);
    };

    const toggleArr = (key: "roles" | "langs", val: string) =>
        setLocalFilters(f => ({ ...f, [key]: f[key].includes(val) ? f[key].filter(x => x !== val) : [...f[key], val] }));

    const front = cards[0];
    const behind = cards[1];

    return (
        <div className="main">
            <aside className="sidebar">
                <div className="sidebar-section">Фильтры поиска</div>
                <div className="filter-panel">
                    <div>
                        <div className="filter-label">Ранг (мин–макс)</div>
                        <div className="range-row">
                            <select className="filter-select" style={{ padding: "6px 8px" }}
                                    value={localFilters.rankMin} onChange={e => setLocalFilters(f => ({ ...f, rankMin: e.target.value }))}>
                                {RANKS.map(r => <option key={r}>{r}</option>)}
                            </select>
                            <span style={{ color: "var(--text3)" }}>—</span>
                            <select className="filter-select" style={{ padding: "6px 8px" }}
                                    value={localFilters.rankMax} onChange={e => setLocalFilters(f => ({ ...f, rankMax: e.target.value }))}>
                                {RANKS.map(r => <option key={r}>{r}</option>)}
                            </select>
                        </div>
                    </div>
                    <div>
                        <div className="filter-label">Роль</div>
                        <div className="chips">
                            {ROLES.map(r => (
                                <span key={r} className={`chip ${localFilters.roles.includes(r) ? "on" : ""}`}
                                      onClick={() => toggleArr("roles", r)}>{ROLE_ICONS[r]} {r}</span>
                            ))}
                        </div>
                    </div>
                    <div>
                        <div className="filter-label">Язык</div>
                        <div className="chips">
                            {LANGS.map(l => (
                                <span key={l} className={`chip teal ${localFilters.langs.includes(l) ? "on" : ""}`}
                                      onClick={() => toggleArr("langs", l)}>{LANG_FLAGS[l]} {l}</span>
                            ))}
                        </div>
                    </div>
                    <button className="apply-btn" onClick={() => {
                        setFilters({ rank_min: localFilters.rankMin as Rank, rank_max: localFilters.rankMax as Rank });
                        loadFeed();
                    }}>Применить</button>
                </div>
            </aside>

            <div className="feed-area">
                {isLoading ? (
                    <div className="empty"><div className="empty-icon">⏳</div><div className="empty-title">Загружаем...</div></div>
                ) : cards.length === 0 ? (
                    <div className="empty">
                        <div className="empty-icon">🎮</div>
                        <div className="empty-title">Карточки закончились</div>
                        <p>Расширь фильтры или зайди позже</p>
                        <button className="apply-btn" style={{ marginTop: 12, padding: "10px 24px" }} onClick={loadFeed}>Обновить</button>
                    </div>
                ) : (
                    <>
                        <div className="card-stack">
                            {behind && <PlayerCard card={behind} isFront={false} animClass="" />}
                            {front && (
                                <div ref={cardRef} style={{ position: "absolute", top: 0, left: 0, zIndex: 2, width: "100%", height: "100%" }}>
                                    <PlayerCard
                                        card={front} isFront={true}
                                        animClass={animDir === "left" ? "anim-left" : animDir === "right" ? "anim-right" : ""}
                                        onDragStart={handleDragStart(front)}
                                    />
                                </div>
                            )}
                        </div>
                        <div className="action-row">
                            <button className="btn-dis" onClick={() => front && doSwipe(front, "left")}>✕</button>
                            <button className="btn-skip">Пропустить</button>
                            <button className="btn-lik" onClick={() => front && doSwipe(front, "right")}>♥</button>
                        </div>
                    </>
                )}
            </div>

            {newMatch && matchedCard && (
                <MatchPopup matchedCard={matchedCard} onClose={() => { clearMatch(); setMatchedCard(null); }} />
            )}
        </div>
    );
}

// MatchesPage
function MatchesPage() {
    const [matches, setMatches] = useState<any[]>([]);
    useEffect(() => {
        getMatches().catch(() =>
            setMatches([
                { id:1, username:"NightOwl",  rank:"Platinum", role:"Controller", discord:"NightOwl#4269", created_at:"2h ago",   avatar_url:"https://api.dicebear.com/7.x/avataaars/svg?seed=Night&backgroundColor=d1d4f9" },
                { id:2, username:"VoidRift",  rank:"Diamond",  role:"Initiator",  discord:"VoidRift#1337", created_at:"Yesterday", avatar_url:"https://api.dicebear.com/7.x/avataaars/svg?seed=Void&backgroundColor=ffd5dc" },
                { id:3, username:"ZephyrX",   rank:"Ascendant",role:"Controller", discord:"ZephyrX#0001",  created_at:"2d ago",    avatar_url:"https://api.dicebear.com/7.x/avataaars/svg?seed=Zeph&backgroundColor=b6e3f4" },
            ])
        );
    }, []);

    return (
        <div className="page">
            <div className="page-title">Матчи <span style={{ color: "var(--teal)", fontSize: 20 }}>{matches.length}</span></div>
            {matches.map(m => (
                <div key={m.id} className="match-card">
                    <div className="match-card-av"><img src={m.avatar_url} alt={m.username} /></div>
                    <div style={{ flex: 1 }}>
                        <div className="mc-name">{m.username}</div>
                        <div className="mc-meta">
                            <span style={{ color: rc(m.rank) }}>● {m.rank}</span> · {ROLE_ICONS[m.role as Role]} {m.role}
                        </div>
                    </div>
                    <div style={{ textAlign: "right", display: "flex", flexDirection: "column", alignItems: "flex-end", gap: 8 }}>
                        <div className="mc-time">{m.created_at}</div>
                        {m.discord && <button className="discord-btn">Discord →</button>}
                    </div>
                </div>
            ))}
        </div>
    );
}

//ProfilePage
function ProfilePage() {
    const { user, logout } = useAuthStore();
    const [showSetup, setShowSetup] = useState(false);
    const [bio, setBio] = useState("");
    const [saving, setSaving] = useState(false);

    const saveProfile = async () => {
        setSaving(true);
        try { await createProfile({ bio }); setShowSetup(false); }
        catch { /* already exists */ setShowSetup(false); }
        finally { setSaving(false); }
    };

    if (!user) return null;
    return (
        <div className="page">
            <div className="prof-header">
                <div className="prof-av">
                    <img src={`https://api.dicebear.com/7.x/avataaars/svg?seed=${user.username}&backgroundColor=ff5a5f`} alt="me" />
                </div>
                <div>
                    <div className="prof-name">{user.username}</div>
                    <div className="prof-sub">{LANG_FLAGS[user.language] ?? "🌐"} {user.language} · {user.region} · {user.age} лет</div>
                    {user.discord && <div className="prof-sub" style={{ marginTop: 4 }}>🎮 {user.discord}</div>}
                    {user.telegram && <div className="prof-sub">✈️ {user.telegram}</div>}
                </div>
                <button className="prof-edit" onClick={logout}>Выйти</button>
            </div>

            {showSetup ? (
                <div className="section-card">
                    <div className="section-head"><div className="section-heading">Настроить профиль</div></div>
                    <div style={{ display: "flex", flexDirection: "column", gap: 10 }}>
            <textarea
                style={{ background: "var(--surface2)", border: "1px solid var(--border2)", color: "var(--text)", borderRadius: 8, padding: "10px 12px", fontSize: 13, resize: "vertical", minHeight: 80, fontFamily: "Barlow, sans-serif" }}
                placeholder="Расскажи о себе..."
                value={bio} onChange={e => setBio(e.target.value)}
            />
                        <button className="apply-btn" onClick={saveProfile} disabled={saving}>
                            {saving ? "Сохраняем..." : "Сохранить"}
                        </button>
                    </div>
                </div>
            ) : (
                <div className="section-card">
                    <div className="section-head">
                        <div className="section-heading">Профиль</div>
                        <button className="add-btn" onClick={() => setShowSetup(true)}>+ Настроить</button>
                    </div>
                    <p style={{ fontSize: 13, color: "var(--text2)" }}>Добавь описание и аватар — тебя лучше заметят</p>
                </div>
            )}

            <div className="section-card">
                <div className="section-head">
                    <div className="section-heading">Мои карточки</div>
                    <button className="add-btn">+ Добавить</button>
                </div>
                <p style={{ fontSize: 13, color: "var(--text3)" }}>Карточки появятся когда бэкенд реализует /cards</p>
            </div>

            <div className="section-card">
                <div className="section-head"><div className="section-heading">Статистика</div></div>
                <div className="stat-row">
                    <div className="stat-box"><div className="stat-num">—</div><div className="stat-lbl">Свайпов</div></div>
                    <div className="stat-box"><div className="stat-num" style={{ color: "var(--teal)" }}>—</div><div className="stat-lbl">Матчей</div></div>
                    <div className="stat-box"><div className="stat-num" style={{ color: "var(--gold)" }}>—</div><div className="stat-lbl">Игр</div></div>
                </div>
            </div>
        </div>
    );
}

// App
export default function App() {
    const { isAuthed, fetchMe } = useAuthStore();
    const { setGame, gameId } = useFeedStore();
    const [authPage, setAuthPage] = useState<"login" | "register">("login");
    const [page, setPage] = useState<"feed" | "matches" | "profile">("feed");
    const [showGameDrop, setShowGameDrop] = useState(false);

    useEffect(() => {
        const s = document.createElement("style");
        s.textContent = APP_CSS;
        document.head.appendChild(s);
        return () => s.remove();
    }, []);

    useEffect(() => { if (isAuthed) fetchMe(); }, []);

    if (!isAuthed) {
        return authPage === "login"
            ? <LoginPage onSwitch={() => setAuthPage("register")} />
            : <RegisterPage onSwitch={() => setAuthPage("login")} />;
    }

    const NAV = [
        { id: "feed",    icon: "🃏", label: "Фид" },
        { id: "matches", icon: "💚", label: "Матчи" },
        { id: "profile", icon: "👤", label: "Профиль" },
    ] as const;

    return (
        <div className="app">
            <nav className="nav">
                <div className="nav-logo" onClick={() => setPage("feed")}>Game<span>Match</span></div>

                <div style={{ position: "relative" }}>
                    <button className="nav-game-select" onClick={() => setShowGameDrop(v => !v)}>
                        {GAMES.find(g => g.id === gameId)?.icon} {GAMES.find(g => g.id === gameId)?.label} ▾
                    </button>
                    {showGameDrop && (
                        <div className="dropdown">
                            {GAMES.map(g => (
                                <div key={g.id} className={`dropdown-opt ${gameId === g.id ? "active" : ""}`}
                                     onClick={() => { setGame(g.id); setShowGameDrop(false); }}>
                                    {g.icon} {g.label}
                                </div>
                            ))}
                        </div>
                    )}
                </div>

                <div className="nav-spacer" />
                {NAV.map(n => (
                    <button key={n.id} className={`nav-btn ${page === n.id ? "active" : ""}`} onClick={() => setPage(n.id)}>
                        {n.icon} {n.label}
                    </button>
                ))}
            </nav>

            {page === "feed"    && <FeedPage />}
            {page === "matches" && <div className="main"><MatchesPage /></div>}
            {page === "profile" && <div className="main"><ProfilePage /></div>}

            <nav className="bottom-nav">
                {NAV.map(n => (
                    <div key={n.id} className={`bn-btn ${page === n.id ? "active" : ""}`} onClick={() => setPage(n.id)}>
                        <span className="bn-icon">{n.icon}</span>{n.label}
                    </div>
                ))}
            </nav>
        </div>
    );
}