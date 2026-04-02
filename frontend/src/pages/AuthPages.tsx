import { useState } from "react";
import { useAuthStore } from "../store/authStore";
import { login, register } from "../api/auth";

const REGIONS = ["EU Central","EU West","NA East","NA West","CIS","Asia"];
const LANGS = ["RU","EN","DE","KZ","TR","PL"];

// стили
const authCSS = `
.auth-wrap{min-height:100vh;display:flex;align-items:center;justify-content:center;background:var(--bg);padding:20px}
.auth-box{background:var(--surface);border:1px solid var(--border2);border-radius:16px;padding:40px 36px;width:100%;max-width:420px}
.auth-logo{font-family:'Rajdhani',sans-serif;font-size:28px;font-weight:700;text-align:center;margin-bottom:28px}
.auth-logo span{color:var(--accent)}
.auth-title{font-family:'Rajdhani',sans-serif;font-size:22px;font-weight:700;margin-bottom:20px}
.field{display:flex;flex-direction:column;gap:6px;margin-bottom:14px}
.field label{font-size:12px;font-weight:600;color:var(--text3);letter-spacing:1px;text-transform:uppercase;font-family:'Barlow Condensed',sans-serif}
.field input,.field select{background:var(--surface2);border:1px solid var(--border2);color:var(--text);border-radius:8px;padding:10px 12px;font-size:14px;font-family:'Barlow',sans-serif;outline:none;transition:.15s}
.field input:focus,.field select:focus{border-color:var(--accent)}
.field-row{display:grid;grid-template-columns:1fr 1fr;gap:12px}
.auth-btn{width:100%;background:var(--accent);color:#fff;border:none;border-radius:10px;padding:12px;font-size:15px;font-weight:700;cursor:pointer;font-family:'Rajdhani',sans-serif;letter-spacing:.5px;transition:.18s;margin-top:4px}
.auth-btn:hover{background:var(--accent2);box-shadow:var(--glow)}
.auth-btn:disabled{opacity:.5;cursor:not-allowed}
.auth-switch{text-align:center;margin-top:18px;font-size:13px;color:var(--text3)}
.auth-switch button{background:none;border:none;color:var(--accent);cursor:pointer;font-size:13px;font-weight:600;padding:0}
.auth-error{background:rgba(255,70,85,.12);border:1px solid rgba(255,70,85,.3);color:var(--accent);border-radius:8px;padding:10px 14px;font-size:13px;margin-bottom:14px}
`;

function injectCSS(css: string, id: string) {
    if (!document.getElementById(id)) {
        const s = document.createElement("style");
        s.id = id;
        s.textContent = css;
        document.head.appendChild(s);
    }
}

export function LoginPage({ onSwitch }: { onSwitch: () => void }) {
    injectCSS(authCSS, "auth-css");
    const { setAuth } = useAuthStore();
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [error, setError] = useState("");
    const [loading, setLoading] = useState(false);

    const handleSubmit = async () => {
        if (!email || !password) { setError("Заполни все поля"); return; }
        setLoading(true); setError("");
        try {
            const data = await login({ email, password });
            setAuth(data.token, data.user);
        } catch (e: any) {
            setError(e.response?.data?.error ?? "Неверный email или пароль");
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="auth-wrap">
            <div className="auth-box">
                <div className="auth-logo">Game<span>Match</span></div>
                <div className="auth-title">Войти</div>
                {error && <div className="auth-error">{error}</div>}
                <div className="field">
                    <label>Email</label>
                    <input type="email" placeholder="your@email.com" value={email} onChange={e=>setEmail(e.target.value)} />
                </div>
                <div className="field">
                    <label>Пароль</label>
                    <input type="password" placeholder="••••••••" value={password} onChange={e=>setPassword(e.target.value)}
                           onKeyDown={e=>e.key==="Enter"&&handleSubmit()} />
                </div>
                <button className="auth-btn" onClick={handleSubmit} disabled={loading}>
                    {loading ? "Входим..." : "Войти"}
                </button>
                <div className="auth-switch">
                    Нет аккаунта?{" "}
                    <button onClick={onSwitch}>Зарегистрироваться</button>
                </div>
            </div>
        </div>
    );
}

export function RegisterPage({ onSwitch }: { onSwitch: () => void }) {
    injectCSS(authCSS, "auth-css");
    const { setAuth } = useAuthStore();
    const [form, setForm] = useState({
        username:"", email:"", password:"", age:"", language:"RU", region:"EU Central", discord:"", telegram:""
    });
    const [error, setError] = useState("");
    const [loading, setLoading] = useState(false);

    const set = (k: string) => (e: any) => setForm(f=>({...f,[k]:e.target.value}));

    const handleSubmit = async () => {
        if (!form.username || !form.email || !form.password || !form.age) {
            setError("Заполни обязательные поля"); return;
        }
        if (Number(form.age) < 16) { setError("Минимальный возраст — 16 лет"); return; }
        setLoading(true); setError("");
        try {
            const data = await register({
                username: form.username,
                email: form.email,
                password: form.password,
                age: Number(form.age),
                language: form.language,
                region: form.region,
                discord: form.discord || undefined,
                telegram: form.telegram || undefined,
            });
            setAuth(data.token, data.user);
        } catch (e: any) {
            setError(e.response?.data?.error ?? "Ошибка регистрации");
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="auth-wrap">
            <div className="auth-box">
                <div className="auth-logo">Game<span>Match</span></div>
                <div className="auth-title">Создать аккаунт</div>
                {error && <div className="auth-error">{error}</div>}
                <div className="field-row">
                    <div className="field"><label>Никнейм *</label><input placeholder="KazFrag" value={form.username} onChange={set("username")} /></div>
                    <div className="field"><label>Возраст *</label><input type="number" min={16} max={99} placeholder="18" value={form.age} onChange={set("age")} /></div>
                </div>
                <div className="field"><label>Email *</label><input type="email" placeholder="your@email.com" value={form.email} onChange={set("email")} /></div>
                <div className="field"><label>Пароль *</label><input type="password" placeholder="мин. 8 символов" value={form.password} onChange={set("password")} /></div>
                <div className="field-row">
                    <div className="field"><label>Язык</label><select value={form.language} onChange={set("language")}>{LANGS.map(l=><option key={l}>{l}</option>)}</select></div>
                    <div className="field"><label>Регион</label><select value={form.region} onChange={set("region")}>{REGIONS.map(r=><option key={r}>{r}</option>)}</select></div>
                </div>
                <div className="field-row">
                    <div className="field"><label>Discord</label><input placeholder="user#1234" value={form.discord} onChange={set("discord")} /></div>
                    <div className="field"><label>Telegram</label><input placeholder="@username" value={form.telegram} onChange={set("telegram")} /></div>
                </div>
                <button className="auth-btn" onClick={handleSubmit} disabled={loading}>{loading ? "Создаём..." : "Зарегистрироваться"}</button>
                <div className="auth-switch">Уже есть аккаунт? <button onClick={onSwitch}>Войти</button></div>
            </div>
        </div>
    );
}