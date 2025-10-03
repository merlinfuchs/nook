import DashboardPageHeader from "@/components/dashboard/DashboardPageHeader";
import ModuleList from "@/components/dashboard/ModuleList";

export default function DashboardGuildPage() {
  return (
    <div>
      <DashboardPageHeader
        title="Modules"
        description="Configure the modules for your server here."
      />
      <ModuleList />
    </div>
  );
}
