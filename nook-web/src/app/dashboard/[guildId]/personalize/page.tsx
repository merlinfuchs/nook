import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader";
import GuildPersonalizeColor from "@/components/dashboard/GuildPersonalizeColor";
import GuildPersonalizeProfile from "@/components/dashboard/GuildPersonalizeProfile";

export default function DashboardGuildPersonalizePage() {
  return (
    <div>
      <DashboardPageHeader
        title="Personalize"
        description="Personalize the appearance of the bot in your server here."
      />

      <div className="flex flex-col gap-10">
        <GuildPersonalizeProfile />

        <GuildPersonalizeColor />
      </div>
    </div>
  );
}
