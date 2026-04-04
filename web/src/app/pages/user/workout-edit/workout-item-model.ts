import type { GroupConfigValue } from '$ui/exercise-group-config/exercise-group-config';

export interface WorkoutItemModel {
  itemType: string;
  exerciseId: number | null;
  selectedSchemeId: number | null;
  /** Backend ID of the workout section item — set in edit mode, null in create mode. */
  sectionItemId: number | null;
  groupConfig: GroupConfigValue;
}
