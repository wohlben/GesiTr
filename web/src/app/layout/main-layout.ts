import { Component } from '@angular/core';
import { RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';

@Component({
  selector: 'app-main-layout',
  imports: [RouterOutlet, RouterLink, RouterLinkActive],
  template: `
    <div class="flex h-screen flex-col bg-gray-50 dark:bg-gray-950">
      <nav
        class="border-b border-gray-200 bg-white px-6 py-3 dark:border-gray-800 dark:bg-gray-900"
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
      <main class="flex-1 overflow-auto p-6">
        <router-outlet />
      </main>
    </div>
  `,
})
export class MainLayout {
  navLinks = [
    { path: 'exercises', label: 'Exercises' },
    { path: 'equipment', label: 'Equipment' },
    { path: 'exercise-groups', label: 'Exercise Groups' },
  ];
}
