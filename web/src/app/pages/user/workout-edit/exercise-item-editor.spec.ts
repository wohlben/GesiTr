import { Component, signal } from '@angular/core';
import { TestBed } from '@angular/core/testing';
import { provideHttpClient } from '@angular/common/http';
import { provideTanStackQuery, QueryClient } from '@tanstack/angular-query-experimental';
import { provideTranslocoForTest } from '$core/testing/transloco-testing';
import { ExerciseItemEditor } from './exercise-item-editor';
import { EMPTY_GROUP_CONFIG } from '$ui/exercise-group-config/exercise-group-config';
import type { WorkoutItemModel } from './workout-item-model';
import type { Exercise } from '$generated/models';
import type { ExerciseScheme } from '$generated/user-exercisescheme';

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
  {
    id: 2,
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
    names: [{ id: 2, name: 'Squat' }],
    type: 'STRENGTH',
    force: ['PUSH'],
    primaryMuscles: ['QUADS'],
    secondaryMuscles: ['GLUTES'],
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

describe('ExerciseItemEditor', () => {
  @Component({
    selector: 'app-test-host-editor',
    template: ` <app-exercise-item-editor [(value)]="item" [exercises]="exercises" /> `,
    imports: [ExerciseItemEditor],
  })
  class Host {
    item = signal(makeItem());
    exercises = exercises;
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
      fixture.componentInstance.item.set(makeItem(overrides));
    }
    fixture.detectChanges();
    return fixture;
  }

  it('renders exercise combobox', () => {
    const fixture = setup();
    expect(fixture.nativeElement.querySelector('hlm-combobox')).toBeTruthy();
  });

  it('renders scheme selector', () => {
    const fixture = setup({ exerciseId: 1 });
    expect(fixture.nativeElement.querySelector('app-scheme-selector')).toBeTruthy();
  });

  it('updates exerciseId via onExerciseSelected', () => {
    const fixture = setup();
    const editor = fixture.debugElement.children[0].componentInstance as ExerciseItemEditor;
    editor.onExerciseSelected(exercises[0]);
    expect(fixture.componentInstance.item().exerciseId).toBe(1);
  });

  it('updates selectedSchemeId via onSchemeSelected', () => {
    const fixture = setup({ exerciseId: 1 });
    const editor = fixture.debugElement.children[0].componentInstance as ExerciseItemEditor;
    editor.onSchemeSelected(42);
    expect(fixture.componentInstance.item().selectedSchemeId).toBe(42);
  });

  it('clears exerciseId when exercise deselected', () => {
    const fixture = setup({ exerciseId: 1 });
    const editor = fixture.debugElement.children[0].componentInstance as ExerciseItemEditor;
    editor.onExerciseSelected(null);
    expect(fixture.componentInstance.item().exerciseId).toBeNull();
  });

  it('opens dialog in create mode', () => {
    const fixture = setup({ exerciseId: 1 });
    const editor = fixture.debugElement.children[0].componentInstance as ExerciseItemEditor;
    expect(editor.dialogOpen()).toBe(false);
    editor.openDialog(null);
    expect(editor.dialogOpen()).toBe(true);
    expect(editor.editingScheme()).toBeNull();
  });

  it('opens dialog in edit mode with scheme', () => {
    const fixture = setup({ exerciseId: 1 });
    const editor = fixture.debugElement.children[0].componentInstance as ExerciseItemEditor;
    const scheme = makeScheme();
    editor.openDialog(scheme);
    expect(editor.dialogOpen()).toBe(true);
    expect(editor.editingScheme()).toBe(scheme);
  });

  it('closes dialog and clears editing scheme', () => {
    const fixture = setup({ exerciseId: 1 });
    const editor = fixture.debugElement.children[0].componentInstance as ExerciseItemEditor;
    editor.openDialog(makeScheme());
    editor.closeDialog();
    expect(editor.dialogOpen()).toBe(false);
    expect(editor.editingScheme()).toBeNull();
  });

  it('selects scheme and closes dialog on save', () => {
    const fixture = setup({ exerciseId: 1 });
    const editor = fixture.debugElement.children[0].componentInstance as ExerciseItemEditor;
    editor.openDialog(null);
    editor.onSchemeSaved(makeScheme({ id: 99 }));
    expect(fixture.componentInstance.item().selectedSchemeId).toBe(99);
    expect(editor.dialogOpen()).toBe(false);
  });

  it('resolves selectedExercise from exercises list', () => {
    const fixture = setup({ exerciseId: 1 });
    const editor = fixture.debugElement.children[0].componentInstance as ExerciseItemEditor;
    expect(editor.selectedExercise()?.id).toBe(1);
    expect(editor.selectedExercise()?.names?.[0]?.name).toBe('Bench Press');
  });

  it('returns null selectedExercise when no exerciseId', () => {
    const fixture = setup({ exerciseId: null });
    const editor = fixture.debugElement.children[0].componentInstance as ExerciseItemEditor;
    expect(editor.selectedExercise()).toBeNull();
  });
});
