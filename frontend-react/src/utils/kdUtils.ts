export function getKdColorClass(kd: number | null | undefined): string {
  if (kd == null) return "text-[#a3a3a3]";
  if (kd > 1.0) return "text-green-400";
  if (kd < 1.0) return "text-red-400";
  return "text-[#a3a3a3]"; // exactly 1.00
}
