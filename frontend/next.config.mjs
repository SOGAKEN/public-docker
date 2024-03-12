/** @type {import('next').NextConfig} */
const nextConfig = {
	reactStrictMode: true,
	output: 'standalone',
	async rewrites() {
		return [
			{
				source: '/api/:path*',
				destination: 'http://localhost:8081/api/:path*',
			},
		];
	},
};

export default nextConfig;
