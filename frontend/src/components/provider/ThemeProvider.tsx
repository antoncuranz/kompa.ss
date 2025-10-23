"use client"

import * as React from "react"
import {ThemeProvider as NextThemesProvider, useTheme} from "next-themes"
import {useEffect} from "react";

export function ThemeProvider({
  children,
  ...props
}: React.ComponentProps<typeof NextThemesProvider>) {
  return (
    <NextThemesProvider {...props}>
      <ThemeColorUpdater/>
      {children}
    </NextThemesProvider>
  )
}

function ThemeColorUpdater() {
  const { theme } = useTheme()

  useEffect(() => {
    const meta = document.querySelector("meta[name=theme-color]")
    if (!meta) return

    const themeColor = theme == "dark" ? "black" : "white"
    meta.setAttribute("content", themeColor)
  }, [theme])

  return null
}