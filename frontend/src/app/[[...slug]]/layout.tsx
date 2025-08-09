import Navigation from "@/components/navigation/Navigation.tsx";

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {

  return (
    <>
      <Navigation/>
      <main id="root" className="w-full p-2 pt-0 sm:px-6 md:gap-2 relative z-[1]" style={{height: "calc(100dvh - 4rem)"}}>
        {children}
      </main>
    </>
  )
}