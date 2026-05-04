import { describe, expect, it } from "vitest";
import { validateLoginInput } from "../components/LoginForm";
import { validateSignupInput } from "../components/SignupForm";

describe("auth flow validation", () => {
  it("accepts valid login payload", () => {
    expect(validateLoginInput("user@example.com", "password123")).toBeNull();
  });

  it("rejects invalid login payload", () => {
    expect(validateLoginInput("bad", "short")).toContain("valid email");
  });

  it("accepts valid signup payload", () => {
    expect(validateSignupInput("user@example.com", "Mason Lee", "password123", "frontend", "Remote")).toBeNull();
  });

  it("accepts typical personal domains", () => {
    expect(validateSignupInput("tyler@mestery.com", "Tyler Mestery", "password123", "engineer", "Austin, TX")).toBeNull();
    expect(validateLoginInput("tyler@mestery.com", "password123")).toBeNull();
  });

  it("rejects invalid signup payload", () => {
    expect(validateSignupInput("user@example", "Pat", "123", "kw", "here")).toBeTruthy();
  });

  it("rejects short legal name on signup", () => {
    expect(validateSignupInput("user@example.com", "A", "password123", "kw", "here")).toContain("full name");
  });

  it("requires at least one keyword and location", () => {
    expect(validateSignupInput("user@example.com", "Pat Lee", "password123", "", "Remote")).toContain("keyword");
    expect(validateSignupInput("user@example.com", "Pat Lee", "password123", "dev", "")).toContain("location");
  });
});
