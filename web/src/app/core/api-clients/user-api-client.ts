import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { firstValueFrom, map } from 'rxjs';
import { PaginatedResponse } from './paginated-response';
import { Exercise, Equipment } from '$generated/models';
import { ExerciseScheme, ExerciseSchemeSectionItem } from '$generated/user-exercisescheme';
import {
  Workout,
  WorkoutSection,
  WorkoutSectionItem,
  WorkoutLog,
  WorkoutLogSection,
  WorkoutLogExercise,
  WorkoutLogExerciseSet,
  ExerciseLog,
} from '$generated/user-models';
import {
  WorkoutGroup,
  WorkoutGroupMembership,
  WorkoutGroupRole,
} from '$generated/user-workoutgroup';
import {
  WorkoutSchedule,
  SchedulePeriod,
  ScheduleCommitment,
} from '$generated/user-workoutschedule';
import { ExerciseMastery, EquipmentMastery } from '$generated/user-mastery';

@Injectable({ providedIn: 'root' })
export class UserApiClient {
  private http = inject(HttpClient);

  // Mastery
  fetchMasteryList(): Promise<ExerciseMastery[]> {
    return firstValueFrom(this.http.get<ExerciseMastery[]>('/api/user/mastery'));
  }

  fetchMastery(exerciseId: number): Promise<ExerciseMastery> {
    return firstValueFrom(this.http.get<ExerciseMastery>(`/api/user/mastery/${exerciseId}`));
  }

  // Equipment Mastery
  fetchEquipmentMasteryList(): Promise<EquipmentMastery[]> {
    return firstValueFrom(this.http.get<EquipmentMastery[]>('/api/user/equipment-mastery'));
  }

  fetchEquipmentMastery(equipmentId: number): Promise<EquipmentMastery> {
    return firstValueFrom(
      this.http.get<EquipmentMastery>(`/api/user/equipment-mastery/${equipmentId}`),
    );
  }

  fetchUserExercise(id: number): Promise<Exercise> {
    return firstValueFrom(this.http.get<Exercise>(`/api/exercises/${id}`));
  }

  createUserExercise(data: Partial<Exercise>): Promise<Exercise> {
    return firstValueFrom(this.http.post<Exercise>('/api/exercises', data));
  }

  createUserEquipment(data: Partial<Equipment>): Promise<Equipment> {
    return firstValueFrom(this.http.post<Equipment>('/api/equipment', data));
  }

  // Workouts
  fetchWorkouts(filters?: Record<string, string | number | undefined>): Promise<Workout[]> {
    const qp = new URLSearchParams();
    if (filters) {
      for (const [k, v] of Object.entries(filters)) {
        if (v != null) qp.set(k, String(v));
      }
    }
    const qs = qp.toString();
    return firstValueFrom(
      this.http
        .get<PaginatedResponse<Workout>>(`/api/workouts${qs ? '?' + qs : ''}`)
        .pipe(map((res) => res.items)),
    );
  }

  fetchWorkout(id: number): Promise<Workout> {
    return firstValueFrom(this.http.get<Workout>(`/api/workouts/${id}`));
  }

  fetchWorkoutPermissions(id: number): Promise<{ permissions: string[] }> {
    return firstValueFrom(
      this.http.get<{ permissions: string[] }>(`/api/workouts/${id}/permissions`),
    );
  }

  createWorkout(data: Partial<Workout>): Promise<Workout> {
    return firstValueFrom(this.http.post<Workout>('/api/workouts', data));
  }

  updateWorkout(id: number, data: Partial<Workout>): Promise<Workout> {
    return firstValueFrom(this.http.put<Workout>(`/api/workouts/${id}`, data));
  }

  deleteWorkout(id: number): Promise<void> {
    return firstValueFrom(this.http.delete<void>(`/api/workouts/${id}`));
  }

  // Workout Sections
  createWorkoutSection(data: Partial<WorkoutSection>): Promise<WorkoutSection> {
    return firstValueFrom(this.http.post<WorkoutSection>('/api/workout-sections', data));
  }

  deleteWorkoutSection(id: number): Promise<void> {
    return firstValueFrom(this.http.delete<void>(`/api/workout-sections/${id}`));
  }

  // Workout Section Items
  createWorkoutSectionItem(data: Partial<WorkoutSectionItem>): Promise<WorkoutSectionItem> {
    return firstValueFrom(this.http.post<WorkoutSectionItem>('/api/workout-section-items', data));
  }

  deleteWorkoutSectionItem(id: number): Promise<void> {
    return firstValueFrom(this.http.delete<void>(`/api/workout-section-items/${id}`));
  }

