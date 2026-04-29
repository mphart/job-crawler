import { describe, expect, it } from "vitest";
import { fetchProfile, updateProfile } from "../api/profile.api";

describe("profile behavior", () => {
  it("loads profile by user id", async () => {
    const profile = await fetchProfile("u_55");
    expect(profile.id).toBe("u_55");
  });

  it("updates privacy and email settings", async () => {
    const updated = await updateProfile({
      isPrivate: true,
      preferences: {
        keywords: ["frontend"],
        locations: ["Remote"],
        desiredTitles: ["Frontend Engineer"],
        minComp: 120000,
        emailOptIn: false,
        darkMode: false,
      },
    });

    expect(updated.isPrivate).toBe(true);
    expect(updated.preferences.emailOptIn).toBe(false);
  });
});
