import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader";
import GuildSettingsManagers from "@/components/dashboard/GuildSettingsManagers";
import GuildSettingsPrefix from "@/components/dashboard/GuildSettingsPrefix";

export default function DashboardGuildSettingsPage() {
  return (
    <div>
      <DashboardPageHeader
        title="Settings"
        description="Configure the settings for your server here."
      />

      <div className="flex flex-col gap-10">
        <GuildSettingsManagers />

        <GuildSettingsPrefix />
      </div>
    </div>
  );
}
