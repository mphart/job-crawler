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
    expect(validateSignupInput("user@example.com", "mason", "password123")).toBeNull();
  });

  it("rejects invalid signup payload", () => {
    expect(validateSignupInput("user@example", "m", "123")).toBeTruthy();
  });
});
