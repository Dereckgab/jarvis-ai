"use client";

import React, { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { DashboardLayout } from "@/components/DashboardLayout";
import { useAuth } from "@/context/AuthContext";
import { apiClient } from "@/lib/api-client";

interface SystemInfo {
  cpu_name: string;
  cpu_cores: number;
  cpu_threads: number;
  cpu_frequency: number;
  cpu_percent: number;
  total_memory_gb: number;
  used_memory_gb: number;
  memory_percent: number;
  total_disk_gb: number;
  used_disk_gb: number;
  disk_percent: number;
  gpu_name: string;
  gpu_temperature_c: number;
  gpu_utilization_percent: number;
  bytes_sent_sec: number;
  bytes_recv_sec: number;
  os_platform: string;
  os_version: string;
}

function ProgressRing({ percentage, color }: { percentage: number; color: string }) {
  const r = 30;
  const circ = 2 * Math.PI * r;
  const offset = circ - (Math.min(percentage, 100) / 100) * circ;
  return (
    <svg width="72" height="72">
      <circle cx="36" cy="36" r={r} stroke="rgba(255,255,255,0.05)" strokeWidth="5" fill="none" />
      <motion.circle
        cx="36" cy="36" r={r}
        stroke={color} strokeWidth="5" fill="none"
        strokeLinecap="round"
        initial={{ strokeDashoffset: circ }}
        animate={{ strokeDashoffset: offset }}
        transition={{ duration: 1, ease: "easeOut" }}
        style={{
          strokeDasharray: circ,
          transform: "rotate(-90deg)",
          transformOrigin: "36px 36px",
          filter: `drop-shadow(0 0 4px ${color})`,
        }}
      />
      <text x="36" y="36" textAnchor="middle" dominantBaseline="middle"
        fontSize="11" fontWeight="bold" fill={color}>
        {percentage.toFixed(0)}%
      </text>
    </svg>
  );
}

function UsageCard({ label, icon, color, glow, percentage, detail }: {
  label: string; icon: string; color: string; glow: string; percentage: number; detail: string;
}) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 16 }} animate={{ opacity: 1, y: 0 }}
      whileHover={{ scale: 1.02 }}
      onMouseEnter={e => (e.currentTarget.style.boxShadow = `0 0 18px ${glow}`)}
      onMouseLeave={e => (e.currentTarget.style.boxShadow = "none")}
      style={{
        background: "rgba(10,22,40,0.85)", backdropFilter: "blur(12px)",
        border: "1px solid rgba(0,212,255,0.15)", borderRadius: 12, padding: 20,
        display: "flex", alignItems: "center", gap: 16, transition: "box-shadow 0.2s",
      }}
    >
      <ProgressRing percentage={percentage} color={color} />
      <div style={{ flex: 1 }}>
        <div style={{ display: "flex", alignItems: "center", gap: 6, marginBottom: 4 }}>
          <span style={{ color, fontSize: 14 }}>{icon}</span>
          <span style={{ color: "#94a3b8", fontSize: 11, textTransform: "uppercase", letterSpacing: "0.1em" }}>{label}</span>
        </div>
        <p style={{ color: "#e2e8f0", fontSize: 13, fontWeight: 500 }}>{detail}</p>
        <div style={{ marginTop: 8, height: 4, borderRadius: 9999, background: "rgba(255,255,255,0.05)", overflow: "hidden" }}>
          <motion.div
            initial={{ width: 0 }}
            animate={{ width: `${Math.min(percentage, 100)}%` }}
            transition={{ duration: 0.8, ease: "easeOut" }}
            style={{ height: 4, borderRadius: 9999, background: color, boxShadow: `0 0 6px ${glow}` }}
          />
        </div>
      </div>
    </motion.div>
  );
}

function SpecRow({ label, value, accent }: { label: string; value: string; accent?: string }) {
  return (
    <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center",
      padding: "9px 0", borderBottom: "1px solid rgba(255,255,255,0.04)" }}>
      <span style={{ color: "#475569", fontSize: 11, textTransform: "uppercase", letterSpacing: "0.08em" }}>{label}</span>
      <span style={{ color: accent || "#e2e8f0", fontSize: 13, fontWeight: 500 }}>{value}</span>
    </div>
  );
}

function SpecCard({ title, color, rows }: { title: string; color: string; rows: { label: string; value: string }[] }) {
  return (
    <div style={{
      background: "rgba(10,22,40,0.85)", backdropFilter: "blur(12px)",
      border: `1px solid ${color}22`, borderRadius: 12, padding: 20,
    }}>
      <h3 style={{ color, fontSize: 11, textTransform: "uppercase", letterSpacing: "0.12em", marginBottom: 12 }}>
        {title}
      </h3>
      {rows.map(r => <SpecRow key={r.label} label={r.label} value={r.value} />)}
    </div>
  );
}

