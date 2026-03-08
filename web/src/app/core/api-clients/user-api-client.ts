import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { firstValueFrom } from 'rxjs';
import { UserExercise, UserEquipment } from '$generated/user-models';

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
}
