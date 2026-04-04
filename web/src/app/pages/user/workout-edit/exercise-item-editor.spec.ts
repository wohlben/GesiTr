import { Component, signal } from '@angular/core';
import { TestBed } from '@angular/core/testing';
import { form } from '@angular/forms/signals';
import { provideHttpClient } from '@angular/common/http';
import { provideTanStackQuery, QueryClient } from '@tanstack/angular-query-experimental';
import { provideTranslocoForTest } from '$core/testing/transloco-testing';
import { ExerciseItemEditor } from './exercise-item-editor';
import { EMPTY_GROUP_CONFIG } from '$ui/exercise-group-config/exercise-group-config';
import type { WorkoutItemModel } from './workout-item-model';
import type { Exercise } from '$generated/models';
import type { ExerciseScheme } from '$generated/user-exercisescheme';

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
    groupConfig: { ...EMPTY_GROUP_CONFIG },
    ...overrides,
  };
}

function makeScheme(overrides: Partial<ExerciseScheme> = {}): ExerciseScheme {
  return {
    id: 10,
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
    owner: 'u1',
    exerciseId: 1,
    measurementType: 'REP_BASED',
    ...overrides,
  };
}

describe('ExerciseItemEditor', () => {
  @Component({
    selector: 'app-test-host-editor',
    template: ` <app-exercise-item-editor [itemField]="itemForm" [exercises]="exercises" /> `,
    imports: [ExerciseItemEditor],
  })
  class Host {
    model = signal(makeItem());
    itemForm = form(this.model);
    exercises = exercises;

    resetModel(overrides: Partial<WorkoutItemModel> = {}) {
      this.model.set(makeItem(overrides));
    }
  }

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [Host],
      providers: [
        provideTranslocoForTest(),
        provideHttpClient(),
        provideTanStackQuery(new QueryClient({ defaultOptions: { queries: { retry: false } } })),
      ],
    });
  });

  function setup(overrides: Partial<WorkoutItemModel> = {}) {
    const fixture = TestBed.createComponent(Host);
    if (Object.keys(overrides).length > 0) {
      fixture.componentInstance.resetModel(overrides);
    }
    fixture.detectChanges();
    return fixture;
  }

  function getEditor(fixture: ReturnType<typeof setup>): ExerciseItemEditor {
    return fixture.debugElement.children[0].componentInstance as ExerciseItemEditor;
  }

  it('renders exercise combobox with formField binding', () => {
    const fixture = setup();
    expect(fixture.nativeElement.querySelector('hlm-combobox')).toBeTruthy();
  });

  it('renders scheme selector when exerciseId is set', () => {
    const fixture = setup({ exerciseId: 1 });
    expect(fixture.nativeElement.querySelector('app-scheme-selector')).toBeTruthy();
  });

  it('updates selectedSchemeId in form via onSchemeSelected', () => {
    const fixture = setup({ exerciseId: 1 });
    getEditor(fixture).onSchemeSelected(42);
    expect(fixture.componentInstance.model().selectedSchemeId).toBe(42);
  });

  it('opens dialog in create mode', () => {
    const fixture = setup({ exerciseId: 1 });
    const editor = getEditor(fixture);
    expect(editor.dialogOpen()).toBe(false);
    editor.openDialog(null);
    expect(editor.dialogOpen()).toBe(true);
    expect(editor.editingScheme()).toBeNull();
  });

  it('opens dialog in edit mode with scheme', () => {
    const fixture = setup({ exerciseId: 1 });
    const editor = getEditor(fixture);
    const scheme = makeScheme();
    editor.openDialog(scheme);
    expect(editor.dialogOpen()).toBe(true);
    expect(editor.editingScheme()).toBe(scheme);
  });

  it('closes dialog and clears editing scheme', () => {
    const fixture = setup({ exerciseId: 1 });
    const editor = getEditor(fixture);
    editor.openDialog(makeScheme());
    editor.closeDialog();
    expect(editor.dialogOpen()).toBe(false);
    expect(editor.editingScheme()).toBeNull();
  });

  it('sets selectedSchemeId in form and closes dialog on save', () => {
    const fixture = setup({ exerciseId: 1 });
    const editor = getEditor(fixture);
    editor.openDialog(null);
    editor.onSchemeSaved(makeScheme({ id: 99 }));
    expect(fixture.componentInstance.model().selectedSchemeId).toBe(99);
    expect(editor.dialogOpen()).toBe(false);
  });

  it('resolves exercise name from ID via exerciseIdToString', () => {
    const fixture = setup({ exerciseId: 1 });
    const editor = getEditor(fixture);
    expect(editor.exerciseIdToString(1)).toBe('Bench Press');
  });

  it('returns empty string for unknown exercise ID', () => {
    const fixture = setup();
    const editor = getEditor(fixture);
    expect(editor.exerciseIdToString(999)).toBe('');
  });

  it('filters exercise by name via exerciseIdFilter', () => {
    const fixture = setup();
    const editor = getEditor(fixture);
    expect(editor.exerciseIdFilter(1, 'bench')).toBe(true);
    expect(editor.exerciseIdFilter(1, 'squat')).toBe(false);
  });
});
