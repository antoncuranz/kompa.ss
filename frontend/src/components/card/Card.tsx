import {Card as InternalCard, CardContent, CardHeader, CardTitle} from "@/components/ui/card.tsx";
import {Separator} from "@/components/ui/separator.tsx";
import React from "react";
import {cn} from "@/lib/utils.ts";

export default async function Card({
  title,
  children,
  headerSlot,
  className
}: {
  title: string,
  children: React.ReactNode,
  headerSlot?: React.ReactNode
  className?: string
}) {

  return (
    <InternalCard className={cn("mb-2 overflow-hidden card h-full", className)}>
      <CardHeader className="pb-0 flex-row justify-between" style={{height: "3.375rem"}}>
        <CardTitle>
          {title}
        </CardTitle>
        {headerSlot}
      </CardHeader>
      <CardContent className="p-0 h-full overflow-y-scroll no-scrollbar">
        <Separator className="mt-4"/>
        {children}
      </CardContent>
    </InternalCard>
  )
}