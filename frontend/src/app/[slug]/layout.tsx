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

  return (
      <>
        <Navigation tripId={tripId}/>
        <main id="root" className="w-full p-2 pt-0 sm:px-6 md:gap-2 relative z-1" style={{height: "calc(100dvh - 4rem)"}}>
          {children}
        </main>
      </>
  )
}