import { Component, input, output } from '@angular/core';
import { RouterLink } from '@angular/router';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideExternalLink } from '@ng-icons/lucide';
import { Exercise } from '$generated/models';
import { ExerciseMastery } from '$generated/user-mastery';
import { HlmIcon } from '$ui/spartan/icon/src/lib/hlm-icon';
import { SlugifyPipe } from '$ui/pipes/slugify';

@Component({
  selector: 'tr[app-exercise-list-item]',
  imports: [RouterLink, SlugifyPipe, NgIcon, HlmIcon],
  providers: [provideIcons({ lucideExternalLink })],
  template: `
    <td class="whitespace-nowrap px-4 py-2 text-sm font-medium text-gray-900 dark:text-gray-100">
      <span class="inline-flex items-center gap-1">
        {{ displayName() }}
        <a
          [routerLink]="['/compendium/exercises', exercise().id, displayName() | slugify]"
          class="inline-flex items-center justify-center rounded-full p-1 text-gray-400 hover:text-gray-600 dark:hover:text-gray-200"
          (click)="nameClicked.emit(displayName())"
        >
          <ng-icon hlm name="lucideExternalLink" size="sm" />
        </a>
      </span>
      @for (name of matchingNames(); track name) {
        <div class="text-xs text-gray-400 dark:text-gray-500 pl-1">
          <a
            [routerLink]="['/compendium/exercises', exercise().id, name | slugify]"
            class="hover:text-gray-600 dark:hover:text-gray-300 cursor-pointer"
            (click)="nameClicked.emit(name)"
            >{{ name }}</a
          >
        </div>
      }
    </td>
    <td class="whitespace-nowrap px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      @if (mastery(); as m) {
        <span
          class="inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium"
          [class]="tierBadgeClass(m.tier)"
        >
          Lv.{{ m.level }}
        </span>
      }
    </td>
    <td class="whitespace-nowrap px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ exercise().type }}
    </td>
    <td class="whitespace-nowrap px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ exercise().technicalDifficulty }}
    </td>
    <td class="whitespace-nowrap px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ exercise().force?.join(', ') }}
    </td>
    <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ exercise().primaryMuscles?.join(', ') }}
    </td>
    <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ exercise().secondaryMuscles?.join(', ') }}
    </td>
    <td class="whitespace-nowrap px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ exercise().bodyWeightScaling }}
    </td>
    <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ exercise().suggestedMeasurementParadigms?.join(', ') }}
    </td>
    <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ exercise().description }}
    </td>
    <td class="px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ exercise().names?.map(n => n.name)?.join(', ') }}
    </td>
    <td class="whitespace-nowrap px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ exercise().authorName }}
    </td>
    <td class="whitespace-nowrap px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ exercise().version }}
    </td>
    <td class="whitespace-nowrap px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ exercise().owner }}
    </td>
    <td class="whitespace-nowrap px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ exercise().createdAt }}
    </td>
    <td class="whitespace-nowrap px-4 py-2 text-sm text-gray-600 dark:text-gray-400">
      {{ exercise().updatedAt }}
    </td>
  `,
})
export class ExerciseListItem {
  exercise = input.required<Exercise>();
  displayName = input.required<string>();
  matchingNames = input<string[]>([]);
  mastery = input<ExerciseMastery | undefined>(undefined);
  nameClicked = output<string>();

  tierBadgeClass(tier: string): string {
    switch (tier) {
      case 'mastered':
        return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-300';
      case 'master':
        return 'bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-300';
      case 'adept':
        return 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300';
      case 'journeyman':
        return 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300';
      default:
        return 'bg-gray-100 text-gray-600 dark:bg-gray-800 dark:text-gray-400';
    }
  }
}
