const fs = require("fs");

const file = process.env["GITHUB_OUTPUT"];
var stream = fs.createWriteStream(file, { flags: "a" });

const allowed = [
  "PULUMI_CONFIG",
  "PULUMI_KUBERNETES_PROVIDER",
  "PULUMI_KUBERNETES_PROVIDER_VERSION",
  "PULUMI_KUBERNETES_PROVIDER_DOWNLOAD_URL",
  "PULUMI_KUBERNETES_PROVIDER_CHECKSUM",
];

for (const name of allowed) {
  const value = process.env[name];
  if (value === undefined) {
    continue;
  }
  try {
    stream.write(`${name}<<EEEOOOFFF\n${value}\nEEEOOOFFF\n`); // << syntax accommodates multiline strings.
  } catch (err) {
    console.log(`error: failed to set output for ${name}: ${err.message}`);
  }
}

stream.end();
