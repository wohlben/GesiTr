import { Component, signal } from '@angular/core';
import { TestBed } from '@angular/core/testing';
import { form } from '@angular/forms/signals';
import { provideHttpClient } from '@angular/common/http';
import { provideTanStackQuery, QueryClient } from '@tanstack/angular-query-experimental';
import { provideTranslocoForTest } from '$core/testing/transloco-testing';
import { SectionItemEditor } from './section-item-editor';
import { EMPTY_GROUP_CONFIG } from '$ui/exercise-group-config/exercise-group-config';
import type { WorkoutItemModel } from './workout-item-model';
import type { Exercise } from '$generated/models';
import type { ExerciseGroup } from '$generated/user-models';

// brn-select / combobox use ResizeObserver
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

function makeItem(overrides: Partial<WorkoutItemModel> = {}): WorkoutItemModel {
  return {
    itemType: 'exercise',
    exerciseId: null,
    selectedSchemeId: null,
    sectionItemId: null,
    groupConfig: { ...EMPTY_GROUP_CONFIG },
    ...overrides,
  };
}

@Component({
  selector: 'app-test-host',
  template: `
    <app-section-item-editor
      [itemField]="itemForm"
      [exercises]="exercises"
      [exerciseGroups]="exerciseGroups"
      [itemLabel]="'Exercise 1'"
      (removed)="removedCalled = true"
    />
  `,
  imports: [SectionItemEditor],
})
class TestHost {
  model = signal(makeItem());
  itemForm = form(this.model);
  exercises = exercises;
  exerciseGroups: ExerciseGroup[] = [];
  removedCalled = false;

  resetModel(overrides: Partial<WorkoutItemModel> = {}) {
    this.model.set(makeItem(overrides));
  }
}

describe('SectionItemEditor', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [TestHost],
      providers: [
        provideTranslocoForTest(),
        provideHttpClient(),
        provideTanStackQuery(new QueryClient({ defaultOptions: { queries: { retry: false } } })),
      ],
    });
  });

  function setup(overrides: Partial<WorkoutItemModel> = {}) {
    const fixture = TestBed.createComponent(TestHost);
    if (Object.keys(overrides).length > 0) {
      fixture.componentInstance.resetModel(overrides);
    }
    fixture.detectChanges();
    return fixture;
  }

  it('shows item label', () => {
    const fixture = setup();
    expect(fixture.nativeElement.textContent).toContain('Exercise 1');
  });

  it('shows remove button', () => {
    const fixture = setup();
    expect(fixture.nativeElement.textContent).toContain('common.remove');
  });

  it('emits removed when remove clicked', () => {
    const fixture = setup();
    const removeBtn = Array.from(
      fixture.nativeElement.querySelectorAll('button') as NodeListOf<HTMLButtonElement>,
    ).find((b) => b.textContent?.includes('common.remove'));
    removeBtn?.click();
    expect(fixture.componentInstance.removedCalled).toBe(true);
  });

  it('shows type button group with two buttons', () => {
    const fixture = setup();
    const buttons = fixture.nativeElement.querySelectorAll('button');
    const typeButtons = Array.from(buttons as NodeListOf<HTMLButtonElement>).filter(
      (b) =>
        b.textContent?.includes('enums.workoutSectionItemType.exercise') ||
        b.textContent?.includes('enums.workoutSectionItemType.exercise_group'),
    );
    expect(typeButtons.length).toBe(2);
  });

  it('renders exercise-item-editor when itemType is exercise', () => {
    const fixture = setup({ itemType: 'exercise' });
    expect(fixture.nativeElement.querySelector('app-exercise-item-editor')).toBeTruthy();
    expect(fixture.nativeElement.querySelector('app-exercise-group-config')).toBeNull();
  });

  it('renders exercise-group-config when itemType is exercise_group', () => {
    const fixture = setup({ itemType: 'exercise_group' });
    expect(fixture.nativeElement.querySelector('app-exercise-group-config')).toBeTruthy();
    expect(fixture.nativeElement.querySelector('app-exercise-item-editor')).toBeNull();
  });

  it('switches itemType via form field when type button clicked', () => {
    const fixture = setup({ itemType: 'exercise' });
    const buttons = Array.from(
      fixture.nativeElement.querySelectorAll('button') as NodeListOf<HTMLButtonElement>,
    );
    const groupBtn = buttons.find((b) =>
      b.textContent?.includes('enums.workoutSectionItemType.exercise_group'),
    );
    groupBtn?.click();
    expect(fixture.componentInstance.model().itemType).toBe('exercise_group');
  });

  it('preserves other fields when switching type', () => {
    const fixture = setup({ itemType: 'exercise', exerciseId: 1, selectedSchemeId: 5 });
    const buttons = Array.from(
      fixture.nativeElement.querySelectorAll('button') as NodeListOf<HTMLButtonElement>,
    );
    const groupBtn = buttons.find((b) =>
      b.textContent?.includes('enums.workoutSectionItemType.exercise_group'),
    );
    groupBtn?.click();
    const item = fixture.componentInstance.model();
    expect(item.itemType).toBe('exercise_group');
    expect(item.exerciseId).toBe(1);
    expect(item.selectedSchemeId).toBe(5);
  });
});