  // Exercise Schemes
  fetchExerciseSchemes(params?: { exerciseId?: number }): Promise<ExerciseScheme[]> {
    const qp = new URLSearchParams();
    if (params?.exerciseId != null) qp.set('exerciseId', String(params.exerciseId));
    const qs = qp.toString();
    return firstValueFrom(
      this.http.get<ExerciseScheme[]>(`/api/user/exercise-schemes${qs ? '?' + qs : ''}`),
    );
  }

  // Exercise Scheme Section Items (join table)
  fetchSchemeSectionItems(workoutSectionItemIds: number[]): Promise<ExerciseSchemeSectionItem[]> {
    const ids = workoutSectionItemIds.join(',');
    return firstValueFrom(
      this.http.get<ExerciseSchemeSectionItem[]>(
        `/api/user/exercise-scheme-section-items?workoutSectionItemIds=${ids}`,
      ),
    );
  }

  upsertSchemeSectionItem(data: {
    exerciseSchemeId: number;
    workoutSectionItemId: number;
  }): Promise<ExerciseSchemeSectionItem> {
    return firstValueFrom(
      this.http.put<ExerciseSchemeSectionItem>('/api/user/exercise-scheme-section-items', data),
    );
  }

  deleteSchemeSectionItem(id: number): Promise<void> {
    return firstValueFrom(this.http.delete<void>(`/api/user/exercise-scheme-section-items/${id}`));
  }

  fetchExerciseScheme(id: number): Promise<ExerciseScheme> {
    return firstValueFrom(this.http.get<ExerciseScheme>(`/api/user/exercise-schemes/${id}`));
  }

  createExerciseScheme(data: Partial<ExerciseScheme>): Promise<ExerciseScheme> {
    return firstValueFrom(this.http.post<ExerciseScheme>('/api/user/exercise-schemes', data));
  }

  updateExerciseScheme(id: number, data: Partial<ExerciseScheme>): Promise<ExerciseScheme> {
    return firstValueFrom(this.http.put<ExerciseScheme>(`/api/user/exercise-schemes/${id}`, data));
  }

  deleteExerciseScheme(id: number): Promise<void> {
    return firstValueFrom(this.http.delete<void>(`/api/user/exercise-schemes/${id}`));
  }

  // Workout Logs
  fetchWorkoutLogs(params?: {
    workoutId?: number;
    status?: string;
    periodId?: number;
  }): Promise<WorkoutLog[]> {
    const qp = new URLSearchParams();
    if (params?.workoutId != null) qp.set('workoutId', String(params.workoutId));
    if (params?.status) qp.set('status', params.status);
    if (params?.periodId != null) qp.set('periodId', String(params.periodId));
    const qs = qp.toString();
    return firstValueFrom(
      this.http.get<WorkoutLog[]>(`/api/user/workout-logs${qs ? '?' + qs : ''}`),
    );
  }

  fetchWorkoutLog(id: number): Promise<WorkoutLog> {
    return firstValueFrom(this.http.get<WorkoutLog>(`/api/user/workout-logs/${id}`));
  }

  createWorkoutLog(data: Partial<WorkoutLog>): Promise<WorkoutLog> {
    return firstValueFrom(this.http.post<WorkoutLog>('/api/user/workout-logs', data));
  }

  updateWorkoutLog(id: number, data: Partial<WorkoutLog>): Promise<WorkoutLog> {
    return firstValueFrom(this.http.patch<WorkoutLog>(`/api/user/workout-logs/${id}`, data));
  }

  createWorkoutLogSection(data: Partial<WorkoutLogSection>): Promise<WorkoutLogSection> {
    return firstValueFrom(
      this.http.post<WorkoutLogSection>('/api/user/workout-log-sections', data),
    );
  }

  updateWorkoutLogSection(
    id: number,
    data: Partial<WorkoutLogSection>,
  ): Promise<WorkoutLogSection> {
    return firstValueFrom(
      this.http.patch<WorkoutLogSection>(`/api/user/workout-log-sections/${id}`, data),
    );
  }

  deleteWorkoutLogSection(id: number): Promise<void> {
    return firstValueFrom(this.http.delete<void>(`/api/user/workout-log-sections/${id}`));
  }

  createWorkoutLogExercise(data: Partial<WorkoutLogExercise>): Promise<WorkoutLogExercise> {
    return firstValueFrom(
      this.http.post<WorkoutLogExercise>('/api/user/workout-log-exercises', data),
    );
  }

  updateWorkoutLogExercise(
    id: number,
    data: Partial<WorkoutLogExercise>,
  ): Promise<WorkoutLogExercise> {
    return firstValueFrom(
      this.http.patch<WorkoutLogExercise>(`/api/user/workout-log-exercises/${id}`, data),
    );
  }

