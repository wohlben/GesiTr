import { HttpClient, HttpParams } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { firstValueFrom } from 'rxjs';
import { Exercise, Equipment, ExerciseGroup } from '$generated/models';

function buildParams(filters: Record<string, string | undefined>): HttpParams {
  let params = new HttpParams();
  for (const [key, value] of Object.entries(filters)) {
    if (value) params = params.set(key, value);
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
  }): Promise<Exercise[]> {
    return firstValueFrom(
      this.http.get<Exercise[]>('/api/exercises', { params: buildParams(filters) }),
    );
  }

  fetchEquipment(filters: {
    q?: string;
    category?: string;
  }): Promise<Equipment[]> {
    return firstValueFrom(
      this.http.get<Equipment[]>('/api/equipment', { params: buildParams(filters) }),
    );
  }

  fetchExerciseGroups(filters: {
    q?: string;
  }): Promise<ExerciseGroup[]> {
    return firstValueFrom(
      this.http.get<ExerciseGroup[]>('/api/exercise-groups', { params: buildParams(filters) }),
    );
  }
}
