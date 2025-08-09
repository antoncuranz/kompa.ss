import React from "react";
import {cn} from "@/lib/utils.ts";
import {GlowContainer} from "@/components/ui/glow-container.tsx";

export default function Card({
  title,
  children,
  headerSlot,
  className,
}: {
  title?: string | null,
  children: React.ReactNode,
  headerSlot?: React.ReactNode
  className?: string
}) {

  return (
    <div className={cn("w-full rounded-2xl sm:rounded-3xl shadow-xl shadow-black/[0.1] dark:shadow-white/[0.05]", className)}>
      <GlowContainer className="flex flex-col h-full sm:p-2 rounded-2xl sm:rounded-3xl">
        {(title || headerSlot) &&
          <div className="flex flex-row p-3 sm:pb-5 border-b">
            <h3 className="flex-grow font-semibold text-xl/[2rem] sm:text-2xl">{title}</h3>
            {headerSlot}
          </div>
        }
        <div className="h-full rounded-2xl no-scrollbar overflow-hidden overflow-y-scroll ">
          {children}
        </div>
      </GlowContainer>
    </div>
  )
}