  deleteWorkoutLogExercise(id: number): Promise<void> {
    return firstValueFrom(this.http.delete<void>(`/api/user/workout-log-exercises/${id}`));
  }

  updateWorkoutLogExerciseSet(
    id: number,
    data: Partial<WorkoutLogExerciseSet> & {
      actualReps?: number;
      actualWeight?: number;
      actualDuration?: number;
      actualDistance?: number;
      actualTime?: number;
    },
  ): Promise<WorkoutLogExerciseSet> {
    return firstValueFrom(
      this.http.patch<WorkoutLogExerciseSet>(`/api/user/workout-log-exercise-sets/${id}`, data),
    );
  }

  // Exercise Logs
  fetchExerciseLogs(params?: {
    exerciseId?: number;
    measurementType?: string;
    isRecord?: boolean;
    from?: string;
    to?: string;
  }): Promise<ExerciseLog[]> {
    const qp = new URLSearchParams();
    if (params?.exerciseId != null) qp.set('exerciseId', String(params.exerciseId));
    if (params?.measurementType) qp.set('measurementType', params.measurementType);
    if (params?.isRecord != null) qp.set('isRecord', String(params.isRecord));
    if (params?.from) qp.set('from', params.from);
    if (params?.to) qp.set('to', params.to);
    const qs = qp.toString();
    return firstValueFrom(
      this.http.get<ExerciseLog[]>(`/api/user/exercise-logs${qs ? '?' + qs : ''}`),
    );
  }

  fetchExerciseLog(id: number): Promise<ExerciseLog> {
    return firstValueFrom(this.http.get<ExerciseLog>(`/api/user/exercise-logs/${id}`));
  }

  createExerciseLog(data: Partial<ExerciseLog>): Promise<ExerciseLog> {
    return firstValueFrom(this.http.post<ExerciseLog>('/api/user/exercise-logs', data));
  }

  updateExerciseLog(id: number, data: Partial<ExerciseLog>): Promise<ExerciseLog> {
    return firstValueFrom(this.http.patch<ExerciseLog>(`/api/user/exercise-logs/${id}`, data));
  }

  deleteExerciseLog(id: number): Promise<void> {
    return firstValueFrom(this.http.delete<void>(`/api/user/exercise-logs/${id}`));
  }

  startWorkoutLog(id: number): Promise<WorkoutLog> {
    return firstValueFrom(this.http.post<WorkoutLog>(`/api/user/workout-logs/${id}/start`, {}));
  }

  abandonWorkoutLog(id: number): Promise<WorkoutLog> {
    return firstValueFrom(this.http.post<WorkoutLog>(`/api/user/workout-logs/${id}/abandon`, {}));
  }

  startAdhocWorkoutLog(): Promise<WorkoutLog> {
    return firstValueFrom(this.http.post<WorkoutLog>('/api/user/workout-logs/adhoc', {}));
  }

  finishWorkoutLog(id: number): Promise<WorkoutLog> {
    return firstValueFrom(this.http.post<WorkoutLog>(`/api/user/workout-logs/${id}/finish`, {}));
  }

  skipWorkoutLog(id: number): Promise<WorkoutLog> {
    return firstValueFrom(this.http.post<WorkoutLog>(`/api/user/workout-logs/${id}/skip`, {}));
  }

  commitWorkoutLog(id: number): Promise<WorkoutLog> {
    return firstValueFrom(this.http.post<WorkoutLog>(`/api/user/workout-logs/${id}/commit`, {}));
  }

  // Workout Groups
  fetchWorkoutGroups(): Promise<WorkoutGroup[]> {
    return firstValueFrom(this.http.get<WorkoutGroup[]>('/api/user/workout-groups'));
  }

  fetchWorkoutGroup(id: number): Promise<WorkoutGroup> {
    return firstValueFrom(this.http.get<WorkoutGroup>(`/api/user/workout-groups/${id}`));
  }

  createWorkoutGroup(data: { name: string; workoutId: number }): Promise<WorkoutGroup> {
    return firstValueFrom(this.http.post<WorkoutGroup>('/api/user/workout-groups', data));
  }

  updateWorkoutGroup(id: number, data: { name: string }): Promise<WorkoutGroup> {
    return firstValueFrom(this.http.put<WorkoutGroup>(`/api/user/workout-groups/${id}`, data));
  }

  deleteWorkoutGroup(id: number): Promise<void> {
    return firstValueFrom(this.http.delete<void>(`/api/user/workout-groups/${id}`));
  }

  acceptWorkoutGroupInvitation(workoutId: number): Promise<WorkoutGroupMembership> {
    return firstValueFrom(
      this.http.post<WorkoutGroupMembership>(`/api/workouts/${workoutId}/group/accept`, {}),
    );
  }

