"use client";

import React, { useState, useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { DashboardLayout } from "@/components/DashboardLayout";
import { useAuth } from "@/context/AuthContext";
import { apiClient } from "@/lib/api-client";

interface GameAnalysis {
  verdict: "sim" | "nao" | "depende";
  verdict_reason: string;
  performance_level: "baixo" | "medio" | "alto" | "ultra";
  min: { cpu: string; ram: string; gpu: string; storage: string };
  rec: { cpu: string; ram: string; gpu: string; storage: string };
  user: { cpu: string; ram: string; gpu: string; storage: string };
  meets: {
    min_cpu: boolean; min_ram: boolean; min_gpu: boolean; min_storage: boolean;
    rec_cpu: boolean; rec_ram: boolean; rec_gpu: boolean; rec_storage: boolean;
  };
  estimated_fps: string;
  tips: string[];
  genre: string;
  release_year: number;
}

const VERDICT_CONFIG = {
  sim: { label: "RODA!", color: "#10b981", glow: "rgba(16,185,129,0.4)", icon: "✓", bg: "rgba(16,185,129,0.08)", border: "rgba(16,185,129,0.3)" },
  nao: { label: "NÃO RODA", color: "#ef4444", glow: "rgba(239,68,68,0.4)", icon: "✕", bg: "rgba(239,68,68,0.08)", border: "rgba(239,68,68,0.3)" },
  depende: { label: "DEPENDE", color: "#f59e0b", glow: "rgba(245,158,11,0.4)", icon: "!", bg: "rgba(245,158,11,0.08)", border: "rgba(245,158,11,0.3)" },
};

const PERF_CONFIG = {
  baixo: { label: "Configurações Baixas", color: "#ef4444" },
  medio: { label: "Configurações Médias", color: "#f59e0b" },
  alto: { label: "Configurações Altas", color: "#10b981" },
  ultra: { label: "Ultra / 4K", color: "#00d4ff" },
};

const RAWG_KEY = process.env.NEXT_PUBLIC_RAWG_API_KEY ?? "";

interface TrendingGame {
  id: number;
  name: string;
  background_image: string | null;
  rating: number;
  genres: { name: string }[];
}

async function fetchGameCover(gameName: string): Promise<string | null> {
  try {
    const res = await fetch(
      `https://api.rawg.io/api/games?search=${encodeURIComponent(gameName)}&key=${RAWG_KEY}&page_size=1`,
      { signal: AbortSignal.timeout(5000) }
    );
    if (!res.ok) return null;
    const data = await res.json();
    return data?.results?.[0]?.background_image ?? null;
  } catch {
    return null;
  }
}

async function fetchTrendingGames(): Promise<TrendingGame[]> {
  try {
    const today = new Date().toISOString().split("T")[0];
    const weekAgo = new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString().split("T")[0];
    const res = await fetch(
      `https://api.rawg.io/api/games?key=${RAWG_KEY}&dates=${weekAgo},${today}&ordering=-added&page_size=5`,
      { signal: AbortSignal.timeout(6000) }
    );
    if (!res.ok) return [];
    const data = await res.json();
    return data?.results ?? [];
  } catch {
    return [];
  }
}

function StatusDot({ meets }: { meets: boolean | undefined }) {
  if (meets === undefined) return <span style={{ color: "#475569" }}>—</span>;
  return (
    <span style={{
      display: "inline-block", width: 8, height: 8, borderRadius: "50%",
      background: meets ? "#10b981" : "#ef4444",
      boxShadow: meets ? "0 0 6px #10b981" : "0 0 6px #ef4444",
    }} />
  );
}

function GameCover({ name, genre, imageUrl }: { name: string; genre: string; imageUrl: string | null }) {
  const colors = [
    ["#00d4ff", "#0284c7"], ["#10b981", "#047857"], ["#a78bfa", "#7c3aed"],
    ["#f59e0b", "#d97706"], ["#ef4444", "#dc2626"], ["#ec4899", "#db2777"],
  ];
  const idx = name.charCodeAt(0) % colors.length;
  const [c1, c2] = colors[idx];
  const [imgError, setImgError] = useState(false);

  if (imageUrl && !imgError) {
    return (
      <div style={{
        width: "100%", aspectRatio: "16/9", borderRadius: 10,
        overflow: "hidden", border: "1px solid rgba(255,255,255,0.08)",
        position: "relative",
      }}>
        <img
          src={imageUrl}
          alt={name}
          onError={() => setImgError(true)}
          style={{ width: "100%", height: "100%", objectFit: "cover", display: "block" }}
        />
        <div style={{
          position: "absolute", inset: 0,
          background: "linear-gradient(to top, rgba(5,13,26,0.6) 0%, transparent 50%)",
        }} />
      </div>
    );
  }

  return (
    <div style={{
      width: "100%", aspectRatio: "16/9", borderRadius: 10,
      background: `linear-gradient(135deg, ${c1}22, ${c2}44)`,
      border: `1px solid ${c1}33`,
      display: "flex", flexDirection: "column", alignItems: "center", justifyContent: "center",
      position: "relative", overflow: "hidden",
    }}>
      <div style={{
        position: "absolute", inset: 0,
        backgroundImage: `radial-gradient(circle at 30% 50%, ${c1}15 0%, transparent 60%), radial-gradient(circle at 70% 50%, ${c2}15 0%, transparent 60%)`,
      }} />
      <span style={{ fontSize: 48, fontWeight: 900, color: c1, textShadow: `0 0 30px ${c1}`, zIndex: 1 }}>
        {name.charAt(0).toUpperCase()}
      </span>
      <span style={{ fontSize: 10, color: c1 + "99", letterSpacing: "0.2em", zIndex: 1, marginTop: 4 }}>
        {genre?.toUpperCase() || "GAME"}
      </span>
    </div>
  );
}

function ComparisonRow({ label, user, min, rec, meetsMin, meetsRec }: {
  label: string; user: string; min: string; rec: string; meetsMin: boolean; meetsRec: boolean;
}) {
  return (
    <div style={{
      display: "grid", gridTemplateColumns: "80px 1fr 1fr 1fr",
      gap: 8, padding: "10px 0",
      borderBottom: "1px solid rgba(255,255,255,0.04)",
      alignItems: "center",
    }}>
      <span style={{ color: "#475569", fontSize: 11, textTransform: "uppercase", letterSpacing: "0.08em" }}>{label}</span>
      <div style={{ display: "flex", alignItems: "center", gap: 6 }}>
        <StatusDot meets={meetsMin} />
        <span style={{ color: "#e2e8f0", fontSize: 12 }}>{user}</span>
      </div>
      <span style={{ color: meetsMin ? "#10b981" : "#ef4444", fontSize: 12 }}>{min}</span>
      <span style={{ color: meetsRec ? "#10b981" : "#f59e0b", fontSize: 12 }}>{rec}</span>
    </div>
  );
}

export default function GamesPage() {
  const { token } = useAuth();
  const [query, setQuery] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [analysis, setAnalysis] = useState<GameAnalysis | null>(null);
  const [gameName, setGameName] = useState("");
  const [coverUrl, setCoverUrl] = useState<string | null>(null);
  const [error, setError] = useState("");
  const [trending, setTrending] = useState<TrendingGame[]>([]);
  const [trendingLoading, setTrendingLoading] = useState(true);

  useEffect(() => {
    fetchTrendingGames().then((games) => {
      setTrending(games);
      setTrendingLoading(false);
    });
  }, []);

  const runAnalysis = async (name: string) => {
    if (!name.trim() || !token) return;
    setIsLoading(true);
    setAnalysis(null);
    setCoverUrl(null);
    setError("");
    setGameName(name.trim());
    setQuery(name.trim());

    const coverPromise = fetchGameCover(name.trim());
    const prompt = `Analise se o jogo "${name}" roda no PC do usuário. Responda APENAS com JSON válido (sem markdown, sem texto antes ou depois):
{
  "verdict": "sim" ou "nao" ou "depende",
  "verdict_reason": "frase curta de até 10 palavras",
  "performance_level": "baixo" ou "medio" ou "alto" ou "ultra",
  "min": {"cpu": "...", "ram": "X GB", "gpu": "...", "storage": "X GB"},
  "rec": {"cpu": "...", "ram": "X GB", "gpu": "...", "storage": "X GB"},
  "user": {"cpu": "i7-1255U", "ram": "11.7 GB", "gpu": "Sem GPU dedicada", "storage": "474 GB disponível"},
  "meets": {
    "min_cpu": true/false, "min_ram": true/false, "min_gpu": true/false, "min_storage": true/false,
    "rec_cpu": true/false, "rec_ram": true/false, "rec_gpu": true/false, "rec_storage": true/false
  },
  "estimated_fps": "Ex: 30-60 FPS em Low",
  "tips": ["dica 1 em português", "dica 2 em português", "dica 3 em português"],
  "genre": "gênero do jogo",
  "release_year": ano
}`;

    try {
      const [response, cover] = await Promise.all([
        apiClient.chatCompletion(prompt, token),
        coverPromise,
      ]);
      setCoverUrl(cover);
      if (response.success && response.data) {
        const text = response.data.response.trim();
        const jsonStart = text.indexOf("{");
        const jsonEnd = text.lastIndexOf("}") + 1;
        const jsonStr = jsonStart >= 0 ? text.slice(jsonStart, jsonEnd) : text;
        const parsed = JSON.parse(jsonStr) as GameAnalysis;
        setAnalysis(parsed);
      } else {
        setError("Não consegui analisar o jogo. Tente novamente.");
      }
    } catch {
      setError("Erro ao processar a resposta. Tente pesquisar novamente.");
    } finally {
      setIsLoading(false);
    }
  };

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!query.trim() || !token) return;
    await runAnalysis(query.trim());
  };

  const verdict = analysis ? VERDICT_CONFIG[analysis.verdict] : null;
  const perf = analysis ? PERF_CONFIG[analysis.performance_level] : null;

  return (
    <DashboardLayout title="Games">
      {/* Search */}
      <div style={{
        background: "rgba(10,22,40,0.85)", backdropFilter: "blur(12px)",
        border: "1px solid rgba(0,212,255,0.2)", borderRadius: 12, padding: 20, marginBottom: 24,
      }}>
        <form onSubmit={handleSearch} style={{ display: "flex", gap: 12 }}>
          <input
            type="text" placeholder="Ex: GTA V, Cyberpunk 2077, Elden Ring, Minecraft..."
            value={query} onChange={e => setQuery(e.target.value)} disabled={isLoading}
            style={{
              flex: 1, background: "rgba(255,255,255,0.04)",
              border: "1px solid rgba(0,212,255,0.2)", borderRadius: 8,
              padding: "10px 14px", color: "#e2e8f0", fontSize: 14, outline: "none",
            }}
            onFocus={e => (e.target.style.borderColor = "rgba(0,212,255,0.6)")}
            onBlur={e => (e.target.style.borderColor = "rgba(0,212,255,0.2)")}
          />
          <button type="submit" disabled={isLoading || !query.trim()} style={{
            padding: "10px 28px",
            background: isLoading || !query.trim() ? "rgba(0,212,255,0.04)" : "rgba(0,212,255,0.12)",
            border: "1px solid rgba(0,212,255,0.4)", borderRadius: 8,
            color: "#00d4ff", fontSize: 14, fontWeight: 700, cursor: isLoading || !query.trim() ? "not-allowed" : "pointer",
            letterSpacing: "0.08em", minWidth: 110,
          }}>
            {isLoading ? "..." : "VERIFICAR"}
          </button>
        </form>
      </div>

      {/* Loading */}
      {isLoading && (
        <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }}
          style={{ textAlign: "center", padding: "60px 0" }}>
          <div style={{ display: "flex", justifyContent: "center", gap: 8, marginBottom: 20 }}>
            {[0, 0.2, 0.4].map((d, i) => (
              <motion.div key={i} animate={{ scale: [1, 1.4, 1] }}
                transition={{ duration: 0.8, repeat: Infinity, delay: d }}
                style={{ width: 10, height: 10, borderRadius: "50%", background: "#00d4ff",
                  boxShadow: "0 0 10px #00d4ff" }} />
            ))}
          </div>
          <p style={{ color: "#475569", fontSize: 13 }}>
            JARVIS analisando <strong style={{ color: "#94a3b8" }}>{gameName}</strong>...
          </p>
          <p style={{ color: "#334155", fontSize: 11, marginTop: 6 }}>Comparando com suas specs</p>
        </motion.div>
      )}

      {/* Error */}
      {error && (
        <div style={{ background: "rgba(239,68,68,0.08)", border: "1px solid rgba(239,68,68,0.3)",
          borderRadius: 12, padding: 16, color: "#f87171", fontSize: 14 }}>⚠ {error}
        </div>
      )}

      {/* Result */}
      <AnimatePresence>
        {analysis && verdict && (
          <motion.div initial={{ opacity: 0, y: 24 }} animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5 }}>

            <div style={{ display: "grid", gridTemplateColumns: "280px 1fr", gap: 20, marginBottom: 20 }}>
              {/* Left: cover + verdict */}
              <div style={{ display: "flex", flexDirection: "column", gap: 12 }}>
                <GameCover name={gameName} genre={analysis.genre} imageUrl={coverUrl} />

                <motion.div initial={{ scale: 0.8 }} animate={{ scale: 1 }}
                  transition={{ type: "spring", stiffness: 300 }}
                  style={{
                    background: verdict.bg, border: `1px solid ${verdict.border}`,
                    borderRadius: 12, padding: "20px 16px", textAlign: "center",
                    boxShadow: `0 0 20px ${verdict.glow}`,
                  }}>
                  <div style={{ fontSize: 36, fontWeight: 900, color: verdict.color,
                    textShadow: `0 0 20px ${verdict.glow}`, letterSpacing: "0.05em" }}>
                    {verdict.icon} {verdict.label}
                  </div>
                  <p style={{ color: verdict.color + "cc", fontSize: 13, marginTop: 6 }}>
                    {analysis.verdict_reason}
                  </p>
                </motion.div>

                {analysis.verdict !== "nao" && perf && (
                  <div style={{ background: "rgba(10,22,40,0.6)", border: "1px solid rgba(255,255,255,0.06)",
                    borderRadius: 10, padding: "12px 16px", textAlign: "center" }}>
                    <p style={{ color: "#475569", fontSize: 10, textTransform: "uppercase", letterSpacing: "0.1em", marginBottom: 4 }}>Performance Esperada</p>
                    <p style={{ color: perf.color, fontSize: 13, fontWeight: 600 }}>{perf.label}</p>
                    <p style={{ color: "#475569", fontSize: 12, marginTop: 4 }}>{analysis.estimated_fps}</p>
                  </div>
                )}
              </div>

              {/* Right: comparison table */}
              <div style={{ background: "rgba(10,22,40,0.85)", backdropFilter: "blur(12px)",
                border: "1px solid rgba(0,212,255,0.2)", borderRadius: 12, padding: 20 }}>
                <h3 style={{ color: "#00d4ff", fontSize: 12, textTransform: "uppercase",
                  letterSpacing: "0.12em", marginBottom: 16 }}>Comparação de Specs</h3>

                <div style={{ display: "grid", gridTemplateColumns: "80px 1fr 1fr 1fr",
                  gap: 8, marginBottom: 8 }}>
                  {["", "Seu PC", "Mínimo", "Recomendado"].map((h, i) => (
                    <span key={i} style={{ color: "#334155", fontSize: 10,
                      textTransform: "uppercase", letterSpacing: "0.08em" }}>{h}</span>
                  ))}
                </div>

                <ComparisonRow label="CPU" user={analysis.user.cpu}
                  min={analysis.min.cpu} rec={analysis.rec.cpu}
                  meetsMin={analysis.meets.min_cpu} meetsRec={analysis.meets.rec_cpu} />
                <ComparisonRow label="RAM" user={analysis.user.ram}
                  min={analysis.min.ram} rec={analysis.rec.ram}
                  meetsMin={analysis.meets.min_ram} meetsRec={analysis.meets.rec_ram} />
                <ComparisonRow label="GPU" user={analysis.user.gpu}
                  min={analysis.min.gpu} rec={analysis.rec.gpu}
                  meetsMin={analysis.meets.min_gpu} meetsRec={analysis.meets.rec_gpu} />
                <ComparisonRow label="Disco" user={analysis.user.storage}
                  min={analysis.min.storage} rec={analysis.rec.storage}
                  meetsMin={analysis.meets.min_storage} meetsRec={analysis.meets.rec_storage} />

                <div style={{ display: "flex", gap: 16, marginTop: 14 }}>
                  {[["#10b981", "Atende"], ["#ef4444", "Não atende"]].map(([c, l]) => (
                    <div key={l} style={{ display: "flex", alignItems: "center", gap: 5 }}>
                      <div style={{ width: 6, height: 6, borderRadius: "50%", background: c, boxShadow: `0 0 4px ${c}` }} />
                      <span style={{ color: "#475569", fontSize: 10 }}>{l}</span>
                    </div>
                  ))}
                </div>
              </div>
            </div>

            {/* Tips */}
            {analysis.tips?.length > 0 && (
              <motion.div initial={{ opacity: 0, y: 12 }} animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.3 }}
                style={{ background: "rgba(10,22,40,0.85)", backdropFilter: "blur(12px)",
                  border: "1px solid rgba(167,139,250,0.2)", borderRadius: 12, padding: 20 }}>
                <h3 style={{ color: "#a78bfa", fontSize: 12, textTransform: "uppercase",
                  letterSpacing: "0.12em", marginBottom: 14 }}>💡 Dicas do JARVIS</h3>
                <div style={{ display: "flex", flexDirection: "column", gap: 10 }}>
                  {analysis.tips.map((tip, i) => (
                    <div key={i} style={{ display: "flex", gap: 10, alignItems: "flex-start" }}>
                      <span style={{ color: "#a78bfa", fontSize: 14, flexShrink: 0, marginTop: 1 }}>▸</span>
                      <span style={{ color: "#94a3b8", fontSize: 13, lineHeight: 1.6 }}>{tip}</span>
                    </div>
                  ))}
                </div>
              </motion.div>
            )}

            <div style={{ textAlign: "center", marginTop: 20 }}>
              <button onClick={() => { setAnalysis(null); setQuery(""); setGameName(""); setCoverUrl(null); }}
                style={{ background: "transparent", border: "1px solid rgba(255,255,255,0.08)",
                  borderRadius: 8, padding: "8px 24px", color: "#475569",
                  fontSize: 13, cursor: "pointer" }}>
                Verificar outro jogo
              </button>
            </div>
          </motion.div>
        )}
      </AnimatePresence>

      {/* Trending / empty state */}
      {!analysis && !isLoading && !error && (
        <div>
          <div style={{ display: "flex", alignItems: "center", gap: 10, marginBottom: 16 }}>
            <span style={{ color: "#00d4ff", fontSize: 11, textTransform: "uppercase", letterSpacing: "0.12em" }}>
              🔥 Em Alta Esta Semana
            </span>
            <div style={{ flex: 1, height: 1, background: "rgba(0,212,255,0.08)" }} />
            <span style={{ color: "#334155", fontSize: 10 }}>Clique para verificar</span>
          </div>

          {trendingLoading ? (
            <div style={{ display: "grid", gridTemplateColumns: "repeat(5, 1fr)", gap: 12 }}>
              {[...Array(5)].map((_, i) => (
                <motion.div key={i}
                  animate={{ opacity: [0.3, 0.6, 0.3] }}
                  transition={{ duration: 1.5, repeat: Infinity, delay: i * 0.1 }}
                  style={{
                    aspectRatio: "2/3", borderRadius: 10,
                    background: "rgba(255,255,255,0.03)",
                    border: "1px solid rgba(255,255,255,0.05)",
                  }} />
              ))}
            </div>
          ) : trending.length > 0 ? (
            <div style={{ display: "grid", gridTemplateColumns: "repeat(5, 1fr)", gap: 12 }}>
              {trending.map((game, i) => (
                <motion.button
                  key={game.id}
                  initial={{ opacity: 0, y: 16 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ delay: i * 0.07 }}
                  onClick={() => runAnalysis(game.name)}
                  style={{
                    background: "transparent", border: "none", padding: 0,
                    cursor: "pointer", textAlign: "left",
                  }}
                >
                  <div style={{
                    borderRadius: 10, overflow: "hidden",
                    border: "1px solid rgba(255,255,255,0.07)",
                    transition: "border-color 0.2s, transform 0.2s",
                    position: "relative",
                  }}
                    onMouseEnter={e => {
                      (e.currentTarget as HTMLElement).style.borderColor = "rgba(0,212,255,0.4)";
                      (e.currentTarget as HTMLElement).style.transform = "translateY(-3px)";
                    }}
                    onMouseLeave={e => {
                      (e.currentTarget as HTMLElement).style.borderColor = "rgba(255,255,255,0.07)";
                      (e.currentTarget as HTMLElement).style.transform = "translateY(0)";
                    }}
                  >
                    {game.background_image ? (
                      <img src={game.background_image} alt={game.name}
                        style={{ width: "100%", aspectRatio: "16/9", objectFit: "cover", display: "block" }} />
                    ) : (
                      <div style={{
                        width: "100%", aspectRatio: "16/9",
                        background: `linear-gradient(135deg, rgba(0,212,255,0.1), rgba(0,212,255,0.05))`,
                        display: "flex", alignItems: "center", justifyContent: "center",
                        fontSize: 28, fontWeight: 900, color: "#00d4ff",
                      }}>
                        {game.name.charAt(0)}
                      </div>
                    )}
                    <div style={{ padding: "8px 10px", background: "rgba(5,13,26,0.95)" }}>
                      <p style={{
                        color: "#e2e8f0", fontSize: 12, fontWeight: 600,
                        overflow: "hidden", textOverflow: "ellipsis", whiteSpace: "nowrap",
                        marginBottom: 2,
                      }}>
                        {game.name}
                      </p>
                      <div style={{ display: "flex", alignItems: "center", gap: 4 }}>
                        <span style={{ color: "#475569", fontSize: 10 }}>
                          {game.genres?.[0]?.name || "Game"}
                        </span>
                        {game.rating > 0 && (
                          <>
                            <span style={{ color: "#1e293b", fontSize: 9 }}>·</span>
                            <div style={{ display: "flex", gap: 1 }}>
                              {[1, 2, 3, 4, 5].map((star) => (
                                <span key={star} style={{
                                  fontSize: 9,
                                  color: star <= Math.round(game.rating) ? "#f59e0b" : "#1e293b",
                                }}>★</span>
                              ))}
                            </div>
                            <span style={{ color: "#475569", fontSize: 9 }}>{game.rating.toFixed(1)}</span>
                          </>
                        )}
                      </div>
                    </div>
                    <div style={{
                      position: "absolute", top: 6, right: 6,
                      background: "rgba(0,212,255,0.15)", borderRadius: 4,
                      padding: "2px 6px", fontSize: 9, color: "#00d4ff",
                      letterSpacing: "0.08em", fontWeight: 700,
                    }}>
                      #{i + 1}
                    </div>
                  </div>
                </motion.button>
              ))}
            </div>
          ) : (
            <div style={{ textAlign: "center", padding: "40px 0", color: "#334155", fontSize: 13 }}>
              Nenhum jogo em alta encontrado. Digite um nome acima.
            </div>
          )}

          <p style={{ color: "#1e293b", fontSize: 11, textAlign: "center", marginTop: 20 }}>
            JARVIS usa suas specs reais para cada análise
          </p>
        </div>
      )}
    </DashboardLayout>
  );
}
