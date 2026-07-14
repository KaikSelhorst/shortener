// Digit glyphs generated once from figlet's bundled "3D-ASCII" font
// (https://github.com/patorjk/figlet.js, font source: patorjk.com/software/taag).
// Kept as static data instead of a runtime figlet dependency.
const digitGlyphs: Record<string, string[]> = {
  "0": [
    " ________     ",
    "|\\   __  \\    ",
    "\\ \\  \\|\\  \\   ",
    " \\ \\  \\\\\\  \\  ",
    "  \\ \\  \\\\\\  \\ ",
    "   \\ \\_______\\",
    "    \\|_______|",
  ],
  "1": [
    "  _____     ",
    " / __  \\    ",
    "|\\/_|\\  \\   ",
    "\\|/ \\ \\  \\  ",
    "     \\ \\  \\ ",
    "      \\ \\__\\",
    "       \\|__|",
  ],
  "2": [
    "  _______     ",
    " /  ___  \\    ",
    "/__/|_/  /|   ",
    "|__|//  / /   ",
    "    /  /_/__  ",
    "   |\\________\\",
    "    \\|_______|",
  ],
  "3": [
    " ________     ",
    "|\\_____  \\    ",
    "\\|____|\\ /_   ",
    "      \\|\\  \\  ",
    "     __\\_\\  \\ ",
    "    |\\_______\\",
    "    \\|_______|",
  ],
  "4": [
    " ___   ___     ",
    "|\\  \\ |\\  \\    ",
    "\\ \\  \\\\_\\  \\   ",
    " \\ \\______  \\  ",
    "  \\|_____|\\  \\ ",
    "         \\ \\__\\",
    "          \\|__|",
  ],
  "5": [
    " ________      ",
    "|\\   ____\\     ",
    "\\ \\  \\___|_    ",
    " \\ \\_____  \\   ",
    "  \\|____|\\  \\  ",
    "    ____\\_\\  \\ ",
    "   |\\_________\\",
    "   \\|_________|",
  ],
  "6": [
    " ________     ",
    "|\\   ____\\    ",
    "\\ \\  \\___|    ",
    " \\ \\  \\____   ",
    "  \\ \\  ___  \\ ",
    "   \\ \\_______\\",
    "    \\|_______|",
  ],
  "7": [
    " ________  ",
    "|\\_____  \\ ",
    " \\|___/  /|",
    "     /  / /",
    "    /  / / ",
    "   /__/ /  ",
    "   |__|/   ",
  ],
  "8": [
    " ________     ",
    "|\\   __  \\    ",
    "\\ \\  \\|\\  \\   ",
    " \\ \\   __  \\  ",
    "  \\ \\  \\|\\  \\ ",
    "   \\ \\_______\\",
    "    \\|_______|",
  ],
  "9": [
    " ________     ",
    "|\\  ___  \\    ",
    "\\ \\____   \\   ",
    " \\|____|\\  \\  ",
    "     __\\_\\  \\ ",
    "    |\\_______\\",
    "    \\|_______|",
  ],
};

function glyphRowAt(glyph: string[], height: number, row: number): string {
  const offset = height - glyph.length;
  return row >= offset ? (glyph[row - offset] ?? "") : "";
}

export function renderAsciiCode(code: string): string {
  const glyphs = code.split("").map((ch) => digitGlyphs[ch] ?? []);
  const height = Math.max(0, ...glyphs.map((g) => g.length));

  const grid: string[][] = Array.from({ length: height }, () => []);
  let cursor = 0;

  glyphs.forEach((glyph, i) => {
    const width = Math.max(0, ...glyph.map((l) => l.length));

    let overlap = 0;
    if (i > 0) {
      for (let k = 1; k <= width; k++) {
        let collides = false;
        for (let row = 0; row < height && !collides; row++) {
          const glyphRow = glyphRowAt(glyph, height, row);
          for (let c = 0; c < k; c++) {
            const existing = grid[row][cursor - k + c];
            const incoming = glyphRow[c];
            if (existing && existing !== " " && incoming && incoming !== " ") {
              collides = true;
              break;
            }
          }
        }
        if (collides) break;
        overlap = k;
      }
      overlap = Math.max(0, overlap - 1);
    }

    const start = cursor - overlap;
    for (let row = 0; row < height; row++) {
      const glyphRow = glyphRowAt(glyph, height, row);
      for (let col = 0; col < width; col++) {
        const ch = glyphRow[col] ?? " ";
        if (ch !== " ") grid[row][start + col] = ch;
        else if (grid[row][start + col] === undefined) grid[row][start + col] = " ";
      }
    }
    cursor = start + width;
  });

  return grid.map((row) => row.join("").trimEnd()).join("\n");
}
