import { render } from '@testing-library/angular';
import { page } from 'vitest/browser';
import { provideTranslocoForTest } from '$core/testing/transloco-testing';
import { provideHttpClient } from '@angular/common/http';
import { provideTanStackQuery, QueryClient } from '@tanstack/angular-query-experimental';
import { SectionItemEditor } from './section-item-editor';
import { EMPTY_GROUP_CONFIG } from '$ui/exercise-group-config/exercise-group-config';
import type { WorkoutItemModel } from './workout-item-model';
import type { Exercise } from '$generated/models';

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

const exerciseItem: WorkoutItemModel = {
  itemType: 'exercise',
  exerciseId: 1,
  selectedSchemeId: null,
  groupConfig: { ...EMPTY_GROUP_CONFIG },
};

const groupItem: WorkoutItemModel = {
  itemType: 'exercise_group',
  exerciseId: null,
  selectedSchemeId: null,
  groupConfig: {
    exerciseGroupId: 5,
    name: 'Push Variants',
    members: [1],
  },
};

const providers = [
  provideTranslocoForTest(),
  provideHttpClient(),
  provideTanStackQuery(new QueryClient({ defaultOptions: { queries: { retry: false } } })),
];

describe('SectionItemEditor screenshots', () => {
  afterEach(() => {
    document.documentElement.classList.remove('dark');
  });

  describe('exercise type', () => {
    const template = `
      <app-section-item-editor
        [(value)]="value"
        [exercises]="exercises"
        [exerciseGroups]="[]"
        [itemLabel]="'Exercise 1'"
      />
    `;
    const opts = {
      imports: [SectionItemEditor],
      providers,
      componentProperties: { value: exerciseItem, exercises },
    };

    it('light', async () => {
      const { fixture } = await render(template, opts);
      const locator = page.elementLocator(fixture.nativeElement);
      await expect(locator).toMatchScreenshot('exercise-light');
    });

    it('dark', async () => {
      document.documentElement.classList.add('dark');
      const { fixture } = await render(template, opts);
      const locator = page.elementLocator(fixture.nativeElement);
      await expect(locator).toMatchScreenshot('exercise-dark');
    });
  });

  describe('exercise group type', () => {
    const template = `
      <app-section-item-editor
        [(value)]="value"
        [exercises]="exercises"
        [exerciseGroups]="exerciseGroups"
        [itemLabel]="'Exercise 2'"
      />
    `;
    const opts = {
      imports: [SectionItemEditor],
      providers,
      componentProperties: {
        value: groupItem,
        exercises,
        exerciseGroups: [
          {
            id: 5,
            name: 'Push Variants',
            ownershipGroupId: 0,
            createdAt: '2024-01-01T00:00:00Z',
            updatedAt: '2024-01-01T00:00:00Z',
          },
        ],
      },
    };

    it('light', async () => {
      const { fixture } = await render(template, opts);
      const locator = page.elementLocator(fixture.nativeElement);
      await expect(locator).toMatchScreenshot('group-light');
    });

    it('dark', async () => {
      document.documentElement.classList.add('dark');
      const { fixture } = await render(template, opts);
      const locator = page.elementLocator(fixture.nativeElement);
      await expect(locator).toMatchScreenshot('group-dark');
    });
  });
});
