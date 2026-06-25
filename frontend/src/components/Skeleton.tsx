"use client";

import React from "react";
import { motion } from "framer-motion";
import clsx from "clsx";

interface SkeletonProps {
  className?: string;
  count?: number;
}

export function Skeleton({ className, count = 1 }: SkeletonProps) {
  return (
    <>
      {Array.from({ length: count }).map((_, i) => (
        <motion.div
          key={i}
          animate={{ opacity: [0.5, 1, 0.5] }}
          transition={{ duration: 1.5, repeat: Infinity }}
          className={clsx("bg-gray-200 rounded", className)}
        />
      ))}
    </>
  );
}
