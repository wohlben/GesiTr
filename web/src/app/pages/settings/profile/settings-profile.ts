import { Component, inject, isDevMode, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { TranslocoDirective } from '@jsverse/transloco';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { profileKeys } from '$core/query-keys';
import { PageLayout } from '../../../layout/page-layout';
import { DevelopmentUserHeaderService } from '$core/dev/development-user-header.service';

@Component({
  selector: 'app-settings-profile',
  imports: [PageLayout, TranslocoDirective, FormsModule],
  template: `
    <ng-container *transloco="let t">
      <app-page-layout
        [header]="t('settings.profile.title')"
        [isPending]="profileQuery.isPending()"
        [errorMessage]="profileQuery.isError() ? profileQuery.error().message : undefined"
      >
        @if (profileQuery.data(); as profile) {
          <div class="space-y-6">
            <div>
              <span class="block text-sm font-medium text-gray-700 dark:text-gray-300">
                {{ t('settings.profile.userName') }}
              </span>
              <p class="mt-1 text-gray-900 dark:text-gray-100">{{ profile.name }}</p>
            </div>

            @if (isDevMode) {
              <div>
                <label
                  class="block text-sm font-medium text-gray-700 dark:text-gray-300"
                  for="dev-user-id"
                >
                  {{ t('settings.profile.devUserId') }}
                </label>
                <input
                  id="dev-user-id"
                  type="text"
                  [ngModel]="devUserId()"
                  (ngModelChange)="onDevUserIdChange($event)"
                  class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 text-sm shadow-sm focus:border-blue-500 focus:outline-none focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                />
                <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">
                  {{ t('settings.profile.devUserIdHint') }}
                </p>
              </div>
            }
          </div>
        }
      </app-page-layout>
    </ng-container>
  `,
})
export class SettingsProfile {
  private userApi = inject(UserApiClient);
  private devUserService = inject(DevelopmentUserHeaderService, { optional: true });

  isDevMode = isDevMode();
  devUserId = signal(this.devUserService?.userId$.getValue() ?? 'devuser');

  profileQuery = injectQuery(() => ({
    queryKey: profileKeys.mine(),
    queryFn: () => this.userApi.fetchProfile(),
  }));

  onDevUserIdChange(value: string) {
    this.devUserId.set(value);
    this.devUserService?.userId$.next(value);
  }
}
