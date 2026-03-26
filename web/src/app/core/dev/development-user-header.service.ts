import { Injectable, isDevMode } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

const STORAGE_KEY = 'dev-x-user-id';

@Injectable()
export class DevelopmentUserHeaderService {
  readonly userId$ = new BehaviorSubject<string>(sessionStorage.getItem(STORAGE_KEY) ?? 'devuser');

  constructor() {
    if (!isDevMode()) {
      throw new Error('DevelopmentUserHeaderService must not be instantiated in production');
    }
    this.userId$.subscribe((id) => sessionStorage.setItem(STORAGE_KEY, id));
  }
}
