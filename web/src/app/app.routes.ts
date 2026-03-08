import { Routes } from '@angular/router';
import { MainLayout } from './layout/main-layout';

export const routes: Routes = [
  { path: '', redirectTo: '/compendium/exercises', pathMatch: 'full' },
  {
    path: 'compendium',
    component: MainLayout,
    children: [
      {
        path: 'exercises',
        loadComponent: () =>
          import('$features/compendium/exercise-list/exercise-list').then(
            (m) => m.ExerciseList,
          ),
      },
      {
        path: 'equipment',
        loadComponent: () =>
          import('$features/compendium/equipment-list/equipment-list').then(
            (m) => m.EquipmentList,
          ),
      },
      {
        path: 'exercise-groups',
        loadComponent: () =>
          import('$features/compendium/exercise-group-list/exercise-group-list').then(
            (m) => m.ExerciseGroupList,
          ),
      },
    ],
  },
];
