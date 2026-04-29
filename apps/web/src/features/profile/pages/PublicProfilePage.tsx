import { useParams } from "react-router-dom";
import { PageShell } from "../../../shared/components/PageShell";
import { ApplicationHistory } from "../components/ApplicationHistory";
import { ProfileHeader } from "../components/ProfileHeader";
import { useProfileQuery } from "../hooks/useProfileQuery";
export function PublicProfilePage(){const {userId="unknown"}=useParams();const {profile,loading}=useProfileQuery(userId);if(loading||!profile)return <PageShell title="Public Profile"><p>Loading...</p></PageShell>;if(profile.isPrivate)return <PageShell title="Public Profile"><p>This profile is private.</p></PageShell>;return <PageShell title={`Profile: ${profile.username}`}><ProfileHeader profile={profile}/><ApplicationHistory jobs={profile.appliedJobs}/></PageShell>;}
