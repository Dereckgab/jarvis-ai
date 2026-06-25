"use client";

import React from "react";
import { motion } from "framer-motion";
import clsx from "clsx";

interface CardProps {
  className?: string;
  children: React.ReactNode;
  hover?: boolean;
  glow?: boolean;
}

export function Card({ className, children, hover = true, glow = false }: CardProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.4 }}
      whileHover={hover ? { y: -2, boxShadow: "0 0 24px rgba(0,212,255,0.2)" } : {}}
      className={clsx(
        "glass-card transition-all duration-300",
        glow && "glow-pulse",
        className
      )}
    >
      {children}
    </motion.div>
  );
}
