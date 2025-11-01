"use client"

import { cn } from "@/lib/utils.ts"
import { usePrivacy } from "@/components/provider/PrivacyProvider.tsx"
import React, { CSSProperties, MouseEventHandler } from "react"

export default function PrivacyFilter({
  children,
  className,
  onClick,
  style,
  mode = "redact",
}: {
  children: React.ReactNode
  className?: string
  onClick?: MouseEventHandler<HTMLDivElement>
  style?: CSSProperties
  mode?: "redact" | "hide"
}) {
  const { privacyMode } = usePrivacy()

  if (privacyMode && mode === "hide") {
    return null
  }

  return (
    <div
      className={cn(className, privacyMode && "redacted")}
      onClick={onClick}
      style={style}
    >
      {children}
    </div>
  )
}
