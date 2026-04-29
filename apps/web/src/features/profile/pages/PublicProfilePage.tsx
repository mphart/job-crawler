import { useParams } from "react-router-dom";
import { PageShell } from "../../../shared/components/PageShell";
import { ApplicationHistory } from "../components/ApplicationHistory";
import { ProfileHeader } from "../components/ProfileHeader";
import { useProfileQuery } from "../hooks/useProfileQuery";
import { EmptyState } from "../../../shared/components/EmptyState";
import { useSession } from "../../auth/hooks/useSession";

export function PublicProfilePage() {
  const { user } = useSession();
  const { userId = "unknown" } = useParams();
  const { profile, loading, error } = useProfileQuery(userId, user?.token ?? null);

  if (!user) return <PageShell title="Public Profile"><EmptyState title="Not authenticated" description="Please login again." /></PageShell>;
  if (loading || !profile) return <PageShell title="Public Profile"><p>Loading...</p></PageShell>;
  if (error) return <PageShell title="Public Profile"><EmptyState title="Unable to load profile" description={error} /></PageShell>;
  if (profile.isPrivate) return <PageShell title="Public Profile"><EmptyState title="Private profile" description="This user has disabled public visibility." /></PageShell>;

  return (
    <PageShell title={`Profile: ${profile.username}`}>
      <ProfileHeader profile={profile} />
      <ApplicationHistory jobs={profile.appliedJobs} />
    </PageShell>
  );
}
