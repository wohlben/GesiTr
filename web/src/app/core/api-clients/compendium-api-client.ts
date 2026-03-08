import { HttpClient, HttpParams } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { firstValueFrom } from 'rxjs';
import { Exercise, Equipment, ExerciseGroup } from '$generated/models';
import { PaginatedResponse } from './paginated-response';
import { VersionEntry } from './version-entry';

function buildParams(filters: Record<string, string | number | undefined>): HttpParams {
  let params = new HttpParams();
  for (const [key, value] of Object.entries(filters)) {
    if (value !== undefined && value !== '') params = params.set(key, String(value));
  }
  return params;
}

@Injectable({ providedIn: 'root' })
export class CompendiumApiClient {
  private http = inject(HttpClient);

  fetchExercises(
    filters: Record<string, string | number | undefined>,
  ): Promise<PaginatedResponse<Exercise>> {
    return firstValueFrom(
      this.http.get<PaginatedResponse<Exercise>>('/api/exercises', {
        params: buildParams(filters),
      }),
    );
  }

  fetchEquipment(
    filters: Record<string, string | number | undefined>,
  ): Promise<PaginatedResponse<Equipment>> {
    return firstValueFrom(
      this.http.get<PaginatedResponse<Equipment>>('/api/equipment', {
        params: buildParams(filters),
      }),
    );
  }

  fetchExerciseGroups(
    filters: Record<string, string | number | undefined>,
  ): Promise<PaginatedResponse<ExerciseGroup>> {
    return firstValueFrom(
      this.http.get<PaginatedResponse<ExerciseGroup>>('/api/exercise-groups', {
        params: buildParams(filters),
      }),
    );
  }

  fetchExercise(id: number): Promise<Exercise> {
    return firstValueFrom(this.http.get<Exercise>(`/api/exercises/${id}`));
  }

  fetchEquipmentItem(id: number): Promise<Equipment> {
    return firstValueFrom(this.http.get<Equipment>(`/api/equipment/${id}`));
  }

  fetchExerciseGroup(id: number): Promise<ExerciseGroup> {
    return firstValueFrom(this.http.get<ExerciseGroup>(`/api/exercise-groups/${id}`));
  }

  updateExercise(id: number, data: Partial<Exercise>): Promise<Exercise> {
    return firstValueFrom(this.http.put<Exercise>(`/api/exercises/${id}`, data));
  }

  updateEquipment(id: number, data: Partial<Equipment>): Promise<Equipment> {
    return firstValueFrom(this.http.put<Equipment>(`/api/equipment/${id}`, data));
  }

  updateExerciseGroup(id: number, data: Partial<ExerciseGroup>): Promise<ExerciseGroup> {
    return firstValueFrom(this.http.put<ExerciseGroup>(`/api/exercise-groups/${id}`, data));
  }

  fetchExerciseVersions(id: number): Promise<VersionEntry<Exercise>[]> {
    return firstValueFrom(this.http.get<VersionEntry<Exercise>[]>(`/api/exercises/${id}/versions`));
  }

  fetchEquipmentVersions(id: number): Promise<VersionEntry<Equipment>[]> {
    return firstValueFrom(
      this.http.get<VersionEntry<Equipment>[]>(`/api/equipment/${id}/versions`),
    );
  }

  deleteExercise(id: number): Promise<void> {
    return firstValueFrom(this.http.delete<void>(`/api/exercises/${id}`));
  }

  deleteEquipment(id: number): Promise<void> {
    return firstValueFrom(this.http.delete<void>(`/api/equipment/${id}`));
  }

  deleteExerciseGroup(id: number): Promise<void> {
    return firstValueFrom(this.http.delete<void>(`/api/exercise-groups/${id}`));
  }

  createExercise(data: Partial<Exercise>): Promise<Exercise> {
    return firstValueFrom(this.http.post<Exercise>('/api/exercises', data));
  }

  createEquipment(data: Partial<Equipment>): Promise<Equipment> {
    return firstValueFrom(this.http.post<Equipment>('/api/equipment', data));
  }

  createExerciseGroup(data: Partial<ExerciseGroup>): Promise<ExerciseGroup> {
    return firstValueFrom(this.http.post<ExerciseGroup>('/api/exercise-groups', data));
  }
}
