"use client"

import React, {useEffect, useState} from "react";
import {useMap} from "react-map-gl/maplibre";

export default function RenderAfterMap({
  children
}: {
  theme: string,
  children: React.ReactNode,
}) {
  const map = useMap()
  const [canRender, setCanRender] = useState(false)

  // useEffect(() => {
  //   map.current?.setConfigProperty("basemap", "lightPreset", theme);
  // }, [map, theme])

  useEffect(() => {
    map.current?.on('load', () => setCanRender(true))
  }, [map])

  return <>{canRender && children}</>
}
