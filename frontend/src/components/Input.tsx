"use client";

import React from "react";
import { motion } from "framer-motion";
import clsx from "clsx";

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
  helperText?: string;
}

export function Input({
  label,
  error,
  helperText,
  className,
  ...props
}: InputProps) {
  const [isFocused, setIsFocused] = React.useState(false);

  return (
    <motion.div
      initial={{ opacity: 0, y: -10 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.3 }}
      className="w-full"
    >
      {label && (
        <label className="block text-sm font-medium text-slate-400 mb-1.5 tracking-wide uppercase text-xs">
          {label}
        </label>
      )}

      <motion.div
        animate={{
          borderColor: isFocused
            ? "rgba(0,212,255,0.6)"
            : error
            ? "rgba(239,68,68,0.6)"
            : "rgba(0,212,255,0.15)",
          boxShadow: isFocused
            ? "0 0 0 2px rgba(0,212,255,0.1), 0 0 12px rgba(0,212,255,0.15)"
            : "none",
        }}
        transition={{ duration: 0.2 }}
        className="relative border rounded-lg overflow-hidden"
      >
        <input
          onFocus={() => setIsFocused(true)}
          onBlur={() => setIsFocused(false)}
          className={clsx(
            "w-full px-4 py-2.5 bg-slate-900/60 text-slate-100 placeholder-slate-600 outline-none",
            "transition-colors duration-200 text-sm",
            className
          )}
          {...props}
        />

        <motion.div
          initial={{ scaleX: 0 }}
          animate={{ scaleX: isFocused ? 1 : 0 }}
          transition={{ duration: 0.3 }}
          className="absolute bottom-0 left-0 right-0 h-0.5 bg-cyan-400 origin-left"
        />
      </motion.div>

      {error && (
        <motion.p
          initial={{ opacity: 0, y: -5 }}
          animate={{ opacity: 1, y: 0 }}
          className="mt-1 text-xs text-red-400"
        >
          {error}
        </motion.p>
      )}

      {helperText && !error && (
        <p className="mt-1 text-xs text-slate-500">{helperText}</p>
      )}
    </motion.div>
  );
}
