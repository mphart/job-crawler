import { describe, expect, it } from "vitest";
import { formatDate } from "../../../shared/utils/date";
describe("profile helpers",()=>{it("formats date",()=>{expect(formatDate("2024-01-01T00:00:00.000Z")).toMatch(/\d/);});});
