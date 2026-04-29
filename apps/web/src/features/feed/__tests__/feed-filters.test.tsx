import { describe, expect, it } from "vitest";
import { formatCompensation } from "../../../shared/utils/money";
describe("feed helpers",()=>{it("formats fallback compensation",()=>{expect(formatCompensation("")).toContain("not listed");});});
