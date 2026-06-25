"use client";

import React from "react";
import { useRouter, usePathname } from "next/navigation";
import { motion } from "framer-motion";
import { useAuth } from "@/context/AuthContext";
import { Button } from "@/components/Button";

interface DashboardLayoutProps {
  children: React.ReactNode;
  title: string;
}

const navLinks = [
  { href: "/dashboard", label: "Dashboard", icon: "⬡" },
  { href: "/dashboard/chat", label: "Chat", icon: "◈" },
  { href: "/dashboard/games", label: "Games", icon: "◆" },
];

export function DashboardLayout({ children, title }: DashboardLayoutProps) {
  const router = useRouter();
  const pathname = usePathname();
  const { user, logout, isAuthenticated, isLoading } = useAuth();

  React.useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.push("/login");
    }
  }, [isAuthenticated, isLoading, router]);

  if (isLoading) return (
    <div className="min-h-screen flex items-center justify-center" style={{ background: "var(--bg-primary)" }}>
      <div style={{ color: "var(--accent-cyan)" }}>Loading...</div>
    </div>
  );

  if (!isAuthenticated) return null;

  return (
    <div className="min-h-screen" style={{ background: "var(--bg-primary)" }}>
      {/* Header */}
      <header
        className="border-b sticky top-0 z-40"
        style={{
          background: "rgba(5, 13, 26, 0.9)",
          backdropFilter: "blur(12px)",
          borderColor: "var(--border-subtle)",
        }}
      >
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-3 flex justify-between items-center">
          <motion.div
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.4 }}
            className="flex items-center gap-3"
          >
            {/* JARVIS Logo */}
            <div className="relative w-8 h-8">
              <motion.div
                animate={{ rotate: 360 }}
                transition={{ duration: 12, repeat: Infinity, ease: "linear" }}
                className="w-8 h-8 rounded-full border border-cyan-500/50 absolute"
              />
              <motion.div
                animate={{ rotate: -360 }}
                transition={{ duration: 8, repeat: Infinity, ease: "linear" }}
                className="w-5 h-5 rounded-full border border-cyan-400/30 absolute top-1.5 left-1.5"
              />
              <div
                className="w-2 h-2 rounded-full absolute top-3 left-3"
                style={{ background: "var(--accent-cyan)", boxShadow: "0 0 6px var(--accent-cyan)" }}
              />
            </div>
            <div>
              <h1
                className="text-lg font-bold tracking-[0.2em]"
                style={{ color: "var(--accent-cyan)", textShadow: "0 0 10px rgba(0,212,255,0.4)" }}
              >
                JARVIS
              </h1>
              <p className="text-xs" style={{ color: "var(--text-muted)" }}>
                Welcome, {user?.name}
              </p>
            </div>
          </motion.div>

          <Button
            variant="ghost"
            size="sm"
            onClick={() => {
              logout();
              router.push("/login");
            }}
          >
            Logout
          </Button>
        </div>
      </header>

      {/* Navigation */}
      <nav
        className="border-b"
        style={{
          background: "rgba(5, 13, 26, 0.95)",
          borderColor: "var(--border-subtle)",
          position: "sticky",
          top: "57px",
          zIndex: 39,
          backdropFilter: "blur(12px)",
        }}
      >
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex space-x-1">
            {navLinks.map((link) => {
              const isActive = pathname === link.href;
              return (
                <a
                  key={link.href}
                  href={link.href}
                  className="relative px-4 py-3 text-sm font-medium transition-all duration-200 flex items-center gap-2"
                  style={{
                    color: isActive ? "var(--accent-cyan)" : "var(--text-secondary)",
                  }}
                >
                  <span className="text-xs opacity-60">{link.icon}</span>
                  {link.label}
                  {isActive && (
                    <motion.div
                      layoutId="nav-indicator"
                      className="absolute bottom-0 left-0 right-0 h-0.5"
                      style={{ background: "var(--accent-cyan)", boxShadow: "0 0 8px var(--accent-cyan)" }}
                    />
                  )}
                </a>
              );
            })}
          </div>
        </div>
      </nav>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <motion.div
          initial={{ opacity: 0, y: 16 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.4 }}
        >
          <div className="flex items-center gap-3 mb-8">
            <div
              className="h-6 w-1 rounded-full"
              style={{ background: "var(--accent-cyan)", boxShadow: "0 0 8px var(--accent-cyan)" }}
            />
            <h2
              className="text-2xl font-bold tracking-wide"
              style={{ color: "var(--text-primary)" }}
            >
              {title}
            </h2>
          </div>
          {children}
        </motion.div>
      </main>
    </div>
  );
}
