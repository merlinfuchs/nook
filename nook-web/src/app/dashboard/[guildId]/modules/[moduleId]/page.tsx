"use client";

import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader";
import ModuleCommandList from "@/components/dashboard/ModuleCommandList";
import ModuleConfigEditor from "@/components/dashboard/ModuleConfigEditor";
import { BigSwitch } from "@/components/ui/big-switch";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useModuleConfigureMutation } from "@/lib/api/mutations";
import { useModule } from "@/lib/hooks/api";
import { useGuildId, useModuleId } from "@/lib/hooks/params";
import { ModuleCommandOverwriteWire } from "@/lib/types/manage.gen";
import { ConfigUISchema } from "@/lib/types/module.gen";
import { RJSFValidationError } from "@rjsf/utils";
import { LoaderCircleIcon } from "lucide-react";
import { useCallback, useEffect, useMemo, useState } from "react";
import { toast } from "sonner";

function convertUiSchema(uiSchema: ConfigUISchema): any {
  const res: Record<string, any> = {
    ...Object.fromEntries(
      Object.entries(uiSchema.properties ?? {}).map(([key, value]) => [
        key,
        convertUiSchema(value),
      ])
    ),
    items: uiSchema.items ? convertUiSchema(uiSchema.items) : undefined,
  };
  if (uiSchema["ui:widget"]) {
    res["ui:widget"] = uiSchema["ui:widget"];
  }

  if (uiSchema["ui:channel_types"]) {
    res["ui:channel_types"] = uiSchema["ui:channel_types"];
  }

  if (uiSchema["ui:select_values"]) {
    res["ui:select_values"] = uiSchema["ui:select_values"];
  }

  if (uiSchema["ui:allow_multiple"]) {
    res["ui:allow_multiple"] = uiSchema["ui:allow_multiple"];
  }

  if (uiSchema.layout) {
    res["ui:field"] = "LayoutGridField";

    const children: any[] = [...uiSchema.layout.items];
    for (const child of uiSchema.layout.children ?? []) {
      if (child.type === "container") {
        if (child.header) {
          res[child.header] = {
            ...res[child.header],
            "ui:field": "LayoutHeaderField",
          };
          children.push({
            "ui:row": {
              className: "border rounded-lg p-5 bg-card",
              children: [child.header, ...child.items],
            },
          });
        }
      } else if (child.type === "flex_wrap") {
        children.push({
          "ui:row": {
            className: "flex flex-wrap gap-5",
            children: child.items,
          },
        });
      }
    }

    res["ui:layoutGrid"] = {
      "ui:row": {
        className: "flex flex-col gap-5",
        children: children,
      },
    };
  }

  return res;
}

export default function DashboardGuildModulePage() {
  const guildId = useGuildId();
  const moduleId = useModuleId();
  const mod = useModule();

  const configureMutation = useModuleConfigureMutation(guildId, moduleId);

  const toggleEnabled = useCallback(() => {
    configureMutation.mutate(
      {
        enabled: !mod?.enabled,
        command_overwrites: mod?.command_overwrites ?? {},
        config: mod?.config ?? {},
      },
      {
        onSuccess: (res) => {
          if (res.success) {
            if (res.data.enabled) {
              toast.success("Module has been enabled!");
            } else {
              toast.success("Module has been disabled!");
            }
          } else {
            toast.error(`Failed to enable module: ${res.error.message}`);
          }
        },
      }
    );
  }, [configureMutation, mod?.enabled, mod?.config]);

  const uiSchema = useMemo(() => {
    return mod?.metadata.config_ui_schema;
  }, [mod?.metadata.config_ui_schema]);

  const [tab, setTab] = useState<"settings" | "commands">("settings");

  const [unsavedChanges, setUnsavedChanges] = useState(false);
  const [formData, setFormData] = useState<any>();
  const [commandOverwrites, setCommandOverwrites] = useState<
    Record<string, ModuleCommandOverwriteWire>
  >({});
  const [errors, setErrors] = useState<RJSFValidationError[]>([]);

  const onSubmit = useCallback(() => {
    if (errors.length > 0) {
      toast.error("Please fix the errors before saving.");
      return;
    }

    configureMutation.mutate(
      {
        enabled: !!mod?.enabled,
        command_overwrites: mod?.command_overwrites ?? {},
        config: formData,
      },
      {
        onSuccess: (res) => {
          if (res.success) {
            setUnsavedChanges(false);
            toast.success("Settings have been saved!");
          } else {
            toast.error(`Failed to save settings: ${res.error.message}`);
          }
        },
      }
    );
  }, [configureMutation, mod?.enabled, formData, errors]);

  useEffect(() => {
    if (mod?.config) {
      setFormData(mod.config);
    }
    if (mod?.command_overwrites) {
      setCommandOverwrites(mod.command_overwrites);
    }
  }, [mod?.config, mod?.command_overwrites]);

  const onOverwriteChange = useCallback(
    (name: string, overwrite: ModuleCommandOverwriteWire | undefined) => {
      if (!overwrite) {
        delete commandOverwrites[name];
      } else {
        commandOverwrites[name] = overwrite;
      }
      setCommandOverwrites({ ...commandOverwrites });
      setUnsavedChanges(true);
    },
    [commandOverwrites, setCommandOverwrites]
  );

  if (!mod) {
    return (
      <div className="flex justify-center items-center h-full">
        <LoaderCircleIcon className="size-20 animate-spin" />
      </div>
    );
  }

  return (
    <div>
      <DashboardPageHeader
        title={(mod.metadata.name ?? "Unknown") + " Module"}
        description={
          mod.metadata.description ??
          "Configure the module settings for your server here."
        }
        tabs={[
          {
            label: "Settings",
            value: "settings",
            disabled: mod.metadata.config_schema === undefined,
          },
          {
            label: "Commands",
            value: "commands",
            disabled: !mod.commands?.length,
          },
        ]}
        tab={tab}
        onTabChange={(tab) => setTab(tab as "settings" | "commands")}
      >
        <BigSwitch checked={!!mod.enabled} onCheckedChange={toggleEnabled} />
        <Button
          onClick={onSubmit}
          disabled={errors.length > 0 || !unsavedChanges}
        >
          Save Settings
        </Button>
      </DashboardPageHeader>

      <div className="flex flex-col space-y-5">
        {mod.metadata.config_schema && tab === "settings" && (
          <Card className="p-8 bg-transparent">
            <ModuleConfigEditor
              schema={mod.metadata.config_schema}
              uiSchema={uiSchema ? convertUiSchema(uiSchema) : undefined}
              formData={mod.config}
              onChange={(data, errors) => {
                setFormData(data);
                setErrors(errors);
                setUnsavedChanges(true);
              }}
            />
          </Card>
        )}

        {mod.commands && tab === "commands" && (
          <ModuleCommandList
            commands={mod.commands}
            commandOverwrites={commandOverwrites}
            onOverwriteChange={(name, overwrite) => {
              onOverwriteChange(name, overwrite);
            }}
          />
        )}
      </div>
    </div>
  );
}
