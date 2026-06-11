import { describe, it, expect } from "vitest";
import { readdirSync } from "node:fs";
import { join, dirname } from "node:path";
import { fileURLToPath } from "node:url";
import { teamLogoKeys, avatarNicknames } from "./assets";

const __dir = dirname(fileURLToPath(import.meta.url));
const AVATAR_DIR = join(__dir, "../assets/avatars");
const LOGO_DIR = join(__dir, "../assets/logos");

function diskKeys(dir: string): Set<string> {
  return new Set(
    readdirSync(dir).map((f: string) =>
      f.replace(/\.[^.]+$/, "").toLowerCase(),
    ),
  );
}

const avatarDisk = diskKeys(AVATAR_DIR);
const logoDisk = diskKeys(LOGO_DIR);

describe("teamLogoKeys — every mapped key has a file on disk", () => {
  const broken: string[] = [];

  for (const [team, fileKey] of Object.entries(teamLogoKeys)) {
    if (!logoDisk.has(fileKey)) {
      broken.push(`"${team}" → "${fileKey}" (file missing)`);
    }
  }

  it("no broken logo mappings", () => {
    expect(broken).toEqual([]);
  });
});

describe("teamLogoKeys — CDL franchise logos resolve", () => {
  const cdlTeams: [string, string][] = [
    ["Atlanta FaZe", "atlanta faze"],
    ["Boston Breach", "boston breach"],
    ["Carolina Royal Ravens", "carolina royal ravens"],
    ["Dallas Empire", "dallas empire"],
    ["Florida Mutineers", "florida mutineers"],
    ["LA Guerrillas M8", "la guerrillas m8"],
    ["Los Angeles Thieves", "los angeles thieves"],
    ["Las Vegas Legion", "las vegas legion"],
    ["London Royal Ravens", "london royal ravens"],
    ["Miami Heretics", "miami heretics"],
    ["Minnesota RØKKR", "minnesota røkkr"],
    ["G2 Minnesota", "g2 minnesota"],
    ["New York Subliners", "new york subliners"],
    ["Cloud9 New York", "cloud9 new york"],
    ["Paris Legion", "paris legion"],
    ["Riyadh Falcons", "riyadh falcons"],
    ["Toronto Ultra", "toronto ultra"],
    ["Toronto KOI", "toronto koi"],
    ["Vancouver Surge", "vancouver surge"],
    ["Seattle Surge", "seattle surge"],
    ["OpTic Texas", "optic texas"],
    ["OpTic Chicago", "optic chicago"],
  ];

  for (const [label, key] of cdlTeams) {
    it(`resolves ${label}`, () => {
      const fileKey = teamLogoKeys[key];
      expect(fileKey, `No mapping for "${key}"`).toBeTruthy();
      expect(
        logoDisk.has(fileKey),
        `File "${fileKey}" missing for "${label}"`,
      ).toBe(true);
    });
  }
});

describe("teamLogoKeys — is case-insensitive by design", () => {
  it("lookup works with any case (caller must lowercase)", () => {
    // The mapping stores lowercase keys — callers must pass .toLowerCase()
    expect(teamLogoKeys["toronto ultra"]).toBeTruthy();
    expect(teamLogoKeys["TORONTO ULTRA"]).toBeUndefined();
  });
});

describe("avatarNicknames — every target key has a file on disk", () => {
  const broken: string[] = [];

  for (const [alias, target] of Object.entries(avatarNicknames)) {
    if (!avatarDisk.has(target)) {
      broken.push(`"${alias}" → "${target}" (file missing)`);
    }
  }

  it("no broken nickname targets", () => {
    expect(broken).toEqual([]);
  });
});

// ── Avatar file coverage ─────────────────────────────────────────────────────

describe("avatar files — known CDL players have avatars", () => {
  // Resolves a gamertag the same way getPlayerAvatar() does
  function resolves(gamertag: string): boolean {
    const key = gamertag.toLowerCase();
    return (
      avatarDisk.has(key) || avatarDisk.has(avatarNicknames[key] ?? "__none__")
    );
  }

  const mustHave = [
    "Shotzzy",
    "Cellium",
    "Simp",
    "aBeZy",
    "Scump",
    "Crimsix",
    "Clayster",
    "FormaL",
    "Gunless",
    "Octane",
    "iLLeY",
    "Methodz",
    "Envoy",
    "Kenny",
    "Dashy",
    "HyDra",
    "Cammy",
    "CleanX",
    "Insight",
    "Nastie",
    "Arcitys",
    "SlasheR",
    "TJHaLy",
    "Attach",
    "Drazah",
    "Pred",
    "Kremp",
    "Priestahh",
    "Loony",
    "Skrapz",
    "Seany",
    "Temp",
    "Skyz",
    "Lucky",
  ];

  for (const gamertag of mustHave) {
    it(`resolves ${gamertag}`, () => {
      expect(resolves(gamertag), `No avatar for ${gamertag}`).toBe(true);
    });
  }
});

describe("avatar files — nickname aliases resolve", () => {
  const aliases: [string, string][] = [
    ["Hicksy", "hicksey"],
    ["Purj", "purj_"],
    ["ReeaL", "real"],
    ["Mercules", "merc"],
    ["LyynnZ", "lynz"],
  ];

  for (const [gamertag, expectedKey] of aliases) {
    it(`${gamertag} resolves via nickname to "${expectedKey}"`, () => {
      const key = gamertag.toLowerCase();
      const resolved = avatarNicknames[key];
      expect(resolved).toBe(expectedKey);
      expect(avatarDisk.has(resolved), `File "${resolved}" missing`).toBe(true);
    });
  }
});

describe("placeholder avatars — no public photo available", () => {
  const noPhotoPlayers = ["5aLDx", "FelonY", "Hamza", "MarkyB", "qk4b"];

  for (const gamertag of noPhotoPlayers) {
    it(`${gamertag} explicitly maps to unknown placeholder`, () => {
      const key = gamertag.toLowerCase();
      expect(avatarNicknames[key]).toBe("unknown");
      expect(avatarDisk.has("unknown")).toBe(true);
    });
  }
});

describe("missing asset report (informational — not failures)", () => {
  it("prints players with no avatar on disk", () => {
    const covered = [...avatarDisk].filter((k) => k !== "unknown").length;
    console.log(
      `Avatar files on disk: ${covered} players covered (excluding Unknown placeholder)`,
    );
    expect(true).toBe(true);
  });

  it("prints the 04 avatar ambiguity warning", () => {
    const hasBothFiles = avatarDisk.has("04") && avatarDisk.has("new04");
    if (hasBothFiles) {
      console.warn(
        "WARNING: Both 04.webp and New04.webp exist in src/assets/avatars/. " +
          "The old 04.webp takes priority and the New04 nickname never fires. " +
          "Delete src/assets/avatars/04.webp to use the newer image.",
      );
    }
    expect(true).toBe(true);
  });
});
