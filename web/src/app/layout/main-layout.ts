import { Component, inject, isDevMode, signal } from '@angular/core';
import { Router, RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';
import { injectQueryClient } from '@tanstack/angular-query-experimental';
import { TranslocoDirective } from '@jsverse/transloco';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideSettings, lucideUser } from '@ng-icons/lucide';
import { HlmIconImports } from '@spartan-ng/helm/icon';
import { HlmPopoverImports } from '@spartan-ng/helm/popover';
import { DevelopmentUserHeaderService } from '$core/dev/development-user-header.service';

@Component({
  selector: 'app-main-layout',
  imports: [
    RouterOutlet,
    RouterLink,
    RouterLinkActive,
    TranslocoDirective,
    NgIcon,
    HlmIconImports,
    HlmPopoverImports,
  ],
  providers: [provideIcons({ lucideSettings, lucideUser })],
  host: { class: 'block' },
  template: `
    <div *transloco="let t" class="min-h-screen bg-gray-50 dark:bg-gray-950">
      <nav
        class="sticky top-0 z-20 border-b border-gray-200 bg-white dark:border-gray-800 dark:bg-gray-900"
      >
        <div class="flex items-center justify-between px-6 py-3">
          <div class="flex items-center gap-4">
            <span class="text-lg font-semibold text-gray-900 dark:text-gray-100">GesiTr</span>
            <a
              routerLink="/user/workouts"
              routerLinkActive="bg-gray-100 text-gray-900 dark:bg-gray-800 dark:text-gray-100"
              class="rounded-md px-3 py-1.5 text-sm font-medium text-gray-600 transition-colors hover:bg-gray-50 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-gray-800/50 dark:hover:text-gray-200 md:hidden"
            >
              {{ t('nav.workouts') }}
            </a>
          </div>

          <!-- Desktop nav -->
          <div class="hidden items-center gap-1 md:flex">
            @for (link of compendiumLinks; track link.path) {
              <a
                [routerLink]="link.path"
                routerLinkActive="bg-gray-100 text-gray-900 dark:bg-gray-800 dark:text-gray-100"
                class="rounded-md px-3 py-1.5 text-sm font-medium text-gray-600 transition-colors hover:bg-gray-50 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-gray-800/50 dark:hover:text-gray-200"
              >
                {{ t(link.labelKey) }}
              </a>
            }
            <div class="mx-2 h-5 w-px bg-gray-200 dark:bg-gray-700"></div>
            @for (link of userLinks; track link.path) {
              <a
                [routerLink]="link.path"
                routerLinkActive="bg-gray-100 text-gray-900 dark:bg-gray-800 dark:text-gray-100"
                class="rounded-md px-3 py-1.5 text-sm font-medium text-gray-600 transition-colors hover:bg-gray-50 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-gray-800/50 dark:hover:text-gray-200"
              >
                {{ t(link.labelKey) }}
              </a>
            }
          </div>

          <div class="flex items-center gap-1">
            @if (isDevMode) {
              <div
                hlmPopover
                [state]="userPickerOpen() ? 'open' : 'closed'"
                (closed)="userPickerOpen.set(false)"
                align="end"
              >
                <button
                  hlmPopoverTrigger
                  (click)="userPickerOpen.set(!userPickerOpen())"
                  class="rounded-md p-2 text-gray-600 transition-colors hover:bg-gray-50 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-gray-800/50 dark:hover:text-gray-200"
                  aria-label="Switch user"
                >
                  <ng-icon hlm name="lucideUser" size="sm" />
                </button>
                <ng-template hlmPopoverPortal>
                  <div hlmPopoverContent class="w-48 p-2">
                    <div class="mb-2 text-xs font-medium text-gray-500 dark:text-gray-400">
                      Switch user ({{ devUserService.userId$.getValue() }})
                    </div>
                    @for (user of devUsers; track user) {
                      <button
                        (click)="switchUser(user); userPickerOpen.set(false)"
                        class="w-full rounded-md px-3 py-1.5 text-left text-sm transition-colors hover:bg-gray-100 dark:hover:bg-gray-800"
                        [class.font-semibold]="user === devUserService.userId$.getValue()"
                      >
                        {{ user }}
                      </button>
                    }
                  </div>
                </ng-template>
              </div>
            }
            <a
              routerLink="/settings/profile"
              routerLinkActive="bg-gray-100 text-gray-900 dark:bg-gray-800 dark:text-gray-100"
              class="rounded-md p-2 text-gray-600 transition-colors hover:bg-gray-50 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-gray-800/50 dark:hover:text-gray-200"
              [attr.aria-label]="t('nav.settings')"
            >
              <ng-icon hlm name="lucideSettings" size="sm" />
            </a>

            <!-- Mobile burger button -->
            <button
              (click)="menuOpen.set(!menuOpen())"
              class="rounded-md p-2 text-gray-600 hover:bg-gray-100 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-gray-800 dark:hover:text-gray-200 md:hidden"
              [attr.aria-label]="t('nav.toggleMenu')"
            >
              <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                @if (menuOpen()) {
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M6 18L18 6M6 6l12 12"
                  />
                } @else {
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M4 6h16M4 12h16M4 18h16"
                  />
                }
              </svg>
            </button>
          </div>
        </div>

        <!-- Mobile dropdown menu -->
        @if (menuOpen()) {
          <div class="border-t border-gray-200 px-6 py-4 dark:border-gray-800 md:hidden">
            <div class="mb-3">
              <span
                class="text-xs font-medium uppercase tracking-wider text-gray-400 dark:text-gray-500"
                >{{ t('nav.compendium') }}</span
              >
              <div class="mt-1 flex flex-col gap-0.5">
                @for (link of compendiumLinks; track link.path) {
                  <a
                    [routerLink]="link.path"
                    (click)="menuOpen.set(false)"
                    routerLinkActive="bg-gray-100 text-gray-900 dark:bg-gray-800 dark:text-gray-100"
                    class="rounded-md px-3 py-2 text-sm font-medium text-gray-600 transition-colors hover:bg-gray-50 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-gray-800/50 dark:hover:text-gray-200"
                  >
                    {{ t(link.labelKey) }}
                  </a>
                }
              </div>
            </div>
            <div>
              <span
                class="text-xs font-medium uppercase tracking-wider text-gray-400 dark:text-gray-500"
                >{{ t('nav.personal') }}</span
              >
              <div class="mt-1 flex flex-col gap-0.5">
                @for (link of userLinks; track link.path) {
                  <a
                    [routerLink]="link.path"
                    (click)="menuOpen.set(false)"
                    routerLinkActive="bg-gray-100 text-gray-900 dark:bg-gray-800 dark:text-gray-100"
                    class="rounded-md px-3 py-2 text-sm font-medium text-gray-600 transition-colors hover:bg-gray-50 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-gray-800/50 dark:hover:text-gray-200"
                  >
                    {{ t(link.labelKey) }}
                  </a>
                }
              </div>
            </div>
          </div>
        }
      </nav>
      <main class="p-6">
        <router-outlet />
      </main>
    </div>
  `,
})
export class MainLayout {
  menuOpen = signal(false);
  userPickerOpen = signal(false);
  isDevMode = isDevMode();
  devUserService = isDevMode() ? inject(DevelopmentUserHeaderService) : null!;
  private router = inject(Router);
  private queryClient = injectQueryClient();

  devUsers = ['devuser', 'anon', 'sinon', 'alice', 'bob'];

  compendiumLinks = [
    { path: '/compendium/exercises', labelKey: 'nav.exercises' },
    { path: '/compendium/equipment', labelKey: 'nav.equipment' },
  ];

  userLinks = [
    { path: '/user/equipment', labelKey: 'nav.myEquipment' },
    { path: '/user/workouts', labelKey: 'nav.myWorkouts' },
    { path: '/user/calendar', labelKey: 'nav.calendar' },
  ];

  switchUser(userId: string) {
    this.devUserService.userId$.next(userId);
    this.queryClient.invalidateQueries();
    this.router.navigate([], {
      queryParams: { onBehalfOf: userId },
      queryParamsHandling: 'merge',
    });
  }
}
