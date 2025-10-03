import { SerializedMessageFormat } from "@/lib/types/module.gen";
import { Input } from "./ui/input";
import { useCallback } from "react";
import { Textarea } from "./ui/textarea";

export default function MessageFormatEditor({
  data,
  onChange,
}: {
  data: SerializedMessageFormat;
  onChange: (data: SerializedMessageFormat) => void;
}) {
  const updateData = useCallback(
    (change: Partial<SerializedMessageFormat>) => {
      onChange({ ...data, ...change });
    },
    [data, onChange]
  );

  return (
    <div className="border rounded-[8px] border-l-6 border-l-red-500 p-5 flex items-start gap-5">
      <div className="flex flex-col gap-3 flex-auto">
        <Input
          value={data.title ?? ""}
          onChange={(e) => updateData({ title: e.target.value })}
          className="rounded-[10px] font-bold text-foreground"
        />
        <Textarea
          value={data.description ?? ""}
          onChange={(e) => updateData({ description: e.target.value })}
          className="resize-none h-32 rounded-[10px] text-foreground/90"
        />
      </div>
      <div className="flex-none">
        <img
          src="https://cdn.discordapp.com/embed/avatars/0.png"
          className="size-24 rounded-[8px]"
          alt="User Avatar"
        />
      </div>
    </div>
  );
}
