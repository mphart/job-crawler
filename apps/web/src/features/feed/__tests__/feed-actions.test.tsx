import { describe, expect, it } from "vitest";
import { markApplied, rejectPosting } from "../api/feed.api";
describe("feed actions api",()=>{it("resolves apply/reject actions",async()=>{await expect(markApplied("job_1")).resolves.toBeUndefined();await expect(rejectPosting("job_1")).resolves.toBeUndefined();});});
