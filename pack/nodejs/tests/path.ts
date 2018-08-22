import * as assert from "assert";
import * as path from "../path";

describe("path.quoteWindowsPath", () => {
    it("escapes Windows path with drive prefix correctly", () => {
        const p = path.quoteWindowsPath("C:\\Users\\grace hopper\\AppData\\Local\\Temp");
        assert.equal(p, "C:\\Users\\grace hopper\\AppData\\Local\\Temp");
    });
    it("escapes Windows path with no drive prefix correctly", () => {
        const p = path.quoteWindowsPath("\\Users\\grace hopper\\AppData\\Local\\Temp");
        assert.equal(p, "\\Users\\grace hopper\\AppData\\Local\\Temp");
    });
    it("escapes relative Windows path correctly", () => {
        const p = path.quoteWindowsPath("Users\\grace hopper\\AppData\\Local\\Temp");
        assert.equal(p, "Users\\grace hopper\\AppData\\Local\\Temp");
    });
});
