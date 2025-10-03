import logo from "@/assets/logo.svg";
import { Separator } from "@/components/ui/separator";
import env from "@/lib/env/client";
import { GithubIcon } from "lucide-react";
import Image from "next/image";
import Link from "next/link";

const footerLinks = [
  {
    title: "Features",
    href: "#features",
  },
  {
    title: "FAQ",
    href: "#faq",
  },
  {
    title: "Modules",
    href: "#modules",
  },
  {
    title: "Privacy",
    href: "/privacy",
  },
  {
    title: "Terms",
    href: "/terms",
  },
];

export default function HomeFooter() {
  return (
    <footer className="dark:border-t mt-40 dark bg-background text-foreground">
      <div className="max-w-(--breakpoint-xl) mx-auto">
        <div className="py-12 flex flex-col sm:flex-row items-start justify-between gap-x-8 gap-y-10 px-6 xl:px-0">
          <div>
            {/* Logo */}
            <div className="flex items-center gap-2">
              <Image
                src={logo}
                alt="Nook Logo"
                className="w-10 h-10"
                width={64}
                height={64}
              />
              <div className="text-xl font-bold text-primary">Nook</div>
            </div>

            <ul className="mt-6 flex items-center gap-4 flex-wrap">
              {footerLinks.map(({ title, href }) => (
                <li key={title}>
                  <Link
                    href={href}
                    className="text-muted-foreground hover:text-foreground"
                  >
                    {title}
                  </Link>
                </li>
              ))}
            </ul>
          </div>
        </div>
        <Separator />
        <div className="py-8 flex flex-col-reverse sm:flex-row items-center justify-between gap-x-2 gap-y-5 px-6 xl:px-0">
          {/* Copyright */}
          <span className="text-muted-foreground text-center sm:text-start">
            &copy; {new Date().getFullYear()}{" "}
            <Link href="/" target="_blank">
              Nooks.chat
            </Link>
            . All rights reserved.
          </span>

          <div className="flex items-center gap-5 text-muted-foreground">
            <Link href={env.NEXT_PUBLIC_GITHUB_LINK} target="_blank">
              <GithubIcon className="h-5 w-5" />
            </Link>
          </div>
        </div>
      </div>
    </footer>
  );
}
