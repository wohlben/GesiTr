import { Injectable, isDevMode } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

@Injectable()
export class DevelopmentUserHeaderService {
  readonly userId$ = new BehaviorSubject<string>('devuser');

  constructor() {
    if (!isDevMode()) {
      throw new Error('DevelopmentUserHeaderService must not be instantiated in production');
    }
  }
}
