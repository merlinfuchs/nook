import GuildList from "@/components/dashboard/GuildList";
import { Separator } from "@/components/ui/separator";

export default function DashboardPage() {
  return (
    <div className="my-20 max-w-5xl mx-auto px-5">
      <div className="text-3xl font-bold mb-2">Servers</div>
      <div className="text-lg text-muted-foreground">
        This list shows all servers you have access to. You can add Nook to your
        server by selecting the server below.
      </div>
      <Separator className="my-8" />
      <GuildList />
    </div>
  );
}
