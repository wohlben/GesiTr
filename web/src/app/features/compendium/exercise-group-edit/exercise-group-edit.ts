import { Component, inject, computed, effect } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { ReactiveFormsModule, FormGroup, FormControl, Validators } from '@angular/forms';
import { injectQuery, injectMutation, injectQueryClient } from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { exerciseGroupKeys } from '$core/query-keys';
import { PageLayout } from '../../../layout/page-layout';

@Component({
  selector: 'app-exercise-group-edit',
  imports: [PageLayout, ReactiveFormsModule, RouterLink],
  template: `
    <app-page-layout
      header="Edit Exercise Group"
      [isPending]="groupQuery.isPending()"
      [errorMessage]="groupQuery.isError() ? groupQuery.error().message : undefined"
    >
      @if (groupQuery.data(); as group) {
        <form [formGroup]="form" (ngSubmit)="onSubmit()" class="space-y-4">
          <div>
            <label for="name" class="block text-sm font-medium text-gray-700 dark:text-gray-300">Name *</label>
            <input
              id="name"
              formControlName="name"
              class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
            />
          </div>

          <div>
            <label for="description" class="block text-sm font-medium text-gray-700 dark:text-gray-300"
              >Description</label
            >
            <textarea
              id="description"
              formControlName="description"
              rows="4"
              class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
            ></textarea>
          </div>

          <div class="flex gap-2">
            <button
              type="submit"
              [disabled]="form.invalid || mutation.isPending()"
              class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
            >
              Save
            </button>
            <a [routerLink]="['..']" class="rounded-md border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-800">
              Cancel
            </a>
          </div>
        </form>
      }
    </app-page-layout>
  `,
})
export class ExerciseGroupEdit {
  private api = inject(CompendiumApiClient);
  private router = inject(Router);
  private route = inject(ActivatedRoute);
  private queryClient = injectQueryClient();
  private params = toSignal(this.route.paramMap);

  private id = computed(() => Number(this.params()?.get('id')));

  form = new FormGroup({
    name: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    description: new FormControl('', { nonNullable: true }),
  });

  groupQuery = injectQuery(() => ({
    queryKey: exerciseGroupKeys.detail(this.id()),
    queryFn: () => this.api.fetchExerciseGroup(this.id()),
    enabled: !!this.id(),
  }));

  mutation = injectMutation(() => ({
    mutationFn: (data: Parameters<typeof this.api.updateExerciseGroup>[1]) =>
      this.api.updateExerciseGroup(this.id(), data),
    onSuccess: () => {
      this.queryClient.invalidateQueries({ queryKey: exerciseGroupKeys.all() });
      this.router.navigate(['..'], { relativeTo: this.route });
    },
  }));

  constructor() {
    effect(() => {
      const data = this.groupQuery.data();
      if (data) {
        this.form.patchValue({
          name: data.name,
          description: data.description ?? '',
        });
      }
    });
  }

  onSubmit() {
    if (this.form.valid) {
      const val = this.form.getRawValue();
      this.mutation.mutate({
        ...this.groupQuery.data()!,
        name: val.name,
        description: val.description || undefined,
      });
    }
  }
}