function SkeletonCard() {
  return (
    <div style={{ background: "rgba(10,22,40,0.6)", border: "1px solid rgba(0,212,255,0.08)",
      borderRadius: 12, padding: 20, height: 96 }}>
      <motion.div animate={{ opacity: [0.3, 0.6, 0.3] }} transition={{ duration: 1.5, repeat: Infinity }}
        style={{ height: 12, width: "50%", background: "rgba(255,255,255,0.05)", borderRadius: 6, marginBottom: 12 }} />
      <motion.div animate={{ opacity: [0.3, 0.6, 0.3] }} transition={{ duration: 1.5, repeat: Infinity, delay: 0.2 }}
        style={{ height: 4, width: "100%", background: "rgba(255,255,255,0.05)", borderRadius: 6 }} />
    </div>
  );
}

export default function DashboardPage() {
  const { token } = useAuth();
  const [info, setInfo] = useState<SystemInfo | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!token) return;
    const fetch = async () => {
      try {
        const res = await apiClient.getLatestSystemInfo(token);
        if (res.success && res.data) setInfo(res.data);
      } catch {
        /* ignore */
      } finally {
        setLoading(false);
      }
    };
    fetch();
    const iv = setInterval(fetch, 5000);
    return () => clearInterval(iv);
  }, [token]);

  const freeRam = (info?.total_memory_gb ?? 0) - (info?.used_memory_gb ?? 0);
  const freeDisk = (info?.total_disk_gb ?? 0) - (info?.used_disk_gb ?? 0);

  return (
    <DashboardLayout title="Overview">

      {/* Usage cards — 2x2 grid */}
      <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fit, minmax(260px, 1fr))", gap: 16, marginBottom: 24 }}>
        {loading ? (
          [...Array(4)].map((_, i) => <SkeletonCard key={i} />)
        ) : info ? (
          <>
            <UsageCard label="CPU" icon="⬡" color="#00d4ff" glow="rgba(0,212,255,0.35)"
              percentage={info.cpu_percent}
              detail={info.cpu_name || "N/A"} />
            <UsageCard label="Memória" icon="◈" color="#10b981" glow="rgba(16,185,129,0.35)"
              percentage={info.memory_percent}
              detail={`${info.used_memory_gb?.toFixed(1)} / ${info.total_memory_gb?.toFixed(1)} GB`} />
            <UsageCard label="Disco" icon="◉" color="#f59e0b" glow="rgba(245,158,11,0.35)"
              percentage={info.disk_percent}
              detail={`${info.used_disk_gb?.toFixed(0)} GB usados · ${freeDisk.toFixed(0)} GB livres`} />
            <UsageCard label="GPU" icon="◆" color="#a78bfa" glow="rgba(167,139,250,0.35)"
              percentage={info.gpu_utilization_percent}
              detail={info.gpu_name || "Sem GPU dedicada"} />
          </>
        ) : (
          <div style={{ gridColumn: "1 / -1", textAlign: "center", color: "#475569", padding: 32 }}>
            Sem dados de sistema disponíveis.
          </div>
        )}
      </div>

      {/* Spec cards */}
      {info && (
        <motion.div
          initial={{ opacity: 0, y: 16 }} animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2 }}
          style={{ display: "grid", gridTemplateColumns: "repeat(auto-fit, minmax(240px, 1fr))", gap: 16 }}
        >
          <SpecCard title="Processador" color="#00d4ff" rows={[
            { label: "Modelo", value: info.cpu_name || "N/A" },
            { label: "Núcleos", value: String(info.cpu_cores || 0) },
            { label: "Threads", value: String(info.cpu_threads || 0) },
            { label: "Frequência", value: `${info.cpu_frequency?.toFixed(0) || 0} MHz` },
          ]} />

          <SpecCard title="Memória RAM" color="#10b981" rows={[
            { label: "Total", value: `${info.total_memory_gb?.toFixed(1) || 0} GB` },
            { label: "Em uso", value: `${info.used_memory_gb?.toFixed(1) || 0} GB` },
            { label: "Livre", value: `${freeRam.toFixed(1)} GB` },
          ]} />

          <SpecCard title="Armazenamento" color="#f59e0b" rows={[
            { label: "Total", value: `${info.total_disk_gb?.toFixed(0) || 0} GB` },
            { label: "Em uso", value: `${info.used_disk_gb?.toFixed(0) || 0} GB` },
            { label: "Livre", value: `${freeDisk.toFixed(0)} GB` },
          ]} />

          <SpecCard title="GPU & Rede & Sistema" color="#a78bfa" rows={[
            { label: "GPU", value: info.gpu_name || "N/A" },
            { label: "Temp GPU", value: info.gpu_temperature_c ? `${info.gpu_temperature_c} °C` : "N/A" },
            { label: "Upload", value: `${((info.bytes_sent_sec || 0) / 1024).toFixed(1)} KB/s` },
            { label: "Download", value: `${((info.bytes_recv_sec || 0) / 1024).toFixed(1)} KB/s` },
            { label: "Sistema", value: info.os_platform || "N/A" },
          ]} />
        </motion.div>
      )}
    </DashboardLayout>
  );
}
