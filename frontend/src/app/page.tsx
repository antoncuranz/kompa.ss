import React from "react";
import {fetchTrips} from "@/requests.ts";
import {Carousel} from "@/components/ui/cards-carousel.tsx";
import Link from "next/link";
import Card from "@/components/card/Card.tsx";
import Navigation from "@/components/navigation/Navigation.tsx";
import NewTripCard from "@/components/card/NewTripCard.tsx";

export default async function Page() {

  const fallbackColors = ["#0081A7", "#00AFB9", "#FED9B7", "#F07167"]

  const cards = (await fetchTrips()).map((trip, idx) => (
      <Link key={trip.id} href={"/" + trip.id + "/itinerary"}>
        <Card className="h-80 w-56 md:h-[40rem] md:w-96">
          <div className="relative h-full w-full">
            {trip.imageUrl &&
              <div className="pointer-events-none absolute inset-x-0 top-0 z-30 h-full bg-gradient-to-b from-black/50 via-transparent to-transparent" />
            }
            <div className="relative z-40 p-8">
              <div className="mt-2 max-w-xs text-left font-sans text-xl font-semibold [text-wrap:balance] text-white md:text-3xl">
                {trip.name}
              </div>
            </div>
            {trip.imageUrl ?
              <img className="absolute inset-0 z-10 object-cover h-full max-w-none w-auto transition duration-300" src={trip.imageUrl} alt=""/>
            :
              <div className="absolute inset-0 w-full h-full" style={{background: fallbackColors[idx % fallbackColors.length]}}/>
            }
          </div>
        </Card>
      </Link>
  ));

  cards.push(<NewTripCard/>)

  return (
    <>
      <Navigation/>
      <main id="root" className="w-full relative z-[1]" style={{height: "calc(100dvh - 4rem)"}}>
        <div className="flex h-full gap-4">
          <div className="w-full h-full py-6">
            <h2 className="max-w-7xl pl-4 mx-auto text-xl md:text-5xl font-bold text-neutral-800 dark:text-neutral-200 font-sans">
              Hello! Let's manage your trips
            </h2>
            <Carousel items={cards} />
          </div>
        </div>
      </main>
    </>
  )
}