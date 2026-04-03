import { APIRequestContext, expect, request } from '@playwright/test';

/** Strip server-set fields (BaseModel + owner/version) before sending an update payload. */
function stripServerFields(obj: Record<string, unknown>): Record<string, unknown> {
  const { id, createdAt, updatedAt, deletedAt, owner, version, ...rest } = obj;
  return rest;
}

// Playwright config sets X-User-Id via extraHTTPHeaders for the default user.
// To create resources as a different user, we need a direct API context
// that bypasses the proxy and sets its own X-User-Id header.
// In Docker, there's no proxy — PLAYWRIGHT_TEST_BASE_URL points directly at the API.
// Locally, the proxy runs on :4200/:4300 and the API on :9876.
const E2E_API_BASE = process.env['PLAYWRIGHT_TEST_BASE_URL'] ?? 'http://localhost:9876';

export async function createApiContextAs(userId: string): Promise<APIRequestContext> {
  return request.newContext({
    baseURL: E2E_API_BASE,
    extraHTTPHeaders: { 'X-User-Id': userId },
  });
}

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
    names: ['Test Exercise'],
    type: 'STRENGTH',
    technicalDifficulty: 'beginner',
    bodyWeightScaling: 0,
    force: [],
    primaryMuscles: [],
    secondaryMuscles: [],
    suggestedMeasurementParadigms: [],
    description: '',
    instructions: [],
    images: [],
    equipmentIds: [],
    ...overrides,
  };
  const res = await request.post('/api/exercises', { data });
  expect(res.ok(), `Failed to create exercise: ${await res.text()}`).toBeTruthy();
  return res.json();
}

export async function createExerciseAs(userId: string, overrides: Record<string, unknown> = {}) {
  const ctx = await createApiContextAs(userId);
  const data = {
    names: ['Test Exercise'],
    type: 'STRENGTH',
    technicalDifficulty: 'beginner',
    bodyWeightScaling: 0,
    force: [],
    primaryMuscles: [],
    secondaryMuscles: [],
    suggestedMeasurementParadigms: [],
    description: '',
    instructions: [],
    images: [],
    equipmentIds: [],
    ...overrides,
  };
  const res = await ctx.post('/api/exercises', { data });
  expect(res.ok(), `Failed to create exercise as ${userId}: ${await res.text()}`).toBeTruthy();
  const json = await res.json();
  await ctx.dispose();
  return json;
}

export async function deleteExerciseAs(id: number, userId: string) {
  const ctx = await createApiContextAs(userId);
  await ctx.delete(`/api/exercises/${id}`);
  await ctx.dispose();
}

export async function updateExercise(
  request: APIRequestContext,
  id: number,
  data: Record<string, unknown>,
) {
  const getRes = await request.get(`/api/exercises/${id}`);
  const current = await getRes.json();
  const merged = { ...stripServerFields(current), ...data };
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
  const merged = { ...stripServerFields(current), ...data };
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
  const merged = { ...stripServerFields(current), ...data };
  const res = await request.put(`/api/exercise-groups/${id}`, { data: merged });
  expect(res.ok(), `Failed to update exercise group: ${await res.text()}`).toBeTruthy();
  return res.json();
}

export async function deleteExerciseGroup(request: APIRequestContext, id: number) {
  await request.delete(`/api/exercise-groups/${id}`);
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
  const res = await request.post('/api/workouts', { data });
  expect(res.ok(), `Failed to create workout: ${await res.text()}`).toBeTruthy();
  return res.json();
}

export async function deleteWorkout(request: APIRequestContext, id: number) {
  await request.delete(`/api/workouts/${id}`);
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
  const res = await request.post('/api/workout-sections', { data });
  expect(res.ok(), `Failed to create workout section: ${await res.text()}`).toBeTruthy();
  return res.json();
}

export async function deleteWorkoutSection(request: APIRequestContext, id: number) {
  await request.delete(`/api/workout-sections/${id}`);
}

export async function createWorkoutSectionItem(
  request: APIRequestContext,
  overrides: Record<string, unknown> = {},
) {
  const data = {
    type: 'exercise',
    position: 0,
    ...overrides,
  };
  const res = await request.post('/api/workout-section-items', { data });
  expect(res.ok(), `Failed to create workout section item: ${await res.text()}`).toBeTruthy();
  return res.json();
}

export async function deleteWorkoutSectionItem(request: APIRequestContext, id: number) {
  await request.delete(`/api/workout-section-items/${id}`);
}

export async function upsertSchemeSectionItem(
  request: APIRequestContext,
  data: { exerciseSchemeId: number; workoutSectionItemId: number },
) {
  const res = await request.put('/api/user/exercise-scheme-section-items', { data });
  expect(res.ok(), `Failed to upsert scheme section item: ${await res.text()}`).toBeTruthy();
  return res.json();
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

export async function createWorkoutSchedule(
  request: APIRequestContext,
  overrides: Record<string, unknown> = {},
) {
  const data = {
    ...overrides,
  };
  const res = await request.post('/api/user/workout-schedules', { data });
  expect(res.ok(), `Failed to create workout schedule: ${await res.text()}`).toBeTruthy();
  return res.json();
}

export async function createSchedulePeriod(
  request: APIRequestContext,
  overrides: Record<string, unknown> = {},
) {
  const data = {
    type: 'fixed_date',
    ...overrides,
  };
  const res = await request.post('/api/user/schedule-periods', { data });
  expect(res.ok(), `Failed to create schedule period: ${await res.text()}`).toBeTruthy();
  return res.json();
}

export async function createScheduleCommitment(
  request: APIRequestContext,
  overrides: Record<string, unknown> = {},
) {
  const res = await request.post('/api/user/schedule-commitments', { data: overrides });
  expect(res.ok(), `Failed to create schedule commitment: ${await res.text()}`).toBeTruthy();
  return res.json();
}
