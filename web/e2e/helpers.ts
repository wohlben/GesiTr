import { APIRequestContext, expect } from '@playwright/test';

export function toSlug(name: string): string {
  return name
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-|-$/g, '');
}

export async function createExercise(
  request: APIRequestContext,
  overrides: Record<string, unknown> = {},
) {
  const data = {
    name: 'Test Exercise',
    type: 'STRENGTH',
    technicalDifficulty: 'BEGINNER',
    bodyWeightScaling: 0,
    force: [],
    primaryMuscles: [],
    secondaryMuscles: [],
    suggestedMeasurementParadigms: [],
    description: '',
    instructions: [],
    images: [],
    alternativeNames: [],
    equipmentIds: [],
    ...overrides,
  };
  const res = await request.post('/api/exercises', { data });
  expect(res.ok(), `Failed to create exercise: ${await res.text()}`).toBeTruthy();
  return res.json();
}

export async function updateExercise(
  request: APIRequestContext,
  id: number,
  data: Record<string, unknown>,
) {
  const getRes = await request.get(`/api/exercises/${id}`);
  const current = await getRes.json();
  const merged = { ...current, ...data };
  const res = await request.put(`/api/exercises/${id}`, { data: merged });
  expect(res.ok(), `Failed to update exercise: ${await res.text()}`).toBeTruthy();
  return res.json();
}

export async function deleteExercise(request: APIRequestContext, id: number) {
  await request.delete(`/api/exercises/${id}`);
}

export async function createEquipment(
  request: APIRequestContext,
  overrides: Record<string, unknown> = {},
) {
  const data = {
    name: 'test-equipment',
    displayName: 'Test Equipment',
    description: '',
    category: 'free_weights',
    ...overrides,
  };
  const res = await request.post('/api/equipment', { data });
  expect(res.ok(), `Failed to create equipment: ${await res.text()}`).toBeTruthy();
  return res.json();
}

export async function updateEquipment(
  request: APIRequestContext,
  id: number,
  data: Record<string, unknown>,
) {
  const getRes = await request.get(`/api/equipment/${id}`);
  const current = await getRes.json();
  const merged = { ...current, ...data };
  const res = await request.put(`/api/equipment/${id}`, { data: merged });
  expect(res.ok(), `Failed to update equipment: ${await res.text()}`).toBeTruthy();
  return res.json();
}

export async function deleteEquipment(request: APIRequestContext, id: number) {
  await request.delete(`/api/equipment/${id}`);
}

export async function createExerciseGroup(
  request: APIRequestContext,
  overrides: Record<string, unknown> = {},
) {
  const data = {
    name: 'Test Exercise Group',
    description: '',
    ...overrides,
  };
  const res = await request.post('/api/exercise-groups', { data });
  expect(res.ok(), `Failed to create exercise group: ${await res.text()}`).toBeTruthy();
  return res.json();
}

export async function updateExerciseGroup(
  request: APIRequestContext,
  id: number,
  data: Record<string, unknown>,
) {
  const getRes = await request.get(`/api/exercise-groups/${id}`);
  const current = await getRes.json();
  const merged = { ...current, ...data };
  const res = await request.put(`/api/exercise-groups/${id}`, { data: merged });
  expect(res.ok(), `Failed to update exercise group: ${await res.text()}`).toBeTruthy();
  return res.json();
}

export async function deleteExerciseGroup(request: APIRequestContext, id: number) {
  await request.delete(`/api/exercise-groups/${id}`);
}

export async function createUserExercise(
  request: APIRequestContext,
  compendiumExerciseId: string,
  compendiumVersion: number = 0,
) {
  const res = await request.post('/api/user/exercises', {
    data: { compendiumExerciseId, compendiumVersion },
  });
  expect(res.ok(), `Failed to create user exercise: ${await res.text()}`).toBeTruthy();
  return res.json();
}

export async function deleteUserExercise(request: APIRequestContext, id: number) {
  await request.delete(`/api/user/exercises/${id}`);
}

export async function createUserEquipment(
  request: APIRequestContext,
  compendiumEquipmentId: string,
  compendiumVersion: number = 0,
) {
  const res = await request.post('/api/user/equipment', {
    data: { compendiumEquipmentId, compendiumVersion },
  });
  expect(res.ok(), `Failed to create user equipment: ${await res.text()}`).toBeTruthy();
  return res.json();
}

