type Filters = Record<string, string | number | undefined>;

export const exerciseKeys = {
  all: () => ['exercises'] as const,
  list: (filters: Filters) => [...exerciseKeys.all(), 'list', filters] as const,
  detail: (id: number) => [...exerciseKeys.all(), 'detail', id] as const,
  permissions: (id: number) => [...exerciseKeys.all(), 'permissions', id] as const,
  versions: (id: number) => [...exerciseKeys.all(), 'versions', id] as const,
  version: (id: number, version: number) =>
    [...exerciseKeys.all(), 'version', id, version] as const,
};

export const equipmentKeys = {
  all: () => ['equipment'] as const,
  list: (filters: Filters) => [...equipmentKeys.all(), 'list', filters] as const,
  detail: (id: number) => [...equipmentKeys.all(), 'detail', id] as const,
  permissions: (id: number) => [...equipmentKeys.all(), 'permissions', id] as const,
  versions: (id: number) => [...equipmentKeys.all(), 'versions', id] as const,
  version: (id: number, version: number) =>
    [...equipmentKeys.all(), 'version', id, version] as const,
};

export const exerciseGroupKeys = {
  all: () => ['exercise-groups'] as const,
  list: (filters: Filters) => [...exerciseGroupKeys.all(), 'list', filters] as const,
  detail: (id: number) => [...exerciseGroupKeys.all(), 'detail', id] as const,
  permissions: (id: number) => [...exerciseGroupKeys.all(), 'permissions', id] as const,
};

export const exerciseRelationshipKeys = {
  all: () => ['exercise-relationships'] as const,
  list: (filters?: Filters) => [...exerciseRelationshipKeys.all(), 'list', filters] as const,
};

export const equipmentRelationshipKeys = {
  all: () => ['equipment-relationships'] as const,
  list: (filters?: Filters) => [...equipmentRelationshipKeys.all(), 'list', filters] as const,
};

export const masteryKeys = {
  all: () => ['mastery'] as const,
  list: () => [...masteryKeys.all(), 'list'] as const,
  detail: (exerciseId: number) => [...masteryKeys.all(), 'detail', exerciseId] as const,
};

export const equipmentMasteryKeys = {
  all: () => ['equipment-mastery'] as const,
  list: () => [...equipmentMasteryKeys.all(), 'list'] as const,
  detail: (equipmentId: number) => [...equipmentMasteryKeys.all(), 'detail', equipmentId] as const,
};

export const workoutKeys = {
  all: () => ['workouts'] as const,
  list: () => [...workoutKeys.all(), 'list'] as const,
  detail: (id: number) => [...workoutKeys.all(), 'detail', id] as const,
  permissions: (id: number) => [...workoutKeys.all(), 'permissions', id] as const,
};

export const exerciseSchemeKeys = {
  all: () => ['exercise-schemes'] as const,
  list: (filters?: Filters) => [...exerciseSchemeKeys.all(), 'list', filters] as const,
  detail: (id: number) => [...exerciseSchemeKeys.all(), 'detail', id] as const,
};

export const workoutLogKeys = {
  all: () => ['workout-logs'] as const,
  list: (filters?: Filters) => [...workoutLogKeys.all(), 'list', filters] as const,
  detail: (id: number) => [...workoutLogKeys.all(), 'detail', id] as const,
};

export const exerciseLogKeys = {
  all: () => ['exercise-logs'] as const,
  list: (filters?: Filters) => [...exerciseLogKeys.all(), 'list', filters] as const,
  detail: (id: number) => [...exerciseLogKeys.all(), 'detail', id] as const,
};

export const workoutGroupKeys = {
  all: () => ['workout-groups'] as const,
  list: () => [...workoutGroupKeys.all(), 'list'] as const,
  detail: (id: number) => [...workoutGroupKeys.all(), 'detail', id] as const,
  memberships: (groupId: number) => [...workoutGroupKeys.all(), 'memberships', groupId] as const,
};

export const workoutScheduleKeys = {
  all: () => ['workout-schedules'] as const,
  list: (workoutId?: number) => [...workoutScheduleKeys.all(), 'list', workoutId] as const,
  detail: (id: number) => [...workoutScheduleKeys.all(), 'detail', id] as const,
};

export const schedulePeriodKeys = {
  all: () => ['schedule-periods'] as const,
  list: (scheduleId?: number) => [...schedulePeriodKeys.all(), 'list', scheduleId] as const,
};

export const scheduleCommitmentKeys = {
  all: () => ['schedule-commitments'] as const,
  list: (periodId?: number) => [...scheduleCommitmentKeys.all(), 'list', periodId] as const,
};

export const namePreferenceKeys = {
  all: () => ['exercise-name-preferences'] as const,
  list: () => [...namePreferenceKeys.all(), 'list'] as const,
};
