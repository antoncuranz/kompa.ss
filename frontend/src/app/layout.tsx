import type {Metadata, Viewport} from 'next'
import '../index.css'
import {Toaster} from "@/components/ui/sonner.tsx";
import {ThemeProvider} from "@/components/provider/ThemeProvider.tsx";

export const metadata: Metadata = {
  title: "kompa.ss",
  appleWebApp: {
    title: "kompa.ss",
  }
}

export const viewport: Viewport = {
  width: "device-width",
  initialScale: 1.0,
  maximumScale: 1.0,
  themeColor: "white",
}

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="de" suppressHydrationWarning>
    <body>
      <div className="root">
        <ThemeProvider
          attribute="class"
          defaultTheme="system"
          enableSystem
          disableTransitionOnChange
        >
          {children}
        </ThemeProvider>
        <Toaster/>
      </div>
      <script>{`
        function updateThemeColor() {
          const meta = document.querySelector("meta[name=theme-color]")
          if (!meta) return
          const bg = document.querySelector("html").classList.contains("dark") ? "black" : "white"
          meta.setAttribute("content", bg)
        }
        updateThemeColor()
        const observer = new MutationObserver(mutations => {
          for (const m of mutations) {
            if (m.type === "attributes" && m.attributeName === "class") {
              updateThemeColor()
            }
          }
        })
        observer.observe(document.documentElement, { attributes: true })
      `}</script>
    </body>
    </html>
  )
}