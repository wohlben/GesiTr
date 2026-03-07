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

  fetchExercises(filters: {
    q?: string;
    type?: string;
    difficulty?: string;
    force?: string;
    muscle?: string;
    limit?: number;
    offset?: number;
  }): Promise<PaginatedResponse<Exercise>> {
    return firstValueFrom(
      this.http.get<PaginatedResponse<Exercise>>('/api/exercises', { params: buildParams(filters) }),
    );
  }

  fetchEquipment(filters: {
    q?: string;
    category?: string;
    limit?: number;
    offset?: number;
  }): Promise<PaginatedResponse<Equipment>> {
    return firstValueFrom(
      this.http.get<PaginatedResponse<Equipment>>('/api/equipment', { params: buildParams(filters) }),
    );
  }

  fetchExerciseGroups(filters: {
    q?: string;
    limit?: number;
    offset?: number;
  }): Promise<PaginatedResponse<ExerciseGroup>> {
    return firstValueFrom(
      this.http.get<PaginatedResponse<ExerciseGroup>>('/api/exercise-groups', { params: buildParams(filters) }),
    );
  }
}
