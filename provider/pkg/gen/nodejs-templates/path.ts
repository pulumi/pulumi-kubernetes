// Copyright 2016-2020, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import * as shell from "shell-quote";

/** @ignore */ export function quotePath(path: string): string {
    if (process.platform === "win32") {
        return quoteWindowsPath(path);
    } else {
        return shell.quote([path]);
    }
}

/** @ignore */ export function quoteWindowsPath(path: string): string {
    // Unescape paths for Windows. Taken directly from[1], an unmerged, but LGTM'd PR to the
    // official library.
    //
    // [1]: https://github.com/substack/node-shell-quote/pull/34

    path = String(path).replace(/([A-z]:)?([#!"$&'()*,:;<=>?@\[\\\]^`{|}])/g, "$1\\$2");
    path = path.replace(/\\:/g, ":");
    return path.replace(/\\\\/g, "\\");
}
