import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { firstValueFrom } from 'rxjs';
import {
  UserExercise,
  UserEquipment,
  Workout,
  WorkoutSection,
  WorkoutSectionExercise,
  UserExerciseScheme,
  WorkoutLog,
  WorkoutLogSection,
  WorkoutLogExercise,
  WorkoutLogExerciseSet,
  ExerciseLog,
} from '$generated/user-models';

@Injectable({ providedIn: 'root' })
export class UserApiClient {
  private http = inject(HttpClient);

  fetchUserExercises(): Promise<UserExercise[]> {
    return firstValueFrom(this.http.get<UserExercise[]>('/api/user/exercises'));
  }

  fetchUserExercise(id: number): Promise<UserExercise> {
    return firstValueFrom(this.http.get<UserExercise>(`/api/user/exercises/${id}`));
  }

  createUserExercise(data: Partial<UserExercise>): Promise<UserExercise> {
    return firstValueFrom(this.http.post<UserExercise>('/api/user/exercises', data));
  }

  deleteUserExercise(id: number): Promise<void> {
    return firstValueFrom(this.http.delete<void>(`/api/user/exercises/${id}`));
  }

  fetchUserEquipment(): Promise<UserEquipment[]> {
    return firstValueFrom(this.http.get<UserEquipment[]>('/api/user/equipment'));
  }

  fetchUserEquipmentItem(id: number): Promise<UserEquipment> {
    return firstValueFrom(this.http.get<UserEquipment>(`/api/user/equipment/${id}`));
  }

  createUserEquipment(data: Partial<UserEquipment>): Promise<UserEquipment> {
    return firstValueFrom(this.http.post<UserEquipment>('/api/user/equipment', data));
  }

  deleteUserEquipment(id: number): Promise<void> {
    return firstValueFrom(this.http.delete<void>(`/api/user/equipment/${id}`));
  }

  // Workouts
  fetchWorkouts(): Promise<Workout[]> {
    return firstValueFrom(this.http.get<Workout[]>('/api/user/workouts'));
  }

  fetchWorkout(id: number): Promise<Workout> {
    return firstValueFrom(this.http.get<Workout>(`/api/user/workouts/${id}`));
  }

  createWorkout(data: Partial<Workout>): Promise<Workout> {
    return firstValueFrom(this.http.post<Workout>('/api/user/workouts', data));
  }

  updateWorkout(id: number, data: Partial<Workout>): Promise<Workout> {
    return firstValueFrom(this.http.put<Workout>(`/api/user/workouts/${id}`, data));
  }

  deleteWorkout(id: number): Promise<void> {
    return firstValueFrom(this.http.delete<void>(`/api/user/workouts/${id}`));
  }

  // Workout Sections
  createWorkoutSection(data: Partial<WorkoutSection>): Promise<WorkoutSection> {
    return firstValueFrom(this.http.post<WorkoutSection>('/api/user/workout-sections', data));
  }

  deleteWorkoutSection(id: number): Promise<void> {
    return firstValueFrom(this.http.delete<void>(`/api/user/workout-sections/${id}`));
  }

  // Workout Section Exercises
  createWorkoutSectionExercise(
    data: Partial<WorkoutSectionExercise>,
  ): Promise<WorkoutSectionExercise> {
    return firstValueFrom(
      this.http.post<WorkoutSectionExercise>('/api/user/workout-section-exercises', data),
    );
  }

  deleteWorkoutSectionExercise(id: number): Promise<void> {
    return firstValueFrom(this.http.delete<void>(`/api/user/workout-section-exercises/${id}`));
  }

  // Exercise Schemes
  fetchExerciseSchemes(params?: { userExerciseId?: number }): Promise<UserExerciseScheme[]> {
    const qp = new URLSearchParams();
    if (params?.userExerciseId != null) qp.set('userExerciseId', String(params.userExerciseId));
    const qs = qp.toString();
    return firstValueFrom(
      this.http.get<UserExerciseScheme[]>(`/api/user/exercise-schemes${qs ? '?' + qs : ''}`),
    );
  }

  fetchExerciseScheme(id: number): Promise<UserExerciseScheme> {
    return firstValueFrom(this.http.get<UserExerciseScheme>(`/api/user/exercise-schemes/${id}`));
  }

  createExerciseScheme(data: Partial<UserExerciseScheme>): Promise<UserExerciseScheme> {
    return firstValueFrom(this.http.post<UserExerciseScheme>('/api/user/exercise-schemes', data));
  }

  updateExerciseScheme(id: number, data: Partial<UserExerciseScheme>): Promise<UserExerciseScheme> {
    return firstValueFrom(
      this.http.put<UserExerciseScheme>(`/api/user/exercise-schemes/${id}`, data),
    );
  }

  deleteExerciseScheme(id: number): Promise<void> {
    return firstValueFrom(this.http.delete<void>(`/api/user/exercise-schemes/${id}`));
  }

  // Workout Logs
  fetchWorkoutLogs(params?: { workoutId?: number; status?: string }): Promise<WorkoutLog[]> {
    const qp = new URLSearchParams();
    if (params?.workoutId != null) qp.set('workoutId', String(params.workoutId));
    if (params?.status) qp.set('status', params.status);
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
    userExerciseId?: number;
    measurementType?: string;
    isRecord?: boolean;
    from?: string;
    to?: string;
  }): Promise<ExerciseLog[]> {
    const qp = new URLSearchParams();
    if (params?.userExerciseId != null) qp.set('userExerciseId', String(params.userExerciseId));
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
}
