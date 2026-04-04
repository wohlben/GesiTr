import type { GroupConfigValue } from '$ui/exercise-group-config/exercise-group-config';

export interface WorkoutItemModel {
  itemType: string;
  exerciseId: number | null;
  selectedSchemeId: number | null;
  groupConfig: GroupConfigValue;
}
