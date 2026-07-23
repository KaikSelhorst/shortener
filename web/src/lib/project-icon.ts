// Adapted from a standalone visual-hash prototype (opengl.html): a name/seed is
// hashed to an OKLCH hue, then rendered as a Bayer-4x4-dithered lightness
// gradient. Same seed always produces the same icon — a deterministic,
// no-two-projects-look-alike-by-chance avatar, with no image upload needed.
// Simplified from the original for use as a small static UI icon: no
// animation loop (would be wasteful with several icons on screen at once),
// no circular clip (this app's avatars are rounded-md, not rounded-full —
// see DESIGN.md), fixed tuning instead of exposed sliders.

function hashHue(seed: string): number {
  let h = 5381;
  for (let i = 0; i < seed.length; i++) {
    h = (h << 5) + h + seed.charCodeAt(i);
    h |= 0;
  }
  return Math.abs(h) % 360;
}

function oklchToSrgb(l: number, c: number, hDeg: number): [number, number, number] {
  const hRad = (hDeg * Math.PI) / 180;
  const a = c * Math.cos(hRad);
  const b = c * Math.sin(hRad);

  const l_ = l + 0.3963377774 * a + 0.2158037573 * b;
  const m_ = l - 0.1055613458 * a - 0.0638541728 * b;
  const s_ = l - 0.0894841775 * a - 1.291485548 * b;

  const l3 = l_ ** 3;
  const m3 = m_ ** 3;
  const s3 = s_ ** 3;

  let r = 4.0767416621 * l3 - 3.3077115913 * m3 + 0.2309699292 * s3;
  let g = -1.2684380046 * l3 + 2.6097574011 * m3 - 0.3413193965 * s3;
  let bb = -0.0041960863 * l3 - 0.7034186147 * m3 + 1.707614701 * s3;

  const toSrgb = (channel: number) => {
    const clamped = Math.max(0, Math.min(1, channel));
    return clamped <= 0.0031308 ? 12.92 * clamped : 1.055 * clamped ** (1 / 2.4) - 0.055;
  };

  r = toSrgb(r);
  g = toSrgb(g);
  bb = toSrgb(bb);

  return [Math.round(r * 255), Math.round(g * 255), Math.round(bb * 255)];
}

const BAYER_4X4 = [
  [0, 8, 2, 10],
  [12, 4, 14, 6],
  [3, 11, 1, 9],
  [15, 7, 13, 5],
];

function ditherLevel(value: number, x: number, y: number, levels: number): number {
  const threshold = (BAYER_4X4[y % 4][x % 4] + 0.5) / 16;
  const step = 1 / levels;
  const index = Math.max(0, Math.min(levels - 1, Math.floor(value / step + threshold)));
  return index / (levels - 1);
}

const CHROMA = 0.14;
const LIGHTNESS_MIN = 0.35;
const LIGHTNESS_MAX = 0.85;
const LEVELS = 6;
const CELL_SIZE = 2;

export function renderProjectIcon(canvas: HTMLCanvasElement, seed: string, size: number) {
  const ctx = canvas.getContext("2d");
  if (!ctx) return;

  const dpr = Math.min(window.devicePixelRatio || 1, 2);
  canvas.width = size * dpr;
  canvas.height = size * dpr;
  ctx.setTransform(dpr, 0, 0, dpr, 0, 0);

  const hue = hashHue(seed);
  const steps = Math.ceil(size / CELL_SIZE);

  for (let y = 0; y < steps; y++) {
    for (let x = 0; x < steps; x++) {
      const px = x * CELL_SIZE;
      const py = y * CELL_SIZE;
      const nx = px / size;
      const ny = py / size;

      const t = 0.5 + 0.5 * Math.sin((nx + ny) * Math.PI * 0.85);
      const lightness = LIGHTNESS_MIN + (LIGHTNESS_MAX - LIGHTNESS_MIN) * t;
      const dithered = ditherLevel(lightness, x, y, LEVELS);
      const [r, g, b] = oklchToSrgb(dithered, CHROMA, hue);

      ctx.fillStyle = `rgb(${r}, ${g}, ${b})`;
      ctx.fillRect(px, py, CELL_SIZE, CELL_SIZE);
    }
  }
}
