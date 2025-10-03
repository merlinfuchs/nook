import { ModuleCommandOverwriteWire } from "@/lib/types/manage.gen";
import { Card, CardDescription, CardHeader, CardTitle } from "../ui/card";
import { BigSwitch } from "../ui/big-switch";
import { Button } from "../ui/button";
import { PenIcon } from "lucide-react";
import ModuleCommandOverwriteDialog from "./ModuleCommandOverwriteDialog";
import { useCallback } from "react";

export default function ModuleCommandListEntry({
  name,
  description,
  overwrite,
  onOverwriteChange,
}: {
  name: string;
  description: string;
  overwrite?: ModuleCommandOverwriteWire;
  onOverwriteChange: (
    overwrite: ModuleCommandOverwriteWire | undefined
  ) => void;
}) {
  const updateOverwrite = useCallback(
    (ov: Partial<ModuleCommandOverwriteWire>) => {
      onOverwriteChange({
        ...overwrite,
        ...ov,
      });
    },
    [name, overwrite, onOverwriteChange]
  );

  return (
    <Card className="w-full">
      <CardHeader className="flex items-center justify-between">
        <div>
          <CardTitle className="mb-1 text-lg">
            <span className="text-muted-foreground mr-1">/</span>
            {name}
          </CardTitle>
          <CardDescription>{description}</CardDescription>
        </div>
        <div className="flex items-center gap-4">
          <BigSwitch
            checked={!overwrite?.disabled}
            onCheckedChange={(checked) =>
              updateOverwrite({ disabled: !checked })
            }
          />
          {/*<ModuleCommandOverwriteDialog>
            <Button variant="outline" size="icon">
              <PenIcon />
            </Button>
          </ModuleCommandOverwriteDialog>*/}
        </div>
      </CardHeader>
    </Card>
  );
}
