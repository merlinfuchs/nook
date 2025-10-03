import { cn } from "@/lib/utils";
import { CheckIcon } from "lucide-react";

const colorSchemes = [
  {
    name: "Purple",
    value: "purple",
    color: "bg-purple-500!",
  },
  {
    name: "Red",
    value: "red",
    color: "bg-red-500!",
  },
  {
    name: "Blue",
    value: "blue",
    color: "bg-blue-500!",
  },
  {
    name: "Green",
    value: "green",
    color: "bg-green-500!",
  },
  {
    name: "Yellow",
    value: "yellow",
    color: "bg-yellow-400!",
  },
];

export default function GuildColorSchemeSelector({
  value,
  onChange,
}: {
  value: string | null;
  onChange: (value: string) => void;
}) {
  return (
    <div className="flex flex-wrap gap-3">
      {colorSchemes.map((cs) => (
        <button
          key={cs.value}
          className={cn(
            "size-20 rounded-full border flex items-center justify-center hover:scale-105 transition-all",
            cs.color
          )}
          onClick={() => onChange(cs.value)}
        >
          {cs.value === value && <CheckIcon className="size-10" />}
        </button>
      ))}
      {/*<button className="size-20 rounded-full border flex items-center justify-center hover:scale-105 transition-all">
    <PaletteIcon className="size-8" />
  </button>*/}
    </div>
  );
}
