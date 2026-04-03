import { HttpClient, HttpParams } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { firstValueFrom } from 'rxjs';
import {
  Exercise,
  Equipment,
  ExerciseRelationship,
  EquipmentRelationship,
} from '$generated/models';
import { ExerciseGroup, ExerciseGroupMember } from '$generated/user-models';
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

  fetchExerciseVersion(id: number, version: number): Promise<VersionEntry<Exercise>> {
    return firstValueFrom(
      this.http.get<VersionEntry<Exercise>>(`/api/exercises/${id}/versions/${version}`),
    );
  }

  fetchEquipmentVersion(id: number, version: number): Promise<VersionEntry<Equipment>> {
    return firstValueFrom(
      this.http.get<VersionEntry<Equipment>>(`/api/equipment/${id}/versions/${version}`),
    );
  }

  fetchExercisePermissions(id: number): Promise<{ permissions: string[] }> {
    return firstValueFrom(
      this.http.get<{ permissions: string[] }>(`/api/exercises/${id}/permissions`),
    );
  }

  fetchEquipmentPermissions(id: number): Promise<{ permissions: string[] }> {
    return firstValueFrom(
      this.http.get<{ permissions: string[] }>(`/api/equipment/${id}/permissions`),
    );
  }

  fetchExerciseGroupPermissions(id: number): Promise<{ permissions: string[] }> {
    return firstValueFrom(
      this.http.get<{ permissions: string[] }>(`/api/exercise-groups/${id}/permissions`),
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

  deleteExerciseVersion(id: number, version: number): Promise<void> {
    return firstValueFrom(this.http.delete<void>(`/api/exercises/${id}/versions/${version}`));
  }

  deleteAllExerciseVersions(id: number): Promise<void> {
    return firstValueFrom(this.http.delete<void>(`/api/exercises/${id}/versions`));
  }

  fetchExerciseGroupMembers(
    filters: Record<string, string | number | undefined>,
  ): Promise<ExerciseGroupMember[]> {
    return firstValueFrom(
      this.http.get<ExerciseGroupMember[]>('/api/exercise-group-members', {
        params: buildParams(filters),
      }),
    );
  }

  createExerciseGroupMember(data: {
    groupId: number;
    exerciseId: number;
  }): Promise<ExerciseGroupMember> {
    return firstValueFrom(this.http.post<ExerciseGroupMember>('/api/exercise-group-members', data));
  }

  deleteExerciseGroupMember(id: number): Promise<void> {
    return firstValueFrom(this.http.delete<void>(`/api/exercise-group-members/${id}`));
  }

  fetchExerciseRelationships(
    filters: Record<string, string | number | undefined>,
  ): Promise<ExerciseRelationship[]> {
    return firstValueFrom(
      this.http.get<ExerciseRelationship[]>('/api/exercise-relationships', {
        params: buildParams(filters),
      }),
    );
  }

  fetchEquipmentRelationships(
    filters: Record<string, string | number | undefined>,
  ): Promise<EquipmentRelationship[]> {
    return firstValueFrom(
      this.http.get<EquipmentRelationship[]>('/api/equipment-relationships', {
        params: buildParams(filters),
      }),
    );
  }

  fetchDeployStatus(): Promise<{ status: string; title?: string; createdAt?: string }> {
    return firstValueFrom(
      this.http.get<{ status: string; title?: string; createdAt?: string }>('/api/deploy-status'),
    );
  }
}
