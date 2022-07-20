/** @type {import('next').NextConfig} */
const nextConfig = {
  eslint: {
    ignoreDuringBuilds: true,
  },
  typescript: {
    ignoreBuildErrors: true,
  },
  // async redirects() {
  //   return [
  //     {
  //       source: "/api/:path*",
  //       destination: "http://localhost:3000/:path*",
  //       permanent: false,
  //     },
  //   ];
  // },
};

module.exports = nextConfig;
