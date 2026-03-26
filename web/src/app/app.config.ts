import { ApplicationConfig, provideBrowserGlobalErrorListeners, isDevMode } from '@angular/core';
import { provideHttpClient, withInterceptors } from '@angular/common/http';
import { provideRouter } from '@angular/router';
import { provideTanStackQuery, QueryClient } from '@tanstack/angular-query-experimental';
import { withDevtools } from '@tanstack/angular-query-experimental/devtools';
import { provideTransloco } from '@jsverse/transloco';
import { routes } from './app.routes';
import { provideServiceWorker } from '@angular/service-worker';
import { TranslocoHttpLoader } from '$core/transloco-loader';
import { devUserHeaderInterceptor } from '$core/dev/dev-user-header.interceptor';
import { DevelopmentUserHeaderService } from '$core/dev/development-user-header.service';

export const appConfig: ApplicationConfig = {
  providers: [
    provideBrowserGlobalErrorListeners(),
    provideRouter(routes),
    provideHttpClient(...(isDevMode() ? [withInterceptors([devUserHeaderInterceptor])] : [])),
    ...(isDevMode() ? [DevelopmentUserHeaderService] : []),
    provideTanStackQuery(new QueryClient(), withDevtools()),
    provideTransloco({
      config: {
        availableLangs: ['en'],
        defaultLang: 'en',
        reRenderOnLangChange: true,
        prodMode: !isDevMode(),
      },
      loader: TranslocoHttpLoader,
    }),
    provideServiceWorker('ngsw-worker.js', {
      enabled: !isDevMode(),
      registrationStrategy: 'registerWhenStable:30000',
    }),
  ],
};
