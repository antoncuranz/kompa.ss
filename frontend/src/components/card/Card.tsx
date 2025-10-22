import React, {MouseEventHandler} from "react";
import {cn} from "@/lib/utils.ts";

export default function Card({
  title,
  subtitle,
  children,
  headerSlot,
  className,
  onClick
}: {
  title?: string
  subtitle?: string
  children: React.ReactNode
  headerSlot?: React.ReactNode
  className?: string
  onClick?: MouseEventHandler<HTMLDivElement>
}) {

  return (
    <div className={cn("w-full bg-background sm:rounded-3xl sm:shadow-xl shadow-black/10 dark:shadow-white/5", className)} onClick={onClick}>
      <div className="flex flex-col h-full sm:p-2 sm:rounded-3xl sm:border">
        {(title || headerSlot) &&
          <div className="hidden sm:flex flex-row p-3 sm:pb-4 border-b">
            <div className="grow text-xl/[2rem] sm:text-2xl">
              <span className="mr-2">{title}</span>
              <span>{subtitle}</span>
            </div>
            {headerSlot}
          </div>
        }
        <div className="h-full sm:rounded-2xl no-scrollbar overflow-hidden overflow-y-scroll">
          {children}
        </div>
      </div>
    </div>
  )
}