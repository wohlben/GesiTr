import { Component, signal } from '@angular/core';
import { render } from '@testing-library/angular';
import { page } from 'vitest/browser';
import { form } from '@angular/forms/signals';
import { provideTranslocoForTest } from '$core/testing/transloco-testing';
import { WorkoutStartSection } from './workout-start-section';
import type { ExerciseDisplayInfo } from './workout-start.store';
import type { StartModel } from './workout-start.models';

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

const DISPLAY_MAP: Record<number, ExerciseDisplayInfo> = {
  1: {
    name: 'Bench Press',
    summary: '3x10 @ 60kg',
    measurementType: 'REP_BASED',
    sets: [],
  },
  2: {
    name: 'Plank',
    summary: '60s',
    measurementType: 'TIME_BASED',
    sets: [],
  },
  3: {
    name: 'Running',
    summary: '5000m',
    measurementType: 'DISTANCE_BASED',
    sets: [],
  },
};

function makeModel(): StartModel {
  return {
    name: 'Full Workout',
    notes: '',
    sections: [
      {
        id: 1,
        type: 'main',
        label: 'Main Section',
        exercises: [
          {
            id: 1,
            sourceExerciseSchemeId: 10,
            breakAfterSeconds: 90,
            sets: [
              {
                id: 1,
                targetReps: 10,
                targetWeight: 60,
                targetDuration: null,
                targetDistance: null,
                targetTime: null,
                restAfterSeconds: 45,
              },
              {
                id: 2,
                targetReps: 10,
                targetWeight: 60,
                targetDuration: null,
                targetDistance: null,
                targetTime: null,
                restAfterSeconds: 45,
              },
              {
                id: 3,
                targetReps: 10,
                targetWeight: 60,
                targetDuration: null,
                targetDistance: null,
                targetTime: null,
                restAfterSeconds: null,
              },
            ],
          },
          {
            id: 2,
            sourceExerciseSchemeId: 20,
            breakAfterSeconds: 60,
            sets: [
              {
                id: 4,
                targetReps: null,
                targetWeight: null,
                targetDuration: 60,
                targetDistance: null,
                targetTime: null,
                restAfterSeconds: 30,
              },
              {
                id: 5,
                targetReps: null,
                targetWeight: null,
                targetDuration: 60,
                targetDistance: null,
                targetTime: null,
                restAfterSeconds: null,
              },
            ],
          },
          {
            id: 3,
            sourceExerciseSchemeId: 30,
            breakAfterSeconds: null,
            sets: [
              {
                id: 6,
                targetReps: null,
                targetWeight: null,
                targetDuration: null,
                targetDistance: 5000,
                targetTime: 1800,
                restAfterSeconds: null,
              },
            ],
          },
        ],
        pendingGroups: [],
      },
    ],
  };
}

const providers = [provideTranslocoForTest()];

@Component({
  selector: 'app-section-host',
  template: `
    <app-workout-start-section
      [section]="startForm.sections[0]"
      [sectionIndex]="0"
      [exerciseDisplayMap]="displayMap"
      [readonly]="false"
    />
  `,
  imports: [WorkoutStartSection],
})
class SectionHost {
  model = signal(makeModel());
  startForm = form(this.model);
  displayMap = DISPLAY_MAP;
}

@Component({
  selector: 'app-section-host-defaults',
  template: `
    <app-workout-start-section
      [section]="startForm.sections[0]"
      [sectionIndex]="0"
      [exerciseDisplayMap]="displayMap"
      [readonly]="false"
    />
  `,
  imports: [WorkoutStartSection],
})
class SectionHostDefaults {
  model = signal<StartModel>({
    name: 'Workout',
    notes: '',
    sections: [
      {
        id: 1,
        type: 'main',
        label: 'Section',
        exercises: [
          {
            id: 1,
            sourceExerciseSchemeId: 10,
            breakAfterSeconds: null,
            sets: [
              {
                id: 1,
                targetReps: 0,
                targetWeight: 0,
                targetDuration: null,
                targetDistance: null,
                targetTime: null,
                restAfterSeconds: 0,
              },
              {
                id: 2,
                targetReps: 0,
                targetWeight: 0,
                targetDuration: null,
                targetDistance: null,
                targetTime: null,
                restAfterSeconds: null,
              },
            ],
          },
        ],
        pendingGroups: [],
      },
    ],
  });
  startForm = form(this.model);
  displayMap: Record<number, ExerciseDisplayInfo> = {
    1: {
      name: 'Bench Press',
      summary: 'Rep based',
      measurementType: 'REP_BASED',
      sets: [],
    },
  };
}

describe('WorkoutStartSection screenshots', () => {
  afterEach(() => {
    document.documentElement.classList.remove('dark');
  });

  describe('with filled targets across measurement types', () => {
    it('light', async () => {
      const { fixture } = await render(SectionHost, { providers });
      const locator = page.elementLocator(fixture.nativeElement);
      await expect(locator).toMatchScreenshot('filled-light');
    });

    it('dark', async () => {
      document.documentElement.classList.add('dark');
      const { fixture } = await render(SectionHost, { providers });
      const locator = page.elementLocator(fixture.nativeElement);
      await expect(locator).toMatchScreenshot('filled-dark');
    });
  });

  describe('with zero-defaulted targets (scheme had no values)', () => {
    it('light', async () => {
      const { fixture } = await render(SectionHostDefaults, { providers });
      const locator = page.elementLocator(fixture.nativeElement);
      await expect(locator).toMatchScreenshot('defaults-light');
    });

    it('dark', async () => {
      document.documentElement.classList.add('dark');
      const { fixture } = await render(SectionHostDefaults, { providers });
      const locator = page.elementLocator(fixture.nativeElement);
      await expect(locator).toMatchScreenshot('defaults-dark');
    });
  });
});
