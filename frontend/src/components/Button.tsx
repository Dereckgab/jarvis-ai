"use client";

import React, { type ButtonHTMLAttributes, type ReactNode } from "react";
import { motion, type MotionProps } from "framer-motion";
import clsx from "clsx";

type ButtonProps = Omit<ButtonHTMLAttributes<HTMLButtonElement>, "onDrag" | "onDragStart" | "onDragEnd" | "onDragEnter" | "onDragLeave" | "onDragOver"> &
  MotionProps & {
    variant?: "primary" | "secondary" | "danger" | "ghost";
    size?: "sm" | "md" | "lg";
    isLoading?: boolean;
    children: ReactNode;
  };

export function Button({
  variant = "primary",
  size = "md",
  isLoading = false,
  disabled,
  className,
  children,
  ...props
}: ButtonProps) {
  const baseStyles = "font-semibold rounded-lg transition-all duration-200 flex items-center justify-center gap-2 tracking-wide";

  const variantStyles = {
    primary: "bg-cyan-500/10 text-cyan-400 border border-cyan-500/40 hover:bg-cyan-500/20 hover:border-cyan-400 hover:shadow-[0_0_12px_rgba(0,212,255,0.3)]",
    secondary: "bg-slate-700/50 text-slate-300 border border-slate-600/50 hover:bg-slate-600/50 hover:border-slate-500",
    danger: "bg-red-500/10 text-red-400 border border-red-500/40 hover:bg-red-500/20 hover:border-red-400",
    ghost: "bg-transparent text-slate-400 hover:text-cyan-400 hover:bg-cyan-500/10",
  };

  const sizeStyles = {
    sm: "px-3 py-1.5 text-sm",
    md: "px-4 py-2 text-base",
    lg: "px-6 py-3 text-lg",
  };

  return (
    <motion.button
      whileHover={{ scale: 1.02 }}
      whileTap={{ scale: 0.97 }}
      disabled={disabled || isLoading}
      className={clsx(
        baseStyles,
        variantStyles[variant],
        sizeStyles[size],
        (disabled || isLoading) && "opacity-40 cursor-not-allowed",
        className
      )}
      {...props}
    >
      {isLoading && (
        <motion.div
          animate={{ rotate: 360 }}
          transition={{ duration: 1, repeat: Infinity, ease: "linear" }}
          className="w-4 h-4 border-2 border-cyan-400 border-t-transparent rounded-full"
        />
      )}
      {children}
    </motion.button>
  );
}
