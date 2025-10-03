import { ModuleCommandOverwriteWire } from "@/lib/types/manage.gen";
import ModuleCommandListEntry from "./ModuleCommandListEntry";

export default function ModuleCommandList({
  commands,
  commandOverwrites,
  onOverwriteChange,
}: {
  commands: { name: string; description: string }[];
  commandOverwrites: Record<string, ModuleCommandOverwriteWire>;
  onOverwriteChange: (
    name: string,
    overwrite: ModuleCommandOverwriteWire | undefined
  ) => void;
}) {
  return (
    <div className="flex flex-col gap-5">
      {commands.map((cmd) => (
        <ModuleCommandListEntry
          {...cmd}
          key={cmd.name}
          overwrite={commandOverwrites[cmd.name]}
          onOverwriteChange={(overwrite) => {
            onOverwriteChange(cmd.name, overwrite);
          }}
        />
      ))}
    </div>
  );
}
