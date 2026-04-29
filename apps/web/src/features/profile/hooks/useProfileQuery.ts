import { useEffect, useState } from "react";
import { fetchProfile } from "../api/profile.api";
import { Profile } from "../model/profile.types";

export function useProfileQuery(userId: string, token: string | null) {
  const [profile, setProfile] = useState<Profile | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!token) {
      setProfile(null);
      setLoading(false);
      setError("Missing session token.");
      return;
    }

    const controller = new AbortController();
    setLoading(true);
    setError(null);

    fetchProfile(userId, token, controller.signal)
      .then((data) => {
        if (!controller.signal.aborted) {
          setProfile(data);
        }
      })
      .catch((err: unknown) => {
        if (!controller.signal.aborted) {
          setError(err instanceof Error ? err.message : "Failed to load profile.");
        }
      })
      .finally(() => {
        if (!controller.signal.aborted) {
          setLoading(false);
        }
      });

    return () => controller.abort();
  }, [userId, token]);

  return { profile, loading, error, setProfile };
}
