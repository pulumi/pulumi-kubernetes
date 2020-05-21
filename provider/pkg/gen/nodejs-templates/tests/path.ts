
import * as assert from "assert";
import * as path from "../path";

describe("path.quoteWindowsPath", () => {
    it("escapes Windows path with drive prefix correctly", () => {
        const p = path.quoteWindowsPath("C:\\Users\\grace hopper\\AppData\\Local\\Temp");
        assert.strictEqual(p, "C:\\Users\\grace hopper\\AppData\\Local\\Temp");
    });
    it("escapes Windows path with no drive prefix correctly", () => {
        const p = path.quoteWindowsPath("\\Users\\grace hopper\\AppData\\Local\\Temp");
        assert.strictEqual(p, "\\Users\\grace hopper\\AppData\\Local\\Temp");
    });
    it("escapes relative Windows path correctly", () => {
        const p = path.quoteWindowsPath("Users\\grace hopper\\AppData\\Local\\Temp");
        assert.strictEqual(p, "Users\\grace hopper\\AppData\\Local\\Temp");
    });
    it("escapes Windows repo URL correctly", () => {
        const p = path.quoteWindowsPath("https\://gcsweb.istio.io/gcs/istio-release/releases/1.1.2/charts/");
        assert.strictEqual(p, "https://gcsweb.istio.io/gcs/istio-release/releases/1.1.2/charts/");
    });
});
