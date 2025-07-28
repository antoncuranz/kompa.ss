"use client"

import React, {CSSProperties, MouseEventHandler} from "react";
import {useStore} from "@/store.ts";
import {cn} from "@/lib/utils.ts";

export default function PrivacyFilter({
  children,
  className,
  onClick,
  style
}: {
  children: React.ReactNode,
  className?: string,
  onClick?: MouseEventHandler<HTMLDivElement>,
  style?: CSSProperties
}) {
  const { privacyMode } = useStore();

  return (
    <div className={cn(className, privacyMode && "redacted")} onClick={onClick} style={style}>
      {children}
    </div>
  )
}
