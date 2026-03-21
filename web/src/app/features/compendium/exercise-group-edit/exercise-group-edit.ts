import { Component, inject, computed, effect } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { ReactiveFormsModule, FormGroup, FormControl, Validators } from '@angular/forms';
import {
  injectQuery,
  injectMutation,
  injectQueryClient,
} from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { exerciseGroupKeys } from '$core/query-keys';
import { SlugifyPipe } from '$ui/pipes/slugify';
import { PageLayout } from '../../../layout/page-layout';
import { HlmInput } from '@spartan-ng/helm/input';
import { HlmTextarea } from '@spartan-ng/helm/textarea';

@Component({
  selector: 'app-exercise-group-edit',
  imports: [PageLayout, ReactiveFormsModule, RouterLink, HlmInput, HlmTextarea],
  template: `
    <app-page-layout
      [header]="isCreateMode() ? 'New Exercise Group' : 'Edit Exercise Group'"
      [isPending]="!isCreateMode() && groupQuery.isPending()"
      [errorMessage]="
        !isCreateMode() && groupQuery.isError() ? groupQuery.error().message : undefined
      "
    >
      @if (isCreateMode() || groupQuery.data()) {
        <form [formGroup]="form" (ngSubmit)="onSubmit()" class="space-y-4">
          <div>
            <label for="name" class="block text-sm font-medium text-gray-700 dark:text-gray-300"
              >Name *</label
            >
            <input id="name" formControlName="name" hlmInput class="mt-1" />
          </div>

          <div>
            <label
              for="description"
              class="block text-sm font-medium text-gray-700 dark:text-gray-300"
              >Description</label
            >
            <textarea
              id="description"
              formControlName="description"
              rows="4"
              hlmTextarea
              class="mt-1"
            ></textarea>
          </div>

          <div class="flex gap-2">
            <button
              type="submit"
              [disabled]="form.invalid || mutation.isPending() || createMutation.isPending()"
              class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
            >
              Save
            </button>
            <a
              [routerLink]="isCreateMode() ? ['/compendium/exercise-groups'] : ['..']"
              class="rounded-md border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-800"
            >
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
  private slugify = new SlugifyPipe();
  private params = toSignal(this.route.paramMap);

  private id = computed(() => Number(this.params()?.get('id')));
  isCreateMode = computed(() => !this.params()?.get('id'));

  form = new FormGroup({
    name: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    description: new FormControl('', { nonNullable: true }),
  });

  groupQuery = injectQuery(() => ({
    queryKey: exerciseGroupKeys.detail(this.id()),
    queryFn: () => this.api.fetchExerciseGroup(this.id()),
    enabled: !!this.id() && !this.isCreateMode(),
  }));

  mutation = injectMutation(() => ({
    mutationFn: (data: Parameters<typeof this.api.updateExerciseGroup>[1]) =>
      this.api.updateExerciseGroup(this.id(), data),
    onSuccess: () => {
      this.queryClient.invalidateQueries({ queryKey: exerciseGroupKeys.all() });
      this.router.navigate(['..'], { relativeTo: this.route });
    },
  }));

  createMutation = injectMutation(() => ({
    mutationFn: (data: Parameters<typeof this.api.createExerciseGroup>[0]) =>
      this.api.createExerciseGroup(data),
    onSuccess: (result) => {
      this.queryClient.invalidateQueries({ queryKey: exerciseGroupKeys.all() });
      this.router.navigate([
        '/compendium/exercise-groups',
        result.id,
        this.slugify.transform(result.name),
      ]);
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
      const payload = {
        ...(this.isCreateMode() ? {} : this.groupQuery.data()!),
        name: val.name,

        description: val.description || undefined,
      };
      if (this.isCreateMode()) {
        this.createMutation.mutate(payload);
      } else {
        this.mutation.mutate(payload);
      }
    }
  }
}
