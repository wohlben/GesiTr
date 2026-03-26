import { HttpInterceptorFn } from '@angular/common/http';
import { inject, isDevMode } from '@angular/core';
import { DevelopmentUserHeaderService } from './development-user-header.service';

export const devUserHeaderInterceptor: HttpInterceptorFn = (req, next) => {
  if (!isDevMode() || !req.url.startsWith('/api')) {
    return next(req);
  }
  const service = inject(DevelopmentUserHeaderService);
  return next(
    req.clone({
      setHeaders: { 'X-User-Id': service.userId$.getValue() },
    }),
  );
};
