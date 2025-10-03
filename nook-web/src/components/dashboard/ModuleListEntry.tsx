import { useGuildId } from "@/lib/hooks/params";
import { ModuleWire } from "@/lib/types/manage.gen";
import { CheckCircle2Icon } from "lucide-react";
import { DynamicIcon } from "lucide-react/dynamic";
import Link from "next/link";
import { Button } from "../ui/button";
import {
  Card,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "../ui/card";

export default function ModuleListEntry(module: ModuleWire) {
  const guildId = useGuildId();

  return (
    <Card className="flex flex-col p-0">
      <CardHeader className="flex flex-row gap-4 p-5 pb-0 flex-auto">
        <div className="h-10 w-10 bg-primary/40 flex-none rounded-md flex items-center justify-center">
          <DynamicIcon
            name={module.metadata.icon as any}
            className="w-6 h-6 text-primary"
          />
        </div>
        <div className="flex-auto">
          <div className="flex items-start justify-between">
            <CardTitle className="mb-2 text-2xl">
              {module.metadata.name}
            </CardTitle>

            {module.enabled && (
              <div>
                <CheckCircle2Icon className="-mt-1 w-6 h-6 text-primary" />
              </div>
            )}
          </div>
          <CardDescription className="text-base">
            {module.metadata.description}
          </CardDescription>
        </div>
      </CardHeader>
      <CardFooter className="p-5 pt-1">
        <Button variant="outline" asChild>
          <Link href={`/dashboard/${guildId}/modules/${module.id}`}>
            Configure
          </Link>
        </Button>
      </CardFooter>
    </Card>
  );
}
