export interface StartSetModel {
  id: number | null;
  targetReps: number | null;
  targetWeight: number | null;
  targetDuration: number | null;
  targetDistance: number | null;
  targetTime: number | null;
  restAfterSeconds: number | null;
}

export interface StartExerciseModel {
  id: number | null;
  sourceExerciseSchemeId: number;
  breakAfterSeconds: number | null;
  sets: StartSetModel[];
}

export interface PendingGroupModel {
  groupId: number;
  groupName: string;
  members: { id: number; name: string }[];
  position: number;
}

export interface StartSectionModel {
  id: number | null;
  type: string;
  label: string;
  exercises: StartExerciseModel[];
  pendingGroups: PendingGroupModel[];
}

export interface StartModel {
  name: string;
  notes: string;
  sections: StartSectionModel[];
}
