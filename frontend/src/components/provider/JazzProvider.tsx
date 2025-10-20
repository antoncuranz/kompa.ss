"use client";

import {JazzReactProvider} from "jazz-tools/react";
import {JazzFestAccount} from "@/schema";
import {Auth} from "@/components/Auth.tsx";
import {JazzInspector} from "jazz-tools/inspector";

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
      <JazzInspector/>
      <Auth/>
      {children}
    </JazzReactProvider>
  );
}