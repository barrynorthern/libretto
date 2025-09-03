import * as React from 'react'

export function Card({ className = '', ...props }: React.HTMLAttributes<HTMLDivElement>) {
  return <div className={`rounded border border-zinc-700 bg-zinc-800 ${className}`} {...props} />
}
export function CardHeader({ className = '', ...props }: React.HTMLAttributes<HTMLDivElement>) {
  return <div className={`border-b border-zinc-700 p-3 ${className}`} {...props} />
}
export function CardTitle({ className = '', ...props }: React.HTMLAttributes<HTMLDivElement>) {
  return <div className={`text-lg font-medium ${className}`} {...props} />
}
export function CardContent({ className = '', ...props }: React.HTMLAttributes<HTMLDivElement>) {
  return <div className={`p-3 text-sm text-zinc-300 ${className}`} {...props} />
}

