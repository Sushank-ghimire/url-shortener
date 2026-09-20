export type LinkStatus = "active" | "expired" | "inactive";

export interface ShortLink {
  id: string;
  shortCode: string;
  longUrl: string;
  title?: string;
  clicks: number;
  uniqueClicks: number;
  status: LinkStatus;
  createdAt: string;
  expiresAt?: string;
}

export const links: ShortLink[] = [
  {
    id: "1",
    shortCode: "a8K2x",
    longUrl: "https://github.com/sushankghimire",
    title: "My GitHub",
    clicks: 4231,
    uniqueClicks: 2912,
    status: "active",
    createdAt: "2026-08-30",
  },
  {
    id: "2",
    shortCode: "b91Lm",
    longUrl: "https://example.com/my-portfolio",
    title: "Portfolio",
    clicks: 2821,
    uniqueClicks: 2104,
    status: "active",
    createdAt: "2026-08-28",
  },
  {
    id: "3",
    shortCode: "x72Qp",
    longUrl: "https://example.com/project",
    title: "My Project",
    clicks: 1982,
    uniqueClicks: 1501,
    status: "active",
    createdAt: "2026-08-25",
  },
  {
    id: "4",
    shortCode: "k29Fd",
    longUrl: "https://nextjs.org",
    title: "Next.js",
    clicks: 921,
    uniqueClicks: 703,
    status: "active",
    createdAt: "2026-08-20",
  },
  {
    id: "5",
    shortCode: "m31Px",
    longUrl: "https://www.typescriptlang.org",
    title: "TypeScript",
    clicks: 621,
    uniqueClicks: 492,
    status: "inactive",
    createdAt: "2026-08-17",
  },
  {
    id: "6",
    shortCode: "q82Za",
    longUrl: "https://react.dev",
    title: "React",
    clicks: 381,
    uniqueClicks: 291,
    status: "expired",
    createdAt: "2026-08-10",
    expiresAt: "2026-08-25",
  },
];

export const clicksData = [
  { date: "Aug 01", clicks: 120 },
  { date: "Aug 03", clicks: 210 },
  { date: "Aug 05", clicks: 180 },
  { date: "Aug 07", clicks: 340 },
  { date: "Aug 09", clicks: 280 },
  { date: "Aug 11", clicks: 410 },
  { date: "Aug 13", clicks: 390 },
  { date: "Aug 15", clicks: 520 },
  { date: "Aug 17", clicks: 460 },
  { date: "Aug 19", clicks: 610 },
  { date: "Aug 21", clicks: 580 },
  { date: "Aug 23", clicks: 690 },
  { date: "Aug 25", clicks: 630 },
  { date: "Aug 27", clicks: 760 },
  { date: "Aug 29", clicks: 820 },
  { date: "Aug 31", clicks: 910 },
];

export const countryData = [
  { name: "Nepal", value: 32 },
  { name: "United States", value: 24 },
  { name: "India", value: 18 },
  { name: "United Kingdom", value: 11 },
  { name: "Germany", value: 7 },
  { name: "Other", value: 8 },
];

export const deviceData = [
  { name: "Mobile", value: 61 },
  { name: "Desktop", value: 35 },
  { name: "Tablet", value: 4 },
];

export const browserData = [
  { name: "Chrome", value: 64 },
  { name: "Safari", value: 21 },
  { name: "Firefox", value: 9 },
  { name: "Edge", value: 6 },
];

export const referrerData = [
  { name: "Direct", clicks: 4212 },
  { name: "Google", clicks: 3120 },
  { name: "Twitter / X", clicks: 1820 },
  { name: "LinkedIn", clicks: 1120 },
  { name: "Reddit", clicks: 840 },
];
