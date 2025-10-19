"use client";

import {JazzReactProvider} from "jazz-tools/react";
import {JazzFestAccount} from "@/schema";
import {Auth} from "@/components/Auth.tsx";

export function JazzProvider({children}: {
  children: React.ReactNode
}) {
  return (
    <JazzReactProvider
      sync={{
        peer: "ws://127.0.0.1:4200"
      }}
      guestMode={false}
      AccountSchema={JazzFestAccount}
    >
      <Auth/>
      {children}
    </JazzReactProvider>
  );
}