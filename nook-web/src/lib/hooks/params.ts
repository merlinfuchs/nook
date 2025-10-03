import { useParams } from "next/navigation";

export function useGuildId() {
  const params = useParams();
  return params.guildId as string;
}

export function useModuleId() {
  const params = useParams();
  return params.moduleId as string;
}
