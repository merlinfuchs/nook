"use client";

import { useGuildSettingsUpdateMutation } from "@/lib/api/mutations";
import { useGuildSettings } from "@/lib/hooks/api";
import { useGuildId } from "@/lib/hooks/params";
import { DialogDescription } from "@radix-ui/react-dialog";
import { ReactNode, useCallback, useEffect, useMemo, useState } from "react";
import { toast } from "sonner";
import { Button } from "../ui/button";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "../ui/dialog";
import { Input } from "../ui/input";
import ColorSchemeSelector from "./GuildColorSchemeSelector";

export default function QuickstartDialog({
  children,
}: {
  children: ReactNode;
}) {
  const defaultOpen = window.location.search.includes("quickstart");
  const [open, setOpen] = useState<boolean>(defaultOpen);

  const settings = useGuildSettings();
  const updateMutation = useGuildSettingsUpdateMutation(useGuildId());

  const [step, setStep] = useState<number>(0);
  const [colorScheme, setColorScheme] = useState<string | null>(null);
  const [commandPrefix, setCommandPrefix] = useState<string | null>(null);

  useEffect(() => {
    if (settings) {
      setColorScheme(settings.color_scheme ?? settings.default.color_scheme);
      setCommandPrefix(settings.command_prefix);
    }
  }, [settings]);

  const onSubmit = useCallback(() => {
    updateMutation.mutate(
      {
        command_prefix: commandPrefix ?? undefined,
        color_scheme: colorScheme ?? undefined,
      },
      {
        onSuccess: (res) => {
          if (res.success) {
            toast.success("Settings have been saved!");
            setStep(0);
            setOpen(false);
          } else {
            toast.error(`Failed to save settings: ${res.error.message}`);
          }
        },
      }
    );
  }, [updateMutation, commandPrefix, colorScheme]);

  const steps = useMemo(() => {
    return [
      <QuickstartWeblcome key="welcome" />,
      <QuickstartColorScheme
        key="color-scheme"
        value={colorScheme}
        onChange={setColorScheme}
      />,
      <QuickstartCommandPrefix
        key="command-prefix"
        placeholder={settings?.default.command_prefix ?? ""}
        value={commandPrefix}
        onChange={setCommandPrefix}
      />,
    ];
  }, [colorScheme, commandPrefix, settings]);

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>{children}</DialogTrigger>
      <DialogContent>
        {steps[step]}

        <DialogFooter className="flex items-center gap-3 justify-end">
          {step === 0 ? (
            <Button variant="outline" onClick={() => setOpen(false)}>
              Cancel
            </Button>
          ) : (
            <Button variant="outline" onClick={() => setStep(step - 1)}>
              Previous Step
            </Button>
          )}
          {step < steps.length - 1 ? (
            <Button onClick={() => setStep(step + 1)}>Next Step</Button>
          ) : (
            <Button onClick={onSubmit}>Save Settings</Button>
          )}
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

function QuickstartWeblcome() {
  return (
    <>
      <DialogHeader>
        <DialogTitle>Welcome to Nook!</DialogTitle>
        <DialogDescription>
          Nook is the friendliest way to level up your Discord server.
        </DialogDescription>

        <div className="text-muted-foreground my-3">
          This is a quickstart guide will help you configure your server to get
          the most out of Nook. <br />
          You can always change these settings later in the settings menu.
        </div>
      </DialogHeader>
    </>
  );
}

function QuickstartColorScheme({
  value,
  onChange,
}: {
  value: string | null;
  onChange: (value: string) => void;
}) {
  return (
    <>
      <DialogHeader>
        <DialogTitle>Color Scheme</DialogTitle>
        <DialogDescription>
          Chose a color scheme. This is used for response messages across all
          modules.
        </DialogDescription>
      </DialogHeader>

      <div className="my-5">
        <ColorSchemeSelector value={value} onChange={onChange} />
      </div>
    </>
  );
}

function QuickstartCommandPrefix({
  placeholder,
  value,
  onChange,
}: {
  placeholder: string;
  value: string | null;
  onChange: (value: string) => void;
}) {
  return (
    <>
      <DialogHeader>
        <DialogTitle>Command Prefix</DialogTitle>
        <DialogDescription>
          Chose a command prefix. This is used for all commands.
        </DialogDescription>
      </DialogHeader>

      <Input
        value={value ?? ""}
        placeholder={placeholder}
        onChange={(e) => onChange(e.target.value)}
      />
    </>
  );
}
