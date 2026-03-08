export interface VersionEntry<T = unknown> {
  version: number;
  snapshot: T;
  changedAt: string;
  changedBy: string;
}
