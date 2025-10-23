'use client';

import { Button } from '@/components/ui/button';
import { ConnectButton } from '@rainbow-me/rainbowkit';
import { Shield, Lock, Database, Activity, Key, Wallet, ChevronRight, Github, Twitter } from 'lucide-react';
import { useAccount } from 'wagmi';
import { useRouter } from 'next/navigation';
import { useState } from 'react';

export default function HomePage() {
  const { isConnected } = useAccount();
  const router = useRouter();
  const [activeFeature, setActiveFeature] = useState(0);

  const handleGetStarted = () => {
    if (isConnected) {
      router.push('/dashboard');
    }
  };

  const features = [
    {
      icon: Lock,
      title: "True Zero-Knowledge",
      description: "Your passwords are encrypted in your browser. Our servers never see plaintext. Even we can't decrypt them.",
      color: "from-purple-500 to-pink-500"
    },
    {
      icon: Key,
      title: "Split-Key Security",
      description: "Your encryption key is split into 3 parts. Vault, blockchain, and your device. All 3 needed to decrypt.",
      color: "from-blue-500 to-cyan-500"
    },
    {
      icon: Wallet,
      title: "Wallet Authentication",
      description: "Use MetaMask or any Web3 wallet. No passwords to remember. Your crypto wallet is your identity.",
      color: "from-green-500 to-emerald-500"
    },
    {
      icon: Activity,
      title: "Blockchain Audit Trail",
      description: "Every access is logged on Hyperledger Fabric. Immutable, tamper-proof audit logs.",
      color: "from-orange-500 to-red-500"
    },
    {
      icon: Shield,
      title: "Enterprise Ready",
      description: "Kubernetes deployment, high availability, SOC2 compliant. Built for production from day one.",
      color: "from-violet-500 to-purple-500"
    },
    {
      icon: Database,
      title: "Open Source",
      description: "Fully open source. Audit the code yourself. No backdoors, no secrets, complete transparency.",
      color: "from-indigo-500 to-blue-500"
    }
  ];

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900">
      {/* Navigation */}
      <nav className="fixed w-full z-50 bg-slate-900/80 backdrop-blur-md border-b border-purple-500/20">
        <div className="container mx-auto px-4 py-4 flex justify-between items-center">
          <div className="flex items-center space-x-2">
            <Shield className="h-8 w-8 text-purple-400" />
            <span className="text-2xl font-bold text-white">Pass Chain</span>
          </div>
          <div className="hidden md:flex items-center space-x-6">
            <a href="#features" className="text-gray-300 hover:text-purple-400 transition">Features</a>
            <a href="#how-it-works" className="text-gray-300 hover:text-purple-400 transition">How It Works</a>
            <a href="#security" className="text-gray-300 hover:text-purple-400 transition">Security</a>
            <a href="/docs" className="text-gray-300 hover:text-purple-400 transition">Docs</a>
            <ConnectButton />
          </div>
        </div>
      </nav>

      {/* Hero Section */}
      <section className="pt-32 pb-20 px-4">
        <div className="container mx-auto text-center">
          <div className="mb-8 animate-float">
            <Shield className="h-32 w-32 mx-auto text-purple-400 drop-shadow-[0_0_30px_rgba(168,85,247,0.5)]" />
          </div>
          <h1 className="text-6xl md:text-8xl font-bold mb-6 bg-gradient-to-r from-purple-400 via-pink-400 to-purple-400 text-transparent bg-clip-text animate-pulse">
            Password Management
            <br />
            Reimagined
          </h1>
          <p className="text-xl md:text-2xl text-gray-300 mb-8 max-w-3xl mx-auto">
            The first password manager where <span className="text-purple-400 font-bold">nobody</span> can decrypt your passwords. 
            Not us. Not hackers. Not even the NSA.
          </p>
          <div className="flex flex-col md:flex-row gap-4 justify-center">
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
                    className="bg-purple-600 hover:bg-purple-700 text-white px-8 py-6 text-lg transform hover:scale-105 transition"
                  >
                    {connected ? (
                      <>Go to Dashboard <ChevronRight className="ml-2" /></>
                    ) : (
                      '🚀 Get Started'
                    )}
                  </Button>
                );
              }}
            </ConnectButton.Custom>
            <Button
              size="lg"
              variant="outline"
              onClick={() => document.getElementById('features')?.scrollIntoView({ behavior: 'smooth' })}
              className="border-purple-400 text-purple-400 hover:bg-purple-950 px-8 py-6 text-lg"
            >
              Learn More
            </Button>
          </div>
        </div>
      </section>

      {/* Features */}
      <section id="features" className="py-20 px-4 bg-slate-900/50">
        <div className="container mx-auto">
          <h2 className="text-5xl font-bold text-center mb-16 text-white">Why Pass Chain?</h2>
          <div className="grid md:grid-cols-3 gap-8">
            {features.map((feature, idx) => (
              <div
                key={idx}
                className="bg-slate-800/50 backdrop-blur-sm border border-purple-500/20 rounded-xl p-8 hover:border-purple-500/40 transition transform hover:scale-105 cursor-pointer"
                onMouseEnter={() => setActiveFeature(idx)}
              >
                <div className={`bg-gradient-to-r ${feature.color} p-4 rounded-lg inline-block mb-4`}>
                  <feature.icon className="h-8 w-8 text-white" />
                </div>
                <h3 className="text-2xl font-bold mb-4 text-white">{feature.title}</h3>
                <p className="text-gray-300">{feature.description}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* How It Works */}
      <section id="how-it-works" className="py-20 px-4">
        <div className="container mx-auto">
          <h2 className="text-5xl font-bold text-center mb-16 text-white">How It Works</h2>
          <div className="max-w-4xl mx-auto space-y-8">
            {[
              {
                step: 1,
                title: "Connect Your Wallet",
                description: "Use MetaMask, WalletConnect, or any Web3 wallet. Your wallet address is your identity."
              },
              {
                step: 2,
                title: "Encrypt in Your Browser",
                description: "Your password is encrypted using military-grade XChaCha20-Poly1305 encryption. This happens entirely in your browser - we never see your plaintext password."
              },
              {
                step: 3,
                title: "Split the Key",
                description: "The encryption key is split into 3 parts using Shamir Secret Sharing: Part 1 → HashiCorp Vault, Part 2 → Hyperledger Fabric, Part 3 → Your device (backup)"
              },
              {
                step: 4,
                title: "Retrieve Securely",
                description: "To decrypt, you need: Your wallet signature + any 2 of 3 key parts. Even if Vault OR blockchain is compromised, your passwords stay safe."
              }
            ].map((item) => (
              <div key={item.step} className="flex gap-6 items-start">
                <div className="bg-purple-600 text-white rounded-full w-12 h-12 flex items-center justify-center font-bold text-xl flex-shrink-0">
                  {item.step}
                </div>
                <div>
                  <h3 className="text-2xl font-bold mb-2 text-white">{item.title}</h3>
                  <p className="text-gray-300">{item.description}</p>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Security */}
      <section id="security" className="py-20 px-4 bg-slate-900/50">
        <div className="container mx-auto">
          <h2 className="text-5xl font-bold text-center mb-16 text-white">Security First</h2>
          <div className="grid md:grid-cols-2 gap-12 max-w-5xl mx-auto">
            <div>
              <h3 className="text-3xl font-bold mb-6 text-red-400">What We Can't Do</h3>
              <ul className="space-y-4 text-lg">
                {[
                  "Decrypt your passwords",
                  "Access your credentials",
                  "See your plaintext data",
                  "Recover your passwords without your wallet"
                ].map((item, idx) => (
                  <li key={idx} className="flex items-start gap-3">
                    <span className="text-red-400 text-2xl">❌</span>
                    <span className="text-gray-300">{item}</span>
                  </li>
                ))}
              </ul>
            </div>
            <div>
              <h3 className="text-3xl font-bold mb-6 text-green-400">What You Get</h3>
              <ul className="space-y-4 text-lg">
                {[
                  "Client-side encryption",
                  "Split-key architecture",
                  "Blockchain audit trail",
                  "Complete control over your data"
                ].map((item, idx) => (
                  <li key={idx} className="flex items-start gap-3">
                    <span className="text-green-400 text-2xl">✅</span>
                    <span className="text-gray-300">{item}</span>
                  </li>
                ))}
              </ul>
            </div>
          </div>
        </div>
      </section>

      {/* Tech Stack */}
      <section className="py-20 px-4">
        <div className="container mx-auto text-center">
          <h2 className="text-5xl font-bold mb-16 text-white">Built With Modern Tech</h2>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-6 max-w-4xl mx-auto">
            {[
              { emoji: "⚛️", name: "React/Next.js" },
              { emoji: "🔵", name: "Go Lang" },
              { emoji: "🔐", name: "HashiCorp Vault" },
              { emoji: "⛓️", name: "Hyperledger Fabric" },
              { emoji: "☸️", name: "Kubernetes" },
              { emoji: "🐘", name: "PostgreSQL" },
              { emoji: "🌐", name: "Web3.js" },
              { emoji: "🐳", name: "Docker" }
            ].map((tech, idx) => (
              <div key={idx} className="bg-slate-800/50 p-6 rounded-xl border border-purple-500/20 hover:border-purple-500/40 transition">
                <div className="text-4xl mb-2">{tech.emoji}</div>
                <div className="font-bold text-white">{tech.name}</div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* CTA */}
      <section className="py-20 px-4 bg-gradient-to-r from-purple-900 to-pink-900">
        <div className="container mx-auto text-center">
          <h2 className="text-5xl font-bold mb-6 text-white">Ready to Secure Your Passwords?</h2>
          <p className="text-xl text-gray-200 mb-8">
            Join the future of password management. No credit card required.
          </p>
          <div className="flex flex-col md:flex-row gap-4 justify-center">
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
                    className="bg-white text-purple-900 hover:bg-gray-100 px-8 py-6 text-lg font-bold transform hover:scale-105 transition"
                  >
                    {connected ? '📊 Go to Dashboard' : '🚀 Launch App'}
                  </Button>
                );
              }}
            </ConnectButton.Custom>
            <Button
              size="lg"
              variant="outline"
              onClick={() => window.open('https://github.com/yourusername/pass-chain', '_blank')}
              className="border-white text-white hover:bg-white/10 px-8 py-6 text-lg"
            >
              <Github className="mr-2 h-5 w-5" />
              Star on GitHub
            </Button>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="py-12 px-4 bg-slate-900 border-t border-purple-500/20">
        <div className="container mx-auto text-center">
          <div className="flex justify-center space-x-8 mb-6">
            <a href="/docs" className="text-gray-400 hover:text-purple-400 transition">Documentation</a>
            <a href="https://github.com/yourusername/pass-chain" className="text-gray-400 hover:text-purple-400 transition flex items-center gap-2">
              <Github className="h-4 w-4" />
              GitHub
            </a>
            <a href="#" className="text-gray-400 hover:text-purple-400 transition flex items-center gap-2">
              <Twitter className="h-4 w-4" />
              Twitter
            </a>
          </div>
          <p className="text-gray-400">
            © 2025 Pass Chain. Open Source. MIT License.
          </p>
          <p className="text-gray-500 mt-2 font-bold">
            AUUUUFFFF! 🔥
          </p>
        </div>
      </footer>

      <style jsx>{`
        @keyframes float {
          0%, 100% { transform: translateY(0px); }
          50% { transform: translateY(-20px); }
        }
        .animate-float {
          animation: float 3s ease-in-out infinite;
        }
      `}</style>
    </div>
  );
}
