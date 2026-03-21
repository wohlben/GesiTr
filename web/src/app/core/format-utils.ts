import { WorkoutLogExerciseSet, WorkoutLogItemStatusFinished } from '$generated/user-models';

export function formatBreak(seconds?: number | null): string {
  if (seconds == null) return '';
  if (seconds >= 60) {
    const min = Math.floor(seconds / 60);
    const sec = seconds % 60;
    return sec > 0 ? `${min}m ${sec}s` : `${min}m`;
  }
  return `${seconds}s`;
}

export function formatTarget(set: WorkoutLogExerciseSet, measurementType: string): string {
  if (measurementType === 'REP_BASED') {
    const parts: string[] = [];
    if (set.targetReps != null) parts.push(`${set.targetReps} reps`);
    if (set.targetWeight != null) parts.push(`${set.targetWeight}kg`);
    return parts.join(' @ ') || '-';
  }
  if (measurementType === 'TIME_BASED') {
    if (set.targetDuration != null) return `${set.targetDuration}s`;
    return '-';
  }
  if (measurementType === 'DISTANCE_BASED') {
    if (set.targetDistance != null) return `${set.targetDistance}m`;
    return '-';
  }
  return '-';
}

export function formatActual(set: WorkoutLogExerciseSet, measurementType: string): string {
  if (set.status !== WorkoutLogItemStatusFinished || !set.exerciseLog) return '-';
  const log = set.exerciseLog;
  if (measurementType === 'REP_BASED') {
    const parts: string[] = [];
    if (log.reps != null) parts.push(`${log.reps} reps`);
    if (log.weight != null) parts.push(`${log.weight}kg`);
    return parts.join(' @ ') || '-';
  }
  if (measurementType === 'TIME_BASED') {
    if (log.duration != null) return `${log.duration}s`;
    return '-';
  }
  if (measurementType === 'DISTANCE_BASED') {
    if (log.distance != null) return `${log.distance}m`;
    return '-';
  }
  return '-';
}

export function formatSetValue(set: WorkoutLogExerciseSet, measurementType: string): string {
  if (set.status === WorkoutLogItemStatusFinished) return formatActual(set, measurementType);
  return formatTarget(set, measurementType);
}

export function formatCountdown(seconds: number): string {
  if (seconds >= 60) {
    const min = Math.floor(seconds / 60);
    const sec = seconds % 60;
    return `${min}:${sec.toString().padStart(2, '0')}`;
  }
  return `${seconds}s`;
}
