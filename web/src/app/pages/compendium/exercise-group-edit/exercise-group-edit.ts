import { Component, inject, computed, effect, signal } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { form, required, FormField } from '@angular/forms/signals';
import { injectQuery, injectMutation, QueryClient } from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { exerciseGroupKeys } from '$core/query-keys';
import { TranslocoDirective } from '@jsverse/transloco';
import { SlugifyPipe } from '$ui/pipes/slugify';
import { PageLayout } from '../../../layout/page-layout';
import { HlmInput } from '@spartan-ng/helm/input';
import { HlmTextarea } from '@spartan-ng/helm/textarea';

@Component({
  selector: 'app-exercise-group-edit',
  imports: [PageLayout, FormField, RouterLink, HlmInput, HlmTextarea, TranslocoDirective],
  template: `
    <ng-container *transloco="let t">
      <app-page-layout
        [header]="
          isCreateMode()
            ? t('compendium.exerciseGroups.newTitle')
            : t('compendium.exerciseGroups.editTitle')
        "
        [isPending]="!isCreateMode() && groupQuery.isPending()"
        [errorMessage]="
          !isCreateMode() && groupQuery.isError() ? groupQuery.error().message : undefined
        "
      >
        @if (isCreateMode() || groupQuery.data()) {
          <form (submit)="onSubmit(); $event.preventDefault()" class="space-y-4">
            <div>
              <label for="name" class="block text-sm font-medium text-gray-700 dark:text-gray-300"
                >{{ t('fields.name') }} *</label
              >
              <input id="name" [formField]="groupForm.name" hlmInput class="mt-1" />
            </div>

            <div>
              <label
                for="description"
                class="block text-sm font-medium text-gray-700 dark:text-gray-300"
                >{{ t('fields.description') }}</label
              >
              <textarea
                id="description"
                [formField]="groupForm.description"
                rows="4"
                hlmTextarea
                class="mt-1"
              ></textarea>
            </div>

            <div class="flex gap-2">
              <button
                type="submit"
                [disabled]="
                  !groupForm().valid() || mutation.isPending() || createMutation.isPending()
                "
                class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
              >
                {{ t('common.save') }}
              </button>
              <a
                [routerLink]="isCreateMode() ? ['/compendium/exercise-groups'] : ['..']"
                class="rounded-md border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-800"
              >
                {{ t('common.cancel') }}
              </a>
            </div>
          </form>
        }
      </app-page-layout>
    </ng-container>
  `,
})
export class ExerciseGroupEdit {
  private api = inject(CompendiumApiClient);
  private router = inject(Router);
  private route = inject(ActivatedRoute);
  private queryClient = inject(QueryClient);
  private slugify = new SlugifyPipe();
  private params = toSignal(this.route.paramMap);

  private id = computed(() => Number(this.params()?.get('id')));
  isCreateMode = computed(() => !this.params()?.get('id'));

  model = signal({ name: '', description: '' });
  groupForm = form(this.model, (f) => {
    required(f.name);
  });

  groupQuery = injectQuery(() => ({
    queryKey: exerciseGroupKeys.detail(this.id()),
    queryFn: () => this.api.fetchExerciseGroup(this.id()),
    enabled: !!this.id() && !this.isCreateMode(),
  }));

  permissionsQuery = injectQuery(() => ({
    queryKey: exerciseGroupKeys.permissions(this.id()),
    queryFn: () => this.api.fetchExerciseGroupPermissions(this.id()),
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
      const perms = this.permissionsQuery.data();
      if (perms && !perms.permissions.includes('MODIFY')) {
        this.router.navigate(['..'], { relativeTo: this.route });
      }
    });

    effect(() => {
      const data = this.groupQuery.data();
      if (data) {
        this.model.set({
          name: data.name,
          description: data.description ?? '',
        });
      }
    });
  }

  onSubmit() {
    if (this.groupForm().valid()) {
      const val = this.model();
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
