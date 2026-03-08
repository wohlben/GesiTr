import { HttpClient, HttpParams } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { firstValueFrom } from 'rxjs';
import { Exercise, Equipment, ExerciseGroup } from '$generated/models';
import { PaginatedResponse } from './paginated-response';

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
}
