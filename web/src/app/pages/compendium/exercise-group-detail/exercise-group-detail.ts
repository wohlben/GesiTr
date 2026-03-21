import { Component, inject, computed, signal } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import {
  injectQuery,
  injectMutation,
  injectQueryClient,
} from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { exerciseGroupKeys } from '$core/query-keys';
import { PageLayout } from '../../../layout/page-layout';
import { ConfirmDialog } from '$ui/confirm-dialog/confirm-dialog';

@Component({
  selector: 'app-exercise-group-detail',
  imports: [PageLayout, RouterLink, ConfirmDialog],
  template: `
    <app-page-layout
      [header]="groupQuery.data()?.name ?? 'Exercise Group'"
      [isPending]="groupQuery.isPending()"
      [errorMessage]="groupQuery.isError() ? groupQuery.error().message : undefined"
    >
      <div actions class="flex gap-2">
        <a
          routerLink="./edit"
          class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
          >Edit</a
        >
        <button
          type="button"
          (click)="showDeleteDialog.set(true)"
          class="rounded-md bg-red-600 px-4 py-2 text-sm font-medium text-white hover:bg-red-700"
        >
          Delete
        </button>
      </div>
      <app-confirm-dialog
        [open]="showDeleteDialog()"
        title="Delete Exercise Group"
        [message]="deleteMessage()"
        [isPending]="deleteMutation.isPending()"
        (confirmed)="deleteMutation.mutate()"
        (cancelled)="showDeleteDialog.set(false)"
      />
      @if (groupQuery.data(); as group) {
        <dl class="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div class="sm:col-span-2">
            <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Description</dt>
            <dd class="text-sm text-gray-900 dark:text-gray-100">{{ group.description }}</dd>
          </div>
          <div>
            <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Created By</dt>
            <dd class="text-sm text-gray-900 dark:text-gray-100">{{ group.createdBy }}</dd>
          </div>
        </dl>
      }
    </app-page-layout>
  `,
})
export class ExerciseGroupDetail {
  private api = inject(CompendiumApiClient);
  private router = inject(Router);
  private queryClient = injectQueryClient();
  private params = toSignal(inject(ActivatedRoute).paramMap);

  private id = computed(() => Number(this.params()?.get('id')));

  showDeleteDialog = signal(false);
  deleteMessage = computed(
    () => `Are you sure you want to delete '${this.groupQuery.data()?.name ?? ''}'?`,
  );

  groupQuery = injectQuery(() => ({
    queryKey: exerciseGroupKeys.detail(this.id()),
    queryFn: () => this.api.fetchExerciseGroup(this.id()),
    enabled: !!this.id(),
  }));

  deleteMutation = injectMutation(() => ({
    mutationFn: () => this.api.deleteExerciseGroup(this.id()),
    onSuccess: () => {
      this.queryClient.invalidateQueries({ queryKey: exerciseGroupKeys.all() });
      this.router.navigate(['/compendium/exercise-groups']);
    },
  }));
}
