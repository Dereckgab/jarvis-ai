"use client";

import React, { useState } from "react";
import { useRouter } from "next/navigation";
import { motion } from "framer-motion";
import { useAuth } from "@/context/AuthContext";
import { Button } from "@/components/Button";
import { Input } from "@/components/Input";

export default function LoginPage() {
  const router = useRouter();
  const { login } = useAuth();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setIsLoading(true);

    try {
      const success = await login(email, password);
      if (success) {
        router.push("/dashboard");
      } else {
        setError("Invalid credentials. Access denied.");
      }
    } catch (err) {
      setError("Connection error. Retrying...");
      console.error(err);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div
      className="min-h-screen flex items-center justify-center p-4 relative overflow-hidden"
      style={{ background: "var(--bg-primary)" }}
    >
      {/* Background grid */}
      <div
        className="absolute inset-0 opacity-[0.03]"
        style={{
          backgroundImage:
            "linear-gradient(var(--accent-cyan) 1px, transparent 1px), linear-gradient(90deg, var(--accent-cyan) 1px, transparent 1px)",
          backgroundSize: "40px 40px",
        }}
      />

      {/* Ambient glow */}
      <div
        className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[500px] h-[500px] rounded-full opacity-5 blur-3xl pointer-events-none"
        style={{ background: "var(--accent-cyan)" }}
      />

      <motion.div
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        transition={{ duration: 0.5 }}
        className="w-full max-w-md relative z-10"
      >
        {/* Logo */}
        <div className="text-center mb-10">
          <motion.div
            initial={{ opacity: 0, scale: 0.5 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.6 }}
            className="inline-flex items-center justify-center mb-6 relative"
          >
            {[80, 56, 32].map((size, i) => (
              <motion.div
                key={i}
                animate={{ rotate: i % 2 === 0 ? 360 : -360 }}
                transition={{ duration: 8 + i * 3, repeat: Infinity, ease: "linear" }}
                className="rounded-full border absolute"
                style={{
                  width: size,
                  height: size,
                  borderColor: `rgba(0,212,255,${0.5 - i * 0.12})`,
                }}
              />
            ))}
            <div
              className="w-4 h-4 rounded-full relative z-10"
              style={{
                background: "var(--accent-cyan)",
                boxShadow: "0 0 16px var(--accent-cyan), 0 0 40px rgba(0,212,255,0.3)",
              }}
            />
          </motion.div>

          <motion.h1
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.2 }}
            className="text-4xl font-bold tracking-[0.3em] mb-2"
            style={{ color: "var(--accent-cyan)", textShadow: "0 0 20px rgba(0,212,255,0.5)" }}
          >
            JARVIS
          </motion.h1>
          <motion.p
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ delay: 0.35 }}
            className="text-sm tracking-widest uppercase"
            style={{ color: "var(--text-muted)" }}
          >
            Identification Required
          </motion.p>
        </div>

        {/* Card */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2 }}
          className="glass-card p-8"
        >
          <form onSubmit={handleSubmit} className="space-y-5">
            <Input
              label="Email"
              type="email"
              placeholder="user@domain.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
              disabled={isLoading}
            />

            <Input
              label="Password"
              type="password"
              placeholder="••••••••"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              disabled={isLoading}
            />

            {error && (
              <motion.div
                initial={{ opacity: 0, y: -8 }}
                animate={{ opacity: 1, y: 0 }}
                className="p-3 rounded-lg text-sm"
                style={{
                  background: "rgba(239,68,68,0.08)",
                  border: "1px solid rgba(239,68,68,0.3)",
                  color: "#f87171",
                }}
              >
                ⚠ {error}
              </motion.div>
            )}

            <Button
              type="submit"
              variant="primary"
              size="lg"
              isLoading={isLoading}
              disabled={isLoading}
              className="w-full mt-2 tracking-widest"
            >
              {isLoading ? "Authenticating..." : "Access System"}
            </Button>
          </form>

          <div className="mt-6 text-center">
            <p className="text-sm" style={{ color: "var(--text-muted)" }}>
              No access?{" "}
              <a
                href="/register"
                className="font-semibold transition-colors hover:opacity-80"
                style={{ color: "var(--accent-cyan)" }}
              >
                Request Account
              </a>
            </p>
          </div>
        </motion.div>

        {/* Bottom decoration */}
        <motion.p
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.6 }}
          className="text-center text-xs mt-6 tracking-widest"
          style={{ color: "var(--text-muted)" }}
        >
          JARVIS AI • SECURE CONNECTION
        </motion.p>
      </motion.div>
    </div>
  );
}
