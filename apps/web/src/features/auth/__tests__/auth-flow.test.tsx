import { describe, expect, it } from "vitest";
import { isEmail } from "../../../shared/utils/validation";
describe("auth validation",()=>{it("accepts valid email",()=>{expect(isEmail("user@example.com")).toBe(true);});});
