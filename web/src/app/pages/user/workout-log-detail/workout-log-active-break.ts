import { Component, input, output } from '@angular/core';
import { TranslocoDirective } from '@jsverse/transloco';
import { formatBreak, formatCountdown } from '$core/format-utils';
import { ViewItemBreak } from './workout-log-view-items';

@Component({
  selector: 'app-workout-log-active-break',
  imports: [TranslocoDirective],
  template: `
    <ng-container *transloco="let t">
      @if (peeked()) {
        <div
          class="my-3 flex flex-col items-center justify-center rounded-xl border-2 border-amber-500 bg-amber-50/50 p-5 dark:border-amber-400 dark:bg-amber-950/20"
        >
          <div
            class="mb-2 text-xs font-semibold tracking-wider text-amber-600 uppercase dark:text-amber-400"
          >
            {{ t('user.workoutLog.rest') }}
          </div>
          <div class="mb-3 text-5xl font-bold tabular-nums text-amber-700 dark:text-amber-300">
            {{ formatBreak(data().seconds) }}
          </div>
          <div class="mb-3 text-sm text-gray-500 dark:text-gray-400">
            {{ t('user.workoutLog.next', { label: data().label }) }}
          </div>
          <button
            type="button"
            (click)="togglePeek.emit(); $event.stopPropagation()"
            class="flex w-full items-center justify-center py-1 text-gray-400 hover:text-gray-600 dark:text-gray-500 dark:hover:text-gray-300"
          >
            <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
              <path
                fill-rule="evenodd"
                d="M9.47 6.47a.75.75 0 011.06 0l4.25 4.25a.75.75 0 11-1.06 1.06L10 8.06l-3.72 3.72a.75.75 0 01-1.06-1.06l4.25-4.25z"
                clip-rule="evenodd"
              />
            </svg>
          </button>
        </div>
      } @else if (data().role === 'active-timer') {
        <div
          class="my-3 flex flex-col items-center justify-center rounded-xl border-2 border-amber-500 bg-amber-50/50 p-5 dark:border-amber-400 dark:bg-amber-950/20"
        >
          <div
            class="mb-2 text-xs font-semibold tracking-wider text-amber-600 uppercase dark:text-amber-400"
          >
            {{ t('user.workoutLog.rest') }}
          </div>
          <div class="mb-3 text-5xl font-bold tabular-nums text-amber-700 dark:text-amber-300">
            {{ formatCountdown(remainingSeconds()) }}
          </div>
          <div class="mb-5 text-sm text-gray-500 dark:text-gray-400">
            {{ t('user.workoutLog.next', { label: data().label }) }}
          </div>
          <button
            type="button"
            (click)="skip.emit()"
            class="w-full rounded-lg border-2 border-amber-500 bg-transparent px-4 py-3 text-base font-semibold text-amber-700 hover:bg-amber-100 active:bg-amber-200 dark:border-amber-400 dark:text-amber-300 dark:hover:bg-amber-900/30 dark:active:bg-amber-900/50"
          >
            {{ t('common.skip') }}
          </button>
        </div>
      } @else {
        <div
          class="relative flex h-0 cursor-pointer items-center justify-center overflow-visible z-10"
          role="button"
          tabindex="0"
          (click)="togglePeek.emit()"
          (keydown.enter)="togglePeek.emit()"
        >
          <div
            class="absolute inset-x-0 top-1/2 border-t border-dashed border-gray-200 dark:border-gray-700"
          ></div>
          <div
            class="relative z-10 flex items-center gap-1 rounded-full bg-gray-50 px-2 py-0.5 text-xs text-gray-400 dark:bg-gray-900 dark:text-gray-500"
          >
            <svg class="h-3 w-3" viewBox="0 0 20 20" fill="currentColor">
              <path
                fill-rule="evenodd"
                d="M10 18a8 8 0 100-16 8 8 0 000 16zm.75-13a.75.75 0 00-1.5 0v5c0 .414.336.75.75.75h4a.75.75 0 000-1.5h-3.25V5z"
                clip-rule="evenodd"
              />
            </svg>
            {{ formatBreak(data().seconds) }}
          </div>
        </div>
      }
    </ng-container>
  `,
})
export class WorkoutLogActiveBreak {
  data = input.required<ViewItemBreak>();
  peeked = input(false);
  remainingSeconds = input<number>(0);
  skip = output<void>();
  togglePeek = output<void>();

  formatBreak = formatBreak;
  formatCountdown = formatCountdown;
}
