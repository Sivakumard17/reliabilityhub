import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import Link from 'next/link'
import { AlertTriangle, LayoutDashboard, Target, Server } from 'lucide-react'
import './globals.css'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'ReliabilityHub',
  description: 'SRE Control Plane for Kubernetes',
}

const navItems = [
  { href: '/',          label: 'Dashboard', icon: LayoutDashboard },
  { href: '/incidents', label: 'Incidents', icon: AlertTriangle },
  { href: '/slos',      label: 'SLOs',      icon: Target },
  { href: '/clusters',  label: 'Clusters',  icon: Server },
]

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <div className="flex h-screen bg-gray-50">
          <aside className="w-60 bg-gray-900 text-white flex flex-col">
            <div className="px-6 py-5 border-b border-gray-700">
              <div className="flex items-center gap-2">
                <div className="w-7 h-7 bg-blue-500 rounded flex items-center justify-center">
                  <span className="text-xs font-bold">R</span>
                </div>
                <div>
                  <p className="text-sm font-semibold">ReliabilityHub</p>
                  <p className="text-xs text-gray-400">SRE Control Plane</p>
                </div>
              </div>
            </div>
            <nav className="flex-1 px-3 py-4 space-y-1">
              {navItems.map(({ href, label, icon: Icon }) => (
                <Link
                  key={href}
                  href={href}
                  className="flex items-center gap-3 px-3 py-2 rounded-lg text-sm text-gray-300 hover:bg-gray-800 hover:text-white transition-colors"
                >
                  <Icon className="h-4 w-4" />
                  {label}
                </Link>
              ))}
            </nav>
            <div className="px-6 py-4 border-t border-gray-700">
              <p className="text-xs text-gray-500">v0.1.0 — dev</p>
            </div>
          </aside>
          <main className="flex-1 overflow-auto">
            {children}
          </main>
        </div>
      </body>
    </html>
  )
}
