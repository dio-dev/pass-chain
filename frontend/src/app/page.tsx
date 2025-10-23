'use client';

import { Button } from '@/components/ui/button';
import { ConnectButton } from '@rainbow-me/rainbowkit';
import { Shield, Lock, Database, Activity } from 'lucide-react';
import { useAccount } from 'wagmi';
import { useRouter } from 'next/navigation';

export default function HomePage() {
  const { isConnected } = useAccount();
  const router = useRouter();

  const handleGetStarted = () => {
    if (isConnected) {
      router.push('/dashboard');
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900">
      {/* Navigation */}
      <nav className="container mx-auto px-4 py-6 flex justify-between items-center">
        <div className="flex items-center space-x-2">
          <Shield className="h-8 w-8 text-purple-400" />
          <span className="text-2xl font-bold text-white">Pass Chain</span>
        </div>
        <ConnectButton />
      </nav>

      {/* Hero Section */}
      <main className="container mx-auto px-4 py-20">
        <div className="text-center max-w-4xl mx-auto">
          <h1 className="text-6xl font-bold text-white mb-6">
            Secure Your Passwords
            <br />
            <span className="text-purple-400">On The Blockchain</span>
          </h1>
          <p className="text-xl text-gray-300 mb-8">
            Revolutionary password management using split-key encryption across
            HashiCorp Vault and Hyperledger Fabric blockchain.
          </p>
          <div className="flex justify-center gap-4">
            <ConnectButton.Custom>
              {({ account, chain, openConnectModal, mounted }) => {
                const ready = mounted;
                const connected = ready && account && chain;

                return (
                  <Button
                    size="lg"
                    onClick={() => {
                      if (connected) {
                        router.push('/dashboard');
                      } else {
                        openConnectModal();
                      }
                    }}
                    className="bg-purple-600 hover:bg-purple-700 text-white px-8 py-6 text-lg"
                  >
                    {connected ? 'Go to Dashboard' : 'Get Started'}
                  </Button>
                );
              }}
            </ConnectButton.Custom>
            <Button
              size="lg"
              variant="outline"
              onClick={() => {
                document.getElementById('features')?.scrollIntoView({ behavior: 'smooth' });
              }}
              className="border-purple-400 text-purple-400 hover:bg-purple-950 px-8 py-6 text-lg"
            >
              Learn More
            </Button>
          </div>
        </div>

        {/* Features */}
        <div id="features" className="grid grid-cols-1 md:grid-cols-3 gap-8 mt-20">
          <FeatureCard
            icon={<Lock className="h-12 w-12 text-purple-400" />}
            title="Split-Key Security"
            description="Your credentials are split between Vault and blockchain. Neither can decrypt alone."
          />
          <FeatureCard
            icon={<Database className="h-12 w-12 text-purple-400" />}
            title="Immutable Audit Trail"
            description="Every access is logged on blockchain. Complete transparency and accountability."
          />
          <FeatureCard
            icon={<Activity className="h-12 w-12 text-purple-400" />}
            title="Pay Per Use"
            description="Pay once to store, small fee per access. No subscriptions, complete control."
          />
        </div>
      </main>
    </div>
  );
}

function FeatureCard({
  icon,
  title,
  description,
}: {
  icon: React.ReactNode;
  title: string;
  description: string;
}) {
  return (
    <div className="bg-slate-800/50 backdrop-blur-sm border border-purple-500/20 rounded-lg p-6 text-center">
      <div className="flex justify-center mb-4">{icon}</div>
      <h3 className="text-xl font-bold text-white mb-2">{title}</h3>
      <p className="text-gray-400">{description}</p>
    </div>
  );
}




