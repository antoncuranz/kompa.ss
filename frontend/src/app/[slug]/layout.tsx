import Navigation from "@/components/navigation/Navigation.tsx";
import {fetchTrip} from "@/requests.ts";

export async function generateMetadata({
  params
} : {
  params: Promise<{ slug: string }>
}) {
  const tripId = parseInt((await params).slug)
  const trip = await fetchTrip(tripId)

  return {
    title: "kompa.ss - " + trip.name
  };
}

export default async function RootLayout({
  params, children,
}: {
  params: Promise<{ slug: string }>
  children: React.ReactNode
}) {
  const tripId = parseInt((await params).slug)
  const trip = await fetchTrip(tripId)

  return (
      <>
        <Navigation trip={trip}/>
        <main id="root" className="w-full sm:px-4 md:px-6 md:gap-2 relative z-1 h-[calc(100dvh-5.5rem)] sm:h-[calc(100dvh-4.5rem)]">
          {children}
        </main>
      </>
  )
}