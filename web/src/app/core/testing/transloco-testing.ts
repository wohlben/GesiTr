import { Injectable } from '@angular/core';
import { provideTransloco, Translation, TranslocoLoader } from '@jsverse/transloco';
import { of } from 'rxjs';

@Injectable()
class InlineTranslocoLoader implements TranslocoLoader {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  getTranslation(lang: string) {
    return of({} as Translation);
  }
}

export function provideTranslocoForTest() {
  return provideTransloco({
    config: {
      availableLangs: ['en'],
      defaultLang: 'en',
      prodMode: true,
    },
    loader: InlineTranslocoLoader,
  });
}
