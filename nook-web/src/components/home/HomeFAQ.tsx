import {
  BadgeQuestionMarkIcon,
  CheckIcon,
  DollarSignIcon,
  PlusIcon,
  ShieldIcon,
  Truck,
} from "lucide-react";

const faq = [
  {
    icon: DollarSignIcon,
    question: "Is Nook free to use?",
    answer:
      "Yes, Nook is free to use. No hidden fees, no paywall, no subscription. Just a great experience.",
  },
  {
    icon: PlusIcon,
    question: "How do I add Nook to my server?",
    answer:
      "You can add Nook to your server by logging in to the Nook dashboard and selecting the server you want to add Nook to.",
  },
  {
    icon: Truck,
    question: "Does Nook have this or that feature?",
    answer:
      "Maybe. Nook is fully modular, so you can pick and choose the features you need. If you don't see the feature you want, you can always ask for it.",
  },
  {
    icon: BadgeQuestionMarkIcon,
    question: "How do I get support?",
    answer:
      "You can get support by joining the Nook Discord server and asking for help in the #support channel.",
  },
  {
    icon: ShieldIcon,
    question: "How is Nook related to Xenon?",
    answer:
      "Nook operates independently from Xenon, but is being built by the same team. So you can expect the same level of quality and support.",
  },
  {
    icon: CheckIcon,
    question: "Why should I choose Nook?",
    answer:
      "Nook is free to use and has all the features you need to start your server. It's built to be easy to use and fully customizable.",
  },
];

export default function HomeFAQ() {
  return (
    <div
      id="faq"
      className="min-h-screen flex items-center justify-center px-6 py-12 xs:py-20"
    >
      <div className="max-w-(--breakpoint-lg)">
        <h2 className="text-3xl xs:text-4xl md:text-5xl leading-[1.15]! font-bold tracking-tight text-center">
          Frequently Asked Questions
        </h2>
        <p className="mt-3 xs:text-lg text-center text-muted-foreground">
          Quick answers to common questions about Nook and its features.
        </p>

        <div className="mt-12 grid md:grid-cols-2 bg-background rounded-xl overflow-hidden outline-[1px] outline-border outline-offset-[-1px]">
          {faq.map(({ question, answer, icon: Icon }) => (
            <div key={question} className="border p-6 -mt-px -ml-px">
              <div className="h-8 w-8 xs:h-10 xs:w-10 flex items-center justify-center rounded-full bg-accent text-accent-foreground">
                <Icon className="h-4 w-4 xs:h-6 xs:w-6" />
              </div>
              <div className="mt-3 mb-2 flex items-start gap-2 text-lg xs:text-[1.35rem] font-semibold tracking-tight">
                <span>{question}</span>
              </div>
              <p className="text-sm xs:text-base">{answer}</p>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
