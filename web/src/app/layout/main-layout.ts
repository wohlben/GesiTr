import { Component } from '@angular/core';
import { RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';

@Component({
  selector: 'app-main-layout',
  imports: [RouterOutlet, RouterLink, RouterLinkActive],
  host: { class: 'block' },
  template: `
    <div class="min-h-screen bg-gray-50 dark:bg-gray-950">
      <nav
        class="sticky top-0 z-20 border-b border-gray-200 bg-white px-6 py-3 dark:border-gray-800 dark:bg-gray-900"
      >
        <div class="flex items-center gap-6">
          <span class="text-lg font-semibold text-gray-900 dark:text-gray-100">GesiTr</span>
          <div class="flex gap-1">
            @for (link of navLinks; track link.path) {
              <a
                [routerLink]="link.path"
                routerLinkActive="bg-gray-100 text-gray-900 dark:bg-gray-800 dark:text-gray-100"
                class="rounded-md px-3 py-1.5 text-sm font-medium text-gray-600 transition-colors hover:bg-gray-50 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-gray-800/50 dark:hover:text-gray-200"
              >
                {{ link.label }}
              </a>
            }
          </div>
        </div>
      </nav>
      <main class="p-6">
        <router-outlet />
      </main>
    </div>
  `,
})
export class MainLayout {
  navLinks = [
    { path: '/compendium/exercises', label: 'Exercises' },
    { path: '/compendium/equipment', label: 'Equipment' },
    { path: '/compendium/exercise-groups', label: 'Exercise Groups' },
    { path: '/user/exercises', label: 'My Exercises' },
    { path: '/user/equipment', label: 'My Equipment' },
    { path: '/user/workouts', label: 'My Workouts' },
  ];
}
