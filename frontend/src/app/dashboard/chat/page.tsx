"use client";

import React, { useRef, useEffect, useState, useCallback } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { DashboardLayout } from "@/components/DashboardLayout";
import { Button } from "@/components/Button";
import { Input } from "@/components/Input";
import { useAuth } from "@/context/AuthContext";
import { apiClient } from "@/lib/api-client";

interface Message {
  id: string;
  role: "user" | "assistant";
  content: string;
  timestamp: Date;
  showFollowUp?: boolean;
  followUpDone?: boolean;
}

const GAME_KEYWORDS = [
  "jogo", "game", "roda", "rodar", "compatível", "compatibilidade",
  "fps", "requisitos", "mínimo", "recomendado", "gráficos", "gpu", "pc",
  "instalar", "jogar", "league", "minecraft", "fortnite", "gta", "valorant",
  "steam", "epic", "processador", "memória", "placa de vídeo", "funcioná",
  "consegue rodar", "consegue jogar", "meu pc roda", "meu computador",
];

function isGameRelated(text: string): boolean {
  const lower = text.toLowerCase();
  return GAME_KEYWORDS.some((kw) => lower.includes(kw));
}

function MessageBubble({
  message,
  onFollowUp,
}: {
  message: Message;
  onFollowUp: (id: string) => void;
}) {
  const isUser = message.role === "user";
  return (
    <motion.div
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3 }}
    >
      <div className={`flex ${isUser ? "justify-end" : "justify-start"}`}>
        {!isUser && (
          <div
            style={{
              width: 28, height: 28, borderRadius: "50%",
              border: "1px solid rgba(0,212,255,0.3)",
              background: "rgba(0,212,255,0.08)",
              display: "flex", alignItems: "center", justifyContent: "center",
              marginRight: 8, marginTop: 2, flexShrink: 0,
            }}
          >
            <span style={{ color: "#00d4ff", fontSize: 11, fontWeight: 700 }}>J</span>
          </div>
        )}
        <div style={{ maxWidth: "75%" }}>
          <div
            style={
              isUser
                ? {
                    background: "rgba(0,212,255,0.1)",
                    border: "1px solid rgba(0,212,255,0.25)",
                    color: "#e2e8f0",
                    borderRadius: "16px 16px 4px 16px",
                    padding: "10px 14px",
                    fontSize: 14,
                    lineHeight: 1.6,
                  }
                : {
                    background: "rgba(255,255,255,0.03)",
                    border: "1px solid rgba(255,255,255,0.07)",
                    color: "#94a3b8",
                    borderRadius: "16px 16px 16px 4px",
                    padding: "10px 14px",
                    fontSize: 14,
                    lineHeight: 1.6,
                    whiteSpace: "pre-wrap",
                  }
            }
          >
            {message.content.split("\n").map((line, i) => {
              const isVerdict = /^[✅❌⚠️]/.test(line.trim());
              return (
                <p key={i} style={{
                  lineHeight: 1.6, margin: 0,
                  fontSize: isVerdict ? 15 : 14,
                  fontWeight: isVerdict ? 700 : 400,
                  color: isVerdict
                    ? line.includes("✅") ? "#10b981"
                    : line.includes("❌") ? "#ef4444"
                    : "#f59e0b"
                    : undefined,
                  marginBottom: isVerdict ? 6 : 2,
                }}>
                  {line || " "}
                </p>
              );
            })}
          </div>
          <p style={{ color: "#334155", fontSize: 10, marginTop: 4, textAlign: isUser ? "right" : "left" }}>
            {message.timestamp.toLocaleTimeString("pt-BR", { hour: "2-digit", minute: "2-digit" })}
          </p>

          {/* Follow-up prompt */}
          {message.showFollowUp && !message.followUpDone && (
            <motion.div
              initial={{ opacity: 0, y: 6 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.4 }}
              style={{
                marginTop: 10,
                background: "rgba(0,212,255,0.05)",
                border: "1px solid rgba(0,212,255,0.18)",
                borderRadius: 10,
                padding: "10px 14px",
              }}
            >
              <p style={{ color: "#94a3b8", fontSize: 13, marginBottom: 10 }}>
                Quer saber mais detalhes sobre a compatibilidade?
              </p>
              <div style={{ display: "flex", gap: 8 }}>
                <button
                  onClick={() => onFollowUp(message.id)}
                  style={{
                    padding: "6px 18px",
                    background: "rgba(0,212,255,0.12)",
                    border: "1px solid rgba(0,212,255,0.4)",
                    borderRadius: 6, color: "#00d4ff",
                    fontSize: 13, fontWeight: 600, cursor: "pointer",
                    letterSpacing: "0.04em",
                  }}
                >
                  Sim
                </button>
                <button
                  onClick={() => onFollowUp("dismiss-" + message.id)}
                  style={{
                    padding: "6px 18px",
                    background: "transparent",
                    border: "1px solid rgba(255,255,255,0.08)",
                    borderRadius: 6, color: "#475569",
                    fontSize: 13, cursor: "pointer",
                  }}
                >
                  Não
                </button>
              </div>
            </motion.div>
          )}
        </div>
      </div>
    </motion.div>
  );
}

