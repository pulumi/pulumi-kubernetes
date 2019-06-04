/**
 * Returns the version of the package containing this file, obtained from the package.json
 * of this package.
 */
export function getVersion(): string {
    let version: string = require("./package.json").version;
    // Node allows for the version to be prefixed by a "v", while semver doesn't.
    // If there is a v, strip it off.
    if (version.indexOf("v") === 0) {
        version = version.slice(1);
    }
    return version;
}
