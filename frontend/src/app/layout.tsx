import type {Metadata, Viewport} from 'next'
import '../index.css'
import {Toaster} from "@/components/ui/toaster.tsx";
import {ThemeProvider} from "@/components/provider/ThemeProvider.tsx";
import Navigation from "@/components/navigation/Navigation.tsx";
import {UserProvider} from "@/components/provider/UserProvider.tsx";
import {getCurrentUser} from "@/requests.ts";

export const metadata: Metadata = {
  title: "travel-planner",
}

export const viewport: Viewport = {
  width: "device-width",
  initialScale: 1.0,
  maximumScale: 1.0
}

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  const username = await getCurrentUser()

  return (
    <html lang="en" suppressHydrationWarning>
    <body>
      <UserProvider username={username}>
        <ThemeProvider
          attribute="class"
          defaultTheme="system"
          enableSystem
          disableTransitionOnChange
        >
          <Navigation/>
          <main id="root" className="w-full p-4 pt-0 sm:px-6 md:gap-2" style={{height: "calc(100dvh - 4rem)"}}>
            {children}
          </main>
          <Toaster/>
        </ThemeProvider>
      </UserProvider>
    </body>
    </html>
  )
}