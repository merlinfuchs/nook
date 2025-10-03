import {
  Blocks,
  Code2Icon,
  DollarSignIcon,
  Settings2Icon,
  UsersIcon,
  WandSparklesIcon,
} from "lucide-react";

const features = [
  {
    icon: Blocks,
    title: "Fully Modular",
    description:
      "Enable or disable modules as you need them. No need to deal with bloated features you don't want.",
  },
  {
    icon: WandSparklesIcon,
    title: "Customizable Appearance",
    description:
      "Customize the appearance of Nook in your server to match your server's theme.",
  },
  {
    icon: Settings2Icon,
    title: "Tweak Everything",
    description:
      "Configure everything to your liking. From the appearance to the behavior, Nook is fully customizable.",
  },
  {
    icon: UsersIcon,
    title: "Collaborative",
    description:
      "Collaborate with your friends to build the perfect server. Nook is allows you to share access.",
  },
  {
    icon: DollarSignIcon,
    title: "Free Forever",
    description:
      "Nook is free forever. No hidden fees, no paywall, no subscription. Just a great experience.",
  },
  {
    icon: Code2Icon,
    title: "Open Source",
    description:
      "Nook is open source. You can view the source code on GitHub and contribute to the project.",
  },
];

export default function HomeFeatures() {
  return (
    <div id="features" className="w-full py-12 xs:py-20 px-6">
      <h2 className="text-3xl xs:text-4xl sm:text-5xl font-bold tracking-tight text-center">
        Unleash Your Creativity
      </h2>
      <div className="w-full max-w-(--breakpoint-lg) mx-auto mt-10 sm:mt-16 grid sm:grid-cols-2 lg:grid-cols-3 gap-6">
        {features.map((feature) => (
          <div
            key={feature.title}
            className="flex flex-col bg-background border rounded-xl py-6 px-5"
          >
            <div className="mb-3 h-10 w-10 flex items-center justify-center bg-muted rounded-full">
              <feature.icon className="h-6 w-6" />
            </div>
            <span className="text-lg font-semibold">{feature.title}</span>
            <p className="mt-1 text-foreground/80 text-[15px]">
              {feature.description}
            </p>
          </div>
        ))}
      </div>
    </div>
  );
}
