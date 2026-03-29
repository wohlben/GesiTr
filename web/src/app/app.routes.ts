import { Routes } from '@angular/router';
import { MainLayout } from './layout/main-layout';

export const routes: Routes = [
  { path: '', redirectTo: '/compendium/exercises', pathMatch: 'full' },
  {
    path: 'settings',
    component: MainLayout,
    children: [
      {
        path: 'profile',
        loadComponent: () =>
          import('$pages/settings/profile/settings-profile').then((m) => m.SettingsProfile),
      },
    ],
  },
  {
    path: 'user',
    component: MainLayout,
    children: [
      {
        path: 'exercises/:id/track',
        loadComponent: () =>
          import('$pages/user/exercise-track/exercise-track').then((m) => m.ExerciseTrack),
      },
      {
        path: 'exercises/:id',
        loadComponent: () =>
          import('$pages/user/user-exercise-detail/user-exercise-detail').then(
            (m) => m.UserExerciseDetail,
          ),
      },
      {
        path: 'equipment/:id',
        loadComponent: () =>
          import('$pages/user/user-equipment-detail/user-equipment-detail').then(
            (m) => m.UserEquipmentDetail,
          ),
      },
      {
        path: 'equipment',
        loadComponent: () =>
          import('$pages/user/user-equipment-list/user-equipment-list').then(
            (m) => m.UserEquipmentList,
          ),
      },
      {
        path: 'workouts/new',
        loadComponent: () =>
          import('$pages/user/workout-edit/workout-edit').then((m) => m.WorkoutEdit),
      },
      {
        path: 'workouts/:id/schedules/new',
        loadComponent: () =>
          import('$pages/user/workout-schedule-edit/workout-schedule-edit').then(
            (m) => m.WorkoutScheduleEdit,
          ),
      },
      {
        path: 'workouts/:id/schedules/:scheduleId/periods/:periodId',
        loadComponent: () =>
          import('$pages/user/workout-schedule-period/workout-schedule-period').then(
            (m) => m.WorkoutSchedulePeriod,
          ),
      },
      {
        path: 'workouts/:id/schedules/:scheduleId/edit',
        loadComponent: () =>
          import('$pages/user/workout-schedule-edit/workout-schedule-edit').then(
            (m) => m.WorkoutScheduleEdit,
          ),
      },
      {
        path: 'workouts/:id/schedules/:scheduleId',
        redirectTo: 'workouts/:id/schedules/:scheduleId/edit',
      },
      {
        path: 'workouts/:id/schedules',
        loadComponent: () =>
          import('$pages/user/workout-schedule-list/workout-schedule-list').then(
            (m) => m.WorkoutScheduleList,
          ),
      },
      {
        path: 'workouts/:id/group',
        loadComponent: () =>
          import('$pages/user/workout-group/workout-group').then((m) => m.WorkoutGroup),
      },
      {
        path: 'workouts/:id/start',
        loadComponent: () =>
          import('$pages/user/workout-start/workout-start').then((m) => m.WorkoutStart),
      },
      {
        path: 'workouts/:id/edit',
        loadComponent: () =>
          import('$pages/user/workout-edit/workout-edit').then((m) => m.WorkoutEdit),
      },
      {
        path: 'workouts/:id',
        redirectTo: 'workouts/:id/edit',
      },
      {
        path: 'workout-logs/:id',
        loadComponent: () =>
          import('$pages/user/workout-log-detail/workout-log-detail').then(
            (m) => m.WorkoutLogDetail,
          ),
      },
      {
        path: 'workouts',
        loadComponent: () =>
          import('$pages/user/workout-list/workout-list').then((m) => m.WorkoutList),
      },
      {
        path: 'calendar',
        loadComponent: () => import('$pages/user/calendar/calendar').then((m) => m.Calendar),
      },
    ],
  },
  {
    path: 'compendium',
    component: MainLayout,
    children: [
      {
        path: 'exercises/new',
        loadComponent: () =>
          import('$pages/compendium/exercise-edit/exercise-edit').then((m) => m.ExerciseEdit),
      },
      {
        path: 'exercises/:id/:slug/history',
        loadComponent: () =>
          import('$pages/compendium/exercise-history/exercise-history').then(
            (m) => m.ExerciseHistory,
          ),
      },
      {
        path: 'exercises/:id/:slug/edit',
        loadComponent: () =>
          import('$pages/compendium/exercise-edit/exercise-edit').then((m) => m.ExerciseEdit),
      },
      {
        path: 'exercises/:id/:slug',
        loadComponent: () =>
          import('$pages/compendium/exercise-detail/exercise-detail').then((m) => m.ExerciseDetail),
      },
      {
        path: 'exercises',
        loadComponent: () =>
          import('$pages/compendium/exercise-list/exercise-list').then((m) => m.ExerciseList),
      },
      {
        path: 'equipment/new',
        loadComponent: () =>
          import('$pages/compendium/equipment-edit/equipment-edit').then((m) => m.EquipmentEdit),
      },
      {
        path: 'equipment/:id/:slug/history',
        loadComponent: () =>
          import('$pages/compendium/equipment-history/equipment-history').then(
            (m) => m.EquipmentHistory,
          ),
      },
      {
        path: 'equipment/:id/:slug/edit',
        loadComponent: () =>
          import('$pages/compendium/equipment-edit/equipment-edit').then((m) => m.EquipmentEdit),
      },
      {
        path: 'equipment/:id/:slug',
        loadComponent: () =>
          import('$pages/compendium/equipment-detail/equipment-detail').then(
            (m) => m.EquipmentDetail,
          ),
      },
      {
        path: 'equipment',
        loadComponent: () =>
          import('$pages/compendium/equipment-list/equipment-list').then((m) => m.EquipmentList),
      },
      {
        path: 'exercise-groups/new',
        loadComponent: () =>
          import('$pages/compendium/exercise-group-edit/exercise-group-edit').then(
            (m) => m.ExerciseGroupEdit,
          ),
      },
      {
        path: 'exercise-groups/:id/:slug/edit',
        loadComponent: () =>
          import('$pages/compendium/exercise-group-edit/exercise-group-edit').then(
            (m) => m.ExerciseGroupEdit,
          ),
      },
      {
        path: 'exercise-groups/:id/:slug',
        loadComponent: () =>
          import('$pages/compendium/exercise-group-detail/exercise-group-detail').then(
            (m) => m.ExerciseGroupDetail,
          ),
      },
      {
        path: 'exercise-groups',
        loadComponent: () =>
          import('$pages/compendium/exercise-group-list/exercise-group-list').then(
            (m) => m.ExerciseGroupList,
          ),
      },
    ],
  },
];
