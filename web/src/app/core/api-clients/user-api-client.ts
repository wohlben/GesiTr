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
  fetchWorkoutLog(id: number): Promise<WorkoutLog> {
    return firstValueFrom(this.http.get<WorkoutLog>(`/api/user/workout-logs/${id}`));
  }

  createWorkoutLog(data: Partial<WorkoutLog>): Promise<WorkoutLog> {
    return firstValueFrom(this.http.post<WorkoutLog>('/api/user/workout-logs', data));
  }

  createWorkoutLogSection(data: Partial<WorkoutLogSection>): Promise<WorkoutLogSection> {
    return firstValueFrom(
      this.http.post<WorkoutLogSection>('/api/user/workout-log-sections', data),
    );
  }

  createWorkoutLogExercise(data: Partial<WorkoutLogExercise>): Promise<WorkoutLogExercise> {
    return firstValueFrom(
      this.http.post<WorkoutLogExercise>('/api/user/workout-log-exercises', data),
    );
  }
}
