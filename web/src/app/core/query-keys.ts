type Filters = Record<string, string | number | undefined>;

export const exerciseKeys = {
  all: () => ['exercises'] as const,
  list: (filters: Filters) => [...exerciseKeys.all(), 'list', filters] as const,
  detail: (id: number) => [...exerciseKeys.all(), 'detail', id] as const,
  versions: (id: number) => [...exerciseKeys.all(), 'versions', id] as const,
  version: (templateId: string, version: number) =>
    [...exerciseKeys.all(), 'version', templateId, version] as const,
};

export const equipmentKeys = {
  all: () => ['equipment'] as const,
  list: (filters: Filters) => [...equipmentKeys.all(), 'list', filters] as const,
  detail: (id: number) => [...equipmentKeys.all(), 'detail', id] as const,
  versions: (id: number) => [...equipmentKeys.all(), 'versions', id] as const,
  version: (templateId: string, version: number) =>
    [...equipmentKeys.all(), 'version', templateId, version] as const,
};

export const exerciseGroupKeys = {
  all: () => ['exercise-groups'] as const,
  list: (filters: Filters) => [...exerciseGroupKeys.all(), 'list', filters] as const,
  detail: (id: number) => [...exerciseGroupKeys.all(), 'detail', id] as const,
};

export const userExerciseKeys = {
  all: () => ['user-exercises'] as const,
  list: () => [...userExerciseKeys.all(), 'list'] as const,
  detail: (id: number) => [...userExerciseKeys.all(), 'detail', id] as const,
};

export const userEquipmentKeys = {
  all: () => ['user-equipment'] as const,
  list: () => [...userEquipmentKeys.all(), 'list'] as const,
  detail: (id: number) => [...userEquipmentKeys.all(), 'detail', id] as const,
};
