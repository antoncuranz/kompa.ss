/** @type {import('next').NextConfig} */
const nextConfig = {
    output: "standalone",
    rewrites: async () => [
        {
            source: "/api/:path*",
            destination: process.env.BACKEND_URL + "/api/:path*",
        },
    ],
    experimental: {
        swcPlugins: [["superjson-next", {router: "APP"}]],
    }
}

export default nextConfig