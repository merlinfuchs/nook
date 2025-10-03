import QueryClientProvider from "@/components/QueryClientProvider";
import { TooltipProvider } from "@/components/ui/tooltip";
import type { Metadata } from "next";
import { ThemeProvider } from "next-themes";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import { Toaster } from "@/components/ui/sonner";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Nook - The friendliest way to start your Discord",
  description:
    "Nook is the easiest way to add a bot to your Discord server. Perfect for beginners, it comes ready to use out of the box, with just the right features to get your community started — and space to grow as you do.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <ThemeProvider
          enableSystem={false}
          defaultTheme="dark"
          attribute="class"
        >
          <TooltipProvider>
            <QueryClientProvider>
              <Toaster />
              {children}
            </QueryClientProvider>
          </TooltipProvider>
        </ThemeProvider>
      </body>
    </html>
  );
}
