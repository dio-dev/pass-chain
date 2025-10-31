/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  swcMinify: true,
  output: 'standalone',
  eslint: {
    // Disable ESLint during builds (warnings prevent build)
    ignoreDuringBuilds: true,
  },
  env: {
    NEXT_PUBLIC_API_URL: 'http://localhost:8080',
  },
  webpack: (config) => {
    config.resolve.fallback = { 
      fs: false, 
      net: false, 
      tls: false,
      // Fix MetaMask SDK issues
      '@react-native-async-storage/async-storage': false,
      'react-native': false,
    };
    config.externals.push('pino-pretty', 'lokijs', 'encoding');
    
    // Ignore MetaMask SDK warnings
    config.ignoreWarnings = [
      { module: /node_modules\/@metamask\/sdk/ },
      { file: /node_modules\/@metamask\/sdk/ },
    ];
    
    return config;
  },
};

module.exports = nextConfig;
