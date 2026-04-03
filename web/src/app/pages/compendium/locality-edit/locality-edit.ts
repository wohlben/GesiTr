import { Component, inject, computed, effect, signal } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { form, required, FormField } from '@angular/forms/signals';
import { injectQuery, injectMutation, QueryClient } from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { localityKeys } from '$core/query-keys';
import { TranslocoDirective } from '@jsverse/transloco';
import { SlugifyPipe } from '$ui/pipes/slugify';
import { PageLayout } from '../../../layout/page-layout';
import { HlmInput } from '@spartan-ng/helm/input';

@Component({
  selector: 'app-locality-edit',
  imports: [PageLayout, FormField, RouterLink, HlmInput, TranslocoDirective],
  template: `
    <ng-container *transloco="let t">
      <app-page-layout
        [header]="
          isCreateMode()
            ? t('compendium.localities.newTitle')
            : t('compendium.localities.editTitle')
        "
        [isPending]="!isCreateMode() && localityQuery.isPending()"
        [errorMessage]="
          !isCreateMode() && localityQuery.isError() ? localityQuery.error().message : undefined
        "
      >
        @if (isCreateMode() || localityQuery.data()) {
          <form (submit)="onSubmit(); $event.preventDefault()" class="space-y-4">
            <div>
              <label for="name" class="block text-sm font-medium text-gray-700 dark:text-gray-300"
                >{{ t('fields.name') }} *</label
              >
              <input id="name" [formField]="localityForm.name" hlmInput class="mt-1" />
            </div>

            <div class="flex gap-2">
              <button
                type="submit"
                [disabled]="
                  !localityForm().valid() || mutation.isPending() || createMutation.isPending()
                "
                class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
              >
                {{ t('common.save') }}
              </button>
              <a
                [routerLink]="isCreateMode() ? ['/compendium/localities'] : ['..']"
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
export class LocalityEdit {
  private api = inject(CompendiumApiClient);
  private router = inject(Router);
  private route = inject(ActivatedRoute);
  private queryClient = inject(QueryClient);
  private slugify = new SlugifyPipe();
  private params = toSignal(this.route.paramMap);

  private id = computed(() => Number(this.params()?.get('id')));
  isCreateMode = computed(() => !this.params()?.get('id'));

  model = signal({ name: '' });
  localityForm = form(this.model, (f) => {
    required(f.name);
  });

  localityQuery = injectQuery(() => ({
    queryKey: localityKeys.detail(this.id()),
    queryFn: () => this.api.fetchLocality(this.id()),
    enabled: !!this.id() && !this.isCreateMode(),
  }));

  permissionsQuery = injectQuery(() => ({
    queryKey: localityKeys.permissions(this.id()),
    queryFn: () => this.api.fetchLocalityPermissions(this.id()),
    enabled: !!this.id() && !this.isCreateMode(),
  }));

  mutation = injectMutation(() => ({
    mutationFn: (data: Parameters<typeof this.api.updateLocality>[1]) =>
      this.api.updateLocality(this.id(), data),
    onSuccess: () => {
      this.queryClient.invalidateQueries({ queryKey: localityKeys.all() });
      this.router.navigate(['..'], { relativeTo: this.route });
    },
  }));

  createMutation = injectMutation(() => ({
    mutationFn: (data: Parameters<typeof this.api.createLocality>[0]) =>
      this.api.createLocality(data),
    onSuccess: (result) => {
      this.queryClient.invalidateQueries({ queryKey: localityKeys.all() });
      this.router.navigate([
        '/compendium/localities',
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
      const data = this.localityQuery.data();
      if (data) {
        this.model.set({ name: data.name });
      }
    });
  }

  onSubmit() {
    if (this.localityForm().valid()) {
      const val = this.model();
      const data = this.localityQuery.data();
      const payload = {
        ...(this.isCreateMode() ? {} : { public: data!.public }),
        name: val.name,
      };
      if (this.isCreateMode()) {
        this.createMutation.mutate(payload);
      } else {
        this.mutation.mutate(payload);
      }
    }
  }
}
