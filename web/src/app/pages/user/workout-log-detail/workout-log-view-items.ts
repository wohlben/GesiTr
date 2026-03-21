import { Pipe, PipeTransform } from '@angular/core';
import {
  WorkoutLogExercise,
  WorkoutLogExerciseSet,
  WorkoutLogSection,
} from '$generated/user-models';

export interface ViewItemHeader {
  type: 'header';
  id: string;
  exerciseName: string;
}

export interface ViewItemSet {
  type: 'set';
  id: string;
  set: WorkoutLogExerciseSet;
  exercise: WorkoutLogExercise;
  section: WorkoutLogSection;
  exerciseName: string;
  role: 'completed' | 'active' | 'upcoming';
  setCount: number;
  isNaturalNext?: boolean;
  isOverride?: boolean;
}

export interface ViewItemBreak {
  type: 'break';
  id: string;
  seconds: number;
  label: string;
  role: 'elapsed' | 'active-timer' | 'upcoming';
}

export type ViewItem = ViewItemHeader | ViewItemSet | ViewItemBreak;

@Pipe({ name: 'asHeader' })
export class AsHeaderPipe implements PipeTransform {
  transform(item: ViewItem): ViewItemHeader | undefined {
    return item.type === 'header' ? item : undefined;
  }
}

@Pipe({ name: 'asSet' })
export class AsSetPipe implements PipeTransform {
  transform(item: ViewItem): ViewItemSet | undefined {
    return item.type === 'set' ? item : undefined;
  }
}

@Pipe({ name: 'asBreak' })
export class AsBreakPipe implements PipeTransform {
  transform(item: ViewItem): ViewItemBreak | undefined {
    return item.type === 'break' ? item : undefined;
  }
}
