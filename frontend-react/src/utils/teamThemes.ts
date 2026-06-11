export interface TeamTheme {
  primary: string;
  bg: string;
  glow: string;
}

const THEMES: Record<string, TeamTheme> = {
  "optic texas": {
    primary: "#78BE20",
    bg: "#111111",
    glow: "rgba(120,190,32,0.18)",
  },
  "optic chicago": {
    primary: "#78BE20",
    bg: "#111111",
    glow: "rgba(120,190,32,0.18)",
  },
  "dallas empire": {
    primary: "#A8272A",
    bg: "#111111",
    glow: "rgba(168,39,42,0.18)",
  },
  "atlanta faze": {
    primary: "#E8002D",
    bg: "#111111",
    glow: "rgba(232,0,45,0.18)",
  },
  "faze vegas": {
    primary: "#E8002D",
    bg: "#111111",
    glow: "rgba(232,0,45,0.18)",
  },
  "los angeles thieves": {
    primary: "#E63329",
    bg: "#0B0B0B",
    glow: "rgba(230,51,41,0.18)",
  },
  "100 thieves": {
    primary: "#E63329",
    bg: "#0B0B0B",
    glow: "rgba(230,51,41,0.18)",
  },
  "boston breach": {
    primary: "#3BAA35",
    bg: "#111111",
    glow: "rgba(59,170,53,0.18)",
  },
  "toronto ultra": {
    primary: "#8B1FFF",
    bg: "#111111",
    glow: "rgba(139,31,255,0.18)",
  },
  "toronto koi": {
    primary: "#6A2DFF",
    bg: "#111111",
    glow: "rgba(106,45,255,0.18)",
  },
  "los angeles guerrillas": {
    primary: "#4CAF50",
    bg: "#111111",
    glow: "rgba(76,175,80,0.18)",
  },
  "la guerrillas m8": {
    primary: "#B8A7FF",
    bg: "#2B2B2B",
    glow: "rgba(184,167,255,0.18)",
  },
  "miami heretics": {
    primary: "#FF6B35",
    bg: "#111111",
    glow: "rgba(255,107,53,0.18)",
  },
  "team heretics": {
    primary: "#FF6B35",
    bg: "#111111",
    glow: "rgba(255,107,53,0.18)",
  },
  "new york subliners": {
    primary: "#F6EB14",
    bg: "#111111",
    glow: "rgba(246,235,20,0.18)",
  },
  "cloud9 new york": {
    primary: "#00AEEF",
    bg: "#101820",
    glow: "rgba(0,174,239,0.18)",
  },
  "seattle surge": {
    primary: "#00FF87",
    bg: "#111111",
    glow: "rgba(0,255,135,0.18)",
  },
  "vancouver surge": {
    primary: "#00A7E1",
    bg: "#111111",
    glow: "rgba(0,167,225,0.18)",
  },
  "minnesota røkkr": {
    primary: "#5B2D8E",
    bg: "#111111",
    glow: "rgba(91,45,142,0.18)",
  },
  "g2 minnesota": {
    primary: "#FF4C00",
    bg: "#111111",
    glow: "rgba(255,76,0,0.18)",
  },
  "carolina royal ravens": {
    primary: "#1B4F8E",
    bg: "#111111",
    glow: "rgba(27,79,142,0.18)",
  },
  "london royal ravens": {
    primary: "#1B4F8E",
    bg: "#111111",
    glow: "rgba(27,79,142,0.18)",
  },
  "paris legion": {
    primary: "#003087",
    bg: "#111111",
    glow: "rgba(0,48,135,0.18)",
  },
  "florida mutineers": {
    primary: "#FFCD00",
    bg: "#111111",
    glow: "rgba(255,205,0,0.18)",
  },
  "riyadh falcons": {
    primary: "#00A86B",
    bg: "#111111",
    glow: "rgba(0,168,107,0.18)",
  },
  "las vegas legion": {
    primary: "#E4B062",
    bg: "#111111",
    glow: "rgba(228,176,98,0.18)",
  },
};

const FALLBACK: TeamTheme = {
  primary: "#4a4a5a",
  bg: "#111111",
  glow: "rgba(74,74,90,0.12)",
};

export function getTeamTheme(teamName: string): TeamTheme {
  return THEMES[teamName.toLowerCase()] ?? FALLBACK;
}
