import { describe, expect, it, vi } from "vitest";
import { fetchProfile, updateProfile } from "../api/profile.api";

global.fetch = vi.fn(async () => new Response(JSON.stringify({ id: "u_55", preferences: { emailOptIn: false } }), { status: 200 })) as unknown as typeof fetch;

describe("profile behavior", () => {
  it("loads profile by user id", async () => {
    const profile = await fetchProfile("u_55", "test-token");
    expect(profile.id).toBe("u_55");
  });

  it("updates privacy and email settings", async () => {
    global.fetch = vi.fn(async () =>
      new Response(
        JSON.stringify({
          isPrivate: true,
          preferences: {
            keywords: ["frontend"],
            locations: ["Remote"],
            desiredTitles: ["Frontend Engineer"],
            preferredCompanies: ["Stripe"],
            minComp: 120000,
            emailOptIn: false,
            darkMode: false,
          },
        }),
        { status: 200 }
      )
    ) as unknown as typeof fetch;

    const updated = await updateProfile(
      {
        isPrivate: true,
        preferences: {
          keywords: ["frontend"],
          locations: ["Remote"],
          desiredTitles: ["Frontend Engineer"],
          preferredCompanies: ["Stripe"],
          minComp: 120000,
          emailOptIn: false,
          darkMode: false,
        },
      },
      "test-token"
    );

    expect(updated.isPrivate).toBe(true);
    expect(updated.preferences.emailOptIn).toBe(false);
  });
});
