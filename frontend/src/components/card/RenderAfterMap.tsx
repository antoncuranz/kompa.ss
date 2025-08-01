"use client"

import React, {useEffect, useState} from "react";
import {useMap} from "react-map-gl/mapbox";

export default function RenderAfterMap({
  children
}: {
  children: React.ReactNode,
}) {

    const map = useMap()
    const [canRender, setCanRender] = useState(false)

    useEffect(() => {
        map.current?.on('load', () => setCanRender(true))
    }, [map])

    return <>{canRender && children}</>
}
