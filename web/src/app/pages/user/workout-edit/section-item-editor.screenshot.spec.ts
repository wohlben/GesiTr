import { Component, signal } from '@angular/core';
import { render } from '@testing-library/angular';
import { page } from 'vitest/browser';
import { form } from '@angular/forms/signals';
import { provideTranslocoForTest } from '$core/testing/transloco-testing';
import { provideHttpClient } from '@angular/common/http';
import { provideTanStackQuery, QueryClient } from '@tanstack/angular-query-experimental';
import { SectionItemEditor } from './section-item-editor';
import { EMPTY_GROUP_CONFIG } from '$ui/exercise-group-config/exercise-group-config';
import type { WorkoutItemModel } from './workout-item-model';
import type { Exercise } from '$generated/models';
import type { ExerciseGroup } from '$generated/user-models';

// brn-select uses ResizeObserver
beforeAll(() => {
  globalThis.ResizeObserver ??= class {
    // eslint-disable-next-line @typescript-eslint/no-empty-function
    observe() {}
    // eslint-disable-next-line @typescript-eslint/no-empty-function
    unobserve() {}
    // eslint-disable-next-line @typescript-eslint/no-empty-function
    disconnect() {}
  } as unknown as typeof ResizeObserver;
});

const exercises: Exercise[] = [
  {
    id: 1,
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
    names: [{ id: 1, name: 'Bench Press' }],
    type: 'STRENGTH',
    force: ['PUSH'],
    primaryMuscles: ['CHEST'],
    secondaryMuscles: ['TRICEPS'],
    technicalDifficulty: 'intermediate',
    bodyWeightScaling: 0,
    suggestedMeasurementParadigms: ['REP_BASED'],
    description: '',
    instructions: [],
    images: [],
    ownershipGroupId: 0,
    public: true,
    version: 1,
    equipmentIds: [],
  },
];

const exerciseGroups: ExerciseGroup[] = [
  {
    id: 5,
    name: 'Push Variants',
    ownershipGroupId: 0,
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
  },
];

const providers = [
  provideTranslocoForTest(),
  provideHttpClient(),
  provideTanStackQuery(new QueryClient({ defaultOptions: { queries: { retry: false } } })),
];

@Component({
  selector: 'app-exercise-host',
  template: `
    <app-section-item-editor
      [itemField]="itemForm"
      [exercises]="exercises"
      [exerciseGroups]="[]"
      [itemLabel]="'Exercise 1'"
    />
  `,
  imports: [SectionItemEditor],
})
class ExerciseHost {
  model = signal<WorkoutItemModel>({
    itemType: 'exercise',
    exerciseId: 1,
    selectedSchemeId: null,
    sectionItemId: null,
    groupConfig: { ...EMPTY_GROUP_CONFIG },
  });
  itemForm = form(this.model);
  exercises = exercises;
}

@Component({
  selector: 'app-group-host',
  template: `
    <app-section-item-editor
      [itemField]="itemForm"
      [exercises]="exercises"
      [exerciseGroups]="exerciseGroups"
      [itemLabel]="'Exercise 2'"
    />
  `,
  imports: [SectionItemEditor],
})
class GroupHost {
  model = signal<WorkoutItemModel>({
    itemType: 'exercise_group',
    exerciseId: null,
    selectedSchemeId: null,
    sectionItemId: null,
    groupConfig: { exerciseGroupId: 5, name: 'Push Variants', members: [1] },
  });
  itemForm = form(this.model);
  exercises = exercises;
  exerciseGroups = exerciseGroups;
}

describe('SectionItemEditor screenshots', () => {
  afterEach(() => {
    document.documentElement.classList.remove('dark');
  });

  describe('exercise type', () => {
    it('light', async () => {
      const { fixture } = await render(ExerciseHost, { providers });
      const locator = page.elementLocator(fixture.nativeElement);
      await expect(locator).toMatchScreenshot('exercise-light');
    });

    it('dark', async () => {
      document.documentElement.classList.add('dark');
      const { fixture } = await render(ExerciseHost, { providers });
      const locator = page.elementLocator(fixture.nativeElement);
      await expect(locator).toMatchScreenshot('exercise-dark');
    });
  });

  describe('exercise group type', () => {
    it('light', async () => {
      const { fixture } = await render(GroupHost, { providers });
      const locator = page.elementLocator(fixture.nativeElement);
      await expect(locator).toMatchScreenshot('group-light');
    });

    it('dark', async () => {
      document.documentElement.classList.add('dark');
      const { fixture } = await render(GroupHost, { providers });
      const locator = page.elementLocator(fixture.nativeElement);
      await expect(locator).toMatchScreenshot('group-dark');
    });
  });
});