export default function ChatPage() {
  const { token } = useAuth();
  const [messages, setMessages] = useState<Message[]>([]);
  const [input, setInput] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const messagesEndRef = useRef<HTMLDivElement>(null);
  const lastUserPrompt = useRef<string>("");

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const sendPrompt = useCallback(
    async (prompt: string, opts?: { showFollowUpOnResponse?: boolean }) => {
      if (!token) return;
      setIsLoading(true);

      try {
        const response = await apiClient.chatCompletion(prompt, token);
        if (response.success && response.data) {
          const text = response.data.response;
          const shouldShowFollowUp =
            opts?.showFollowUpOnResponse && isGameRelated(prompt) && isGameRelated(text);

          const assistantMsg: Message = {
            id: Date.now().toString(),
            role: "assistant",
            content: text,
            timestamp: new Date(),
            showFollowUp: shouldShowFollowUp,
            followUpDone: false,
          };
          setMessages((prev) => [...prev, assistantMsg]);
        }
      } catch (err) {
        console.error("Chat error:", err);
      } finally {
        setIsLoading(false);
      }
    },
    [token]
  );

  const handleSendMessage = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!input.trim() || !token || isLoading) return;

    const text = input.trim();
    lastUserPrompt.current = text;

    const userMsg: Message = {
      id: Date.now().toString(),
      role: "user",
      content: text,
      timestamp: new Date(),
    };
    setMessages((prev) => [...prev, userMsg]);
    setInput("");

    // If game-related question, wrap prompt to force SIM/NÃO/DEPENDE verdict first
    const gameQuestion = isGameRelated(text);
    const finalPrompt = gameQuestion
      ? `${text}\n\n[INSTRUÇÃO: Responda começando obrigatoriamente com "✅ SIM", "❌ NÃO" ou "⚠️ DEPENDE" em destaque na primeira linha, depois dê uma explicação curta e direta em português.]`
      : text;

    await sendPrompt(finalPrompt, { showFollowUpOnResponse: gameQuestion });
  };

  const handleFollowUp = useCallback(
    async (messageId: string) => {
      // Dismiss case
      if (messageId.startsWith("dismiss-")) {
        const originalId = messageId.replace("dismiss-", "");
        setMessages((prev) =>
          prev.map((m) => (m.id === originalId ? { ...m, followUpDone: true } : m))
        );
        return;
      }

      // Mark as done
      setMessages((prev) =>
        prev.map((m) => (m.id === messageId ? { ...m, followUpDone: true } : m))
      );

      // Send follow-up as a visible user message
      const followUpText =
        "Me dê uma análise detalhada: quais componentes do meu PC atendem os requisitos mínimos e recomendados do jogo, o que gera gargalo, configurações ideais para ter o melhor desempenho possível, e se vale a pena jogar assim.";

      const userMsg: Message = {
        id: Date.now().toString(),
        role: "user",
        content: "Sim, quero mais detalhes!",
        timestamp: new Date(),
      };
      setMessages((prev) => [...prev, userMsg]);

      await sendPrompt(followUpText, { showFollowUpOnResponse: false });
    },
    [sendPrompt]
  );

  const suggestedQuestions = [
    "GTA V roda no meu PC?",
    "Qual FPS vou ter no Fortnite?",
    "Minecraft roda bem aqui?",
    "O que devo atualizar primeiro?",
  ];

  return (
    <DashboardLayout title="Chat">
      <div
        style={{
          background: "rgba(10,22,40,0.85)",
          backdropFilter: "blur(12px)",
          border: "1px solid rgba(0,212,255,0.15)",
          borderRadius: 14,
          display: "flex",
          flexDirection: "column",
          height: "calc(100vh - 220px)",
          minHeight: 500,
          overflow: "hidden",
        }}
      >
        {/* Messages */}
        <div style={{ flex: 1, overflowY: "auto", padding: "20px 24px", display: "flex", flexDirection: "column", gap: 16 }}>
          {messages.length === 0 ? (
            <div style={{ display: "flex", flexDirection: "column", alignItems: "center", justifyContent: "center", height: "100%", gap: 12 }}>
              <motion.div
                animate={{ opacity: [0.4, 1, 0.4] }}
                transition={{ duration: 2.5, repeat: Infinity }}
                style={{
                  width: 48, height: 48, borderRadius: "50%",
                  border: "1px solid rgba(0,212,255,0.3)",
                  background: "rgba(0,212,255,0.05)",
                  display: "flex", alignItems: "center", justifyContent: "center",
                  fontSize: 20, color: "#00d4ff",
                }}
              >
                ◈
              </motion.div>
              <p style={{ color: "#475569", fontSize: 14, letterSpacing: "0.05em" }}>
                JARVIS está pronto. Comece uma conversa.
              </p>

              {/* Suggested questions */}
              <div style={{ display: "flex", flexWrap: "wrap", gap: 8, justifyContent: "center", marginTop: 12 }}>
                {suggestedQuestions.map((q) => (
                  <button
                    key={q}
                    onClick={() => setInput(q)}
                    style={{
                      padding: "6px 14px",
                      background: "rgba(0,212,255,0.05)",
                      border: "1px solid rgba(0,212,255,0.15)",
                      borderRadius: 20, color: "#475569",
                      fontSize: 12, cursor: "pointer",
                    }}
                    onMouseEnter={e => { (e.target as HTMLButtonElement).style.borderColor = "rgba(0,212,255,0.4)"; (e.target as HTMLButtonElement).style.color = "#94a3b8"; }}
                    onMouseLeave={e => { (e.target as HTMLButtonElement).style.borderColor = "rgba(0,212,255,0.15)"; (e.target as HTMLButtonElement).style.color = "#475569"; }}
                  >
                    {q}
                  </button>
                ))}
              </div>
            </div>
          ) : (
            messages.map((message) => (
              <MessageBubble key={message.id} message={message} onFollowUp={handleFollowUp} />
            ))
          )}

          {isLoading && (
            <motion.div initial={{ opacity: 0, y: 8 }} animate={{ opacity: 1, y: 0 }}
              className="flex justify-start items-center gap-2">
              <div style={{
                width: 28, height: 28, borderRadius: "50%",
                border: "1px solid rgba(0,212,255,0.3)",
                background: "rgba(0,212,255,0.08)",
                display: "flex", alignItems: "center", justifyContent: "center",
                flexShrink: 0,
              }}>
                <span style={{ color: "#00d4ff", fontSize: 11, fontWeight: 700 }}>J</span>
              </div>
              <div style={{
                background: "rgba(255,255,255,0.03)",
                border: "1px solid rgba(255,255,255,0.07)",
                borderRadius: "16px 16px 16px 4px",
                padding: "12px 16px",
              }}>
                <div style={{ display: "flex", gap: 6 }}>
                  {[0, 0.15, 0.3].map((delay, i) => (
                    <motion.div key={i}
                      animate={{ y: [0, -6, 0] }}
                      transition={{ duration: 0.7, repeat: Infinity, delay }}
                      style={{ width: 6, height: 6, borderRadius: "50%", background: "#00d4ff" }}
                    />
                  ))}
                </div>
              </div>
            </motion.div>
          )}
          <div ref={messagesEndRef} />
        </div>

        {/* Divider */}
        <div style={{ borderTop: "1px solid rgba(255,255,255,0.05)" }} />

        {/* Input */}
        <form onSubmit={handleSendMessage} style={{ padding: "14px 20px", display: "flex", gap: 10 }}>
          <Input
            type="text"
            placeholder="Pergunte ao JARVIS sobre jogos, seu PC, ou qualquer coisa..."
            value={input}
            onChange={(e) => setInput(e.target.value)}
            disabled={isLoading}
            className="flex-1"
          />
          <Button type="submit" variant="primary" disabled={isLoading || !input.trim()} isLoading={isLoading}>
            Enviar
          </Button>
        </form>
      </div>
    </DashboardLayout>
  );
}
