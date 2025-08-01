/** @type {import('next').NextConfig} */
const nextConfig = {
    output: "standalone",
    rewrites: async () => [
        {
            source: "/api/:path*",
            destination: process.env.BACKEND_URL + "/api/:path*",
        },
    ]
}

export default nextConfig