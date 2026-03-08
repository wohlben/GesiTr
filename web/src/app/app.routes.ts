import { Routes } from '@angular/router';
import { MainLayout } from './layout/main-layout';

export const routes: Routes = [
  { path: '', redirectTo: '/compendium/exercises', pathMatch: 'full' },
  {
    path: 'compendium',
    component: MainLayout,
    children: [
      {
        path: 'exercises/:id/:slug/edit',
        loadComponent: () =>
          import('$features/compendium/exercise-edit/exercise-edit').then(
            (m) => m.ExerciseEdit,
          ),
      },
      {
        path: 'exercises/:id/:slug',
        loadComponent: () =>
          import('$features/compendium/exercise-detail/exercise-detail').then(
            (m) => m.ExerciseDetail,
          ),
      },
      {
        path: 'exercises',
        loadComponent: () =>
          import('$features/compendium/exercise-list/exercise-list').then(
            (m) => m.ExerciseList,
          ),
      },
      {
        path: 'equipment/:id/:slug/edit',
        loadComponent: () =>
          import('$features/compendium/equipment-edit/equipment-edit').then(
            (m) => m.EquipmentEdit,
          ),
      },
      {
        path: 'equipment/:id/:slug',
        loadComponent: () =>
          import('$features/compendium/equipment-detail/equipment-detail').then(
            (m) => m.EquipmentDetail,
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
        path: 'exercise-groups/:id/:slug/edit',
        loadComponent: () =>
          import('$features/compendium/exercise-group-edit/exercise-group-edit').then(
            (m) => m.ExerciseGroupEdit,
          ),
      },
      {
        path: 'exercise-groups/:id/:slug',
        loadComponent: () =>
          import('$features/compendium/exercise-group-detail/exercise-group-detail').then(
            (m) => m.ExerciseGroupDetail,
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
