export function formatDailyReport(report) {
  const riskPoints = report?.riskPoints || report?.risk_points || [];
  const needsConfirmation = report?.needsConfirmation || report?.needs_confirmation || [];
  return {
    title: `Daily Review ${report?.date || ""}`.trim(),
    summary: report?.summary || "No report summary.",
    riskPoints,
    needsConfirmation,
    dataTime: report?.dataTime || report?.data_time || "",
    empty: riskPoints.length === 0,
  };
}
