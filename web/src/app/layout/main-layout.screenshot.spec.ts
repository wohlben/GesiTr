import { render } from '@testing-library/angular';
import { page } from 'vitest/browser';
import { Component } from '@angular/core';
import { provideRouter } from '@angular/router';
import { MainLayout } from './main-layout';

@Component({
  selector: 'app-test-page',
  template: '<p class="text-gray-700 dark:text-gray-300">Page content goes here</p>',
})
class TestPage {}

describe('MainLayout screenshots', () => {
  afterEach(() => {
    document.documentElement.classList.remove('dark');
  });

  const providers = [
    provideRouter([
      { path: '**', component: TestPage },
    ]),
  ];

  it('light', async () => {
    const { fixture } = await render(MainLayout, { providers });
    await fixture.whenStable();
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('light');
  });

  it('dark', async () => {
    document.documentElement.classList.add('dark');
    const { fixture } = await render(MainLayout, { providers });
    await fixture.whenStable();
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('dark');
  });
});
