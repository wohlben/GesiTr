export function formatSchemeSummary(scheme: {
  measurementType: string;
  sets?: number | null;
  reps?: number | null;
  weight?: number | null;
  duration?: number | null;
  distance?: number | null;
  targetTime?: number | null;
}): string {
  if (scheme.measurementType === 'REP_BASED') {
    const parts: string[] = [];
    if (scheme.sets) parts.push(`${scheme.sets}x`);
    if (scheme.reps) parts.push(`${scheme.reps}`);
    const setsReps = parts.join('');
    if (scheme.weight) return `${setsReps} @ ${scheme.weight}kg`;
    return setsReps || 'Rep based';
  }
  if (scheme.measurementType === 'TIME_BASED') {
    if (scheme.duration) return `${scheme.duration}s`;
    return 'Time based';
  }
  if (scheme.measurementType === 'DISTANCE_BASED') {
    if (scheme.distance) return `${scheme.distance}m`;
    return 'Distance based';
  }
  return scheme.measurementType;
}
