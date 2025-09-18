import type {Metadata, Viewport} from 'next'
import '../index.css'
import {Toaster} from "@/components/ui/toaster.tsx";
import {ThemeProvider} from "@/components/provider/ThemeProvider.tsx";
import {UserProvider} from "@/components/provider/UserProvider.tsx";
import {getCurrentUser} from "@/requests.ts";

export const metadata: Metadata = {
  title: "kompa.ss",
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
          {children}
          <Toaster/>
        </ThemeProvider>
      </UserProvider>
    </body>
    </html>
  )
}