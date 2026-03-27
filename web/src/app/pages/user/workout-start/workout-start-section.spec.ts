import { Component, signal } from '@angular/core';
import { TestBed } from '@angular/core/testing';
import { form } from '@angular/forms/signals';
import { provideTranslocoForTest } from '$core/testing/transloco-testing';
import { WorkoutStartSection } from './workout-start-section';
import { ExerciseDisplayInfo } from './workout-start.store';
import { StartModel } from './workout-start.models';

// brn-select uses ResizeObserver which isn't available in the unit test environment
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
};

function makeModel(overrides: { withPendingGroup?: boolean } = {}): StartModel {
  return {
    name: 'Test',
    notes: '',
    sections: [
      {
        id: 1,
        type: 'main',
        label: 'Chest',
        exercises: [
          {
            id: 1,
            sourceExerciseSchemeId: 10,
            breakAfterSeconds: 60,
            sets: [
              {
                id: 1,
                targetReps: 10,
                targetWeight: 60,
                targetDuration: null,
                targetDistance: null,
                targetTime: null,
                restAfterSeconds: null,
              },
            ],
          },
        ],
        pendingGroups: overrides.withPendingGroup
          ? [
              {
                groupId: 5,
                groupName: 'Push Variants',
                members: [
                  { id: 100, name: 'Incline Press' },
                  { id: 101, name: 'Decline Press' },
                ],
                position: 1,
              },
            ]
          : [],
      },
    ],
  };
}

@Component({
  selector: 'app-test-host-section',
  template: `
    <app-workout-start-section
      [section]="startForm.sections[0]"
      [sectionIndex]="0"
      [exerciseDisplayMap]="displayMap"
      [readonly]="isReadonly"
      (removed)="removedCalled = true"
      (addExerciseRequested)="addExerciseCalled = true"
      (groupExercisePicked)="lastGroupPick = $event"
    />
  `,
  imports: [WorkoutStartSection],
})
class TestHost {
  model = signal(makeModel());
  startForm = form(this.model);
  displayMap: Record<number, ExerciseDisplayInfo> = DISPLAY_MAP;
  isReadonly = false;
  removedCalled = false;
  addExerciseCalled = false;
  lastGroupPick: { groupIndex: number; exerciseId: number } | null = null;

  setModel(m: StartModel) {
    this.model.set(m);
  }
}

describe('WorkoutStartSection', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [TestHost],
      providers: [provideTranslocoForTest()],
    });
  });

  function setup(overrides: { readonly?: boolean; withPendingGroup?: boolean } = {}) {
    const fixture = TestBed.createComponent(TestHost);
    if (overrides.withPendingGroup) {
      fixture.componentInstance.setModel(makeModel({ withPendingGroup: true }));
    }
    if (overrides.readonly != null) fixture.componentInstance.isReadonly = overrides.readonly;
    fixture.detectChanges();
    return fixture;
  }

  describe('editable mode', () => {
    it('shows section drag handle', () => {
      const fixture = setup({ readonly: false });
      expect(fixture.nativeElement.querySelector('[cdkDragHandle]')).toBeTruthy();
    });

    it('shows remove button', () => {
      const fixture = setup({ readonly: false });
      expect(fixture.nativeElement.textContent).toContain('common.remove');
    });

    it('shows type as brn-select', () => {
      const fixture = setup({ readonly: false });
      expect(fixture.nativeElement.querySelector('brn-select')).toBeTruthy();
    });

    it('shows label as input', () => {
      const fixture = setup({ readonly: false });
      // The label input has [formField]="section().label"
      const inputs = fixture.nativeElement.querySelectorAll('input[hlminput]');
      expect(inputs.length).toBeGreaterThan(0);
    });

    it('shows add exercise button', () => {
      const fixture = setup({ readonly: false });
      expect(fixture.nativeElement.textContent).toContain('user.workouts.addExercise');
    });
  });

  describe('readonly mode', () => {
    it('hides section drag handle', () => {
      const fixture = setup({ readonly: true });
      expect(fixture.nativeElement.querySelector('[cdkDragHandle]')).toBeNull();
    });

    it('hides remove button', () => {
      const fixture = setup({ readonly: true });
      // Only buttons in readonly should be none (remove is hidden, add exercise is hidden)
      const buttons = fixture.nativeElement.querySelectorAll('button');
      const removeBtn = Array.from(buttons as NodeListOf<HTMLButtonElement>).find(
        (b) => b.textContent?.trim() === 'common.remove',
      );
      expect(removeBtn).toBeUndefined();
    });

    it('shows type as plain text instead of select', () => {
      const fixture = setup({ readonly: true });
      expect(fixture.nativeElement.querySelector('brn-select')).toBeNull();
      expect(fixture.nativeElement.textContent).toContain('enums.workoutSectionType.main');
    });

    it('shows label as plain text instead of input', () => {
      const fixture = setup({ readonly: true });
      expect(fixture.nativeElement.textContent).toContain('Chest');
      // Should show em-dash for empty labels, but here we have 'Chest'
    });

    it('hides add exercise button', () => {
      const fixture = setup({ readonly: true });
      expect(fixture.nativeElement.textContent).not.toContain('user.workouts.addExercise');
    });

    it('still renders exercise items', () => {
      const fixture = setup({ readonly: true });
      const exerciseItems = fixture.nativeElement.querySelectorAll(
        'app-workout-start-exercise-item',
      );
      expect(exerciseItems.length).toBe(1);
    });

    it('passes readonly to exercise items', () => {
      const fixture = setup({ readonly: true });
      const exerciseItem = fixture.nativeElement.querySelector('app-workout-start-exercise-item');
      expect(exerciseItem).toBeTruthy();
      // Exercise item should not have drag handles when parent is readonly
      expect(exerciseItem.querySelector('[cdkDragHandle]')).toBeNull();
    });
  });

  describe('pending groups', () => {
    it('shows pending groups in readonly mode', () => {
      const fixture = setup({ readonly: true, withPendingGroup: true });
      const el = fixture.nativeElement as HTMLElement;

      expect(el.querySelector('[data-testid="pending-group"]')).toBeTruthy();
      expect(el.textContent).toContain('Push Variants');
    });

    it('shows pending groups in editable mode', () => {
      const fixture = setup({ readonly: false, withPendingGroup: true });
      const el = fixture.nativeElement as HTMLElement;

      expect(el.querySelector('[data-testid="pending-group"]')).toBeTruthy();
      expect(el.textContent).toContain('Push Variants');
    });

    it('shows group member options in select', () => {
      const fixture = setup({ withPendingGroup: true });
      const el = fixture.nativeElement as HTMLElement;

      const options = el.querySelectorAll('[data-testid="pending-group"] option');
      // First option is the placeholder, then 2 members
      expect(options.length).toBe(3);
      expect(options[1].textContent).toContain('Incline Press');
      expect(options[2].textContent).toContain('Decline Press');
    });

    it('emits groupExercisePicked when a member is selected', () => {
      const fixture = setup({ withPendingGroup: true });
      const el = fixture.nativeElement as HTMLElement;

      const select = el.querySelector('[data-testid="pending-group"] select') as HTMLSelectElement;
      select.value = '100';
      select.dispatchEvent(new Event('change'));
      fixture.detectChanges();

      expect(fixture.componentInstance.lastGroupPick).toEqual({ groupIndex: 0, exerciseId: 100 });
    });
  });
});