export async function deleteUserEquipment(request: APIRequestContext, id: number) {
  await request.delete(`/api/user/equipment/${id}`);
}

export async function createWorkout(
  request: APIRequestContext,
  overrides: Record<string, unknown> = {},
) {
  const data = {
    name: 'Test Workout',
    notes: '',
    ...overrides,
  };
  const res = await request.post('/api/user/workouts', { data });
  expect(res.ok(), `Failed to create workout: ${await res.text()}`).toBeTruthy();
  return res.json();
}

export async function deleteWorkout(request: APIRequestContext, id: number) {
  await request.delete(`/api/user/workouts/${id}`);
}

export async function createExerciseScheme(
  request: APIRequestContext,
  overrides: Record<string, unknown> = {},
) {
  const data = {
    measurementType: 'REP_BASED',
    sets: 3,
    reps: 10,
    weight: 60,
    restBetweenSets: 90,
    ...overrides,
  };
  const res = await request.post('/api/user/exercise-schemes', { data });
  expect(res.ok(), `Failed to create exercise scheme: ${await res.text()}`).toBeTruthy();
  return res.json();
}

export async function deleteExerciseScheme(request: APIRequestContext, id: number) {
  await request.delete(`/api/user/exercise-schemes/${id}`);
}

export async function createWorkoutSection(
  request: APIRequestContext,
  overrides: Record<string, unknown> = {},
) {
  const data = {
    type: 'main',
    position: 0,
    ...overrides,
  };
  const res = await request.post('/api/user/workout-sections', { data });
  expect(res.ok(), `Failed to create workout section: ${await res.text()}`).toBeTruthy();
  return res.json();
}

export async function deleteWorkoutSection(request: APIRequestContext, id: number) {
  await request.delete(`/api/user/workout-sections/${id}`);
}

export async function createWorkoutSectionExercise(
  request: APIRequestContext,
  overrides: Record<string, unknown> = {},
) {
  const data = {
    position: 0,
    ...overrides,
  };
  const res = await request.post('/api/user/workout-section-exercises', { data });
  expect(
    res.ok(),
    `Failed to create workout section exercise: ${await res.text()}`,
  ).toBeTruthy();
  return res.json();
}

export async function deleteWorkoutSectionExercise(
  request: APIRequestContext,
  id: number,
) {
  await request.delete(`/api/user/workout-section-exercises/${id}`);
}

export async function fetchWorkoutLogs(
  request: APIRequestContext,
  params: Record<string, string | number> = {},
) {
  const qp = new URLSearchParams();
  for (const [key, value] of Object.entries(params)) {
    qp.set(key, String(value));
  }
  const qs = qp.toString();
  const res = await request.get(`/api/user/workout-logs${qs ? '?' + qs : ''}`);
  return res.json();
}

export async function createWorkoutLog(
  request: APIRequestContext,
  overrides: Record<string, unknown> = {},
) {
  const data = {
    name: 'Test Log',
    ...overrides,
  };
  const res = await request.post('/api/user/workout-logs', { data });
  expect(res.ok(), `Failed to create workout log: ${await res.text()}`).toBeTruthy();
  return res.json();
}

export async function startWorkoutLog(request: APIRequestContext, id: number) {
  const res = await request.post(`/api/user/workout-logs/${id}/start`);
  expect(res.ok(), `Failed to start workout log: ${await res.text()}`).toBeTruthy();
  return res.json();
}

export async function abandonWorkoutLog(request: APIRequestContext, id: number) {
  const res = await request.post(`/api/user/workout-logs/${id}/abandon`);
  expect(res.ok(), `Failed to abandon workout log: ${await res.text()}`).toBeTruthy();
  return res.json();
}

export async function startAdhocWorkoutLog(request: APIRequestContext) {
  const res = await request.post('/api/user/workout-logs/adhoc');
  expect(res.ok(), `Failed to start adhoc workout log: ${await res.text()}`).toBeTruthy();
  return res.json();
}

export async function finishWorkoutLog(request: APIRequestContext, id: number) {
  const res = await request.post(`/api/user/workout-logs/${id}/finish`);
  expect(res.ok(), `Failed to finish workout log: ${await res.text()}`).toBeTruthy();
  return res.json();
}

export async function deleteWorkoutLog(request: APIRequestContext, id: number) {
  await request.delete(`/api/user/workout-logs/${id}`);
}