  // Workout Group Memberships
  fetchWorkoutGroupMemberships(params?: { groupId?: number }): Promise<WorkoutGroupMembership[]> {
    const qp = new URLSearchParams();
    if (params?.groupId != null) qp.set('groupId', String(params.groupId));
    const qs = qp.toString();
    return firstValueFrom(
      this.http.get<WorkoutGroupMembership[]>(
        `/api/user/workout-group-memberships${qs ? '?' + qs : ''}`,
      ),
    );
  }

  createWorkoutGroupMembership(data: {
    groupId: number;
    userId: string;
    role: WorkoutGroupRole;
  }): Promise<WorkoutGroupMembership> {
    return firstValueFrom(
      this.http.post<WorkoutGroupMembership>('/api/user/workout-group-memberships', data),
    );
  }

  updateWorkoutGroupMembership(
    id: number,
    data: { role: WorkoutGroupRole },
  ): Promise<WorkoutGroupMembership> {
    return firstValueFrom(
      this.http.put<WorkoutGroupMembership>(`/api/user/workout-group-memberships/${id}`, data),
    );
  }

  deleteWorkoutGroupMembership(id: number): Promise<void> {
    return firstValueFrom(this.http.delete<void>(`/api/user/workout-group-memberships/${id}`));
  }

  // Workout Schedules
  fetchWorkoutSchedules(params?: { workoutId?: number }): Promise<WorkoutSchedule[]> {
    const qp = new URLSearchParams();
    if (params?.workoutId != null) qp.set('workoutId', String(params.workoutId));
    const qs = qp.toString();
    return firstValueFrom(
      this.http.get<WorkoutSchedule[]>(`/api/user/workout-schedules${qs ? '?' + qs : ''}`),
    );
  }

  fetchWorkoutSchedule(id: number): Promise<WorkoutSchedule> {
    return firstValueFrom(this.http.get<WorkoutSchedule>(`/api/user/workout-schedules/${id}`));
  }

  createWorkoutSchedule(data: Partial<WorkoutSchedule>): Promise<WorkoutSchedule> {
    return firstValueFrom(this.http.post<WorkoutSchedule>('/api/user/workout-schedules', data));
  }

  updateWorkoutSchedule(id: number, data: Partial<WorkoutSchedule>): Promise<WorkoutSchedule> {
    return firstValueFrom(
      this.http.patch<WorkoutSchedule>(`/api/user/workout-schedules/${id}`, data),
    );
  }

  deleteWorkoutSchedule(id: number): Promise<void> {
    return firstValueFrom(this.http.delete<void>(`/api/user/workout-schedules/${id}`));
  }

  // Schedule Periods
  createSchedulePeriod(data: Partial<SchedulePeriod>): Promise<SchedulePeriod> {
    return firstValueFrom(this.http.post<SchedulePeriod>('/api/user/schedule-periods', data));
  }

  // Schedule Commitments
  createScheduleCommitment(data: Partial<ScheduleCommitment>): Promise<ScheduleCommitment> {
    return firstValueFrom(
      this.http.post<ScheduleCommitment>('/api/user/schedule-commitments', data),
    );
  }

  fetchScheduleCommitments(params?: { periodId?: number }): Promise<ScheduleCommitment[]> {
    const qp = new URLSearchParams();
    if (params?.periodId != null) qp.set('periodId', String(params.periodId));
    const qs = qp.toString();
    return firstValueFrom(
      this.http.get<ScheduleCommitment[]>(`/api/user/schedule-commitments${qs ? '?' + qs : ''}`),
    );
  }

  deleteScheduleCommitment(id: number): Promise<void> {
    return firstValueFrom(this.http.delete<void>(`/api/user/schedule-commitments/${id}`));
  }

  fetchSchedulePeriods(params?: { scheduleId?: number }): Promise<SchedulePeriod[]> {
    const qp = new URLSearchParams();
    if (params?.scheduleId != null) qp.set('scheduleId', String(params.scheduleId));
    const qs = qp.toString();
    return firstValueFrom(
      this.http.get<SchedulePeriod[]>(`/api/user/schedule-periods${qs ? '?' + qs : ''}`),
    );
  }

  // Exercise Name Preferences
  fetchExerciseNamePreferences(): Promise<{ exerciseId: number; exerciseNameId: number }[]> {
    return firstValueFrom(
      this.http.get<{ exerciseId: number; exerciseNameId: number }[]>(
        '/api/user/exercise-name-preferences',
      ),
    );
  }

  setExerciseNamePreference(exerciseId: number, exerciseNameId: number): Promise<void> {
    return firstValueFrom(
      this.http.put<void>(`/api/user/exercise-name-preferences/${exerciseId}`, { exerciseNameId }),
    );
  }
}
