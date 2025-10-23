import type { MetadataRoute } from 'next'

export default function manifest(): MetadataRoute.Manifest {
  return {
    name: "kompa.ss travel planner",
    short_name: "kompa.ss",
    description: "kompa.ss travel planner",
    start_url: "/",
    display: "standalone",
    background_color: "#ffffff",
    theme_color: "chocolate",
    orientation: "portrait",
    icons: [
      {
        src: "/icon.png",
        sizes: "160x160",
        type: "image/png",
      },
    ],
  }
}