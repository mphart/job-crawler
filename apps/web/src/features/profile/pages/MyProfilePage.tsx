import React from "react";
import { useNavigate } from "react-router-dom";
import { PageShell } from "../../../shared/components/PageShell";
import { Button } from "../../../shared/components/Button";
import { NotificationSettingsPanel } from "../../notifications/components/NotificationSettingsPanel";
import { ProfileHeader } from "../components/ProfileHeader";
import { useProfileQuery } from "../hooks/useProfileQuery";
import { ResumeUploader } from "../components/ResumeUploader";
import { EmailOptInToggle } from "../components/EmailOptInToggle";
import { PrivacyToggle } from "../components/PrivacyToggle";
import { updateProfile } from "../api/profile.api";
import { ApplicationHistory } from "../components/ApplicationHistory";
import { EmptyState } from "../../../shared/components/EmptyState";
import { useSession } from "../../auth/hooks/useSession";
import { PreferredCompaniesPicker } from "../components/PreferredCompaniesPicker";

export function MyProfilePage() {
  const navigate = useNavigate();
  const { user } = useSession();
  const { profile, loading, error, setProfile } = useProfileQuery(user?.id ?? "", user?.token ?? null);

  if (!user) {
    return <PageShell title="My Profile"><EmptyState title="Not authenticated" description="Please login again." /></PageShell>;
  }

  if (error) {
    return <PageShell title="My Profile"><EmptyState title="Unable to load profile" description={error} /></PageShell>;
  }

  if (loading || !profile) {
    return <PageShell title="My Profile"><p>Loading profile...</p></PageShell>;
  }

  return (
    <PageShell title="My Profile" actions={<Button variant="secondary" onClick={() => navigate("/feed")}>Homepage</Button>}>
      <ProfileHeader profile={profile} />
      <div style={{ display: "flex", gap: "1rem", marginBottom: "1rem", flexWrap: "wrap" }}>
        <EmailOptInToggle value={profile.preferences.emailOptIn} onChange={async (value) => setProfile(await updateProfile({ preferences: { ...profile.preferences, emailOptIn: value } }, user.token))} />
        <PrivacyToggle value={profile.isPrivate} onChange={async (isPrivate) => setProfile(await updateProfile({ isPrivate }, user.token))} />
      </div>
      <ResumeUploader
        current={profile.resumeFileName}
        onUpload={async ({ fileName, contentBase64 }) => {
          const updated = await updateProfile({ resumeFileName: fileName, resumeContentBase64: contentBase64 }, user.token);
          setProfile(updated);
        }}
      />
      <PreferredCompaniesPicker
        token={user.token}
        selected={profile.preferences.preferredCompanies}
        onChange={async (preferredCompanies) => {
          const updated = await updateProfile({ preferences: { ...profile.preferences, preferredCompanies } }, user.token);
          setProfile(updated);
        }}
      />
      <NotificationSettingsPanel emailOptIn={profile.preferences.emailOptIn} frequency="daily" token={user.token} />
      <h3>Application History</h3>
      <ApplicationHistory jobs={profile.appliedJobs} />
    </PageShell>
  );
}
