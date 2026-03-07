import { Routes } from '@angular/router';
import { MainLayout } from './layout/main-layout';
import { ExerciseList } from '$features/compendium/exercise-list/exercise-list';
import { EquipmentList } from '$features/compendium/equipment-list/equipment-list';
import { ExerciseGroupList } from '$features/compendium/exercise-group-list/exercise-group-list';

export const routes: Routes = [
  { path: '', redirectTo: '/compendium/exercises', pathMatch: 'full' },
  {
    path: 'compendium',
    component: MainLayout,
    children: [
      { path: 'exercises', component: ExerciseList },
      { path: 'equipment', component: EquipmentList },
      { path: 'exercise-groups', component: ExerciseGroupList },
    ],
  },
];
